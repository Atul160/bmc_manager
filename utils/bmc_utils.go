package utils

import (
	"fmt"
	"strings"
)

func GetBMCType(ipAddress string) (string, error) {

	var bmctype = "unknown"
	var url = "https://" + ipAddress + "/redfish/v1"

	// if !PingIP(ipAddress) {
	// 	return "", fmt.Errorf("Unable to ping %s", ipAddress)
	// }

	resp, err := InvokeRestAPI(url, "GET", nil, "", "", nil)
	if err != nil {
		if TCPPing(ipAddress, "22") {
			return "lenovoimm", nil
		} else {
			return bmctype, err
		}
	}

	_, resultstring, err := ReadResponseBody(resp)
	if err != nil {
		// return "", err
		fmt.Print(err)
	}

	if strings.Contains(resultstring, "redfish") {
		if strings.Contains(resultstring, "HP RESTful") {
			bmctype = "hpe"
		} else if strings.Contains(resultstring, "Dell") {
			bmctype = "dell"
		} else if strings.Contains(resultstring, "Lenovo") {
			bmctype = "lenovoxcc"
		} else {
			bmctype = "nutanix"
		}
	}

	return bmctype, err
}

func ValidatePowerOptions(input string, bmctype string) bool {
	var allowedOptions []string
	if bmctype == "dell" {
		allowedOptions = []string{
			"CheckStatus",
			"On",
			"Off",
			"ForceRestart",
			"GracefulRestart",
			"GracefulShutdown",
			"PushPowerButton",
			"Nmi",
			"PowerCycle",
		}
	} else if bmctype == "hpe" || bmctype == "lenovo" {
		allowedOptions = []string{
			"CheckStatus",
			"On",
			"Off",
			"ForceRestart",
			"PushPowerButton",
			"Nmi",
		}

	} else {
		allowedOptions = []string{
			"CheckStatus",
			"On",
			"Off",
			"Cycle",
		}

	}
	for _, option := range allowedOptions {
		if strings.EqualFold(input, option) {
			return true
		}
	}
	return false
}

/*
// Firmware update
func downloadImagePayload(image string) error {
	startTime = getCurrentTime() // Use a library for getting the current time

	log.Printf("- INFO, downloading \"%s\" image, this may take a few minutes depending on the size of the payload\n", image)

	var err error
	var response *http.Response

	// Get the UpdateService URI
	url := fmt.Sprintf(redfishBaseURL, idracIP) + updateServicePath
	if authToken != "" {
		response, err = http.Get(url, nil)
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(username, password)
		response, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("GET  request failed to get UpdateService: %s", string(body))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return err
	}

	httpPushURI, ok := data["HttpPushUri"].(string)
	if !ok {
		httpPushURI = "/redfish/v1/UpdateService/FirmwareInventory"
	}

	// Get the image details
	imageUrl := fmt.Sprintf(redfishBaseURL, idracIP) + httpPushURI

	filename := image
	imagePath := filepath.Join(imageLocation, filename)

	var req *http.Request
	if authToken != "" {
		req, err = http.NewRequest("GET", imageUrl, nil)
		if err != nil {
			return err
		}
		req.Header.Set("X-Auth-Token", authToken)
	} else {
		req, err = http.NewRequest("GET", imageUrl, nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(username, password)
	}

	response, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GET request failed to download image details: %s", response.Status)
	}

	etag := response.Header.Get("ETag") // Extract the ETag header

	// Prepare the multipart form data for uploading the image
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(
		map[string][]string{"Content-Disposition": {fmt.Sprintf("form-data; filename=%s", filename)}},
	)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	var headers map[string]string
	if authToken != "" {
		headers = map[string]string{"X-Auth-Token": authToken, "If-Match": etag}
	} else {
		headers = map[string]string{"If-Match": etag}
	}

	url = fmt.Sprintf(redfishBaseURL, idracIP) + httpPushURI
	response, err = http.Post(url, "multipart/form-data", body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("POST command failed to download image payload: %s", string(body))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Println("- FAIL, SimpleUpdate using HttpPushUri is not supported on this iDRAC release")
		return nil
	}

	availableEntry = data["Id"].(string)
	log.Printf("- INFO, AVAILABLE entry created for download image \"%s\" is \"%s\"\n", image, availableEntry)

	return nil
}
*/
/*
func installImagePayload() error {
	url := fmt.Sprintf(redfishBaseURL, idracIP) + "/UpdateService/Actions/UpdateService.SimpleUpdate"

	var payload map[string]interface{}
	if args["reboot"] {
		payload = map[string]interface{}{
			"ImageURI":                    fmt.Sprintf("%s/%s", httpPushURI, availableEntry),
			"@Redfish.OperationApplyTime": "Immediate",
		}
	} else {
		payload = map[string]interface{}{
			"ImageURI":                    fmt.Sprintf("%s/%s", httpPushURI, availableEntry),
			"@Redfish.OperationApplyTime": "OnReset",
		}
	}

	var req *http.Request
	if authToken != "" {
		req, err := http.NewRequest("POST", url, bytes.NewReader(MarshalJSON(payload)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Token", authToken)
	} else {
		req, err := http.NewRequest("POST", url, bytes.NewReader(MarshalJSON(payload)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(username, password)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("POST command failed to check job status: %s", string(body))
	}

	jobID = response.Header.Get("Location").Split("/")[-1]
	log.Printf("- PASS, update job ID %s successfully created, script will now loop polling the job status\n", jobID)

	return nil
}
*/

/*
func main() {
	// Replace with your actual values
	idracIP = "your_idrac_ip"
	username = "your_username"
	password = "your_password"
	authToken = "your_auth_token"
	verifyCert = true
	imageLocation = "path/to/image_location"

	image := "image_name.bin"

	err := downloadImagePayload(image)
	if err != nil {
		log.Fatalf("Error downloading image payload: %v", err)
	}

	err = installImagePayload()
	if err != nil {
		log.Fatalf("Error installing image payload: %v", err)
	}
}
*/
