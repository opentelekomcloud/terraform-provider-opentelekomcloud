package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	var cceName = fmt.Sprintf("cce-test-%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cce_cluster_v3.clusters"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3DataSourceBasic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", cceName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "Available"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_type", "VirtualMachine"),
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

func testAccCCEClusterV3DataSourceBasic(cceName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  container_network_type = "overlay_l2"
}

data "opentelekomcloud_cce_cluster_v3" "clusters" {
  name = opentelekomcloud_cce_cluster_v3.cluster_1.name
}
`, common.DataSourceSubnet, cceName)
}
