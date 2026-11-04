// Command youtube downloads videos (yt-dlp) and edits them (ffmpeg) for the site's static/videos/, both downloaded into
// .build/ on first use. Each command is a file of its own and a task of .mise/tasks/youtube.toml; its arguments and
// flags are internal/youtube's, passed through (mise run youtube:download <url> -output ...).
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/youtube"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "info":
		info()
	case "download":
		download()
	case "list":
		list()
	case "refresh":
		refresh()
	case "crop":
		crop()
	case "trim":
		trim()
	case "scale":
		scale()
	case "compress":
		compress()
	case "probe":
		probe()
	case "overlay":
		overlay()
	case "gif":
		gif()
	case "thumbnail":
		thumbnail()
	case "speed":
		speed()
	case "rotate":
		rotate()
	case "fade":
		fade()
	case "concat":
		concat()
	case "mute":
		mute()
	case "audio":
		audio()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/youtube's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return youtube.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
