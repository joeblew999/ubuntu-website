package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

// tools is upgrade tools, the task upgrade:tools: every tool mise.toml pins to its latest version (mise upgrade
// --bump), and Go to its latest (go.mod's toolchain line), all installed. Hugo is left alone, as it goes through a
// Cloudflare preview first (mise run upgrade:hugo, which is said when there's a newer one).
func tools() {
	list, err := output("mise", "ls", "--no-global", "--current", "--json")
	if err != nil {
		fail("mise ls: " + err.Error())
	}
	var current map[string]json.RawMessage
	if err := json.Unmarshal([]byte(list), &current); err != nil {
		fail("mise ls: " + err.Error())
	}
	var bump []string
	for name := range current {
		if name != "go" && name != "hugo-extended" {
			bump = append(bump, name)
		}
	}
	slices.Sort(bump)
	if err := run("mise", append([]string{"upgrade", "--bump"}, bump...)...); err != nil {
		fail("mise upgrade: " + err.Error())
	}

	gomod, err := os.ReadFile("go.mod")
	if err != nil {
		fail(err.Error())
	}
	pin := ""
	if m := regexp.MustCompile(`(?m)^toolchain go(\S+)$`).FindSubmatch(gomod); m != nil {
		pin = string(m[1])
	}
	latest, err := output("mise", "latest", "go")
	if err != nil {
		fail("mise latest go: " + err.Error())
	}
	if pin != latest {
		if err := run("go", "mod", "edit", "-toolchain=go"+latest); err != nil {
			fail("go mod edit: " + err.Error())
		}
		if err := run("go", "mod", "tidy"); err != nil {
			fail("go mod tidy: " + err.Error())
		}
		fmt.Printf("toolchain go%s → go%s (go.mod)\n", pin, latest)
	}
	if err := run("mise", "install"); err != nil {
		fail("mise install: " + err.Error())
	}

	hugoPin, _ := output("mise", "config", "get", "-f", "mise.toml", "tools.hugo-extended")
	hugoLatest, _ := output("mise", "latest", "hugo-extended")
	if hugoPin != "" && hugoLatest != "" && hugoPin != hugoLatest {
		fmt.Printf("\nHugo %s → %s available: mise run upgrade:hugo\n", hugoPin, strings.TrimSpace(hugoLatest))
	}
	fmt.Println("\nReview: git diff -- mise.toml go.mod; then mise run root:ci and commit.")
}
