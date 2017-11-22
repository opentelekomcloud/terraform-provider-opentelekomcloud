package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccELBHealth_basic(t *testing.T) {
	var health healthcheck.Health

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccELBHealthConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBHealthExists(t, "opentelekomcloud_elb_healthcheck.health_1", &health),
				),
			},
			resource.TestStep{
				Config: TestAccLBV2MonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_healthcheck.health_1", "name", "health1_updated"),
					resource.TestCheckResourceAttr("opentelekomcloud_elb_healthcheck.health_1", "delay", "30"),
					resource.TestCheckResourceAttr("opentelekomcloud_elb_healthcheck.health_1", "timeout", "15"),
				),
			},
		},
	})
}

func testAccCheckELBHealthDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_elb_healthcheck" {
			continue
		}

		_, err := healthcheck.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Health still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBHealthExists(t *testing.T, n string, health *healthcheck.Health) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := healthcheck.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Health not found")
		}

		*health = *found

		return nil
	}
}

var TestAccELBHealthConfig_basic = fmt.Sprintf(`
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
  //vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
  vpc_id = "%s"
  type = "External"
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.loadbalancer_1.id}"
}


resource "opentelekomcloud_elb_health" "health_1" {
  name = "health_1"
  type = "PING"
  delay = 20
  timeout = 10
  max_retries = 5
  // pool_id = "${opentelekomcloud_lb_pool_v2.pool_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

const TestAccELBHealthConfig_update = `
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
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.loadbalancer_1.id}"
}


resource "opentelekomcloud_elb_health" "health_1" {
  name = "health_1_updated"
  type = "PING"
  delay = 30
  timeout = 15
  max_retries = 10
  admin_state_up = "true"
  //pool_id = "${opentelekomcloud_lb_pool_v2.pool_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
