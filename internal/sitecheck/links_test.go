package sitecheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLinks: a link lands on a file or a folder with an index.html, quoted or not (the minifier unquotes), with its
// query or fragment ignored and its escapes undone; links off the site (another host, protocol-relative) aren't
// checked; each missing target is named once, with every page linking to it.
func TestLinks(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		path := filepath.Join(dir, filepath.FromSlash(name))
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, []byte(body), 0o644)
	}
	write("index.html", `<a href=/de/fleet>x</a><a href="/">home</a><img src=/images/logo.svg><a href="/doc/a%20b.pdf#p2">`+
		`<a href=https://example.com/x><a href=//cdn.example.com/y.js><a href="/about?x=1">`)
	write("de/fleet/index.html", `<a href=/de/fleet/ class=x><a href=/about>`)
	write("images/logo.svg", "<svg/>")
	write("doc/a b.pdf", "%PDF")

	var out strings.Builder
	missing, err := Links(dir, &out)
	if err != nil {
		t.Fatal(err)
	}
	if missing != 1 || !strings.Contains(out.String(), "✗ /about: no such page or file, linked from /de/fleet/index.html, /index.html") {
		t.Errorf("missing %d:\n%s", missing, out.String())
	}

	write("about/index.html", "")
	out.Reset()
	if missing, err := Links(dir, &out); err != nil || missing != 0 || !strings.Contains(out.String(), "every one lands") {
		t.Errorf("missing %d, %v:\n%s", missing, err, out.String())
	}

	if _, err := Links(t.TempDir(), &out); err == nil {
		t.Error("an empty folder passed: no site was built")
	}
}
