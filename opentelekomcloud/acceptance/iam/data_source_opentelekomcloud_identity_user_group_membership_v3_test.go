package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceUGMName = "data.opentelekomcloud_identity_user_group_membership_v3.membership_1"

func TestAccIdentityUserGroupMembershipV3DS_basic(t *testing.T) {
	var groupName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var userName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserGroupMembershipV3DSBasic(userName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserGroupMembershipExists(dataSourceUGMName, []string{groupName}),
					resource.TestCheckResourceAttr(dataSourceUGMName, "groups.0.name", groupName),
				),
			},
		},
	})
}

func testAccIdentityUserGroupMembershipV3DSBasic(userName, groupName string) string {
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

data "opentelekomcloud_identity_user_group_membership_v3" "membership_1" {
  user       = opentelekomcloud_identity_user_v3.user_1.id
  depends_on = [opentelekomcloud_identity_user_group_membership_v3.membership_1]
}
`, userName, groupName)
}
