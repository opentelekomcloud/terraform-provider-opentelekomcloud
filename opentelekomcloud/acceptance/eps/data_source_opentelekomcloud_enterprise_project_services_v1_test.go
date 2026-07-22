package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataEnterpriseProjectServices_basic(t *testing.T) {
	var (
		dataSourceName = "data.opentelekomcloud_enterprise_project_services_v1.test"
		dc             = common.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectServices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.service_i18n_display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.0.resource_type_i18n_display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.0.regions.#"),

					resource.TestCheckOutput("opentelekomcloud_enterprise_project_services_v1_locale", "true"),
					resource.TestCheckOutput("opentelekomcloud_enterprise_project_services_v1_service", "true"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectServices_basic = `
data "opentelekomcloud_enterprise_project_services_v1" "test" {}

output "opentelekomcloud_enterprise_project_services_v1" {
  value = length(data.opentelekomcloud_enterprise_project_services_v1.test.services) > 0
}

data "opentelekomcloud_enterprise_project_services_v1" "test_locale" {
  locale  = "en-us"
}

output "opentelekomcloud_enterprise_project_services_v1_locale" {
  value = length(data.opentelekomcloud_enterprise_project_services_v1.test_locale.services) > 0
}

data "opentelekomcloud_enterprise_project_services_v1" "test_service" {
  service  = "vpc"
}

output "opentelekomcloud_enterprise_project_services_v1_service" {
  value = length(data.opentelekomcloud_enterprise_project_services_v1.test_service.services) > 0
}
`
