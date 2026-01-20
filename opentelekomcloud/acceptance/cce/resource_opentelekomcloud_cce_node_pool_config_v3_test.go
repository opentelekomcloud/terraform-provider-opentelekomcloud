package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCENodePoolConfigV3_basic(t *testing.T) {
	clusterId := os.Getenv("OS_CLUSTER_ID")
	nodepoolId := os.Getenv("OS_NODE_POOL_ID")
	if clusterId == "" || nodepoolId == "" {
		t.Skip("OS_CLUSTER_ID or OS_NODE_POOL_ID required for the test are not set")
	}
	nodePoolResourceName := "opentelekomcloud_cce_node_pool_config_v3.node_pool_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolConfigV3Basic(clusterId, nodepoolId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(nodePoolResourceName, "name", "configuration"),
				),
			},
		},
	})
}

func testAccCCENodePoolConfigV3Basic(clusterId, nodepoolId string) string {
	return fmt.Sprintf(`


resource "opentelekomcloud_cce_node_pool_config_v3" "node_pool_config" {
  cluster_id  = "%s"
  nodepool_id = "%s"
  name        = "configuration"

  packages {
    name = "kubelet"
    configurations {
      name  = "system-reserved-mem"
      value = 600
    }
  }
}`, clusterId, nodepoolId)
}
