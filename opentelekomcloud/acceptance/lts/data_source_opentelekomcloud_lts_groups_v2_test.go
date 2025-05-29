package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceGroups_basic(t *testing.T) {
	var (
		dataSource = "data.opentelekomcloud_lts_groups_v2.test"
		rName      = fmt.Sprintf("lts_groups%s", acctest.RandString(5))
		dc         = common.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.ttl_in_days"),
					resource.TestMatchResourceAttr(dataSource, "groups.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_exist_log_group", "true"),
					resource.TestCheckOutput("is_eps_return_and_matched", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_lts_groups_v2" "test" {
  depends_on = [
    opentelekomcloud_lts_group_v2.group
  ]
}

output "is_exist_log_group" {
  value = contains(data.opentelekomcloud_lts_groups_v2.test.groups[*].id, opentelekomcloud_lts_group_v2.group.id)
}

locals {
  eps_filter_result = [for v in data.opentelekomcloud_lts_groups_v2.test.groups :
  v.enterprise_project_id == opentelekomcloud_lts_group_v2.group.enterprise_project_id if v.id == opentelekomcloud_lts_group_v2.group.id]
}

output "is_eps_return_and_matched" {
  value = length(local.eps_filter_result) > 0 && alltrue(local.eps_filter_result)
}
`, testAccLtsGroup_basic(name, 30))
}
