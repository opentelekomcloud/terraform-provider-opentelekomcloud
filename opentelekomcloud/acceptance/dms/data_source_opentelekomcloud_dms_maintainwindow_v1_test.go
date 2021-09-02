package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataMaintainWindowName = "data.opentelekomcloud_dms_maintainwindow_v1.maintainwindow1"

func TestAccDmsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsMaintainWindowV1DataSourceID(dataMaintainWindowName),
					resource.TestCheckResourceAttr(dataMaintainWindowName, "seq", "1"),
					resource.TestCheckResourceAttr(dataMaintainWindowName, "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDmsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find dms maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("dms maintainwindow data source ID not set")
		}

		return nil
	}
}

const testAccDmsMaintainWindowV1DataSourceBasic = `
data "opentelekomcloud_dms_maintainwindow_v1" "maintainwindow1" {
  seq = 1
}
`
