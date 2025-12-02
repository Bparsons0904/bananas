package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bananas/internal/logger"
	"bananas/internal/services"
	"bananas/internal/templates"
)

type TemplController struct {
	Service *services.Service
	Logger  logger.Logger
}

func NewTemplController(service *services.Service, log logger.Logger) *TemplController {
	return &TemplController{
		Service: service,
		Logger:  log,
	}
}

// HomePage renders the main testing interface
func (c *TemplController) HomePage(w http.ResponseWriter, r *http.Request) {
	component := templates.Home()
	err := component.Render(r.Context(), w)
	if err != nil {
		c.Logger.Er("failed to render home template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RunTest handles HTMX test requests and returns HTML results
func (c *TemplController) RunTest(w http.ResponseWriter, r *http.Request) {
	// Get framework from context (set by framework middleware)
	frameworkName := "unknown"
	if fw := r.Context().Value("framework"); fw != nil {
		frameworkName = fw.(string)
	}

	// Get parameters from query string
	frameworkPort := r.URL.Query().Get("framework")
	orm := r.URL.Query().Get("orm")
	endpoint := r.URL.Query().Get("endpoint")

	if frameworkPort == "" || endpoint == "" {
		c.renderError(w, r, "Missing required parameters", frameworkName, orm)
		return
	}

	// Build the target URL
	targetURL := fmt.Sprintf("http://localhost:%s%s", frameworkPort, endpoint)

	// Add ORM parameter if endpoint is database test
	if endpoint == "/api/test/database?limit=10" && orm != "" {
		targetURL += "&orm=" + orm
	}

	c.Logger.Info("running test", "framework", frameworkName, "target", targetURL, "orm", orm)

	// Make the request and measure duration
	startTime := time.Now()
	resp, err := http.Get(targetURL)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		c.Logger.Er("request failed", err)
		c.renderError(w, r, err.Error(), frameworkName, orm)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.Er("failed to read response", err)
		c.renderError(w, r, err.Error(), frameworkName, orm)
		return
	}

	// Parse JSON response
	var responseData interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		// If not JSON, treat as plain text
		responseData = string(body)
	}

	// Get framework name from port
	frameworkDisplayName := getFrameworkName(frameworkPort)

	// Render results template
	result := templates.TestResult{
		Framework: frameworkDisplayName,
		ORM:       getORMName(orm),
		Response:  responseData,
		Duration:  float64(duration),
		Error:     "",
	}

	component := templates.Results(result)
	err = component.Render(r.Context(), w)
	if err != nil {
		c.Logger.Er("failed to render results template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *TemplController) renderError(w http.ResponseWriter, r *http.Request, errorMsg, framework, orm string) {
	result := templates.TestResult{
		Framework: framework,
		ORM:       getORMName(orm),
		Response:  nil,
		Duration:  0,
		Error:     errorMsg,
	}

	component := templates.Results(result)
	err := component.Render(r.Context(), w)
	if err != nil {
		c.Logger.Er("failed to render error template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getFrameworkName(port string) string {
	frameworks := map[string]string{
		"8081": "Standard Library",
		"8082": "Gin",
		"8083": "Fiber",
		"8084": "Echo",
		"8085": "Chi",
		"8086": "Gorilla Mux",
	}
	if name, ok := frameworks[port]; ok {
		return name
	}
	return "Unknown"
}

func getORMName(value string) string {
	orms := map[string]string{
		"sql":  "database/sql",
		"gorm": "GORM",
		"sqlx": "SQLx",
		"pgx":  "PGX",
	}
	if name, ok := orms[value]; ok {
		return name
	}
	return value
}
