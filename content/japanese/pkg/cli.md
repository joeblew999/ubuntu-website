---
title: "CLI"
import_path: www.ubuntusoftware.net/pkg/cli
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Ubuntu Software ツール向けの共有 CLI フレームワーク。標準フラグ、出力フォーマット、GitHub issue 用の Markdown を提供します。"
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/cli
license: MIT
author: Gerard Webb
created_at: 2025-12-16T00:00:00Z
updated_at: 2025-12-16T00:00:00Z
has_binary: false
---

## 機能

- **標準フラグ** - `--github-issue`、`--verbose`、`--version` をすぐに利用可能
- **出力フォーマット** - テーブル、リスト、ヘッダー、コードブロック
- **Markdown モード** - `--github-issue` による GitHub Issue に適した出力
- **クロスプラットフォーム** - macOS、Linux、Windows でのブラウザ起動
- **コンテキスト処理** - タイムアウトとキャンセル

## 使い方

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

## 出力モード

**通常の出力:**

```
Items
=====
NAME    VALUE
foo     bar
baz     qux
```

**`--github-issue` フラグ使用時:**

```markdown
## Items

| NAME | VALUE |
|--------|--------|
| foo | bar |
| baz | qux |
```

## コンテキストメソッド

| メソッド | 説明 |
|--------|-------------|
| `c.Header(title)` | セクションヘッダーを出力 |
| `c.KeyValue(key, val)` | キーと値のペアを出力 |
| `c.NewTable(headers...)` | テーブルを作成 |
| `c.List(items...)` | 箇条書きリストを出力 |
| `c.Code(code)` | コードブロックを出力 |
| `c.Success(msg)` | 成功メッセージを出力 |
| `c.Warning(msg)` | 警告を出力 |
| `c.Open(url)` | ブラウザで URL を開く |
