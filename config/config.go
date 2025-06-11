package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// BMCConfig holds configuration for a specific BMC type
type BMCConfig struct {
	Username string
	Password string
}

// Config holds the application-wide configuration values
type Config struct {
	ServerPort         string
	Env                string
	LogLevel           string
	JWTSecret          string
	LDAPConfig         LDAPConfig
	DellConfig         BMCConfig
	DCHPEConfig        BMCConfig
	vertical1HPEConfig BMCConfig
	LenovoXCCConfig    BMCConfig
	LenovoIMMConfig    BMCConfig
	DCNutanixConfig    BMCConfig
}

// LDAPConfig holds LDAP configuration
type LDAPConfig struct {
	Server string
	BaseDN string
	// BindPass string
}

// Load loads configuration from environment variables
func Load() Config {
	err := godotenv.Load("./secrets/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	return Config{
		ServerPort: getEnv("SERVER_PORT", "8086"),
		Env:        getEnv("ENV", "CERT"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		JWTSecret:  os.Getenv("JWTSecret"),
		LDAPConfig: LDAPConfig{
			Server: os.Getenv("LDAP_SERVER"),
			BaseDN: os.Getenv("LDAP_BASE_DN"),
			// BindPass: os.Getenv("LDAP_BIND_PASS"),
		},
		DellConfig: BMCConfig{
			Username: os.Getenv("DELL_BMC_USERNAME"),
			Password: os.Getenv("DELL_BMC_PASSWORD"),
		},
		DCHPEConfig: BMCConfig{
			Username: os.Getenv("HPE_BMC_USERNAME"),
			Password: os.Getenv("HPE_BMC_PASSWORD"),
		},
		vertical1HPEConfig: BMCConfig{
			Username: os.Getenv("HPE_BMC_USERNAME"),
			Password: os.Getenv("HPE_BMC_PASSWORD"),
		},
		LenovoXCCConfig: BMCConfig{
			Username: os.Getenv("XCC_USERNAME"),
			Password: os.Getenv("XCC_PASSWORD"),
		},
		LenovoIMMConfig: BMCConfig{
			Username: os.Getenv("IMM_USERNAME"),
			Password: os.Getenv("IMM_PASSWORD"),
		},
		DCNutanixConfig: BMCConfig{
			Username: os.Getenv("vertical2_IPMI_USERNAME"),
			Password: os.Getenv("vertical2_IPMI_PASSWORD"),
		},
	}
}

// getEnv retrieves environment variables, or returns a fallback default value
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
