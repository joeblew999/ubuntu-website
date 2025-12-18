package youtube

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/ulikunitz/xz"
)

// defaultOutputTemplate uses yt-dlp template syntax with sanitization:
// - %(title)s gets the video title
// We post-process the filename to: lowercase, replace spaces with hyphens, remove special chars
const defaultOutputTemplate = "%(title)s.%(ext)s"

// Client wraps yt-dlp with a small, opinionated API.
type Client struct {
	proxy       string
	autoInstall bool
}

// NewClient returns a Client configured for proxy use and binary management.
func NewClient(proxy string, autoInstall bool) *Client {
	return &Client{
		proxy:       proxy,
		autoInstall: autoInstall,
	}
}

// DownloadOptions controls how a download is performed.
type DownloadOptions struct {
	Output    string
	Format    string
	Quality   string
	AudioOnly bool
	Force     bool
}

// DownloadResult contains the download summary.
type DownloadResult struct {
	Filename string
	Info     *InfoResult
}

// VideoManifest is saved as a JSON sidecar file for each downloaded video.
// Enables re-downloading with the same settings.
type VideoManifest struct {
	SourceURL    string    `json:"source_url"`
	VideoID      string    `json:"video_id"`
	Title        string    `json:"title"`
	Uploader     string    `json:"uploader,omitempty"`
	Quality      string    `json:"quality,omitempty"`
	AudioOnly    bool      `json:"audio_only,omitempty"`
	DownloadedAt time.Time `json:"downloaded_at"`
	Filename     string    `json:"filename"`
}

// InfoResult is a simplified view of yt-dlp metadata.
type InfoResult struct {
	ID        string
	Title     string
	URL       string
	Uploader  string
	Filename  string
	Duration  time.Duration
	Formats   []FormatInfo
	Extractor string
}

// FormatInfo holds the minimal format data we want to display.
type FormatInfo struct {
	ID         string
	Ext        string
	Resolution string
	FPS        string
	VCodec     string
	ACodec     string
	Note       string
	Size       string
}

// exeExt returns ".exe" on Windows, empty string otherwise.
func exeExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// Info fetches metadata without downloading content.
func (c *Client) Info(ctx context.Context, url string) (*InfoResult, error) {
	ytdlpPath, err := c.ensureYtdlp(ctx)
	if err != nil {
		return nil, fmt.Errorf("ensure yt-dlp: %w", err)
	}

	args := []string{"--no-progress", "--dump-single-json", "--skip-download"}
	if c.proxy != "" {
		args = append(args, "--proxy", c.proxy)
	}
	args = append(args, url)

	cmd := exec.CommandContext(ctx, ytdlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("youtube info: %w", err)
	}

	var raw rawInfo
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse info: %w", err)
	}

	return simplifyInfo(&raw), nil
}

// Download downloads a video (or audio) using yt-dlp.
// If a manifest for this URL already exists and the video file is present, download is skipped.
func (c *Client) Download(ctx context.Context, url string, opts DownloadOptions) (*DownloadResult, error) {
	// Check if video already exists by looking for matching manifest in output directory
	if !opts.Force && opts.Output != "" {
		outputDir := filepath.Dir(opts.Output)
		if existing := findExistingDownload(outputDir, url); existing != nil {
			fmt.Printf("Already downloaded: %s\n", existing.Filename)
			fmt.Printf("  (use -force to re-download)\n")
			return &DownloadResult{
				Filename: filepath.Join(outputDir, existing.Filename),
				Info: &InfoResult{
					ID:       existing.VideoID,
					Title:    existing.Title,
					URL:      existing.SourceURL,
					Uploader: existing.Uploader,
					Filename: existing.Filename,
				},
			}, nil
		}
	}

	ytdlpPath, err := c.ensureYtdlp(ctx)
	if err != nil {
		return nil, fmt.Errorf("ensure yt-dlp: %w", err)
	}

	// Ensure ffmpeg is available for merging video+audio
	ffmpegPath, err := c.ensureFFmpeg(ctx)
	if err != nil {
		return nil, fmt.Errorf("ensure ffmpeg: %w", err)
	}

	format := buildFormatString(opts.Format, opts.Quality, opts.AudioOnly)
	output := opts.Output
	if output == "" {
		output = defaultOutputTemplate
	}

	args := []string{
		"--print-json",
		"-f", format,
		"-o", output,
		"--ffmpeg-location", filepath.Dir(ffmpegPath),
		"--merge-output-format", "mp4",
	}
	if c.proxy != "" {
		args = append(args, "--proxy", c.proxy)
	}
	if opts.Force {
		args = append(args, "--force-overwrites")
	} else {
		args = append(args, "--no-overwrites")
	}
	args = append(args, url)

	cmd := exec.CommandContext(ctx, ytdlpPath, args...)
	cmdOutput, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("youtube download: %w\nstderr: %s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("youtube download: %w", err)
	}

	var raw rawInfo
	if err := json.Unmarshal(cmdOutput, &raw); err != nil {
		// Download may have succeeded without JSON output
		return &DownloadResult{Filename: output}, nil
	}

	info := simplifyInfo(&raw)
	filename := info.Filename
	if filename == "" {
		filename = output
	}

	// Rename to sanitized filename
	if filename != "" && filename != output {
		sanitizedName := sanitizeFilename(filepath.Base(filename))
		sanitizedPath := filepath.Join(filepath.Dir(filename), sanitizedName)

		if sanitizedPath != filename {
			if err := os.Rename(filename, sanitizedPath); err != nil {
				// Log but don't fail - file was downloaded successfully
				fmt.Printf("Warning: could not rename to %s: %v\n", sanitizedName, err)
			} else {
				filename = sanitizedPath
			}
		}
	}

	// Save manifest JSON sidecar file
	if filename != "" && info != nil {
		manifest := VideoManifest{
			SourceURL:    url,
			VideoID:      info.ID,
			Title:        info.Title,
			Uploader:     info.Uploader,
			Quality:      opts.Quality,
			AudioOnly:    opts.AudioOnly,
			DownloadedAt: time.Now(),
			Filename:     filepath.Base(filename),
		}
		manifestPath := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".json"
		if manifestData, err := json.MarshalIndent(manifest, "", "  "); err == nil {
			if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
				fmt.Printf("Warning: could not save manifest: %v\n", err)
			}
		}
	}

	return &DownloadResult{
		Filename: filename,
		Info:     info,
	}, nil
}

// getBinDir returns the directory for downloaded binaries.
// Uses .build/ in the current working directory to keep things project-local.
func getBinDir() string {
	// Try to find project root by looking for go.mod
	cwd, err := os.Getwd()
	if err != nil {
		// Fallback to .build in current dir
		return ".build"
	}

	// Walk up to find go.mod (project root)
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, ".build")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, use cwd
			return filepath.Join(cwd, ".build")
		}
		dir = parent
	}
}

// ensureYtdlp ensures yt-dlp is available, downloading if necessary.
func (c *Client) ensureYtdlp(ctx context.Context) (string, error) {
	// Check in .build/ first (project-local)
	binDir := getBinDir()
	ytdlpPath := filepath.Join(binDir, "yt-dlp"+exeExt())
	if _, err := os.Stat(ytdlpPath); err == nil {
		return ytdlpPath, nil
	}

	// Check if yt-dlp is in PATH (system-installed)
	if path, err := exec.LookPath("yt-dlp" + exeExt()); err == nil {
		return path, nil
	}

	if !c.autoInstall {
		return "", errors.New("yt-dlp not found (use -no-install=false to auto-download)")
	}

	// Download from GitHub releases to .build/
	if err := downloadYtdlp(ctx, ytdlpPath); err != nil {
		return "", err
	}

	return ytdlpPath, nil
}

// ensureFFmpeg ensures ffmpeg is available, downloading if necessary.
func (c *Client) ensureFFmpeg(ctx context.Context) (string, error) {
	// Check in .build/ first (project-local)
	binDir := getBinDir()
	ffmpegPath := filepath.Join(binDir, "ffmpeg"+exeExt())
	if _, err := os.Stat(ffmpegPath); err == nil {
		return ffmpegPath, nil
	}

	// Check if ffmpeg is in PATH (system-installed)
	if path, err := exec.LookPath("ffmpeg" + exeExt()); err == nil {
		return path, nil
	}

	if !c.autoInstall {
		return "", errors.New("ffmpeg not found (use -no-install=false to auto-download)")
	}

	// Download ffmpeg to .build/
	if err := downloadFFmpeg(ctx, binDir); err != nil {
		return "", err
	}

	return ffmpegPath, nil
}

// getYtdlpURL returns the download URL for yt-dlp based on OS/arch.
func getYtdlpURL() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// yt-dlp release URLs: https://github.com/yt-dlp/yt-dlp/releases
	switch {
	case goos == "darwin":
		// Universal macOS binary (works on both Intel and Apple Silicon)
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos", nil
	case goos == "linux" && goarch == "amd64":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux", nil
	case goos == "linux" && goarch == "arm64":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64", nil
	case goos == "windows" && goarch == "amd64":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe", nil
	case goos == "windows" && goarch == "386":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_x86.exe", nil
	default:
		return "", fmt.Errorf("unsupported platform: %s/%s - please install yt-dlp manually", goos, goarch)
	}
}

func downloadYtdlp(ctx context.Context, destPath string) error {
	url, err := getYtdlpURL()
	if err != nil {
		return err
	}

	fmt.Println("Downloading yt-dlp from", url)

	if err := downloadBinary(ctx, url, destPath); err != nil {
		return fmt.Errorf("download yt-dlp: %w", err)
	}

	fmt.Println("yt-dlp installed to", destPath)
	return nil
}

// ffmpegURLs holds download URLs for ffmpeg and optionally ffprobe (if separate).
type ffmpegURLs struct {
	ffmpeg  string
	ffprobe string // Only set for macOS where they're separate downloads
}

// getFFmpegURLs returns the download URL(s) for ffmpeg/ffprobe based on OS/arch.
func getFFmpegURLs() (*ffmpegURLs, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch {
	case goos == "darwin" && goarch == "arm64":
		// Martin Riedl's FFmpeg Build Server - macOS ARM64 (separate ffmpeg/ffprobe downloads)
		return &ffmpegURLs{
			ffmpeg:  "https://ffmpeg.martin-riedl.de/redirect/latest/macos/arm64/release/ffmpeg.zip",
			ffprobe: "https://ffmpeg.martin-riedl.de/redirect/latest/macos/arm64/release/ffprobe.zip",
		}, nil
	case goos == "darwin" && goarch == "amd64":
		// Martin Riedl's FFmpeg Build Server - macOS Intel (separate ffmpeg/ffprobe downloads)
		return &ffmpegURLs{
			ffmpeg:  "https://ffmpeg.martin-riedl.de/redirect/latest/macos/amd64/release/ffmpeg.zip",
			ffprobe: "https://ffmpeg.martin-riedl.de/redirect/latest/macos/amd64/release/ffprobe.zip",
		}, nil
	case goos == "linux" && goarch == "amd64":
		// yt-dlp FFmpeg builds - Linux x64 (includes ffprobe in archive)
		return &ffmpegURLs{ffmpeg: "https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz"}, nil
	case goos == "linux" && goarch == "arm64":
		// yt-dlp FFmpeg builds - Linux ARM64 (includes ffprobe in archive)
		return &ffmpegURLs{ffmpeg: "https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linuxarm64-gpl.tar.xz"}, nil
	case goos == "windows" && goarch == "amd64":
		// yt-dlp FFmpeg builds - Windows x64 (includes ffprobe in archive)
		return &ffmpegURLs{ffmpeg: "https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"}, nil
	case goos == "windows" && goarch == "arm64":
		// yt-dlp FFmpeg builds - Windows ARM64 (includes ffprobe in archive)
		return &ffmpegURLs{ffmpeg: "https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-winarm64-gpl.zip"}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s/%s - please install ffmpeg manually", goos, goarch)
	}
}

func downloadFFmpeg(ctx context.Context, binDir string) error {
	urls, err := getFFmpegURLs()
	if err != nil {
		return err
	}

	// Ensure bin directory exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}

	// Download ffmpeg
	if err := downloadAndExtractBinary(ctx, urls.ffmpeg, binDir, "ffmpeg"); err != nil {
		return err
	}

	// Download ffprobe separately if URL provided (macOS case)
	if urls.ffprobe != "" {
		if err := downloadAndExtractBinary(ctx, urls.ffprobe, binDir, "ffprobe"); err != nil {
			return err
		}
	}

	// Remove quarantine attribute on macOS (required for unsigned binaries)
	if runtime.GOOS == "darwin" {
		_ = exec.Command("xattr", "-dr", "com.apple.quarantine", filepath.Join(binDir, "ffmpeg")).Run()
		_ = exec.Command("xattr", "-dr", "com.apple.quarantine", filepath.Join(binDir, "ffprobe")).Run()
	}

	return nil
}

// downloadAndExtractBinary downloads an archive and extracts the specified binary.
func downloadAndExtractBinary(ctx context.Context, url, binDir, binaryName string) error {
	fmt.Printf("Downloading %s from %s\n", binaryName, url)

	// Download the archive
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download %s: %w", binaryName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: status %d", binaryName, resp.StatusCode)
	}

	// Read entire response into memory
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read %s archive: %w", binaryName, err)
	}

	// Extract based on archive type
	if strings.HasSuffix(url, ".zip") {
		if err := extractBinaryFromZip(data, binDir, binaryName); err != nil {
			return err
		}
	} else if strings.HasSuffix(url, ".tar.xz") {
		if err := extractFFmpegFromTarXz(data, binDir); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported archive format: %s", url)
	}

	return nil
}

// extractBinaryFromZip extracts a specific binary from a zip archive.
func extractBinaryFromZip(data []byte, binDir, binaryName string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}

	ext := exeExt()
	found := false

	for _, f := range zipReader.File {
		name := filepath.Base(f.Name)

		// Match the target binary
		if name == binaryName+ext || name == binaryName {
			destPath := filepath.Join(binDir, binaryName+ext)
			if err := extractZipFile(f, destPath); err != nil {
				return fmt.Errorf("extract %s: %w", binaryName, err)
			}
			fmt.Printf("%s installed to %s\n", binaryName, destPath)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%s not found in zip archive", binaryName)
	}

	return nil
}

func extractFFmpegFromTarXz(data []byte, binDir string) error {
	// Decompress xz
	xzReader, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("open xz: %w", err)
	}

	// Read tar archive
	tarReader := tar.NewReader(xzReader)

	ext := exeExt()
	ffmpegFound := false
	ffprobeFound := false

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		// Skip non-regular files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		name := filepath.Base(header.Name)

		// Match ffmpeg binary (in bin/ subdirectory typically)
		if name == "ffmpeg"+ext || name == "ffmpeg" {
			destPath := filepath.Join(binDir, "ffmpeg"+ext)
			if err := extractTarFile(tarReader, destPath, header.Mode); err != nil {
				return fmt.Errorf("extract ffmpeg: %w", err)
			}
			fmt.Println("ffmpeg installed to", destPath)
			ffmpegFound = true
		}

		// Match ffprobe binary
		if name == "ffprobe"+ext || name == "ffprobe" {
			destPath := filepath.Join(binDir, "ffprobe"+ext)
			if err := extractTarFile(tarReader, destPath, header.Mode); err != nil {
				return fmt.Errorf("extract ffprobe: %w", err)
			}
			fmt.Println("ffprobe installed to", destPath)
			ffprobeFound = true
		}
	}

	if !ffmpegFound {
		return errors.New("ffmpeg not found in tar.xz archive")
	}
	if !ffprobeFound {
		fmt.Println("Warning: ffprobe not found in tar.xz archive")
	}

	return nil
}

func extractZipFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

func extractTarFile(r io.Reader, destPath string, mode int64) error {
	// Use mode from tar header, but ensure executable
	fileMode := os.FileMode(mode)
	if fileMode == 0 {
		fileMode = 0755
	}

	outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, r)
	return err
}

// extractFile is kept for backwards compatibility but uses extractZipFile internally
func extractFile(f *zip.File, destPath string) error {
	return extractZipFile(f, destPath)
}

func downloadBinary(ctx context.Context, url, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	// Ensure bin directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}

	// Write to file
	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// rawInfo is the JSON structure from yt-dlp --dump-single-json
type rawInfo struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	WebpageURL string      `json:"webpage_url"`
	Uploader   string      `json:"uploader"`
	Filename   string      `json:"filename"`
	Duration   float64     `json:"duration"`
	Extractor  string      `json:"extractor"`
	Formats    []rawFormat `json:"formats"`
}

type rawFormat struct {
	FormatID       string  `json:"format_id"`
	Ext            string  `json:"ext"`
	Resolution     string  `json:"resolution"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	FPS            float64 `json:"fps"`
	VCodec         string  `json:"vcodec"`
	ACodec         string  `json:"acodec"`
	FormatNote     string  `json:"format_note"`
	FileSize       int64   `json:"filesize"`
	FileSizeApprox int64   `json:"filesize_approx"`
}

func buildFormatString(format, quality string, audioOnly bool) string {
	if format != "" {
		return format
	}
	if audioOnly {
		return "bestaudio/best"
	}
	// With ffmpeg available, we can merge best video + best audio
	if strings.TrimSpace(quality) != "" {
		q := strings.TrimSpace(quality)
		return fmt.Sprintf("bestvideo[height<=?%s]+bestaudio/best[height<=?%s]/best", q, q)
	}
	return "bestvideo+bestaudio/best"
}

func simplifyInfo(raw *rawInfo) *InfoResult {
	if raw == nil {
		return nil
	}

	result := &InfoResult{
		ID:        raw.ID,
		Title:     raw.Title,
		URL:       raw.WebpageURL,
		Uploader:  raw.Uploader,
		Filename:  raw.Filename,
		Extractor: raw.Extractor,
		Duration:  time.Duration(raw.Duration * float64(time.Second)),
	}

	for _, f := range raw.Formats {
		result.Formats = append(result.Formats, simplifyFormat(&f))
	}

	return result
}

func simplifyFormat(f *rawFormat) FormatInfo {
	if f == nil {
		return FormatInfo{}
	}

	res := f.Resolution
	if res == "" {
		if f.Width > 0 && f.Height > 0 {
			res = fmt.Sprintf("%dx%d", f.Width, f.Height)
		} else if f.Height > 0 {
			res = fmt.Sprintf("%dp", f.Height)
		}
	}

	fpsStr := ""
	if f.FPS > 0 {
		if math.Mod(f.FPS, 1) == 0 {
			fpsStr = fmt.Sprintf("%.0ffps", f.FPS)
		} else {
			fpsStr = fmt.Sprintf("%.1ffps", f.FPS)
		}
	}

	size := ""
	if f.FileSize > 0 {
		size = humanBytes(f.FileSize)
	} else if f.FileSizeApprox > 0 {
		size = "~" + humanBytes(f.FileSizeApprox)
	}

	return FormatInfo{
		ID:         f.FormatID,
		Ext:        f.Ext,
		Resolution: res,
		FPS:        fpsStr,
		VCodec:     f.VCodec,
		ACodec:     f.ACodec,
		Note:       f.FormatNote,
		Size:       size,
	}
}

func humanBytes(size int64) string {
	const unit = 1024
	value := float64(size)
	suffix := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	idx := 0
	for value >= unit && idx < len(suffix)-1 {
		value /= unit
		idx++
	}

	if idx == 0 {
		return fmt.Sprintf("%.0f%s", value, suffix[idx])
	}

	return fmt.Sprintf("%.1f%s", value, suffix[idx])
}

// Unused import placeholder for gzip (may be needed for future .tar.gz support)
var _ = gzip.NewReader

// LoadManifest reads a video manifest from a JSON file.
func LoadManifest(path string) (*VideoManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m VideoManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ListManifests finds all .json manifest files in a directory.
func ListManifests(dir string) ([]*VideoManifest, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var manifests []*VideoManifest
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		m, err := LoadManifest(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue // Skip invalid manifests
		}
		manifests = append(manifests, m)
	}
	return manifests, nil
}

// findExistingDownload checks if a video for the given URL already exists in the directory.
// Returns the manifest if both manifest and video file exist, nil otherwise.
func findExistingDownload(dir, url string) *VideoManifest {
	manifests, err := ListManifests(dir)
	if err != nil {
		return nil
	}

	for _, m := range manifests {
		if m.SourceURL == url {
			// Check if the video file also exists
			videoPath := filepath.Join(dir, m.Filename)
			if _, err := os.Stat(videoPath); err == nil {
				return m
			}
		}
	}
	return nil
}

// RefreshAll re-downloads all videos in a directory based on their manifests.
func (c *Client) RefreshAll(ctx context.Context, dir string) error {
	manifests, err := ListManifests(dir)
	if err != nil {
		return fmt.Errorf("list manifests: %w", err)
	}

	for _, m := range manifests {
		fmt.Printf("Refreshing: %s\n", m.Title)
		outputTemplate := filepath.Join(dir, "%(title)s.%(ext)s")
		_, err := c.Download(ctx, m.SourceURL, DownloadOptions{
			Output:    outputTemplate,
			Quality:   m.Quality,
			AudioOnly: m.AudioOnly,
			Force:     true,
		})
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		fmt.Printf("  Done\n")
	}
	return nil
}

// sanitizeFilename converts a filename to a clean, URL-friendly format:
// - lowercase
// - spaces and underscores become hyphens
// - keep only ASCII letters, digits, and hyphens
// - collapse multiple hyphens
// - trim leading/trailing hyphens
func sanitizeFilename(name string) string {
	// Separate extension
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	// Lowercase
	base = strings.ToLower(base)

	// Replace spaces and underscores with hyphens
	base = strings.ReplaceAll(base, " ", "-")
	base = strings.ReplaceAll(base, "_", "-")

	// Remove non-ASCII and special characters (keep only a-z, 0-9, and hyphens)
	var result strings.Builder
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	base = result.String()

	// Collapse multiple hyphens
	multiHyphen := regexp.MustCompile(`-+`)
	base = multiHyphen.ReplaceAllString(base, "-")

	// Trim leading/trailing hyphens
	base = strings.Trim(base, "-")

	// Ensure we have something
	if base == "" {
		base = "video"
	}

	return base + strings.ToLower(ext)
}
