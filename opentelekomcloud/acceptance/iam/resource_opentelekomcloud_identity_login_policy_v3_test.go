package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/security"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceLoginPolicyName = "opentelekomcloud_identity_login_policy_v3.pol_1"

func TestAccIdentityV3Login_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3LoginPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3LoginPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3LoginPolicyExists(resourceLoginPolicyName),
					resource.TestCheckResourceAttr(resourceLoginPolicyName, "session_timeout", "1396"),
				),
			},
			{
				Config: testAccIdentityV3LoginPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3LoginPolicyExists(resourceLoginPolicyName),
					resource.TestCheckResourceAttr(resourceLoginPolicyName, "session_timeout", "1395"),
				),
			},
			{
				ResourceName:      resourceLoginPolicyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIdentityV3LoginPolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcloud IdentityV3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_login_policy_v3" {
			continue
		}

		_, err := security.GetLoginAuthPolicy(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching the IAM account login policy")
		}
	}

	return nil
}

func testAccCheckIdentityV3LoginPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.IdentityV30Client()
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomcloud IdentityV3 client: %w", err)
		}

		_, err = security.GetLoginAuthPolicy(client, rs.Primary.ID)
		return err
	}
}

const testAccIdentityV3LoginPolicyBasic = `
resource "opentelekomcloud_identity_login_policy_v3" "pol_1" {
  custom_info_for_login      = ""
  period_with_login_failures = 60
  lockout_duration           = 15
  account_validity_period    = 0
  login_failed_times         = 3
  session_timeout            = 1396
  show_recent_login_info     = false
}
`

const testAccIdentityV3LoginPolicyUpdate = `
resource "opentelekomcloud_identity_login_policy_v3" "pol_1" {
  custom_info_for_login      = ""
  period_with_login_failures = 60
  lockout_duration           = 15
  account_validity_period    = 0
  login_failed_times         = 3
  session_timeout            = 1395
  show_recent_login_info     = false
}
`
