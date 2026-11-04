package main

import (
	"fmt"
	"regexp"
	"time"
)

// auth is google-mcp auth, the task google-mcp:auth: a Google account signed in to the Google MCP server (it opens the
// browser), then check.
func auth() {
	srv, err := startServer()
	if err != nil {
		fail("google-mcp-server didn't start: " + err.Error())
	}
	defer srv.stop()
	srv.initialize()
	srv.send(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"accounts_add","arguments":{}}}`)
	signIn := regexp.MustCompile(`https://accounts\.google\.com/o/oauth2/auth\S*`)
	line, ok := srv.wait(srv.errLines, signIn, 60*time.Second)
	if !ok {
		fail("no sign-in URL from google-mcp-server:\n" + srv.tail(srv.errLines, 20))
	}
	url := signIn.FindString(line)
	wait := seconds("GOOGLE_AUTH_TIMEOUT", 300)
	fmt.Printf("Sign in with Google in the browser (listening on localhost:8080 for %d s).\n", wait)
	fmt.Println("Unverified-app warning: Advanced → Go to <app> (unsafe) → Allow.")
	fmt.Println("  " + url)
	browse(url)
	if _, ok := srv.wait(srv.errLines, regexp.MustCompile(`Successfully added account`), time.Duration(wait)*time.Second); !ok {
		fail(fmt.Sprintf("no sign-in within %d s", wait))
	}
	fmt.Println("✓ account signed in")
	srv.stop()
	check()
}
