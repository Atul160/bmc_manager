package api

import (
	"bmc_manager/services"
	"bmc_manager/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SystemInfoHandler handles system info retrieval requests
// @Summary System Info for BMC
// @Description This endpoint fetches the system info of a BMC device.
// @Tags info
// @Accept json
// @Produce json
// @Param system_info body SystemInfoRequest true "System Info request parameters"
// @Success 200 {object} map[string]bool "Success response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure	401	{string} string	"Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /systeminfo [post]
func SystemInfoHandler(c *gin.Context) {
	var req SystemInfoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid system info request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the system info retrieval service
	// info, err := services.GetSystemInfo(c.Query("bmc_type"), c.Query("ip_address"))
	info, err := services.GetSystemInfo(string(req.BMCType), req.IPAddress)
	if err != nil {
		utils.LogError("System info retrieval failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"system_info": info})
}
