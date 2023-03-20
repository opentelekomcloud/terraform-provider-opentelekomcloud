package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceRuleName = "opentelekomcloud_lb_rule_v3.this"

func TestAccLBV3Rule_basic(t *testing.T) {
	var rule rules.ForwardingRule

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LbPolicy, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3RuleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3RuleExists(resourceRuleName, &rule),
					resource.TestCheckResourceAttr(resourceRuleName, "value", "^.+$"),
				),
			},
			{
				Config: testAccLBV3RuleConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRuleName, "value", "^.*$"),
				),
			},
		},
	})
}

func TestAccLBV3Rule_condition(t *testing.T) {
	var rule rules.ForwardingRule

	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LbPolicy, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3RuleConfigCondition,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3RuleExists(resourceRuleName, &rule),
					resource.TestCheckResourceAttr(resourceRuleName, "conditions.0.value", "/"),
				),
			},
			{
				Config: testAccLBV3RuleConfigConditionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRuleName, "conditions.0.value", "/home"),
				),
			},
		},
	})
}

func TestAccLBRuleV3_import(t *testing.T) {
	t.Parallel()
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LbPolicy, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.LbListener, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3RuleConfigBasic,
			},
			{
				ResourceName:      resourceRuleName,
				ImportStateVerify: true,
				ImportState:       true,
			},
		},
	})
}

func testAccCheckLBV3RuleDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_rule_v3" {
			continue
		}

		part := strings.Split(rs.Primary.ID, "/")
		_, err := rules.Get(client, part[0], part[1]).Extract()
		if err == nil {
			return fmt.Errorf("rule still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV3RuleExists(n string, rule *rules.ForwardingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		part := strings.Split(rs.Primary.ID, "/")
		if rs.Primary.ID == "" || len(part) != 2 {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elbv3.ErrCreateClient, err)
		}

		found, err := rules.Get(client, part[0], part[1]).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.Attributes["rule_id"] {
			return fmt.Errorf("rule not found")
		}

		rule = found

		return nil
	}
}

var testAccLBV3RuleConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol        = "HTTP"
  protocol_port   = 8080
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  action           = "REDIRECT_TO_POOL"
  listener_id      = opentelekomcloud_lb_listener_v3.this.id
  redirect_pool_id = opentelekomcloud_lb_pool_v3.this.id
  position         = 37
}

resource "opentelekomcloud_lb_rule_v3" "this" {
  type         = "PATH"
  compare_type = "REGEX"
  value        = "^.+$"
  policy_id    = opentelekomcloud_lb_policy_v3.this.id
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3RuleConfigUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol        = "HTTP"
  protocol_port   = 8080
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  name             = "policy_updated"
  description      = "some interesting description"
  action           = "REDIRECT_TO_POOL"
  listener_id      = opentelekomcloud_lb_listener_v3.this.id
  redirect_pool_id = opentelekomcloud_lb_pool_v3.this.id
  position         = 37
}

resource "opentelekomcloud_lb_rule_v3" "this" {
  type         = "PATH"
  compare_type = "REGEX"
  value        = "^.*$"
  policy_id    = opentelekomcloud_lb_policy_v3.this.id
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3RuleConfigCondition = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id     = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol            = "HTTP"
  protocol_port       = 8080
  advanced_forwarding = true
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  action           = "REDIRECT_TO_POOL"
  listener_id      = opentelekomcloud_lb_listener_v3.this.id
  redirect_pool_id = opentelekomcloud_lb_pool_v3.this.id
  position         = 37
}

resource "opentelekomcloud_lb_rule_v3" "this" {
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/test"
  policy_id    = opentelekomcloud_lb_policy_v3.this.id
  conditions {
    value = "/"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3RuleConfigConditionUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id     = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol            = "HTTP"
  protocol_port       = 8080
  advanced_forwarding = true
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  name             = "policy_updated"
  description      = "some interesting description"
  action           = "REDIRECT_TO_POOL"
  listener_id      = opentelekomcloud_lb_listener_v3.this.id
  redirect_pool_id = opentelekomcloud_lb_pool_v3.this.id
  position         = 37
}

resource "opentelekomcloud_lb_rule_v3" "this" {
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/test"
  policy_id    = opentelekomcloud_lb_policy_v3.this.id

  conditions {
    value = "/home"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
