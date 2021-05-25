package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCENodePoolV3ImportBasic(t *testing.T) {
	resourceName := "opentelekomcloud_cce_node_pool_v3.node_pool"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckCCENodePoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodePoolV3ImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"max_node_count", "min_node_count", "priority", "scale_down_cooldown_time",
				},
			},
		},
	})
}

func testAccCCENodePoolV3ImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		var nodePoolID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_cce_cluster_v3" {
				clusterID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_cce_node_pool_v3" {
				nodePoolID = rs.Primary.ID
			}
		}
		if clusterID == "" || nodePoolID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, nodePoolID)
		}
		return fmt.Sprintf("%s/%s", clusterID, nodePoolID), nil
	}
}
