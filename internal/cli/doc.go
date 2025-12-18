// Package cli provides a shared CLI framework for Ubuntu Software tools.
//
// Features:
//   - Standard flags: --github-issue, --verbose, --version
//   - Consistent output formatting (tables, lists, code blocks)
//   - Markdown output for GitHub issues
//   - Context with timeout
//   - Cross-platform browser opening
//
// Example:
//
//	package main
//
//	import (
//	    "os"
//	    "www.ubuntusoftware.net/pkg/cli"
//	)
//
//	func main() {
//	    app := cli.New("myapp", "v1.0.0")
//
//	    err := app.Run(os.Args[1:], func(c *cli.Context) error {
//	        if len(c.Args) == 0 {
//	            return fmt.Errorf("command required")
//	        }
//
//	        switch c.Args[0] {
//	        case "list":
//	            c.Header("Items")
//	            table := c.NewTable("NAME", "VALUE")
//	            table.Row("foo", "bar")
//	            table.Row("baz", "qux")
//	            table.Flush()
//
//	        case "info":
//	            c.Header("Info")
//	            c.KeyValue("Version", "1.0.0")
//	            c.KeyValue("Author", "Gerard Webb")
//	        }
//
//	        return nil
//	    })
//
//	    if err != nil {
//	        fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
//	        os.Exit(1)
//	    }
//	}
//
// Output modes:
//
// Normal output:
//
//	Items
//	=====
//	NAME    VALUE
//	foo     bar
//	baz     qux
//
// With --github-issue flag:
//
//	## Items
//
//	| NAME | VALUE |
//	|--------|--------|
//	| foo | bar |
//	| baz | qux |
package cli
