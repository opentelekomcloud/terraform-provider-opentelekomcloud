package opentelekomcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/smn/v2/topics"
)

func TestAccSMNV2Topic_basic(t *testing.T) {
	var topic topics.TopicGet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMNTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2TopicConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2TopicExists("opentelekomcloud_smn_topic_v2.topic_1", &topic, OS_TENANT_NAME),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_topic_v2.topic_1", "name", "topic_1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_topic_v2.topic_1", "display_name",
						"The display name of topic_1"),
				),
			},
			{
				Config: TestAccSMNV2TopicConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_topic_v2.topic_1", "display_name",
						"The update display name of topic_1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_topic_v2.topic_1", "name", "topic_1"),
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
	OS_TENANT_NAME = projectName2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMNTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_projectName(OS_TENANT_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2TopicExists("opentelekomcloud_smn_topic_v2.topic_1", &topic, OS_TENANT_NAME),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_smn_topic_v2.topic_1", "project_name", OS_TENANT_NAME),
				),
			},
		},
	})
	OS_TENANT_NAME = getTenantName()
}

func testAccCheckSMNTopicV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	smnClient, err := config.SmnV2Client(OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud smn: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_smn_topic_v2" {
			continue
		}

		_, err := topics.Get(smnClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Topic still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2TopicExists(n string, topic *topics.TopicGet, projectName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		smnClient, err := config.SmnV2Client(projectName)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud smn client: %s", err)
		}

		found, err := topics.Get(smnClient, rs.Primary.ID).ExtractGet()
		if err != nil {
			return err
		}

		if found.TopicUrn != rs.Primary.ID {
			return fmt.Errorf("Topic not found")
		}

		*topic = *found

		return nil
	}
}

var TestAccSMNV2TopicConfig_basic = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  = "topic_1"
  display_name    = "The display name of topic_1"
}
`

var TestAccSMNV2TopicConfig_update = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  = "topic_1"
  display_name    = "The update display name of topic_1"
}
`

func testAccSMNV2TopicConfig_projectName(projectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		   = "topic_1"
  display_name = "The display name of topic_1"
  project_name = "%s"
}
`, projectName)
}
