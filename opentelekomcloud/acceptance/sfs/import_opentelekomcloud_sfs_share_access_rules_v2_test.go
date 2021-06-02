package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccSFSShareAccessRulesV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_sfs_share_access_rules_v2.sfs_rules"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSFSShareAccessRulesV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareAccessRulesV2_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
