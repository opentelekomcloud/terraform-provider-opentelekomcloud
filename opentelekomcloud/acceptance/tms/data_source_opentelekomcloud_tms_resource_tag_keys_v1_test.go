package tms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceTmsTagKeys_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_tms_resource_tag_keys_v1.test"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTmsTagKeys_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "keys.#"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTmsTagKeys_basic() string {
	return `
data "opentelekomcloud_tms_resource_tag_keys_v1" "test" {}

data "opentelekomcloud_tms_resource_tag_keys_v1" "region_id_filter" {
  region_id = "eu-de"
}
output "region_id_filter_is_useful" {
  value = length(data.opentelekomcloud_tms_resource_tag_keys_v1.region_id_filter.keys) > 0
}
`
}
