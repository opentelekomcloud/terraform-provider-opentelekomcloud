package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	vpc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/vpc"
)

const flavorDataName = "data.opentelekomcloud_rds_flavors_v1.flavor"

func TestAccOpenTelekomCloudRdsFlavorV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudRdsFlavorV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV1DataSourceID(flavorDataName),
					resource.TestCheckResourceAttrSet(flavorDataName, "name"),
					resource.TestCheckResourceAttrSet(flavorDataName, "id"),
					resource.TestCheckResourceAttrSet(flavorDataName, "speccode"),
				),
			},
		},
	})
}

func TestAccOpenTelekomCloudRdsFlavorV1DataSource_speccode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudRdsFlavorV1DataSourceSpeccode,
				Check: resource.ComposeTestCheckFunc(
					vpc.TestAccCheckNetworkingNetworkV2DataSourceID(flavorDataName),
					resource.TestCheckResourceAttr(flavorDataName, "name", "OTC_PGCM_XLARGE"),
					resource.TestCheckResourceAttr(flavorDataName, "speccode", "rds.pg.s1.xlarge"),
				),
			},
		},
	})
}

func testAccCheckRdsFlavorV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find rds data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("rds data source ID not set ")
		}

		return nil
	}
}

const testAccOpenTelekomCloudRdsFlavorV1DataSourceBasic = `
data "opentelekomcloud_rds_flavors_v1" "flavor" {
  region            = "eu-de"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
}
`

const testAccOpenTelekomCloudRdsFlavorV1DataSourceSpeccode = `
data "opentelekomcloud_rds_flavors_v1" "flavor" {
  region            = "eu-de"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.xlarge"
}
`
