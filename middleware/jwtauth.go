// middleware/jwt.go
package middleware

import (
	"ecc-bmc/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTSecret is the secret key used to sign the JWT tokens
var JWTSecret []byte = []byte(os.Getenv("JWTSecret"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWTAuthMiddleware checks if the JWT token is valid for incoming requests
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from Authorization header (e.g., "Bearer <token>")
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims := &Claims{}

		// Parse and validate JWT token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		utils.LogInfo("JWT token verified for: " + claims.Username)

		// Set user in the context
		c.Set("user", claims.Username)
		c.Next()
	}
}
