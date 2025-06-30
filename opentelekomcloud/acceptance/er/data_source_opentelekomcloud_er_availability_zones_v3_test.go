package er

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceAZs_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_er_availability_zones_v3.zones"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceErAZs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
		},
	})
}

const testAccDataSourceErAZs_basic string = `data "opentelekomcloud_er_availability_zones_v3" "zones" {}`
