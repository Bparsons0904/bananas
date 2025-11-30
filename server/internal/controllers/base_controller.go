package controllers

import (
	"bananas/internal/logger"
	"bananas/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type BaseController struct {
	Service *services.Service
	Logger  logger.Logger
}

func New(service *services.Service) *BaseController {
	return &BaseController{
		Service: service,
		Logger:  logger.New("controller"),
	}
}

// Common response methods
func (c *BaseController) WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (c *BaseController) WriteError(w http.ResponseWriter, status int, message string) error {
	return c.WriteJSON(w, status, map[string]string{"error": message})
}

// Test endpoints
func (c *BaseController) SimpleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	err := c.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Simple request successful",
		"framework": r.Context().Value("framework").(string),
	})
	
	if err != nil {
		c.Logger.Er("failed to write response", err)
		return
	}
	
	executionMs := time.Since(start).Milliseconds()
	c.logTestResult("simple_request", executionMs, true)
}

func (c *BaseController) DatabaseQuery(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	ormType := r.URL.Query().Get("orm")
	if ormType == "" {
		ormType = "sql"
	}

	results, err := c.Service.GetTestResults(r.Context(), ormType, limit)
	if err != nil {
		c.WriteError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	err = c.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"results":   results,
		"count":     len(results),
		"orm":       ormType,
		"framework": r.Context().Value("framework"),
	})

	if err != nil {
		c.Logger.Er("failed to write response", err)
		return
	}

	executionMs := time.Since(start).Milliseconds()
	c.Logger.Info("Database query completed - ORM: %s, Time: %dms", ormType, executionMs)
	c.logTestResult("database_query", executionMs, true)
}

func (c *BaseController) JsonResponse(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	data := map[string]interface{}{
		"message": "JSON response successful",
		"framework": r.Context().Value("framework").(string),
		"timestamp": time.Now().Unix(),
		"data": []map[string]interface{}{
			{"id": 1, "name": "Item 1", "value": 100.5},
			{"id": 2, "name": "Item 2", "value": 200.3},
			{"id": 3, "name": "Item 3", "value": 150.7},
		},
	}
	
	err := c.WriteJSON(w, http.StatusOK, data)
	if err != nil {
		c.Logger.Er("failed to write response", err)
		return
	}
	
	executionMs := time.Since(start).Milliseconds()
	c.logTestResult("json_response", executionMs, true)
}

func (c *BaseController) FrameworkInfo(w http.ResponseWriter, r *http.Request) {
	framework := r.Context().Value("framework").(string)
	
	err := c.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"framework": framework,
		"type": "backend",
		"endpoints": []string{
			"/api/test/simple",
			"/api/test/database",
			"/api/test/json",
			"/api/info",
		},
	})
	
	if err != nil {
		c.Logger.Er("failed to write response", err)
		return
	}
}

func (c *BaseController) logTestResult(testType string, executionMs int64, success bool) {
	// This would typically be called via service
	c.Logger.Info("Test completed - Type: %s, Time: %dms, Success: %v", 
		testType, executionMs, success)
}