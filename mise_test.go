package website

import (
	"os"
	"path/filepath"
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
}
