package middleware

import (
	"net/http"

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
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
			c.Abort()
			return
		}

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
