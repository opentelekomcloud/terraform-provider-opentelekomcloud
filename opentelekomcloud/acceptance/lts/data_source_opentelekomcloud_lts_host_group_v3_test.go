package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccHostGroupV3DataSource_byId(t *testing.T) {
	var (
		dataSource = "data.opentelekomcloud_lts_host_group_v3.test"
		rName      = fmt.Sprintf("lts_group%s", acctest.RandString(3))
		dc         = common.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestLtsPreCheckLts(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostGroupV3_byId(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "host_groups.0.id"),
					resource.TestCheckResourceAttr(dataSource, "host_groups.0.name", rName),
				),
			},
		},
	})
}

func TestAccHostGroupV3DataSource_byName(t *testing.T) {
	var (
		dataSource = "data.opentelekomcloud_lts_host_group_v3.test"
		rName      = fmt.Sprintf("lts_group%s", acctest.RandString(3))
		dc         = common.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestLtsPreCheckLts(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostGroupV3_byName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "host_groups.0.id"),
					resource.TestCheckResourceAttr(dataSource, "host_groups.0.name", rName),
				),
			},
		},
	})
}

func testDataSourceHostGroupV3_byId(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_host_group_v3" "hg" {
  name = "%s"
  type = "linux"
}

data "opentelekomcloud_lts_host_group_v3" "test" {
  id = opentelekomcloud_lts_host_group_v3.hg.id
}
`, name)
}

func testDataSourceHostGroupV3_byName(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_host_group_v3" "hg" {
  name = "%s"
  type = "linux"
}

data "opentelekomcloud_lts_host_group_v3" "test" {
  name = opentelekomcloud_lts_host_group_v3.hg.name
}
`, name)
}
