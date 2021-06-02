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

func TestAccIdentityV3GroupMembership_basic(t *testing.T) {
	var groupName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var userName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var userName2 = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3GroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3GroupMembership_basic(groupName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupMembershipExists("opentelekomcloud_identity_group_membership_v3.membership_1", []string{userName}),
				),
			},
			{
				Config: testAccIdentityV3GroupMembership_update(groupName, userName, userName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupMembershipExists("opentelekomcloud_identity_group_membership_v3.membership_1", []string{userName, userName2}),
				),
			},
			{
				Config: testAccIdentityV3GroupMembership_updatedown(groupName, userName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupMembershipExists("opentelekomcloud_identity_group_membership_v3.membership_1", []string{userName2}),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupMembershipDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_group_membership_v3" {
			continue
		}

		_, err := users.ListInGroup(identityClient, rs.Primary.Attributes["group"], nil).AllPages()

		if err == nil {
			return fmt.Errorf("User still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3GroupMembershipExists(n string, us []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}
		group := rs.Primary.Attributes["group"]
		if group == "" {
			return fmt.Errorf("No group is set")
		}

		pages, err := users.ListInGroup(identityClient, group, nil).AllPages()
		if err != nil {
			return err
		}

		founds, err := users.ExtractUsers(pages)
		if err != nil {
			return err
		}

		uc := len(us)
		for _, u := range us {
			for _, f := range founds {
				if f.Name == u {
					uc--
				}
			}
		}

		if uc > 0 {
			return fmt.Errorf("Bad group membership compare, excepted(%d), but(%d)", len(us), len(founds))
		}

		return nil
	}
}

func testAccIdentityV3GroupMembership_basic(groupName, userName string) string {
	return fmt.Sprintf(`
    resource "opentelekomcloud_identity_group_v3" "group_1" {
      name = "%s"
    }

    resource "opentelekomcloud_identity_user_v3" "user_1" {
      name = "%s"
      password = "password123@#"
      enabled = true
    }

    resource "opentelekomcloud_identity_group_membership_v3" "membership_1" {
        group = opentelekomcloud_identity_group_v3.group_1.id
        users = [opentelekomcloud_identity_user_v3.user_1.id]
    }
  `, groupName, userName)
}

func testAccIdentityV3GroupMembership_update(groupName, userName string, userName2 string) string {
	return fmt.Sprintf(`
    resource "opentelekomcloud_identity_group_v3" "group_1" {
      name = "%s"
    }

    resource "opentelekomcloud_identity_user_v3" "user_1" {
      name = "%s"
      password = "password123@#"
      enabled = true
    }

    resource "opentelekomcloud_identity_user_v3" "user_2" {
      name = "%s"
      password = "password123@#"
      enabled = true
    }


    resource "opentelekomcloud_identity_group_membership_v3" "membership_1" {
        group = opentelekomcloud_identity_group_v3.group_1.id
        users = [opentelekomcloud_identity_user_v3.user_1.id,
                opentelekomcloud_identity_user_v3.user_2.id]
    }
  `, groupName, userName, userName2)
}

func testAccIdentityV3GroupMembership_updatedown(groupName, userName string) string {
	return fmt.Sprintf(`
    resource "opentelekomcloud_identity_group_v3" "group_1" {
      name = "%s"
    }

    resource "opentelekomcloud_identity_user_v3" "user_2" {
      name = "%s"
      password = "password123@#"
      enabled = true
    }


    resource "opentelekomcloud_identity_group_membership_v3" "membership_1" {
        group = opentelekomcloud_identity_group_v3.group_1.id
        users = [opentelekomcloud_identity_user_v3.user_2.id]
    }
  `, groupName, userName)
}
