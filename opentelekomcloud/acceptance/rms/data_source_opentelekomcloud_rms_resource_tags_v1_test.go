package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceRmsResourceTags_basic(t *testing.T) {
	t.Skip("You are not authorized with rms:resources:list impossible to run within CI")
	dataSource := "data.opentelekomcloud_rms_resource_tags_v1.tags_1"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsResourceTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
				),
			},
		},
	})
}

var testDataSourceDataSourceRmsResourceTags_basic = `
data "opentelekomcloud_rms_resource_tags_v1" "tags_1" {}
`
