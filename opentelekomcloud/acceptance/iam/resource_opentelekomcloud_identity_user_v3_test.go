package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceName = "opentelekomcloud_identity_user_v3.user_1"

func TestAccIdentityV3User_basic(t *testing.T) {
	var user users.User
	var userName = fmt.Sprintf("tf-user-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3UserBasic(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(resourceName, &user),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &user.Name),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email", "test@acme.org"),
				),
			},
			{
				Config: testAccIdentityV3UserUpdate(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(resourceName, &user),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &user.Name),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", "test2@acme.org"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcloud IdentityV3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_user_v3" {
			continue
		}

		_, err := users.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("user still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3UserExists(n string, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomcloud IdentityV3 client: %w", err)
		}

		found, err := users.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("user not found")
		}

		*user = *found

		return nil
	}
}

func testAccIdentityV3UserBasic(userName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "%s"
  password = "password123@!"
  enabled  = true
  email    = "test@acme.org"
}
  `, userName)
}

func testAccIdentityV3UserUpdate(userName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "%s"
  enabled  = false
  password = "password123@!"
  email    = "tEst2@acme.org"
}
  `, userName)
}
