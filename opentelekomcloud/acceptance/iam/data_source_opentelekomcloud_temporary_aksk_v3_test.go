package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccIdentityTemporaryAKSKV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityTemporaryAKSKV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityTemporaryAKSKV3DataSourceID("data.opentelekomcloud_identity_temporary_aksk_v3.creds"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "access"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "secret"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "security_token"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "expires_at"),
				),
			},
		},
	})
}

func TestAccIdentityTemporaryAKSKV3DataSource_duration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityTemporaryAKSKV3DataSource_duration,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityTemporaryAKSKV3DataSourceID("data.opentelekomcloud_identity_temporary_aksk_v3.creds"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "duration_seconds", "3600"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "access"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "secret"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "security_token"),
					resource.TestCheckResourceAttrSet(
						"data.opentelekomcloud_identity_temporary_aksk_v3.creds", "expires_at"),
				),
			},
		},
	})
}

func testAccCheckIdentityTemporaryAKSKV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find temporary credentials data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("temporary credentials data source ID not set")
		}

		return nil
	}
}

const testAccIdentityTemporaryAKSKV3DataSource_basic = `
data "opentelekomcloud_identity_temporary_aksk_v3" "creds" {}
`

const testAccIdentityTemporaryAKSKV3DataSource_duration = `
data "opentelekomcloud_identity_temporary_aksk_v3" "creds" {
  duration_seconds = 3600
}
`
