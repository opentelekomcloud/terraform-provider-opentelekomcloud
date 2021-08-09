package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	vpc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/vpc"
)

func TestAccOpenTelekomCloudRdsFlavorV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudRdsFlavorV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV1DataSourceID("data.opentelekomcloud_rds_flavors_v1.flavor"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_rds_flavors_v1.flavor", "name"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_rds_flavors_v1.flavor", "id"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_rds_flavors_v1.flavor", "speccode"),
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
				Config: testAccOpenTelekomCloudRdsFlavorV1DataSource_speccode,
				Check: resource.ComposeTestCheckFunc(
					vpc.TestAccCheckNetworkingNetworkV2DataSourceID("data.opentelekomcloud_rds_flavors_v1.flavor"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_rds_flavors_v1.flavor", "name", "OTC_PGCM_XLARGE"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_rds_flavors_v1.flavor", "speccode", "rds.pg.s1.xlarge"),
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

var testAccOpenTelekomCloudRdsFlavorV1DataSource_basic = `

data "opentelekomcloud_rds_flavors_v1" "flavor" {
    region = "eu-de"
	datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
}
`

var testAccOpenTelekomCloudRdsFlavorV1DataSource_speccode = `

data "opentelekomcloud_rds_flavors_v1" "flavor" {
    region = "eu-de"
	datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.xlarge"
}
`
