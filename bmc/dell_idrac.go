package bmc

import (
	"bytes"
	"ecc-bmc/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

type DellIDRACClient struct {
	IPAddress string
	Username  string
	Password  string
}

// GetLogs implements BMCClient.
func (c *DellIDRACClient) GetLogs(logType string) ([]map[string]interface{}, error) {
	logURIs := map[string]string{
		"system":     fmt.Sprintf("https://%s/redfish/v1/Managers/iDRAC.Embedded.1/LogServices/Sel/Entries", c.IPAddress),
		"management": fmt.Sprintf("https://%s/redfish/v1/Managers/iDRAC.Embedded.1/LogServices/Lclog/Entries", c.IPAddress),
		"fault":      fmt.Sprintf("https://%s/redfish/v1/Managers/iDRAC.Embedded.1/LogServices/FaultList/Entries", c.IPAddress),
	}
	uri, exists := logURIs[logType]
	if !exists {
		return nil, fmt.Errorf("invalid log type")
	}

	resp, err := utils.InvokeRestAPI(uri, "GET", nil, c.Username, c.Password, nil)
	if err != nil {
		return nil, err
	}
	result, _, _, err := utils.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	members, _ := result["Members"].([]interface{})

	var logs []map[string]interface{}

	for _, entry := range members {
		if entryMap, ok := entry.(map[string]interface{}); ok {
			logs = append(logs, entryMap)
		}
	}

	var mu sync.Mutex // Mutex to avoid concurrent slice modification
	if count, ok := result["Members@odata.count"].(float64); ok && int(count) > 50 {
		var wg sync.WaitGroup
		errChan := make(chan error, (int(count)-50)/50) // Channel to capture errors

		for i := 50; i < 1000; i += 50 {
			wg.Add(1)
			go func(skip int) {
				defer wg.Done()
				newURI := fmt.Sprintf("%s?$skip=%d", uri, skip)

				// Call API in parallel
				chunkResp, err := utils.InvokeRestAPI(newURI, "GET", nil, c.Username, c.Password, nil)
				if err != nil {
					errChan <- err
					return
				}

				chunkresult, _, _, err := utils.ReadResponseBody(chunkResp)
				if err != nil {
					errChan <- err
					return
				}

				if chunkMembers, ok := chunkresult["Members"].([]interface{}); ok {
					var tempLogs []map[string]interface{} // Temporary local slice
					for _, chunkEntry := range chunkMembers {
						if entryMap, ok := chunkEntry.(map[string]interface{}); ok {
							tempLogs = append(tempLogs, entryMap)
						}
					}
					mu.Lock()
					logs = append(logs, tempLogs...) // Append results safely
					mu.Unlock()
				}
			}(i) // Pass `i` as argument to prevent race conditions
		}

		wg.Wait()
		close(errChan)

		// Check for errors from goroutines
		for err := range errChan {
			if err != nil {
				log.Println("Error fetching logs:", err)
				return nil, err
			}
		}
	}

	return logs, nil
}

// Connect implements BMCClient.
func (c *DellIDRACClient) Connect() (map[string]string, error) {
	return nil, fmt.Errorf("this operation is unimplemented")
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
		return errors.New("invalid power action")
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
	result, _, _, err := utils.ReadResponseBody(response)
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
	result, _, _, err := utils.ReadResponseBody(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *DellIDRACClient) UpdateFirmware(firmwarePath string) error {
	return errors.New("firmware update not implemented for Dell iDRAC")
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
