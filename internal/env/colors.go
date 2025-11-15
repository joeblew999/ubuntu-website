package env

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"

	// Bold variants
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
)

// Colorize wraps text with color codes
func Colorize(text, color string) string {
	return color + text + ColorReset
}

// Success formats success message in green
func Success(text string) string {
	return Colorize("✓ "+text, ColorGreen)
}

// Error formats error message in red
func Error(text string) string {
	return Colorize("✗ "+text, ColorRed)
}

// Warning formats warning message in yellow
func Warning(text string) string {
	return Colorize("⚠ "+text, ColorYellow)
}

// Info formats info message in cyan
func Info(text string) string {
	return Colorize("ℹ "+text, ColorCyan)
}

// Skipped formats skipped message in gray
func Skipped(text string) string {
	return Colorize("⊘ "+text, ColorGray)
}
