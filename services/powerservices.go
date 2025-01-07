package services

import (
	"ecc-bmc/bmc"
	"ecc-bmc/utils"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ManagePower manages power operations for different BMC types
func ManagePower(bmcType, ipAddress, action string) (bool, error) {
	if bmcType == "" {
		var err = errors.New("")
		bmcType, err = utils.GetBMCType(ipAddress)
		if err != nil {
			fmt.Println(err)
		}
	}

	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		return false, err
	}

	allowedActions := []string{"on", "off", "reset", "bmcreset"}
	if !slices.Contains(allowedActions, strings.ToLower(action)) {
		return false, errors.New("invalid power action")
	}

	// Create a channel to receive the result and error
	// resultChan := make(chan bool)
	errorChan := make(chan error)

	// Call the GetSystemInfo function concurrently
	go func() {
		if action == "bmcreset" {
			err = client.BMCReset()
		} else {
			err = client.SetPower(action)
		}

		if err != nil {
			errstring := err.Error()

			if strings.Contains(errstring, "Power is on") || strings.Contains(errstring, "Server is already powered ON") {
				errorChan <- fmt.Errorf("server is already powered on")
			} else if strings.Contains(errstring, "Power is off") || strings.Contains(errstring, "Server is already powered OFF") {
				errorChan <- fmt.Errorf("server is already powered off")
			} else {
				errorChan <- err
			}
		} else {
			errorChan <- nil
		}
	}()

	if errorChan != nil {
		return false, <-errorChan
	} else {
		return true, nil
	}
}
