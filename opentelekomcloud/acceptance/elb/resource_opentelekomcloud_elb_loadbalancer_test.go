package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccELBLoadBalancer_basic(t *testing.T) {
	var lb loadbalancer_elbs.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBLoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists("opentelekomcloud_elb_loadbalancer.loadbalancer_1", &lb),
				),
			},
			{
				Config: testAccELBLoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_loadbalancer.loadbalancer_1", "name", "loadbalancer_1_updated"),
				),
			},
		},
	})
}

func testAccCheckELBLoadBalancerDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.ElbV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_elb_loadbalancer" {
			continue
		}

		_, err := loadbalancer_elbs.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("loadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBLoadBalancerExists(
	n string, lb *loadbalancer_elbs.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.ElbV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}
		found, err := loadbalancer_elbs.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			log.Printf("[#####ERR#####] : %v", err)

			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}
		*lb = *found

		return nil
	}
}

var testAccELBLoadBalancerConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name      = "loadbalancer_1"
  vpc_id    = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  type      = "External"
  bandwidth = "5"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)

var testAccELBLoadBalancerConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name           = "loadbalancer_1_updated"
  admin_state_up = "true"
  vpc_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  type           = "External"
  bandwidth      = 3

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, common.DataSourceSubnet)
