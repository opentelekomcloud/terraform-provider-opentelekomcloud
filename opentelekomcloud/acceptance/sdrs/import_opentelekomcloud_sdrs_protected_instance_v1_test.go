package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccSdrsProtectedInstanceV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_sdrs_protected_instance_v1.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccSdrsProtectedInstanceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectedInstanceV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_target_eip", "delete_target_server",
				},
			},
		},
	})
}
