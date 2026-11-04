package cli

import (
	"flag"
	"slices"
	"testing"
)

func TestInterspersed(t *testing.T) {
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	fs.Bool("github-issue", false, "")
	fs.String("provider", "", "")
	for _, tc := range []struct {
		args, flags, rest []string
	}{
		{[]string{"content", "missing"}, nil, []string{"content", "missing"}},
		{[]string{"content", "missing", "--github-issue"}, []string{"--github-issue"}, []string{"content", "missing"}},
		{[]string{"-github-issue", "content"}, []string{"-github-issue"}, []string{"content"}},
		{[]string{"missing", "--provider", "deepl", "vi"}, []string{"--provider", "deepl"}, []string{"missing", "vi"}},
		{[]string{"missing", "--provider=deepl", "vi"}, []string{"--provider=deepl"}, []string{"missing", "vi"}},
		{[]string{"send", "--to", "x"}, nil, []string{"send", "--to", "x"}},
		{[]string{"a", "--", "--github-issue"}, nil, []string{"a", "--", "--github-issue"}},
	} {
		flags, rest := Interspersed(fs, tc.args)
		if !slices.Equal(flags, tc.flags) || !slices.Equal(rest, tc.rest) {
			t.Errorf("Interspersed(%q) = %q, %q; want %q, %q", tc.args, flags, rest, tc.flags, tc.rest)
		}
	}
}
