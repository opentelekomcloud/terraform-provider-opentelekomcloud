package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/dnatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const dnatRuleResourceName = "opentelekomcloud_private_nat_dnat_rule_v3.rule_1"

func getPrivateDnatRule(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 Client: %s", err)
	}
	getResp, err := dnatrules.Get(client, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching DNAT rule: %s", err)
	}
	return getResp.DnatRule, nil
}

func TestAccPrivateNatDnatRuleV3_basic(t *testing.T) {
	transitIpId := os.Getenv("OS_TRANSIT_IP_ID")
	networkInterfaceId := os.Getenv("OS_NIC_ID")
	if transitIpId == "" || networkInterfaceId == "" {
		t.Skip("OS_TRANSIT_IP_ID or OS_NIC_ID is missing but test requires using existing network")
	}

	var gateway dnatrules.PrivateDnat
	rc := common.InitResourceCheck(
		dnatRuleResourceName,
		&gateway,
		getPrivateDnatRule,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatDnatRuleV3Basic(transitIpId, networkInterfaceId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnatRuleResourceName, "description", "created"),
				),
			},
			{
				Config: testAccPrivateNatDnatRuleV3Update(transitIpId, networkInterfaceId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dnatRuleResourceName, "description", "updated"),
				),
			},
			{
				ResourceName:      dnatRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivateNatDnatRuleV3Basic(transitIpId, networkInterfaceId string) string {
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
`, common.DataSourceSubnet, transitIpId, networkInterfaceId)
}

func testAccPrivateNatDnatRuleV3Update(transitIpId, networkInterfaceId string) string {
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
  description          = "updated"
  transit_ip_id        = "%s"
  network_interface_id = "%s"
  gateway_id           = opentelekomcloud_private_nat_gateway_v3.gateway_1.id
}
`, common.DataSourceSubnet, transitIpId, networkInterfaceId)
}
