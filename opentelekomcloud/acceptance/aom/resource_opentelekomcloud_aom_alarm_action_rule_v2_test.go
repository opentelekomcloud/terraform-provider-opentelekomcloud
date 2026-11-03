package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/aom/v2/alarm"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/aom"
)

const resourceAlarmActionRuleName = "opentelekomcloud_aom_alarm_action_rule_v2.rule"

func getAlarmActionRuleFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.AomV2Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM v2 client: %s", err)
	}
	return aom.GetAlarmActionRule(client, state.Primary.ID)
}

func TestAccAomAlarmActionRuleV2_basic(t *testing.T) {
	rand5 := acctest.RandString(5)
	ruleName := fmt.Sprintf("terraform_acc_%s", rand5)
	topicName := fmt.Sprintf("topic_acc_%s", rand5)

	rc := common.InitResourceCheck(
		resourceAlarmActionRuleName,
		&alarm.ActionRule{},
		getAlarmActionRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAomAlarmActionRuleV2Basic(ruleName, topicName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "type", "1"),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "smn_topic.#", "1"),
				),
			},
			{
				Config: testAccAomAlarmActionRuleV2Updated(ruleName, topicName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "type", "1"),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "description", "updated by acceptance test"),
					resource.TestCheckResourceAttr(resourceAlarmActionRuleName, "smn_topic.#", "1"),
				),
			},
			{
				ResourceName:      resourceAlarmActionRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAomAlarmActionRuleV2Basic(ruleName, topicName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic" {
  name = "%[1]s"
}

resource "opentelekomcloud_aom_alarm_action_rule_v2" "rule" {
  name                  = "%[2]s"
  type                  = "1"
  description           = "created by acceptance test"
  notification_template = "aom.built-in.template.en"

  smn_topic {
    topic_urn   = opentelekomcloud_smn_topic_v2.topic.topic_urn
    name        = opentelekomcloud_smn_topic_v2.topic.name
    push_policy = 0
  }
}
`, topicName, ruleName)
}

func testAccAomAlarmActionRuleV2Updated(ruleName, topicName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic" {
  name = "%[1]s"
}

resource "opentelekomcloud_aom_alarm_action_rule_v2" "rule" {
  name                  = "%[2]s"
  type                  = "1"
  description           = "updated by acceptance test"
  notification_template = "aom.built-in.template.en"

  smn_topic {
    topic_urn   = opentelekomcloud_smn_topic_v2.topic.topic_urn
    name        = opentelekomcloud_smn_topic_v2.topic.name
    push_policy = 0
  }
}
`, topicName, ruleName)
}
