package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceRmsResourceRelationships_basic(t *testing.T) {
	t.Skip("You are not authorized with rms:resources:list impossible to run within CI.")
	dataSource := "data.opentelekomcloud_rms_resource_relationships_v1.relations_1"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsResourceRelationships_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.relation_type"),
				),
			},
		},
	})
}

var testDataSourceDataSourceRmsResourceRelationships_basic = `
data "opentelekomcloud_rms_resource_relationships_v1" "relations_1" {
  resource_id = "c3ce2ac4-3c03-44c0-9433-11c7bd390662"
  direction   = "in"
}
`
