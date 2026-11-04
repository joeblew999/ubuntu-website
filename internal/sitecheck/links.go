package sitecheck

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// link is an href or src on the site itself: one slash, then the path up to a query, a fragment or the attribute's end
// (Hugo's minifier leaves attributes unquoted where it can).
var link = regexp.MustCompile(`(?:href|src)="?(/[^/"\s>#?][^"\s>#?]*|/)(?:["\s>#?]|$)`)

// Links checks that every link of the built site in dir (href or src on the site itself, in every HTML page) lands on
// something the build made: a file, or a folder with an index.html. It lists each missing target with the pages that
// link to it, and returns how many are missing.
func Links(dir string, stdout io.Writer) (int, error) {
	missing := map[string][]string{}
	pages, links := 0, 0
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}
		html, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		pages++
		page := "/" + filepath.ToSlash(strings.TrimPrefix(path, dir+string(filepath.Separator)))
		for _, m := range link.FindAllSubmatch(html, -1) {
			links++
			target := string(m[1])
			if p, err := url.PathUnescape(target); err == nil {
				target = p
			}
			if !built(dir, target) && !contains(missing[target], page) {
				missing[target] = append(missing[target], page)
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if pages == 0 {
		return 0, fmt.Errorf("no HTML in %s: build the site first (mise run root:build)", dir)
	}
	targets := make([]string, 0, len(missing))
	for t := range missing {
		targets = append(targets, t)
	}
	sort.Strings(targets)
	for _, t := range targets {
		sort.Strings(missing[t])
		fmt.Fprintf(stdout, "✗ %s: no such page or file, linked from %s\n", t, strings.Join(missing[t], ", "))
	}
	if len(targets) == 0 {
		fmt.Fprintf(stdout, "✓ %d links in %d pages, every one lands\n", links, pages)
	}
	return len(targets), nil
}

// built reports whether the build in dir made target: a file, or a folder with an index.html.
func built(dir, target string) bool {
	path := filepath.Join(dir, filepath.FromSlash(target))
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return true
	}
	_, err := os.Stat(filepath.Join(path, "index.html"))
	return err == nil
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
