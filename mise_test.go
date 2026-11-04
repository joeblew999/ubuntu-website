package website

import (
	"testing"

	misecheck "github.com/joeblew999/mise-check"
)

// TestMise is misecheck (github.com/joeblew999/mise-check) on this repo's mise file and the Go its tasks run. The task
// that runs it says which by its env: mise run root:check-mise.
func TestMise(t *testing.T) { misecheck.Check(t, misecheck.FromEnv()) }
