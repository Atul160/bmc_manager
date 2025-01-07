package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// InitLogger initializes the logger
func InitLogger() {
	file, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
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
