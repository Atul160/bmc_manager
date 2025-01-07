package utils

import (
	"log"
	"net/http"
)

// LogInfo logs an informational message with a timestamp.
func LogInfo(message string) {
	log.Printf("[INFO] : %s", message)
}

// time.Now().Format(time.RFC3339),

// LogError logs an error message with a timestamp.
func LogError(message string, err error) {
	log.Printf("[ERROR] : %s - %v", message, err)
}

// LogRequest logs the HTTP request details
func LogRequest(r *http.Request) {
	log.Printf("Request: %s %s, Headers: %+v", r.Method, r.URL, r.Header)
}

// LogResponse logs the HTTP response details
func LogResponse(status int, msg string) {
	log.Printf("Response: %d, Message: %s", status, msg)
}
