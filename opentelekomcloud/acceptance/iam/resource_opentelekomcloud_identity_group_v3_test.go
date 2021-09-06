package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccIdentityV3Group_basic(t *testing.T) {
	var group groups.Group
	var groupName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Group_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists("opentelekomcloud_identity_group_v3.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_group_v3.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_group_v3.group_1", "description", &group.Description),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_group_v3.group_1", "domain_id", &group.DomainID),
				),
			},
			{
				Config: testAccIdentityV3Group_update(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists("opentelekomcloud_identity_group_v3.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_group_v3.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_group_v3.group_1", "description", &group.Description),
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_identity_group_v3.group_1", "domain_id", &group.DomainID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_group_v3" {
			continue
		}

		_, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("group still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
		}

		found, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("group not found")
		}

		*group = *found

		return nil
	}
}

func testAccIdentityV3Group_basic(groupName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_group_v3" "group_1" {
  name        = "%s"
  description = "A ACC test group"
}
  `, groupName)
}

func testAccIdentityV3Group_update(groupName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_group_v3" "group_1" {
  name        = "%s"
  description = "Some Group"
}
  `, groupName)
}
