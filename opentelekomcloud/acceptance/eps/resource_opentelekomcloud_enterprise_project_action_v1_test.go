package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccEpsAction_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEpsAction_basic(),
			},
			{
				Config: testAccEpsAction_basic_update(),
			},
		},
	})
}

func testAccEpsAction_basic() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_enterprise_project_action_v1" "disable" {
  enterprise_project_id = "%[1]s"
  action                = "disable"
}
`, env.OS_ENTERPRISE_PROJECT_ID)
}

func testAccEpsAction_basic_update() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_enterprise_project_action_v1" "enable" {
  enterprise_project_id = "%[1]s"
  action                = "enable"
}
`, env.OS_ENTERPRISE_PROJECT_ID)
}
