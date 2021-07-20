package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccNetworkingSecGroupRuleIdsV2DataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupIdsV2DataSource_sg,
			},
			{
				Config: testAccNetworkingSecGroupRuleIdsV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "2"),
				),
			},
		},
	})
}

const testAccNetworkingSecGroupIdsV2DataSource_sg = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}
`

var testAccNetworkingSecGroupRuleIdsV2DataSourceBasic = fmt.Sprintf(`
%s
data "opentelekomcloud_networking_secgroup_rule_ids_v2" "secgroup_ids" {
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`, testAccNetworkingSecGroupIdsV2DataSource_sg)
