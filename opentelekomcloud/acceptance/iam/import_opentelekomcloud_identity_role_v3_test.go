package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccIdentityV3Role_importBasic(t *testing.T) {
	importResourceName := "opentelekomcloud_identity_role_v3.import_role"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleV3_import(acctest.RandString(10)),
			},

			{
				ResourceName:      importResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityRoleV3_import(val string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_role_v3" "import_role" {
  description   = "role"
  display_name  = "custom_role%s"
  display_layer = "domain"
  statement {
    effect = "Allow"
    action = ["ecs:*:list*"]
  }
}`, val)
}
