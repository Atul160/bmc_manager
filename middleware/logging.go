package middleware

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// InitLogger initializes the logger
func InitLogger() {
	// Define the log file path
	logFilePath := "./logs/app.log"

	// Extract the directory from the log file path
	logDir := filepath.Dir(logFilePath)

	// Create the logs directory if it doesn't exist
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Open or create the log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set the output of the log package to the file
	log.SetOutput(file)
	log.Println("Logger initialized successfully")
}

// LoggingMiddleware logs the details of each request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log details after request has been processed
		duration := time.Since(startTime)
		log.Printf("Request: %s %s | Status: %d | Duration: %v", c.Request.Method, c.Request.URL, c.Writer.Status(), duration)
	}
}
