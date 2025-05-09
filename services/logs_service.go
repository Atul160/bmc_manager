package services

import (
	"ecc-bmc/bmc"
	"ecc-bmc/utils"
	"fmt"
	"strings"
)

// GetLogsParallel fetches logs info for multiple IPs in parallel
// func GetLogsParallel(bmcType, logType, ipAddresses string) []LogInfo {
// 	ips := strings.Split(ipAddresses, ",") // Split comma-separated IPs
// 	var wg sync.WaitGroup
// 	results := make([]LogInfo, len(ips))

// 	for i, ip := range ips {
// 		wg.Add(1)
// 		go func(index int, ip string) {
// 			defer wg.Done()
// 			data, err := GetLogs(bmcType, logType, ip)
// 			if err != nil {
// 				results[index] = LogInfo{Device: ip, Error: err.Error()}
// 			} else {
// 				results[index] = LogInfo{Device: ip, LogType: logType, Data: data}
// 			}
// 		}(i, ip)
// 	}

// 	wg.Wait() // Wait for all goroutines to complete
// 	return results
// }

func GetLogsParallel(bmcType, logType, ipAddresses string) []LogInfo {
	ips := strings.Split(ipAddresses, ",")
	return utils.ParallelExecute(ips, func(ip string) (LogInfo, error) {
		data, err := GetLogs(bmcType, logType, ip)
		if err != nil {
			return LogInfo{Device: ip, Error: err.Error()}, nil
		}
		return LogInfo{Device: ip, LogType: logType, Data: data}, nil
	})
}

// GetLogs retrieves logs information for a single BMC.
func GetLogs(bmcType, logType, ipAddress string) ([]LogEntry, error) {
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

	logResult, err := client.GetLogs(logType)
	if err != nil {
		utils.LogError("GetLogs Failed for: "+ipAddress, err)
		return []LogEntry{}, err
	}

	// // Check if "Members" exists and is an array
	// members, exists := logResult["Members"].([]interface{})
	// if !exists {
	// 	return nil, fmt.Errorf("members field not found or not an array")
	// }

	var logs []LogEntry

	// Iterate and convert each map to LogEntry struct
	for _, entryMap := range logResult {
		logEntry := LogEntry{
			ID:       getString(entryMap, "Id"),
			Created:  getString(entryMap, "Created"),
			Message:  getString(entryMap, "Message"),
			Severity: getString(entryMap, "Severity"),
			ODataID:  getString(entryMap, "@odata.id"),
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}

// Helper function to safely extract string values
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}
