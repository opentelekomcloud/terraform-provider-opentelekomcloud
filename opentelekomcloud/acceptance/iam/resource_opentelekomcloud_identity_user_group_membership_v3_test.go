package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceUGMName = "opentelekomcloud_identity_user_group_membership_v3.membership_1"

func TestAccIdentityUserGroupMembershipV3_basic(t *testing.T) {
	var groupName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var groupName2 = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var userName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3UserGroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserGroupMembershipV3Basic(userName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserGroupMembershipExists(resourceUGMName, []string{groupName}),
				),
			},
			{
				Config: testAccIdentityUserGroupMembershipV3Update(userName, groupName, groupName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserGroupMembershipExists(resourceUGMName, []string{groupName, groupName2}),
				),
			},
			{
				Config: testAccIdentityUserGroupMembershipV3UpdateDown(userName, groupName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserGroupMembershipExists(resourceUGMName, []string{groupName2}),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserGroupMembershipDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_group_membership_v3" {
			continue
		}

		_, err := users.ListGroups(identityClient, rs.Primary.ID).AllPages()

		if err == nil {
			return fmt.Errorf("user still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3UserGroupMembershipExists(n string, us []string) resource.TestCheckFunc {
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
			return fmt.Errorf("error creating OpenStack identity client: %s", err)
		}
		user := rs.Primary.ID
		if user == "" {
			return fmt.Errorf("no ID is set")
		}

		pages, err := users.ListGroups(identityClient, user).AllPages()
		if err != nil {
			return err
		}

		founds, err := groups.ExtractGroups(pages)
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
			return fmt.Errorf("bad group membership compare, excepted \n%+v\nbut found\n%+v)", us, founds)
		}

		return nil
	}
}

func testAccIdentityUserGroupMembershipV3Basic(userName, groupName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "%s"
}


resource "opentelekomcloud_identity_user_group_membership_v3" "membership_1" {
  user = opentelekomcloud_identity_user_v3.user_1.id
  groups = [
    opentelekomcloud_identity_group_v3.group_1.id,
  ]
}
  `, userName, groupName)
}

func testAccIdentityUserGroupMembershipV3Update(userName, groupName, groupName2 string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "%s"
}

resource "opentelekomcloud_identity_group_v3" "group_2" {
  name = "%s"
}

resource "opentelekomcloud_identity_user_group_membership_v3" "membership_1" {
  user = opentelekomcloud_identity_user_v3.user_1.id
  groups = [
    opentelekomcloud_identity_group_v3.group_1.id,
    opentelekomcloud_identity_group_v3.group_2.id,
  ]
}
  `, userName, groupName, groupName2)
}

func testAccIdentityUserGroupMembershipV3UpdateDown(groupName, userName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}

resource "opentelekomcloud_identity_group_v3" "group_2" {
  name = "%s"
}

resource "opentelekomcloud_identity_user_group_membership_v3" "membership_1" {
  user = opentelekomcloud_identity_group_v3.group_1.id
  groups = [
    opentelekomcloud_identity_group_v3.group_2.id,
  ]
}
  `, groupName, userName)
}
