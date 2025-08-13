# args

A lightweight and flexible command-line argument parser for Go application.

## Features

- Simple and intuitive API
- Support for flags, options, and positional arguments
- Built-in help generation
- Type-safe argument handling
- Subcommand support
- Custom validation

## Installation

```bash
go get github.com/marcomit/args
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/marcomit/args"
)

func main() {
    parser := argparser.New("myapp", "A sample application")
    
    // Add flags
    verbose := parser.Flag("verbose", "v", "Enable verbose output")
    output := parser.String("output", "o", "output.txt", "Output file path")
    count := parser.Int("count", "c", 1, "Number of iterations")
    
    // Parse arguments
    err := parser.Parse()
    if err != nil {
        fmt.Println(err)
        parser.Help()
        return
    }
    
    // Use parsed values
    if *verbose {
        fmt.Println("Verbose mode enabled")
    }
    fmt.Printf("Output file: %s\n", *output)
    fmt.Printf("Count: %d\n", *count)
}
```

## Usage Examples

### Basic Flags and Options

```go
// Boolean flag
debug := parser.Flag("debug", "d", "Enable debug mode")

// String option with default value
config := parser.String("config", "c", "config.json", "Configuration file")

// Integer option
port := parser.Int("port", "p", 8080, "Server port")

// Float option
rate := parser.Float("rate", "r", 1.0, "Processing rate")
```

### Positional Arguments

```go
// Required positional arguments
inputFile := parser.Positional("input", "Input file path")
outputFile := parser.Positional("output", "Output file path")

// Optional positional arguments
optionalArg := parser.PositionalOptional("optional", "Optional argument")
```

### Subcommands

```go
parser := argparser.New("git-clone", "Git-like command structure")

// Add subcommands
cloneCmd := parser.Subcommand("clone", "Clone a repository")
pushCmd := parser.Subcommand("push", "Push changes")

// Add flags to subcommands
cloneUrl := cloneCmd.String("url", "u", "", "Repository URL")
branch := cloneCmd.String("branch", "b", "main", "Branch name")

force := pushCmd.Flag("force", "f", "Force push")
```

### Custom Validation

```go
port := parser.Int("port", "p", 8080, "Server port")

// Add custom validation
parser.Validate(func() error {
    if *port < 1 || *port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    return nil
})
```

## API Reference

### Parser Methods

| Method | Description |
|--------|-------------|
| `New(name, description)` | Create a new parser |
| `Flag(name, short, help)` | Add a boolean flag |
| `String(name, short, default, help)` | Add a string option |
| `Int(name, short, default, help)` | Add an integer option |
| `Float(name, short, default, help)` | Add a float option |
| `Positional(name, help)` | Add required positional argument |
| `PositionalOptional(name, help)` | Add optional positional argument |
| `Subcommand(name, help)` | Add a subcommand |
| `Parse()` | Parse command-line arguments |
| `Help()` | Print help message |
| `Validate(func)` | Add custom validation |

### Help Output

The parser automatically generates help text:

```
Usage: myapp [OPTIONS] <input> <output>

A sample application

Arguments:
  input              Input file path
  output             Output file path

Options:
  -v, --verbose      Enable verbose output
  -o, --output       Output file path (default: output.txt)
  -c, --count        Number of iterations (default: 1)
  -h, --help         Show this help message
```

## Error Handling

The parser provides detailed error messages for common issues:

- Missing required arguments
- Invalid flag formats
- Type conversion errors
- Custom validation failures

```go
err := parser.Parse()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    parser.Help()
    os.Exit(1)
}
```

## Todo
 - [] argparser.New
 - [] argparser.New

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details
