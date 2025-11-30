package main

import (
	"bananas/internal/app"
	"bananas/internal/logger"
	appLogger "bananas/internal/logger"
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func gracefulShutdown(
	app *app.App,
	fiberApp *fiber.App,
	done chan bool,
	log appLogger.Logger,
) {
	log = log.Function("gracefulShutdown")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
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

	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Add middleware
	fiberApp.Use(recover.New())
	
	// Set framework context middleware
	fiberApp.Use(func(c *fiber.Ctx) error {
		c.Context().SetUserValue("framework", "fiber")
		return c.Next()
	})

	// Register routes
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	api := fiberApp.Group("/api")
	{
		test := api.Group("/test")
		{
			test.Get("/simple", func(c *fiber.Ctx) error {
				// Set framework context
				ctx := context.WithValue(context.Background(), "framework", "fiber")
				
				// Create a standard HTTP request from Fiber context
				req := &http.Request{
					Method: c.Method(),
					Header: make(http.Header),
					URL: &url.URL{},
				}
				
				// Copy query parameters
				c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
					req.URL.Query().Add(string(key), string(value))
				})
				
				// Use a simple response writer that writes directly to Fiber
				writer := &fiberResponseWriter{ctx: c}
				
				app.Controllers.SimpleRequest(writer, req.WithContext(ctx))
				return nil
			})
			
			test.Get("/database", func(c *fiber.Ctx) error {
				ctx := context.WithValue(context.Background(), "framework", "fiber")
				
				limitStr := c.Query("limit")
				url := &url.URL{}
				if limitStr != "" {
					url.RawQuery = "limit=" + limitStr
				}
				
				req := &http.Request{
					Method: "GET",
					URL: url,
					Header: make(http.Header),
				}
				
				writer := &fiberResponseWriter{ctx: c}
				app.Controllers.DatabaseQuery(writer, req.WithContext(ctx))
				
				return nil
			})
			
			test.Get("/json", func(c *fiber.Ctx) error {
				ctx := context.WithValue(context.Background(), "framework", "fiber")
				
				req := &http.Request{
					Method: "GET",
					URL: &url.URL{},
					Header: make(http.Header),
				}
				
				writer := &fiberResponseWriter{ctx: c}
				app.Controllers.JsonResponse(writer, req.WithContext(ctx))
				return nil
			})
		}
		
		api.Get("/info", func(c *fiber.Ctx) error {
			ctx := context.WithValue(context.Background(), "framework", "fiber")
			
			req := &http.Request{
				Method: "GET",
				URL: &url.URL{},
				Header: make(http.Header),
			}
			
			writer := &fiberResponseWriter{ctx: c}
			app.Controllers.FrameworkInfo(writer, req.WithContext(ctx))
			return nil
		})
	}

	done := make(chan bool, 1)

	go func() {
		log.Info("Starting Fiber server on port %s", app.Config.ServerPort)
		if err := fiberApp.Listen(":" + app.Config.ServerPort); err != nil {
			log.Er("Server failed to start", err)
			os.Exit(1)
		}
	}()

	go gracefulShutdown(app, fiberApp, done, log)

	<-done
	log.Info("Graceful shutdown complete.")
}

// fiberResponseWriter adapts Fiber's Ctx to http.ResponseWriter
type fiberResponseWriter struct {
	ctx *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	header := make(http.Header)
	// Access response headers through Response()
	w.ctx.Response().Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (w *fiberResponseWriter) Write(data []byte) (int, error) {
	err := w.ctx.Send(data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
}