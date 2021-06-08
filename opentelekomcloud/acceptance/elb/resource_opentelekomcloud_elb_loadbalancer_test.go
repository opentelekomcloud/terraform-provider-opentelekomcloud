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
				Config: testAccELBLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists("opentelekomcloud_elb_loadbalancer.loadbalancer_1", &lb),
				),
			},
			{
				Config: testAccELBLoadBalancerConfig_update,
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
			return fmt.Errorf("Member not found")
		}

		*lb = *found

		return nil
	}
}

var testAccELBLoadBalancerConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = "5"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, env.OS_VPC_ID)

var testAccELBLoadBalancerConfig_update = fmt.Sprintf(`
resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1_updated"
  admin_state_up = "true"
  vpc_id = "%s"
  type = "External"
  bandwidth = 3

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, env.OS_VPC_ID)
