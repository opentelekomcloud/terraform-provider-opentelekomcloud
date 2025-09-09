package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dnatRuleDataSourceName = "data.opentelekomcloud_private_nat_dnat_rule_v3.rule_1_ds"

func TestAccPrivateNatDnatRuleV3DS_basic(t *testing.T) {
	transitIpId := os.Getenv("OS_TRANSIT_IP_ID")
	networkInterfaceId := os.Getenv("OS_NIC_ID")
	if transitIpId == "" || networkInterfaceId == "" {
		t.Skip("OS_TRANSIT_IP_ID or OS_NIC_ID is missing but test requires using existing network")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatDnatRuleV3DSBasic(transitIpId, networkInterfaceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateNATDnatRuleV3DataSourceID(dnatRuleDataSourceName),
					resource.TestCheckResourceAttr(dnatRuleDataSourceName, "dnat_rules.0.description", "created"),
				),
			},
		},
	})
}

func testAccCheckPrivateNATDnatRuleV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Private NAT DNAT rule data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Private NAT DNAT rule data source ID not set")
		}

		return nil
	}
}

func testAccPrivateNatDnatRuleV3DSBasic(transitIpId, networkInterfaceId string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name = "test-acc-nat-gateway"
  spec = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1" {
  description          = "created"
  transit_ip_id        = "%s"
  network_interface_id = "%s"
  gateway_id           = opentelekomcloud_private_nat_gateway_v3.gateway_1.id
}

data "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1_ds" {
  depends_on = ["opentelekomcloud_private_nat_dnat_rule_v3.rule_1"]
}
`, common.DataSourceSubnet, transitIpId, networkInterfaceId)
}

func TestAccPrivateNatDnatRuleV3DS_id(t *testing.T) {
	transitIpId := os.Getenv("OS_TRANSIT_IP_ID")
	networkInterfaceId := os.Getenv("OS_NIC_ID")
	if transitIpId == "" || networkInterfaceId == "" {
		t.Skip("OS_TRANSIT_IP_ID or OS_NIC_ID is missing but test requires using existing network")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatDnatRuleV3DS_id(transitIpId, networkInterfaceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dnatRuleDataSourceName, "dnat_rules.0.description", "created_id"),
				),
			},
		},
	})
}

func testAccPrivateNatDnatRuleV3DS_id(transitIpId, networkInterfaceId string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name = "test-acc-nat-gateway"
  spec = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1" {
  description          = "created_id"
  transit_ip_id        = "%s"
  network_interface_id = "%s"
  gateway_id           = opentelekomcloud_private_nat_gateway_v3.gateway_1.id
}

data "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1_ds" {
  depends_on = ["opentelekomcloud_private_nat_dnat_rule_v3.rule_1"]
  id         = opentelekomcloud_private_nat_dnat_rule_v3.rule_1.id
}
`, common.DataSourceSubnet, transitIpId, networkInterfaceId)
}
