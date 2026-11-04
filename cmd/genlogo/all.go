package main

import (
	"fmt"
	"os"
)

// all is genlogo all, the task genlogo:all: every asset, then the places to re-upload them by hand, which take
// copies, not links.
func all() {
	if code := run(append([]string{"-asset", "all"}, os.Args[2:]...)); code != 0 {
		os.Exit(code)
	}
	fmt.Println("\nAssets generated. Next steps:")
	fmt.Println("  1. Deploy the website (push): the favicon, logos and og-image go with it")
	fmt.Println("  2. Gmail signature: re-upload email-logo.png at https://mail.google.com/mail/u/0/#settings/general")
	fmt.Println("  3. Bluesky profile: re-upload the avatar and banner at https://bsky.app/profile/ubuntusoftware.net")
}
