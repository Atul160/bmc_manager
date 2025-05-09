package services

import (
	"ecc-bmc/bmc"
	"ecc-bmc/utils"
	"strings"
)

func GetSystemInfoParallel(bmcType, ipAddresses string) []SystemInfo {
	ips := strings.Split(ipAddresses, ",")
	return utils.ParallelExecute(ips, func(ip string) (SystemInfo, error) {
		data, err := GetSystemInfo(bmcType, ip)
		if err != nil {
			return SystemInfo{Device: ip, Error: err.Error()}, nil
		}
		return data, nil
	})
}

// GetSystemInfo retrieves system information for a single BMC.
func GetSystemInfo(bmcType, ipAddress string) (SystemInfo, error) {
	if bmcType == "" {
		var err error
		bmcType, err = utils.GetBMCType(ipAddress)
		if err != nil {
			utils.LogError("Failed to get BMC type", err)
			return SystemInfo{}, err
		}
	}

	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		utils.LogError("Failed to create BMC client", err)
		return SystemInfo{}, err
	}

	inforesult, err := client.GetSystemInfo()
	if err != nil {
		utils.LogError("GetSystemInfo Failed", err)
		return SystemInfo{}, err
	}

	return parseSystemInfo(inforesult)
}

// parseSystemInfo parses the system information from the result map.
func parseSystemInfo(inforesult map[string]interface{}) (SystemInfo, error) {
	var result SystemInfo

	if health, ok := inforesult["Status"].(map[string]interface{}); ok {
		result.Health = safeString(health["Health"])
	} else {
		result.Health = safeString(inforesult["Health"])
	}

	result.Device = safeString(inforesult["device"])
	result.HostName = safeString(inforesult["HostName"])
	result.Manufacturer = safeString(inforesult["Manufacturer"])
	result.Model = safeString(inforesult["Model"])
	result.PowerState = safeString(inforesult["PowerState"])
	result.BiosVersion = safeString(inforesult["BiosVersion"])
	result.SerialNumber = safeString(inforesult["SerialNumber"])
	result.Memory = inforesult["MemorySummary"]
	result.CPU = inforesult["ProcessorSummary"]

	return result, nil
}

// safeString safely asserts a value to a string.
func safeString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
