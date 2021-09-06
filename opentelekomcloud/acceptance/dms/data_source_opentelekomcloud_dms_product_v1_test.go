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
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
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
)
