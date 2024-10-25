package main

import (
	"bmc_manager/api"
	"bmc_manager/config"
	"bmc_manager/middleware"
	"log"

	// docs "bmc_manager/docs" // This is needed to load the swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// func init() {
// 	// this will fill the placeholders dynamically based on environment.
// 	docs.SwaggerInfo.Title = fmt.Sprintf("ECC ServiceNow Incident Endpoint [%s]", strings.ToUpper(snowhelper.AppInstance))
// 	docs.SwaggerInfo.Description = fmt.Sprintf("ECC ServiceNow Incident Endpoint & This is an %s instance", snowhelper.AppInstance)

// }

func main() {
	// gin.SetMode(gin.ReleaseMode)
	// Initialize the Gin router
	router := gin.Default()
	// router := echo.New()

	// Load configuration (e.g., BMC credentials, API keys)
	config.Load()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// logger = logging.NewLogger()

	// Apply global middleware (authentication, logging)
	router.Use(middleware.LoggingMiddleware())

	// Intialize file logging
	middleware.InitLogger()

	// Define API routes
	router.GET("/bmc/systeminfo", api.SystemInfoHandler)
	router.GET("/bmc/firmwareinfo", middleware.JWTAuthMiddleware(), api.FirmwareInfoHandler)
	router.POST("/bmc/power", api.PowerHandler)
	router.POST("/bmc/firmware", middleware.JWTAuthMiddleware(), api.FirmwareUpdateHandler)

	// Serve Swagger UI files
	router.StaticFile("/swagger-docs/doc.json", "./docs/swagger.json")

	// Handle Swagger UI requests
	router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("http://localhost:8081/swagger-docs/doc.json")))

	// Start the server
	log.Println("Starting BMC Manager API server...")
	if err := router.Run("localhost:8081"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
