package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/natgateway"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const natGatewayResourceName = "opentelekomcloud_private_nat_gateway_v3.gateway_1"

func getPrivateNatGateway(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 Client: %s", err)
	}
	getResp, err := natgateway.Get(client, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching NAT Gateway: %s", err)
	}
	return getResp.Gateway, nil
}

func TestAccPrivateNatGatewayV3_basic(t *testing.T) {
	var gateway natgateway.PrivateNATGateway
	rc := common.InitResourceCheck(
		natGatewayResourceName,
		&gateway,
		getPrivateNatGateway,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatGatewayV3Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(natGatewayResourceName, "name", "test-acc-nat-gateway"),
					resource.TestCheckResourceAttr(natGatewayResourceName, "description", "created"),
				),
			},
			{
				Config: testAccPrivateNatGatewayV3Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(natGatewayResourceName, "name", "test-acc-nat-gateway-updated"),
					resource.TestCheckResourceAttr(natGatewayResourceName, "description", "updated"),
				),
			},
			{
				ResourceName:      natGatewayResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccPrivateNatGatewayV3Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name        = "test-acc-nat-gateway"
  description = "created"
  spec        = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet)

var testAccPrivateNatGatewayV3Update = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name        = "test-acc-nat-gateway-updated"
  description = "updated"
  spec        = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet)
