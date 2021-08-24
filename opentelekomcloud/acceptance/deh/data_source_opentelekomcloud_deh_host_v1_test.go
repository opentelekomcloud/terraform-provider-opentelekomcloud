package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccOTCDedicatedHostV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCDedicatedHostV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostV1DataSourceID("data.opentelekomcloud_deh_host_v1.hosts"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_host_v1.hosts", "name", "test-deh-1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_host_v1.hosts", "auto_placement", "on"),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find deh data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("deh data source ID not set ")
		}

		return nil
	}
}

var testAccOTCDedicatedHostV1DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
	 availability_zone= "%s"
     auto_placement= "on"
     host_type= "h1"
	 name = "test-deh-1"
}
data "opentelekomcloud_deh_host_v1" "hosts" {
  id = opentelekomcloud_deh_host_v1.deh1.id
}
`, env.OsAvailabilityZone)
