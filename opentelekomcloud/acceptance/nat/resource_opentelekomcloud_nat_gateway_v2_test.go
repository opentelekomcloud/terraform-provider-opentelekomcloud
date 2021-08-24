package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/natgateways"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	vpc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/vpc"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNatGateway_basic(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2Gateway_basic,
				Check: resource.ComposeTestCheckFunc(
					vpc.TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.network_1", &network),
					vpc.TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					vpc.TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
					vpc.TestAccCheckNetworkingV2RouterInterfaceExists("opentelekomcloud_networking_router_interface_v2.int_1"),
					testAccCheckNatV2GatewayExists("opentelekomcloud_nat_gateway_v2.nat_1"),
				),
			},
			{
				Config: testAccNatV2Gateway_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_nat_gateway_v2.nat_1", "name", "nat_1_updated"),
					resource.TestCheckResourceAttr("opentelekomcloud_nat_gateway_v2.nat_1", "description", "nat_1 updated"),
					resource.TestCheckResourceAttr("opentelekomcloud_nat_gateway_v2.nat_1", "spec", "2"),
				),
			},
		},
	})
}

func testAccCheckNatV2GatewayDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	natClient, err := config.NatV2Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_nat_gateway_v2" {
			continue
		}

		_, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("nat gateway still exists")
		}
	}

	return nil
}

func testAccCheckNatV2GatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		natClient, err := config.NatV2Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud nat client: %s", err)
		}

		found, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("nat gateway not found")
		}

		return nil
	}
}

const testAccNatV2Gateway_basic = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
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

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name   = "nat_1"
  description = "test for terraform"
  spec = "1"
  internal_network_id = opentelekomcloud_networking_network_v2.network_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_1"]
}
`

const testAccNatV2Gateway_update = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
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

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name   = "nat_1_updated"
  description = "nat_1 updated"
  spec = "2"
  internal_network_id = opentelekomcloud_networking_network_v2.network_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_1"]
}
`
