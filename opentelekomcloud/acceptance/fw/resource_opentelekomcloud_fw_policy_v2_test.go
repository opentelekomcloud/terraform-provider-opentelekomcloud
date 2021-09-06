package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/fwaas_v2/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccFWPolicyV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"opentelekomcloud_fw_policy_v2.policy_1", "", "", 0),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_addRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_addRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"opentelekomcloud_fw_policy_v2.policy_1", "policy_1", "terraform acceptance test", 2),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_deleteRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_deleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"opentelekomcloud_fw_policy_v2.policy_1", "policy_1", "terraform acceptance test", 1),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_timeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"opentelekomcloud_fw_policy_v2.policy_1", "", "", 0),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_removeSingleRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2SingleRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("opentelekomcloud_fw_policy_v2.policy_1", "policy_1", "", 1),
				),
			},
			{
				Config: testAccFWPolicyV2SingleRule_removeRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("opentelekomcloud_fw_policy_v2.policy_1", "policy_1", "", 0),
				),
			},
		},
	})
}

func testAccCheckFWPolicyV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_fw_policy_v2" {
			continue
		}
		_, err = policies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("firewall policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWPolicyV2Exists(n, name, description string, ruleCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		var found *policies.Policy
		for i := 0; i < 5; i++ {
			// Firewall policy creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = policies.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case name != found.Name:
			err = fmt.Errorf("expected name <%s>, but found <%s>", name, found.Name)
		case description != found.Description:
			err = fmt.Errorf("expected description <%s>, but found <%s>", description, found.Description)
		case ruleCount != len(found.Rules):
			err = fmt.Errorf("expected rule count <%d>, but found <%d>", ruleCount, len(found.Rules))
		}

		if err != nil {
			return err
		}

		return nil
	}
}

const testAccFWPolicyV2_basic = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
}
`

const testAccFWPolicyV2_addRules = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
  description =  "terraform acceptance test"
  rules = [
    opentelekomcloud_fw_rule_v2.udp_deny.id,
    opentelekomcloud_fw_rule_v2.tcp_allow.id
  ]
}

resource "opentelekomcloud_fw_rule_v2" "tcp_allow" {
  protocol = "tcp"
  action = "allow"
}

resource "opentelekomcloud_fw_rule_v2" "udp_deny" {
  protocol = "udp"
  action = "deny"
}
`

const testAccFWPolicyV2_deleteRules = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
  description =  "terraform acceptance test"
  rules = [
    opentelekomcloud_fw_rule_v2.udp_deny.id
  ]
}

resource "opentelekomcloud_fw_rule_v2" "udp_deny" {
  protocol = "udp"
  action = "deny"
}
`

const testAccFWPolicyV2_timeout = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  timeouts {
    create = "5m"
  }
}
`

const testAccFWPolicyV2SingleRule = `
resource "opentelekomcloud_fw_rule_v2" "rule_1" {
  name = "rule_1"
  action = "allow"
  description = "allow-all"
  protocol = "any"
  source_ip_address = "0.0.0.0/0"
  destination_ip_address = "0.0.0.0/0"
  source_port = ""
  destination_port = ""
  ip_version = "4"
}

resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
  rules = [opentelekomcloud_fw_rule_v2.rule_1.id]
}
`

const testAccFWPolicyV2SingleRule_removeRule = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
  rules = []
}
`
