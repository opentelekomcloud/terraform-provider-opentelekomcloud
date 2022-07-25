package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

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

	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceBasicSubnet,
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

	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceBasicPort,
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

	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceTimeout,
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

func TestAccNetworkingV2RouterInterface_port(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceDeleteHangingPort,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.network_1", &network),
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
					TestAccCheckNetworkingV2RouterInterfaceExists("opentelekomcloud_networking_router_interface_v2.interface_1"),
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
			return fmt.Errorf("router interface still exists")
		}
	}

	return nil
}

const testAccNetworkingV2RouterInterfaceBasicSubnet = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_ri_sn"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_ri_sn"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2RouterInterfaceBasicPort = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_bp"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_ri_bp"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1_ri_bp"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}
`

const testAccNetworkingV2RouterInterfaceTimeout = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_ri_t"
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
  name           = "network_1_ri_t"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2RouterInterfaceDeleteHangingPort = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr       = "172.22.34.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "interface_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name       = "instance_1"
  image_name = "Standard_Debian_10_latest"

  network {
    port = opentelekomcloud_networking_port_v2.instance_port.id
  }
}

resource "opentelekomcloud_networking_port_v2" "instance_port" {
  name           = "port_1"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = true

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "172.22.34.120"
  }
}
`
