// Package genlogo generates Ubuntu Software logo assets using Go graphics.
//
// SINGLE SOURCE OF TRUTH for all logo/branding assets.
//
// This package generates all brand assets from a single set of inputs (brand symbol,
// company name, colors, fonts). When branding needs to change, update the BRAND
// INPUTS section and regenerate all assets with: task genlogo:all
//
// Output types:
//   - SVG: For Hugo website (scalable, crisp text at any resolution)
//   - PNG: For social media, email, favicon (raster format required by platforms)
//
// Asset destinations:
//   - assets/images/: Hugo-processed assets (optimized during build)
//   - static/images/: Served directly without processing (email, social)
package genlogo

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

// =============================================================================
// OUTPUT PATHS - All generated assets go to these locations
// =============================================================================
//
// Why two directories?
// - assets/images/: Hugo processes these (optimization, fingerprinting, etc.)
//   Used for website logos, favicon, OG image - things Hugo needs to reference.
// - static/images/: Served as-is without Hugo processing.
//   Used for external platforms (email clients, Bluesky) that need stable URLs.
const (
	// Hugo website assets (processed by Hugo)
	pathFavicon     = "assets/images/favicon.png"       // 512x512 - browser tab icon, large for high-DPI
	pathLogoSVG     = "assets/images/logo.svg"          // SVG - Hugo header (light mode backgrounds)
	pathLogoDarkSVG = "assets/images/logo-darkmode.svg" // SVG - Hugo header (dark mode backgrounds)
	pathOGImage     = "static/images/og-image.png"      // 1200x630 - social sharing (Twitter, LinkedIn, etc.)

	// Static files (served directly, not processed by Hugo)
	pathEmailLogo     = "static/images/email-logo.png"     // 300x47 - Gmail signature (must be externally accessible)
	pathBlueskyAvatar = "static/images/bluesky-avatar.png" // 400x400 - Bluesky profile picture (square)
	pathBlueskyBanner = "static/images/bluesky-banner.png" // 1500x500 - Bluesky header banner (3:1 ratio)
)

// =============================================================================
// BRAND INPUTS - Change these to update all assets
// =============================================================================
//
// To rebrand: Update these values and run `task genlogo:all`
// All generated assets will automatically use the new branding.
const (
	// Brand identity
	brandSymbol = "[U|S]"           // Logo mark - appears in favicon, avatar, and before company name
	brandName   = "Ubuntu Software" // Company name - appears alongside logo mark in headers/banners

	// Background colors
	bgColorDark  = "#1a1a2e" // Dark navy - favicon, avatar (icon-only assets)
	bgColorLight = "#FFFFFF" // White - banner, OG image (full logo assets)

	// Text colors
	textColor     = "#58a6ff" // Ubuntu blue - [U|S] logo mark (all assets)
	textColorDark = "#121212" // Near-black - "Ubuntu Software" on light backgrounds

	// Fonts (macOS system fonts - these are available on all Macs)
	// Note: On Linux/Windows, these would need to be changed or fonts bundled
	fontMono = "/System/Library/Fonts/Menlo.ttc"     // Monospace for [U|S] - gives it a "code" feel
	fontSans = "/System/Library/Fonts/Helvetica.ttc" // Sans-serif for company name - clean, professional
)

// Run executes the genlogo CLI with the given arguments.
// Returns exit code: 0 for success, non-zero for errors.
func Run(args []string, version string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("genlogo", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		outputDir string
		asset     string
		ver       bool
	)
	fs.StringVar(&outputDir, "dir", ".", "Output directory (project root)")
	fs.StringVar(&asset, "asset", "all", "Asset to generate: avatar, banner, favicon, logo-svg, email, og, or all")
	fs.BoolVar(&ver, "version", false, "Print version and exit")

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if ver {
		fmt.Fprintf(stdout, "genlogo %s\n", version)
		return 0
	}

	gen := &Generator{
		outputDir: outputDir,
		stdout:    stdout,
		stderr:    stderr,
	}

	switch asset {
	case "avatar":
		gen.generateAvatar()
	case "banner":
		gen.generateBanner()
	case "favicon":
		gen.generateFavicon()
	case "logo-svg":
		gen.generateLogoSVG(false)
		gen.generateLogoSVG(true)
	case "email":
		gen.generateEmailLogo()
	case "og":
		gen.generateOGImage()
	case "all":
		gen.generateFavicon()
		gen.generateLogoSVG(false)
		gen.generateLogoSVG(true)
		gen.generateEmailLogo()
		gen.generateOGImage()
		gen.generateAvatar()
		gen.generateBanner()
	default:
		fmt.Fprintf(stderr, "Unknown asset: %s\nValid: avatar, banner, favicon, logo-svg, email, og, all\n", asset)
		return 1
	}

	return 0
}

// Generator holds configuration for asset generation
type Generator struct {
	outputDir string
	stdout    io.Writer
	stderr    io.Writer
}

// generateAvatar creates a 400x400 square avatar for Bluesky profile picture.
// Shows only the logo mark [U|S] centered on dark background with rounded corners.
// Bluesky crops to circle, so we use rounded rectangle that looks good either way.
func (g *Generator) generateAvatar() {
	size := 400
	dc := gg.NewContext(size, size)

	// Background
	dc.SetHexColor(bgColorDark)
	dc.DrawRoundedRectangle(0, 0, float64(size), float64(size), 40)
	dc.Fill()

	// Load font
	if err := dc.LoadFontFace(fontMono, 140); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}

	// Draw logo mark
	dc.SetHexColor(textColor)
	w, h := dc.MeasureString(brandSymbol)
	x := (float64(size) - w) / 2
	y := (float64(size) + h) / 2
	dc.DrawString(brandSymbol, x, y)

	g.savePNG(dc, filepath.Join(g.outputDir, pathBlueskyAvatar), "400x400")
}

// generateBanner creates a 1500x500 banner for Bluesky profile header.
// Shows logo mark [U|S] + company name "Ubuntu Software" centered horizontally.
// Uses WHITE background with blue [U|S] and black "Ubuntu Software" - matches website/email branding.
// Bluesky requires 3:1 aspect ratio (1500x500).
func (g *Generator) generateBanner() {
	width, height := 1500, 500
	dc := gg.NewContext(width, height)

	// White background (matches email logo and website light mode)
	dc.SetHexColor(bgColorLight)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// Load font for logo mark
	if err := dc.LoadFontFace(fontMono, 100); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}

	// Measure logo mark
	dc.SetHexColor(textColor)
	symbolW, h := dc.MeasureString(brandSymbol)

	// Load font for company name
	if err := dc.LoadFontFace(fontSans, 60); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load Helvetica: %v\n", err)
	}

	// Measure company name
	nameW, _ := dc.MeasureString(brandName)

	// Calculate total width and center everything
	gap := 30.0 // gap between symbol and name
	totalW := symbolW + gap + nameW
	startX := (float64(width) - totalW) / 2
	y := (float64(height) + h) / 2

	// Draw logo mark in BLUE
	if err := dc.LoadFontFace(fontMono, 100); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}
	dc.SetHexColor(textColor) // blue #58a6ff
	dc.DrawString(brandSymbol, startX, y)

	// Draw company name in BLACK
	if err := dc.LoadFontFace(fontSans, 60); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load Helvetica: %v\n", err)
	}
	dc.SetHexColor(textColorDark) // black #121212
	dc.DrawString(brandName, startX+symbolW+gap, y)

	g.savePNG(dc, filepath.Join(g.outputDir, pathBlueskyBanner), "1500x500")
}

// generateFavicon creates a 512x512 favicon for browser tabs.
// Large size (512x512) ensures crisp display on high-DPI screens.
// Shows only logo mark [U|S] - company name would be too small to read.
// Browsers will scale down as needed for smaller contexts.
func (g *Generator) generateFavicon() {
	size := 512
	dc := gg.NewContext(size, size)

	// Background
	dc.SetHexColor(bgColorDark)
	dc.DrawRoundedRectangle(0, 0, float64(size), float64(size), 80)
	dc.Fill()

	// Load font
	if err := dc.LoadFontFace(fontMono, 180); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}

	// Draw logo mark
	dc.SetHexColor(textColor)
	w, h := dc.MeasureString(brandSymbol)
	x := (float64(size) - w) / 2
	y := (float64(size) + h) / 2
	dc.DrawString(brandSymbol, x, y)

	g.savePNG(dc, filepath.Join(g.outputDir, pathFavicon), "512x512")
}

// generateLogoSVG creates SVG logo for Hugo website header.
// SVG is preferred for web because it scales perfectly at any resolution.
// Two variants:
//   - Light mode (darkMode=false): Dark text for light backgrounds
//   - Dark mode (darkMode=true): Light text for dark backgrounds
//
// The logo mark [U|S] is always Ubuntu blue; only company name color changes.
func (g *Generator) generateLogoSVG(darkMode bool) {
	width, height := 265, 50 // Tight viewBox - matches 220px mobile constraint well

	// Text color depends on mode
	nameColor := textColorDark
	if darkMode {
		nameColor = textColor
	}

	// No width/height attributes - viewBox only for responsive scaling on mobile
	// Larger font sizes than original for better visibility
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d">
  <text x="0" y="36" font-family="Menlo, Monaco, 'Courier New', monospace" font-size="28" font-weight="bold" fill="%s">%s</text>
  <text x="95" y="36" font-family="Helvetica, Arial, sans-serif" font-size="23" fill="%s">%s</text>
</svg>`, width, height, textColor, brandSymbol, nameColor, brandName)

	var outPath string
	if darkMode {
		outPath = filepath.Join(g.outputDir, pathLogoDarkSVG)
	} else {
		outPath = filepath.Join(g.outputDir, pathLogoSVG)
	}
	g.saveSVG(svg, outPath)
}

// generateEmailLogo creates a logo optimized for Gmail signatures.
// Slightly smaller than website logo (300x47) to fit well in email.
// Uses dark text since most email clients have white backgrounds.
// Must be externally accessible URL for Gmail to display it.
func (g *Generator) generateEmailLogo() {
	width, height := 300, 47
	dc := gg.NewContext(width, height)

	// Transparent background
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	// Load font for [U|S]
	if err := dc.LoadFontFace(fontMono, 22); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}

	// Draw logo mark
	dc.SetHexColor(textColor)
	dc.DrawString(brandSymbol, 0, 30)

	// Load font for company name
	if err := dc.LoadFontFace(fontSans, 18); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load Helvetica: %v\n", err)
	}

	// Draw company name (dark for email - most email backgrounds are white)
	dc.SetHexColor(textColorDark)
	dc.DrawString(brandName, 75, 30)

	g.savePNG(dc, filepath.Join(g.outputDir, pathEmailLogo), "300x47")
}

// generateOGImage creates a 1200x630 Open Graph image for social sharing.
// This appears as the preview image when the website URL is shared on
// Twitter, LinkedIn, Facebook, Slack, etc.
// Uses HORIZONTAL layout with "[U|S] Ubuntu Software" on single line to avoid
// cropping issues when messaging apps display previews in narrower aspect ratios.
// WHITE background with blue [U|S] and black company name.
func (g *Generator) generateOGImage() {
	width, height := 1200, 630
	dc := gg.NewContext(width, height)

	// White background (matches website/email/banner branding)
	dc.SetHexColor(bgColorLight)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// Load font for logo mark - smaller to fit horizontal layout
	if err := dc.LoadFontFace(fontMono, 80); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}

	// Measure logo mark
	dc.SetHexColor(textColor)
	symbolW, h := dc.MeasureString(brandSymbol)

	// Load font for company name
	if err := dc.LoadFontFace(fontSans, 50); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load Helvetica: %v\n", err)
	}

	// Measure company name
	nameW, _ := dc.MeasureString(brandName)

	// Calculate total width and center everything horizontally
	gap := 25.0 // gap between symbol and name
	totalW := symbolW + gap + nameW
	startX := (float64(width) - totalW) / 2
	y := (float64(height) + h) / 2

	// Draw logo mark in BLUE
	if err := dc.LoadFontFace(fontMono, 80); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load font: %v\n", err)
	}
	dc.SetHexColor(textColor) // blue #58a6ff
	dc.DrawString(brandSymbol, startX, y)

	// Draw company name in BLACK
	if err := dc.LoadFontFace(fontSans, 50); err != nil {
		fmt.Fprintf(g.stderr, "Warning: Could not load Helvetica: %v\n", err)
	}
	dc.SetHexColor(textColorDark) // black #121212
	dc.DrawString(brandName, startX+symbolW+gap, y)

	g.savePNG(dc, filepath.Join(g.outputDir, pathOGImage), "1200x630")
}

// savePNG writes the graphics context to a PNG file, creating directories as needed.
func (g *Generator) savePNG(dc *gg.Context, path, size string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Fprintf(g.stderr, "Error creating directory: %v\n", err)
		return
	}
	if err := dc.SavePNG(path); err != nil {
		fmt.Fprintf(g.stderr, "Error saving PNG: %v\n", err)
		return
	}
	fmt.Fprintf(g.stdout, "Generated: %s (%s)\n", path, size)
}

// saveSVG writes SVG content to a file, creating directories as needed.
func (g *Generator) saveSVG(content, path string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Fprintf(g.stderr, "Error creating directory: %v\n", err)
		return
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(g.stderr, "Error saving SVG: %v\n", err)
		return
	}
	fmt.Fprintf(g.stdout, "Generated: %s (SVG)\n", path)
}
