package opentelekomcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elbaas/backendmember"
)

func TestAccELBBackend_basic(t *testing.T) {
	var backend backendmember.Backend

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccELBBackendConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBBackendExists("opentelekomcloud_elb_backend.backend_1", &backend),
				),
			},
		},
	})
}

func testAccCheckELBBackendDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.elbV1Client(OS_REGION_NAME)
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
		client, err := config.elbV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := backendmember.Get(client, rs.Primary.Attributes["listener_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		f := found[0]
		log.Printf("[DEBUG] testAccCheckELBBackendExists found %+v.\n", found)

		if f.ID != rs.Primary.ID {
			return fmt.Errorf("Backend member not found")
		}

		*backend = f

		return nil
	}
}

var TestAccELBBackendConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = 5
  admin_state_up = true
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
  address = "${opentelekomcloud_compute_instance_v2.vm_1.network.0.fixed_ip_v4}"
  listener_id = "${opentelekomcloud_elb_listener.listener_1.id}"
  server_id = "${opentelekomcloud_compute_instance_v2.vm_1.id}"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
