package loghound

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

// mockLogger is a mock implementation of the Logger interface for testing
type mockLogger struct {
	messages []string
}

func (m *mockLogger) Error(message string, args ...any) {
	m.messages = append(m.messages, message)
}

func TestSetLogger(t *testing.T) {
	mock := &mockLogger{}
	SetLogger(mock)

	if globalLogger != mock {
		t.Errorf("expected globalLogger to be set to mockLogger, got %T", globalLogger)
	}
}

func TestLogError(t *testing.T) {
	mock := &mockLogger{}
	SetLogger(mock)

	err := os.Setenv("GOROOT", "/usr/local/go")
	if err != nil {
		t.Fatalf("failed to set GOROOT: %v", err)
	}

	testError := "test error"
	LogError(errors.New(testError), map[string]interface{}{"key": "value"})

	if len(mock.messages) == 0 {
		t.Fatalf("expected messages to be logged, got none")
	}

	if !strings.Contains(mock.messages[0], testError) {
		t.Errorf("expected log message to contain %q, got %q", testError, mock.messages[0])
	}
}

func TestDefaultLogger(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	defaultLogger := &DefaultLogger{}
	defaultLogger.Error("test message")

	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("expected log output to contain %q, got %q", "test message", buf.String())
	}
}
