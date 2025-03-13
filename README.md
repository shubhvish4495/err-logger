# LogHound

LogHound is a simple, thread-safe error logging library for Go that provides stack trace capture and structured logging capabilities.

## Features

- Thread-safe error logging
- Customizable logger interface
- Automatic stack trace capture
- Structured logging with additional context
- Filtering of internal Go runtime frames
- Default logger implementation included

## Installation

```bash
go get github.com/shubhvish4495/loghound
```

## Usage

### Basic Usage

```go
package main

import "github.com/shubhvish4495/loghound"

func main() {
    err := someFunction()
    if err != nil {
        loghound.LogError(err, map[string]interface{}{
            "context": "main function",
            "status": "failed",
        })
    }
}
```

### Custom Logger

You can implement your own logger by implementing the `Logger` interface:

```go
type CustomLogger struct {
    // your fields here
}

func (c *CustomLogger) Error(message string, args ...any) {
    // your custom logging implementation
}

// Set the custom logger
loghound.SetLogger(&CustomLogger{})
```

## API Reference

### Types

- `Logger` - Interface for custom loggers
- `DefaultLogger` - Basic implementation using standard library
- `errorLog` - Structure containing error details and stack trace

### Functions

- `SetLogger(logger Logger)` - Set a custom logger implementation
- `LogError(err error, args map[string]interface{})` - Log an error with context

## Stack Trace Output

The stack trace output excludes Go runtime frames and includes:
- Function name
- File location
- Line number

## Thread Safety

All operations in LogHound are thread-safe and can be used concurrently.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Your License Here]