package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/monitors"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	resourceName := "opentelekomcloud_lb_monitor_v2.monitor_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2MonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "monitor_port", "112"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", ""),
				),
			},
			{
				Config: TestAccLBV2MonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "monitor_1_updated"),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "monitor_port", "120"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.test.com"),
				),
			},
		},
	})
}

func TestAccLBV2Monitor_minConfig(t *testing.T) {
	var monitor monitors.Monitor
	resourceName := "opentelekomcloud_lb_monitor_v2.monitor_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2MonitorConfig_minConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
				),
			},
			{
				Config: TestAccLBV2MonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "monitor_1_updated"),
					resource.TestCheckResourceAttr(resourceName, "monitor_port", "120"),
				),
			},
		},
	})
}

func testAccCheckLBV2MonitorDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_monitor_v2" {
			continue
		}

		_, err := monitors.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("monitor not found")
		}

		*monitor = *found

		return nil
	}
}

var TestAccLBV2MonitorConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_1" {
  name         = "monitor_1"
  type         = "TCP"
  delay        = 20
  timeout      = 10
  max_retries  = 5
  pool_id      = opentelekomcloud_lb_pool_v2.pool_1.id
  monitor_port = 112

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, env.OS_SUBNET_ID)

var TestAccLBV2MonitorConfig_update = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_1" {
  name           = "monitor_1_updated"
  type           = "TCP"
  delay          = 30
  timeout        = 15
  max_retries    = 10
  admin_state_up = "true"
  pool_id        = opentelekomcloud_lb_pool_v2.pool_1.id
  monitor_port   = 120
  domain_name    = "www.test.com"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, env.OS_SUBNET_ID)

var TestAccLBV2MonitorConfig_minConfig = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_1" {
  type         = "TCP"
  delay        = 20
  timeout      = 10
  max_retries  = 5
  pool_id      = opentelekomcloud_lb_pool_v2.pool_1.id
}
`, env.OS_SUBNET_ID)
