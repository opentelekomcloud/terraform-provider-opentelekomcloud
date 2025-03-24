package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCEClustersDataSource_basic(t *testing.T) {
	dsName := "data.opentelekomcloud_cce_clusters_v3.test"
	dc := common.InitDataSourceCheck(dsName)
	rName := fmt.Sprintf("cce-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClustersDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dsName, "clusters.0.name", rName),
					resource.TestCheckResourceAttr(dsName, "clusters.0.status", "Available"),
					resource.TestCheckResourceAttr(dsName, "clusters.0.cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCCEClustersDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_cce_clusters_v3" "test" {
  name = opentelekomcloud_cce_cluster_v3.cluster_1.name

  depends_on = [opentelekomcloud_cce_cluster_v3.cluster_1]
}
`, testAccCCEClusterV3Basic(rName))
}
