package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataFlavorName = "data.opentelekomcloud_compute_flavor_v2.flavor_1"

func TestAccComputeV2FlavorDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID(dataFlavorName),
					resource.TestCheckResourceAttr(dataFlavorName, "name", "hl1.8xlarge.8"),
					resource.TestCheckResourceAttr(dataFlavorName, "ram", "262144"),
					resource.TestCheckResourceAttr(dataFlavorName, "disk", "40"),
					resource.TestCheckResourceAttr(dataFlavorName, "vcpus", "32"),
					resource.TestCheckResourceAttr(dataFlavorName, "rx_tx_factor", "1"),
				),
			},
		},
	})
}

func TestAccComputeV2FlavorDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorDataSourceQueryDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID(dataFlavorName),
					resource.TestCheckResourceAttr(dataFlavorName, "name", "hl1.8xlarge.8"),
					resource.TestCheckResourceAttr(dataFlavorName, "ram", "262144"),
					resource.TestCheckResourceAttr(dataFlavorName, "disk", "40"),
					resource.TestCheckResourceAttr(dataFlavorName, "vcpus", "32"),
					resource.TestCheckResourceAttr(dataFlavorName, "rx_tx_factor", "1"),
				),
			},
			{
				Config: testAccComputeV2FlavorDataSourceQueryMinDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID(dataFlavorName),
					resource.TestCheckResourceAttr(dataFlavorName, "name", "hl1.8xlarge.8"),
					resource.TestCheckResourceAttr(dataFlavorName, "ram", "262144"),
					resource.TestCheckResourceAttr(dataFlavorName, "disk", "40"),
					resource.TestCheckResourceAttr(dataFlavorName, "vcpus", "32"),
					resource.TestCheckResourceAttr(dataFlavorName, "rx_tx_factor", "1"),
				),
			},
			{
				Config: testAccComputeV2FlavorDataSourceQueryMinRAM,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID(dataFlavorName),
					resource.TestCheckResourceAttr(dataFlavorName, "name", "hl1.8xlarge.8"),
					resource.TestCheckResourceAttr(dataFlavorName, "ram", "262144"),
					resource.TestCheckResourceAttr(dataFlavorName, "disk", "40"),
					resource.TestCheckResourceAttr(dataFlavorName, "vcpus", "32"),
					resource.TestCheckResourceAttr(dataFlavorName, "rx_tx_factor", "1"),
				),
			},
			{
				Config: testAccComputeV2FlavorDataSourceQueryVCPUs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID(dataFlavorName),
					resource.TestCheckResourceAttr(dataFlavorName, "name", "hl1.8xlarge.8"),
					resource.TestCheckResourceAttr(dataFlavorName, "ram", "262144"),
					resource.TestCheckResourceAttr(dataFlavorName, "disk", "40"),
					resource.TestCheckResourceAttr(dataFlavorName, "vcpus", "32"),
					resource.TestCheckResourceAttr(dataFlavorName, "rx_tx_factor", "1"),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("flavor data source ID not set")
		}

		return nil
	}
}

const testAccComputeV2FlavorDataSourceBasic = `
data "opentelekomcloud_compute_flavor_v2" "flavor_1" {
  name = "hl1.8xlarge.8"
}
`

const testAccComputeV2FlavorDataSourceQueryDisk = `
data "opentelekomcloud_compute_flavor_v2" "flavor_1" {
  disk = 40
}
`

const testAccComputeV2FlavorDataSourceQueryMinDisk = `
data "opentelekomcloud_compute_flavor_v2" "flavor_1" {
  name     = "hl1.8xlarge.8"
  min_disk = 40
}
`

const testAccComputeV2FlavorDataSourceQueryMinRAM = `
data "opentelekomcloud_compute_flavor_v2" "flavor_1" {
  name    = "hl1.8xlarge.8"
  min_ram = 262144
}
`

const testAccComputeV2FlavorDataSourceQueryVCPUs = `
data "opentelekomcloud_compute_flavor_v2" "flavor_1" {
  name = "hl1.8xlarge.8"
  vcpus = 32
}
`
