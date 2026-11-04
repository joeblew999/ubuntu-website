---
title: "cli"
import_path: www.ubuntusoftware.net/pkg/cli
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Framework CLI dùng chung cho các công cụ của Ubuntu Software. Các cờ tiêu chuẩn, định dạng đầu ra và markdown cho issue trên GitHub."
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/cli
license: MIT
author: Gerard Webb
created_at: 2025-12-16T00:00:00Z
updated_at: 2025-12-16T00:00:00Z
has_binary: false
---

## Tính năng

- **Cờ tiêu chuẩn** - `--github-issue`, `--verbose`, `--version` có sẵn ngay từ đầu
- **Định dạng đầu ra** - Bảng, danh sách, tiêu đề, khối mã
- **Chế độ Markdown** - Đầu ra thân thiện với GitHub issue bằng `--github-issue`
- **Đa nền tảng** - Mở trình duyệt trên macOS, Linux, Windows
- **Xử lý context** - Thời gian chờ và hủy bỏ

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

## Các chế độ đầu ra

**Đầu ra thông thường:**

```
Items
=====
NAME    VALUE
foo     bar
baz     qux
```

**Với cờ `--github-issue`:**

```markdown
## Items

| NAME | VALUE |
|--------|--------|
| foo | bar |
| baz | qux |
```

## Các phương thức của Context

| Phương thức | Mô tả |
|--------|-------------|
| `c.Header(title)` | In tiêu đề mục |
| `c.KeyValue(key, val)` | In cặp khóa-giá trị |
| `c.NewTable(headers...)` | Tạo bảng |
| `c.List(items...)` | In danh sách gạch đầu dòng |
| `c.Code(code)` | In khối mã |
| `c.Success(msg)` | In thông báo thành công |
| `c.Warning(msg)` | In cảnh báo |
| `c.Open(url)` | Mở URL trong trình duyệt |
