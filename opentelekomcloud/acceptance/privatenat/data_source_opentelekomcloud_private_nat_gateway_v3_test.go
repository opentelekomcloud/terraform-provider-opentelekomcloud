package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const natGatewayDataSourceName = "data.opentelekomcloud_private_nat_gateway_v3.gateway_1"

func TestAccPrivateNatGatewayV3DS_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatGatewayV3DSBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateNATGatewayV3DataSourceID(natGatewayDataSourceName),
					resource.TestCheckResourceAttr(natGatewayDataSourceName, "gateways.0.name", "test-acc-nat-gateway"),
				),
			},
		},
	})
}

func TestAccPrivateNatGatewayV3DS_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatGatewayV3DSId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(natGatewayDataSourceName, "gateways.0.name", "test-acc-nat-gateway"),
				),
			},
		},
	})
}

func TestAccPrivateNatGatewayV3DS_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatGatewayV3DSName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(natGatewayDataSourceName, "gateways.0.description", "test_success"),
				),
			},
		},
	})
}

func testAccCheckPrivateNATGatewayV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Private NAT Gateway data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Private NAT Gateway data source ID not set")
		}

		return nil
	}
}

var testAccPrivateNatGatewayV3DSBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway" {
  name        = "test-acc-nat-gateway"
  description = "test_success"
  spec        = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

data "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  depends_on = ["opentelekomcloud_private_nat_gateway_v3.gateway"]
}
`, common.DataSourceSubnet)

var testAccPrivateNatGatewayV3DSId = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway" {
  name        = "test-acc-nat-gateway"
  description = "test_success"
  spec        = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

data "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  depends_on = ["opentelekomcloud_private_nat_gateway_v3.gateway"]
  id         = opentelekomcloud_private_nat_gateway_v3.gateway.id
}
`, common.DataSourceSubnet)

var testAccPrivateNatGatewayV3DSName = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway" {
  name        = "test-acc-nat-gateway"
  description = "test_success"
  spec        = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

data "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  depends_on = ["opentelekomcloud_private_nat_gateway_v3.gateway"]
  name       = "test-acc-nat-gateway"
}
`, common.DataSourceSubnet)
