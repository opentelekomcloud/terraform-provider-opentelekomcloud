package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackIdentityV3CustomRoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3CustomRoleDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DataSourceID("data.opentelekomcloud_identity_role_custom_v3.role_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_role_custom_v3.role_1", "display_name", "custom_role-terraform"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_role_custom_v3.role_1", "description", "role"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_role_custom_v3.role_1", "type", "domain"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_role_custom_v3.role_1", "statement.0.effect", "Allow"),
				),
			},
		},
	})
}

const testAccOpenStackIdentityV3CustomRoleDataSource_basic = `
data "opentelekomcloud_identity_role_custom_v3" "role_1" {
  display_name = opentelekomcloud_identity_role_v3.role.display_name
}

resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "role"
  display_name  = "custom_role-terraform"
  display_layer = "domain"
  statement {
    effect = "Allow"
    action = ["ecs:*:list*"]
  }
}
`

func TestAccOpenStackIdentityV3CustomRoleDataSource_agencyResource(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_identity_role_custom_v3.role_1"
	name := "custom-role" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3CustomRoleDataSource_agencyResource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "display_name", name),
					resource.TestCheckResourceAttr(dataSourceName, "statement.0.action.0", "iam:agencies:assume"),
					resource.TestCheckResourceAttr(dataSourceName, "statement.0.resource.#", "1"),
					testAccCheckIdentityRoleV3AgencyURI(dataSourceName, "statement.0.resource.0",
						"opentelekomcloud_identity_agency_v3.agency_1"),
				),
			},
		},
	})
}

func testAccOpenStackIdentityV3CustomRoleDataSource_agencyResource(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_identity_role_custom_v3" "role_1" {
  display_name = opentelekomcloud_identity_role_v3.role.display_name
}
`, testAccIdentityRoleV3_agencyResource(name))
}
