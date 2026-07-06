package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCcCentralNetworkCapabilitiesV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_cc_central_network_capabilities_v3.test"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcCentralNetworkCapabilitiesV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "capabilities.#"),
				),
			},
		},
	})
}

const testAccCcCentralNetworkCapabilitiesV3DataSource_basic = `
data "opentelekomcloud_cc_central_network_capabilities_v3" "test" {}
`
