package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDDSInstanceV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSInstanceV3DataSourceID("data.opentelekomcloud_dds_instance_v3.instance"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instance", "name", ""),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instance", "name", ""),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instance", "name", ""),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instance", "name", ""),
				),
			},
		},
	})
}

func testAccCheckDDSInstanceV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find instances data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Node data source ID not set ")
		}

		return nil
	}
}

var testAccDDSInstanceV3DataSource_basic = fmt.Sprint(`
`)
