package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePolicyName = "opentelekomcloud_lb_l7policy_v2.l7policy_1"

func TestAccLBV2L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2L7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7PolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists(resourcePolicyName, &l7Policy),
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "test"),
					resource.TestCheckResourceAttr(resourcePolicyName, "description", "test description"),
					resource.TestCheckResourceAttr(resourcePolicyName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourcePolicyName, "position", "1"),
					resource.TestMatchResourceAttr(resourcePolicyName, "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
		},
	})
}

func testAccCheckLBV2L7PolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	lbClient, err := config.ElbV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elb.ErrCreationV2Client, err)
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
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		lbClient, err := config.ElbV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elb.ErrCreationV2Client, err)
		}

		found, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("policy not found")
		}

		*l7Policy = *found

		return nil
	}
}

var testAccCheckLBV2L7PolicyConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_l7policy_v2" "l7policy_1" {
  name             = "test"
  action           = "REDIRECT_TO_POOL"
  description      = "test description"
  position         = 1
  listener_id      = opentelekomcloud_lb_listener_v2.listener_1.id
  redirect_pool_id = opentelekomcloud_lb_pool_v2.pool_1.id
}
`, common.DataSourceSubnet)
