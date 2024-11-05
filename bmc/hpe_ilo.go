package bmc

import (
	"bytes"
	"ecc-bmc/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HPEILOClient struct {
	IPAddress string
	Username  string
	Password  string
}

func NewHPEILOClient(ipAddress, username, password string) *HPEILOClient {
	return &HPEILOClient{
		IPAddress: ipAddress,
		Username:  username,
		Password:  password,
	}
}

func (c *HPEILOClient) Connect() (map[string]string, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/SessionService/Sessions/", c.IPAddress)
	body := map[string]string{
		"UserName": c.Username,
		"Password": c.Password,
	}
	payload, _ := json.Marshal(body)

	response, err := utils.InvokeRestAPI(url, "POST", nil, "", "", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	result := map[string]string{
		"X-Auth-Token": response.Header["X-Auth-Token"][0],
		"Location":     response.Header["Location"][0],
	}
	return result, nil
}

// SetPower sends a Redfish request to execute the input power action on the server
func (c *HPEILOClient) SetPower(action string) error {
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
func (c *HPEILOClient) BMCReset() error {
	url := fmt.Sprintf("https://%s/redfish/v1/Managers/1/Actions/Manager.Reset", c.IPAddress)
	body := map[string]string{
		"ResetType": "GracefulRestart",
	}
	return c.sendRedfishRequest(url, body)
}

// GetSystemInfo retrieves system information via Redfish
func (c *HPEILOClient) GetSystemInfo() (map[string]interface{}, error) {

	url := fmt.Sprintf("https://%s/redfish/v1/Systems/1", c.IPAddress)

	var response *http.Response
	var err error

	response, err = utils.InvokeRestAPI(url, "GET", nil, c.Username, c.Password, nil)
	if err != nil {
		session, _ := c.Connect()
		headers := map[string]string{
			"X-Auth-Token": session["X-Auth-Token"],
		}
		response, err = utils.InvokeRestAPI(url, "GET", headers, "", "", nil)
		if err != nil {
			return nil, err
		}
	}

	result, err := utils.ReadResponseBody(response)
	if err != nil {
		return nil, err
	}
	result["device"] = c.IPAddress
	return result, nil
}

// GetFirmwareInfo implements BMCClient.
func (c *HPEILOClient) GetFirmwareInfo() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/redfish/v1/UpdateService/?$expand=", c.IPAddress)

	var response *http.Response
	var err error

	response, err = utils.InvokeRestAPI(url, "GET", nil, c.Username, c.Password, nil)
	if err != nil {
		session, _ := c.Connect()
		headers := map[string]string{
			"X-Auth-Token": session["X-Auth-Token"],
		}
		response, err = utils.InvokeRestAPI(url, "GET", headers, "", "", nil)
		if err != nil {
			return nil, err
		}
	}

	result, err := utils.ReadResponseBody(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateFirmware implements BMCClient.
func (c *HPEILOClient) UpdateFirmware(filePath string) error {
	panic("unimplemented")
}

// sendRedfishRequest is a helper method for sending Redfish actions
func (c *HPEILOClient) sendRedfishRequest(url string, body map[string]string) error {
	// Marshal the request body to JSON
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	// var resp *http.Response

	_, err = utils.InvokeRestAPI(url, "POST", nil, c.Username, c.Password, bytes.NewBuffer(payload))
	if err != nil {
		session, _ := c.Connect()
		headers := map[string]string{
			"X-Auth-Token": session["X-Auth-Token"],
		}
		_, err = utils.InvokeRestAPI(url, "POST", headers, "", "", bytes.NewBuffer(payload))
		if err != nil {
			return err
		}
	}

	return err
}
