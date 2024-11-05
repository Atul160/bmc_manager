package bmc

import (
	"bytes"
	"ecc-bmc/utils"
	"encoding/json"
	"errors"
	"fmt"
)

type DellIDRACClient struct {
	IPAddress string
	Username  string
	Password  string
}

// Connect implements BMCClient.
func (c *DellIDRACClient) Connect() (map[string]string, error) {
	return nil, fmt.Errorf("This operation is unimplemented")
}

func NewDellIDRACClient(ipAddress, username, password string) *DellIDRACClient {
	return &DellIDRACClient{
		IPAddress: ipAddress,
		Username:  username,
		Password:  password,
	}
}

// SetPower sends a Redfish request to execute the input power action on the server
func (c *DellIDRACClient) SetPower(action string) error {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset", c.IPAddress)
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
func (c *DellIDRACClient) BMCReset() error {
	url := fmt.Sprintf("https://%s/redfish/v1/Managers/iDRAC.Embedded.1/Actions/Manager.Reset", c.IPAddress)
	body := map[string]string{
		"ResetType": "GracefulRestart",
	}
	return c.sendRedfishRequest(url, body)
}

// GetSystemInfo retrieves system information using Redfish API
func (c *DellIDRACClient) GetSystemInfo() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/Systems/System.Embedded.1", c.IPAddress)
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
func (c *DellIDRACClient) GetFirmwareInfo() (map[string]interface{}, error) {
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

	return result, nil
}

func (c *DellIDRACClient) UpdateFirmware(firmwarePath string) error {
	return errors.New("Firmware update not implemented for Dell iDRAC")
}

// sendRedfishRequest is a helper method for sending Redfish actions
func (c *DellIDRACClient) sendRedfishRequest(url string, body map[string]string) error {
	// Marshal the request body to JSON
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	_, err = utils.InvokeRestAPI(url, "POST", nil, c.Username, c.Password, bytes.NewBuffer(payload))
	return err
}
