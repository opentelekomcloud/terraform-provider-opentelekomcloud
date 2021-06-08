package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2RouterInterface_basic_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_basic_subnet,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.network_1", &network),
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
					TestAccCheckNetworkingV2RouterInterfaceExists("opentelekomcloud_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_basic_port(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_basic_port,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.network_1", &network),
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.port_1", &port),
					TestAccCheckNetworkingV2RouterInterfaceExists("opentelekomcloud_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_timeout(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_timeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.network_1", &network),
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
					TestAccCheckNetworkingV2RouterInterfaceExists("opentelekomcloud_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterInterfaceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_router_interface_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router interface still exists")
		}
	}

	return nil
}

const testAccNetworkingV2RouterInterface_basic_subnet = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2RouterInterface_basic_port = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}
`

const testAccNetworkingV2RouterInterface_timeout = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`
