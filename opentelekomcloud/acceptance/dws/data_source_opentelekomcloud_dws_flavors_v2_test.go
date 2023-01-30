package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDwsFlavorsDataSource_basic(t *testing.T) {
	resourceName := "data.opentelekomcloud_dws_flavors_v2.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "32"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_memory(t *testing.T) {
	resourceName := "data.opentelekomcloud_dws_flavors_v2.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_memory,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.memory", "512"),
				),
			},
		},
	})
}

func testAccCheckDwsFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find dws flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DWS flavors data source ID not set")
		}

		return nil
	}
}

const testAccDwsFlavorsDataSource_basic = `
data "opentelekomcloud_dws_flavors_v2" "test" {
  vcpus = 32
}
`

const testAccDwsFlavorsDataSource_memory = `
data "opentelekomcloud_dws_flavors_v2" "test" {
  memory = 512
}
`
