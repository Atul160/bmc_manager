package services

import (
	"bmc_manager/bmc"
	"bmc_manager/utils"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ManagePower manages power operations for different BMC types
func ManagePower(bmcType, ipAddress, action string) (bool, error) {
	if !utils.PingIP(ipAddress) {
		return false, fmt.Errorf("Unable to ping %s", ipAddress)
	}

	client, err := bmc.NewBMCClient(bmcType, ipAddress)
	if err != nil {
		return false, err
	}

	allowedActions := []string{"on", "off", "reset", "bmcreset"}
	if !slices.Contains(allowedActions, strings.ToLower(action)) {
		return false, errors.New("Invalid Power Action")
	}

	if action == "bmcreset" {
		err = client.BMCReset()
	} else {
		err = client.SetPower(action)
	}

	if err != nil {
		errstring := err.Error()
		if strings.Contains(errstring, "Power is on") || strings.Contains(errstring, "Server is already powered ON") {
			return false, fmt.Errorf("Server is Already Powered On")
		} else if strings.Contains(errstring, "Power is off") || strings.Contains(errstring, "Server is already powered OFF") {
			return false, fmt.Errorf("Server is Already Powered Off")
		} else {
			return false, err
		}
	}

	return true, nil
}
