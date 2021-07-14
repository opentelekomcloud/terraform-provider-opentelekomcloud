package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/topicattributes"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

// NOTE: Sometimes SMN resources are not created XAAS-8823
func TestAccSMNV2TopicAttribute_basic(t *testing.T) {
	resourceName := "opentelekomcloud_smn_topic_attribute_v2.attribute_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNTopicAttributeV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2TopicAttributeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2TopicAttributeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "attribute_name", "access_policy"),
				),
			},
		},
	})
}

func TestAccSMNV2TopicAttribute_import(t *testing.T) {
	resourceName := "opentelekomcloud_smn_topic_attribute_v2.attribute_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNTopicAttributeV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2TopicAttributeConfigBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSMNTopicAttributeV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SmnV2Client(env.OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SMNv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_smn_topic_attribute_v2" {
			continue
		}

		listOpts := topicattributes.ListOpts{
			Name: rs.Primary.Attributes["attribute_name"],
		}

		topicAttrs, err := topicattributes.List(client, rs.Primary.Attributes["topic_urn"], listOpts).Extract()
		if err == nil {
			if topicAttrs["access_policy"] == "" {
				return nil
			}
			return fmt.Errorf("SMN Topic attributes still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2TopicAttributeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.SmnV2Client(env.OS_TENANT_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud SMNv2 client: %w", err)
		}
		listOpts := topicattributes.ListOpts{
			Name: rs.Primary.Attributes["attribute_name"],
		}

		_, err = topicattributes.List(client, rs.Primary.Attributes["topic_urn"], listOpts).Extract()
		if err != nil {
			return err
		}

		return nil
	}
}

var TestAccSMNV2TopicAttributeConfigBasic = fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		   = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_smn_topic_attribute_v2" "attribute_1" {
  topic_urn       = opentelekomcloud_smn_topic_v2.topic_1.topic_urn
  attribute_name  = "access_policy"
  topic_attribute = <<EOF
{
  "Version": "2016-09-07",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__service_pub_0",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "OBS"
        ]
      },
      "Action": [
        "SMN:Publish",
        "SMN:QueryTopicDetail"
      ],
      "Resource": "${opentelekomcloud_smn_topic_v2.topic_1.topic_urn}"
    }
  ]
}
EOF
}
`)
