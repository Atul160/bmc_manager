package api

import (
	"ecc-bmc/services"
	"ecc-bmc/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// FirmwareInfoHandler handles firmware info requests
// @Summary Firmware Info for BMC
// @Description This endpoint fetches the firmware info of a BMC device [dell | hpe | lenovoxcc | lenovoimm | nutanix].
// @Tags Info
// @Accept json
// @Produce json
// @Param firmware_info body FirmwareInfoRequest true "Firmware info request parameters"
// @Success 200 {object} map[string]bool "Success response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure	401	{string} string	"Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bmc/firmwareinfo [post]
// @Security		JWT
func FirmwareInfoHandler(c *gin.Context) {
	var req FirmwareInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid firmware info request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	utils.LogInfo("Get Firmware Info - " + req.IPAddress + ", initiated by: " + c.GetString("user"))

	info := services.GetFirmwareInfoParallel(strings.ToLower(string(req.BMCType)), strings.ToLower(req.IPAddress))

	c.JSON(http.StatusOK, gin.H{"firmware_info": info})
}

// FirmwareUpdateHandler handles firmware update requests
// @Summary Firmware Update for BMC
// @Description This endpoint updates the firmware of a BMC device.
// @Tags Patch
// @Accept json
// @Produce json
// @Param firmware_request body FirmwareUpdateRequest true "Firmware Update request parameters"
// @Success 200 {object} map[string]bool "Success response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bmc/firmwareupdate [post]
// @Security		JWT
// @Deprecated true
func FirmwareUpdateHandler(c *gin.Context) {
	var req FirmwareUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid firmware update request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the firmware update service
	success, err := services.UpdateFirmware(strings.ToLower(string(req.BMCType)), strings.ToLower(req.IPAddress), strings.ToLower(req.FirmwarePath))
	if err != nil {
		utils.LogError("Firmware update failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}
