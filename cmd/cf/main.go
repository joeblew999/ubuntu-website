// Command cf is the Cloudflare side of the site that has no wrangler command of its own: the CDN cache purged, Email
// Routing (its rules, and adding or removing a forward), and the latest preview deployment. It calls Cloudflare's API
// with CLOUDFLARE_API_TOKEN (.env locally). Each command is a file of its own and a task of .mise/tasks/cf.toml
// (mise run cf:purge, mise run cf:email-status); wrangler's own (deployments, secrets, R2, Workers) are tool tasks.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	accountID = "7384af54e33b8a54ff240371ea368440"
	zoneID    = "867fb5996096b895b47114b6564ee0fb" // ubuntusoftware.net
	project   = "ubuntu-website"                   // the Pages project
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: cf purge | email-status | email-add <from> <to> | email-remove <from> | preview")
	}
	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "purge":
		purge()
	case "email-status":
		emailStatus()
	case "email-add":
		emailAdd(args)
	case "email-remove":
		emailRemove(args)
	case "preview":
		preview()
	default:
		fail("no command " + strconv.Quote(cmd))
	}
}

// fail says what went wrong and exits 1.
func fail(msg string) {
	fmt.Println("✗ " + msg)
	os.Exit(1)
}

// api calls Cloudflare's API (path under /client/v4), decoding its result into out; a request Cloudflare refuses is
// an error, with its messages.
func api(method, path string, body, out any) error {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		return fmt.Errorf("CLOUDFLARE_API_TOKEN isn't set (.env)")
	}
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, "https://api.cloudflare.com/client/v4"+path, r)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var env struct {
		Success bool                       `json:"success"`
		Errors  []struct{ Message string } `json:"errors"`
		Result  json.RawMessage            `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return fmt.Errorf("%s %s: %s", method, path, resp.Status)
	}
	if !env.Success {
		msg := resp.Status
		for _, e := range env.Errors {
			msg += ": " + e.Message
		}
		return fmt.Errorf("%s %s: %s", method, path, msg)
	}
	if out != nil {
		return json.Unmarshal(env.Result, out)
	}
	return nil
}
