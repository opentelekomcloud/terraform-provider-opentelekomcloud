package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceSmnTopicSubscriptions_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_smn_topic_subscription_v2.test"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnTopicSubscription_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoint"),
					resource.TestCheckResourceAttrSet(dataSource, "topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "owner"),
					resource.TestCheckResourceAttrSet(dataSource, "remark"),
				),
			},
		},
	})
}

var testDataSourceSmnTopicSubscription_basic = `
data "opentelekomcloud_smn_topic_subscription_v2" "test" {
  topic_urn  = opentelekomcloud_smn_topic_v2.topic_1.id
  protocol   = "email"
  depends_on = [opentelekomcloud_smn_subscription_v2.subscription_1]
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_smn_subscription_v2" "subscription_1" {
  topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}
`
