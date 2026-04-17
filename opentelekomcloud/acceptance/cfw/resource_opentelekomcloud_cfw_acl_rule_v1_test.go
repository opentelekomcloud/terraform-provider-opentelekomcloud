package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/acl"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const aclRuleResourceName = "opentelekomcloud_cfw_acl_rule_v1.rule_1"

func getACLRuleFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	return acl.GetACLRule(client, state.Primary.Attributes["object_id"], state.Primary.Attributes["name"])
}

func TestAccCFWACLRuleV1_basic(t *testing.T) {
	var aclRule acl.ACLRule
	rc := common.InitResourceCheck(
		aclRuleResourceName,
		&aclRule,
		getACLRuleFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWACLRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclRuleResourceName, "name", "test-acc-tf-acl-rule"),
					resource.TestCheckResourceAttrSet(aclRuleResourceName, "id"),
				),
			},
			{
				Config: testAccCFWACLRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclRuleResourceName, "name", "test-acc-tf-acl-rule-updated"),
				),
			},
			{
				ResourceName:      aclRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCFWACLRuleV1ImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"sequence",
					"applications",
					"applications_json_string",
					"source.0.predefined_group",
					"destination.0.predefined_group",
					"service.0.predefined_group",
				},
			},
		},
	})
}

func testAccCFWACLRuleV1ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var objectId string
		var name string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cfw_acl_rule_v1" {
				objectId = rs.Primary.Attributes["object_id"]
				name = rs.Primary.Attributes["name"]
			}
		}
		if objectId == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s", objectId, name)
		}
		return fmt.Sprintf("%s/%s", objectId, name), nil
	}
}

var testAccCFWACLRuleV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_acl_rule_v1" "rule_1" {
  object_id = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  type      = 0
  name      = "test-acc-tf-acl-rule"
  sequence {
    top = 1
  }
  address_type        = 0
  action_type         = 0
  status              = 1
  long_connect_enable = 0
  direction           = 0
  source {
    type    = 0
    address = "1.1.1.1"
  }
  destination {
    type    = 0
    address = "2.2.2.2"
  }
  service {
    type     = 0
    protocol = -1
  }
}
`

var testAccCFWACLRuleV1Update = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_acl_rule_v1" "rule_1" {
  object_id = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  type      = 0
  name      = "test-acc-tf-acl-rule-updated"
  sequence {
    top = 1
  }
  address_type        = 0
  action_type         = 0
  status              = 1
  long_connect_enable = 0
  direction           = 0
  source {
    type    = 0
    address = "1.1.1.1"
  }
  destination {
    type    = 0
    address = "2.2.2.2"
  }
  service {
    type     = 0
    protocol = -1
  }
}
`
