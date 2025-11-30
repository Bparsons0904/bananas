package main

import (
	"bananas/internal/app"
	"bananas/internal/logger"
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func gracefulShutdown(
	servers []*http.Server,
	app *app.App,
	done chan bool,
	shutdownChan chan struct{},
	log logger.Logger,
	wg *sync.WaitGroup,
) {
	log = log.Function("gracefulShutdown")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutting down gracefully, press Ctrl+C again to force")

	close(shutdownChan)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wgShutdown sync.WaitGroup
	for i, server := range servers {
		wgShutdown.Add(1)
		go func(srv *http.Server, idx int) {
			defer wgShutdown.Done()
			if err := srv.Shutdown(ctx); err != nil {
				log.Er("Server %d forced to shutdown", err, idx)
			} else {
				log.Info("Server %d shutdown gracefully", idx)
			}
		}(server, i)
	}
	wgShutdown.Wait()

	if err := app.Close(); err != nil {
		log.Er("failed to close app", err)
	}

	log.Info("All servers stopped")
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

	var servers []*http.Server
	var wg sync.WaitGroup
	shutdownChan := make(chan struct{})

	// 1. Standard Library Server (Port 8081)
	{
		mux := http.NewServeMux()
		
		frameworkMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "framework", "standard")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}

		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK - Standard Library"))
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
			Addr:    ":8081",
			Handler: frameworkMiddleware(mux),
		}
		servers = append(servers, server)
		
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting Standard Library server on port 8081")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Er("Standard Library server failed to start", err)
			}
		}()
	}

	// 2. Gin Server (Port 8082)
	{
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
		r.Use(func(c *gin.Context) {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "framework", "gin"))
			c.Next()
		})

		r.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "OK - Gin")
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
			Addr:    ":8082",
			Handler: r,
		}
		servers = append(servers, server)
		
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting Gin server on port 8082")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Er("Gin server failed to start", err)
			}
		}()
	}

	// 3. Fiber Server (Port 8083)
	{
		fiberApp := fiber.New(fiber.Config{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})

		fiberApp.Use(recover.New())
		fiberApp.Use(func(c *fiber.Ctx) error {
			c.Context().SetUserValue("framework", "fiber")
			return c.Next()
		})

		fiberApp.Get("/health", func(c *fiber.Ctx) error {
			return c.SendString("OK - Fiber")
		})

		api := fiberApp.Group("/api")
		{
			test := api.Group("/test")
			{
				test.Get("/simple", func(c *fiber.Ctx) error {
					ctx := context.WithValue(context.Background(), "framework", "fiber")
					writer := &fiberResponseWriter{ctx: c}

					parsedURL, _ := url.Parse(c.OriginalURL())
					req := &http.Request{
						Method: c.Method(),
						URL:    parsedURL,
						Header: make(http.Header),
					}
					app.Controllers.SimpleRequest(writer, req.WithContext(ctx))
					return nil
				})
				test.Get("/database", func(c *fiber.Ctx) error {
					ctx := context.WithValue(context.Background(), "framework", "fiber")
					writer := &fiberResponseWriter{ctx: c}

					parsedURL, _ := url.Parse(c.OriginalURL())
					req := &http.Request{
						Method: c.Method(),
						URL:    parsedURL,
						Header: make(http.Header),
					}
					app.Controllers.DatabaseQuery(writer, req.WithContext(ctx))
					return nil
				})
				test.Get("/json", func(c *fiber.Ctx) error {
					ctx := context.WithValue(context.Background(), "framework", "fiber")
					writer := &fiberResponseWriter{ctx: c}

					parsedURL, _ := url.Parse(c.OriginalURL())
					req := &http.Request{
						Method: c.Method(),
						URL:    parsedURL,
						Header: make(http.Header),
					}
					app.Controllers.JsonResponse(writer, req.WithContext(ctx))
					return nil
				})
			}
			api.Get("/info", func(c *fiber.Ctx) error {
				ctx := context.WithValue(context.Background(), "framework", "fiber")
				writer := &fiberResponseWriter{ctx: c}

				parsedURL, _ := url.Parse(c.OriginalURL())
				req := &http.Request{
					Method: c.Method(),
					URL:    parsedURL,
					Header: make(http.Header),
				}
				app.Controllers.FrameworkInfo(writer, req.WithContext(ctx))
				return nil
			})
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting Fiber server on port 8083")
			if err := fiberApp.Listen(":8083"); err != nil {
				log.Er("Fiber server failed to start", err)
			}
		}()

		go func() {
			<-shutdownChan
			log.Info("Shutting down Fiber server...")
			if err := fiberApp.Shutdown(); err != nil {
				log.Er("Fiber server shutdown error", err)
			}
		}()
	}

	// 4. Echo Server (Port 8084)
	{
		e := echo.New()
		e.HideBanner = true
		e.Use(echomiddleware.Logger())
		e.Use(echomiddleware.Recover())
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "framework", "echo")))
				return next(c)
			}
		})

		e.GET("/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK - Echo")
		})

		api := e.Group("/api")
		{
			test := api.Group("/test")
			{
				test.GET("/simple", func(c echo.Context) error {
					app.Controllers.SimpleRequest(c.Response(), c.Request())
					return nil
				})
				test.GET("/database", func(c echo.Context) error {
					app.Controllers.DatabaseQuery(c.Response(), c.Request())
					return nil
				})
				test.GET("/json", func(c echo.Context) error {
					app.Controllers.JsonResponse(c.Response(), c.Request())
					return nil
				})
			}
			api.GET("/info", func(c echo.Context) error {
				app.Controllers.FrameworkInfo(c.Response(), c.Request())
				return nil
			})
		}

		server := &http.Server{
			Addr:    ":8084",
			Handler: e,
		}
		servers = append(servers, server)
		
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting Echo server on port 8084")
			e.Start(":8084")
		}()
	}

	// 5. Chi Server (Port 8085)
	{
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "framework", "chi")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK - Chi"))
		})

		r.Route("/api", func(r chi.Router) {
			r.Route("/test", func(r chi.Router) {
				r.Get("/simple", func(w http.ResponseWriter, r *http.Request) {
					app.Controllers.SimpleRequest(w, r)
				})
				r.Get("/database", func(w http.ResponseWriter, r *http.Request) {
					app.Controllers.DatabaseQuery(w, r)
				})
				r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
					app.Controllers.JsonResponse(w, r)
				})
			})
			r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
				app.Controllers.FrameworkInfo(w, r)
			})
		})

		server := &http.Server{
			Addr:    ":8085",
			Handler: r,
		}
		servers = append(servers, server)
		
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting Chi server on port 8085")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Er("Chi server failed to start", err)
			}
		}()
	}

	// 6. Gorilla Mux Server (Port 8086)
	{
		r := mux.NewRouter()
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "framework", "gorilla")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK - Gorilla Mux"))
		}).Methods("GET")

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
			Addr:    ":8086",
			Handler: r,
		}
		servers = append(servers, server)
		
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting Gorilla Mux server on port 8086")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Er("Gorilla Mux server failed to start", err)
			}
		}()
	}

	// Wait for all servers to be ready
	log.Info("All frameworks are starting up...")
	log.Info("Standard Library: http://localhost:8081")
	log.Info("Gin: http://localhost:8082")
	log.Info("Fiber: http://localhost:8083")
	log.Info("Echo: http://localhost:8084")
	log.Info("Chi: http://localhost:8085")
	log.Info("Gorilla Mux: http://localhost:8086")

	// Health check verification
	time.Sleep(200 * time.Millisecond)
	healthCheckServers(log)

	done := make(chan bool, 1)
	go gracefulShutdown(servers, app, done, shutdownChan, log, &wg)

	<-done
	log.Info("Graceful shutdown complete.")
	wg.Wait()
}

func healthCheckServers(log logger.Logger) {
	log.Info("Performing health checks on all servers...")

	servers := []struct {
		name string
		url  string
	}{
		{"Standard Library", "http://localhost:8081/health"},
		{"Gin", "http://localhost:8082/health"},
		{"Fiber", "http://localhost:8083/health"},
		{"Echo", "http://localhost:8084/health"},
		{"Chi", "http://localhost:8085/health"},
		{"Gorilla Mux", "http://localhost:8086/health"},
	}

	client := &http.Client{Timeout: 2 * time.Second}
	allHealthy := true

	for _, server := range servers {
		resp, err := client.Get(server.url)
		if err != nil {
			log.Er("Health check failed for %s", err, server.name)
			allHealthy = false
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Info("✓ %s is healthy", server.name)
		} else {
			log.Er("Health check failed for %s with status %d", nil, server.name, resp.StatusCode)
			allHealthy = false
		}
	}

	if allHealthy {
		log.Info("All servers are healthy and ready!")
	} else {
		log.Er("Some servers failed health checks", nil)
	}
}

// fiberResponseWriter adapts Fiber's Ctx to http.ResponseWriter
type fiberResponseWriter struct {
	ctx *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	header := make(http.Header)
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