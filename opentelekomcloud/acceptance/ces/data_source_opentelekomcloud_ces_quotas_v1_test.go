package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const quotasdataSourceName = "data.opentelekomcloud_ces_quotas_v1.quotas_1"

func TestAccCESQuotasV1DataSource_basic(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCESQuotasBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCESQuotasV1DataSourceID(quotasdataSourceName),
					resource.TestCheckResourceAttr(quotasdataSourceName, "quotas.0.resources.0.type", "alarm"),
				),
			},
		},
	})
}

func testAccCheckCESQuotasV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CES Quotas data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CES quotas data source ID not set")
		}

		return nil
	}
}

var testCESQuotasBasic = `
data "opentelekomcloud_ces_quotas_v1" "quotas_1" {}
`
