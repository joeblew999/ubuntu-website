package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// Printf writes formatted output.
func (c *Context) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.config.Output, format, a...)
}

// Println writes a line.
func (c *Context) Println(a ...interface{}) {
	fmt.Fprintln(c.config.Output, a...)
}

// Errorf writes formatted error output.
func (c *Context) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(c.config.ErrOutput, format, a...)
}

// Errorln writes an error line.
func (c *Context) Errorln(a ...interface{}) {
	fmt.Fprintln(c.config.ErrOutput, a...)
}

// Header prints a header line.
func (c *Context) Header(title string) {
	if c.GitHubIssue() {
		c.Printf("## %s\n\n", title)
	} else {
		c.Println(title)
		c.Println(strings.Repeat("=", len(title)))
	}
}

// SubHeader prints a sub-header line.
func (c *Context) SubHeader(title string) {
	if c.GitHubIssue() {
		c.Printf("### %s\n\n", title)
	} else {
		c.Println(title)
		c.Println(strings.Repeat("-", len(title)))
	}
}

// Separator prints a separator line.
func (c *Context) Separator() {
	if !c.GitHubIssue() {
		c.Println(strings.Repeat("─", 60))
	}
}

// KeyValue prints a key-value pair.
func (c *Context) KeyValue(key string, value interface{}) {
	if c.GitHubIssue() {
		c.Printf("- **%s:** %v\n", key, value)
	} else {
		c.Printf("  %s: %v\n", key, value)
	}
}

// Table provides a simple table writer.
type Table struct {
	ctx     *Context
	w       *tabwriter.Writer
	headers []string
}

// NewTable creates a new table.
func (c *Context) NewTable(headers ...string) *Table {
	t := &Table{
		ctx:     c,
		headers: headers,
	}

	if c.GitHubIssue() {
		// Markdown table header
		c.Printf("| %s |\n", strings.Join(headers, " | "))
		c.Printf("|%s|\n", strings.Repeat("--------|", len(headers)))
	} else {
		t.w = tabwriter.NewWriter(c.config.Output, 0, 0, 2, ' ', 0)
		fmt.Fprintln(t.w, strings.Join(headers, "\t"))
	}

	return t
}

// Row adds a row to the table.
func (t *Table) Row(values ...interface{}) {
	strs := make([]string, len(values))
	for i, v := range values {
		strs[i] = fmt.Sprintf("%v", v)
	}

	if t.ctx.GitHubIssue() {
		t.ctx.Printf("| %s |\n", strings.Join(strs, " | "))
	} else {
		fmt.Fprintln(t.w, strings.Join(strs, "\t"))
	}
}

// Flush flushes the table output.
func (t *Table) Flush() {
	if t.w != nil {
		t.w.Flush()
	}
}

// List prints a bulleted list.
func (c *Context) List(items ...string) {
	for _, item := range items {
		if c.GitHubIssue() {
			c.Printf("- %s\n", item)
		} else {
			c.Printf("  • %s\n", item)
		}
	}
}

// Code prints a code block.
func (c *Context) Code(code string) {
	if c.GitHubIssue() {
		c.Println("```")
		c.Println(code)
		c.Println("```")
	} else {
		c.Printf("  %s\n", code)
	}
}

// CodeBlock prints a code block with language.
func (c *Context) CodeBlock(lang, code string) {
	if c.GitHubIssue() {
		c.Printf("```%s\n", lang)
		c.Println(code)
		c.Println("```")
	} else {
		c.Printf("  %s\n", code)
	}
}

// Success prints a success message.
func (c *Context) Success(msg string) {
	if c.GitHubIssue() {
		c.Printf("✅ %s\n", msg)
	} else {
		c.Printf("✓ %s\n", msg)
	}
}

// Warning prints a warning message.
func (c *Context) Warning(msg string) {
	if c.GitHubIssue() {
		c.Printf("⚠️ %s\n", msg)
	} else {
		c.Printf("Warning: %s\n", msg)
	}
}

// Error prints an error message.
func (c *Context) Error(msg string) {
	if c.GitHubIssue() {
		c.Errorf("❌ %s\n", msg)
	} else {
		c.Errorf("ERROR: %s\n", msg)
	}
}

// Link prints a link.
func (c *Context) Link(text, url string) {
	if c.GitHubIssue() {
		c.Printf("[%s](%s)", text, url)
	} else {
		c.Printf("%s: %s", text, url)
	}
}

// Writer returns a raw io.Writer for custom output.
func (c *Context) Writer() io.Writer {
	return c.config.Output
}
