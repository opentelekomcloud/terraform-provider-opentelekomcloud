package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2Subnet_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_basic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.100"),
				),
			},
			{
				Config: testAccNetworkingV2Subnet_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "gateway_ip", "192.168.199.1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "enable_dhcp", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.150"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_enableDHCP(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_enableDHCP,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "enable_dhcp", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_noGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_noGateway,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "gateway_ip", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_impliedGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_impliedGateway,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_subnet_v2.subnet_1", "gateway_ip", "192.168.199.1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_timeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SubnetDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_subnet_v2" {
			continue
		}

		_, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}

const testAccNetworkingV2Subnet_basic = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = opentelekomcloud_networking_network_v2.network_1.id

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_update = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  network_id = opentelekomcloud_networking_network_v2.network_1.id

  allocation_pools {
    start = "192.168.199.150"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_enableDHCP = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  enable_dhcp = true
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2Subnet_disableDHCP = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  enable_dhcp = false
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2Subnet_noGateway = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2Subnet_impliedGateway = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2Subnet_timeout = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = opentelekomcloud_networking_network_v2.network_1.id

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
