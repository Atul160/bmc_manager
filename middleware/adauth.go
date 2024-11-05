package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"ecc-bmc/utils"

	"github.com/gin-gonic/gin"
)

// Config for Active Directory Authentication
type Config struct {
	Server string // Active Directory server URL
	BaseDN string // Base DN for the search
}

// NewConfig initializes configuration
// func NewConfig() *Config {
// 	return &Config{
// 		Server: "ldap://your.ad.server:389", // Replace with your AD server
// 		BaseDN: "dc=example,dc=com",         // Replace with your base DN
// 	}
// }

// ADAuthMiddleware checks for AD authentication
func ADAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch the Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Extract the Basic Auth credentials
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Decode the base64 encoded credentials
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to decode credentials"})
			c.Abort()
			return
		}

		// Split username and password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		username := credentials[0]
		password := credentials[1]

		// Authenticate the user
		if err := utils.AuthenticateADUser(username, password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid AD credentials"})
			c.Abort()
			return
		}

		// If authentication is successful, proceed to the next handler
		c.Next()
	}
}
