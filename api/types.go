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

func validatePowerAction(action PowerAction) error {
	validActions := []PowerAction{On, Off, Reset, BMCReset}
	for _, validAction := range validActions {
		if action == validAction {
			return nil
		}
	}
	return fmt.Errorf("invalid power action: %s", action)
}

func validateBMCType(bmcType BMCType) error {
	validBMCTypes := []BMCType{Dell, HPE, LenovoXCC, LenovoIMM, Nutanix}
	for _, validBMCType := range validBMCTypes {
		if bmcType == validBMCType {
			return nil
		}
	}
	return fmt.Errorf("invalid bmc type: %s", bmcType)
}
