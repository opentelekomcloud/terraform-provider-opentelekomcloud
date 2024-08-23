package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccDatasourceDnatRules_basic(t *testing.T) {
	var (
		name           = fmt.Sprintf("acc_nat_dnat_%s", acctest.RandString(5))
		baseConfig     = testAccDnatRulesDataSource_base(name)
		dataSourceName = "data.opentelekomcloud_nat_dnat_rules_v2.test"
		dc             = common.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_rule_id"
		dcByRuleId = common.InitDataSourceCheck(byRuleId)

		byGatewayId   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_gateway_id"
		dcByGatewayId = common.InitDataSourceCheck(byGatewayId)

		byProtocol   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_protocol"
		dcByProtocol = common.InitDataSourceCheck(byProtocol)

		byInternalServicePort   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_internal_service_port"
		dcByInternalServicePort = common.InitDataSourceCheck(byInternalServicePort)

		byPortId   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_port_id"
		dcByPortId = common.InitDataSourceCheck(byPortId)

		byPrivateIp   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_private_ip"
		dcByPrivateIp = common.InitDataSourceCheck(byPrivateIp)

		byStatus   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_status"
		dcByStatus = common.InitDataSourceCheck(byStatus)

		byFloatingIpAddress   = "data.opentelekomcloud_nat_dnat_rules_v2.filter_by_floating_ip_address"
		dcByFloatingIpAddress = common.InitDataSourceCheck(byFloatingIpAddress)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDnatRules_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),

					dcByInternalServicePort.CheckResourceExists(),
					resource.TestCheckOutput("internal_service_port_filter_is_useful", "true"),

					dcByPortId.CheckResourceExists(),
					resource.TestCheckOutput("port_id_filter_is_useful", "true"),

					dcByPrivateIp.CheckResourceExists(),
					resource.TestCheckOutput("private_ip_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByFloatingIpAddress.CheckResourceExists(),
					resource.TestCheckOutput("floating_ip_address_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDnatRulesDataSource_base(name string) string {
	return fmt.Sprintf(`

%s

%s

resource "opentelekomcloud_nat_gateway_v2" "this" {
  name                = "%[4]s"
  spec                = "1"
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

resource "opentelekomcloud_networking_floatingip_v2" "eip" {}

resource "opentelekomcloud_networking_port_v2" "this" {
  name       = "dnat_rule_port"
  network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  fixed_ip {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%[3]s"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id

  network {
    port = opentelekomcloud_networking_port_v2.this.id
  }
}

resource "opentelekomcloud_nat_dnat_rule_v2" "test" {
  floating_ip_id        = opentelekomcloud_networking_floatingip_v2.eip.id
  nat_gateway_id        = opentelekomcloud_nat_gateway_v2.this.id
  external_service_port = 80
  protocol              = "tcp"
  port_id               = opentelekomcloud_networking_port_v2.this.id
  internal_service_port = 80
  depends_on            = [opentelekomcloud_compute_instance_v2.instance_1]
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name)
}

func testAccDatasourceDnatRules_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_nat_dnat_rules_v2" "test" {
  depends_on = [
    opentelekomcloud_nat_dnat_rule_v2.test
  ]
}

locals {
  rule_id = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].id
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_rule_id" {
  rule_id = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  gateway_id = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].gateway_id
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  gateway_id_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_gateway_id.rules[*].gateway_id :
    v == local.gateway_id
  ]
}

output "gateway_id_filter_is_useful" {
  value = alltrue(local.gateway_id_filter_result) && length(local.gateway_id_filter_result) > 0
}

locals {
  protocol = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].protocol
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_protocol" {
  protocol = local.protocol
}

locals {
  protocol_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_protocol.rules[*].protocol : v == local.protocol
  ]
}

output "protocol_filter_is_useful" {
  value = alltrue(local.protocol_filter_result) && length(local.protocol_filter_result) > 0
}

locals {
  internal_service_port = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].internal_service_port
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_internal_service_port" {
  internal_service_port = local.internal_service_port
}

locals {
  internal_service_port_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_internal_service_port.rules[*].internal_service_port :
    v == local.internal_service_port
  ]
}

output "internal_service_port_filter_is_useful" {
  value = alltrue(local.internal_service_port_filter_result) && length(local.internal_service_port_filter_result) > 0
}

locals {
  port_id = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].port_id
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_port_id" {
  port_id = local.port_id
}

locals {
  port_id_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_port_id.rules[*].port_id : v == local.port_id
  ]
}

output "port_id_filter_is_useful" {
  value = alltrue(local.port_id_filter_result) && length(local.port_id_filter_result) > 0
}

locals {
  private_ip = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].private_ip
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_private_ip" {
  private_ip = local.private_ip
}

locals {
  private_ip_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_private_ip.rules[*].private_ip : v == local.private_ip
  ]
}

output "private_ip_filter_is_useful" {
  value = alltrue(local.private_ip_filter_result) && length(local.private_ip_filter_result) > 0
}

locals {
  status = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].status
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}


locals {
  floating_ip_address = data.opentelekomcloud_nat_dnat_rules_v2.test.rules[0].floating_ip_address
}

data "opentelekomcloud_nat_dnat_rules_v2" "filter_by_floating_ip_address" {
  floating_ip_address = local.floating_ip_address
}

locals {
  floating_ip_address_filter_result = [
    for v in data.opentelekomcloud_nat_dnat_rules_v2.filter_by_floating_ip_address.rules[*].floating_ip_address :
    v == local.floating_ip_address
  ]
}

output "floating_ip_address_filter_is_useful" {
  value = alltrue(local.floating_ip_address_filter_result) && length(local.floating_ip_address_filter_result) > 0
}
`, baseConfig)
}
