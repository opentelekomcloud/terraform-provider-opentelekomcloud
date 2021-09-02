package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataProductName = "data.opentelekomcloud_dms_product_v1.product1"

func TestAccDmsProductV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID(dataProductName),
					resource.TestCheckResourceAttr(dataProductName, "engine", "kafka"),
					resource.TestCheckResourceAttr(dataProductName, "partition_num", "300"),
					resource.TestCheckResourceAttr(dataProductName, "storage", "600"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqSingle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSourceRabbitmqSingle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID(dataProductName),
					resource.TestCheckResourceAttr(dataProductName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(dataProductName, "node_num", "3"),
					resource.TestCheckResourceAttr(dataProductName, "io_type", "normal"),
					resource.TestCheckResourceAttr(dataProductName, "storage", "100"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSourceRabbitmqCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID(dataProductName),
					resource.TestCheckResourceAttr(dataProductName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(dataProductName, "node_num", "5"),
					resource.TestCheckResourceAttr(dataProductName, "storage", "500"),
					resource.TestCheckResourceAttr(dataProductName, "io_type", "high"),
					resource.TestCheckResourceAttr(dataProductName, "storage_spec_code", "dms.physical.storage.high"),
				),
			},
		},
	})
}

func testAccCheckDmsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find dms product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("dms product data source ID not set")
		}

		return nil
	}
}

const (
	testAccDmsProductV1DataSourceBasic = `
data "opentelekomcloud_dms_product_v1" "product1" {
  engine            = "kafka"
  version           = "2.3.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}
`
	testAccDmsProductV1DataSourceRabbitmqSingle = `
data "opentelekomcloud_dms_product_v1" "product1" {
  engine            = "rabbitmq"
  version           = "3.7.0"
  instance_type     = "single"
  node_num          = 3
  storage           = 100
  storage_spec_code = "dms.physical.storage.normal"
}
`
	testAccDmsProductV1DataSourceRabbitmqCluster = `
data "opentelekomcloud_dms_product_v1" "product1" {
  engine            = "rabbitmq"
  version           = "3.7.0"
  instance_type     = "cluster"
  node_num          = 5
  storage           = 500
  storage_spec_code = "dms.physical.storage.high"
}
`
)
