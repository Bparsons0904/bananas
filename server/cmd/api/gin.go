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

	"github.com/gin-gonic/gin"
)

func gracefulShutdown(
	app *app.App,
	server *http.Server,
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
	
	if err := server.Shutdown(ctx); err != nil {
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

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	
	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// Set framework context middleware
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "framework", "gin"))
		c.Next()
	})

	// Register routes
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := r.Group("/api")
	{
		test := api.Group("/test")
		{
			test.GET("/simple", func(c *gin.Context) {
				app.Controllers.SimpleRequest(c.Writer, c.Request)
			})
			
			test.GET("/database", func(c *gin.Context) {
				app.Controllers.DatabaseQuery(c.Writer, c.Request)
			})
			
			test.GET("/json", func(c *gin.Context) {
				app.Controllers.JsonResponse(c.Writer, c.Request)
			})
		}
		
		api.GET("/info", func(c *gin.Context) {
			app.Controllers.FrameworkInfo(c.Writer, c.Request)
		})
	}

	server := &http.Server{
		Addr:    ":" + app.Config.ServerPort,
		Handler: r,
	}

	done := make(chan bool, 1)

	go func() {
		log.Info("Starting Gin server on port %s", app.Config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Er("Server failed to start", err)
			os.Exit(1)
		}
	}()

	go gracefulShutdown(app, server, done, log)

	<-done
	log.Info("Graceful shutdown complete.")
}