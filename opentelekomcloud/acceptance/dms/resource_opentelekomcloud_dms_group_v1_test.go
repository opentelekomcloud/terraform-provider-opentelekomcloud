package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/groups"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGroupName = "opentelekomcloud_dms_group_v1.group_1"

func TestAccDmsGroupsV1_basic(t *testing.T) {
	var group groups.Group
	var groupName = fmt.Sprintf("dms_group_%s", acctest.RandString(5))
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1GroupBasic(groupName, queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1GroupExists(resourceGroupName, group),
					resource.TestCheckResourceAttr(resourceGroupName, "name", groupName),
				),
			},
		},
	})
}

func testAccCheckDmsV1GroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DmsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dms_group_v1" {
			continue
		}

		queueID := rs.Primary.Attributes["queue_id"]
		page, err := groups.List(client, queueID, false).AllPages()
		if err == nil {
			groupsList, err := groups.ExtractGroups(page)
			if err != nil {
				return fmt.Errorf("error getting groups in queue %s: %s", queueID, err)
			}
			if len(groupsList) > 0 {
				for _, group := range groupsList {
					if group.ID == rs.Primary.ID {
						return fmt.Errorf("dms group still exists")
					}
				}
			}
		}
	}
	return nil
}

func testAccCheckDmsV1GroupExists(n string, group groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DmsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
		}

		queueID := rs.Primary.Attributes["queue_id"]
		page, err := groups.List(client, queueID, false).AllPages()
		if err != nil {
			return fmt.Errorf("error getting groups in queue %s: %s", queueID, err)
		}

		groupsList, err := groups.ExtractGroups(page)
		if err != nil {
			return fmt.Errorf("error extracting groups: %w", err)
		}
		if len(groupsList) > 0 {
			for _, found := range groupsList {
				if found.ID == rs.Primary.ID {
					group = found
					return nil
				}
			}
		}
		return fmt.Errorf("dms group not found")
	}
}

func testAccDmsV1GroupBasic(groupName string, queueName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_queue_v1" "queue_1" {
  name = "%s"
}

resource "opentelekomcloud_dms_group_v1" "group_1" {
  name     = "%s"
  queue_id = opentelekomcloud_dms_queue_v1.queue_1.id
}
`, queueName, groupName)
}
