package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataConnectionName = "data.opentelekomcloud_enterprise_vpn_connection_v5.conn_1"

func TestAccConnectionDataSource_basic(t *testing.T) {
	connId := os.Getenv("OS_VPN_CONNECTION_ID")
	if connId == "" {
		t.Skip("OS_VPN_CONNECTION_ID is required for the test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEvpnConnectionDataSource_basic(connId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEnterpriseConnectionDataSourceID(dataConnectionName),
					resource.TestCheckResourceAttrSet(dataConnectionName, "name"),
					resource.TestCheckResourceAttrSet(dataConnectionName, "customer_gateway_id"),
				),
			},
		},
	})
}

func testAccCheckVPNEnterpriseConnectionDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find enterprise vpn connection data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("enterprise vpn connection data source ID not set ")
		}

		return nil
	}
}

func testEvpnConnectionDataSource_basic(connId string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_enterprise_vpn_connection_v5" "conn_1" {
  id = "%s"
}
`, connId)
}
