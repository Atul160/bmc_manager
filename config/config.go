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
	ServerPort      string
	Env             string
	LogLevel        string
	JWTSecret       string
	LDAPConfig      LDAPConfig
	DellConfig      BMCConfig
	DCHPEConfig     BMCConfig
	StoresHPEConfig BMCConfig
	LenovoXCCConfig BMCConfig
	LenovoIMMConfig BMCConfig
	DCNutanixConfig BMCConfig
}

// LDAPConfig holds LDAP configuration
type LDAPConfig struct {
	Server string
	BaseDN string
	// BindPass string
}

// Load loads configuration from environment variables
func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Env:        getEnv("ENV", "development"),
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
		StoresHPEConfig: BMCConfig{
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
			Username: os.Getenv("DC_IPMI_USERNAME"),
			Password: os.Getenv("DC_IPMI_PASSWORD"),
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
