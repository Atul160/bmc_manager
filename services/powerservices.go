package services

import (
	"ecc-bmc/bmc"
	"ecc-bmc/utils"
	"fmt"
	"slices"
	"strings"
	"sync"
)

// ManagePower manages power operations for multiple IPs concurrently
func ManagePower(bmcType, ipAddresses, action string) []PowerResponse {
	ipList := strings.Split(ipAddresses, ",") // Split comma-separated IPs
	var wg sync.WaitGroup
	results := make(chan PowerResponse, len(ipList)) // Channel to collect results

	for _, ip := range ipList {
		wg.Add(1)
		go func(ipAddress string) {
			defer wg.Done()
			result := PowerResponse{IPAddress: ipAddress, Action: action}

			if bmcType == "" {
				var err error
				bmcType, err = utils.GetBMCType(ipAddress)
				if err != nil {
					result.Error = fmt.Sprintf("Error detecting BMC type: %v", err)
					results <- result
					return
				}
			}

			client, err := bmc.NewBMCClient(bmcType, ipAddress)
			if err != nil {
				result.Error = fmt.Sprintf("Error creating BMC client: %v", err)
				results <- result
				return
			}

			allowedActions := []string{"on", "off", "reset", "bmcreset"}
			if !slices.Contains(allowedActions, strings.ToLower(action)) {
				result.Error = "Invalid power action"
				results <- result
				return
			}

			// Execute power action
			var powerErr error
			if action == "bmcreset" {
				powerErr = client.BMCReset()
			} else {
				powerErr = client.SetPower(action)
			}

			// Handle errors with detailed messages
			if powerErr != nil {
				errMsg := powerErr.Error()
				if strings.Contains(errMsg, "Power is on") || strings.Contains(errMsg, "Server is already powered ON") {
					result.Error = "Server is already powered on"
				} else if strings.Contains(errMsg, "Power is off") || strings.Contains(errMsg, "Server is already powered OFF") {
					result.Error = "Server is already powered off"
				} else {
					result.Error = errMsg
				}
				results <- result
				return
			}

			result.Success = true
			results <- result
		}(strings.TrimSpace(ip)) // Trim spaces to avoid errors
	}

	// Wait for all go routines to finish
	wg.Wait()
	close(results)

	// Collect all results
	var response []PowerResponse
	for res := range results {
		response = append(response, res)
	}

	return response
}
