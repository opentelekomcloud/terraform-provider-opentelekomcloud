package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"regexp"
)

// PASS
func TestAccELBLoadBalancer_basic(t *testing.T) {
	var lb loadbalancer_elbs.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccELBLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists("opentelekomcloud_elb_loadbalancer.loadbalancer_1", &lb),
				),
			},
			resource.TestStep{
				Config: testAccELBLoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", "name", "loadbalancer_1_updated"),
					resource.TestMatchResourceAttr(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", "vip_port_id",
						regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
		},
	})
}

// PASS
func TestAccELBLoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancer_elbs.LoadBalancer
	var sg_1, sg_2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccELBLoadBalancer_secGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_1", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_1", &sg_2),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", "security_group_ids.#", "1"),
					testAccCheckELBLoadBalancerHasSecGroup(&lb, &sg_1),
				),
			},
			resource.TestStep{
				Config: testAccLBV2LoadBalancer_secGroup_update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", &lb),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", "security_group_ids.#", "2"),
					testAccCheckELBLoadBalancerHasSecGroup(&lb, &sg_1),
					testAccCheckELBLoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
			resource.TestStep{
				Config: testAccELBLoadBalancer_secGroup_update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", &lb),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", "security_group_ids.#", "1"),
					testAccCheckELBLoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
		},
	})
}

func testAccCheckELBLoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_elb_loadbalancer" {
			continue
		}

		_, err := loadbalancer_elbs.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBLoadBalancerExists(
	n string, lb *loadbalancer_elbs.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := loadbalancer_elbs.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*lb = *found

		return nil
	}
}

func testAccCheckELBLoadBalancerHasSecGroup(
	lb *loadbalancer_elbs.LoadBalancer, sg *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		_ /*networkingClient,*/, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		return nil
	}
}

const testAccELBLoadBalancerConfig_basic = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${opentelekomcloud_networking_network_v2.network_1.id}"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccELBLoadBalancerConfig_update = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${opentelekomcloud_networking_network_v2.network_1.id}"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1_updated"
  #loadbalancer_provider = "haproxy"
  admin_state_up = "true"
  vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccELBLoadBalancer_secGroup = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${opentelekomcloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
    security_group_ids = [
      "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
    ]
}
`

const testAccELBLoadBalancer_secGroup_update1 = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${opentelekomcloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
    security_group_ids = [
      "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}",
      "${opentelekomcloud_networking_secgroup_v2.secgroup_2.id}"
    ]
}
`

const testAccELBLoadBalancer_secGroup_update2 = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${opentelekomcloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
    security_group_ids = [
      "${opentelekomcloud_networking_secgroup_v2.secgroup_2.id}"
    ]
    depends_on = ["opentelekomcloud_networking_secgroup_v2.secgroup_1"]
}
`
