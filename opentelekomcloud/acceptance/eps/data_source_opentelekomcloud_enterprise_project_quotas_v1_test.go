package eps

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataEnterpriseProjectQuotas_basic(t *testing.T) {
	all := "data.opentelekomcloud_enterprise_project_quotas_v1.test"
	dc := common.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testaccdataenterpriseprojectquotasBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_eps_quota_configured", "true"),
					resource.TestCheckOutput("is_total_quota_number_configured", "true"),
					resource.TestCheckOutput("is_type_configured", "true"),
					resource.TestCheckOutput("is_quota_used_number_configured", "true"),
				),
			},
		},
	})
}

const testaccdataenterpriseprojectquotasBasic = `
data "opentelekomcloud_enterprise_project_quotas_v1" "test" {}

locals {
  enterprise_project_quota = try(element([for o in data.opentelekomcloud_enterprise_project_quotas_v1.test.quotas : o if
  o.type == "enterprise_project"], 0), null)
}

output "is_eps_quota_configured" {
  value = local.enterprise_project_quota != null
}

output "is_total_quota_number_configured" {
  value = lookup(local.enterprise_project_quota, "quota") > 0
}

output "is_type_configured" {
  value = lookup(local.enterprise_project_quota, "type") != ""
}

output "is_quota_used_number_configured" {
  value = lookup(local.enterprise_project_quota, "used") >= 0
}
`
