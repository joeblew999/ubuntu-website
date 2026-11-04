package main

import (
	"bytes"
	"fmt"
	"os"
)

// check is theme check, the task theme:check (root:ci runs it): assets/css/generated-theme.css is what data/theme.json
// makes, so the CSS the builds use as committed is current.
func check() {
	want, err := generated()
	if err != nil {
		fail(err.Error())
	}
	if have, _ := os.ReadFile(themeCSS); !bytes.Equal(have, want) {
		fail(fmt.Sprintf("%s isn't what %s makes: mise run theme:css, and commit it", themeCSS, themeJSON))
	}
	fmt.Printf("✓ %s is what %s makes\n", themeCSS, themeJSON)
}
