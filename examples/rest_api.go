//go:build examples
// +build examples

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coredds/godateparser"
)

// REST API Example demonstrates using godateparser in a web service

type ParseRequest struct {
	DateString   string   `json:"date_string"`
	DateOrder    string   `json:"date_order,omitempty"`
	Languages    []string `json:"languages,omitempty"`
	PreferFuture bool     `json:"prefer_future,omitempty"`
}

type ParseResponse struct {
	Success   bool      `json:"success"`
	Input     string    `json:"input"`
	Parsed    time.Time `json:"parsed,omitempty"`
	Formatted string    `json:"formatted,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
	Error     string    `json:"error,omitempty"`
}

type ExtractRequest struct {
	Text      string   `json:"text"`
	DateOrder string   `json:"date_order,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

type ExtractResponse struct {
	Success bool                       `json:"success"`
	Count   int                        `json:"count"`
	Dates   []ExtractedDate            `json:"dates,omitempty"`
	Error   string                     `json:"error,omitempty"`
}

type ExtractedDate struct {
	Text       string    `json:"text"`
	Position   int       `json:"position"`
	Parsed     time.Time `json:"parsed"`
	Formatted  string    `json:"formatted"`
	Confidence float64   `json:"confidence"`
}

func main() {
	fmt.Println("=== REST API Example ===")
	fmt.Println()

	// Setup routes
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/parse", handleParse)
	http.HandleFunc("/extract", handleExtract)
	http.HandleFunc("/health", handleHealth)

	port := ":8080"
	fmt.Printf("Starting date parser API server on http://localhost%s\n\n", port)
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /          - API documentation")
	fmt.Println("  POST /parse     - Parse a single date string")
	fmt.Println("  POST /extract   - Extract dates from text")
	fmt.Println("  GET  /health    - Health check")
	fmt.Println()
	fmt.Println("Example requests:")
	fmt.Println(`
  curl -X POST http://localhost:8080/parse \
    -H "Content-Type: application/json" \
    -d '{"date_string":"December 31, 2024"}'

  curl -X POST http://localhost:8080/extract \
    -H "Content-Type: application/json" \
    -d '{"text":"Meeting on Dec 31 and follow-up in 2 weeks"}'
	`)

	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	docs := map[string]interface{}{
		"name":    "Date Parser API",
		"version": godateparser.Version,
		"endpoints": map[string]interface{}{
			"/parse": map[string]string{
				"method":      "POST",
				"description": "Parse a single date string",
				"example": `{
  "date_string": "December 31, 2024",
  "date_order": "MDY",
  "languages": ["en"],
  "prefer_future": true
}`,
			},
			"/extract": map[string]string{
				"method":      "POST",
				"description": "Extract all dates from text",
				"example": `{
  "text": "Meeting on Dec 31 and follow-up in 2 weeks",
  "date_order": "MDY",
  "languages": ["en", "es"]
}`,
			},
			"/health": map[string]string{
				"method":      "GET",
				"description": "Health check endpoint",
			},
		},
	}

	json.NewEncoder(w).Encode(docs)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type": "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ParseResponse{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ParseResponse{
			Success: false,
			Error:   "Invalid JSON: " + err.Error(),
		})
		return
	}

	if req.DateString == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ParseResponse{
			Success: false,
			Error:   "date_string is required",
		})
		return
	}

	// Build settings
	settings := &godateparser.Settings{
		RelativeBase: time.Now(),
	}

	if req.DateOrder != "" {
		settings.DateOrder = req.DateOrder
	}

	if len(req.Languages) > 0 {
		settings.Languages = req.Languages
	}

	if req.PreferFuture {
		settings.PreferDatesFrom = "future"
	}

	// Parse the date
	parsed, err := godateparser.ParseDate(req.DateString, settings)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ParseResponse{
			Success: false,
			Input:   req.DateString,
			Error:   err.Error(),
		})
		return
	}

	// Return success
	json.NewEncoder(w).Encode(ParseResponse{
		Success:   true,
		Input:     req.DateString,
		Parsed:    parsed,
		Formatted: parsed.Format("Monday, January 2, 2006 15:04:05 MST"),
		Timestamp: parsed.Unix(),
	})
}

func handleExtract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type: application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ExtractResponse{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return
	}

	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExtractResponse{
			Success: false,
			Error:   "Invalid JSON: " + err.Error(),
		})
		return
	}

	if req.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExtractResponse{
			Success: false,
			Error:   "text is required",
		})
		return
	}

	// Build settings
	settings := &godateparser.Settings{
		RelativeBase: time.Now(),
	}

	if req.DateOrder != "" {
		settings.DateOrder = req.DateOrder
	}

	if len(req.Languages) > 0 {
		settings.Languages = req.Languages
	}

	// Extract dates
	dates, err := godateparser.ExtractDates(req.Text, settings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ExtractResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var extracted []ExtractedDate
	for _, d := range dates {
		extracted = append(extracted, ExtractedDate{
			Text:       d.MatchedText,
			Position:   d.Position,
			Parsed:     d.Date,
			Formatted:  d.Date.Format("2006-01-02 15:04:05"),
			Confidence: d.Confidence,
		})
	}

	json.NewEncoder(w).Encode(ExtractResponse{
		Success: true,
		Count:   len(extracted),
		Dates:   extracted,
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   godateparser.Version,
	}

	json.NewEncoder(w).Encode(health)
}

/*
Example usage with curl:

# Parse a simple date
curl -X POST http://localhost:8080/parse \
  -H "Content-Type: application/json" \
  -d '{"date_string":"December 31, 2024"}'

# Parse with custom settings
curl -X POST http://localhost:8080/parse \
  -H "Content-Type: application/json" \
  -d '{"date_string":"31/12/2024","date_order":"DMY","languages":["en"]}'

# Parse relative date
curl -X POST http://localhost:8080/parse \
  -H "Content-Type: application/json" \
  -d '{"date_string":"next Monday","prefer_future":true}'

# Extract dates from text
curl -X POST http://localhost:8080/extract \
  -H "Content-Type: application/json" \
  -d '{"text":"Meeting on December 31, 2024 and follow-up on 2025-01-15"}'

# Extract with multiple languages
curl -X POST http://localhost:8080/extract \
  -H "Content-Type: application/json" \
  -d '{"text":"Event: 31 diciembre 2024 or December 31","languages":["es","en"]}'

# Health check
curl http://localhost:8080/health

# API documentation
curl http://localhost:8080/

Example responses:

Parse success:
{
  "success": true,
  "input": "December 31, 2024",
  "parsed": "2024-12-31T00:00:00Z",
  "formatted": "Tuesday, December 31, 2024 00:00:00 UTC",
  "timestamp": 1735689600
}

Extract success:
{
  "success": true,
  "count": 2,
  "dates": [
    {
      "text": "December 31, 2024",
      "position": 11,
      "parsed": "2024-12-31T00:00:00Z",
      "formatted": "2024-12-31 00:00:00",
      "confidence": 1.0
    },
    {
      "text": "2025-01-15",
      "position": 47,
      "parsed": "2025-01-15T00:00:00Z",
      "formatted": "2025-01-15 00:00:00",
      "confidence": 1.0
    }
  ]
}
*/

