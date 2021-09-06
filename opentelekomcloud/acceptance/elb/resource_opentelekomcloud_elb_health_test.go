package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/elbaas/healthcheck"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceHealthName = "opentelekomcloud_elb_health.health_1"

func TestAccELBHealth_basic(t *testing.T) {
	var health healthcheck.Health

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckELBHealthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBHealthConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBHealthExists(t, resourceHealthName, &health),
					resource.TestCheckResourceAttr(resourceHealthName, "healthy_threshold", "3"),
					resource.TestCheckResourceAttr(resourceHealthName, "healthcheck_timeout", "10"),
				),
			},
			{
				Config: testAccELBHealthConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHealthName, "healthy_threshold", "5"),
					resource.TestCheckResourceAttr(resourceHealthName, "healthcheck_timeout", "15"),
				),
			},
		},
	})
}

func testAccCheckELBHealthDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.ElbV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_elb_healthcheck" {
			continue
		}

		_, err := healthcheck.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("health still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBHealthExists(t *testing.T, n string, health *healthcheck.Health) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] testAccCheckELBHealthExists resources %+v.\n", s.RootModule().Resources)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := healthcheck.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] testAccCheckELBHealthExists found %+v.\n", found)

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("health not found")
		}

		*health = *found

		return nil
	}
}

var testAccELBHealthConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name      = "loadbalancer_1"
  vpc_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  type      = "External"
  bandwidth = 5
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name             = "listener_1"
  protocol         = "TCP"
  protocol_port    = 8080
  backend_protocol = "TCP"
  backend_port     = 8080
  lb_algorithm     = "roundrobin"
  loadbalancer_id  = opentelekomcloud_elb_loadbalancer.loadbalancer_1.id
}


resource "opentelekomcloud_elb_health" "health_1" {
  listener_id       = opentelekomcloud_elb_listener.listener_1.id
  #healthcheck_protocol = "HTTP"
  healthy_threshold = 3
  #healthcheck_timeout = 10
  #healthcheck_interval = 5

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)

var testAccELBHealthConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name      = "loadbalancer_1"
  vpc_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  type      = "External"
  bandwidth = 5
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name             = "listener_1"
  protocol         = "TCP"
  protocol_port    = 8080
  backend_protocol = "TCP"
  backend_port     = 8080
  lb_algorithm     = "roundrobin"
  loadbalancer_id  = opentelekomcloud_elb_loadbalancer.loadbalancer_1.id
}


resource "opentelekomcloud_elb_health" "health_1" {
  listener_id          = opentelekomcloud_elb_listener.listener_1.id
  healthcheck_protocol = "HTTP"
  healthy_threshold    = 5
  healthcheck_timeout  = 15
  healthcheck_interval = 3

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)
