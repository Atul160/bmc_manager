package api

import (
	"ecc-bmc/services"
	"ecc-bmc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PowerHandler handles power management requests.
// @Summary Manage power for BMC
// @Description This endpoint allows the user to power on/off/reset a BMC device.
// @Tags power
// @Accept json
// @Produce json
// @Param power_request body PowerRequest true "Power request parameters"
// @Success 200 {object} map[string]bool "Success response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure	401	{string} string	"Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bmc/power [post]
// @Security		JWT
func PowerHandler(c *gin.Context) {
	var req PowerRequest
	// var data map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid Power request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := validatePowerAction(req.Action); err != nil {
		// Handle error
	}

	success, err := services.ManagePower(string(req.BMCType), req.IPAddress, string(req.Action))
	if err != nil {
		utils.LogError("Power management failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}
