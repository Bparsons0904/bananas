package main

import (
	"bananas/internal/app"
	"bananas/internal/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func gracefulShutdown(
	app *app.App,
	echoServer *echo.Echo,
	done chan bool,
	log logger.Logger,
) {
	log = log.Function("gracefulShutdown")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := echoServer.Shutdown(ctx); err != nil {
		log.Er("Server forced to shutdown", err)
	}

	if err := app.Close(); err != nil {
		log.Er("failed to close app", err)
	}

	log.Info("Server exiting")
	done <- true
}

func main() {
	log := logger.New("main")

	app, err := app.New()
	if err != nil {
		log.Er("failed to initialize app", err)
		os.Exit(1)
	}
	defer func() {
		if err := app.Close(); err != nil {
			log.Er("failed to close app", err)
		}
	}()

	e := echo.New()
	e.HideBanner = true

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Set framework context middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "framework", "echo")))
			return next(c)
		}
	})

	// Register routes
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	api := e.Group("/api")
	{
		test := api.Group("/test")
		{
			test.GET("/simple", func(c echo.Context) error {
				req := c.Request()
				app.Controllers.SimpleRequest(c.Response(), req)
				return nil
			})
			
			test.GET("/database", func(c echo.Context) error {
				req := c.Request()
				app.Controllers.DatabaseQuery(c.Response(), req)
				return nil
			})
			
			test.GET("/json", func(c echo.Context) error {
				req := c.Request()
				app.Controllers.JsonResponse(c.Response(), req)
				return nil
			})
		}
		
		api.GET("/info", func(c echo.Context) error {
			req := c.Request()
			app.Controllers.FrameworkInfo(c.Response(), req)
			return nil
		})
	}

	done := make(chan bool, 1)

	go func() {
		log.Info("Starting Echo server on port %s", app.Config.ServerPort)
		if err := e.Start(":" + app.Config.ServerPort); err != nil && err != http.ErrServerClosed {
			log.Er("Server failed to start", err)
			os.Exit(1)
		}
	}()

	go gracefulShutdown(app, e, done, log)

	<-done
	log.Info("Graceful shutdown complete.")
}