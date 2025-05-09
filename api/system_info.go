package api

import (
	"ecc-bmc/services"
	"ecc-bmc/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SystemInfoHandler handles system info retrieval requests
// @Summary System Info for BMC
// @Description This endpoint fetches the system info of a BMC device [dell | hpe | lenovoxcc | lenovoimm | nutanix].
// @Tags Info
// @Accept json
// @Produce json
// @Param system_info body SystemInfoRequest true "System Info request parameters"
// @Success 200 {object} SystemInfoResponse "Success response"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure	401	{string} string	"Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bmc/systeminfo [post]
// @Security		JWT
func SystemInfoHandler(c *gin.Context) {
	var req SystemInfoRequest

	// Bind JSON body to SystemInfoRequest object
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid system info request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	utils.LogInfo("Get System Info - " + req.IPAddress + ", initiated by: " + c.GetString("user"))

	info := services.GetSystemInfoParallel(strings.ToLower(string(req.BMCType)), strings.ToLower(req.IPAddress))

	c.JSON(http.StatusOK, gin.H{"system_info": info})
}
