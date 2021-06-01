package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDms(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsMaintainWindowV1DataSourceID("data.opentelekomcloud_dms_maintainwindow_v1.maintainwindow1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_maintainwindow_v1.maintainwindow1", "seq", "1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_maintainwindow_v1.maintainwindow1", "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDmsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dms maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dms maintainwindow data source ID not set")
		}

		return nil
	}
}

var testAccDmsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_dms_maintainwindow_v1" "maintainwindow1" {
  seq = 1
}
`)
