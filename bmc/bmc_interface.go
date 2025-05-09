package bmc

// BMCClient is a common interface that all BMC clients must implement
type BMCClient interface {
	Connect() (map[string]string, error)
	BMCReset() error
	SetPower(string) error
	GetSystemInfo() (map[string]interface{}, error)
	GetLogs(string) ([]map[string]interface{}, error)
	GetFirmwareInfo() (map[string]interface{}, error)
	UpdateFirmware(filePath string) error
}
