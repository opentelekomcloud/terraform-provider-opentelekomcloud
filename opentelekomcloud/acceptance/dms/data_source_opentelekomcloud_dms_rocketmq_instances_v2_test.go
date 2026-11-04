package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQInstancesV2DataSource_basic(t *testing.T) {
	byName := "data.opentelekomcloud_dms_rocketmq_instances_v2.by_name"
	byID := "data.opentelekomcloud_dms_rocketmq_instances_v2.by_id"
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	dc := common.InitDataSourceCheck(byName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQInstancesV2DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(byName, "instances.#", "1"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.id", resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(byName, "instances.0.name", name),
					resource.TestCheckResourceAttr(byName, "instances.0.description", "rocketmq test"),
					resource.TestCheckResourceAttr(byName, "instances.0.engine_version", "5.x"),
					resource.TestCheckResourceAttr(byName, "instances.0.flavor_id", "rocketmq.b1.large.1"),
					resource.TestCheckResourceAttr(byName, "instances.0.type", "single.basic"),
					resource.TestCheckResourceAttr(byName, "instances.0.status", "RUNNING"),
					resource.TestCheckResourceAttr(byName, "instances.0.storage_space", "300"),
					resource.TestCheckResourceAttr(byName, "instances.0.broker_num", "1"),
					resource.TestCheckResourceAttr(byName, "instances.0.storage_spec_code", "dms.physical.storage.ultra.v2"),
					resource.TestCheckResourceAttr(byName, "instances.0.available_zones.#", "1"),
					resource.TestCheckResourceAttr(byName, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.vpc_id", resourceRocketMQInstanceV2Name, "vpc_id"),
					resource.TestCheckResourceAttrPair(byName, "instances.0.namesrv_address", resourceRocketMQInstanceV2Name, "namesrv_address"),
					resource.TestCheckResourceAttr(byID, "instances.#", "1"),
					resource.TestCheckResourceAttrPair(byID, "instances.0.id", resourceRocketMQInstanceV2Name, "id"),
				),
			},
		},
	})
}

func testAccDmsRocketMQInstancesV2DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dms_rocketmq_instances_v2" "by_name" {
  name             = opentelekomcloud_dms_rocketmq_instance_v2.test.name
  exact_match_name = true
}

data "opentelekomcloud_dms_rocketmq_instances_v2" "by_id" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
}
`, testAccDmsRocketMQInstanceV2Basic(name))
}
