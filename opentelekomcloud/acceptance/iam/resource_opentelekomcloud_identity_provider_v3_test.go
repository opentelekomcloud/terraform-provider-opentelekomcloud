package acceptance

import (
	"fmt"
	"os"
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
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProviderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProviderDestroy,
					resource.TestCheckResourceAttr(fullName, "enabled", "true"),
					resource.TestCheckResourceAttr(fullName, "description", providerDescription),
					resource.TestCheckResourceAttr(fullName, "links.%", "2"),
					// sso_type is omitted from the config, so the computed value must
					// reflect the API default rather than an empty string.
					resource.TestCheckResourceAttr(fullName, "sso_type", "virtual_user_sso"),
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

func TestAccIdentityV3ProviderSSOType(t *testing.T) {
	fullName := fmt.Sprintf("%s.%s", providerResource, "provider")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProviderSSOType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProviderDestroy,
					resource.TestCheckResourceAttr(fullName, "sso_type", "virtual_user_sso"),
				),
			},
		},
	})
}

// TestAccIdentityV3ProviderIAMUserSSO verifies that an explicit non-default
// sso_type ("iam_user_sso") is actually transmitted and overrides the API
// default ("virtual_user_sso"). It is opt-in: iam_user_sso requires a
// pre-existing IAM user and is limited to one provider of this type per
// account, so it must not run in the default acceptance suite.
func TestAccIdentityV3ProviderIAMUserSSO(t *testing.T) {
	fullName := fmt.Sprintf("%s.%s", providerResource, "provider")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
			if os.Getenv("OS_IAM_USER_SSO") == "" {
				t.Skip("OS_IAM_USER_SSO must be set (account has an IAM user and a free iam_user_sso slot) to run this test")
			}
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProviderIAMUserSSO,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProviderDestroy,
					resource.TestCheckResourceAttr(fullName, "sso_type", "iam_user_sso"),
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

	testAccIdentityV3ProviderSSOType = fmt.Sprintf(`
resource "opentelekomcloud_identity_provider_v3" "provider" {
  name        = "%s"
  description = "%s"
  enabled     = true
  sso_type    = "virtual_user_sso"
}
`, providerName, providerDescription)

	testAccIdentityV3ProviderIAMUserSSO = fmt.Sprintf(`
resource "opentelekomcloud_identity_provider_v3" "provider" {
  name        = "%s"
  description = "%s"
  enabled     = true
  sso_type    = "iam_user_sso"
}
`, providerName, providerDescription)
)
