package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccDataSourceResourceTypes_basic(t *testing.T) {
	var (
		byServiceName     = "data.opentelekomcloud_tms_resource_types_v1.filter_by_service_name"
		serviceNotFound   = "data.opentelekomcloud_tms_resource_types_v1.not_found"
		dcByServiceName   = common.InitDataSourceCheck(byServiceName)
		dcServiceNotFound = common.InitDataSourceCheck(serviceNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dcByServiceName.CheckResourceExists(),
					resource.TestCheckOutput("is_service_name_filter_useful", "true"),
					dcServiceNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceResourceTypes_basic = `
data "opentelekomcloud_tms_resource_types_v1" "filter_by_service_name" {
  service_name = "dli"
}

data "opentelekomcloud_tms_resource_types_v1" "not_found" {
  service_name = "not_found"
}

locals {
  filter_result = [for v in data.opentelekomcloud_tms_resource_types_v1.filter_by_service_name.types[*].service_name : v == "dli"]
}

output "is_service_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.opentelekomcloud_tms_resource_types_v1.not_found.types) == 0
}
`

func TestAccDataSourceResourceTypes_filterByRegion(t *testing.T) {
	var (
		byRegion         = "data.opentelekomcloud_tms_resource_types_v1.filter_by_region"
		regionNotFound   = "data.opentelekomcloud_tms_resource_types_v1.not_found"
		dcByRegion       = common.InitDataSourceCheck(byRegion)
		dcRegionNotFound = common.InitDataSourceCheck(regionNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTypes_filterByRegion(),
				Check: resource.ComposeTestCheckFunc(
					dcByRegion.CheckResourceExists(),
					resource.TestCheckOutput("is_region_filter_useful", "true"),
					dcRegionNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceResourceTypes_filterByRegion() string {
	return fmt.Sprintf(`
data "opentelekomcloud_tms_resource_types_v1" "filter_by_region" {
  region = "%[1]s"
}

data "opentelekomcloud_tms_resource_types_v1" "not_found" {
  region = "not_found"
}

output "is_region_filter_useful" {
  value = length(data.opentelekomcloud_tms_resource_types_v1.filter_by_region.types) > 0
}

output "not_found_validation_pass" {
  value = length(data.opentelekomcloud_tms_resource_types_v1.not_found.types) == 0
}
`, env.OS_REGION_NAME)
}
