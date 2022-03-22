package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackIdentityV3RoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3RoleDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DataSourceID("data.opentelekomcloud_identity_role_v3.role_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_role_v3.role_1", "name", "admin"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find role data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityV3RoleDataSource_basic = `
data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "admin"
}
`
