package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataNatGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2GatewayBasic,
			},
			{
				Config: testAccDataNatV2GatewayBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.opentelekomcloud_nat_gateway_v2.nat_1", "name", "nat_1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_nat_gateway_v2.nat_1", "description", "test for terraform"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_nat_gateway_v2.nat_1", "spec", "1"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_nat_gateway_v2.nat_1", "tenant_id"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_nat_gateway_v2.nat_1", "internal_network_id"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_nat_gateway_v2.nat_1", "router_id"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_nat_gateway_v2.nat_1", "status"),
				),
			},
		},
	})
}

var testAccDataNatV2GatewayBasic = fmt.Sprintf(`
%s

data "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name           = "nat_1"
  admin_state_up = true
}
`, testAccNatV2GatewayBasic)
