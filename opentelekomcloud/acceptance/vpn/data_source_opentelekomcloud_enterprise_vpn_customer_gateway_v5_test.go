package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataVPNEnterpriseCustomerGatewayName = "data.opentelekomcloud_enterprise_vpn_customer_gateway_v5.gw_1"

func TestAccVpnEnterpriseCustomerGatewayV5DataSource_basic(t *testing.T) {
	gatewayId := os.Getenv("OS_VPN_CUSTOMER_GATEWAY_ID")
	if gatewayId == "" {
		t.Skip("Customer Gateway ID is required for the test")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpnEnterpriseCustomerGatewayV5ConfigBasic(gatewayId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEnterpriseCustomerGatewayDataSourceID(dataVPNEnterpriseCustomerGatewayName),
					resource.TestCheckResourceAttrSet(dataVPNEnterpriseCustomerGatewayName, "name"),
					resource.TestCheckResourceAttrSet(dataVPNEnterpriseCustomerGatewayName, "id_value"),
				),
			},
		},
	})
}

func testAccCheckVPNEnterpriseCustomerGatewayDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Enterprise Customer gateway data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("enterprise customer gateway data source ID not set ")
		}

		return nil
	}
}

func testAccDataSourceVpnEnterpriseCustomerGatewayV5ConfigBasic(gatewayId string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_enterprise_vpn_customer_gateway_v5" "gw_1" {
  id = "%s"
}
`, gatewayId)
}
