package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/nat"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/natgateways"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGatewayName = "opentelekomcloud_nat_gateway_v2.nat_1"

func TestAccNatGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2GatewayBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists(resourceGatewayName),
					resource.TestCheckResourceAttr(resourceGatewayName, "name", "nat_1"),
					resource.TestCheckResourceAttr(resourceGatewayName, "description", "test for terraform"),
					resource.TestCheckResourceAttr(resourceGatewayName, "spec", "1"),
					resource.TestCheckResourceAttr(resourceGatewayName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccNatV2GatewayUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceGatewayName, "name", "nat_1_updated"),
					resource.TestCheckResourceAttr(resourceGatewayName, "description", "nat_1 updated description"),
					resource.TestCheckResourceAttr(resourceGatewayName, "spec", "2"),
					resource.TestCheckResourceAttr(resourceGatewayName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccNatV2Gateway_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2GatewayBasic,
			},
			{
				ResourceName:      resourceGatewayName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNatV2GatewayDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NatV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(nat.ErrCreationClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_nat_gateway_v2" {
			continue
		}

		_, err := natgateways.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("nat gateway still exists")
		}
	}

	return nil
}

func testAccCheckNatV2GatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NatV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(nat.ErrCreationClient, err)
		}

		found, err := natgateways.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("nat gateway not found")
		}

		return nil
	}
}

var testAccNatV2GatewayBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name                = "nat_1"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceSubnet)

var testAccNatV2GatewayUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name                = "nat_1_updated"
  description         = "nat_1 updated description"
  spec                = "2"
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  tags = {
    muh = "value-update"
  }
}
`, common.DataSourceSubnet)
