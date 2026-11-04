package main

import (
	"os/exec"
	"regexp"
	"strings"
)

// branch is url branch, the task url:branch: this branch's Worker preview, which Workers Builds deploys for every
// branch but main (once it's pushed and built): <branch>-ubuntusoftware-net.<subdomain>.workers.dev, the branch name
// lowercased with anything but letters and digits as dashes.
func branch() {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		fail("git rev-parse: " + err.Error())
	}
	name := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(strings.TrimSpace(string(out))), "-")
	if name == "main" {
		fail("main is production: url prod")
	}
	open("https://" + strings.Trim(name, "-") + "-ubuntusoftware-net.gedw99.workers.dev")
}
