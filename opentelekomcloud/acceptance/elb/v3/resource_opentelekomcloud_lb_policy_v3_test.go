package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourcePolicyName = "opentelekomcloud_lb_policy_v3.this"

func TestAccLBV3Policy_basic(t *testing.T) {
	var policy policies.Policy

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
		CheckDestroy:      testAccCheckLBV3PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3PolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3PolicyExists(resourcePolicyName, &policy),
				),
			},
			{
				Config: testAccLBV3PolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "policy_updated"),
				),
			},
		},
	})
}

func TestAccLBV3Policy_fixedResponse(t *testing.T) {
	var policy policies.Policy

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
		CheckDestroy:      testAccCheckLBV3PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3PolicyConfigFixedResponse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3PolicyExists(resourcePolicyName, &policy),
					resource.TestCheckResourceAttr(resourcePolicyName, "priority", "10"),
					resource.TestCheckResourceAttr(resourcePolicyName, "fixed_response_config.0.status_code", "200"),
					resource.TestCheckResourceAttr(resourcePolicyName, "fixed_response_config.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr(resourcePolicyName, "fixed_response_config.0.message_body", "Fixed Response"),
				),
			},
			{
				Config: testAccLBV3PolicyConfigFixedResponseUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "policy_updated"),
					resource.TestCheckResourceAttr(resourcePolicyName, "priority", "11"),
					resource.TestCheckResourceAttr(resourcePolicyName, "fixed_response_config.0.status_code", "202"),
					resource.TestCheckResourceAttr(resourcePolicyName, "fixed_response_config.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr(resourcePolicyName, "fixed_response_config.0.message_body", "Fixed Response update"),
				),
			},
		},
	})
}

func TestAccLBV3Policy_redirectUrl(t *testing.T) {
	var policy policies.Policy

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
		CheckDestroy:      testAccCheckLBV3PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3PolicyConfigRedirectUrl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3PolicyExists(resourcePolicyName, &policy),
					resource.TestCheckResourceAttr(resourcePolicyName, "redirect_url", "https://www.google.com:443"),
					resource.TestCheckResourceAttr(resourcePolicyName, "redirect_url_config.0.status_code", "301"),
					resource.TestCheckResourceAttr(resourcePolicyName, "redirect_url_config.0.query", "name=my_name"),
				),
			},
			{
				Config: testAccLBV3PolicyConfigRedirectUrlUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePolicyName, "name", "policy_updated"),
					resource.TestCheckResourceAttr(resourcePolicyName, "redirect_url", "https://www.google.com:443"),
					resource.TestCheckResourceAttr(resourcePolicyName, "redirect_url_config.0.status_code", "308"),
					resource.TestCheckResourceAttr(resourcePolicyName, "redirect_url_config.0.query", "name=my_name_updated"),
				),
			},
		},
	})
}

func TestAccLBPolicyV3_import(t *testing.T) {
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
		CheckDestroy:      testAccCheckLBV3PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3PolicyConfigBasic,
			},
			{
				ResourceName:      resourcePolicyName,
				ImportStateVerify: true,
				ImportState:       true,
			},
		},
	})
}

func testAccCheckLBV3PolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_policy_v3" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV3PolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elbv3.ErrCreateClient, err)
		}

		found, err := policies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("policy not found")
		}

		policy = found

		return nil
	}
}

var testAccLBV3PolicyConfigBasic = fmt.Sprintf(`
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
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3PolicyConfigUpdate = fmt.Sprintf(`
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

  rules {
    type         = "HOST_NAME"
    compare_type = "EQUAL_TO"
    value        = "abc.com"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3PolicyConfigFixedResponse = fmt.Sprintf(`
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
  action      = "FIXED_RESPONSE"
  listener_id = opentelekomcloud_lb_listener_v3.this.id
  position    = 37
  priority    = 10

  fixed_response_config {
    status_code  = "200"
    content_type = "text/plain"
    message_body = "Fixed Response"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3PolicyConfigFixedResponseUpdate = fmt.Sprintf(`
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
  name        = "policy_updated"
  description = "some interesting description"
  action      = "FIXED_RESPONSE"
  listener_id = opentelekomcloud_lb_listener_v3.this.id
  position    = 37
  priority    = 11

  fixed_response_config {
    status_code  = "202"
    content_type = "text/plain"
    message_body = "Fixed Response update"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3PolicyConfigRedirectUrl = fmt.Sprintf(`
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
  action      = "REDIRECT_TO_URL"
  listener_id = opentelekomcloud_lb_listener_v3.this.id
  position    = 37
  priority    = 10

  redirect_url = "https://www.google.com:443"

  redirect_url_config {
    status_code = "301"
    query       = "name=my_name"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccLBV3PolicyConfigRedirectUrlUpdate = fmt.Sprintf(`
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
  name        = "policy_updated"
  description = "some interesting description"
  action      = "REDIRECT_TO_URL"
  listener_id = opentelekomcloud_lb_listener_v3.this.id
  position    = 37
  priority    = 11

  redirect_url_config {
    status_code = "308"
    query       = "name=my_name_updated"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
