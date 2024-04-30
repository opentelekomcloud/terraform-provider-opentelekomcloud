package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestAccNetworkingV2PortIDsDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_networking_port_ids_v2.ports"
	port1Name := "opentelekomcloud_networking_port_v2.port_1"
	port2Name := "opentelekomcloud_networking_port_v2.port_2"
	t.Parallel()
	quotas.BookOne(t, quotas.SecurityGroup)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortIDsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "2"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ids.0", port1Name, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ids.1", port2Name, "id"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortIDsDataSourceBasic = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "acc_network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_secgroup_v2" "sg_1" {
  name        = "acc_secgroup_1"
  description = "acc_secgroup_1"
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"

  security_group_ids = [
    opentelekomcloud_networking_secgroup_v2.sg_1.id
  ]
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"

  security_group_ids = [
    opentelekomcloud_networking_secgroup_v2.sg_1.id
  ]
}

data "opentelekomcloud_networking_port_ids_v2" "ports" {
  sort_direction = "asc"
  sort_key       = "name"

  network_id = opentelekomcloud_networking_network_v2.network_1.id
  depends_on = [
    opentelekomcloud_networking_port_v2.port_1,
    opentelekomcloud_networking_port_v2.port_2,
  ]
}
`
