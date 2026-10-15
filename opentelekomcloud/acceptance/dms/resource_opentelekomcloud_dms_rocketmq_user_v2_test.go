package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRocketMQUserV2Name = "opentelekomcloud_dms_rocketmq_user_v2.test"

func getDmsRocketMQUserFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %s", err)
	}
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	return users.Get(client, parts[0], parts[1])
}

func TestAccDmsRocketMQUserV2_basic(t *testing.T) {
	var user users.User
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourceRocketMQUserV2Name, &user, getDmsRocketMQUserFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQUserV2Basic(name, "5ecuredPa55w0rd!", "192.168.1.*"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceRocketMQUserV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "access_key", name),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "secret_key", "5ecuredPa55w0rd!"),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "white_remote_address", "192.168.1.*"),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "admin", "true"),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "default_topic_perm", "PUB|SUB"),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "default_group_perm", "SUB"),
				),
			},
			{
				Config: testAccDmsRocketMQUserV2Basic(name, "2ecuredPa22w2rd!", ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "secret_key", "2ecuredPa22w2rd!"),
					resource.TestCheckResourceAttr(resourceRocketMQUserV2Name, "white_remote_address", ""),
				),
			},
			{
				ResourceName:      resourceRocketMQUserV2Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDmsRocketMQUserV2Basic(name, secretKey, whiteRemoteAddress string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

resource "opentelekomcloud_dms_rocketmq_instance_v2" "test" {
  name = "%[3]s"

  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id

  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "5.x"
  flavor_id         = "rocketmq.b1.large.1"
  storage_spec_code = "dms.physical.storage.ultra.v2"
  storage_space     = 300
  broker_num        = 1
  enable_acl        = true
}

resource "opentelekomcloud_dms_rocketmq_user_v2" "test" {
  instance_id          = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  access_key           = "%[3]s"
  secret_key           = "%[4]s"
  white_remote_address = "%[5]s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, name, secretKey, whiteRemoteAddress)
}
