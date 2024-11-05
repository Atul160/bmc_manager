// api/auth.go
package api

import (
	"ecc-bmc/utils"
	"net/http"
	"github.com/gin-gonic/gin"
)

// @Summary		Generate Token
// @Description	generate JWT token with basic authentication
// @Tags			Token
// @Accept			json
// @Produce		json
// @Success		200	{string}	string	"OK"
// @Router			/bmc/auth [POST]
// @Security		BasicAuth
func TokenHandler(c *gin.Context) {
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
