package loghound

import (
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Logger interface allows users to use custom loggers
type Logger interface {
	Error(message string, args ...any)
}

// DefaultLogger is a basic logger using the standard library
type DefaultLogger struct{}

func (d *DefaultLogger) Error(message string, args ...any) {
	log.Println(message)
}

// stackFrame represents a function call in the stack
type stackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// errorLog represents the full error log with context
type errorLog struct {
	Message    string                 `json:"message"`
	Args       map[string]interface{} `json:"args"`
	StackTrace []stackFrame           `json:"stack_trace"`
}

// Global logger instance (default is standard logger)
var (
	globalLogger   Logger = &DefaultLogger{}
	mu             sync.Mutex
	goRootLocation = os.Getenv("GOROOT")
)

// SetLogger allows users to set a custom logger
func SetLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

// captureStackTrace collects the stack trace up to a certain depth
func captureStackTrace(skip int) []stackFrame {
	var pcs [10]uintptr
	n := runtime.Callers(skip, pcs[:]) // Skip first `skip` frames

	var stack []stackFrame
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		if strings.Contains(file, goRootLocation) {
			continue
		}
		stack = append(stack, stackFrame{
			Function: fn.Name(),
			File:     file,
			Line:     line,
		})
	}
	return stack
}

// LogError logs an error with arguments and stack trace using the set logger
func LogError(err error, args map[string]interface{}) {
	logEntry := errorLog{
		Message:    err.Error(),
		Args:       args,
		StackTrace: captureStackTrace(3), // Skip LogError, CaptureStackTrace
	}

	// lock and defer unlock to ensure thread safety
	mu.Lock()
	defer mu.Unlock()
	globalLogger.Error(logEntry.Message, "stack_trace", logEntry.StackTrace, "args", logEntry.Args)
}
