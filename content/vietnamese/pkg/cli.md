---
title: cli
import_path: www.ubuntusoftware.net/pkg/cli
repo_url: https://github.com/joeblew999/ubuntu-website
description: Khung giao diện dòng lệnh (CLI) chung cho các công cụ phần mềm Ubuntu. Các tùy chọn tiêu chuẩn, định dạng đầu ra và định dạng Markdown cho các vấn đề trên GitHub.
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/cli
license: MIT
author: Gerard Webb
created_at: 2025-12-16T00:00:00Z
updated_at: 2025-12-16T00:00:00Z
has_binary: false
---

## Tính năng

- **Cờ tiêu chuẩn* - `--github-issue`, `--verbose`, `--version` có sẵn ngay từ đầu
- **Định dạng đầu ra** - Bảng, danh sách, tiêu đề, khối mã
- **Chế độ Markdown* - Đầu ra tương thích với GitHub Issue khi sử dụng tùy chọn `--github-issue`
- **Đa nền tảng** - Mở trình duyệt trên macOS, Linux, Windows
- **Xử lý ngữ cảnh** - Thời gian chờ và hủy bỏ

## Cách sử dụng

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

## Chế độ đầu ra

**Đầu ra bình thường:**

```
Items
=====
NAME    VALUE
foo     bar
baz     qux
```

**Với tùy chọn `--github-issue`:**

```markdown
## Items

| NAME | VALUE |
|--------|--------|
| foo | bar |
| baz | qux |
```

## Phương pháp bối cảnh

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
