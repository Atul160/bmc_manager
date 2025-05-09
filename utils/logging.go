package utils

import (
	"fmt"
	"log"
	"net/http"
)

// LogInfo logs an informational message with a timestamp.
func LogInfo(message string) {
	fmt.Println("[INFO] :", message)
	log.Println("[INFO] :", message)
}

// time.Now().Format(time.RFC3339),

// LogError logs an error message with a timestamp.
func LogError(message string, err error) {
	fmt.Printf("[ERROR] : %s - %v", message, err)
	log.Printf("[ERROR] : %s - %v", message, err)
}

// LogRequest logs the HTTP request details
func LogRequest(r *http.Request) {
	fmt.Printf("Request: %s %s, Headers: %+v", r.Method, r.URL, r.Header)
	log.Printf("Request: %s %s, Headers: %+v", r.Method, r.URL, r.Header)
}

// LogResponse logs the HTTP response details
func LogResponse(status int, msg string) {
	fmt.Printf("Response: %d, Message: %s", status, msg)
	log.Printf("Response: %d, Message: %s", status, msg)
}
