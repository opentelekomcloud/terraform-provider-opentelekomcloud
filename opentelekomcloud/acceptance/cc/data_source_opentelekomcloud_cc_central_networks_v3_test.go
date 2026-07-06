package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCcCentralNetworksV3DataSource_basic(t *testing.T) {
	name := fmt.Sprintf("cc_acc_cn_ds_%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cc_central_networks_v3.by_name"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcCentralNetworksV3DataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "central_networks.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "central_networks.0.id",
						"opentelekomcloud_cc_central_network_v3.test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "central_networks.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "central_networks.0.state", "AVAILABLE"),
					resource.TestCheckResourceAttrSet(dataSourceName, "central_networks.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "central_networks.0.created_at"),
				),
			},
		},
	})
}

func testAccCcCentralNetworksV3DataSource_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cc_central_network_v3" "test" {
  name        = "%s"
  description = "created by terraform acceptance test"
}

data "opentelekomcloud_cc_central_networks_v3" "by_name" {
  name = opentelekomcloud_cc_central_network_v3.test.name
}
`, name)
}
