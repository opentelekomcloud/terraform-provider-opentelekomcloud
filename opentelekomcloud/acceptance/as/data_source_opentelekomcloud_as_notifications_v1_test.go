package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccASNotificationsV1_basic(t *testing.T) {
	name := common.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccASNotificationPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASNotificationsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASNotificationsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					common.InitDataSourceCheck("data.opentelekomcloud_as_notifications_v1.test").CheckResourceExists(),
					resource.TestCheckResourceAttr("data.opentelekomcloud_as_notifications_v1.test", "notifications.#", "2"),
					resource.TestCheckResourceAttrPair("data.opentelekomcloud_as_notifications_v1.test", "id",
						"opentelekomcloud_as_group_v1.test", "id"),
					resource.TestCheckOutput("notification_fields_match", "true"),
				),
			},
		},
	})
}

func testAccASNotificationsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_as_notifications_v1" "test" {
  depends_on = [
    opentelekomcloud_as_notification_v1.test,
    opentelekomcloud_as_notification_v1.peer,
  ]
  scaling_group_id = opentelekomcloud_as_group_v1.test.id
}

output "notification_fields_match" {
  value = alltrue([
    for notification in data.opentelekomcloud_as_notifications_v1.test.notifications :
    notification.topic_urn != opentelekomcloud_smn_topic_v2.test[0].id || (
      notification.topic_name == opentelekomcloud_as_notification_v1.test.topic_name &&
      notification.topic_scene == opentelekomcloud_as_notification_v1.test.topic_scene
    )
    ]) && contains(
    data.opentelekomcloud_as_notifications_v1.test.notifications[*].topic_urn,
    opentelekomcloud_smn_topic_v2.test[0].id
  )
}
`, testAccASNotificationConfig(name, false))
}
