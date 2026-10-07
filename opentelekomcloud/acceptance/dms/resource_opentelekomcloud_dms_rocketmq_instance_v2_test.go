package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRocketMQInstanceV2Name = "opentelekomcloud_dms_rocketmq_instance_v2.test"

func getDmsRocketMQInstanceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %s", err)
	}
	return instances.Get(client, state.Primary.ID)
}

func TestAccDmsRocketMQInstanceV2_basic(t *testing.T) {
	var instance instances.Instance
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("rocketmq-acc-update-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourceRocketMQInstanceV2Name, &instance, getDmsRocketMQInstanceFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQInstanceV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "name", name),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "description", "rocketmq test"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "engine_version", "5.x"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "flavor_id", "rocketmq.b1.large.1"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "storage_space", "300"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "broker_num", "1"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "status", "RUNNING"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "type", "single.basic"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "configs.#", "1"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "configs.0.name", "fileReservedTime"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "configs.0.value", "72"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceV2Name, "namesrv_address"),
					resource.TestCheckResourceAttrSet(resourceRocketMQInstanceV2Name, "broker_address"),
				),
			},
			{
				Config: testAccDmsRocketMQInstanceV2Update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "name", updateName),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "description", ""),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "storage_space", "400"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "configs.0.value", "73"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceRocketMQInstanceV2Name, "tags.new", "test"),
				),
			},
			{
				ResourceName:      resourceRocketMQInstanceV2Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"configs",
					"used_storage_space",
				},
			},
		},
	})
}

func testAccDmsRocketMQInstanceV2Basic(name string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

resource "opentelekomcloud_dms_rocketmq_instance_v2" "test" {
  name        = "%s"
  description = "rocketmq test"

  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id

  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "5.x"
  flavor_id         = "rocketmq.b1.large.1"
  storage_spec_code = "dms.physical.storage.ultra.v2"
  storage_space     = 300
  broker_num        = 1

  configs {
    name  = "fileReservedTime"
    value = "72"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, name)
}

func testAccDmsRocketMQInstanceV2Update(name string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

resource "opentelekomcloud_dms_rocketmq_instance_v2" "test" {
  name = "%s"

  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id

  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "5.x"
  flavor_id         = "rocketmq.b1.large.1"
  storage_spec_code = "dms.physical.storage.ultra.v2"
  storage_space     = 400
  broker_num        = 1
  enable_acl        = true

  configs {
    name  = "fileReservedTime"
    value = "73"
  }

  tags = {
    foo = "bar_update"
    new = "test"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, name)
}
