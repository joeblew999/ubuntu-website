package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// generated is the CSS Hugo Plate's generator makes of data/theme.json. The generator finds its files from where it
// sits, and knows only upstream's own layouts (exampleSite/, themes/<name>/), not a site that imports the theme as a
// module, so it runs, as upstream wrote it, in a temporary folder shaped like exampleSite/.
func generated() ([]byte, error) {
	js, err := os.ReadFile(generator)
	if err != nil {
		return nil, err
	}
	theme, err := os.ReadFile(themeJSON)
	if err != nil {
		return nil, err
	}
	tmp, err := os.MkdirTemp("", "theme-css-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	for path, b := range map[string][]byte{
		"scripts/themeGenerator.js":   js,
		"exampleSite/data/theme.json": theme,
		"exampleSite/hugo.toml":       nil,
	} {
		path = filepath.Join(tmp, path)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, b, 0o644); err != nil {
			return nil, err
		}
	}
	if out, err := exec.Command("node", filepath.Join(tmp, "scripts", "themeGenerator.js")).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("node %s: %v\n%s", generator, err, out)
	}
	return os.ReadFile(filepath.Join(tmp, "assets", "css", "generated-theme.css"))
}

// writeCSS writes assets/css/generated-theme.css as the generator makes it.
func writeCSS() error {
	b, err := generated()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(themeCSS), 0o755); err != nil {
		return err
	}
	return os.WriteFile(themeCSS, b, 0o644)
}
