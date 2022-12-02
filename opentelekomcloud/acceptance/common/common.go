package common

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/catalog"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const (
	AlternativeProviderAlias           = "opentelekomcloud.alternative"
	AlternativeProviderWithRegionAlias = "opentelekomcloud.region"
)

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
	AlternativeProviderWithRegionConfig string
)

func init() {
	TestAccProvider = opentelekomcloud.Provider()

	err := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err == nil {
		config := TestAccProvider.Meta().(*cfg.Config)
		env.OS_REGION_NAME = config.GetRegion(nil)
	}

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
		"opentelekomcloudregion": func() (*schema.Provider, error) {
			provider := opentelekomcloud.Provider()
			provider.Configure(context.Background(), &terraform.ResourceConfig{
				Config: map[string]interface{}{
					"cloud":  altCloud,
					"region": env.OS_REGION_NAME,
				},
			})
			return provider, nil
		},
	}

	AlternativeProviderWithRegionConfig = fmt.Sprintf(`
provider opentelekomcloud {
  alias = "region"

  region = "%s"
}
`, env.OS_REGION_NAME)
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

func TestAccPreCheckServiceAvailability(t *testing.T, service string, regions []string) diag.Diagnostics {
	t.Logf("Service: %s, Region %s", service, env.OS_REGION_NAME)
	config := TestAccProvider.Meta().(*cfg.Config)
	client, err := config.RegionIdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmterr.Errorf("clientCreationFail", err)
	}
	allPages, err := catalog.List(client).AllPages()
	if err != nil {
		return fmterr.Errorf("error fetching service catalog: %s", err)
	}
	allServices, err := catalog.ExtractServiceCatalog(allPages)
	if err != nil {
		return fmterr.Errorf("error fetching services from catalog: %s", err)
	}
	for _, entry := range allServices {
		// if found in service catalog then ok
		if service == entry.Name {
			for _, region := range regions {
				if env.OS_REGION_NAME == region {
					return nil
				}
			}
		}
	}

	t.Skipf("Service %s not available or configuration differs in %s", service, env.OS_REGION_NAME)
	return fmterr.Errorf("test not valid in region")
}
