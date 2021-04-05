package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccSFSShareAccessRuleV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_sfs_share_access_rule_v2.sfs_rules"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckSFSShareAccessRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareAccessRuleV2_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
