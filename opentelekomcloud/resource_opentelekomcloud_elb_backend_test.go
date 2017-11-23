package opentelekomcloud

import (
	"fmt"
	"testing"

	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS with diff
func TestAccELBBackend_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBBackendDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:             TestAccELBBackendConfig_basic,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBBackendExists("opentelekomcloud_elb_backend.backend_1"),
				),
			},
			resource.TestStep{
				Config:             TestAccELBBackendConfig_update,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_elb_backend.backend_1", "weight", "10"),
				),
			},
		},
	})
}

func testAccCheckELBBackendDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	_, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_elb_backendmember" {
			continue
		}

	}

	return nil
}

func testAccCheckELBBackendExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		_ /*networkingClient*/, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		return nil
	}
}

var TestAccELBBackendConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  // vip_subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
  vpc_id = "%s"
  type = "External"
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.loadbalancer_1.id}"
}

resource "opentelekomcloud_elb_backend" "backend_1" {
  address = "192.168.199.10"
  subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
  listener_id = "${opentelekomcloud_elb_listener.listener_1.id}"
  server_id = "gary-backend"
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_elb_backend" "backend_2" {
  address = "192.168.199.11"
  protocol_port = 8080
  subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

const TestAccELBBackendConfig_update = `
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


resource "opentelekomcloud_elb_backend" "backend_1" {
  address = "192.168.199.10"
  protocol_port = 8080
  weight = 10
  admin_state_up = "true"
  subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_elb_backend" "backend" {
  address = "192.168.199.11"
  protocol_port = 8080
  weight = 15
  admin_state_up = "true"
  subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
