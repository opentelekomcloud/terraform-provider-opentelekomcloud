package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOTCDedicatedHostServerV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccDehServerPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCDedicatedHostServerV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostServerV1DataSourceID("data.opentelekomcloud_deh_server_v1.servers"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_server_v1.servers", "server_id", OS_SERVER_ID),
					resource.TestCheckResourceAttr("data.opentelekomcloud_deh_server_v1.servers", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostServerV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find servers data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server data source ID not set ")
		}

		return nil
	}
}

var testAccOTCDedicatedHostServerV1DataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_deh_server_v1" "servers" {
  dedicated_host_id  = "%s"
  server_id = "%s"
}
`, OS_DEH_HOST_ID, OS_SERVER_ID)
