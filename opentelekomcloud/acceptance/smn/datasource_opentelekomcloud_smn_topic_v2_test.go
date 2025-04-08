package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const datasourceTopicName = "data.opentelekomcloud_smn_topic_v2.topic_1"

func TestAccSMNV2TopicDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestAccSMNV2DataSourceTopicConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNTopicV2DataSourceID(datasourceTopicName),
					resource.TestCheckResourceAttr(datasourceTopicName, "name", "test_topic"),
					resource.TestCheckResourceAttr(datasourceTopicName, "display_name", "Test topic"),
				),
			},
		},
	})
}

func testAccCheckSMNTopicV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find topic data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SMN topic data source ID not set ")
		}

		return nil
	}
}

var TestAccSMNV2DataSourceTopicConfig_basic = `
data "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "test_topic"
}
`
