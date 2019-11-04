package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCceNodeIdsV3DataSource_basic(t *testing.T) {
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCceNodeIdsV3DataSource_ccenode(cceName),
			},
			{
				Config: testAccCceNodeIdsV3DataSource_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCceNodeIdsV3DataSourceID("data.opentelekomcloud_cce_node_ids_v3.node_ids"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_node_ids_v3.node_ids", "ids.#", "1"),
				),
			},
		},
	})
}

func testAccCceNodeIdsV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cce node data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Cce node data source ID not set")
		}

		return nil
	}
}

func testAccCceNodeIdsV3DataSource_ccenode(cceName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
cluster_id = "${opentelekomcloud_cce_cluster_v3.cluster_1.id}"
  name = "%s"
  flavor_id="s2.medium.1"
  availability_zone= "%s"
  key_pair="%s"
  root_volume {
    size= 40
    volumetype= "SATA"
  }
  data_volumes {
    size= 100
    volumetype= "SATA"
  }
}
`, cceName, OS_VPC_ID, OS_NETWORK_ID, cceName, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}

func testAccCceNodeIdsV3DataSource_basic(cceName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}

resource "opentelekomcloud_cce_node_v3" "node_1" {
cluster_id = "${opentelekomcloud_cce_cluster_v3.cluster_1.id}"
  name = "%s"
  flavor_id="s2.medium.1"
  availability_zone= "%s"
  key_pair="%s"
  root_volume {
    size= 40
    volumetype= "SATA"
  }
  data_volumes {
    size= 100
    volumetype= "SATA"
  }
}

data "opentelekomcloud_cce_node_ids_v3" "node_ids" {
	cluster_id = "${opentelekomcloud_cce_cluster_v3.cluster_1.id}"
}
`, cceName, OS_VPC_ID, OS_NETWORK_ID, cceName, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}
