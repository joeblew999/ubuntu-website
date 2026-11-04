package main

import (
	"fmt"
	"os"
	"strings"
)

// status is theme status, the task theme:status: the theme's modules and the commit each is pinned to, as
// _vendor/modules.txt records them.
func status() {
	b, err := os.ReadFile(modulesTxt)
	if err != nil {
		fail(err.Error())
	}
	for line := range strings.Lines(string(b)) {
		if m, ok := strings.CutPrefix(line, "# "); ok {
			fmt.Print(m)
		}
	}
}

// pinned is the version _vendor/modules.txt pins mod to (v0.0.0-<date>-<commit>), or "".
func pinned(mod string) string {
	b, _ := os.ReadFile(modulesTxt)
	for line := range strings.Lines(string(b)) {
		if f := strings.Fields(line); len(f) >= 3 && f[0] == "#" && f[1] == mod {
			return f[2]
		}
	}
	return ""
}

// commit is the commit of a pinned version, its part after the last "-".
func commit(version string) string {
	return version[strings.LastIndex(version, "-")+1:]
}
