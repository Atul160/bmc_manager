package main

import (
	"ecc-bmc/api"
	"ecc-bmc/config"
	"ecc-bmc/middleware"
	"log"
	"strings"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"

	docs "ecc-bmc/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	// this will fill the placeholders dynamically based on environment.
	docs.SwaggerInfo.Title = fmt.Sprintf("ECC BMC Endpoint [%s]", strings.ToUpper(os.Getenv("ENV")))
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

	// Load configuration (e.g., BMC credentials, API keys)
	config.Load()

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
	router.POST("/bmc/firmwareinfo", api.FirmwareInfoHandler)
	router.POST("/bmc/power", middleware.JWTAuthMiddleware(), api.PowerHandler)
	router.POST("/bmc/firmware", api.FirmwareUpdateHandler)

	// Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	log.Println("Starting BMC Manager API server...")
	if err := router.Run("localhost:8081"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
