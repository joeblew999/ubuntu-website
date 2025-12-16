package gmail

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server is an HTTP server for email operations
type Server struct {
	config *Config
	sender Sender
	port   int

	// Draft storage (in-memory for now)
	drafts   map[string]*Draft
	draftsMu sync.RWMutex
}

// Draft represents a pending email draft
type Draft struct {
	ID        string    `json:"id"`
	Email     *Email    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// SendRequest is the request body for /send endpoint
type SendRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// NewServer creates a new email server
func NewServer(config *Config, port int) (*Server, error) {
	sender, err := NewAPISender(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create API sender: %w", err)
	}

	return &Server{
		config: config,
		sender: sender,
		port:   port,
		drafts: make(map[string]*Draft),
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", s.handleInfo)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/send", s.handleSend)
	mux.HandleFunc("/compose", s.handleCompose)
	mux.HandleFunc("/drafts", s.handleDrafts)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Gmail server starting on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  POST /send    - Send email")
	log.Printf("  POST /compose - Create draft")
	log.Printf("  GET  /drafts  - List drafts")
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
		"name":    "gmail-server",
		"version": "1.0.0",
		"from":    s.config.FromAddress,
		"endpoints": map[string]string{
			"/":        "Info (this page)",
			"/health":  "Health check",
			"/send":    "POST - Send email",
			"/compose": "POST - Create draft for review",
			"/drafts":  "GET - List pending drafts",
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

// handleSend sends an email immediately
func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	email := &Email{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
	}

	result, err := s.sender.Send(email)
	if err != nil {
		log.Printf("Send failed: %v", err)
	} else {
		log.Printf("Email sent to %s: %s", email.To, result.MessageID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleCompose creates a draft for review
func (s *Server) handleCompose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	email := &Email{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
		From:    s.config.FromAddress,
	}

	// Validate
	if err := email.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create draft
	draft := &Draft{
		ID:        fmt.Sprintf("draft-%d", time.Now().UnixNano()),
		Email:     email,
		CreatedAt: time.Now(),
	}

	s.draftsMu.Lock()
	s.drafts[draft.ID] = draft
	s.draftsMu.Unlock()

	log.Printf("Draft created: %s (to: %s)", draft.ID, email.To)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draft)
}

// handleDrafts lists or sends drafts
func (s *Server) handleDrafts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listDrafts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listDrafts returns all pending drafts
func (s *Server) listDrafts(w http.ResponseWriter, r *http.Request) {
	s.draftsMu.RLock()
	drafts := make([]*Draft, 0, len(s.drafts))
	for _, d := range s.drafts {
		drafts = append(drafts, d)
	}
	s.draftsMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drafts)
}
