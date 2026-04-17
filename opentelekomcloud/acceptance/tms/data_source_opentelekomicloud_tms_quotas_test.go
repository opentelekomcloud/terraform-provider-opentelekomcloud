package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccTmsQuotasDS_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_tms_quotas_v1.quotas_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTmsQuotasDS_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTMSQuotasV1DsId(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "quotas.0.quota_key"),
				),
			},
		},
	})
}

func testAccCheckTMSQuotasV1DsId(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TMS quotas v1 data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TMS quotas v1 data source ID not set")
		}

		return nil
	}
}

const testTmsQuotasDS_basic = `
data "opentelekomcloud_tms_quotas_v1" "quotas_1" {}
`
