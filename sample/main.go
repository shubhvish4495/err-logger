package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/shubhvish4495/loghound"
)

type OwnLogger struct {
	log *slog.Logger
}

func (o *OwnLogger) Error(message string, args ...interface{}) {
	o.log.Error(message, args...)
}

func divideHandler(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	// Convert to integers
	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Perform division
	result, err := divide(a, b)
	if err != nil {
		loghound.LogError(err, map[string]interface{}{
			"a": a,
			"b": b,
		})
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	ownLogger := &OwnLogger{
		log: log,
	}

	// Use a custom logger
	loghound.SetLogger(ownLogger)

	// Define HTTP routes
	http.HandleFunc("/divide", divideHandler)

	// Start the server
	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Error("Error starting server: %v", "error", err)
	}
}
