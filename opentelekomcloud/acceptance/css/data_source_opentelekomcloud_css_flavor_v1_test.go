package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCSSFlavorV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { common.TestAccPreCheck(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSFlavorV1DataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCSSFlavorV1DataSourceID("data.opentelekomcloud_css_flavor_v1.flavor"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_css_flavor_v1.flavor", "name"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_css_flavor_v1.flavor", "region"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_css_flavor_v1.flavor", "ram"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_css_flavor_v1.flavor", "cpu"),
				),
			},
		},
	})
}

func TestAccCSSFlavorV1DataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { common.TestAccPreCheck(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSFlavorV1DataSourceByName,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCSSFlavorV1DataSourceID("data.opentelekomcloud_css_flavor_v1.flavor"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_css_flavor_v1.flavor", "name"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_css_flavor_v1.flavor", "region"),
				),
			},
		},
	})
}

var (
	testAccCSSFlavorV1DataSource = fmt.Sprintf(`
data "opentelekomcloud_css_flavor_v1" "flavor" {
  min_cpu = 4
  min_ram = 32

  disk_range {
    min_from = 320
    min_to   = 800
  }
}
`)

	testAccCSSFlavorV1DataSourceByName = `
data "opentelekomcloud_css_flavor_v1" "flavor" {
  name = "css.medium.8"
}
`
)

func testAccCheckCSSFlavorV1DataSourceID(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find backup data source: %s ", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}
