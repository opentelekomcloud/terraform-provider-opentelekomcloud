package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenTelekomCloudIdentityAuthScopeV3DataSource_basic(t *testing.T) {
	const dsName = "data.opentelekomcloud_identity_auth_scope_v3.token"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudIdentityAuthScopeV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAuthScopeV3DataSourceID(dsName),
					resource.TestCheckResourceAttrSet(dsName, "user_id"),
					checkAttrEqualsEnvOrSet(dsName, "user_name", "OS_USERNAME"),
					checkAttrEqualsEnvOrSet(dsName, "project_name", "OS_PROJECT_NAME"),
				),
			},
		},
	})
}

func checkAttrEqualsEnvOrSet(name, key, envVar string) resource.TestCheckFunc {
	if expected := os.Getenv(envVar); expected != "" {
		return resource.TestCheckResourceAttr(name, key, expected)
	}
	return resource.TestCheckResourceAttrSet(name, key)
}

func testAccCheckIdentityAuthScopeV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find token data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("token data source ID not set")
		}

		return nil
	}
}

const testAccOpenTelekomCloudIdentityAuthScopeV3DataSource_basic = `
data "opentelekomcloud_identity_auth_scope_v3" "token" {
  name = "my_token"
}
`
