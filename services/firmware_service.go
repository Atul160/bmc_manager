package services

import (
	"bmc_manager/bmc"
	"fmt"
	"os"
)

type FirmeareInfo struct {
	Device         interface{} `json:"device,omitempty"`
	Health         string      `json:"health,omitempty"`
	Manufacturer   string      `json:"manufacturer,omitempty"`
	PowerState     interface{} `json:"powerstate"`
	Model          string      `json:"model,omitempty"`
	BiosVersion    string      `json:"biosversion,omitempty"`
	SerialNumber   string      `json:"serialnumber,omitempty"`
	HostName       interface{} `json:"hostname,omitempty"`
	ResponseStatus string      `json:"responsestatus,omitempty"`
	Memory         interface{} `json:"memory,omitempty"`
	CPU            interface{} `json:"cpu,omitempty"`
}

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
