package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQUsersV2DataSource_basic(t *testing.T) {
	all := "data.opentelekomcloud_dms_rocketmq_users_v2.all"
	byName := "data.opentelekomcloud_dms_rocketmq_users_v2.by_name"
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	dc := common.InitDataSourceCheck(all)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQUsersV2DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("all_contains_users", "true"),
					resource.TestCheckResourceAttr(byName, "users.#", "1"),
					resource.TestCheckResourceAttr(byName, "users.0.access_key", name),
					resource.TestCheckResourceAttr(byName, "users.0.white_remote_address", "192.168.1.*"),
					resource.TestCheckResourceAttr(byName, "users.0.admin", "true"),
					resource.TestCheckResourceAttr(byName, "users.0.default_topic_perm", "PUB|SUB"),
					resource.TestCheckResourceAttr(byName, "users.0.default_group_perm", "SUB"),
				),
			},
		},
	})
}

func testAccDmsRocketMQUsersV2DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_user_v2" "test_2" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  access_key  = "%[2]s-2"
  secret_key  = "5ecuredPa55w0rd!"
}

data "opentelekomcloud_dms_rocketmq_users_v2" "all" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id

  depends_on = [
    opentelekomcloud_dms_rocketmq_user_v2.test,
    opentelekomcloud_dms_rocketmq_user_v2.test_2,
  ]
}

data "opentelekomcloud_dms_rocketmq_users_v2" "by_name" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  access_key  = opentelekomcloud_dms_rocketmq_user_v2.test.access_key

  depends_on = [opentelekomcloud_dms_rocketmq_user_v2.test_2]
}

output "all_contains_users" {
  value = length(setintersection(
    data.opentelekomcloud_dms_rocketmq_users_v2.all.users[*].access_key, ["%[2]s", "%[2]s-2"]
  )) == 2
}
`, testAccDmsRocketMQUserV2Basic(name, "5ecuredPa55w0rd!", "192.168.1.*"), name)
}
