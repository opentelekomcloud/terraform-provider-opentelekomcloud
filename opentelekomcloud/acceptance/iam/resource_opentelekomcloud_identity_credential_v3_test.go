package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccIdentityV3Credential_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			checkAKSKUnset(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccIdentityV3CredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3CredentialBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "user_id"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "description"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "secret"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "access"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "create_time"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "last_use_time"),
					resource.TestCheckResourceAttr("opentelekomcloud_identity_credential_v3.aksk", "status", "active"),
				),
			},
			{
				Config: testAccIdentityV3CredentialUpdateStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "access"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "secret"),
					resource.TestCheckResourceAttr("opentelekomcloud_identity_credential_v3.aksk", "status", "inactive"),
				),
			},
			{
				Config: testAccIdentityV3CredentialUpdateDescription,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_identity_credential_v3.aksk", "description", "This is one and unique test AK/SK 2"),
				),
			},
		},
	})
}

func testAccIdentityV3CredentialDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.IdentityV3Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating identity v3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_credential_v3" {
			continue
		}

		_, err := credentials.Get(client, rs.Primary.ID).Extract()

		if err == nil {
			return fmt.Errorf("AK/SK still exists")
		}
	}

	return nil
}

func checkAKSKUnset(t *testing.T) {
	if env.OsSecretKey != "" && env.OsAccessKey != "" {
		t.Error("AK/SK should not be set for AK/SK creation test")
	}
}

const (
	testAccIdentityV3CredentialBasic = `
resource opentelekomcloud_identity_credential_v3 aksk {
  description = "This is one and unique test AK/SK"
}
`
	testAccIdentityV3CredentialUpdateStatus = `
resource opentelekomcloud_identity_credential_v3 aksk {
  description = "This is one and unique test AK/SK"
  status  = "inactive"
}
`
	testAccIdentityV3CredentialUpdateDescription = `
resource opentelekomcloud_identity_credential_v3 aksk {
  description = "This is one and unique test AK/SK 2"
  status  = "inactive"
}
`
)
