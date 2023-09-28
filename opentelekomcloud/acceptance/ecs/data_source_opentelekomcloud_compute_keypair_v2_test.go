package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const publicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALRzbIOR9HUYNwfKtII/et98eGXDJhf8YxHf9BtRdAU"

func TestAccComputeV2KeyPairDataSource_basic(t *testing.T) {
	resourceName := "data.opentelekomcloud_compute_keypair_v2.key_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeyPairDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeKeyPairV2DataSourceName(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "key_1"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
				),
			},
			{
				Config: testAccComputeV2KeyPairDataSourceRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeKeyPairV2DataSourceName(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "key_1"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
				),
			},
		},
	})
}

func testAccCheckComputeKeyPairV2DataSourceName(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find compute keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("compute keypair data source name not set")
		}

		return nil
	}
}

var testAccComputeV2KeyPairDataSource_keytest = fmt.Sprintf(`
resource "opentelekomcloud_compute_keypair_v2" "kp_1" {
  name       = "key_1"
  public_key = "%s"
}
`, publicKey)

var testAccComputeV2KeyPairDataSourceBasic = fmt.Sprintf(`
%s

data "opentelekomcloud_compute_keypair_v2" "key_1" {
  name = "key_1"

  depends_on = [opentelekomcloud_compute_keypair_v2.kp_1]
}
`, testAccComputeV2KeyPairDataSource_keytest)

var testAccComputeV2KeyPairDataSourceRegex = fmt.Sprintf(`
%s

data "opentelekomcloud_compute_keypair_v2" "key_1" {
  name_regex = "^key.+"

  depends_on = [opentelekomcloud_compute_keypair_v2.kp_1]
}
`, testAccComputeV2KeyPairDataSource_keytest)
