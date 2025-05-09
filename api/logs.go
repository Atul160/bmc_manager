package api

import (
	"ecc-bmc/services"
	"ecc-bmc/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// LogsHandler handles log info retrieval requests
// @Summary Log Info for BMC
// @Description This endpoint fetches the log info of a BMC device [dell | hpe | lenovoxcc | lenovoimm | nutanix].
// @Tags Info
// @Accept json
// @Produce json
// @Param log_info body LogsRequest true "log Info request parameters"
// @Success 200 {object} map[string]bool "Success response"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure	401	{string} string	"Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bmc/logs [post]
// @Security		JWT
func LogsHandler(c *gin.Context) {
	var req LogsRequest

	// Bind JSON body to LogsRequest object
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid system info request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	utils.LogInfo("Get Logs - " + req.IPAddress + ", initiated by: " + c.GetString("user"))

	info := services.GetLogsParallel(strings.ToLower(string(req.BMCType)), strings.ToLower(string(req.LogType)), strings.ToLower(req.IPAddress))

	c.JSON(http.StatusOK, gin.H{"result": info})
}
