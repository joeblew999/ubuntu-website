---
title: "命令行"
import_path: www.ubuntusoftware.net/pkg/cli
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Ubuntu 软件工具的共享 CLI 框架。提供标准参数、输出格式化和 GitHub issue Markdown。"
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/cli
license: MIT
author: Gerard Webb
created_at: 2025-12-16T00:00:00Z
updated_at: 2025-12-16T00:00:00Z
has_binary: false
---

## 功能特性

- **标准标志** - 开箱即用的 `--github-issue`、`--verbose`、`--version`
- **输出格式化** - 表格、列表、标题、代码块
- **Markdown 模式** - 使用 `--github-issue` 生成适合 GitHub issue 的输出
- **跨平台** - 支持在 macOS、Linux、Windows 上打开浏览器
- **上下文处理** - 超时与取消

## 用法

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

## 输出模式

**普通输出：**

```
Items
=====
NAME    VALUE
foo     bar
baz     qux
```

**使用 `--github-issue` 标志：**

```markdown
## Items

| NAME | VALUE |
|--------|--------|
| foo | bar |
| baz | qux |
```

## 上下文方法

| 方法 | 描述 |
|--------|-------------|
| `c.Header(title)` | 打印章节标题 |
| `c.KeyValue(key, val)` | 打印键值对 |
| `c.NewTable(headers...)` | 创建表格 |
| `c.List(items...)` | 打印项目符号列表 |
| `c.Code(code)` | 打印代码块 |
| `c.Success(msg)` | 打印成功消息 |
| `c.Warning(msg)` | 打印警告 |
| `c.Open(url)` | 在浏览器中打开 URL |
