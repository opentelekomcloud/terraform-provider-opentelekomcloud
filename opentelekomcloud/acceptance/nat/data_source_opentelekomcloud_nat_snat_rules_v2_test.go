package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccDatasourceSnatRules_basic(t *testing.T) {
	var (
		name          = fmt.Sprintf("acc_nat_snat_%s", acctest.RandString(5))
		byGatewayId   = "data.opentelekomcloud_nat_snat_rules_v2.filter_by_gateway_id"
		dcByGatewayId = common.InitDataSourceCheck(byGatewayId)
		byEipId       = "data.opentelekomcloud_nat_snat_rules_v2.filter_by_floating_ip_id"
		dcByEipId     = common.InitDataSourceCheck(byEipId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatRulesDataSource_base(name),
			},
			{
				Config: testAccDatasourceSnatRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byGatewayId, "rules.#", "1"),
					dcByEipId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byEipId, "rules.#", "1"),
				),
			},
		},
	})
}

func testAccSnatRulesDataSource_base(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_nat_gateway_v2" "this" {
  name                = "%[3]s"
  spec                = "1"
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

resource "opentelekomcloud_networking_floatingip_v2" "eip" {}

resource "opentelekomcloud_nat_snat_rule_v2" "test" {
  nat_gateway_id = opentelekomcloud_nat_gateway_v2.this.id
  floating_ip_id = opentelekomcloud_networking_floatingip_v2.eip.id
  source_type    = 0
  cidr           = "192.168.0.0/24"
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name)
}

func testAccDatasourceSnatRules_basic(name string) string {
	relatedConfig := testAccSnatRulesDataSource_base(name)
	return fmt.Sprintf(`
%[1]s

locals {
  gateway_id = opentelekomcloud_nat_snat_rule_v2.test.nat_gateway_id
}

data "opentelekomcloud_nat_snat_rules_v2" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  floating_ip_id = opentelekomcloud_nat_snat_rule_v2.test.floating_ip_id
}

data "opentelekomcloud_nat_snat_rules_v2" "filter_by_floating_ip_id" {
  floating_ip_id = local.floating_ip_id
}
`, relatedConfig)
}
