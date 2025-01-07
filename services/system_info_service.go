package services

import (
	"ecc-bmc/bmc"
	"ecc-bmc/utils"
	"errors"
	"fmt"
)

// GetSystemInfo retrieves system information for the specified BMC.
func GetSystemInfo(bmcType, ipAddress string) (SystemInfo, error) {
	if bmcType == "" {
		var err = errors.New("")
		bmcType, err = utils.GetBMCType(ipAddress)
		if err != nil {
			fmt.Println(err)
		}
	}

	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		return SystemInfo{}, err
	}

	// Create a channel to receive the result and error
	resultChan := make(chan map[string]interface{})
	errorChan := make(chan error)

	// Call the GetSystemInfo function concurrently
	go func() {
		result, err := client.GetSystemInfo()
		if err != nil {
			errorChan <- err
			// utils.LogError("GetSystemInfo Failed", err)
			// return SystemInfo{}, err
		} else {
			resultChan <- result
		}
	}()

	var inforesult map[string]interface{}
	select {
	case result := <-resultChan:
		inforesult = result
	case err := <-errorChan:
		return SystemInfo{}, err
	}

	var result SystemInfo
	health, ok := inforesult["Status"].(map[string]interface{})
	if ok {
		result.Health = health["Health"].(string)
	} else {
		result.Health = inforesult["Health"].(string)
	}

	result.Device = inforesult["device"]
	_, ok = inforesult["HostName"]
	if ok {
		result.HostName = inforesult["HostName"].(string)
	}

	result.Manufacturer = inforesult["Manufacturer"].(string)
	result.Model = inforesult["Model"].(string)
	result.PowerState = inforesult["PowerState"].(string)

	_, ok = inforesult["BiosVersion"]
	if ok {
		result.BiosVersion = inforesult["BiosVersion"].(string)
	}
	_, ok = inforesult["SerialNumber"]
	if ok {
		result.SerialNumber = inforesult["SerialNumber"].(string)
	}
	_, ok = inforesult["MemorySummary"]
	if ok {
		result.Memory = inforesult["MemorySummary"]
	}
	_, ok = inforesult["ProcessorSummary"]
	if ok {
		result.CPU = inforesult["ProcessorSummary"]
	}

	return result, nil
}
