# Simple Go Argument Parser

Lightweight, fluent CLI argument parser for Go.

## Quick Start

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    app := New("myapp").
        Flag("verbose", "v", "Enable verbose output").
        Option("config", "c", "Config file path").
        Action(func(r *Result) error {
            fmt.Printf("Config: %s\n", r.Options["config"])
            return nil
        })

    app.Run(os.Args[1:])
}
```

## API

```go
// Create parser
app := New("name")

// Add arguments
app.Flag("name", "short", "help")           // Boolean flag
app.Option("name", "short", "help")         // String option
app.Command("name", "help")                 // Subcommand

// Set handler
app.Action(func(r *Result) error { ... })

// Parse and run
app.Run(os.Args[1:])
```

## Example

```go
app := New("docker")

run := app.Command("run", "Run container").
    Flag("detach", "d", "Run in background").
    Option("name", "", "Container name").
    Action(func(r *Result) error {
        fmt.Printf("Running: %s\n", r.Positionals[0])
        return nil
    })
```

**Usage:** `docker run -d --name=web nginx`

## Result Structure

```go
type Result struct {
    Command     []string          // ["run"]
    Flags       map[string]bool   // {"detach": true}
    Options     map[string]string // {"name": "web"}
    Positionals []string          // ["nginx"]
}
```
