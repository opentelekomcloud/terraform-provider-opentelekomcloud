package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	resourceName := "opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
				),
			},
			{
				Config: testAccLBV2LoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "loadbalancer_1_updated"),
					resource.TestMatchResourceAttr(resourceName, "vip_port_id", regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
		},
	})
}

func testAccCheckLBV2LoadBalancerDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_loadbalancer_v2" {
			continue
		}

		_, err := loadbalancers.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2LoadBalancerExists(n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
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

		found, err := loadbalancers.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*lb = *found

		return nil
	}
}

var testAccLBV2LoadBalancerConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, env.OS_SUBNET_ID)

var testAccLBV2LoadBalancerConfig_update = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name           = "loadbalancer_1_updated"
  admin_state_up = "true"
  vip_subnet_id  = "%s"

  tags = {
    muh = "value-update"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, env.OS_SUBNET_ID)
