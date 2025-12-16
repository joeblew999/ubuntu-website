---
title: cli
import_path: www.ubuntusoftware.net/pkg/cli
repo_url: https://github.com/joeblew999/ubuntu-website
description: Shared CLI framework for Ubuntu Software tools. Standard flags, output formatting, and GitHub issue markdown.
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/cli
license: MIT
author: Gerard Webb
created_at: 2025-12-16T00:00:00Z
updated_at: 2025-12-16T00:00:00Z
has_binary: false
---

## Features

- **Standard flags** - `--github-issue`, `--verbose`, `--version` out of the box
- **Output formatting** - Tables, lists, headers, code blocks
- **Markdown mode** - GitHub issue-friendly output with `--github-issue`
- **Cross-platform** - Browser opening for macOS, Linux, Windows
- **Context handling** - Timeouts and cancellation

## Usage

```go
package main

import (
    "fmt"
    "os"

    "www.ubuntusoftware.net/pkg/cli"
)

func main() {
    app := cli.New("myapp", "v1.0.0")

    err := app.Run(os.Args[1:], func(c *cli.Context) error {
        if len(c.Args) == 0 {
            return fmt.Errorf("command required")
        }

        switch c.Args[0] {
        case "list":
            c.Header("Items")
            table := c.NewTable("NAME", "VALUE")
            table.Row("foo", "bar")
            table.Row("baz", "qux")
            table.Flush()

        case "open":
            return c.Open("https://example.com")
        }

        return nil
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
        os.Exit(1)
    }
}
```

## Output Modes

**Normal output:**

```
Items
=====
NAME    VALUE
foo     bar
baz     qux
```

**With `--github-issue` flag:**

```markdown
## Items

| NAME | VALUE |
|--------|--------|
| foo | bar |
| baz | qux |
```

## Context Methods

| Method | Description |
|--------|-------------|
| `c.Header(title)` | Print section header |
| `c.KeyValue(key, val)` | Print key-value pair |
| `c.NewTable(headers...)` | Create table |
| `c.List(items...)` | Print bulleted list |
| `c.Code(code)` | Print code block |
| `c.Success(msg)` | Print success message |
| `c.Warning(msg)` | Print warning |
| `c.Open(url)` | Open URL in browser |
