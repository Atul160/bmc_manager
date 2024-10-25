package bmc

import (
	"errors"
	"os"
	"strings"
)

func loadEnvVars() map[string]string {
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		key := parts[0]
		value := parts[1]
		envVars[key] = value
	}
	return envVars
}

// NewBMCClient creates a BMC client based on the provided BMC type (Dell, HPE, Lenovo, Nutanix).
func NewBMCClient(bmcType, ipAddress string) (BMCClient, error) {
	envvars := loadEnvVars()
	switch bmcType {
	case "dell":
		return NewDellIDRACClient(ipAddress, envvars["IDRAC_USERNAME"], envvars["IDRAC_PASSWORD"]), nil
	case "hpe":
		return NewHPEILOClient(ipAddress, envvars["Stores_ILO_USERNAME"], envvars["Stores_ILO_PASSWORD"]), nil
	case "lenovoxcc":
		return NewLenovoXCCClient(ipAddress, envvars["XCC_USERNAME"], envvars["XCC_PASSWORD"]), nil
	case "lenovoimm":
		return NewLenovoIMMClient(ipAddress, envvars["IMM_USERNAME"], envvars["IMM_PASSWORD"]), nil
	case "nutanix":
		return NewNutanixIPMIClient(ipAddress, envvars["DC_IPMI_USERNAME"], envvars["DC_IPMI_PASSWORD"]), nil
	default:
		return nil, errors.New("unsupported BMC type")
	}
}
