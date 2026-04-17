package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const datasourceTopicName = "data.opentelekomcloud_smn_topic_v2.topic_1"

func TestAccSMNV2TopicDataSource_basic(t *testing.T) {
	dc := common.InitDataSourceCheck(datasourceTopicName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2DataSourceTopicConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceTopicName, "name", "test_topic"),
					resource.TestCheckResourceAttr(datasourceTopicName, "display_name", "Test topic"),
				),
			},
		},
	})
}

var TestAccSMNV2DataSourceTopicConfig_basic = `
data "opentelekomcloud_smn_topic_v2" "topic_1" {
  name = opentelekomcloud_smn_topic_v2.topic.name
}

resource "opentelekomcloud_smn_topic_v2" "topic" {
  name         = "test_topic"
  display_name = "Test topic"
}
`
