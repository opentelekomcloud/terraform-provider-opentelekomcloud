package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQFlavorsV2DataSource_basic(t *testing.T) {
	all := "data.opentelekomcloud_dms_rocketmq_flavors_v2.all"
	byID := "data.opentelekomcloud_dms_rocketmq_flavors_v2.by_id"
	dc := common.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQFlavorsV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "versions.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestMatchResourceAttr(all, "flavors.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttr(byID, "flavors.#", "1"),
					resource.TestCheckResourceAttr(byID, "flavors.0.id", "rocketmq.b1.large.1"),
					resource.TestCheckResourceAttr(byID, "flavors.0.type", "single.basic"),
					resource.TestCheckResourceAttr(byID, "flavors.0.properties.0.engine_versions", "5.x"),
					resource.TestCheckResourceAttr(byID, "flavors.0.properties.0.broker_num", "1"),
					resource.TestCheckResourceAttrSet(byID, "flavors.0.vm_specification"),
					resource.TestCheckResourceAttrSet(byID, "flavors.0.properties.0.max_storage"),
					resource.TestCheckOutput("filter_validation", "true"),
				),
			},
		},
	})
}

const testAccDmsRocketMQFlavorsV2DataSourceBasic = `
data "opentelekomcloud_dms_rocketmq_flavors_v2" "all" {}

data "opentelekomcloud_dms_rocketmq_flavors_v2" "by_id" {
  flavor_id = "rocketmq.b1.large.1"
}

locals {
  refer = data.opentelekomcloud_dms_rocketmq_flavors_v2.by_id.flavors[0]
}

data "opentelekomcloud_dms_rocketmq_flavors_v2" "filtered" {
  type               = local.refer.type
  arch_type          = local.refer.arch_types[0]
  storage_spec_code  = local.refer.ios[0].storage_spec_code
  availability_zones = local.refer.ios[0].availability_zones
}

output "filter_validation" {
  value = alltrue([
    for f in data.opentelekomcloud_dms_rocketmq_flavors_v2.filtered.flavors :
    f.type == local.refer.type && contains(f.arch_types, local.refer.arch_types[0]) &&
    alltrue([for io in f.ios : io.storage_spec_code == local.refer.ios[0].storage_spec_code])
  ]) && length(data.opentelekomcloud_dms_rocketmq_flavors_v2.filtered.flavors) > 0
}
`
