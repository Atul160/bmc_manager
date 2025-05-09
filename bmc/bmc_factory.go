package bmc

import (
	"ecc-bmc/utils"
	"errors"
	"strings"
)

// NewBMCClient creates a BMC client based on the provided BMC type (Dell, HPE, Lenovo, Nutanix).
func NewBMCClient(bmcType, ipAddress string) (BMCClient, error) {
	envvars := utils.LoadEnvVars()
	fqdn, err := utils.ResolveDNS(ipAddress)
	if err != nil || fqdn == "" {
		// fmt.Print(err)
	} else {
		ipAddress = fqdn
	}

	switch bmcType {
	case "dell":
		if strings.HasPrefix(ipAddress, "vsrv") {
			return NewDellIDRACClient(ipAddress, envvars["vertical1_WIN_IDRAC_USERNAME"], envvars["vertical1_WIN_IDRAC_PASSWORD"]), nil
		} else {
			return NewDellIDRACClient(ipAddress, envvars["vertical1_VM_IDRAC_USERNAME"], envvars["vertical1_VM_IDRAC_PASSWORD"]), nil
		}
	case "hpe":
		if strings.HasPrefix(ipAddress, "nlc") || strings.HasPrefix(ipAddress, "vhst") || strings.HasPrefix(ipAddress, "vsh") {
			return NewHPEILOClient(ipAddress, envvars["vertical1_ILO_USERNAME"], envvars["vertical1_ILO_PASSWORD"]), nil
		} else {
			return NewHPEILOClient(ipAddress, envvars["vertical2_ILO_USERNAME"], envvars["vertical2_ILO_PASSWORD"]), nil
		}
	case "lenovoxcc":
		return NewLenovoXCCClient(ipAddress, envvars["XCC_USERNAME"], envvars["XCC_PASSWORD"]), nil
	case "lenovoimm":
		return NewLenovoIMMClient(ipAddress, envvars["IMM_USERNAME"], envvars["IMM_PASSWORD"]), nil
	case "nutanix":
		return NewNutanixIPMIClient(ipAddress, envvars["vertical2_IPMI_USERNAME"], envvars["vertical2_IPMI_PASSWORD"]), nil
	default:
		return nil, errors.New("unsupported BMC type " + bmcType)
	}
}
