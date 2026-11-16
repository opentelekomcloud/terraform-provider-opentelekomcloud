package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/notifications"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccASNotificationsV1_basic(t *testing.T) {
	var notification notifications.Notification
	name := common.RandomAccResourceName()
	dataSourceName := "data.opentelekomcloud_as_notifications_v1.test"
	rc := common.InitResourceCheck(asNotificationResourceName, &notification, getASNotificationFunc)
	dc := common.InitDataSourceCheck(dataSourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccASNotificationPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASNotificationsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "notifications.#", "2"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id",
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
