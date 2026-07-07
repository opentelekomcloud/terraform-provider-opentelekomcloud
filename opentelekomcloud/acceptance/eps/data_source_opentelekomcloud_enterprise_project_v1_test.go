package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccEnterpriseProjectDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_enterprise_project_v1.test"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testaccenterpriseprojectdatasourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "id", "0"),
				),
			},
		},
	})
}

const testaccenterpriseprojectdatasourceBasic = `
data "opentelekomcloud_enterprise_project_v1" "test" {
  name = "default"
}
`
