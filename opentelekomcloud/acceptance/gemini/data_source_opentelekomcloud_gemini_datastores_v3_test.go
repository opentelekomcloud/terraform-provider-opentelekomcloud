package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceGeminiDatastoresName = "data.opentelekomcloud_gemini_datastores_v3.test"

func TestAccDataSourceGeminiDatastores_cassandra(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiDatastoresCassandra,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDatastoresDataSourceID(dataSourceGeminiDatastoresName),
					resource.TestCheckResourceAttrSet(dataSourceGeminiDatastoresName, "versions.#"),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "storage_engines.0", "rocksDB"),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "versions.0", "3.11"),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "modes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "modes.0", "Cluster"),
				),
			},
		},
	})
}

func TestAccDataSourceGeminiDatastores_influx(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiDatastoresInflux,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDatastoresDataSourceID(dataSourceGeminiDatastoresName),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "versions.0", "1.7"),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "modes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceGeminiDatastoresName, "modes.0", "CloudNativeCluster"),
				),
			},
		},
	})
}

func testAccCheckGeminiDatastoresDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB datastores data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GeminiDB datastores data source ID not set")
		}

		return nil
	}
}

const testDataSourceGeminiDatastoresCassandra = `
data "opentelekomcloud_gemini_datastores_v3" "test" {
  engine_name = "cassandra"
}
`

const testDataSourceGeminiDatastoresInflux = `
data "opentelekomcloud_gemini_datastores_v3" "test" {
  engine_name = "influxdb"
}
`
