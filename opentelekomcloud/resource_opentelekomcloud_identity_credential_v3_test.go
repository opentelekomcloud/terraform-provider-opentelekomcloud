package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
)

func TestAccIdentityV3Credential_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testIdentityCredentialPrecheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIdentityV3CredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccIdentityV3CredentialBasic, OS_USER_ID),
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
				Config: fmt.Sprintf(testAccIdentityV3CredentialUpdate, OS_USER_ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "access"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_identity_credential_v3.aksk", "secret"),
					resource.TestCheckResourceAttr("opentelekomcloud_identity_credential_v3.aksk", "status", "inactive"),
				),
			},
		},
	})
}

func testIdentityCredentialPrecheck(t *testing.T) {
	if OS_USER_ID == "" {
		t.Error("OS_USER_ID is required for credentials test")
	}
}

func testAccIdentityV3CredentialDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
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

const (
	testAccIdentityV3CredentialBasic = `

resource opentelekomcloud_identity_credential_v3 aksk {
  user_id = "%s"
  description = "This is one and unique test AK/SK"
}
`
	testAccIdentityV3CredentialUpdate = `
resource opentelekomcloud_identity_credential_v3 aksk {
  user_id = "%s"
  description = "This is one and unique test AK/SK"
  status  = "inactive"
}
`
)
