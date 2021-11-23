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
		env.OS_REGION_NAME = config.GetRegion(nil)
	}
}

func TestAccPreCheckRequiredEnvVars(t *testing.T) {
	if env.OS_REGION_NAME == "" {
		t.Skip("OS_TENANT_NAME or OS_PROJECT_NAME must be set for acceptance tests")
	}

	if env.OS_AVAILABILITY_ZONE == "" {
		t.Skip("OS_AVAILABILITY_ZONE must be set for acceptance tests")
	}

	if env.OsSubnetName == "" {
		t.Skip("OS_SUBNET_NAME must be set for acceptance tests")
	}
}

func TestAccPreCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)
}

func TestAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_TENANT_ADMIN")
	if v == "" {
		t.Skip("Skipping test because it requires set OS_TENANT_ADMIN")
	}
}

func TestAccVBSBackupShareCheck(t *testing.T) {
	TestAccPreCheckRequiredEnvVars(t)
	if env.OS_TO_TENANT_ID == "" {
		t.Skip("OS_TO_TENANT_ID must be set for acceptance tests")
	}
}
