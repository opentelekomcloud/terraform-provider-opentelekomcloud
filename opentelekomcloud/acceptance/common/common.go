package common

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const AlternativeProviderAlias = "opentelekomcloud.alternative"

var (
	TestAccProviderFactories map[string]func() (*schema.Provider, error)
	TestAccProvider          *schema.Provider

	altCloud                  = os.Getenv("OS_CLOUD_2")
	altProjectID              = os.Getenv("OS_PROJECT_ID_2")
	altProjectName            = os.Getenv("OS_PROJECT_NAME_2")
	AlternativeProviderConfig = fmt.Sprintf(`
provider opentelekomcloud {
  alias = "alternative"

  cloud       = "%s"
  tenant_id   = "%s"
  tenant_name = "%s"
}
`, altCloud, altProjectID, altProjectName)
)

func init() {
	TestAccProvider = opentelekomcloud.Provider()
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"opentelekomcloud": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
		"opentelekomcloudalternative": func() (*schema.Provider, error) {
			provider := opentelekomcloud.Provider()
			provider.Configure(context.Background(), &terraform.ResourceConfig{
				Config: map[string]interface{}{
					"cloud":       altCloud,
					"tenant_id":   altProjectID,
					"tenant_name": altProjectName,
				},
			})
			return provider, nil
		},
	}

	err := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err == nil {
		config := TestAccProvider.Meta().(*cfg.Config)
		env.OsRegionName = config.GetRegion(nil)
	}
}

func TestAccPreCheckRequiredEnvVars(t *testing.T) {
	if env.OsPoolName == "" {
		t.Fatal("OS_POOL_NAME must be set for acceptance tests")
	}

	if env.OsRegionName == "" {
		t.Fatal("OS_TENANT_NAME or OS_PROJECT_NAME must be set for acceptance tests")
	}

	if env.OsFlavorID == "" && env.OsFlavorName == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if env.OsNetworkID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
	}

	if env.OsRouterID == "" {
		t.Fatal("OS_VPC_ID must be set for acceptance tests")
	}

	if env.OsAvailabilityZone == "" {
		t.Fatal("OS_AVAILABILITY_ZONE must be set for acceptance tests")
	}

	if env.OsSubnetID == "" {
		t.Fatal("OS_SUBNET_ID must be set for acceptance tests")
	}

	if env.OsExtGwID == "" {
		t.Fatal("OS_EXTGW_ID must be set for acceptance tests")
	}
}

func TestAccPreCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)

	// Do not run the test if this is a deprecated testing environment.
	if env.OsDeprecatedEnvironment != "" {
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
	if env.OsFlavorID == "" {
		t.Skip("OS_FLAVOR_ID must be set for acceptance tests")
	}
}

func TestAccVBSBackupShareCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)
	if env.OsToTenantID == "" {
		t.Skip("OS_TO_TENANT_ID must be set for acceptance tests")
	}
}
