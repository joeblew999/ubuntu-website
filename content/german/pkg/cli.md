---
title: "cli"
import_path: www.ubuntusoftware.net/pkg/cli
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Gemeinsames CLI-Framework für Ubuntu-Software-Tools. Standard-Flags, Ausgabeformatierung und GitHub-Issue-Markdown."
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/cli
license: MIT
author: Gerard Webb
created_at: 2025-12-16T00:00:00Z
updated_at: 2025-12-16T00:00:00Z
has_binary: false
---

## Funktionen

- **Standard-Flags** – `--github-issue`, `--verbose`, `--version` sofort einsatzbereit
- **Ausgabeformatierung** – Tabellen, Listen, Überschriften, Codeblöcke
- **Markdown-Modus** – Für GitHub-Issues geeignete Ausgabe mit `--github-issue`
- **Plattformübergreifend** – Öffnen des Browsers unter macOS, Linux und Windows
- **Kontextverwaltung** – Timeouts und Abbruch

## Verwendung

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

## Ausgabemodi

**Normale Ausgabe:**

```
Items
=====
NAME    VALUE
foo     bar
baz     qux
```

**Mit dem Flag `--github-issue`:**

```markdown
## Items

| NAME | VALUE |
|--------|--------|
| foo | bar |
| baz | qux |
```

## Kontextmethoden

| Methode | Beschreibung |
|--------|-------------|
| `c.Header(title)` | Abschnittsüberschrift ausgeben |
| `c.KeyValue(key, val)` | Schlüssel-Wert-Paar ausgeben |
| `c.NewTable(headers...)` | Tabelle erstellen |
| `c.List(items...)` | Aufzählungsliste ausgeben |
| `c.Code(code)` | Codeblock ausgeben |
| `c.Success(msg)` | Erfolgsmeldung ausgeben |
| `c.Warning(msg)` | Warnung ausgeben |
| `c.Open(url)` | URL im Browser öffnen |
