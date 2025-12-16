package mailerlite

import "time"

// Subscriber represents a MailerLite subscriber.
type Subscriber struct {
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Status        string            `json:"status"`
	Source        string            `json:"source"`
	Name          string            `json:"name"`
	Fields        map[string]string `json:"fields"`
	Groups        []Group           `json:"groups"`
	SubscribedAt  time.Time         `json:"subscribed_at"`
	UnsubscribedAt *time.Time       `json:"unsubscribed_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// Group represents a MailerLite subscriber group.
type Group struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ActiveCount      int       `json:"active_count"`
	SentCount        int       `json:"sent_count"`
	OpensCount       int       `json:"opens_count"`
	OpenRate         Rate      `json:"open_rate"`
	ClicksCount      int       `json:"clicks_count"`
	ClickRate        Rate      `json:"click_rate"`
	UnsubscribedCount int      `json:"unsubscribed_count"`
	UnconfirmedCount int       `json:"unconfirmed_count"`
	BouncedCount     int       `json:"bounced_count"`
	JunkCount        int       `json:"junk_count"`
	CreatedAt        time.Time `json:"created_at"`
}

// Rate represents a percentage rate.
type Rate struct {
	Float  float64 `json:"float"`
	String string  `json:"string"`
}

// Stats represents account statistics.
type Stats struct {
	SubscribersTotal   int `json:"subscribers_total"`
	SubscribersActive  int `json:"subscribers_active"`
	SubscribersUnsubscribed int `json:"subscribers_unsubscribed"`
	SubscribersBounced int `json:"subscribers_bounced"`
	SubscribersJunk    int `json:"subscribers_junk"`
	GroupsTotal        int `json:"groups_total"`
	CampaignsTotal     int `json:"campaigns_total"`
}

// AddSubscriberOptions contains optional parameters for adding a subscriber.
type AddSubscriberOptions struct {
	// Name is the subscriber's full name.
	Name string

	// Fields contains custom field values.
	Fields map[string]string

	// Groups is a list of group IDs to add the subscriber to.
	Groups []string

	// Status sets the subscriber status (active, unsubscribed, unconfirmed, bounced, junk).
	// Default is "active".
	Status string

	// Resubscribe if true, resubscribes a previously unsubscribed contact.
	Resubscribe bool
}

// ListOptions contains pagination options for list operations.
type ListOptions struct {
	// Limit is the maximum number of items to return.
	Limit int

	// Cursor is the pagination cursor for fetching the next page.
	Cursor string

	// Filter filters results by a specific field value.
	Filter string
}
