package api

import (
	"ecc-bmc/services"
	"ecc-bmc/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// PowerHandler handles power management requests.
// @Summary Manage power for BMC
// @Description This endpoint allows the user to run power actions [on | off | reset | bmcreset] on a BMC device [dell | hpe | lenovoxcc | lenovoimm | nutanix].
// @Tags Power
// @Accept json
// @Produce json
// @Param power_request body PowerRequest true "Power request parameters"
// @Success 200 {object} PowerResponse "Success response"
// @Failure 400 {object} ErrorResponse "Bad request"
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

	if err := ValidatePowerAction(req.Action); err != nil {
		// Handle error
	}

	utils.LogInfo("Power Action: " + string(req.Action) + ", for " + req.IPAddress + ", initiated by: " + c.GetString("user"))

	result := services.ManagePower(strings.ToLower(string(req.BMCType)), strings.ToLower(req.IPAddress), strings.ToLower(string(req.Action)))

	c.JSON(http.StatusOK, gin.H{"result": result})

}
