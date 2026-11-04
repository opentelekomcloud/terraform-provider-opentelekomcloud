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

const resourceAlarmSilenceRuleName = "opentelekomcloud_aom_alarm_silence_rule_v2.rule"

func getAlarmSilenceRuleFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.AomV2Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM v2 client: %s", err)
	}
	return aom.GetAlarmSilenceRule(client, state.Primary.ID)
}

func TestAccAomAlarmSilenceRuleV2_basic(t *testing.T) {
	ruleName := fmt.Sprintf("terraform_acc_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceAlarmSilenceRuleName,
		&alarm.SilenceRule{},
		getAlarmSilenceRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAomAlarmSilenceRuleV2Basic(ruleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceAlarmSilenceRuleName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceAlarmSilenceRuleName, "mute_config.0.type", "DAILY"),
					resource.TestCheckResourceAttr(resourceAlarmSilenceRuleName, "match.0.conditions.0.operate", "EQUALS"),
				),
			},
			{
				Config: testAccAomAlarmSilenceRuleV2Updated(ruleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceAlarmSilenceRuleName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceAlarmSilenceRuleName, "description", "updated by acceptance test"),
					resource.TestCheckResourceAttr(resourceAlarmSilenceRuleName, "mute_config.0.type", "DAILY"),
				),
			},
			{
				ResourceName:      resourceAlarmSilenceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAomAlarmSilenceRuleV2Basic(ruleName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_aom_alarm_silence_rule_v2" "rule" {
  name        = "%s"
  description = "created by acceptance test"
  timezone    = "UTC"

  match {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      value   = ["Info"]
    }
  }

  mute_config {
    type      = "DAILY"
    starts_at = 0
    ends_at   = 86399
  }
}
`, ruleName)
}

func testAccAomAlarmSilenceRuleV2Updated(ruleName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_aom_alarm_silence_rule_v2" "rule" {
  name        = "%s"
  description = "updated by acceptance test"
  timezone    = "UTC"

  match {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      value   = ["Info"]
    }
  }

  mute_config {
    type      = "DAILY"
    starts_at = 0
    ends_at   = 86399
  }
}
`, ruleName)
}
