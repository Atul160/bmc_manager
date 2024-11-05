package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Struct to represent the nested structure of the JSON
type ExtendedInfo struct {
	Message     string `json:"Message,omitempty"`
	MessageArgs string `json:"-"`
	MessageId   string `json:"MessageId,omitempty"`
	Resolution  string `json:"Resolution,omitempty"`
	Severity    string `json:"Severity,omitempty"`
}

type ErrorDetails struct {
	ExtendedInfo []ExtendedInfo `json:"@Message.ExtendedInfo,omitempty"`
	Code         string         `json:"code,omitempty"`
	Message      string         `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error,omitempty"`
}

// InvokeRestAPI makes a REST API call to the specified URL with the provided method, headers, and body
func InvokeRestAPI(url, method string, headers map[string]string, username, password string, body io.Reader) (*http.Response, error) {
	// Create the HTTP client with SSL verification disabled
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set basic authentication
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// // Check the response status code
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	// }

	// Check for non-200 status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		// jsonString := strings.TrimPrefix(string(bodyBytes), "error response: ")
		// body, _ := UnMarshalJSON(bodyBytes)

		// var errorResponse ErrorResponse
		// if err := json.Unmarshal([]byte(jsonString), &errorResponse); err != nil {
		// 	log.Fatalf("Error unmarshaling JSON: %v", err)
		// }

		// // Convert the slice []string field into a joined string for each ExtendedInfo entry
		// for i := range errorResponse.Error.ExtendedInfo {
		// 	// Joining []string fields with comma or other delimiter
		// 	errorResponse.Error.ExtendedInfo[i].MessageArgs = strings.Join([]string{}, ",")
		// }

		// df := dataframe.LoadStructs(errorResponse.Error.ExtendedInfo)
		// errmsg := df.Subset(0).Col("Message").Records()[0]

		return nil, fmt.Errorf(string(bodyBytes))
	}

	return resp, nil
}

func ReadResponseBody(resp *http.Response) (map[string]interface{}, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return data, nil
}

func ReadRequestBody(resp *http.Request) (map[string]interface{}, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return data, nil
}

// Helper function to marshal JSON data with error handling
func MarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
	}
	return data
}

func UnMarshalJSON(jsonData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON: %w", err)
	}
	return data, nil
}
