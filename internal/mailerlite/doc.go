// Package mailerlite provides a Go client library for the MailerLite API.
//
// This package wraps the official MailerLite Go SDK to provide a simplified
// interface for common email marketing operations like managing subscribers,
// groups, and campaigns.
//
// # Installation
//
//	go get www.ubuntusoftware.net/pkg/mailerlite
//
// # Usage
//
//	import "www.ubuntusoftware.net/pkg/mailerlite"
//
//	func main() {
//	    client := mailerlite.NewClient("your-api-key")
//
//	    // Add a subscriber
//	    sub, err := client.AddSubscriber(context.Background(), "user@example.com", "John Doe")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Printf("Added subscriber: %s\n", sub.Email)
//
//	    // List subscribers
//	    subs, err := client.ListSubscribers(context.Background(), 10)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    for _, s := range subs {
//	        fmt.Printf("- %s (%s)\n", s.Email, s.Status)
//	    }
//	}
//
// # Environment Variables
//
// The package can read the API key from the MAILERLITE_API_KEY environment
// variable if not provided directly:
//
//	client := mailerlite.NewClientFromEnv()
package mailerlite
