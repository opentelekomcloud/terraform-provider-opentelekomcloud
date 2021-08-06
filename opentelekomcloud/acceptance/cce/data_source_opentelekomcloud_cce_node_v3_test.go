package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccCCENodesV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCCEKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3DataSourceID("data.opentelekomcloud_cce_node_v3.nodes"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_node_v3.nodes", "name", "test-node"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_node_v3.nodes", "flavor_id", "s1.medium"),
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

var testAccCCENodeV3DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                   = "opentelekomcloud-cce"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = "%s"
  subnet_id              = "%s"
  container_network_type = "overlay_l2"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "test-node"
  flavor_id         = "s1.medium"
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
  cluster_id = opentelekomcloud_cce_cluster_v3.cluster_1.id
  node_id    = opentelekomcloud_cce_node_v3.node_1.id
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)
