package dcaas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestDirectConnectV2Datasource_basic(t *testing.T) {
	var directConnectName = fmt.Sprintf("dc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectConnectV2Datasource_basic(directConnectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.opentelekomcloud_direct_connect_v2.direct_connect", "bandwidth", "100"),
				),
			},
		},
	})
}

func testAccDirectConnectV2Datasource_basic(directConnectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_direct_connect_v2" "direct_connect" {
  name          = "%s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 100
  provider_name = "OTC"
}

data "opentelekomcloud_direct_connect_v2" "direct_connect" {
  id = opentelekomcloud_direct_connect_v2.direct_connect.id
}
`, directConnectName)
}
