package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/cce/shared"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataSourceNodesName = "data.opentelekomcloud_cce_node_v3.nodes"

func TestAccCCENodesV3DataSource_basic(t *testing.T) {
	var cceNodeName = fmt.Sprintf("node-test-%s", acctest.RandString(5))

	t.Parallel()
	shared.BookCluster(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3DataSourceInit(cceNodeName),
			},
			{
				Config: testAccCCENodeV3DataSourceBasic(cceNodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3DataSourceID(dataSourceNodesName),
					resource.TestCheckResourceAttr(dataSourceNodesName, "name", cceNodeName),
					resource.TestCheckResourceAttr(dataSourceNodesName, "flavor_id", "s2.large.2"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "runtime", "docker"),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find nodes data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("node data source ID not set ")
		}

		return nil
	}
}

func testAccCCENodeV3DataSourceInit(cceNodeName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name              = "%s"
  flavor_id         = "s2.large.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, shared.DataSourceCluster, cceNodeName, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)
}

func testAccCCENodeV3DataSourceBasic(cceNodeName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = data.opentelekomcloud_cce_cluster_v3.cluster.id
  name              = "%s"
  flavor_id         = "s2.large.2"
  availability_zone = "%s"
  key_pair          = "%s"
  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
data "opentelekomcloud_cce_node_v3" "nodes" {
  cluster_id = data.opentelekomcloud_cce_cluster_v3.cluster.id
  node_id    = opentelekomcloud_cce_node_v3.node_1.id
}
`, shared.DataSourceCluster, cceNodeName, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)
}
