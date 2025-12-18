package mailerlite

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/mailerlite/mailerlite-go"
)

// Client provides methods for interacting with the MailerLite API.
type Client struct {
	sdk *mailerlite.Client
}

// NewClient creates a new MailerLite client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		sdk: mailerlite.NewClient(apiKey),
	}
}

// NewClientFromEnv creates a new MailerLite client using the
// MAILERLITE_API_KEY environment variable.
func NewClientFromEnv() (*Client, error) {
	apiKey := os.Getenv("MAILERLITE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("MAILERLITE_API_KEY environment variable not set")
	}
	return NewClient(apiKey), nil
}

// AddSubscriber adds a new subscriber with the given email and optional parameters.
func (c *Client) AddSubscriber(ctx context.Context, email string, opts *AddSubscriberOptions) (*Subscriber, error) {
	params := &mailerlite.UpsertSubscriber{
		Email: email,
	}

	if opts != nil {
		if opts.Status != "" {
			params.Status = opts.Status
		}
		if len(opts.Fields) > 0 {
			// Convert map[string]string to map[string]interface{}
			fields := make(map[string]interface{}, len(opts.Fields))
			for k, v := range opts.Fields {
				fields[k] = v
			}
			params.Fields = fields
		}
		if len(opts.Groups) > 0 {
			params.Groups = opts.Groups
		}
	}

	resp, _, err := c.sdk.Subscriber.Upsert(ctx, params)
	if err != nil {
		return nil, err
	}

	return convertSubscriber(&resp.Data), nil
}

// GetSubscriber retrieves a subscriber by email address.
func (c *Client) GetSubscriber(ctx context.Context, email string) (*Subscriber, error) {
	opts := &mailerlite.GetSubscriberOptions{
		Email: email,
	}
	resp, _, err := c.sdk.Subscriber.Get(ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertSubscriber(&resp.Data), nil
}

// GetSubscriberByID retrieves a subscriber by their ID.
func (c *Client) GetSubscriberByID(ctx context.Context, subscriberID string) (*Subscriber, error) {
	opts := &mailerlite.GetSubscriberOptions{
		SubscriberID: subscriberID,
	}
	resp, _, err := c.sdk.Subscriber.Get(ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertSubscriber(&resp.Data), nil
}

// ListSubscribers lists subscribers with optional limit.
func (c *Client) ListSubscribers(ctx context.Context, opts *ListOptions) ([]Subscriber, error) {
	params := &mailerlite.ListSubscriberOptions{}
	if opts != nil && opts.Limit > 0 {
		params.Limit = opts.Limit
	} else {
		params.Limit = 25
	}

	resp, _, err := c.sdk.Subscriber.List(ctx, params)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(resp.Data))
	for i, s := range resp.Data {
		subscribers[i] = *convertSubscriber(&s)
	}
	return subscribers, nil
}

// DeleteSubscriber removes a subscriber by ID.
func (c *Client) DeleteSubscriber(ctx context.Context, subscriberID string) error {
	_, err := c.sdk.Subscriber.Delete(ctx, subscriberID)
	return err
}

// ListGroups lists all subscriber groups.
func (c *Client) ListGroups(ctx context.Context, opts *ListOptions) ([]Group, error) {
	params := &mailerlite.ListGroupOptions{}
	if opts != nil && opts.Limit > 0 {
		params.Limit = opts.Limit
	} else {
		params.Limit = 25
	}

	resp, _, err := c.sdk.Group.List(ctx, params)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, len(resp.Data))
	for i, g := range resp.Data {
		groups[i] = convertGroup(&g)
	}
	return groups, nil
}

// CreateGroup creates a new subscriber group.
func (c *Client) CreateGroup(ctx context.Context, name string) (*Group, error) {
	resp, _, err := c.sdk.Group.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	g := convertGroup(&resp.Data)
	return &g, nil
}

// UpdateGroup updates a group's name.
func (c *Client) UpdateGroup(ctx context.Context, groupID, name string) (*Group, error) {
	resp, _, err := c.sdk.Group.Update(ctx, groupID, name)
	if err != nil {
		return nil, err
	}
	g := convertGroup(&resp.Data)
	return &g, nil
}

// DeleteGroup removes a group by ID.
func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := c.sdk.Group.Delete(ctx, groupID)
	return err
}

// AssignSubscriberToGroup adds a subscriber to a group.
func (c *Client) AssignSubscriberToGroup(ctx context.Context, groupID, subscriberID string) error {
	_, _, err := c.sdk.Group.Assign(ctx, groupID, subscriberID)
	return err
}

// UnassignSubscriberFromGroup removes a subscriber from a group.
func (c *Client) UnassignSubscriberFromGroup(ctx context.Context, groupID, subscriberID string) error {
	_, err := c.sdk.Group.UnAssign(ctx, groupID, subscriberID)
	return err
}

// GetSubscribersInGroup lists subscribers in a specific group.
func (c *Client) GetSubscribersInGroup(ctx context.Context, groupID string, opts *ListOptions) ([]Subscriber, error) {
	params := &mailerlite.ListGroupSubscriberOptions{
		GroupID: groupID,
	}
	if opts != nil && opts.Limit > 0 {
		params.Limit = opts.Limit
	} else {
		params.Limit = 25
	}

	resp, _, err := c.sdk.Group.Subscribers(ctx, params)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(resp.Data))
	for i, s := range resp.Data {
		subscribers[i] = *convertSubscriber(&s)
	}
	return subscribers, nil
}

// helper functions to convert SDK types to our types

func convertSubscriber(s *mailerlite.Subscriber) *Subscriber {
	sub := &Subscriber{
		ID:     s.ID,
		Email:  s.Email,
		Status: s.Status,
		Source: s.Source,
	}

	// Convert fields from map[string]interface{} to map[string]string
	if len(s.Fields) > 0 {
		sub.Fields = make(map[string]string, len(s.Fields))
		for k, v := range s.Fields {
			if str, ok := v.(string); ok {
				sub.Fields[k] = str
			}
		}
		// Extract name from fields if present
		if name, ok := s.Fields["name"].(string); ok {
			sub.Name = name
		}
	}

	if s.SubscribedAt != "" {
		if t, err := time.Parse(time.RFC3339, s.SubscribedAt); err == nil {
			sub.SubscribedAt = t
		}
	}
	if s.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, s.CreatedAt); err == nil {
			sub.CreatedAt = t
		}
	}
	if s.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, s.UpdatedAt); err == nil {
			sub.UpdatedAt = t
		}
	}

	return sub
}

func convertGroup(g *mailerlite.Group) Group {
	group := Group{
		ID:                g.ID,
		Name:              g.Name,
		ActiveCount:       g.ActiveCount,
		SentCount:         g.SentCount,
		OpensCount:        g.OpensCount,
		ClicksCount:       g.ClicksCount,
		UnsubscribedCount: g.UnsubscribedCount,
		UnconfirmedCount:  g.UnconfirmedCount,
		BouncedCount:      g.BouncedCount,
		JunkCount:         g.JunkCount,
	}

	group.OpenRate = Rate{
		Float:  g.OpenRate.Float,
		String: g.OpenRate.String,
	}
	group.ClickRate = Rate{
		Float:  g.ClickRate.Float,
		String: g.ClickRate.String,
	}

	if g.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, g.CreatedAt); err == nil {
			group.CreatedAt = t
		}
	}

	return group
}
