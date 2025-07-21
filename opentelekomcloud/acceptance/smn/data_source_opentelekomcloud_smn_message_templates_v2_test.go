package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceSmnMessageTemplates_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_smn_message_templates_v2.test"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnMessageTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "templates.0.protocol", "default"),
					resource.TestCheckResourceAttr(dataSource, "templates.0.name", "test-message-template"),
				),
			},
		},
	})
}

var testDataSourceSmnMessageTemplates_basic = `
data "opentelekomcloud_smn_message_templates_v2" "test" {
  name     = opentelekomcloud_smn_message_template_v2.test.name
  protocol = "default"
}

resource "opentelekomcloud_smn_message_template_v2" "test" {
  name     = "test-message-template"
  protocol = "default"
  content  = "Test content, contains {content1} and {content2}"
}
`
