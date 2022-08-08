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
				Config: testAccIdentityV3CredentialNoPGP,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "access"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "secret"),
					resource.TestCheckResourceAttr("opentelekomcloud_identity_credential_v3.aksk", "status", "active"),
				),
			},
			{
				Config: testAccIdentityV3CredentialUpdateStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "access"),
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
	client, err := config.IdentityV3Client(env.OS_REGION_NAME)
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
	if env.OS_SECRET_KEY != "" && env.OS_ACCESS_KEY != "" {
		t.Skip("AK/SK should not be set for AK/SK creation test")
	}
}

const (
	testAccIdentityV3CredentialBasic = `
resource opentelekomcloud_identity_credential_v3 aksk {
  description = "This is one and unique test AK/SK"
  pgp_key = "mDMEYs6CMxYJKwYBBAHaRw8BAQdAoF5tLs+wvwvF+A5+mAfdGLrwz3bKdKsUA3iuxQcDUZe0JFZsYWRpbWlyIFZzaGl2a292IDxlbnJyb3VAZ21haWwuY29tPoiZBBMWCgBBFiEEDdv9KvPC2UDkSDiCKWyGs27oaG0FAmLOgjMCGwMFCQPCZwAFCwkIBwICIgIGFQoJCAsCBBYCAwECHgcCF4AACgkQKWyGs27oaG1VKgD8CBsP+kSZsuVXsa+OV1l3bOu1Ql1Ep7iTljk6ih+AIUUBAMIt2Xbm4UjN1RW2Z6lOOJqiMTDKZQ3y37bXdC74UKYOuDgEYs6CMxIKKwYBBAGXVQEFAQEHQE8bDiDqHGlFvwadjLXQ/FSYCWt4Hk7uOAoAdRsR+L8pAwEIB4h+BBgWCgAmFiEEDdv9KvPC2UDkSDiCKWyGs27oaG0FAmLOgjMCGwwFCQPCZwAACgkQKWyGs27oaG0FpQEApm+XVsFSaWA/cfoJadrWQwZzG3+Ifnbw/+yWk9FTxZIA/1FTrsAp/5ajz5O+knRQp6gop1lToK1HzFVgQa12gCQL"
}
`
	testAccIdentityV3CredentialNoPGP = `
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
