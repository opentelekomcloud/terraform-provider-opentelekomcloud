package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataVPNEnterpriseGatewayName = "data.opentelekomcloud_enterprise_vpn_gateway_v5.gw_1"

func TestAccVpnEnterpriseGatewayV5DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpnEnterpriseGatewayV5ConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEnterpriseGatewayDataSourceID(dataVPNEnterpriseGatewayName),
					resource.TestCheckResourceAttrSet(dataVPNEnterpriseGatewayName, "name"),
					resource.TestCheckResourceAttrSet(dataVPNEnterpriseGatewayName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckVPNEnterpriseGatewayDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Enterprise VPN gateway data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Enterprise VPN gateway data source ID not set ")
		}

		return nil
	}
}

const testAccDataSourceVpnEnterpriseGatewayV5ConfigBasic = `
data "opentelekomcloud_enterprise_vpn_gateway_v5" "gw_1" {
  id = "a7c33913-0eb3-40e2-a473-53c699592235"
}
`
