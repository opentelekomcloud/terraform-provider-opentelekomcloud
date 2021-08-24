package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/snatrules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNatSnatRule_basic(t *testing.T) {
	resourceName := "opentelekomcloud_nat_gateway_v2.nat_1"
	snatResourceName := "opentelekomcloud_nat_snat_rule_v2.snat_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatV2SnatRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2SnatRuleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists(resourceName),
					testAccCheckNatV2SnatRuleExists(snatResourceName),
				),
			},
		},
	})
}

func testAccCheckNatV2SnatRuleDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NatV2Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NATv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_nat_snat_rule_v2" {
			continue
		}

		_, err := snatrules.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("snat rule still exists")
		}
	}

	return nil
}

func testAccCheckNatV2SnatRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NatV2Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NATv2 client: %w", err)
		}

		found, err := snatrules.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("snat rule not found")
		}

		return nil
	}
}

func testAccNatV2SnatRuleBasic() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name                = "nat_1"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = "%s"
  router_id           = "%s"
}

resource "opentelekomcloud_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = opentelekomcloud_nat_gateway_v2.nat_1.id
  floating_ip_id = opentelekomcloud_networking_floatingip_v2.fip_1.id
  cidr           = "192.168.0.0/24"
  source_type    = 0
}
`, env.OsNetworkID, env.OsRouterID)
}
