package website

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	misecheck "github.com/joeblew999/mise-check"
)

// TestMise is misecheck (github.com/joeblew999/mise-check) on every mise file here, run by mise run root:check-mise:
// mise.toml, the control plane, and each task file, .mise/tasks/<binary>.toml, strictly against its own binary,
// cmd/<binary>/: its commands a file each, run by its tasks, and the task names its Go prints all tasks of the file. A
// task file with no binary of that name (one that only runs tools) is checked for the rest.
func TestMise(t *testing.T) {
	t.Run("mise.toml", func(t *testing.T) { misecheck.Check(t, misecheck.Config{MiseFile: "mise.toml"}) })
	files, _ := filepath.Glob(filepath.Join(".mise", "tasks", "*.toml"))
	if len(files) == 0 {
		t.Fatal("no .mise/tasks/*.toml: did the task files move?")
	}
	for _, file := range files {
		c := misecheck.Config{MiseFile: file}
		binary := strings.TrimSuffix(filepath.Base(file), ".toml")
		if st, err := os.Stat(filepath.Join("cmd", binary)); err == nil && st.IsDir() {
			c.GoCmd, c.GoLayout, c.GoBinaries, c.GoCallers = "cmd", "file", binary, filepath.Join("cmd", binary)
		}
		t.Run(filepath.Base(file), func(t *testing.T) { misecheck.Check(t, c) })
	}
	t.Run("task names in Go", func(t *testing.T) { checkTaskNames(t, append([]string{"mise.toml"}, files...)) })
}

// checkTaskNames: every task the Go anywhere here names to run ("mise run translate:content", in a hint it prints, in
// internal/ as much as in cmd/) is a task, or an alias of one, of the mise files. misecheck holds each binary's own
// cmd/<binary> to its file; this is the rest.
func checkTaskNames(t *testing.T, miseFiles []string) {
	header := regexp.MustCompile(`(?m)^\[(?:tasks\.)?"?([a-z0-9][a-z0-9:-]*)"?\]$`)
	alias := regexp.MustCompile(`(?m)^alias\s*=\s*"([^"]+)"`)
	tasks := map[string]bool{}
	for _, f := range miseFiles {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range header.FindAllStringSubmatch(string(b), -1) {
			tasks[m[1]] = true
		}
		for _, m := range alias.FindAllStringSubmatch(string(b), -1) {
			tasks[m[1]] = true
		}
	}
	named := regexp.MustCompile(`mise run ([a-z0-9][a-z0-9-]*:[a-z0-9][a-z0-9-]*)`)
	checked := 0
	for _, root := range []string{"cmd", "internal"} {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}
			b, _ := os.ReadFile(path)
			for _, m := range named.FindAllStringSubmatch(string(b), -1) {
				checked++
				if !tasks[m[1]] {
					t.Errorf("%s names mise run %s, which isn't a task", path, m[1])
				}
			}
			return nil
		})
	}
	t.Logf("%d task names in the Go of cmd/ and internal/, against %d tasks", checked, len(tasks))
}
