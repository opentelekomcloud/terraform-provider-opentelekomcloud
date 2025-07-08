package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	var (
		all = "data.opentelekomcloud_er_quotas_v3.test"
		dc  = common.InitDataSourceCheck(all)

		byType   = "data.opentelekomcloud_er_quotas_v3.filter_by_type"
		dcByType = common.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_no_filter_useful", "true"),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.available"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.unit"),
					dcByType.CheckResourceExists(),
					resource.TestCheckResourceAttr(byType, "quotas.#", "12"),
					resource.TestMatchResourceAttr(byType, "quotas.0.used", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func testAccDataSourceQuotas_basic() string {
	var (
		name        = fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
		bgpAsNum    = acctest.RandIntRange(64512, 65534)
		baseConfig  = testAccStaticRoute_basic_step1(name, bgpAsNum)
		randUUID, _ = uuid.GenerateUUID()
	)

	return fmt.Sprintf(`
%[1]s

locals {
  all_used_quota_types         = ["er_instance", "route_table", "vpc_attachment", "static_route"]
  instance_used_quota_types    = ["route_table", "vpc_attachment", "static_route"]
  route_table_used_quota_types = ["static_route"]
}

data "opentelekomcloud_er_quotas_v3" "test" {
  depends_on = [opentelekomcloud_er_static_route_v3.source_self]
}

# Expression interpretation:
# 1. [for _, v in quotaList: v if v.used > 0]: Filter the list to find quotas with usage greater than 0.
# 2. setintersection(quotaListA, quotaListB): Find the union of two lists.
# 3. length(setsubtract(quotaListA, quotaListB)) == 0: There are no different elements between two lists.
output "is_no_filter_useful" {
  value = length(setsubtract(setintersection([for _, v in data.opentelekomcloud_er_quotas_v3.test.quotas : v.type if v.used > 0],
  local.all_used_quota_types), local.all_used_quota_types)) == 0
}

# Filter by type
data "opentelekomcloud_er_quotas_v3" "filter_by_type" {
  depends_on = [opentelekomcloud_er_instance_v3.test]

  type = "er_instance"
}
`, baseConfig, randUUID)
}
