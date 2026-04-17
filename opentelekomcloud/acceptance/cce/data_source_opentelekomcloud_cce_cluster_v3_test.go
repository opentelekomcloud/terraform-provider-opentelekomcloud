package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	var cceName = fmt.Sprintf("cce-test-%s", acctest.RandString(5))
	t.Parallel()
	quotas.BookOne(t, quotas.CCEClusterQuota)

	var (
		byName   = "data.opentelekomcloud_cce_cluster_v3.by_name"
		dcByName = common.InitDataSourceCheck(byName)

		byId   = "data.opentelekomcloud_cce_cluster_v3.by_id"
		dcById = common.InitDataSourceCheck(byName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3DataSourceBasic(cceName),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttr(byName, "name", cceName),
					resource.TestCheckResourceAttr(byName, "status", "Available"),
					resource.TestCheckResourceAttr(byName, "cluster_type", "VirtualMachine"),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "name", cceName),
					resource.TestCheckResourceAttr(byId, "status", "Available"),
					resource.TestCheckResourceAttr(byId, "cluster_type", "VirtualMachine"),
				),
			},
		},
	})
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

data "opentelekomcloud_cce_cluster_v3" "by_name" {
  name = opentelekomcloud_cce_cluster_v3.cluster_1.name
}

data "opentelekomcloud_cce_cluster_v3" "by_id" {
  id = opentelekomcloud_cce_cluster_v3.cluster_1.id
}
`, common.DataSourceSubnet, cceName)
}
