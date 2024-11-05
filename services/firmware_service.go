package services

import (
	"ecc-bmc/bmc"
	"fmt"
	"os"
)

// Get Firmware Info retrieves Firmware information for the specified BMC.
func GetFirmwareInfo(bmcType, ipAddress string) (map[string]interface{}, error) {
	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		return nil, err
	}

	return client.GetFirmwareInfo()
}

// UpdateFirmware performs firmware updates for the specified BMC type.
func UpdateFirmware(bmcType, ipAddress, firmwarePath string) (bool, error) {
	// Check if the firmware file exists
	if _, err := os.Stat(firmwarePath); os.IsNotExist(err) {
		return false, fmt.Errorf("firmware file does not exist at path: %s", firmwarePath)
	}

	// Create the appropriate BMC client
	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		return false, err
	}

	// Call the UpdateFirmware method on the client
	err = client.UpdateFirmware(firmwarePath)
	if err != nil {
		return false, err
	}

	return true, nil
}
