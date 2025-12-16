package mailerlite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mailerlite/mailerlite-go"
)

// Web3FormPayload represents the incoming webhook payload from Web3Forms.
type Web3FormPayload struct {
	// Standard Web3Forms fields
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	Subject   string `json:"subject"`
	AccessKey string `json:"access_key"`

	// Custom fields from Get Started form
	Company  string `json:"company"`
	Platform string `json:"platform"`
	Industry string `json:"industry"`
	UseCase  string `json:"usecase"`
}

// ServerConfig holds configuration for the webhook server.
type ServerConfig struct {
	Port    int
	GroupID string // Optional: auto-assign subscribers to this group
}

// handleServer starts a webhook server to receive Web3Forms submissions.
func (c *CLI) handleServer(args []string) error {
	cfg := ServerConfig{
		Port: 8086,
	}

	// Parse args for PORT= and GROUP_ID=
	for _, arg := range args {
		if strings.HasPrefix(arg, "PORT=") {
			fmt.Sscanf(strings.TrimPrefix(arg, "PORT="), "%d", &cfg.Port)
		} else if strings.HasPrefix(arg, "GROUP_ID=") {
			cfg.GroupID = strings.TrimPrefix(arg, "GROUP_ID=")
		}
	}

	return c.StartServer(cfg)
}

// StartServer starts the webhook server with the given configuration.
func (c *CLI) StartServer(cfg ServerConfig) error {
	c.println("╔══════════════════════════════════════════════════════════════╗")
	c.println("║         MailerLite Webhook Server                            ║")
	c.println("╚══════════════════════════════════════════════════════════════╝")
	c.println()
	c.printf("Starting server on port %d...\n", cfg.Port)
	c.println()
	c.println("Webhook URL (for Web3Forms):")
	c.printf("  http://localhost:%d/webhook\n", cfg.Port)
	c.println()
	c.println("For production, use a tunnel service like:")
	c.printf("  ngrok http %d\n", cfg.Port)
	c.printf("  cloudflared tunnel --url http://localhost:%d\n", cfg.Port)
	c.println()
	if cfg.GroupID != "" {
		c.printf("Auto-assigning to group: %s\n", cfg.GroupID)
	}
	c.println("Press Ctrl+C to stop")
	c.println(strings.Repeat("─", 60))

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Webhook endpoint
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		c.handleWebhookRequest(w, r, cfg.GroupID)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	return server.ListenAndServe()
}

func (c *CLI) handleWebhookRequest(w http.ResponseWriter, r *http.Request, groupID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	var payload Web3FormPayload

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			c.printf("[ERROR] Failed to parse JSON: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			c.printf("[ERROR] Failed to parse form: %v\n", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		payload = Web3FormPayload{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Message:  r.FormValue("message"),
			Subject:  r.FormValue("subject"),
			Company:  r.FormValue("company"),
			Platform: r.FormValue("platform"),
			Industry: r.FormValue("industry"),
			UseCase:  r.FormValue("usecase"),
		}
	}

	// Validate required fields
	if payload.Email == "" {
		c.println("[WARN] Received webhook without email")
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	// Log the submission
	c.printf("\n[%s] New submission\n", time.Now().Format("15:04:05"))
	c.printf("  Email:    %s\n", payload.Email)
	if payload.Name != "" {
		c.printf("  Name:     %s\n", payload.Name)
	}
	if payload.Company != "" {
		c.printf("  Company:  %s\n", payload.Company)
	}
	if payload.Platform != "" {
		c.printf("  Platform: %s\n", payload.Platform)
	}
	if payload.Industry != "" {
		c.printf("  Industry: %s\n", payload.Industry)
	}

	// Add subscriber to MailerLite
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	subscriber := &mailerlite.UpsertSubscriber{
		Email: payload.Email,
		Fields: map[string]interface{}{
			"name":    payload.Name,
			"company": payload.Company,
		},
	}

	result, _, err := c.client.sdk.Subscriber.Upsert(ctx, subscriber)
	if err != nil {
		c.printf("  [ERROR] Failed to add to MailerLite: %v\n", err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Received (MailerLite error logged)")
		return
	}

	c.printf("  [OK] Added to MailerLite: ID=%s, Status=%s\n", result.Data.ID, result.Data.Status)

	// Auto-assign to group if specified
	if groupID != "" {
		_, _, err := c.client.sdk.Group.Assign(ctx, groupID, result.Data.ID)
		if err != nil {
			c.printf("  [WARN] Failed to assign to group: %v\n", err)
		} else {
			c.printf("  [OK] Assigned to group: %s\n", groupID)
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
