package gcal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server is an HTTP server for calendar operations
type Server struct {
	config *Config
	client *APIClient
	port   int
}

// CreateRequest is the request body for /create endpoint
type CreateRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Location    string   `json:"location,omitempty"`
	Start       string   `json:"start"` // RFC3339 format
	End         string   `json:"end"`   // RFC3339 format
	Attendees   []string `json:"attendees,omitempty"`
}

// ListRequest is the request body for /list endpoint
type ListRequest struct {
	Start      string `json:"start,omitempty"`      // RFC3339 format, defaults to now
	End        string `json:"end,omitempty"`        // RFC3339 format, defaults to end of day
	MaxResults int    `json:"max_results,omitempty"` // defaults to 10
}

// NewServer creates a new calendar server
func NewServer(config *Config, port int) (*Server, error) {
	client, err := NewAPIClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &Server{
		config: config,
		client: client,
		port:   port,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", s.handleInfo)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/create", s.handleCreate)
	mux.HandleFunc("/list", s.handleList)
	mux.HandleFunc("/today", s.handleToday)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Calendar server starting on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  POST /create  - Create event")
	log.Printf("  GET  /list    - List events")
	log.Printf("  GET  /today   - List today's events")
	log.Printf("  GET  /health  - Health check")

	return http.ListenAndServe(addr, mux)
}

// handleInfo shows server info
func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	info := map[string]interface{}{
		"name":       "calendar-server",
		"version":    "1.0.0",
		"calendar":   s.config.CalendarID,
		"endpoints": map[string]string{
			"/":       "Info (this page)",
			"/health": "Health check",
			"/create": "POST - Create event",
			"/list":   "GET/POST - List events",
			"/today":  "GET - List today's events",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleCreate creates a new calendar event
func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Parse times
	start, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid start time: %v", err), http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, req.End)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid end time: %v", err), http.StatusBadRequest)
		return
	}

	event := &Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Start:       start,
		End:         end,
		Attendees:   req.Attendees,
	}

	result, err := s.client.Create(event)
	if err != nil {
		log.Printf("Create failed: %v", err)
	} else {
		log.Printf("Event created: %s - %s", result.EventID, event.Title)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleList lists calendar events
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	var req ListRequest

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Default time range: now to end of day
	now := time.Now()
	start := now
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Override with request values
	if req.Start != "" {
		var err error
		start, err = time.Parse(time.RFC3339, req.Start)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid start time: %v", err), http.StatusBadRequest)
			return
		}
	}
	if req.End != "" {
		var err error
		end, err = time.Parse(time.RFC3339, req.End)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid end time: %v", err), http.StatusBadRequest)
			return
		}
	}

	maxResults := req.MaxResults
	if maxResults <= 0 {
		maxResults = 10
	}

	result, err := s.client.List(start, end, maxResults)
	if err != nil {
		log.Printf("List failed: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleToday lists today's events
func (s *Server) handleToday(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	result, err := s.client.List(start, end, 20)
	if err != nil {
		log.Printf("List today failed: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
