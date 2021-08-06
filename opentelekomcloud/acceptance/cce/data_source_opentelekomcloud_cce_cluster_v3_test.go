package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
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
			return fmt.Errorf("can't find cluster data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

var testAccCCEClusterV3DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                   = "opentelekomcloud-cce"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = "%s"
  subnet_id              = "%s"
  container_network_type = "overlay_l2"
}

data "opentelekomcloud_cce_cluster_v3" "clusters" {
  name = opentelekomcloud_cce_cluster_v3.cluster_1.name
}
`, env.OS_VPC_ID, env.OS_NETWORK_ID)
