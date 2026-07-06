package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCcCentralNetworkConnectionsV3DataSource_basic(t *testing.T) {
	name := fmt.Sprintf("cc_acc_conn_ds_%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cc_central_network_connections_v3.test"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcCentralNetworkConnectionsV3DataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.#"),
				),
			},
		},
	})
}

func testAccCcCentralNetworkConnectionsV3DataSource_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cc_central_network_v3" "test" {
  name        = "%s"
  description = "created by terraform acceptance test"
}

data "opentelekomcloud_cc_central_network_connections_v3" "test" {
  central_network_id = opentelekomcloud_cc_central_network_v3.test.id
  is_cross_region    = "false"
}
`, name)
}
