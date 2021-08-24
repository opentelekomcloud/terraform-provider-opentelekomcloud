package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/subscriptions"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSMNV2Subscription_basic(t *testing.T) {
	var subscription1 subscriptions.SubscriptionGet
	var subscription2 subscriptions.SubscriptionGet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNSubscriptionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2SubscriptionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2SubscriptionExists(
						"opentelekomcloud_smn_subscription_v2.subscription_1", &subscription1, env.OsTenantName),
					testAccCheckSMNV2SubscriptionExists(
						"opentelekomcloud_smn_subscription_v2.subscription_2", &subscription2, env.OsTenantName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_subscription_v2.subscription_1", "endpoint", "mailtest@gmail.com"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_subscription_v2.subscription_2", "endpoint", "13600000000"),
				),
			},
		},
	})
}

func TestAccSMNV2Subscription_schemaProjectName(t *testing.T) {
	var subscription1 subscriptions.SubscriptionGet

	var projectName2 = os.Getenv("OS_PROJECT_NAME_2")
	if projectName2 == "" {
		t.Skip("OS_PROJECT_NAME_2 should be set in order to run test")
	}
	env.OsTenantName = cfg.ProjectName(projectName2)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNSubscriptionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2SubscriptionConfig_projectName(env.OsTenantName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2SubscriptionExists(
						"opentelekomcloud_smn_subscription_v2.subscription_1", &subscription1, env.OsTenantName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_subscription_v2.subscription_1", "project_name", string(env.OsTenantName)),
				),
			},
		},
	})
	env.OsTenantName = env.GetTenantName()
}

func testAccCheckSMNSubscriptionV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	smnClient, err := config.SmnV2Client(env.OsTenantName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud smn: %s", err)
	}
	var subscription *subscriptions.SubscriptionGet
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_smn_subscription_v2" {
			continue
		}
		foundList, err := subscriptions.List(smnClient).Extract()
		if err != nil {
			return err
		}
		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription != nil {
			return fmt.Errorf("subscription still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2SubscriptionExists(n string, subscription *subscriptions.SubscriptionGet, projectName cfg.ProjectName) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		smnClient, err := config.SmnV2Client(projectName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud smn client: %s", err)
		}

		foundList, err := subscriptions.List(smnClient).Extract()
		if err != nil {
			return err
		}
		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription == nil {
			return fmt.Errorf("subscription not found")
		}

		return nil
	}
}

var TestAccSMNV2SubscriptionConfig_basic = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		   = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_smn_subscription_v2" "subscription_1" {
  topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}

resource "opentelekomcloud_smn_subscription_v2" "subscription_2" {
  topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "O&M"
}
`

func testAccSMNV2SubscriptionConfig_projectName(projectName cfg.ProjectName) string {
	return fmt.Sprintf(`
locals {
  project_name = "%s"
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
  project_name = local.project_name
}

resource "opentelekomcloud_smn_subscription_v2" "subscription_1" {
  topic_urn    = opentelekomcloud_smn_topic_v2.topic_1.id
  endpoint     = "mailtest@gmail.com"
  protocol     = "email"
  remark       = "O&M"
  project_name = local.project_name
}
`, projectName)
}
