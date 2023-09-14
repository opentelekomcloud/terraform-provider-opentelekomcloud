package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

type testPortWithExtensions struct {
	ports.Port
	portsecurity.PortSecurityExt
}

const resourceNwPortName = "opentelekomcloud_networking_port_v2.port_1"

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists(resourceNwSubnetName, &subnet),
					TestAccCheckNetworkingV2NetworkExists(resourceNwNetworkName, &network),
					testAccCheckNetworkingV2PortExists(resourceNwPortName, &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_importBasic(t *testing.T) {
	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortBasic,
			},

			{
				ResourceName:      resourceNwPortName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingV2Port_noip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet
	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortNoIP,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists(resourceNwSubnetName, &subnet),
					TestAccCheckNetworkingV2NetworkExists(resourceNwNetworkName, &network),
					testAccCheckNetworkingV2PortExists(resourceNwPortName, &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_allowedAddressPairs(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrpPort1, vrrpPort2, instancePort ports.Port
	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAllowedAddressPairs,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.vrrp_subnet", &subnet),
					TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.instance_port", &instancePort),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_portSecurity_enabled(t *testing.T) {
	var port testPortWithExtensions
	resourceName := "opentelekomcloud_networking_port_v2.port_1"
	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists(resourceName, &port),
					resource.TestCheckResourceAttr(resourceName, "port_security_enabled", "true"),
					testAccCheckNetworkingV2PortPortSecurity(&port, true),
				),
			},
			{
				Config: testAccNetworkingV2PortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists(resourceName, &port),
					resource.TestCheckResourceAttr(resourceName, "port_security_enabled", "false"),
					testAccCheckNetworkingV2PortPortSecurity(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noPortSecurityNoSecurityGroups(t *testing.T) {
	var port testPortWithExtensions
	resourceName := "opentelekomcloud_networking_port_v2.port_1"
	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists(resourceName, &port),
					resource.TestCheckResourceAttr(resourceName, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "0"),
					testAccCheckNetworkingV2PortPortSecurity(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet
	t.Parallel()
	qts := subnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SubnetExists("opentelekomcloud_networking_subnet_v2.subnet_1", &subnet),
					TestAccCheckNetworkingV2NetworkExists("opentelekomcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_port_v2" {
			continue
		}

		_, err := ports.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortExists(n string, port *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
		}

		found, err := ports.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingV2PortWithExtensionsExists(n string, port *testPortWithExtensions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		var found testPortWithExtensions
		err = ports.Get(client, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("port not found")
		}

		*port = found

		return nil
	}
}

func testAccCheckNetworkingV2PortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmt.Errorf("expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortPortSecurity(port *testPortWithExtensions, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if port.PortSecurityEnabled != expected {
			return fmt.Errorf("port has wrong port_security_enabled. Expected %t, got %t", expected, port.PortSecurityEnabled)
		}

		return nil
	}
}

const testAccNetworkingV2PortBasic = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_secgroup_v2" "group_1" {
  name        = "terraform-sg-1"
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }
  security_group_ids = [
    opentelekomcloud_networking_secgroup_v2.group_1.id
  ]
}
`

const testAccNetworkingV2PortNoIP = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairs = `
resource "opentelekomcloud_networking_network_v2" "vrrp_network" {
  name           = "vrrp_network"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "vrrp_subnet" {
  name       = "vrrp_subnet"
  cidr       = "10.0.0.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.vrrp_network.id

  allocation_pools {
    start = "10.0.0.2"
    end   = "10.0.0.200"
  }
}

resource "opentelekomcloud_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "opentelekomcloud_networking_router_interface_v2" "vrrp_interface" {
  router_id = opentelekomcloud_networking_router_v2.vrrp_router.id
  subnet_id = opentelekomcloud_networking_subnet_v2.vrrp_subnet.id
}

resource "opentelekomcloud_networking_port_v2" "vrrp_port_1" {
  name           = "vrrp_port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.vrrp_network.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.vrrp_subnet.id
    ip_address = "10.0.0.202"
  }
}

resource "opentelekomcloud_networking_port_v2" "vrrp_port_2" {
  name           = "vrrp_port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.vrrp_network.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.vrrp_subnet.id
    ip_address = "10.0.0.201"
  }
}

resource "opentelekomcloud_networking_port_v2" "instance_port" {
  name           = "instance_port"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.vrrp_network.id

  allowed_address_pairs {
    ip_address  = opentelekomcloud_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address
    mac_address = opentelekomcloud_networking_port_v2.vrrp_port_1.mac_address
  }

  allowed_address_pairs {
    ip_address  = opentelekomcloud_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address
    mac_address = opentelekomcloud_networking_port_v2.vrrp_port_2.mac_address
  }
}
`

const testAccNetworkingV2PortSecurityDisabled = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
}
resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name                  = "port_1"
  network_id            = opentelekomcloud_networking_network_v2.network_1.id
  port_security_enabled = false
  no_security_groups    = true
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortSecurityEnabled = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
}
resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name                  = "port_1"
  network_id            = opentelekomcloud_networking_network_v2.network_1.id
  port_security_enabled = true
  no_security_groups    = false
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortTimeout = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
