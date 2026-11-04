package main

import "fmt"

// css is theme css, the task theme:css: assets/css/generated-theme.css made from data/theme.json by Hugo Plate's own
// generator. Commit it: the site's builds use it as committed, and theme:check, in CI, proves it current.
func css() {
	if err := writeCSS(); err != nil {
		fail(err.Error())
	}
	fmt.Printf("✓ %s (from %s)\n", themeCSS, themeJSON)
}
