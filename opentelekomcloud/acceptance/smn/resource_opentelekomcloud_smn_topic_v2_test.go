package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/topics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceTopicName = "opentelekomcloud_smn_topic_v2.topic_1"

func TestAccSMNV2Topic_basic(t *testing.T) {
	var topic topics.TopicGet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2TopicConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2TopicExists(resourceTopicName, &topic, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(resourceTopicName, "name", "topic_1"),
					resource.TestCheckResourceAttr(resourceTopicName, "display_name", "The display name of topic_1"),
					resource.TestCheckResourceAttr(resourceTopicName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceTopicName, "tags.key", "value"),
				),
			},
			{
				Config: TestAccSMNV2TopicConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceTopicName, "display_name", "The update display name of topic_1"),
					resource.TestCheckResourceAttr(resourceTopicName, "name", "topic_1"),
					resource.TestCheckResourceAttr(resourceTopicName, "tags.foo", "bar_ch"),
				),
			},
		},
	})
}

func TestAccSMNV2Topic_schemaProjectName(t *testing.T) {
	var topic topics.TopicGet
	var projectName2 = os.Getenv("OS_PROJECT_NAME_2")
	if projectName2 == "" {
		t.Skip("OS_PROJECT_NAME_2 should be set in order to run test")
	}
	env.OS_TENANT_NAME = cfg.ProjectName(projectName2)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_projectName(env.OS_TENANT_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2TopicExists(resourceTopicName, &topic, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(resourceTopicName, "project_name", string(env.OS_TENANT_NAME)),
				),
			},
		},
	})
	env.OS_TENANT_NAME = env.GetTenantName()
}

func testAccCheckSMNTopicV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	smnClient, err := config.SmnV2Client(env.OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud smn: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_smn_topic_v2" {
			continue
		}

		_, err := topics.Get(smnClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("topic still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2TopicExists(n string, topic *topics.TopicGet, projectName cfg.ProjectName) resource.TestCheckFunc {
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

		found, err := topics.Get(smnClient, rs.Primary.ID).ExtractGet()
		if err != nil {
			return err
		}

		if found.TopicUrn != rs.Primary.ID {
			return fmt.Errorf("topic not found")
		}

		*topic = *found

		return nil
	}
}

var TestAccSMNV2TopicConfig_basic = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`

var TestAccSMNV2TopicConfig_update = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The update display name of topic_1"

  tags = {
    foo = "bar_ch"
  }
}
`

func testAccSMNV2TopicConfig_projectName(projectName cfg.ProjectName) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
  project_name = "%s"
}
`, projectName)
}
