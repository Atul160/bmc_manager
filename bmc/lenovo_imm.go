package bmc

import (
	"bytes"
	"ecc-bmc/utils"
	"encoding/json"
	"fmt"
	"strings"
)

type LenovoIMMClient struct {
	IPAddress string
	Username  string
	Password  string
}

// Connect implements BMCClient.
func (c *LenovoIMMClient) Connect() (map[string]string, error) {
	return nil, fmt.Errorf("This operation is unimplemented")
}

func NewLenovoIMMClient(ipAddress, username, password string) *LenovoIMMClient {
	return &LenovoIMMClient{
		IPAddress: ipAddress,
		Username:  username,
		Password:  password,
	}
}

// SetPower sends a Redfish request to execute the input power action on the server
func (c *LenovoIMMClient) SetPower(action string) error {

	client, err := utils.CreateSSHConnection(c.IPAddress, c.Username, c.Password)
	if err != nil {
		return err
	}

	//execute ssh cmds for old imm to perform power actions
	if action == "On" {
		_, err = utils.ExecuteSSHCommand(client, "power on")
	} else if action == "Off" {
		_, err = utils.ExecuteSSHCommand(client, "power off")
	} else if action == "Cycle" {
		_, err = utils.ExecuteSSHCommand(client, "power cycle")
	}

	if err != nil {
		return fmt.Errorf("Error performing Action:%s on %s:%s", action, c.IPAddress, err)
	} else {
		return nil
	}
}

// BMC Reset sends a Redfish request to reset the server BMC
func (c *LenovoIMMClient) BMCReset() error {
	url := fmt.Sprintf("https://%s/redfish/v1/Managers/1/Actions/Manager.Reset", c.IPAddress)
	body := map[string]string{
		"ResetType": "GracefulRestart",
	}
	return c.sendRedfishRequest(url, body)
}

// GetSystemInfo retrieves system information via Redfish
func (c *LenovoIMMClient) GetSystemInfo() (map[string]interface{}, error) {
	sshconn, err := utils.CreateSSHConnection(c.IPAddress, c.Username, c.Password)
	if err != nil {
		return nil, err
	}

	out, err := utils.ExecuteSSHCommand(sshconn, "syshealth summary")
	if err != nil {
		return nil, fmt.Errorf("Unable to execute SSH command")
	}

	str := string(out)
	var PowerState string

	if strings.Contains(str, "On") {
		PowerState = "On"
	} else {
		PowerState = "Off"
	}
	Info := map[string]interface{}{
		"device":       c.IPAddress,
		"Health":       "NA | Old IMM",
		"Manufacturer": "Lenovo",
		"PowerState":   PowerState,
		"Model":        "Old M5",
	}

	return Info, nil
}

// GetFirmwareInfo implements BMCClient.
func (c *LenovoIMMClient) GetFirmwareInfo() (map[string]interface{}, error) {
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

// UpdateFirmware implements BMCClient.
func (c *LenovoIMMClient) UpdateFirmware(filePath string) error {
	panic("unimplemented")
}

// sendRedfishRequest is a helper method for sending Redfish actions
func (c *LenovoIMMClient) sendRedfishRequest(url string, body map[string]string) error {
	// Marshal the request body to JSON
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	_, err = utils.InvokeRestAPI(url, "POST", nil, c.Username, c.Password, bytes.NewBuffer(payload))
	return err
}
