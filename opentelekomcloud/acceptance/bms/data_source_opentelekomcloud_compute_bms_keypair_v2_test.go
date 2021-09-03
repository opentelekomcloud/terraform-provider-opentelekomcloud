package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataKeyPairName = "data.opentelekomcloud_compute_bms_keypairs_v2.keypair"

func TestAccOTCBMSV2KeyPairDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccBmsKeyPairPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSV2KeyPairDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSV2KeyPairDataSourceID(dataKeyPairName),
					resource.TestCheckResourceAttr(dataKeyPairName, "name", env.OS_KEYPAIR_NAME),
				),
			},
		},
	})
}

func testAccCheckBMSV2KeyPairDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("keypair data source ID not set")
		}

		return nil
	}
}

var testAccBMSV2KeyPairDataSourceBasic = fmt.Sprintf(`
data "opentelekomcloud_compute_bms_keypairs_v2" "keypair" {
  name = "%s"
}
`, env.OS_KEYPAIR_NAME)

func testAccBmsKeyPairPreCheck(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)
	if env.OS_KEYPAIR_NAME == "" {
		t.Skip("Provide the key pair name")
	}
}
