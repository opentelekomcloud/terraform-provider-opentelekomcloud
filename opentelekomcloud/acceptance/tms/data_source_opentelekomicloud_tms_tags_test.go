package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccTmsTagDS_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_tms_tags_v1.tags_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTmsTagDS_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTMSTagsV1DsId(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.0.key", "foo"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.1.key", "foo1"),
				),
			},
		},
	})
}

func testAccCheckTMSTagsV1DsId(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TMS tags v1 data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TMS tags v1 data source ID not set")
		}

		return nil
	}
}

const testTmsTagDS_basic = `
resource "opentelekomcloud_tms_tags_v1" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
  tags {
    key   = "foo1"
    value = "bar1"
  }
  tags {
    key   = "one"
    value = "two"
  }
}

data "opentelekomcloud_tms_tags_v1" "tags_1" {
  key = "foo"

  depends_on = ["opentelekomcloud_tms_tags_v1.test"]
}
`
