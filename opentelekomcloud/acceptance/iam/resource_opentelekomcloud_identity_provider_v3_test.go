package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const providerResource = "opentelekomcloud_identity_provider_v3"

func TestAccIdentityV3ProviderBasic(t *testing.T) {
	fullName := fmt.Sprintf("%s.%s", providerResource, "provider")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckIdentityV3ProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProviderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProviderDestroy,
					resource.TestCheckResourceAttr(fullName, "enabled", "true"),
					resource.TestCheckResourceAttr(fullName, "description", providerDescription),
					resource.TestCheckResourceAttr(fullName, "links.%", "2"),
				),
			},
			{
				Config: testAccIdentityV3ProviderUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProviderDestroy,
					resource.TestCheckResourceAttr(fullName, "enabled", "false"),
					resource.TestCheckResourceAttr(fullName, "description", providerDescriptionUpdated),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProviderDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity v3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != providerResource {
			continue
		}

		_, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("provider still exists")
		}
	}

	return nil
}

var (
	providerName = tools.RandomString("tf-test-", 4)

	providerDescription        = tools.RandomString("Provider for ", 20)
	providerDescriptionUpdated = tools.RandomString("Updated provider for ", 20)

	testAccIdentityV3ProviderBasic = fmt.Sprintf(`
resource "opentelekomcloud_identity_provider_v3" "provider" {
  name        = "%s"
  description = "%s"
  enabled     = true
}
`, providerName, providerDescription)

	testAccIdentityV3ProviderUpdated = fmt.Sprintf(`
resource "opentelekomcloud_identity_provider_v3" "provider" {
  name        = "%s"
  description = "%s"
  enabled     = false
}
`, providerName, providerDescriptionUpdated)
)
