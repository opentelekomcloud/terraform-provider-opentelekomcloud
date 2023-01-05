package swr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/organizations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/swr"
)

func TestSwrOrganizationPermissionsV2_basic(t *testing.T) {
	userID := os.Getenv("OS_USER_ID_2")
	username := os.Getenv("OS_USERNAME_2")

	if userID == "" || username == "" {
		t.Skip("no OS_USER_ID_2 and OS_USERNAME_2 provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrOrganizationPermissionsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testSwrOrganizationPermissionV2BasicTemplate, name, userID, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePermissionsName, "auth", "3"),
				),
			},
		},
	})
}

func testSwrOrganizationPermissionsV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SwrV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(swr.ClientError, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_swr_organization_permissions_v2" {
			continue
		}

		org := rs.Primary.Attributes["organization"]
		perms, err := organizations.GetPermissions(client, org)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil // no permissions at all exist for the organization
			}
			return fmt.Errorf("error retrieving organization permissions: %w", err)
		}

		for _, a := range perms.OthersAuth {
			if a.UserID == rs.Primary.ID {
				return fmt.Errorf("expected permission to be deleted, but it exist")
			}
		}
	}

	return nil
}

const (
	resourcePermissionsName = "opentelekomcloud_swr_organization_permissions_v2.user_1"

	testSwrOrganizationPermissionV2BasicTemplate = `
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "%s"
}

resource opentelekomcloud_swr_organization_permissions_v2 user_1 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name

  user_id  = "%s"
  username = "%s"
  auth     = 3
}
`
)
