package main

import (
	"ecc-bmc/api"
	"ecc-bmc/config"
	"ecc-bmc/middleware"
	"fmt"
	"log"
	"os"
	"strings"

	docs "ecc-bmc/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	// Load env configuration (e.g., ServerPort, API keys)
	envconfig := config.Load()

	// this will fill the placeholders dynamically based on environment.
	docs.SwaggerInfo.Title = fmt.Sprintf("ECC BMC Endpoint [%s]", strings.ToUpper(envconfig.Env))
	docs.SwaggerInfo.Description = fmt.Sprintf("ECC BMC Endpoint & This is an %s instance", os.Getenv("ENV"))

}

// @version					1.0
// @Security     BasicAuth || JWT
// @securityDefinitions.basic	BasicAuth
// @securityDefinitions.apiKey	JWT
// @in							header
// @name						Authorization
// @schema						Bearer token
// @description JWT Authorization header using the Bearer schema. Example: "Authorization: Bearer {token}"
func main() {
	// gin.SetMode(gin.ReleaseMode)
	// Initialize the Gin router
	router := gin.Default()

	// Load env configuration (e.g., ServerPort, API keys)
	envconfig := config.Load()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// logger = logging.NewLogger()

	// Apply global middleware (authentication, logging)
	router.Use(middleware.LoggingMiddleware())

	// Intialize file logging
	middleware.InitLogger()

	// API route for Authentication
	router.POST("/bmc/auth", api.TokenHandler)

	// Define API routes
	router.POST("/bmc/systeminfo", middleware.JWTAuthMiddleware(), api.SystemInfoHandler)
	router.POST("/bmc/firmwareinfo", middleware.JWTAuthMiddleware(), api.FirmwareInfoHandler)
	router.POST("/bmc/logs", middleware.JWTAuthMiddleware(), api.LogsHandler)
	router.POST("/bmc/power", middleware.JWTAuthMiddleware(), api.PowerHandler)
	router.POST("/bmc/firmwareupdate", middleware.JWTAuthMiddleware(), api.FirmwareUpdateHandler)

	// Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	log.Printf("Starting BMC Manager API server [%s]...", envconfig.Env)
	if err := router.Run(":" + envconfig.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
