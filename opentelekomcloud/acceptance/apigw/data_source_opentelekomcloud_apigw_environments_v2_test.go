package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccDataSourceEnvironments_basic(t *testing.T) {
	var (
		dataSourceName = "data.opentelekomcloud_apigw_environments_v2.test"
		dc             = common.InitDataSourceCheck(dataSourceName)
		name           = common.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironments_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataSourceEnvironments_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_apigw_environment_v2" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by script"
}

data "opentelekomcloud_apigw_environments_v2" "test" {
  depends_on = [opentelekomcloud_apigw_environment_v2.test]

  instance_id = "%[1]s"
  name        = opentelekomcloud_apigw_environment_v2.test.name
}
`, env.OS_APIGW_GATEWAY_ID, name)
}
