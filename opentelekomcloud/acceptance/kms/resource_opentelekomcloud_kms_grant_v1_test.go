package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/kms/v1/grants"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/kms"
)

func TestAccKmsGrantV1Basic(t *testing.T) {
	var grant grants.Grant
	resourceName := "opentelekomcloud_kms_grant_v1.grant_1"

	granteePrincipal := os.Getenv("OS_USER_ID")
	if granteePrincipal == "" {
		t.Skip("OS_USER_ID must be set for acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckKmsV1GrantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsGrantV1Basic(granteePrincipal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsV1GrantExists(resourceName, &grant),
					resource.TestCheckResourceAttr(resourceName, "key_id", env.OsKmsID),
					resource.TestCheckResourceAttr(resourceName, "name", "my_grant"),
				),
			},
		},
	})
}

func testAccCheckKmsV1GrantDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.KmsKeyV1Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud KMSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_kms_grant_v1" {
			continue
		}
		kmsID, grantID, err := kms.ResourceKMSGrantV1ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}
		listOpts := grants.ListOpts{
			KeyID: kmsID,
		}

		grantsList, err := grants.List(client, listOpts).Extract()
		if err != nil {
			return err
		}
		var found *grants.Grant
		for _, grant := range grantsList.Grants {
			if grant.GrantID == grantID {
				found = &grant
				break
			}
		}

		if found != nil {
			return fmt.Errorf("grant still exists")
		}
	}
	return nil
}

func testAccCheckKmsV1GrantExists(n string, grant *grants.Grant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.KmsKeyV1Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud KMSv1 client: %w", err)
		}

		kmsID, grantID, err := kms.ResourceKMSGrantV1ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}
		listOpts := grants.ListOpts{
			KeyID: kmsID,
		}

		grantsList, err := grants.List(client, listOpts).Extract()
		if err != nil {
			return err
		}
		var found *grants.Grant
		for _, grant := range grantsList.Grants {
			if grant.GrantID == grantID {
				found = &grant
				break
			}
		}

		if found == nil {
			return fmt.Errorf("grant not found")
		}

		*grant = *found
		return nil
	}
}

func testAccKmsGrantV1Basic(granteePrincipal string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_kms_grant_v1" "grant_1" {
  key_id            = "%s"
  name              = "my_grant"
  grantee_principal = "%s"
  operations        = ["describe-key", "create-datakey", "encrypt-datakey"]
}
`, env.OsKmsID, granteePrincipal)
}
