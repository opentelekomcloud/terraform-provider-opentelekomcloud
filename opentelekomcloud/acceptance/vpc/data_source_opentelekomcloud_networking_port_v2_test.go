package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	t.Parallel()
	qts := subnetQuotas()
	qts = append(qts, &quotas.ExpectedQuota{Q: quotas.SecurityGroup, Count: 1})
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.opentelekomcloud_networking_port_v2.port_1", "id",
						"opentelekomcloud_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.opentelekomcloud_networking_port_v2.port_2", "id",
						"opentelekomcloud_networking_port_v2.port_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.opentelekomcloud_networking_port_v2.port_3", "id",
						"opentelekomcloud_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_port_v2.port_3", "all_fixed_ips.#", "1"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortDataSourceBasic = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_ds_port"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1_ds_port"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  cidr       = "10.0.0.0/24"
  ip_version = 4
}

data "opentelekomcloud_networking_secgroup_v2" "default" {
  name = "default"
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_ds_1"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"

  security_group_ids = [
    data.opentelekomcloud_networking_secgroup_v2.default.id,
  ]

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name               = "port_ds_2"
  network_id         = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"

  security_group_ids = [
    data.opentelekomcloud_networking_secgroup_v2.default.id,
  ]

  allowed_address_pairs {
    ip_address  = "10.0.0.201"
    mac_address = "fa:16:3e:f8:ab:da"
  }

  allowed_address_pairs {
    ip_address  = "10.0.0.202"
    mac_address = "fa:16:3e:ab:4b:58"
  }
}

data "opentelekomcloud_networking_port_v2" "port_1" {
  name           = opentelekomcloud_networking_port_v2.port_1.name
  admin_state_up = opentelekomcloud_networking_port_v2.port_1.admin_state_up

  security_group_ids = [
    data.opentelekomcloud_networking_secgroup_v2.default.id,
  ]
}

data "opentelekomcloud_networking_port_v2" "port_2" {
  name           = opentelekomcloud_networking_port_v2.port_2.name
  admin_state_up = opentelekomcloud_networking_port_v2.port_2.admin_state_up
}

data "opentelekomcloud_networking_port_v2" "port_3" {
  fixed_ip = opentelekomcloud_networking_port_v2.port_1.all_fixed_ips.0
}
`
