package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccDataSourceGroups_basic(t *testing.T) {
	var (
		all = "data.opentelekomcloud_apigw_groups_v2.test"
		dc  = common.InitDataSourceCheck(all)

		byId   = "data.opentelekomcloud_apigw_groups_v2.filter_by_id"
		dcById = common.InitDataSourceCheck(byId)

		byName   = "data.opentelekomcloud_apigw_groups_v2.filter_by_name"
		dcByName = common.InitDataSourceCheck(byName)

		byNotFoundName   = "data.opentelekomcloud_apigw_groups_v2.filter_by_not_found_name"
		dcByNotFoundName = common.InitDataSourceCheck(byNotFoundName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`[1-9]\d*`)),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "groups.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "groups.0.id", "opentelekomcloud_apigw_group_v2.group", "id"),
					resource.TestCheckResourceAttrPair(byId, "groups.0.name", "opentelekomcloud_apigw_group_v2.group", "name"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.status"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.sl_domain"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.updated_at"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.on_sell_status"),
					resource.TestCheckResourceAttr(byId, "groups.0.sl_domains.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.description"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.is_default"),
					resource.TestCheckResourceAttr(byId, "groups.0.environment.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "groups.0.environment.0.environment_id", "opentelekomcloud_apigw_environment_v2.env", "id"),
					resource.TestCheckResourceAttr(byId, "groups.0.environment.0.variable.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.environment.0.variable.0.name"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.environment.0.variable.0.value"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.environment.0.variable.0.id"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGroups_base() string {
	name := fmt.Sprintf("gr_%s", acctest.RandString(10))

	return fmt.Sprintf(`
resource "opentelekomcloud_apigw_environment_v2" "env"{
  name        = "%[2]s"
  instance_id = "%[1]s"
  description = "test description"
}

resource "opentelekomcloud_apigw_group_v2" "group"{
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "test description"

  environment {
	variable {
	  name  = "test-name"
	  value = "test-value"
	}
  	environment_id = opentelekomcloud_apigw_environment_v2.env.id
  }
}
`, env.OS_APIGW_GATEWAY_ID, "group_"+name)
}

func testAccDataSourceGroups_basic() string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_apigw_groups_v2" "test" {
  depends_on = [opentelekomcloud_apigw_group_v2.group]

  instance_id = "%[2]s"
}

# Filter by ID
locals {
  group_id = opentelekomcloud_apigw_group_v2.group.id
}

data "opentelekomcloud_apigw_groups_v2" "filter_by_id" {
  instance_id = "%[2]s"
  group_id    = local.group_id
}

locals {
  id_filter_result = [
    for v in data.opentelekomcloud_apigw_groups_v2.filter_by_id.groups[*].id : v == local.group_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  group_name = opentelekomcloud_apigw_group_v2.group.name
}

data "opentelekomcloud_apigw_groups_v2" "filter_by_name" {
  depends_on = [opentelekomcloud_apigw_group_v2.group]

  instance_id = "%[2]s"
  name        = local.group_name
}

locals {
  name_filter_result = [
    for v in data.opentelekomcloud_apigw_groups_v2.filter_by_name.groups[*].name : v == local.group_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name and the name is not exist
data "opentelekomcloud_apigw_groups_v2" "filter_by_not_found_name" {
  depends_on = [opentelekomcloud_apigw_group_v2.group]

  instance_id = "%[2]s"
  name        = "not_found_name"
}

output "is_not_found_name_filter_useful" {
  value = length(data.opentelekomcloud_apigw_groups_v2.filter_by_not_found_name.groups) == 0
}
`, testAccDataSourceGroups_base(), env.OS_APIGW_GATEWAY_ID)
}
