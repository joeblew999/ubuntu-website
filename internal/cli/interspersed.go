package cli

import (
	"flag"
	"strings"
)

// Interspersed separates the flags fs defines from the rest of args, wherever they are, each side in its order. A mise
// task passes its arguments after its run line (translate content missing --github-issue), and Go's flag package stops
// at the first argument that isn't a flag, so a binary's own flags would be lost after its command. A flag's value goes
// with it (--provider deepl) unless it's a bool flag or written --name=value; from "--" on, everything is the rest.
func Interspersed(fs *flag.FlagSet, args []string) (flags, rest []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			return flags, append(rest, args[i:]...)
		}
		f, inline := lookup(fs, a)
		if f == nil {
			rest = append(rest, a)
			continue
		}
		flags = append(flags, a)
		if b, ok := f.Value.(interface{ IsBoolFlag() bool }); !inline && !(ok && b.IsBoolFlag()) && i+1 < len(args) {
			i++
			flags = append(flags, args[i])
		}
	}
	return flags, rest
}

// lookup is the flag of fs that a names (-name, --name, -name=value or --name=value), and whether its value is in a.
func lookup(fs *flag.FlagSet, a string) (*flag.Flag, bool) {
	name, ok := strings.CutPrefix(a, "-")
	if !ok {
		return nil, false
	}
	name = strings.TrimPrefix(name, "-")
	name, _, inline := strings.Cut(name, "=")
	if name == "" || strings.HasPrefix(name, "-") {
		return nil, false
	}
	return fs.Lookup(name), inline
}
