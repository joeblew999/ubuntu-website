package main

import (
	"flag"
	"fmt"

	"github.com/joeblew999/ubuntu-website/internal/hugobuild"
)

// hugo is upgrade hugo [--to <version>], the task upgrade:hugo: Hugo pinned to that version (the latest without --to)
// in mise.toml, the one place local, CI and Cloudflare's builds take it from, and a production build proved with it.
// Then a branch's Cloudflare preview build proves it there, before main.
func hugo(args []string) {
	fs := flag.NewFlagSet("hugo", flag.ExitOnError)
	to := fs.String("to", "", "the Hugo version to pin; default, the latest")
	fs.Parse(args)

	v := *to
	if v == "" {
		latest, err := output("mise", "latest", "hugo-extended")
		if err != nil {
			fail("mise latest hugo-extended: " + err.Error())
		}
		v = latest
	}
	pin, _ := output("mise", "config", "get", "-f", "mise.toml", "tools.hugo-extended")
	fmt.Printf("Hugo %s → %s\n", pin, v)
	if err := run("mise", "use", "--pin", "--path", "mise.toml", "hugo-extended@"+v); err != nil {
		fail("mise use: " + err.Error())
	}
	if err := run("mise", "install", "hugo-extended@"+v); err != nil {
		fail("mise install: " + err.Error())
	}
	if err := hugobuild.Production(); err != nil {
		fail(fmt.Sprintf("%s on Hugo %s; mise.toml is left at %s, back with: mise run upgrade:hugo --to %s", err, v, v, pin))
	}
	fmt.Println("\nNext: commit on a branch and push it; when its Cloudflare preview build passes, merge it to main")
}
