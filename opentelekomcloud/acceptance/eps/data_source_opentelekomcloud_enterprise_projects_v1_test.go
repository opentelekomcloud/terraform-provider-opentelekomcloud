package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccEnterpriseProjectsDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_enterprise_projects_v1.test"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProjectsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.0.type"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccEnterpriseProjectsDataSource_basic() string {
	return `
data "opentelekomcloud_enterprise_projects_v1" "test" {}

locals {
  enterprise_project_id = data.opentelekomcloud_enterprise_projects_v1.test.enterprise_projects[0].id
}
data "opentelekomcloud_enterprise_projects_v1" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.opentelekomcloud_enterprise_projects_v1.enterprise_project_id_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.opentelekomcloud_enterprise_projects_v1.enterprise_project_id_filter.enterprise_projects[*].id : v == local.enterprise_project_id]
  )
}

locals {
  name = data.opentelekomcloud_enterprise_projects_v1.test.enterprise_projects[0].name
}
data "opentelekomcloud_enterprise_projects_v1" "name_filter" {
  name = local.name
}
output "name_filter_is_useful" {
  value = length(data.opentelekomcloud_enterprise_projects_v1.name_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.opentelekomcloud_enterprise_projects_v1.name_filter.enterprise_projects[*].name : v == local.name]
  )
}

locals {
  status = data.opentelekomcloud_enterprise_projects_v1.test.enterprise_projects[0].status
}
data "opentelekomcloud_enterprise_projects_v1" "status_filter" {
  status = local.status
}
output "status_filter_is_useful" {
  value = length(data.opentelekomcloud_enterprise_projects_v1.status_filter.enterprise_projects) > 0 && alltrue(
	[for v in data.opentelekomcloud_enterprise_projects_v1.status_filter.enterprise_projects[*].status : v == local.status]
  )
}`
}
