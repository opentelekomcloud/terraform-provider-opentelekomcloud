package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3DataSource_cluster,
			},
			{
				Config: testAccCCEClusterV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3DataSourceID("data.opentelekomcloud_cce_cluster_v3.clusters"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v3.clusters", "name", "opentelekomcloud-cce"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v3.clusters", "status", "Available"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cce_cluster_v3.clusters", "cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

var testAccCCEClusterV3DataSource_cluster = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "opentelekomcloud-cce"
  cluster_type = "VirtualMachine"
  flavor_id = "cce.s1.small"
  cluster_version = "v1.9.2-r2"
  vpc_id = "%s"
  subnet_id = "%s"
  container_network_type = "overlay_l2"
}`, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3DataSource_basic = fmt.Sprintf(`
%s
data "opentelekomcloud_cce_cluster_v3" "clusters" {
  name = "${opentelekomcloud_cce_cluster_v3.cluster_1.name}"
}
`, testAccCCEClusterV3DataSource_cluster)
