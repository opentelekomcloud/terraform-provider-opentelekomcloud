package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackIdentityV3GroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3GroupDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DataSourceID("data.opentelekomcloud_identity_group_v3.group_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_group_v3.group_1", "name", "admins"),
				),
			},
		},
	})
}

const testAccOpenStackIdentityV3GroupDataSource_basic = `
data "opentelekomcloud_identity_group_v3" "group_1" {
  name = "admins"
}
`
