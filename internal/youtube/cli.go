package youtube

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// Run is the CLI entry point. args should be os.Args.
func Run(args []string, version string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, usage())
		return 1
	}

	switch args[1] {
	case "-h", "--help", "help":
		fmt.Fprintln(stdout, usage())
		return 0
	case "-v", "--version", "version":
		fmt.Fprintf(stdout, "youtube %s\n", version)
		return 0
	case "info":
		return runInfo(args[2:], stdout, stderr)
	case "download":
		return runDownload(args[2:], stdout, stderr)
	case "list":
		return runList(args[2:], stdout, stderr)
	case "refresh":
		return runRefresh(args[2:], stdout, stderr)
	case "crop":
		return runCrop(args[2:], stdout, stderr)
	case "trim":
		return runTrim(args[2:], stdout, stderr)
	case "scale":
		return runScale(args[2:], stdout, stderr)
	case "compress":
		return runCompress(args[2:], stdout, stderr)
	case "probe":
		return runProbe(args[2:], stdout, stderr)
	case "overlay":
		return runOverlay(args[2:], stdout, stderr)
	case "gif":
		return runGif(args[2:], stdout, stderr)
	case "thumbnail":
		return runThumbnail(args[2:], stdout, stderr)
	case "speed":
		return runSpeed(args[2:], stdout, stderr)
	case "rotate":
		return runRotate(args[2:], stdout, stderr)
	case "fade":
		return runFade(args[2:], stdout, stderr)
	case "concat":
		return runConcat(args[2:], stdout, stderr)
	case "mute":
		return runMute(args[2:], stdout, stderr)
	case "audio":
		return runReplaceAudio(args[2:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n%s", args[1], usage())
		return 1
	}
}

func runInfo(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube info", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var proxy string
	var noInstall bool
	var githubIssue bool

	fs.StringVar(&proxy, "proxy", "", "Proxy URL passed to yt-dlp")
	fs.BoolVar(&noInstall, "no-install", false, "Require yt-dlp/ffmpeg to already be present (no auto-download)")
	fs.BoolVar(&githubIssue, "github-issue", false, "Output markdown suitable for GitHub Issues")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube info [flags] <url>")
		return 1
	}

	client := NewClient(proxy, !noInstall)
	ctx := context.Background()

	info, err := client.Info(ctx, fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	printInfo(stdout, info, githubIssue)
	return 0
}

func runDownload(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube download", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var proxy string
	var noInstall bool
	var githubIssue bool
	var format string
	var quality string
	var output string
	var audioOnly bool
	var force bool

	fs.StringVar(&proxy, "proxy", "", "Proxy URL passed to yt-dlp")
	fs.BoolVar(&noInstall, "no-install", false, "Require yt-dlp/ffmpeg to already be present (no auto-download)")
	fs.BoolVar(&githubIssue, "github-issue", false, "Output markdown suitable for GitHub Issues")
	fs.StringVar(&format, "format", "", "yt-dlp format selector (overrides other selectors)")
	fs.StringVar(&quality, "quality", "", "Target max height (e.g. 720, 1080). Ignored when -format is set")
	fs.StringVar(&output, "output", "", "Output path/template (Taskfile usually sets this)")
	fs.BoolVar(&audioOnly, "audio-only", false, "Download audio only")
	fs.BoolVar(&force, "force", false, "Overwrite existing files")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube download [flags] <url>")
		return 1
	}

	client := NewClient(proxy, !noInstall)
	ctx := context.Background()

	result, err := client.Download(ctx, fs.Arg(0), DownloadOptions{
		Output:    output,
		Format:    format,
		Quality:   quality,
		AudioOnly: audioOnly,
		Force:     force,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	printDownload(stdout, result, githubIssue)
	return 0
}

func printInfo(w io.Writer, info *InfoResult, githubIssue bool) {
	if info == nil {
		fmt.Fprintln(w, "no info available")
		return
	}

	if githubIssue {
		fmt.Fprintf(w, "## YouTube Info: %s\n\n", info.Title)
		fmt.Fprintf(w, "- ID: `%s`\n", info.ID)
		if info.URL != "" {
			fmt.Fprintf(w, "- URL: %s\n", info.URL)
		}
		if info.Uploader != "" {
			fmt.Fprintf(w, "- Uploader: %s\n", info.Uploader)
		}
		if info.Duration > 0 {
			fmt.Fprintf(w, "- Duration: %s\n", humanDuration(info.Duration))
		}
		if info.Extractor != "" {
			fmt.Fprintf(w, "- Extractor: %s\n", info.Extractor)
		}

		if len(info.Formats) > 0 {
			fmt.Fprintln(w, "\n| id | ext | res | fps | vcodec | acodec | note | size |")
			fmt.Fprintln(w, "|---|---|---|---|---|---|---|---|")
			for _, f := range info.Formats {
				fmt.Fprintf(
					w,
					"| %s | %s | %s | %s | %s | %s | %s | %s |\n",
					emptyDash(f.ID),
					emptyDash(f.Ext),
					emptyDash(f.Resolution),
					emptyDash(f.FPS),
					emptyDash(f.VCodec),
					emptyDash(f.ACodec),
					emptyDash(f.Note),
					emptyDash(f.Size),
				)
			}
		}
		return
	}

	fmt.Fprintf(w, "Title:    %s\n", info.Title)
	fmt.Fprintf(w, "ID:       %s\n", info.ID)
	if info.URL != "" {
		fmt.Fprintf(w, "URL:      %s\n", info.URL)
	}
	if info.Uploader != "" {
		fmt.Fprintf(w, "Uploader: %s\n", info.Uploader)
	}
	if info.Duration > 0 {
		fmt.Fprintf(w, "Duration: %s\n", humanDuration(info.Duration))
	}
	if info.Extractor != "" {
		fmt.Fprintf(w, "Extractor: %s\n", info.Extractor)
	}

	if len(info.Formats) == 0 {
		return
	}

	fmt.Fprintln(w, "\nFormats:")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tEXT\tRES\tFPS\tVCodec\tACodec\tNOTE\tSIZE")
	for _, f := range info.Formats {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			f.ID, f.Ext, f.Resolution, f.FPS, f.VCodec, f.ACodec, f.Note, f.Size)
	}
	_ = tw.Flush()
}

func runList(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube list", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var dir string
	fs.StringVar(&dir, "dir", ".", "Directory to scan for video manifests")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	manifests, err := ListManifests(dir)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	if len(manifests) == 0 {
		fmt.Fprintln(stdout, "No videos found (no .json manifests)")
		return 0
	}

	tw := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FILENAME\tQUALITY\tSOURCE")
	for _, m := range manifests {
		quality := m.Quality
		if quality == "" {
			quality = "best"
		}
		if m.AudioOnly {
			quality = "audio"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", m.Filename, quality, m.SourceURL)
	}
	_ = tw.Flush()
	return 0
}

func runRefresh(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube refresh", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var dir string
	var proxy string
	var noInstall bool

	fs.StringVar(&dir, "dir", ".", "Directory containing video manifests")
	fs.StringVar(&proxy, "proxy", "", "Proxy URL passed to yt-dlp")
	fs.BoolVar(&noInstall, "no-install", false, "Require yt-dlp/ffmpeg to already be present")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	client := NewClient(proxy, !noInstall)
	ctx := context.Background()

	if err := client.RefreshAll(ctx, dir); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runCrop(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube crop", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var width, height, x, y int
	var output string
	var force bool

	fs.IntVar(&width, "w", 0, "Output width (required)")
	fs.IntVar(&height, "h", 0, "Output height (required)")
	fs.IntVar(&x, "x", 0, "X offset from left")
	fs.IntVar(&y, "y", 0, "Y offset from top")
	fs.StringVar(&output, "o", "", "Output file (default: input-cropped.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube crop -w WIDTH -h HEIGHT [-x X] [-y Y] <input>")
		return 1
	}

	if width <= 0 || height <= 0 {
		fmt.Fprintln(stderr, "error: -w and -h are required and must be positive")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Crop(ctx, CropOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Width:       width,
		Height:      height,
		X:           x,
		Y:           y,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runTrim(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube trim", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var start, end, output string
	var force bool

	fs.StringVar(&start, "start", "", "Start time (e.g., 00:00:10 or 10)")
	fs.StringVar(&end, "end", "", "End time (e.g., 00:01:30 or 90)")
	fs.StringVar(&output, "o", "", "Output file (default: input-trimmed.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube trim [-start TIME] [-end TIME] <input>")
		return 1
	}

	if start == "" && end == "" {
		fmt.Fprintln(stderr, "error: at least -start or -end is required")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Trim(ctx, TrimOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Start:       start,
		End:         end,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runScale(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube scale", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var width, height int
	var output string
	var force bool

	fs.IntVar(&width, "w", -1, "Target width (-1 for auto)")
	fs.IntVar(&height, "h", -1, "Target height (-1 for auto)")
	fs.StringVar(&output, "o", "", "Output file (default: input-scaled.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube scale [-w WIDTH] [-h HEIGHT] <input>")
		return 1
	}

	if width <= 0 && height <= 0 {
		fmt.Fprintln(stderr, "error: at least -w or -h must be positive")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Scale(ctx, ScaleOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Width:       width,
		Height:      height,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runCompress(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube compress", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var crf int
	var preset, output string
	var force bool

	fs.IntVar(&crf, "crf", 23, "Quality (18-28, lower=better, default 23)")
	fs.StringVar(&preset, "preset", "medium", "Speed preset (ultrafast/fast/medium/slow/veryslow)")
	fs.StringVar(&output, "o", "", "Output file (default: input-compressed.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube compress [-crf N] [-preset PRESET] <input>")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Compress(ctx, CompressOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		CRF:         crf,
		Preset:      preset,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runProbe(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube probe", flag.ContinueOnError)
	fs.SetOutput(stderr)

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube probe <input>")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	info, err := client.Probe(ctx, fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "File:       %s\n", fs.Arg(0))
	fmt.Fprintf(stdout, "Resolution: %dx%d\n", info.Width, info.Height)
	fmt.Fprintf(stdout, "Codec:      %s\n", info.Codec)
	fmt.Fprintf(stdout, "Duration:   %s seconds\n", info.Duration)
	fmt.Fprintf(stdout, "Size:       %s\n", humanBytes(info.Size))

	return 0
}

func runOverlay(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube overlay", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var image, position, output string
	var marginPercent, scalePercent float64
	var force bool

	fs.StringVar(&image, "image", "", "Path to overlay image (PNG with transparency recommended)")
	fs.StringVar(&position, "pos", "bottomright", "Position: topleft, topright, bottomleft, bottomright, center")
	fs.Float64Var(&marginPercent, "margin", 2.0, "Margin from edge as % of video width")
	fs.Float64Var(&scalePercent, "scale", 10.0, "Logo size as % of video width")
	fs.StringVar(&output, "o", "", "Output file (default: input-edited.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube overlay -image LOGO.png [-pos POSITION] [-scale %] [-margin %] <input>")
		return 1
	}

	if image == "" {
		fmt.Fprintln(stderr, "error: -image is required")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Overlay(ctx, OverlayOptions{
		EditOptions:   EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Image:         image,
		Position:      position,
		MarginPercent: marginPercent,
		ScalePercent:  scalePercent,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runGif(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube gif", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var start, output string
	var duration float64
	var width, fps int
	var force bool

	fs.StringVar(&start, "start", "0", "Start time (e.g., 00:00:10 or 10)")
	fs.Float64Var(&duration, "duration", 5, "Duration in seconds")
	fs.IntVar(&width, "w", 480, "Output width (height auto)")
	fs.IntVar(&fps, "fps", 10, "Frames per second")
	fs.StringVar(&output, "o", "", "Output file (default: input.gif)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube gif [-start TIME] [-duration SECS] [-w WIDTH] [-fps FPS] <input>")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Gif(ctx, GifOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Start:       start,
		Duration:    duration,
		Width:       width,
		FPS:         fps,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runThumbnail(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube thumbnail", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var time, output string
	var force bool

	fs.StringVar(&time, "time", "00:00:01", "Time to extract frame")
	fs.StringVar(&output, "o", "", "Output file (default: input.jpg)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube thumbnail [-time TIME] <input>")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Thumbnail(ctx, ThumbnailOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Time:        time,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runSpeed(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube speed", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var factor float64
	var output string
	var force bool

	fs.Float64Var(&factor, "factor", 1.0, "Speed multiplier (0.5=half, 2.0=double)")
	fs.StringVar(&output, "o", "", "Output file (default: input-Nx.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube speed -factor N <input>")
		return 1
	}

	if factor <= 0 {
		fmt.Fprintln(stderr, "error: -factor must be positive")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Speed(ctx, SpeedOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Factor:      factor,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runRotate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube rotate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var degrees int
	var output string
	var force bool

	fs.IntVar(&degrees, "deg", 90, "Rotation degrees (90, 180, 270)")
	fs.StringVar(&output, "o", "", "Output file (default: input-rotN.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube rotate -deg DEGREES <input>")
		return 1
	}

	if degrees != 90 && degrees != 180 && degrees != 270 {
		fmt.Fprintln(stderr, "error: -deg must be 90, 180, or 270")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Rotate(ctx, RotateOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		Degrees:     degrees,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runFade(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube fade", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var fadeIn, fadeOut float64
	var output string
	var force bool

	fs.Float64Var(&fadeIn, "in", 0, "Fade in duration (seconds)")
	fs.Float64Var(&fadeOut, "out", 0, "Fade out duration (seconds)")
	fs.StringVar(&output, "o", "", "Output file (default: input-fade.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube fade [-in SECS] [-out SECS] <input>")
		return 1
	}

	if fadeIn <= 0 && fadeOut <= 0 {
		fmt.Fprintln(stderr, "error: at least -in or -out must be specified")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Fade(ctx, FadeOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		FadeIn:      fadeIn,
		FadeOut:     fadeOut,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runConcat(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube concat", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var output string
	var force bool

	fs.StringVar(&output, "o", "", "Output file (default: first-input-concat.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() < 2 {
		fmt.Fprintln(stderr, "usage: youtube concat <input1> <input2> [input3...] [-o output]")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Concat(ctx, ConcatOptions{
		EditOptions: EditOptions{Output: output, Force: force},
		Inputs:      fs.Args(),
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runMute(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube mute", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var output string
	var force bool

	fs.StringVar(&output, "o", "", "Output file (default: input-muted.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube mute <input>")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.Mute(ctx, MuteOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func runReplaceAudio(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("youtube audio", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var audioFile, output string
	var force bool

	fs.StringVar(&audioFile, "file", "", "Audio file to use (required)")
	fs.StringVar(&output, "o", "", "Output file (default: input-newaudio.mp4)")
	fs.BoolVar(&force, "force", false, "Overwrite output if exists")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: youtube audio -file AUDIO <input>")
		return 1
	}

	if audioFile == "" {
		fmt.Fprintln(stderr, "error: -file is required")
		return 1
	}

	client := NewClient("", true)
	ctx := context.Background()

	err := client.ReplaceAudio(ctx, ReplaceAudioOptions{
		EditOptions: EditOptions{Input: fs.Arg(0), Output: output, Force: force},
		AudioFile:   audioFile,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func printDownload(w io.Writer, result *DownloadResult, githubIssue bool) {
	if result == nil || result.Info == nil {
		fmt.Fprintln(w, "download complete")
		return
	}

	if githubIssue {
		fmt.Fprintf(w, "## Downloaded: %s\n\n", result.Info.Title)
		fmt.Fprintf(w, "- File: `%s`\n", result.Filename)
		if result.Info.URL != "" {
			fmt.Fprintf(w, "- URL: %s\n", result.Info.URL)
		}
		if result.Info.Duration > 0 {
			fmt.Fprintf(w, "- Duration: %s\n", humanDuration(result.Info.Duration))
		}
		return
	}

	fmt.Fprintf(w, "Downloaded %s -> %s\n", result.Info.Title, result.Filename)
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func humanDuration(d time.Duration) string {
	if d <= 0 {
		return ""
	}
	// Round to nearest second for readability.
	rounded := d.Truncate(time.Second)
	if rounded == 0 {
		rounded = time.Second
	}
	return rounded.String()
}

func usage() string {
	return `youtube <command> [flags]

Download Commands:
  info <url>      Print metadata and available formats
  download <url>  Download a video (or audio with -audio-only)
  list            List videos from manifest files in a directory
  refresh         Re-download all videos from their manifests

Edit Commands (FFmpeg-based post-processing):
  probe <file>    Show video info (resolution, codec, duration, size)
  crop <file>     Crop video (-w WIDTH -h HEIGHT [-x X] [-y Y])
  trim <file>     Trim video (-start TIME and/or -end TIME)
  scale <file>    Resize video (-w WIDTH and/or -h HEIGHT, -1 for auto)
  compress <file> Re-encode for smaller size (-crf N -preset PRESET)
  overlay <file>  Add logo/watermark (-image LOGO.png, scales to % of video width)
  gif <file>      Extract animated GIF (-start TIME -duration SECS -w WIDTH -fps N)
  thumbnail <file> Extract poster frame (-time TIME, outputs .jpg)
  speed <file>    Change playback speed (-factor N, e.g., 0.5 or 2.0)
  rotate <file>   Rotate video (-deg 90/180/270)
  fade <file>     Add fade in/out (-in SECS -out SECS)
  concat <files>  Join multiple videos (input1 input2 ...)
  mute <file>     Remove audio track
  audio <file>    Replace audio track (-file AUDIO)

Download flags:
  -proxy          Proxy URL passed to yt-dlp
  -no-install     Require yt-dlp/ffmpeg to already be present
  -github-issue   Output markdown for GitHub Issues
  -format         yt-dlp format selector (overrides quality/audio flags)
  -quality        Target max height (e.g. 720, 1080) when -format is not set
  -audio-only     Download audio only
  -output         Output path/template (Taskfile typically sets this)
  -force          Overwrite existing files
  -dir            Directory for list/refresh commands (default: current dir)

Edit flags (common):
  -o              Output file (default: adds appropriate suffix)
  -force          Overwrite output if exists

Edit flags (specific):
  -w, -h          Width/height for crop/scale/gif
  -x, -y          X/Y offset for crop (default: 0)
  -start, -end    Time positions for trim (e.g., "00:00:10" or "10")
  -crf            Quality for compress (18-28, lower=better, default 23)
  -preset         Speed for compress (ultrafast/fast/medium/slow/veryslow)
  -image          Path to overlay image (PNG with transparency recommended)
  -pos            Overlay position (topleft/topright/bottomleft/bottomright/center)
  -scale          Logo size as % of video width (default: 10)
  -margin         Margin as % of video width (default: 2)
  -duration       GIF duration in seconds (default: 5)
  -fps            GIF frames per second (default: 10)
  -time           Time for thumbnail extraction (default: 00:00:01)
  -factor         Speed multiplier (0.5=half, 2.0=double)
  -deg            Rotation degrees (90, 180, 270)
  -in, -out       Fade in/out duration in seconds
  -file           Audio file for replacement
`
}
