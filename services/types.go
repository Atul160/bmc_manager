package services

type SystemInfo struct {
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

type FirmeareInfo struct {
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