package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/snatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const snatRuleResourceName = "opentelekomcloud_private_nat_snat_rule_v3.rule_1"

func getPrivateSnatRule(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 Client: %s", err)
	}
	getResp, err := snatrules.Get(client, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching SNAT rule: %s", err)
	}
	return getResp.SnatRule, nil
}

func TestAccPrivateNatSnatRuleV3_basic(t *testing.T) {
	transitIpId := os.Getenv("OS_TRANSIT_IP_ID")
	if transitIpId == "" {
		t.Skip("OS_TRANSIT_IP_ID is missing but test requires using existing network")
	}

	var gateway snatrules.PrivateSnat
	rc := common.InitResourceCheck(
		snatRuleResourceName,
		&gateway,
		getPrivateSnatRule,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNatSnatRuleV3Basic(transitIpId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(snatRuleResourceName, "description", "created"),
				),
			},
			{
				Config: testAccPrivateNatSnatRuleV3Update(transitIpId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(snatRuleResourceName, "description", "updated"),
				),
			},
			{
				ResourceName:      snatRuleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivateNatSnatRuleV3Basic(transitIpId string) string {
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
`, common.DataSourceSubnet, transitIpId)
}

func testAccPrivateNatSnatRuleV3Update(transitIpId string) string {
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
  description    = "updated"
  transit_ip_ids = ["%s"]
  virsubnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}
`, common.DataSourceSubnet, transitIpId)
}
