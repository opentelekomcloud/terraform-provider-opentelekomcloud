package env

import (
	"os"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var (
	OsDeprecatedEnvironment = os.Getenv("OS_DEPRECATED_ENVIRONMENT")
	OsExtGwID               = os.Getenv("OS_EXTGW_ID")
	OsFlavorID              = os.Getenv("OS_FLAVOR_ID")
	OsFlavorName            = os.Getenv("OS_FLAVOR_NAME")
	OsImageID               = os.Getenv("OS_IMAGE_ID")
	OsPoolName              = os.Getenv("OS_POOL_NAME")
	OsRegionName            string
	OsAccessKey             = os.Getenv("OS_ACCESS_KEY")
	OsSecretKey             = os.Getenv("OS_SECRET_KEY")
	OsMrsEnvironment        = os.Getenv("OS_MRS_ENVIRONMENT")
	OsDmsEnvironment        = os.Getenv("OS_DMS_ENVIRONMENT")
	OsAvailabilityZone      = os.Getenv("OS_AVAILABILITY_ZONE")
	OsRouterID              = os.Getenv("OS_ROUTER_ID")
	OsNetworkID             = os.Getenv("OS_NETWORK_ID")
	OsSubnetID              = os.Getenv("OS_SUBNET_ID")
	OsKeypairName           = os.Getenv("OS_KEYPAIR_NAME")
	OsKmsID                 = os.Getenv("OS_KMS_ID")
	OsBmsFlavorName         = os.Getenv("OS_BMS_FLAVOR_NAME")
	OsNicID                 = os.Getenv("OS_NIC_ID")
	OsToTenantID            = os.Getenv("OS_TO_TENANT_ID")
	OsTenantName            = GetTenantName()
	OsTenantID              = os.Getenv("OS_TENANT_ID")
)

func GetTenantName() cfg.ProjectName {
	tn := os.Getenv("OS_TENANT_NAME")
	if tn == "" {
		tn = os.Getenv("OS_PROJECT_NAME")
	}
	return cfg.ProjectName(tn)
}
