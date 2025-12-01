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

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Set framework context middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "framework", "gorilla")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Register routes
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		app.Controllers.FrameworkInfo(w, r)
	}).Methods("GET")

	test := api.PathPrefix("/test").Subrouter()
	test.HandleFunc("/simple", func(w http.ResponseWriter, r *http.Request) {
		app.Controllers.SimpleRequest(w, r)
	}).Methods("GET")

	test.HandleFunc("/database", func(w http.ResponseWriter, r *http.Request) {
		app.Controllers.DatabaseQuery(w, r)
	}).Methods("GET")

	test.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		app.Controllers.JsonResponse(w, r)
	}).Methods("GET")

	server := &http.Server{
		Addr:    ":" + app.Config.ServerPort,
		Handler: r,
	}

	done := make(chan bool, 1)

	go func() {
		log.Info("Starting Gorilla Mux server on port %s", app.Config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Er("Server failed to start", err)
			os.Exit(1)
		}
	}()

	go gracefulShutdown(app, server, done, log)

	<-done
	log.Info("Graceful shutdown complete.")
}