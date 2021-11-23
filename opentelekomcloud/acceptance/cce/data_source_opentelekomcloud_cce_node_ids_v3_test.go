package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccCceNodeIdsV3DataSource_basic(t *testing.T) {
	var cceName = fmt.Sprintf("cce-test-%s", acctest.RandString(5))
	var cceNodeName = fmt.Sprintf("node-test-%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cce_node_ids_v3.node_ids"

	t.Parallel()
	quotas.BookOne(t, quotas.CCEClusterQuota)
	quotas.BookMany(t, singleNodeQuotas)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCceNodeIdsV3DataSourceBasic(cceName, cceNodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCceNodeIdsV3DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "1"),
				),
			},
		},
	})
}

func testAccCceNodeIdsV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CCE Node data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCceNodeIdsV3DataSourceBasic(cceName string, cceNodeName string) string {
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

resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_1.id
  name              = "%s"
  flavor_id         = "s2.medium.1"
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

data "opentelekomcloud_cce_node_ids_v3" "node_ids" {
  cluster_id = opentelekomcloud_cce_cluster_v3.cluster_1.id

  depends_on = [opentelekomcloud_cce_node_v3.node_1]
}
`, common.DataSourceSubnet, cceName, cceNodeName, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)
}
