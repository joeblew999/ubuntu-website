package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// modules is upgrade modules [--force] [--ignore <regexp>], the task upgrade:modules: the Go module dependencies, major
// versions too, with go-mod-upgrade (asking which, unless --force; held back a day after a release), then go mod tidy
// and the direct dependencies that are deprecated or retracted, which go-mod-upgrade doesn't move off: their
// replacements are named so they can be switched to.
func modules(args []string) {
	fs := flag.NewFlagSet("modules", flag.ExitOnError)
	force := fs.Bool("force", false, "upgrade every module without asking")
	ignore := fs.String("ignore", "", "regexp of modules to leave alone")
	fs.Parse(args)

	upgrade := []string{"--cooldown", "1d"}
	if *force {
		upgrade = append(upgrade, "-f")
	}
	if *ignore != "" {
		upgrade = append(upgrade, "-i", *ignore)
	}
	if err := run("go-mod-upgrade", upgrade...); err != nil {
		fail("go-mod-upgrade: " + err.Error())
	}
	if err := run("go", "mod", "tidy"); err != nil {
		fail("go mod tidy: " + err.Error())
	}
	deprecated, err := deprecatedModules()
	if err != nil {
		fail(err.Error())
	}
	if len(deprecated) == 0 {
		fmt.Println("✓ No deprecated or retracted direct dependencies")
		return
	}
	fmt.Println("⚠ Deprecated/retracted direct dependencies: switch to the replacement named here")
	for _, d := range deprecated {
		fmt.Println("  " + d)
	}
}

// deprecatedModules are the direct dependencies that are deprecated (a moved module, say) or on a retracted version,
// each as path@version: why.
func deprecatedModules() ([]string, error) {
	out, err := exec.Command("go", "list", "-m", "-u", "-json", "all").Output()
	if err != nil {
		return nil, fmt.Errorf("go list -m -u: %w", err)
	}
	var found []string
	dec := json.NewDecoder(strings.NewReader(string(out)))
	for {
		var m struct {
			Path, Version, Deprecated string
			Indirect                  bool
			Retracted                 []string
		}
		if err := dec.Decode(&m); errors.Is(err, io.EOF) {
			return found, nil
		} else if err != nil {
			return nil, fmt.Errorf("go list -m -u: %w", err)
		}
		switch {
		case m.Indirect:
		case m.Deprecated != "":
			found = append(found, fmt.Sprintf("%s@%s: %s", m.Path, m.Version, m.Deprecated))
		case len(m.Retracted) > 0:
			found = append(found, fmt.Sprintf("%s@%s: retracted: %s", m.Path, m.Version, strings.Join(m.Retracted, "; ")))
		}
	}
}
