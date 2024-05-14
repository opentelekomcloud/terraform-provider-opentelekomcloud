package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceApigwHistName = "data.opentelekomcloud_apigw_api_history_v2.test"

func TestAccDataApigwApiHistory_basic(t *testing.T) {
	dc := common.InitDataSourceCheck(dataSourceApigwHistName)
	rName := fmt.Sprintf("apigw_acc_api_hist%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApigwApiHistory_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceApigwHistName, "history.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataApigwApiHistory_basic(rName string) string {
	relatedConfig := testAccApiPublishment_basic(rName)

	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_apigw_api_history_v2" "test" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id

  depends_on = ["opentelekomcloud_apigw_api_publishment_v2.pub"]
}
`, relatedConfig)
}
