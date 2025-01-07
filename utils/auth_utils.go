package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	ldap "github.com/go-ldap/ldap/v3"
)

// Credentials structure to bind AD login credentials
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// var cfg = config.Load()

// AuthenticateUser authenticates a user against Active Directory
func AuthenticateADUser(username, password string) error {
	// ldapconfig := cfg.LDAPConfig
	// fmt.Println(username, password, os.Getenv("LDAP_SERVER"))
	// l, err := ldap.Dial("tcp", ldapconfig.Server)
	l, err := ldap.DialURL(os.Getenv("LDAP_SERVER"))
	if err != nil {
		return err
	}
	defer l.Close()

	// Attempt to bind to the AD server using the user's credentials
	err = l.Bind(username, password)
	if err != nil {
		return err
	}

	LogInfo("AD authentication successful for: " + username)
	// If bind is successful, return true
	return nil
}

// JWTSecret is the secret key used to sign the JWT tokens
var JWTSecret []byte = []byte(os.Getenv("JWTSecret"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token after successful AD authentication
func GenerateJWT(username string) (string, string, error) {
	tokenexpiry := time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      tokenexpiry,
		})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, time.Unix(tokenexpiry, 0).String(), nil
}

// Function to extract username and password from a JWT token
func ExtractJWT(tokenString string) (string, string, error) {
	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("error parsing token: %v", err)
	}

	// Extract the claims and validate the token
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Return the username and password from the claims
		return claims.Username, claims.Password, nil
	}

	return "", "", fmt.Errorf("invalid token")
}
