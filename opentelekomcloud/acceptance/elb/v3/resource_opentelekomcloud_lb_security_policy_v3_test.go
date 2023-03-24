package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/security_policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceSecurityPolicyName = "opentelekomcloud_lb_security_policy_v3.this"

func TestAccLBV3SecurityPolicy_basic(t *testing.T) {
	var policy security_policy.SecurityPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3SecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3SecurityPolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3SecurityPolicyExists(resourceSecurityPolicyName, &policy),
				),
			},
			{
				Config: testAccLBV3SecurityPolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSecurityPolicyName, "name", "updated-security-policy"),
					resource.TestCheckResourceAttr(resourceSecurityPolicyName, "description", "test-description-updated"),
				),
			},
		},
	})
}

func TestAccLBV3SecurityPolicy_assignment(t *testing.T) {
	var policy security_policy.SecurityPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3SecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3SecurityPolicyAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3SecurityPolicyExists(resourceSecurityPolicyName, &policy),
					resource.TestCheckResourceAttrSet(resourceListenerName, "security_policy_id"),
				),
			},
		},
	})
}

func TestAccLBSecurityPolicyV3_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3SecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3SecurityPolicyConfigBasic,
			},
			{
				ResourceName:      resourceSecurityPolicyName,
				ImportStateVerify: true,
				ImportState:       true,
			},
		},
	})
}

func testAccCheckLBV3SecurityPolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_security_policy_v3" {
			continue
		}

		_, err := security_policy.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV3SecurityPolicyExists(n string, policy *security_policy.SecurityPolicy) resource.TestCheckFunc {
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

		found, err := security_policy.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.SecurityPolicy.ID != rs.Primary.ID {
			return fmt.Errorf("policy not found")
		}

		policy = found

		return nil
	}
}

var testAccLBV3SecurityPolicyConfigBasic = `
resource "opentelekomcloud_lb_security_policy_v3" "this" {
  name        = "new-security-policy"
  description = "test-description"
  protocols   = ["TLSv1", "TLSv1.1"]
  ciphers     = ["ECDHE-ECDSA-AES128-SHA", "ECDHE-RSA-AES128-SHA"]
}
`

var testAccLBV3SecurityPolicyConfigUpdate = `
resource "opentelekomcloud_lb_security_policy_v3" "this" {
  name        = "updated-security-policy"
  description = "test-description-updated"
  protocols   = ["TLSv1", "TLSv1.1"]
  ciphers     = ["ECDHE-ECDSA-AES128-SHA"]
}`

var testAccLBV3SecurityPolicyAssignment = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%s"]
}

resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1"
  type        = "server"
  private_key = %s
  certificate = %s
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name                      = "listener_1"
  description               = "some interesting description"
  loadbalancer_id           = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol                  = "HTTPS"
  protocol_port             = 443
  default_tls_container_ref = opentelekomcloud_lb_certificate_v3.certificate_1.id
  security_policy_id        = opentelekomcloud_lb_security_policy_v3.this.id

  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }
}
resource "opentelekomcloud_lb_security_policy_v3" "this" {
  name        = "assignmend-security-policy"
  description = "test-description"
  protocols   = ["TLSv1", "TLSv1.1"]
  ciphers     = ["ECDHE-ECDSA-AES128-SHA", "ECDHE-RSA-AES128-SHA"]
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, privateKey, certificate)
