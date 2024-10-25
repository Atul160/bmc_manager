// api/auth.go
package api

import (
	"bmc_manager/utils"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary AD Authentication
// @Description Authenticates a user Basic Auth against Active Directory and generates a JWT token
// @Tags auth
// @Success 200 {string} string "Authenticated successfully"
// @Failure 401 {string} string "Authentication failed"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth [get]
// TokenHandler generates a JWT token after validating AD credentials
func TokenHandler(c *gin.Context) {
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

	// Authenticate using AD
	if err := utils.AuthenticateADUser(username, password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid AD credentials"})
		return
	}

	// Generate JWT token
	token, expiry, err := utils.GenerateJWT(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return JWT token
	c.JSON(http.StatusOK, gin.H{"token": token, "expiry": expiry})
}
