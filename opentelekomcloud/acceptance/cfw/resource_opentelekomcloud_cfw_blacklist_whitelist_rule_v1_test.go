package acceptance

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	list "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/blackwhitelist"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const blacklistWhitelistRuleResourceName = "opentelekomcloud_cfw_blacklist_whitelist_rule_v1.rule_1"

func getBlacklistWhitelistRuleFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	listType, err := strconv.Atoi(state.Primary.Attributes["list_type"])
	if err != nil {
		return nil, fmt.Errorf("error converting list type to integer: %s", err)
	}
	objectId := state.Primary.Attributes["object_id"]
	address := state.Primary.Attributes["address"]
	return list.GetBlacklistOrWhitelistRule(client, objectId, listType, address)
}

func TestAccCFWBlacklistWhitelistRuleV1_basic(t *testing.T) {
	var blacklistWhitelistRule list.BlackWhiteListRecord
	rc := common.InitResourceCheck(
		blacklistWhitelistRuleResourceName,
		&blacklistWhitelistRule,
		getBlacklistWhitelistRuleFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWBlacklistWhitelistRuleV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(blacklistWhitelistRuleResourceName, "description", "Test111161"),
					resource.TestCheckResourceAttr(blacklistWhitelistRuleResourceName, "port", "1"),
					resource.TestCheckResourceAttrSet(blacklistWhitelistRuleResourceName, "id"),
				),
			},
			{
				Config: testAccCFWBlacklistWhitelistRuleV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(blacklistWhitelistRuleResourceName, "description", "Test111162"),
					resource.TestCheckResourceAttr(blacklistWhitelistRuleResourceName, "port", "2"),
				),
			},
			{
				ResourceName:      blacklistWhitelistRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCFWBlacklistWhitelistRuleV1ImportStateIdFunc(),
			},
		},
	})
}

func testAccCFWBlacklistWhitelistRuleV1ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var objectId, listType, address string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cfw_blacklist_whitelist_rule_v1" {
				objectId = rs.Primary.Attributes["object_id"]
				listType = rs.Primary.Attributes["list_type"]
				address = rs.Primary.Attributes["address"]
			}
		}
		if objectId == "" || listType == "" || address == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", objectId, listType, address)
		}
		return fmt.Sprintf("%s/%s/%s", objectId, listType, address), nil
	}
}

var testAccCFWBlacklistWhitelistRuleV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_blacklist_whitelist_rule_v1" "rule_1" {
  object_id    = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  list_type    = 5
  direction    = 0
  address_type = 0
  address      = "1.1.1.1"
  protocol     = 6
  port         = "1"
  description  = "Test111161"
}
`

var testAccCFWBlacklistWhitelistRuleV1Update = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_blacklist_whitelist_rule_v1" "rule_1" {
  object_id    = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  list_type    = 5
  direction    = 0
  address_type = 0
  address      = "1.1.1.1"
  protocol     = 6
  port         = "2"
  description  = "Test111162"
}
`
