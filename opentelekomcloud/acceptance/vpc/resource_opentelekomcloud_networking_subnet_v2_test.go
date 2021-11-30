package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceNwSubnetName = "opentelekomcloud_networking_subnet_v2.subnet_1"

func TestAccNetworkingV2Subnet_basic(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	quotas.BookMany(t, subnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists(resourceNwSubnetName, &subnet),
					resource.TestCheckResourceAttr("opentelekomcloud_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.100"),
				),
			},
			{
				Config: testAccNetworkingV2SubnetUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNwSubnetName, "name", "subnet_1"),
					resource.TestCheckResourceAttr(resourceNwSubnetName, "gateway_ip", "192.168.199.1"),
					resource.TestCheckResourceAttr(resourceNwSubnetName, "enable_dhcp", "true"),
					resource.TestCheckResourceAttr(resourceNwSubnetName, "allocation_pools.0.start", "192.168.199.150"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_import(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, subnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetImport,
			},

			{
				ResourceName:      resourceNwSubnetName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2Subnet_enableDHCP(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	quotas.BookMany(t, subnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetEnableDHCP,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(resourceNwSubnetName, "enable_dhcp", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_noGateway(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	quotas.BookMany(t, subnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetNoGateway,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(resourceNwSubnetName, "gateway_ip", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_impliedGateway(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	quotas.BookMany(t, subnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetImpliedGateway,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(resourceNwSubnetName, "gateway_ip", "192.168.199.1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_timeout(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	quotas.BookMany(t, subnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetTimeout,
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
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_subnet_v2" {
			continue
		}

		_, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("subnet still exists")
		}
	}

	return nil
}

const testAccNetworkingV2SubnetBasic = `
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
const testAccNetworkingV2SubnetImport = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1_i"
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

const testAccNetworkingV2SubnetUpdate = `
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

const testAccNetworkingV2SubnetEnableDHCP = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1_dhcp"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1_dhcp"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  enable_dhcp = true
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2SubnetNoGateway = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1_no_gw"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1_no_gw"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2SubnetImpliedGateway = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1_i_gw"
  admin_state_up = "true"
}
resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1_i_gw"
  cidr = "192.168.199.0/24"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
`

const testAccNetworkingV2SubnetTimeout = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1_t"
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
