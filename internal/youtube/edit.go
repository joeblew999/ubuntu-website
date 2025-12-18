package youtube

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// EditOptions contains settings for video editing operations.
type EditOptions struct {
	Input  string
	Output string
	Force  bool
}

// CropOptions defines how to crop a video.
type CropOptions struct {
	EditOptions
	Width  int // Output width
	Height int // Output height
	X      int // X offset from left
	Y      int // Y offset from top
}

// TrimOptions defines how to trim a video.
type TrimOptions struct {
	EditOptions
	Start string // Start time (e.g., "00:00:10" or "10")
	End   string // End time (e.g., "00:01:30" or "90")
}

// ScaleOptions defines how to resize a video.
type ScaleOptions struct {
	EditOptions
	Width  int // Target width (-1 for auto based on height)
	Height int // Target height (-1 for auto based on width)
}

// CompressOptions defines how to re-encode for smaller file size.
type CompressOptions struct {
	EditOptions
	CRF     int    // Constant Rate Factor (18-28, lower = better quality, default 23)
	Preset  string // Encoding preset (ultrafast, fast, medium, slow, veryslow)
	MaxSize int    // Target max file size in MB (0 = no limit)
}

// OverlayOptions defines how to add an image overlay (logo/watermark) to a video.
type OverlayOptions struct {
	EditOptions
	Image         string  // Path to overlay image (PNG with transparency recommended)
	Position      string  // Position: topleft, topright, bottomleft, bottomright, center (default: bottomright)
	MarginPercent float64 // Margin from edge as % of video width (default: 2.0)
	ScalePercent  float64 // Scale overlay to % of video width (default: 10.0)
	Opacity       float64 // Opacity 0.0-1.0 (default: 1.0, requires re-encoding if < 1.0)
}

// GifOptions defines how to extract an animated GIF from video.
type GifOptions struct {
	EditOptions
	Start    string  // Start time
	Duration float64 // Duration in seconds (default: 5)
	Width    int     // Output width (default: 480, height auto)
	FPS      int     // Frames per second (default: 10)
}

// ThumbnailOptions defines how to extract a poster frame.
type ThumbnailOptions struct {
	EditOptions
	Time string // Time to extract frame (default: "00:00:01")
}

// SpeedOptions defines how to change video speed.
type SpeedOptions struct {
	EditOptions
	Factor float64 // Speed multiplier (0.5 = half speed, 2.0 = double speed)
}

// RotateOptions defines how to rotate video.
type RotateOptions struct {
	EditOptions
	Degrees int // Rotation: 90, 180, 270
}

// FadeOptions defines fade in/out transitions.
type FadeOptions struct {
	EditOptions
	FadeIn  float64 // Fade in duration in seconds (0 = no fade in)
	FadeOut float64 // Fade out duration in seconds (0 = no fade out)
}

// ConcatOptions defines how to join multiple videos.
type ConcatOptions struct {
	EditOptions
	Inputs []string // List of input files to concatenate
}

// MuteOptions defines audio removal.
type MuteOptions struct {
	EditOptions
}

// ReplaceAudioOptions defines audio replacement.
type ReplaceAudioOptions struct {
	EditOptions
	AudioFile string // Path to audio file
}

// Crop creates a cropped version of the video.
func (c *Client) Crop(ctx context.Context, opts CropOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-cropped")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// ffmpeg -i input.mp4 -vf "crop=w:h:x:y" -c:a copy output.mp4
	cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d", opts.Width, opts.Height, opts.X, opts.Y)
	args := []string{
		"-i", opts.Input,
		"-vf", cropFilter,
		"-c:a", "copy",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Cropping %s -> %s (crop=%dx%d+%d+%d)\n", opts.Input, output, opts.Width, opts.Height, opts.X, opts.Y)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Trim creates a trimmed version of the video.
func (c *Client) Trim(ctx context.Context, opts TrimOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-trimmed")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// ffmpeg -i input.mp4 -ss START -to END -c copy output.mp4
	// Using -c copy for fast trimming (no re-encoding)
	args := []string{"-i", opts.Input}
	if opts.Start != "" {
		args = append(args, "-ss", opts.Start)
	}
	if opts.End != "" {
		args = append(args, "-to", opts.End)
	}
	args = append(args, "-c", "copy")
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Trimming %s -> %s (start=%s, end=%s)\n", opts.Input, output, opts.Start, opts.End)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Scale resizes the video.
func (c *Client) Scale(ctx context.Context, opts ScaleOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-scaled")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// ffmpeg -i input.mp4 -vf "scale=w:h" -c:a copy output.mp4
	// Use -1 for auto-calculate maintaining aspect ratio
	w := strconv.Itoa(opts.Width)
	h := strconv.Itoa(opts.Height)
	if opts.Width <= 0 {
		w = "-1"
	}
	if opts.Height <= 0 {
		h = "-1"
	}
	scaleFilter := fmt.Sprintf("scale=%s:%s", w, h)

	args := []string{
		"-i", opts.Input,
		"-vf", scaleFilter,
		"-c:a", "copy",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Scaling %s -> %s (scale=%sx%s)\n", opts.Input, output, w, h)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Compress re-encodes the video with a smaller file size.
func (c *Client) Compress(ctx context.Context, opts CompressOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-compressed")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// Default CRF for good balance
	crf := opts.CRF
	if crf <= 0 {
		crf = 23
	}

	preset := opts.Preset
	if preset == "" {
		preset = "medium"
	}

	// ffmpeg -i input.mp4 -c:v libx264 -crf CRF -preset PRESET -c:a aac output.mp4
	args := []string{
		"-i", opts.Input,
		"-c:v", "libx264",
		"-crf", strconv.Itoa(crf),
		"-preset", preset,
		"-c:a", "aac",
		"-b:a", "128k",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Compressing %s -> %s (crf=%d, preset=%s)\n", opts.Input, output, crf, preset)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Overlay adds an image overlay (logo/watermark) to the video.
// Logo size and margin are calculated as percentages of video width for resolution independence.
func (c *Client) Overlay(ctx context.Context, opts OverlayOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-edited")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// Validate image exists
	if _, err := os.Stat(opts.Image); err != nil {
		return fmt.Errorf("overlay image not found: %s", opts.Image)
	}

	// Probe video to get dimensions
	info, err := c.Probe(ctx, opts.Input)
	if err != nil {
		return fmt.Errorf("probe video: %w", err)
	}

	// Default position
	position := opts.Position
	if position == "" {
		position = "bottomright"
	}

	// Default percentages (of video width)
	scalePercent := opts.ScalePercent
	if scalePercent <= 0 {
		scalePercent = 10.0 // Logo is 10% of video width
	}

	marginPercent := opts.MarginPercent
	if marginPercent <= 0 {
		marginPercent = 2.0 // Margin is 2% of video width
	}

	// Calculate actual pixel values from video dimensions
	logoWidth := int(float64(info.Width) * scalePercent / 100.0)
	margin := int(float64(info.Width) * marginPercent / 100.0)

	// Build overlay position expression
	// FFmpeg overlay filter: overlay=x:y
	var overlayX, overlayY string
	switch position {
	case "topleft":
		overlayX = fmt.Sprintf("%d", margin)
		overlayY = fmt.Sprintf("%d", margin)
	case "topright":
		overlayX = fmt.Sprintf("W-w-%d", margin)
		overlayY = fmt.Sprintf("%d", margin)
	case "bottomleft":
		overlayX = fmt.Sprintf("%d", margin)
		overlayY = fmt.Sprintf("H-h-%d", margin)
	case "bottomright":
		overlayX = fmt.Sprintf("W-w-%d", margin)
		overlayY = fmt.Sprintf("H-h-%d", margin)
	case "center":
		overlayX = "(W-w)/2"
		overlayY = "(H-h)/2"
	default:
		return fmt.Errorf("invalid position: %s (use topleft, topright, bottomleft, bottomright, center)", position)
	}

	// Build filter complex - always scale logo to percentage of video width
	filterComplex := fmt.Sprintf("[1:v]scale=%d:-1[logo];[0:v][logo]overlay=%s:%s",
		logoWidth, overlayX, overlayY)

	// ffmpeg -i video.mp4 -i logo.png -filter_complex "overlay=..." -c:a copy output.mp4
	args := []string{
		"-i", opts.Input,
		"-i", opts.Image,
		"-filter_complex", filterComplex,
		"-c:a", "copy",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Adding overlay %s -> %s\n", opts.Input, output)
	fmt.Printf("  Video: %dx%d, Logo: %dpx wide (%.0f%%), Margin: %dpx (%.0f%%)\n",
		info.Width, info.Height, logoWidth, scalePercent, margin, marginPercent)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Gif extracts an animated GIF from the video.
func (c *Client) Gif(ctx context.Context, opts GifOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = strings.TrimSuffix(opts.Input, filepath.Ext(opts.Input)) + ".gif"
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// Defaults
	duration := opts.Duration
	if duration <= 0 {
		duration = 5
	}
	width := opts.Width
	if width <= 0 {
		width = 480
	}
	fps := opts.FPS
	if fps <= 0 {
		fps = 10
	}
	start := opts.Start
	if start == "" {
		start = "0"
	}

	// ffmpeg -i input.mp4 -ss START -t DURATION -vf "fps=FPS,scale=WIDTH:-1:flags=lanczos" output.gif
	filter := fmt.Sprintf("fps=%d,scale=%d:-1:flags=lanczos", fps, width)
	args := []string{
		"-i", opts.Input,
		"-ss", start,
		"-t", fmt.Sprintf("%.2f", duration),
		"-vf", filter,
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Creating GIF %s -> %s (start=%s, duration=%.1fs, %dpx, %dfps)\n",
		opts.Input, output, start, duration, width, fps)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Thumbnail extracts a single frame as an image.
func (c *Client) Thumbnail(ctx context.Context, opts ThumbnailOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = strings.TrimSuffix(opts.Input, filepath.Ext(opts.Input)) + ".jpg"
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	time := opts.Time
	if time == "" {
		time = "00:00:01"
	}

	// ffmpeg -i input.mp4 -ss TIME -frames:v 1 output.jpg
	args := []string{
		"-i", opts.Input,
		"-ss", time,
		"-frames:v", "1",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Extracting thumbnail %s -> %s (time=%s)\n", opts.Input, output, time)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Speed changes the playback speed of the video.
func (c *Client) Speed(ctx context.Context, opts SpeedOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		suffix := fmt.Sprintf("-%.1fx", opts.Factor)
		output = addSuffix(opts.Input, suffix)
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	if opts.Factor <= 0 {
		return fmt.Errorf("speed factor must be positive")
	}

	// Video: setpts=PTS/FACTOR (lower = faster)
	// Audio: atempo=FACTOR (only supports 0.5-2.0, chain for more)
	videoPTS := 1.0 / opts.Factor
	videoFilter := fmt.Sprintf("setpts=%.4f*PTS", videoPTS)

	// For audio, atempo only supports 0.5-2.0, so chain multiple filters if needed
	var audioFilter string
	tempo := opts.Factor
	if tempo >= 0.5 && tempo <= 2.0 {
		audioFilter = fmt.Sprintf("atempo=%.4f", tempo)
	} else if tempo > 2.0 {
		// Chain atempo filters for speeds > 2x
		audioFilter = "atempo=2.0"
		tempo = tempo / 2.0
		for tempo > 2.0 {
			audioFilter += ",atempo=2.0"
			tempo = tempo / 2.0
		}
		if tempo > 1.0 {
			audioFilter += fmt.Sprintf(",atempo=%.4f", tempo)
		}
	} else {
		// Chain atempo filters for speeds < 0.5x
		audioFilter = "atempo=0.5"
		tempo = tempo / 0.5
		for tempo < 0.5 {
			audioFilter += ",atempo=0.5"
			tempo = tempo / 0.5
		}
		if tempo < 1.0 {
			audioFilter += fmt.Sprintf(",atempo=%.4f", tempo)
		}
	}

	args := []string{
		"-i", opts.Input,
		"-filter:v", videoFilter,
		"-filter:a", audioFilter,
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Changing speed %s -> %s (%.1fx)\n", opts.Input, output, opts.Factor)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Rotate rotates the video by 90, 180, or 270 degrees.
func (c *Client) Rotate(ctx context.Context, opts RotateOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, fmt.Sprintf("-rot%d", opts.Degrees))
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	var filter string
	switch opts.Degrees {
	case 90:
		filter = "transpose=1" // 90 clockwise
	case 180:
		filter = "transpose=1,transpose=1" // 180
	case 270:
		filter = "transpose=2" // 90 counter-clockwise
	default:
		return fmt.Errorf("rotation must be 90, 180, or 270 degrees")
	}

	args := []string{
		"-i", opts.Input,
		"-vf", filter,
		"-c:a", "copy",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Rotating %s -> %s (%d degrees)\n", opts.Input, output, opts.Degrees)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Fade adds fade in and/or fade out effects.
func (c *Client) Fade(ctx context.Context, opts FadeOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-fade")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	if opts.FadeIn <= 0 && opts.FadeOut <= 0 {
		return fmt.Errorf("at least one of fade_in or fade_out must be specified")
	}

	// Get video duration for fade out
	info, err := c.Probe(ctx, opts.Input)
	if err != nil {
		return fmt.Errorf("probe video: %w", err)
	}

	duration, _ := strconv.ParseFloat(info.Duration, 64)

	var videoFilters []string
	var audioFilters []string

	if opts.FadeIn > 0 {
		videoFilters = append(videoFilters, fmt.Sprintf("fade=t=in:st=0:d=%.2f", opts.FadeIn))
		audioFilters = append(audioFilters, fmt.Sprintf("afade=t=in:st=0:d=%.2f", opts.FadeIn))
	}

	if opts.FadeOut > 0 {
		fadeStart := duration - opts.FadeOut
		if fadeStart < 0 {
			fadeStart = 0
		}
		videoFilters = append(videoFilters, fmt.Sprintf("fade=t=out:st=%.2f:d=%.2f", fadeStart, opts.FadeOut))
		audioFilters = append(audioFilters, fmt.Sprintf("afade=t=out:st=%.2f:d=%.2f", fadeStart, opts.FadeOut))
	}

	args := []string{
		"-i", opts.Input,
		"-vf", strings.Join(videoFilters, ","),
		"-af", strings.Join(audioFilters, ","),
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Adding fade %s -> %s (in=%.1fs, out=%.1fs)\n", opts.Input, output, opts.FadeIn, opts.FadeOut)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Concat joins multiple videos into one.
func (c *Client) Concat(ctx context.Context, opts ConcatOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	if len(opts.Inputs) < 2 {
		return fmt.Errorf("concat requires at least 2 input files")
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Inputs[0], "-concat")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// Create concat file list
	listFile := output + ".txt"
	var lines []string
	for _, input := range opts.Inputs {
		absPath, _ := filepath.Abs(input)
		lines = append(lines, fmt.Sprintf("file '%s'", absPath))
	}
	if err := os.WriteFile(listFile, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("write concat list: %w", err)
	}
	defer os.Remove(listFile)

	// ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp4
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Concatenating %d files -> %s\n", len(opts.Inputs), output)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// Mute removes the audio track from the video.
func (c *Client) Mute(ctx context.Context, opts MuteOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-muted")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	// ffmpeg -i input.mp4 -an -c:v copy output.mp4
	args := []string{
		"-i", opts.Input,
		"-an",
		"-c:v", "copy",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Removing audio %s -> %s\n", opts.Input, output)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// ReplaceAudio replaces the video's audio track with a different audio file.
func (c *Client) ReplaceAudio(ctx context.Context, opts ReplaceAudioOptions) error {
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return fmt.Errorf("ensure ffmpeg: %w", err)
	}

	output := opts.Output
	if output == "" {
		output = addSuffix(opts.Input, "-newaudio")
	}

	if err := checkOverwrite(output, opts.Force); err != nil {
		return err
	}

	if opts.AudioFile == "" {
		return fmt.Errorf("audio file is required")
	}

	// ffmpeg -i video.mp4 -i audio.mp3 -c:v copy -map 0:v:0 -map 1:a:0 -shortest output.mp4
	args := []string{
		"-i", opts.Input,
		"-i", opts.AudioFile,
		"-c:v", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-shortest",
	}
	if opts.Force {
		args = append(args, "-y")
	}
	args = append(args, output)

	fmt.Printf("Replacing audio %s + %s -> %s\n", opts.Input, opts.AudioFile, output)
	return runFFmpeg(ctx, ffmpegPath, args)
}

// ProbeInfo contains video metadata from ffprobe.
type ProbeInfo struct {
	Width    int
	Height   int
	Duration string
	Codec    string
	Size     int64
}

// Probe gets video information using ffprobe.
func (c *Client) Probe(ctx context.Context, input string) (*ProbeInfo, error) {
	ffprobePath, err := c.ensureFFprobe(ctx)
	if err != nil {
		return nil, fmt.Errorf("ensure ffprobe: %w", err)
	}

	// Get video stream info
	// Output order matches -show_entries order: codec_name,width,height,duration
	args := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name,width,height,duration",
		"-of", "csv=p=0",
		input,
	}

	cmd := exec.CommandContext(ctx, ffprobePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w", err)
	}

	// Parse output: codec_name,width,height,duration
	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	info := &ProbeInfo{}

	if len(parts) >= 1 {
		info.Codec = parts[0]
	}
	if len(parts) >= 2 {
		info.Width, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		info.Height, _ = strconv.Atoi(parts[2])
	}
	if len(parts) >= 4 {
		info.Duration = parts[3]
	}

	// Get file size
	if fi, err := os.Stat(input); err == nil {
		info.Size = fi.Size()
	}

	return info, nil
}

// ensureFFprobe ensures ffprobe is available (downloaded with ffmpeg).
func (c *Client) ensureFFprobe(ctx context.Context) (string, error) {
	binDir := getBinDir()
	ffprobePath := filepath.Join(binDir, "ffprobe"+exeExt())
	if _, err := os.Stat(ffprobePath); err == nil {
		return ffprobePath, nil
	}

	// Check if ffprobe is in PATH
	if path, err := exec.LookPath("ffprobe" + exeExt()); err == nil {
		return path, nil
	}

	// FFprobe should have been downloaded with ffmpeg
	// Try to download ffmpeg (which includes ffprobe)
	if err := downloadFFmpeg(ctx, binDir); err != nil {
		return "", fmt.Errorf("download ffmpeg/ffprobe: %w", err)
	}

	if _, err := os.Stat(ffprobePath); err == nil {
		return ffprobePath, nil
	}

	return "", fmt.Errorf("ffprobe not found after download")
}

// runFFmpeg executes ffmpeg with the given arguments.
func runFFmpeg(ctx context.Context, ffmpegPath string, args []string) error {
	cmd := exec.CommandContext(ctx, ffmpegPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// addSuffix adds a suffix before the file extension.
func addSuffix(path, suffix string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	return base + suffix + ext
}

// checkOverwrite checks if output file exists and returns error if Force is false.
func checkOverwrite(output string, force bool) error {
	if _, err := os.Stat(output); err == nil && !force {
		return fmt.Errorf("output file exists: %s (use -force to overwrite)", output)
	}
	return nil
}
