package bmc

import (
	"bmc_manager/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type NutanixIPMIClient struct {
	IPAddress string
	Username  string
	Password  string
}

// Connect implements BMCClient.
func (c *NutanixIPMIClient) Connect() (map[string]string, error) {
	return nil, fmt.Errorf("This operation is unimplemented")
}

func NewNutanixIPMIClient(ipAddress, username, password string) *NutanixIPMIClient {
	return &NutanixIPMIClient{
		IPAddress: ipAddress,
		Username:  username,
		Password:  password,
	}
}

// SetPower sends a Redfish request to execute the input power action on the server
func (c *NutanixIPMIClient) SetPower(action string) error {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/1/Actions/ComputerSystem.Reset", c.IPAddress)
	var resettype string
	switch action {
	case "on":
		resettype = "On"
	case "off":
		resettype = "ForceOff"
	case "reset":
		resettype = "ForceRestart"
	default:
		return errors.New("Invalid Power Action")
	}
	body := map[string]string{
		"ResetType": resettype,
	}
	return c.sendRedfishRequest(url, body)
}

// BMC Reset sends a Redfish request to reset the server BMC
func (c *NutanixIPMIClient) BMCReset() error {
	url := fmt.Sprintf("https://%s/redfish/v1/Managers/1/Actions/Manager.Reset", c.IPAddress)
	body := map[string]string{
		"ResetType": "GracefulRestart",
	}
	return c.sendRedfishRequest(url, body)
}

// GetSystemInfo retrieves system information via Redfish
func (c *NutanixIPMIClient) GetSystemInfo() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/1", c.IPAddress)
	response, err := utils.InvokeRestAPI(url, "GET", nil, c.Username, c.Password, nil)
	if err != nil {
		return nil, err
	}
	result, err := utils.ReadResponseBody(response)
	if err != nil {
		return nil, err
	}

	result["device"] = c.IPAddress
	return result, nil
}

// GetFirmwareInfo implements BMCClient.
func (c *NutanixIPMIClient) GetFirmwareInfo() (map[string]interface{}, error) {
	// panic("unimplemented")
	url := fmt.Sprintf("https://%s/redfish/v1/UpdateService/FirmwareInventory?$expand=*($levels=1)", c.IPAddress)
	response, err := utils.InvokeRestAPI(url, "GET", nil, c.Username, c.Password, nil)
	if err != nil {
		return nil, err
	}
	result, err := utils.ReadResponseBody(response)
	if err != nil {
		return nil, err
	}
	// var result map[string]interface{}
	// if err := json.Unmarshal(response, &result); err != nil {
	// 	return nil, err
	// }
	return result, nil
}

// UpdateFirmware implements BMCClient.
func (c *NutanixIPMIClient) UpdateFirmware(filePath string) error {
	panic("unimplemented")
}

// sendRedfishRequest is a helper method for sending Redfish actions
func (c *NutanixIPMIClient) sendRedfishRequest(url string, body map[string]string) error {
	// Marshal the request body to JSON
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	_, err = utils.InvokeRestAPI(url, "POST", nil, c.Username, c.Password, bytes.NewBuffer(payload))
	return err
}