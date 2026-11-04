package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Upstream packages we deliberately don't use (process-compose runs the dev processes instead).
var skipPackages = []string{"concurrently"}

// update is theme update [--ref <commit>], the task theme:update: Hugo Plate and the gethugothemes modules to the latest
// upstream (or Hugo Plate to the commit --ref names, to pin or roll back), repeatably. Hugo Modules pin the exact
// commits (_vendor/modules.txt) and hugo mod vendor makes each change a diff to review; on top of that it does what
// upstream's own themeUpdate.js (curl main | tar) doesn't:
//   - syncs cmd/theme/themeGenerator.js from the pinned commit, unmodified
//   - adds or bumps the npm packages upstream's package.json needs (tailwind-bootstrap-grid once broke the build)
//   - makes assets/css/generated-theme.css again, and runs a production build
//   - 3-way merges upstream's changes into the templates our layouts/ override: clean merges applied, for review;
//     conflicts written to .theme-merge/ (gitignored)
//
// go.mod and go.sum are put back: they belong to the Go code (go mod tidy strips Hugo's entries anyway), and
// _vendor/modules.txt is the theme's pin. Run again with nothing new upstream, it changes nothing.
func update(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	ref := fs.String("ref", "", "pin Hugo Plate to this commit (or roll back); default, the latest")
	fs.Parse(args)

	before := pinned(module)
	fmt.Printf("==> Modules (now: hugoplate %s)\n", commit(before))
	if err := updateModules(*ref); err != nil {
		fail("hugo mod get/vendor: " + err.Error())
	}
	after := pinned(module)
	sha := commit(after)
	fmt.Printf("    hugoplate %s → %s\n", commit(before), sha)

	fmt.Println("==> " + generator)
	changed, err := syncGenerator(after, sha)
	if err != nil {
		fail(err.Error())
	}
	if changed {
		fmt.Println("    updated from upstream @ " + sha)
	} else {
		fmt.Println("    unchanged")
	}

	fmt.Printf("==> npm packages (vs upstream package.json @ %s)\n", sha)
	add, err := packagesToAdd()
	if err != nil {
		fail(err.Error())
	}
	if len(add) == 0 {
		fmt.Println("    up to date")
	} else {
		fmt.Println("    bun add -d " + strings.Join(add, " "))
		if err := run("bun", append([]string{"add", "-d"}, add...)...); err != nil {
			fail("bun add: " + err.Error())
		}
	}

	fmt.Println("==> Theme CSS")
	if err := writeCSS(); err != nil {
		fail(err.Error())
	}
	fmt.Printf("    %s (from %s)\n", themeCSS, themeJSON)

	fmt.Println("==> Upstream changes to templates our layouts/ override (3-way merge, base = previously vendored)")
	conflicts, err := mergeOverrides()
	if err != nil {
		fail(err.Error())
	}

	fmt.Println("==> Production build")
	if err := productionBuild(); err != nil {
		fail(err.Error())
	}

	fmt.Println()
	if conflicts > 0 {
		fail(fmt.Sprintf("%d template conflict(s): resolve the files in .theme-merge/ into layouts/, then run it again", conflicts))
	}
	untracked, _ := output("git", "ls-files", "--others", "--exclude-standard", "--", "_vendor")
	if run("git", "diff", "--quiet", "--", "_vendor", "package.json", "bun.lock", generator, themeCSS, "layouts") == nil && untracked == "" {
		fmt.Printf("✓ Already up to date (hugoplate %s)\n", sha)
		return
	}
	fmt.Println("✓ Updated. Review: git diff --stat -- _vendor package.json bun.lock cmd/theme assets/css layouts")
	fmt.Println("  Ship: push a branch, check the Cloudflare preview build, then merge to main.")
}

// updateModules is hugo mod get -u (then Hugo Plate at ref, if one is given) and hugo mod vendor, with go.mod and
// go.sum put back as they were, whatever happens.
func updateModules(ref string) error {
	gomod, err := os.ReadFile("go.mod")
	if err != nil {
		return err
	}
	gosum, err := os.ReadFile("go.sum")
	if err != nil {
		return err
	}
	defer func() {
		os.WriteFile("go.mod", gomod, 0o644)
		os.WriteFile("go.sum", gosum, 0o644)
	}()
	if err := run("hugo", "mod", "get", "-u"); err != nil {
		return err
	}
	if ref != "" {
		if err := run("hugo", "mod", "get", module+"@"+ref); err != nil {
			return err
		}
	}
	return run("hugo", "mod", "vendor")
}

// syncGenerator writes upstream's scripts/themeGenerator.js, at the pinned version, to cmd/theme/: from the module
// cache if it's there (_vendor/ holds only what Hugo mounts), else from GitHub. It says whether that changed it.
func syncGenerator(version, sha string) (bool, error) {
	gomodcache, err := output("go", "env", "GOMODCACHE")
	if err != nil {
		return false, fmt.Errorf("go env GOMODCACHE: %w", err)
	}
	js, err := os.ReadFile(filepath.Join(strings.TrimSpace(gomodcache), module+"@"+version, "scripts", "themeGenerator.js"))
	if err != nil {
		url := "https://raw.githubusercontent.com/zeon-studio/hugoplate/" + sha + "/scripts/themeGenerator.js"
		resp, err := http.Get(url)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("%s: %s", url, resp.Status)
		}
		if js, err = io.ReadAll(resp.Body); err != nil {
			return false, err
		}
	}
	old, _ := os.ReadFile(generator)
	if bytes.Equal(old, js) {
		return false, nil
	}
	return true, os.WriteFile(generator, js, 0o644)
}

// packagesToAdd are the npm packages of upstream's package.json, as vendored, that ours lacks or has older:
// name@range, upstream's range, for bun add -d.
func packagesToAdd() ([]string, error) {
	type packageJSON struct{ Dependencies, DevDependencies map[string]string }
	read := func(path string) (map[string]string, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var p packageJSON
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		deps := map[string]string{}
		for name, r := range p.Dependencies {
			deps[name] = r
		}
		for name, r := range p.DevDependencies {
			deps[name] = r
		}
		return deps, nil
	}
	upstream, err := read(filepath.Join("_vendor", module, "package.json"))
	if err != nil {
		return nil, err
	}
	ours, err := read("package.json")
	if err != nil {
		return nil, err
	}
	var add []string
	for _, name := range slices.Sorted(maps.Keys(upstream)) {
		if slices.Contains(skipPackages, name) {
			continue
		}
		r, have := upstream[name], ours[name]
		if have == "" || newer(strings.TrimLeft(r, "^~"), strings.TrimLeft(have, "^~")) {
			add = append(add, name+"@"+r)
		}
	}
	return add, nil
}

// newer reports whether version a is newer than b, comparing their dot-separated parts as numbers where both are (as
// sort -V does).
func newer(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) && i < len(pb); i++ {
		na, errA := strconv.Atoi(pa[i])
		nb, errB := strconv.Atoi(pb[i])
		switch {
		case errA == nil && errB == nil && na != nb:
			return na > nb
		case (errA != nil || errB != nil) && pa[i] != pb[i]:
			return pa[i] > pb[i]
		}
	}
	return len(pa) > len(pb)
}

// mergeOverrides 3-way merges upstream's changes to each theme template our layouts/ override: base, the template as
// vendored before (HEAD); theirs, as vendored now. A clean merge is applied (review it: a clean merge can still drop
// something we rely on); a conflict is written to .theme-merge/ for hand resolution, and counted.
func mergeOverrides() (conflicts int, err error) {
	os.RemoveAll(".theme-merge")
	changed, err := output("git", "diff", "--name-only", "HEAD", "--", "_vendor/*/layouts/*")
	if err != nil {
		return 0, fmt.Errorf("git diff: %w", err)
	}
	added, err := output("git", "ls-files", "--others", "--exclude-standard", "--", "_vendor/*/layouts/*")
	if err != nil {
		return 0, fmt.Errorf("git ls-files: %w", err)
	}
	tmp, err := os.MkdirTemp("", "theme-merge-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tmp)
	overridden := 0
	for _, theirs := range strings.Fields(changed + added) {
		_, rel, _ := strings.Cut(theirs, "/layouts/")
		ours := ""
		for _, c := range []string{rel, swapPrefix(rel, "partials/", "_partials/"), swapPrefix(rel, "_partials/", "partials/")} {
			if st, err := os.Stat(filepath.Join("layouts", c)); err == nil && !st.IsDir() {
				ours = filepath.Join("layouts", c)
				break
			}
		}
		if ours == "" {
			continue
		}
		overridden++
		if _, err := os.Stat(theirs); err != nil {
			fmt.Printf("    %s: upstream removed or moved %s (is our override still used?)\n", ours, rel)
			continue
		}
		base := filepath.Join(tmp, "base")
		b, _ := exec.Command("git", "show", "HEAD:"+theirs).Output() // none: the template is new upstream
		if err := os.WriteFile(base, b, 0o644); err != nil {
			return conflicts, err
		}
		merged, err := exec.Command("git", "merge-file", "-p", ours, base, theirs).Output()
		var exit *exec.ExitError
		switch {
		case err == nil:
			if mine, _ := os.ReadFile(ours); bytes.Equal(merged, mine) {
				fmt.Printf("    %s: nothing to apply\n", ours)
			} else if err := os.WriteFile(ours, merged, 0o644); err != nil {
				return conflicts, err
			} else {
				fmt.Printf("    %s: merged upstream changes\n", ours)
			}
		case errors.As(err, &exit) && exit.ExitCode() > 0 && exit.ExitCode() < 128: // the number of conflicts
			out := filepath.Join(".theme-merge", strings.ReplaceAll(strings.TrimPrefix(filepath.ToSlash(ours), "layouts/"), "/", "_"))
			if err := os.MkdirAll(".theme-merge", 0o755); err != nil {
				return conflicts, err
			}
			if err := os.WriteFile(out, merged, 0o644); err != nil {
				return conflicts, err
			}
			fmt.Printf("    %s: CONFLICT, resolve from %s\n", ours, out)
			conflicts++
		default:
			return conflicts, fmt.Errorf("git merge-file %s: %w", ours, err)
		}
	}
	if overridden == 0 {
		fmt.Println("    none")
	}
	return conflicts, nil
}

// swapPrefix is s with the prefix from replaced by to, or s.
func swapPrefix(s, from, to string) string {
	if rest, ok := strings.CutPrefix(s, from); ok {
		return to + rest
	}
	return s
}

// productionBuild is hugo -e production --gc --minify, into a temporary folder: its warnings and its summary shown,
// its errors if it fails.
func productionBuild() error {
	tmp, err := os.MkdirTemp("", "theme-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	log, err := exec.Command("hugo", "-e", "production", "--gc", "--minify", "-d", filepath.Join(tmp, "public")).CombinedOutput()
	lines := strings.Split(string(log), "\n")
	if err != nil {
		shown := 0
		for _, l := range lines {
			if (strings.HasPrefix(l, "ERROR") || strings.HasPrefix(l, "Error")) && shown < 5 {
				fmt.Println("    " + l)
				shown++
			}
		}
		return fmt.Errorf("the production build fails (hugo -e production): fix it before committing")
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "WARN") || strings.Contains(l, "Pages ") || strings.Contains(l, "Total in") {
			fmt.Println("    " + l)
		}
	}
	return nil
}
