package common

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var (
	TestAccProviderFactories map[string]func() (*schema.Provider, error)
	TestAccProvider          *schema.Provider
)

func init() {
	TestAccProvider = opentelekomcloud.Provider()
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"opentelekomcloud": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}

	err := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err == nil {
		config := TestAccProvider.Meta().(*cfg.Config)
		env.OS_REGION_NAME = config.GetRegion(nil)
	}
}

func TestAccPreCheckRequiredEnvVars(t *testing.T) {
	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}

	if env.OS_POOL_NAME == "" {
		t.Fatal("OS_POOL_NAME must be set for acceptance tests")
	}

	if env.OS_REGION_NAME == "" {
		t.Fatal("OS_TENANT_NAME or OS_PROJECT_NAME must be set for acceptance tests")
	}

	if env.OS_FLAVOR_ID == "" && env.OS_FLAVOR_NAME == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if env.OS_NETWORK_ID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
	}

	if env.OS_VPC_ID == "" {
		t.Fatal("OS_VPC_ID must be set for acceptance tests")
	}

	if env.OS_AVAILABILITY_ZONE == "" {
		t.Fatal("OS_AVAILABILITY_ZONE must be set for acceptance tests")
	}

	if env.OS_SUBNET_ID == "" {
		t.Fatal("OS_SUBNET_ID must be set for acceptance tests")
	}

	if env.OS_EXTGW_ID == "" {
		t.Fatal("OS_EXTGW_ID must be set for acceptance tests")
	}

}

func TestAccPreCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)

	// Do not run the test if this is a deprecated testing environment.
	if env.OS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}
}

func TestAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_TENANT_ADMIN")
	if v == "" {
		t.Skip("Skipping test because it requires set OS_TENANT_ADMIN")
	}
}

func TestAccFlavorPreCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)
	if env.OS_FLAVOR_ID == "" {
		t.Skip("OS_FLAVOR_ID must be set for acceptance tests")
	}
}

func TestAccVBSBackupShareCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)
	if env.OS_TO_TENANT_ID == "" {
		t.Skip("OS_TO_TENANT_ID must be set for acceptance tests")
	}

}
