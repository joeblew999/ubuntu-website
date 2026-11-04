package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// rotate is google-mcp rotate, the task google-mcp:rotate: you click Add secret on the client page and copy the new
// secret; this picks it up from the clipboard, stores it, proves it at Google, and clears the clipboard. If the new
// secret never works, the old one goes back.
func rotate() {
	if runtime.GOOS != "darwin" {
		fail("rotate reads the clipboard with pbpaste: macOS only")
	}
	client, err := secretValue("GOOGLE_CLIENT_ID")
	if err != nil {
		fail(err.Error())
	}
	old, err := secretValue("GOOGLE_CLIENT_SECRET")
	if err != nil {
		fail(err.Error())
	}
	url := fmt.Sprintf("https://console.cloud.google.com/auth/clients/%s?project=%s", client, os.Getenv("UBUNTU_WEBSITE_GCP_PROJECT"))
	wait := seconds("GOOGLE_ROTATE_TIMEOUT", 600)
	fmt.Printf("On the client page click Add secret, then copy the NEW secret (waiting %d s):\n  %s\n", wait, url)
	browse(url)
	fresh := ""
	for i := 0; i < wait && fresh == ""; i++ {
		out, _ := exec.Command("pbpaste").Output()
		if c := strings.Join(strings.Fields(string(out)), ""); strings.HasPrefix(c, "GOCSPX-") && c != old {
			fresh = c
			break
		}
		time.Sleep(time.Second)
	}
	if fresh == "" {
		fail("no new secret copied")
	}
	clear := exec.Command("pbcopy")
	clear.Stdin = strings.NewReader("")
	_ = clear.Run()
	if err := storeSecret("GOOGLE_CLIENT_SECRET", fresh); err != nil {
		fail("storing the new secret: " + err.Error())
	}
	fmt.Println("✓ new secret stored in the keychain (clipboard cleared)")
	for try := 1; try <= 6; try++ {
		err := checkAll()
		if err == nil {
			fmt.Println("✓ rotated. Now delete the OLD secret on the client page, then run: mise run google-mcp:check")
			return
		}
		fmt.Println(err)
		if try < 6 {
			fmt.Println("  retrying in 30 s (new secrets can take a moment at Google)…")
			time.Sleep(30 * time.Second)
		}
	}
	_ = storeSecret("GOOGLE_CLIENT_SECRET", old)
	fail("the new secret never worked; the old one is back in the keychain")
}
