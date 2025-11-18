# Via Framework Patterns Guide

This guide documents architectural patterns and best practices for building pages with the Via reactive web framework in the `internal/env/web` package.

## Table of Contents

1. [Understanding Via's Initialization Behavior](#understanding-vias-initialization-behavior)
2. [The Three Page Patterns](#the-three-page-patterns)
3. [LazyLoader Helper](#lazyloader-helper)
4. [Best Practices](#best-practices)

---

## Understanding Via's Initialization Behavior

**Critical Concept**: Via validates all page init functions at startup by running them.

When you register a page with `v.Page("/path", func(c *via.Context) { ... })`, Via immediately calls this function during server startup to validate the page structure. This means:

- ✅ **Safe at init time**: Creating signals, actions, and views
- ❌ **Unsafe at init time**: API calls, database queries, file I/O
- ❌ **Result**: Expensive operations run before any user visits the page

### Example of the Problem

```go
// ❌ BAD: API call happens at server startup
v.Page("/zones", func(c *via.Context) {
    zones, _ := env.ListZones(token, accountID) // Runs during Via validation!

    c.View(func() h.H {
        // Render zones...
    })
})
```

### The Solution: Lazy Loading

Move expensive operations inside `c.View()` where they only run when the page is actually rendered:

```go
// ✅ GOOD: API call happens only when user visits the page
v.Page("/zones", func(c *via.Context) {
    zonesLoader := NewLazyLoader(func() ([]env.Zone, error) {
        return env.ListZones(token, accountID)
    })

    c.View(func() h.H {
        zones, _ := zonesLoader.Get() // Runs on first render only
        // Render zones...
    })
})
```

---

## The Three Page Patterns

Analysis of `internal/env/web` route files reveals three distinct architectural patterns:

### Pattern 1: Simple Form Pages (No Lazy Loading)

**When to use**: Pages with only form fields and simple validation, no expensive operations.

**Characteristics**:
- Form fields backed by signals
- Validation happens on user action
- No API calls at page load time

**Examples**:
- [`route_cloudflare.go`](route_cloudflare.go) - Token input only
- [`route_cloudflare_step2.go`](route_cloudflare_step2.go) - Account ID input
- [`route_claude.go`](route_claude.go) - Claude API key input

**Pattern**:
```go
func simplePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
    svc := env.NewService(mockMode)
    fields := CreateFormFields(c, cfg, []string{env.KeySomeField})
    saveMessage := c.Signal("")

    saveAction := c.Action(func() {
        // Validate and save
    })

    c.View(func() h.H {
        return h.Main(
            // Render form
        )
    })
}
```

### Pattern 2: Lazy Loading Pages (LazyLoader)

**When to use**: Pages that load external data (API calls) that should only happen when the page is visited.

**Characteristics**:
- Uses `LazyLoader[T]` helper for deferred loading
- Loader function contains API calls or expensive operations
- Data cached after first load
- Handles errors gracefully with user feedback

**Examples**:
- [`route_cloudflare_step3.go`](route_cloudflare_step3.go) - Loads zones from Cloudflare API
- [`route_cloudflare_step4.go`](route_cloudflare_step4.go) - Loads Pages projects from Cloudflare API

**Pattern**:
```go
func lazyLoadPage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
    svc := env.NewService(mockMode)
    fields := CreateFormFields(c, cfg, []string{env.KeySomeField})
    dataMessage := c.Signal("") // For loading status/errors

    // LazyLoader - only runs when Get() is called in View
    dataLoader := NewLazyLoader(func() ([]DataType, error) {
        token := cfg.Get(env.KeyAPIToken)

        if mockMode {
            return []DataType{{ID: "mock"}}, nil
        }

        if token == "" || env.IsPlaceholder(token) {
            return []DataType{}, nil
        }

        return api.LoadData(token)
    })

    c.View(func() h.H {
        // Load data on first render
        data, err := dataLoader.Get()
        if err != nil {
            log.Printf("Failed to load data: %v", err)
            dataMessage.SetValue("error:" + err.Error())
        }

        return h.Main(
            // Render with data
        )
    })
}
```

### Pattern 3: Conditional Expensive Operations

**When to use**: Pages with operations that should run conditionally based on runtime state.

**Characteristics**:
- Expensive operations guarded by conditionals
- May run on user actions rather than page load
- Often involves real-time command execution

**Examples**:
- [`route_deploy.go`](route_deploy.go) - Hugo build and Wrangler deploy only on button click

**Pattern**:
```go
func conditionalPage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
    outputMessage := c.Signal("")

    expensiveAction := c.Action(func() {
        // Only runs when user clicks button
        result := performExpensiveOperation()
        outputMessage.SetValue(result)
        c.Sync()
    })

    c.View(func() h.H {
        return h.Main(
            h.Button(h.Text("Run Operation"), expensiveAction.OnClick()),
            // Show output
        )
    })
}
```

---

## LazyLoader Helper

The `LazyLoader[T]` generic helper in [`components.go`](components.go) abstracts the lazy loading pattern.

### API

```go
// Create a lazy loader with a function that returns data
loader := NewLazyLoader(func() ([]DataType, error) {
    // Expensive operation here
    return fetchData()
})

// Get data (calls loader function only once, then caches)
data, err := loader.Get()
```

### How It Works

1. **Initialization**: Store the loader function without calling it
2. **First Get()**: Call loader function, cache result, return data
3. **Subsequent Get()**: Return cached data immediately

### Benefits

- ✅ **Prevents startup overhead**: Expensive operations deferred until needed
- ✅ **Caching**: Data loaded once and reused across View renders
- ✅ **Type-safe**: Go generics provide compile-time type checking
- ✅ **Error handling**: Propagates errors to caller for appropriate UI feedback
- ✅ **Testable**: Easy to provide mock data in loader function

### Example Usage

```go
// In page init (runs at startup - this is cheap)
zonesLoader := NewLazyLoader(func() ([]env.Zone, error) {
    token := cfg.Get(env.KeyCloudflareAPIToken)
    accountID := cfg.Get(env.KeyCloudflareAccountID)

    // Mock data for testing
    if mockMode {
        return []env.Zone{{ID: "mock-1", Name: "example.com"}}, nil
    }

    // Return empty if credentials not set
    if token == "" || env.IsPlaceholder(token) {
        return []env.Zone{}, nil
    }

    // Real API call
    return env.ListZones(token, accountID)
})

// In c.View() (runs when user visits page - API call happens here)
c.View(func() h.H {
    zones, err := zonesLoader.Get() // First call: runs loader, subsequent: returns cache
    if err != nil {
        // Handle error in UI
    }

    // Build UI with zones...
})
```

---

## Best Practices

### 1. Choose the Right Pattern

- **Simple forms** → Pattern 1 (no LazyLoader needed)
- **Load external data** → Pattern 2 (use LazyLoader)
- **User-triggered operations** → Pattern 3 (action-based)

### 2. Signal Usage

Via signals only support simple types (strings, booleans, numbers):

```go
// ✅ GOOD: Simple types in signals
message := c.Signal("")
count := c.Signal(0)
isLoading := c.Signal(false)

// ❌ BAD: Complex types in signals
zones := c.Signal([]env.Zone{}) // Will cause compilation errors!
```

For complex data types, use closure variables or LazyLoader:

```go
// ✅ GOOD: Complex types with LazyLoader
zonesLoader := NewLazyLoader(func() ([]env.Zone, error) {
    // ...
})
```

### 3. Error Handling

Always handle errors from expensive operations and show them in the UI:

```go
data, err := loader.Get()
if err != nil {
    log.Printf("Failed to load: %v", err) // Server logs
    dataMessage.SetValue("error:" + err.Error()) // User feedback
}
```

### 4. Mock Mode Support

Support mock mode in all LazyLoaders for testing:

```go
dataLoader := NewLazyLoader(func() ([]DataType, error) {
    if mockMode {
        return []DataType{{ID: "mock", Name: "Test Data"}}, nil
    }
    // Real implementation
})
```

### 5. Credential Checks

Check for valid credentials before making API calls:

```go
if token == "" || env.IsPlaceholder(token) {
    return []DataType{}, nil // Return empty, not error
}
```

### 6. Logging

Use structured logging for debugging:

```go
log.Printf("Failed to fetch zones: %v", err) // Context + error
```

The enhanced logging in [`logging.go`](logging.go) automatically adds page context to errors.

---

## Migration Guide

### Before: Manual Closure Variables

```go
var zonesCache []env.Zone
var zonesLoaded bool

c.View(func() h.H {
    if !zonesLoaded {
        token := cfg.Get(env.KeyCloudflareAPIToken)
        if token != "" {
            zonesCache, _ = env.ListZones(token, accountID)
        }
        zonesLoaded = true
    }
    // Use zonesCache...
})
```

### After: LazyLoader

```go
zonesLoader := NewLazyLoader(func() ([]env.Zone, error) {
    token := cfg.Get(env.KeyCloudflareAPIToken)
    if token == "" {
        return []env.Zone{}, nil
    }
    return env.ListZones(token, accountID)
})

c.View(func() h.H {
    zones, _ := zonesLoader.Get()
    // Use zones...
})
```

**Benefits**: Less boilerplate, clearer intent, type-safe, reusable pattern.

---

## References

- Via Framework: https://github.com/go-via/via
- LazyLoader implementation: [`components.go`](components.go)
- Enhanced logging: [`logging.go`](logging.go)
- Example pages: All files matching `route_*.go` in this directory
