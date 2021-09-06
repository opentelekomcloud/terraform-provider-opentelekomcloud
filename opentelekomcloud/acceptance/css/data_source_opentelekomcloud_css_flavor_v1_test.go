package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataFlavorName = "data.opentelekomcloud_css_flavor_v1.flavor"

func TestAccCSSFlavorV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSFlavorV1DataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCSSFlavorV1DataSourceID(dataFlavorName),
					resource.TestCheckResourceAttrSet(dataFlavorName, "name"),
					resource.TestCheckResourceAttrSet(dataFlavorName, "region"),
					resource.TestCheckResourceAttrSet(dataFlavorName, "ram"),
					resource.TestCheckResourceAttrSet(dataFlavorName, "cpu"),
				),
			},
		},
	})
}

func TestAccCSSFlavorV1DataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSFlavorV1DataSourceByName,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCSSFlavorV1DataSourceID(dataFlavorName),
					resource.TestCheckResourceAttrSet(dataFlavorName, "name"),
					resource.TestCheckResourceAttrSet(dataFlavorName, "region"),
				),
			},
		},
	})
}

func testAccCheckCSSFlavorV1DataSourceID(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("can't find backup data source: %s ", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

const (
	testAccCSSFlavorV1DataSource = `
data "opentelekomcloud_css_flavor_v1" "flavor" {
  min_cpu = 4
  min_ram = 32

  disk_range {
    min_from = 320
    min_to   = 800
  }
}
`

	testAccCSSFlavorV1DataSourceByName = `
data "opentelekomcloud_css_flavor_v1" "flavor" {
  name = "css.medium.8"
}
`
)
