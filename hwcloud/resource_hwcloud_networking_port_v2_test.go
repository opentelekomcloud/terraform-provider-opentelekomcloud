package hwcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

// PASS
func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_basic,
				//ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

// PASS
func TestAccNetworkingV2Port_noip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_noip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

// KNOWN problem (#7)
/*
func TestAccNetworkingV2Port_multipleNoIP(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_multipleNoIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 3),
				),
			},
		},
	})
}
*/

// PASS
func TestAccNetworkingV2Port_allowedAddressPairs(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrp_port_1, vrrp_port_2, instance_port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:             testAccNetworkingV2Port_allowedAddressPairs,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.instance_port", &instance_port),
				),
			},
		},
	})
}

// KNOWN problem (#7)
/*
func TestAccNetworkingV2Port_multipleFixedIPs(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_multipleFixedIPs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 3),
				),
			},
		},
	})
}
*/

// PASS
func TestAccNetworkingV2Port_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

// KNOWN problem (#7)
/*
func TestAccNetworkingV2Port_fixedIPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_fixedIPs,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hwcloud_networking_port_v2.port_1", "all_fixed_ips.0", "192.168.199.23"),
					resource.TestCheckResourceAttr(
						"hwcloud_networking_port_v2.port_1", "all_fixed_ips.1", "192.168.199.24"),
				),
			},
		},
	})
}
*/

// KNOWN Problem (#15)
/*
func TestAccNetworkingV2Port_updateSecurityGroups(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var security_group groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"hwcloud_networking_secgroup_v2.secgroup_1", &security_group),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"hwcloud_networking_secgroup_v2.secgroup_1", &security_group),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"hwcloud_networking_secgroup_v2.secgroup_1", &security_group),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("hwcloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("hwcloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("hwcloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"hwcloud_networking_secgroup_v2.secgroup_1", &security_group),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}
*/

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HWCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hwcloud_networking_port_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortExists(n string, port *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HWCloud networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingV2PortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmt.Errorf("Expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("Expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

const testAccNetworkingV2Port_basic = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_noip = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
  }
}
`

const testAccNetworkingV2Port_multipleNoIP = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
  }
}
`

const testAccNetworkingV2Port_allowedAddressPairs = `
resource "hwcloud_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "hwcloud_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "hwcloud_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${hwcloud_networking_router_v2.vrrp_router.id}"
  subnet_id = "${hwcloud_networking_subnet_v2.vrrp_subnet.id}"
}

resource "hwcloud_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "hwcloud_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "hwcloud_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${hwcloud_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${hwcloud_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${hwcloud_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${hwcloud_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2Port_multipleFixedIPs = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.40"
  }
}
`

const testAccNetworkingV2Port_timeout = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Port_fixedIPs = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.24"
  }

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_1 = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_2 = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
  security_group_ids = ["${hwcloud_networking_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_3 = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "terraform security group acceptance test"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
  security_group_ids = ["${hwcloud_networking_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_4 = `
resource "hwcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "hwcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
}

resource "hwcloud_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "hwcloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${hwcloud_networking_network_v2.network_1.id}"
	security_group_ids = []

  fixed_ip {
    subnet_id =  "${hwcloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`
