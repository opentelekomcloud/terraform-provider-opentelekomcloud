package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/alarm"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getKeywordsAlarmRuleResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v2 client: %s", err)
	}

	requestResp, err := alarm.ListKeywordRules(client)
	if err != nil {
		return nil, err
	}
	if len(requestResp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var ruleResult alarm.KeywordRule
	for _, acc := range requestResp {
		if acc.ID == state.Primary.ID {
			ruleResult = acc
		}
	}
	if ruleResult.ID == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return ruleResult, nil
}

func TestAccKeywordsAlarmRule_basic(t *testing.T) {
	var (
		rule  alarm.KeywordRule
		name  = fmt.Sprintf("lts_keyword%s", acctest.RandString(3))
		rName = "opentelekomcloud_lts_keywords_alarm_rule_v2.test"
	)

	rc := common.InitResourceCheck(
		rName,
		&rule,
		getKeywordsAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeywordsAlarmRuleBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "severity", "CRITICAL"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
				),
			},
			{
				Config: testKeywordsAlarmRuleBasic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "severity", "INFO"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notification_rule",
				},
			},
		},
	})
}

func testKeywordsAlarmRuleBasic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic" {
  name         = "topic_%[1]s"
  display_name = "The display name of topic"
}

resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_keywords_alarm_rule_v2" "test" {
  name        = "%[1]s"
  description = "created by terraform"
  severity    = "CRITICAL"

  notification_frequency = 5

  keywords_requests {
    keyword                = "%[1]s_key_words"
    condition              = ">"
    number                 = 100
    log_group_id           = opentelekomcloud_lts_group_v2.group.id
    log_stream_id          = opentelekomcloud_lts_stream_v2.stream.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }

  notification_rule {
    language  = "en-us"
    timezone  = "xx/xx"
    user_name = "test"

    topics {
      name      = opentelekomcloud_smn_topic_v2.topic.name
      topic_urn = opentelekomcloud_smn_topic_v2.topic.topic_urn
    }
  }
}
`, name)
}

func testKeywordsAlarmRuleBasic_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic" {
  name         = "topic_%[1]s"
  display_name = "The display name of topic"
}

resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_keywords_alarm_rule_v2" "test" {
  name        = "%[1]s"
  description = ""
  severity    = "INFO"

  notification_frequency = 5

  keywords_requests {
    keyword                = "%[1]s_key_words"
    condition              = ">"
    number                 = 100
    log_group_id           = opentelekomcloud_lts_group_v2.group.id
    log_stream_id          = opentelekomcloud_lts_stream_v2.stream.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
  }

  notification_rule {
    language  = "en-us"
    timezone  = "xx/xx"
    user_name = "test"

    topics {
      name      = opentelekomcloud_smn_topic_v2.topic.name
      topic_urn = opentelekomcloud_smn_topic_v2.topic.topic_urn
    }
  }
}
`, name)
}
