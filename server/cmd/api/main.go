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

	// Create a router for the standard library
	mux := http.NewServeMux()
	
	// Add middleware to set framework context
	frameworkMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "framework", "standard")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	// Register routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/test/simple", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "framework", "standard")
		app.Controllers.SimpleRequest(w, r.WithContext(ctx))
	})

	mux.HandleFunc("/api/test/database", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "framework", "standard")
		app.Controllers.DatabaseQuery(w, r.WithContext(ctx))
	})

	mux.HandleFunc("/api/test/json", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "framework", "standard")
		app.Controllers.JsonResponse(w, r.WithContext(ctx))
	})

	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "framework", "standard")
		app.Controllers.FrameworkInfo(w, r.WithContext(ctx))
	})

	server := &http.Server{
		Addr:    ":" + app.Config.ServerPort,
		Handler: frameworkMiddleware(mux),
	}

	done := make(chan bool, 1)

	go func() {
		log.Info("Starting standard library server on port %s", app.Config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Er("Server failed to start", err)
			os.Exit(1)
		}
	}()

	go gracefulShutdown(app, server, done, log)

	<-done
	log.Info("Graceful shutdown complete.")
}