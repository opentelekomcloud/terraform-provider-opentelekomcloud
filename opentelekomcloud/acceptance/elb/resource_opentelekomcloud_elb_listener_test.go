package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/elbaas/listeners"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceListenerV1Name = "opentelekomcloud_elb_listener.listener_1"

func TestAccELBListener_basic(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckELBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBListenerExists(resourceListenerV1Name, &listener),
				),
			},
			{
				Config: testAccELBListenerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceListenerV1Name, "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(resourceListenerV1Name, "backend_port", "8088"),
				),
			},
		},
	})
}

func testAccCheckELBListenerDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.ElbV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_listener_v2" {
			continue
		}

		_, err := listeners.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		found, err := listeners.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*listener = *found

		return nil
	}
}

var testAccELBListenerConfigBasic = fmt.Sprintf(`
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

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)

var testAccELBListenerConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name      = "loadbalancer_1"
  vpc_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  type      = "External"
  bandwidth = 5
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name             = "listener_1_updated"
  protocol         = "TCP"
  protocol_port    = 8080
  backend_protocol = "TCP"
  backend_port     = 8088
  lb_algorithm     = "roundrobin"
  loadbalancer_id  = opentelekomcloud_elb_loadbalancer.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)
