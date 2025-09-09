package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const snatRuleDataSourceName = "data.opentelekomcloud_private_nat_snat_rule_v3.rule_1_ds"

func TestAccPrivateNatSnatRuleV3DS_basic(t *testing.T) {
	transitIpId := os.Getenv("OS_TRANSIT_IP_ID")
	if transitIpId == "" {
		t.Skip("OS_TRANSIT_IP_ID is missing but test requires using existing network")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatSnatRuleV3DSBasic(transitIpId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateNATSnatRuleV3DataSourceID(snatRuleDataSourceName),
					resource.TestCheckResourceAttr(snatRuleDataSourceName, "snat_rules.0.description", "created"),
				),
			},
		},
	})
}

func testAccCheckPrivateNATSnatRuleV3DataSourceID(n string) resource.TestCheckFunc {
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

func testAccPrivateNatSnatRuleV3DSBasic(transitIpId string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name = "test-acc-nat-gateway"
  spec = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_private_nat_snat_rule_v3" "rule_1" {
  gateway_id     = opentelekomcloud_private_nat_gateway_v3.gateway_1.id
  description    = "created"
  transit_ip_ids = ["%s"]
  virsubnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

data "opentelekomcloud_private_nat_snat_rule_v3" "rule_1_ds" {
  depends_on = ["opentelekomcloud_private_nat_snat_rule_v3.rule_1"]
}
`, common.DataSourceSubnet, transitIpId)
}

func TestAccPrivateNatSnatRuleV3DS_id(t *testing.T) {
	transitIpId := os.Getenv("OS_TRANSIT_IP_ID")
	if transitIpId == "" {
		t.Skip("OS_TRANSIT_IP_ID is missing but test requires using existing network")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatSnatRuleV3DS_id(transitIpId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(snatRuleDataSourceName, "snat_rules.0.description", "created_id"),
				),
			},
		},
	})
}

func testAccPrivateNatSnatRuleV3DS_id(transitIpId string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name = "test-acc-nat-gateway"
  spec = "Small"

  downlink_vpcs {
    virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_private_nat_snat_rule_v3" "rule_1" {
  gateway_id     = opentelekomcloud_private_nat_gateway_v3.gateway_1.id
  description    = "created_id"
  transit_ip_ids = ["%s"]
  virsubnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

data "opentelekomcloud_private_nat_snat_rule_v3" "rule_1_ds" {
  depends_on = ["opentelekomcloud_private_nat_snat_rule_v3.rule_1"]
  id         = opentelekomcloud_private_nat_snat_rule_v3.rule_1.id
}
`, common.DataSourceSubnet, transitIpId)
}
