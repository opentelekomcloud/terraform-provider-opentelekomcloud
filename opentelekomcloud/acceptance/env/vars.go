package env

import (
	"os"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var (
	OS_EXTGW_ID          = os.Getenv("OS_EXTGW_ID")
	OS_FLAVOR_ID         = os.Getenv("OS_FLAVOR_ID")
	OS_FLAVOR_NAME       = os.Getenv("OS_FLAVOR_NAME")
	OS_IMAGE_ID          = os.Getenv("OS_IMAGE_ID")
	OsImageName          = imageName()
	OS_NETWORK_ID        = os.Getenv("OS_NETWORK_ID")
	OsNetworkName        = os.Getenv("OS_NETWORK_NAME")
	OsExtNetworkName     = extNetworkName()
	OS_REGION_NAME       string
	OS_ACCESS_KEY        = os.Getenv("OS_ACCESS_KEY")
	OS_SECRET_KEY        = os.Getenv("OS_SECRET_KEY")
	OS_AVAILABILITY_ZONE = os.Getenv("OS_AVAILABILITY_ZONE")
	OS_VPC_ID            = os.Getenv("OS_VPC_ID")
	OsRouterName         = os.Getenv("OS_ROUTER_NAME")
	OS_SUBNET_ID         = os.Getenv("OS_SUBNET_ID")
	OsSubnetName         = os.Getenv("OS_SUBNET_NAME")
	OS_KEYPAIR_NAME      = os.Getenv("OS_KEYPAIR_NAME")
	OS_KMS_ID            = os.Getenv("OS_KMS_ID")
	OsKmsName            = os.Getenv("OS_KMS_NAME")
	OS_BMS_FLAVOR_NAME   = os.Getenv("OS_BMS_FLAVOR_NAME")
	OS_TO_TENANT_ID      = os.Getenv("OS_TO_TENANT_ID")
	OS_TENANT_NAME       = GetTenantName()
	OS_TENANT_ID         = os.Getenv("OS_TENANT_ID")
)

func extNetworkName() string {
	if nw := os.Getenv("OS_EXT_NETWORK_NAME"); nw != "" {
		return nw
	}
	return "admin_external_net" // value valid for OTC PROD, both eu-de and eu-nl
}

func imageName() string {
	if image := os.Getenv("OS_IMAGE_NAME"); image != "" {
		return image
	}
	return "Standard_Debian_10_latest" // value valid for OTC PROD, both eu-de and eu-nl
}

func GetTenantName() cfg.ProjectName {
	tn := os.Getenv("OS_TENANT_NAME")
	if tn == "" {
		tn = os.Getenv("OS_PROJECT_NAME")
	}
	return cfg.ProjectName(tn)
}
