package services

// PowerResponse holds the result for each IP processed
type PowerResponse struct {
	IPAddress string `json:"ip_address"`
	Action    string `json:"action"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// SystemInfo represents the system information structure.
type SystemInfo struct {
	Device         interface{} `json:"device,omitempty"`
	Health         string      `json:"health,omitempty"`
	Manufacturer   string      `json:"manufacturer,omitempty"`
	PowerState     interface{} `json:"powerstate,omitempty"`
	Model          string      `json:"model,omitempty"`
	BiosVersion    string      `json:"biosversion,omitempty"`
	SerialNumber   string      `json:"serialnumber,omitempty"`
	HostName       interface{} `json:"hostname,omitempty"`
	ResponseStatus string      `json:"responsestatus,omitempty"`
	Memory         interface{} `json:"memory,omitempty"`
	CPU            interface{} `json:"cpu,omitempty"`
	Error          string      `json:"error,omitempty"`
}

// FirmwareResult holds the firmware info for each IP address
type FirmwareResult struct {
	IPAddress string                 `json:"ip_address"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// LogInfo holds the logs info for each IP address
type LogInfo struct {
	Device  string     `json:"device"`
	LogType string     `json:"log_type"`
	Data    []LogEntry `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}

type LogEntry struct {
	ID       string `json:"Id"`
	Message  string `json:"Message"`
	Created  string `json:"Created"`
	Severity string `json:"Severity"`
	ODataID  string `json:"@odata.id"`
}

type LogResponse struct {
	Members          []LogEntry `json:"Members"`
	MemberODataCount int        `json:"Members@odata.count"`
}
