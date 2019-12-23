package opentelekomcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func TestAccLBV2L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2L7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7PolicyConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("opentelekomcloud_lb_l7policy_v2.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_lb_l7policy_v2.l7policy_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_lb_l7policy_v2.l7policy_1", "description", "test description"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_lb_l7policy_v2.l7policy_1", "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_lb_l7policy_v2.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"opentelekomcloud_lb_l7policy_v2.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
		},
	})
}

func testAccCheckLBV2L7PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpentelekomCloud load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_l7policy_v2" {
			continue
		}

		_, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2L7PolicyExists(n string, l7Policy *l7policies.L7Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpentelekomCloud load balancing client: %s", err)
		}

		found, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7Policy = *found

		return nil
	}
}

var testAccCheckLBV2L7PolicyConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "opentelekomcloud_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_POOL"
  description  = "test description"
  position     = 1
  listener_id  = "${opentelekomcloud_lb_listener_v2.listener_1.id}"
  redirect_pool_id = "${opentelekomcloud_lb_pool_v2.pool_1.id}"
}
`, OS_SUBNET_ID)
