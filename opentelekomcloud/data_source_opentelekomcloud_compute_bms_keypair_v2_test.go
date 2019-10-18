package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOTCBMSV2KeyPairDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccBmsKeyPairPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSV2KeyPairDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSV2KeyPairDataSourceID("data.opentelekomcloud_compute_bms_keypairs_v2.keypair"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_compute_bms_keypairs_v2.keypair", "name", OS_KEYPAIR_NAME),
				),
			},
		},
	})
}

func testAccCheckBMSV2KeyPairDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Keypair data source ID not set")
		}

		return nil
	}
}

var testAccBMSV2KeyPairDataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_compute_bms_keypairs_v2" "keypair" {
  name = "%s"
}
`, OS_KEYPAIR_NAME)
