package services

import (
	"bmc_manager/bmc"
	"bmc_manager/utils"
)

type SystemInfo struct {
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

// GetSystemInfo retrieves system information for the specified BMC.
func GetSystemInfo(bmcType, ipAddress string) (SystemInfo, error) {
	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		return SystemInfo{}, err
	}

	inforesult, err := client.GetSystemInfo()
	if err != nil {
		utils.LogError("GetSystemInfo Failed", err)
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