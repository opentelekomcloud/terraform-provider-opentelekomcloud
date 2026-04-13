package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCINetworksV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("cci-net-%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cci_networks_v2.test"

	dc := common.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCINetworksV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "namespace", rName),
					resource.TestCheckResourceAttr(dataSourceName, "networks.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "networks.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "networks.0.namespace", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.status.0.status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "networks.0.subnets.0.subnet_id",
						"opentelekomcloud_vpc_subnet_v1.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "networks.0.security_group_ids.0",
						"opentelekomcloud_networking_secgroup_v2.test", "id"),
				),
			},
		},
	})
}

func testAccCCINetworksV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_cci_networks_v2" "test" {
  depends_on = [opentelekomcloud_cci_network_v2.test]
  namespace  = opentelekomcloud_cci_namespace_v2.test.name
}
`, testAccV2Network_basic(rName))
}
