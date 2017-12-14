package opentelekomcloud

import (
	"fmt"
	"testing"

	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

// PASS with diff
func TestAccELBBackend_basic(t *testing.T) {
	var backend backendmember.Backend

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBBackendDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:             TestAccELBBackendConfig_basic,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false, unfinished elb?
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBBackendExists("opentelekomcloud_elb_backend.backend_1", &backend),
				),
			},
		},
	})
}

func testAccCheckELBBackendDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_elb_healthcheck" {
			continue
		}

		_, err := backendmember.Get(client, rs.Primary.Attributes["listener_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Backend member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBBackendExists(n string, backend *backendmember.Backend) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := backendmember.Get(client, rs.Primary.Attributes["listener_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] testAccCheckELBBackendExists found %+v.\n", found)

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Backend member not found")
		}

		*backend = *found

		return nil
	}
}

var TestAccELBBackendConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = 5
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.loadbalancer_1.id}"
}


resource "opentelekomcloud_elb_health" "health_1" {
  listener_id = "${opentelekomcloud_elb_listener.listener_1.id}"
  healthcheck_protocol = "HTTP"
  healthy_threshold = 3
  healthcheck_timeout = 10
  healthcheck_interval = 5

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_elb_backend" "backend_1" {
  address = "%s"
  listener_id = "${opentelekomcloud_elb_listener.listener_1.id}"
  server_id = "%s"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID, OS_SERVER_ADDRESS, OS_SERVER_ID)
