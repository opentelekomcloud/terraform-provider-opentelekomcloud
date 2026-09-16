package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceGeminiFlavorsName = "data.opentelekomcloud_gemini_flavors_v3.test"

func TestAccDataSourceGeminiFlavors_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiFlavorsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiFlavorsDataSourceID(dataSourceGeminiFlavorsName),
					resource.TestCheckResourceAttrSet(dataSourceGeminiFlavorsName, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSourceGeminiFlavorsName, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceGeminiFlavorsName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSourceGeminiFlavorsName, "flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSourceGeminiFlavorsName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(dataSourceGeminiFlavorsName, "flavors.0.engine_name", "cassandra"),
				),
			},
		},
	})
}

func TestAccDataSourceGeminiFlavors_influx(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiFlavorsInflux,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiFlavorsDataSourceID(dataSourceGeminiFlavorsName),
					resource.TestCheckResourceAttrSet(dataSourceGeminiFlavorsName, "flavors.#"),
					resource.TestCheckResourceAttr(dataSourceGeminiFlavorsName, "flavors.0.engine_name", "influxdb"),
					resource.TestCheckResourceAttr(dataSourceGeminiFlavorsName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttr(dataSourceGeminiFlavorsName, "flavors.0.spec_code", "geminidb.influxdb-geminifs.xlarge.4"),
					resource.TestCheckResourceAttr(dataSourceGeminiFlavorsName, "flavors.0.engine_version", "1.7"),
				),
			},
		},
	})
}

// GeminiDB Influx is only offered with cloud native storage, so querying its flavors without an
// explicit `mode` must return the same specifications as querying them with `CloudNativeCluster`.
func TestAccDataSourceGeminiFlavors_influxDefaultMode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiFlavorsInfluxDefaultMode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiFlavorsDataSourceID(dataSourceGeminiFlavorsName),
					resource.TestCheckResourceAttr(dataSourceGeminiFlavorsName, "flavors.0.engine_name", "influxdb"),
					resource.TestCheckResourceAttrPair(
						dataSourceGeminiFlavorsName, "flavors.0.spec_code",
						"data.opentelekomcloud_gemini_flavors_v3.explicit_mode", "flavors.0.spec_code"),
				),
			},
		},
	})
}

func testAccCheckGeminiFlavorsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GeminiDB flavors data source ID not set")
		}

		return nil
	}
}

const testDataSourceGeminiFlavorsBasic = `
data "opentelekomcloud_gemini_flavors_v3" "test" {
  engine_name = "cassandra"
}
`

const testDataSourceGeminiFlavorsInflux = `
data "opentelekomcloud_gemini_flavors_v3" "test" {
  engine_name = "influxdb"
  mode        = "CloudNativeCluster"
  vcpus       = 4
}
`

const testDataSourceGeminiFlavorsInfluxDefaultMode = `
data "opentelekomcloud_gemini_flavors_v3" "test" {
  engine_name = "influxdb"
}

data "opentelekomcloud_gemini_flavors_v3" "explicit_mode" {
  engine_name = "influxdb"
  mode        = "CloudNativeCluster"
}
`
