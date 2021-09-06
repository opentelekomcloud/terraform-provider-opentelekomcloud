package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenStackIdentityV3UserDataSource_basic(t *testing.T) {
	userName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	userPassword := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityUserV3DataSource_user(userName, userPassword),
			},
			{
				Config: testAccOpenStackIdentityUserV3DataSource_basic(userName, userPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityUserV3DataSourceID("data.opentelekomcloud_identity_user_v3.user_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_user_v3.user_1", "name", userName),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_user_v3.user_1", "domain_id", "default"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_identity_user_v3.user_1", "enabled", "true"),
					testAccCheckIdentityUserV3DataSourceDefaultProjectID(
						"data.opentelekomcloud_identity_user_v3.user_1", "opentelekomcloud_identity_project_v3.project_1"),
				),
			},
		},
	})
}

func testAccCheckIdentityUserV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find user data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("user data source ID not set")
		}

		return nil
	}
}

func testAccCheckIdentityUserV3DataSourceDefaultProjectID(n1, n2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds1, ok := s.RootModule().Resources[n1]
		if !ok {
			return fmt.Errorf("can't find user data source: %s", n1)
		}

		if ds1.Primary.ID == "" {
			return fmt.Errorf("user data source ID not set")
		}

		rs2, ok := s.RootModule().Resources[n2]
		if !ok {
			return fmt.Errorf("can't find project resource: %s", n2)
		}

		if rs2.Primary.ID == "" {
			return fmt.Errorf("project resource ID not set")
		}

		if rs2.Primary.ID != ds1.Primary.Attributes["default_project_id"] {
			return fmt.Errorf("project id and user default_project_id don't match")
		}

		return nil
	}
}

func testAccOpenStackIdentityUserV3DataSource_user(name, password string) string {
	return fmt.Sprintf(`
	%s

resource "opentelekomcloud_identity_user_v3" "user_1" {
  name               = "%s"
  password           = "%s"
  default_project_id = opentelekomcloud_identity_project_v3.project_1.id
}
`, testAccOpenStackIdentityProjectV3DataSource_project(fmt.Sprintf("%s_project", name), acctest.RandString(20)), name, password)
}

func testAccOpenStackIdentityUserV3DataSource_basic(name, password string) string {
	return fmt.Sprintf(`
	%s

data "opentelekomcloud_identity_user_v3" "user_1" {
  name = opentelekomcloud_identity_user_v3.user_1.name
}
`, testAccOpenStackIdentityUserV3DataSource_user(name, password))
}
