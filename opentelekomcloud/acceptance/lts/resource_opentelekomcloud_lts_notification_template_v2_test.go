package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	message_template "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/message-template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getNotificationTemplateResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v2 client: %s", err)
	}

	requestResp, err := message_template.List(client, client.DomainID)
	if err != nil {
		return nil, err
	}
	if len(requestResp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var ruleResult message_template.MessageTemplateResponse
	for _, t := range requestResp {
		if t.Name == state.Primary.ID {
			ruleResult = t
		}
	}
	if ruleResult.Name == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return ruleResult, nil
}

func TestAccNotificationTemplate_basic(t *testing.T) {
	var (
		template message_template.TemplateResponse
		name     = fmt.Sprintf("lts_template%s", acctest.RandString(3))
		rName    = "opentelekomcloud_lts_notification_template_v2.test"
	)

	rc := common.InitResourceCheck(
		rName,
		&template,
		getNotificationTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testNotificationTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source", "LTS"),
					resource.TestCheckResourceAttr(rName, "language", "en-us"),
					resource.TestCheckResourceAttr(rName, "description", "acc test"),
					resource.TestCheckResourceAttr(rName, "templates.#", "1"),
					resource.TestCheckResourceAttr(rName, "templates.0.sub_type", "sms"),
				),
			},
			{
				Config: testNotificationTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source", "LTS"),
					resource.TestCheckResourceAttr(rName, "language", "en-us"),
					resource.TestCheckResourceAttr(rName, "description", "acc test update"),
					resource.TestCheckResourceAttr(rName, "templates.#", "2"),
					resource.TestCheckResourceAttr(rName, "templates.0.sub_type", "sms"),
					resource.TestCheckResourceAttr(rName, "templates.1.sub_type", "email"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testNotificationTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_notification_template_v2" "test" {
  name        = "%s"
  source      = "LTS"
  language    = "en-us"
  description = "acc test"

  templates {
    sub_type = "sms"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }
}
`, name)
}

func testNotificationTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_notification_template_v2" "test" {
  name        = "%s"
  description = "acc test update"
  source      = "LTS"
  language    = "en-us"

  templates {
    sub_type = "sms"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }

  templates {
    sub_type = "email"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }
}
`, name)
}
