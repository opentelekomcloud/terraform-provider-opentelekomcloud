package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func TestAccLBV2Listener_basic(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2ListenerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("opentelekomcloud_lb_listener_v2.listener_1", &listener),
					/* resource.TestCheckResourceAttr(
					"opentelekomcloud_lb_listener_v2.listener_1", "connection_limit", "-1"), */
				),
			},
			{
				Config: TestAccLBV2ListenerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					/* resource.TestCheckResourceAttr(
					"opentelekomcloud_lb_listener_v2.listener_1", "connection_limit", "100"), */
				),
			},
		},
	})
}

func testAccCheckLBV2ListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_listener_v2" {
			continue
		}

		_, err := listeners.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2ListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := listeners.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*listener = *found

		return nil
	}
}

var TestAccLBV2ListenerConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id}"

	timeouts {
		create = "5m"
		update = "5m"
		delete = "5m"
	}
}
`, OS_SUBNET_ID)

var TestAccLBV2ListenerConfig_update = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  #connection_limit = 100
  admin_state_up = "true"
  loadbalancer_id = "${opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id}"

	timeouts {
		create = "5m"
		update = "5m"
		delete = "5m"
	}
}
`, OS_SUBNET_ID)
