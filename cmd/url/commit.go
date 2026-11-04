package main

import (
	"os/exec"
	"strings"
)

// commit is url commit, the task url:commit: the Cloudflare preview deployment of this commit (once it's pushed and
// built).
func commit() {
	sha, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		fail("git rev-parse: " + err.Error())
	}
	open("https://" + strings.TrimSpace(string(sha)) + ".ubuntu-website.pages.dev")
}
