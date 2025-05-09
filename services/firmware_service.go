package services

import (
	"ecc-bmc/bmc"
	"ecc-bmc/utils"
	"fmt"
	"os"
	"strings"
)

func GetFirmwareInfoParallel(bmcType, ipAddresses string) []FirmwareResult {
	ips := strings.Split(ipAddresses, ",")
	return utils.ParallelExecute(ips, func(ip string) (FirmwareResult, error) {
		data, err := GetFirmwareInfo(bmcType, ip)
		if err != nil {
			return FirmwareResult{IPAddress: ip, Error: err.Error()}, nil
		}
		return FirmwareResult{IPAddress: ip, Data: data}, nil
	})
}

// GetFirmwareInfo retrieves Firmware information for a single BMC.
func GetFirmwareInfo(bmcType, ipAddress string) (map[string]interface{}, error) {
	if bmcType == "" {
		var err error
		bmcType, err = utils.GetBMCType(ipAddress)
		if err != nil {
			fmt.Println(err)
		}
	}

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
