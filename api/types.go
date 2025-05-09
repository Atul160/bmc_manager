package api

import "fmt"

type PowerAction string

const (
	On       PowerAction = "on"
	Off      PowerAction = "off"
	Reset    PowerAction = "reset"
	BMCReset PowerAction = "bmcreset"
)

type BMCType string

const (
	Dell      BMCType = "dell"
	HPE       BMCType = "hpe"
	LenovoXCC BMCType = "lenovoxcc"
	LenovoIMM BMCType = "lenovoimm"
	Nutanix   BMCType = "nutanix"
)

type LogType string

const (
	System     LogType = "system"
	Management LogType = "management"
	Fault      LogType = "fault"
)

type PowerRequest struct {
	BMCType   BMCType     `json:"bmc_type" binding:"omitempty"`  // Dell, HPE, Lenovo, Nutanix, etc.
	IPAddress string      `json:"ip_address" binding:"required"` // BMC IP address
	Action    PowerAction `json:"action" binding:"required"`     // Power action (on, off, reset, bmcreset)
}

// SystemInfoRequest defines the request body for querying system info
type SystemInfoRequest struct {
	BMCType   BMCType `json:"bmc_type" binding:"omitempty"`
	IPAddress string  `json:"ip_address" binding:"required"`
}

// LogsRequest defines the request body for querying system info
type LogsRequest struct {
	BMCType   BMCType `json:"bmc_type" binding:"omitempty"`
	LogType   LogType `json:"log_type" binding:"omitempty"`
	IPAddress string  `json:"ip_address" binding:"required"`
}

// FirmwareUpdateRequest defines the request body for firmware update operations
type FirmwareInfoRequest struct {
	BMCType   BMCType `json:"bmc_type" binding:"omitempty"`
	IPAddress string  `json:"ip_address" binding:"required"`
}

// FirmwareUpdateRequest defines the request body for firmware update operations
type FirmwareUpdateRequest struct {
	BMCType      BMCType `json:"bmc_type" binding:"omitempty"`
	IPAddress    string  `json:"ip_address" binding:"required"`
	FirmwarePath string  `json:"firmware_path" binding:"required"` // Path to the firmware file
}

// PowerResponse represents the success response.
type PowerResponse struct {
	Status  bool   `json:"status" example:"true"` // Status of the action
	Message string `json:"message" example:"Power action executed successfully"`
}

type SystemInfoResponse struct {
	Device         interface{} `json:"device,omitempty"`
	Health         string      `json:"health,omitempty"`
	Manufacturer   string      `json:"manufacturer,omitempty"`
	PowerState     interface{} `json:"powerstate"`
	Model          string      `json:"model,omitempty"`
	BiosVersion    string      `json:"biosversion,omitempty"`
	SerialNumber   string      `json:"serialnumber,omitempty"`
	HostName       interface{} `json:"hostname,omitempty"`
	ResponseStatus string      `json:"responsestatus,omitempty"`
	Memory         interface{} `json:"memory,omitempty"`
	CPU            interface{} `json:"cpu,omitempty"`
}

// ErrorResponse represents the error response structure.
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`            // HTTP status code
	Message string `json:"message" example:"Bad Request"` // Error message
}

func ValidatePowerAction(action PowerAction) error {
	validActions := []PowerAction{On, Off, Reset, BMCReset}
	for _, validAction := range validActions {
		if action == validAction {
			return nil
		}
	}
	return fmt.Errorf("invalid power action: %s", action)
}

func ValidateBMCType(bmcType BMCType) error {
	validBMCTypes := []BMCType{Dell, HPE, LenovoXCC, LenovoIMM, Nutanix}
	for _, validBMCType := range validBMCTypes {
		if bmcType == validBMCType {
			return nil
		}
	}
	return fmt.Errorf("invalid bmc type: %s", bmcType)
}
