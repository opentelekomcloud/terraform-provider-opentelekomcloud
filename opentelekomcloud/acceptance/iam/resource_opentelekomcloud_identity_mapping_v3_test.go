package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/mappings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccIdentityV3MappingBasic(t *testing.T) {
	resourceName := "opentelekomcloud_identity_mapping_v3.mapping"
	mappingID := tools.RandomString("mapping-", 3)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3MappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3MappingBasic(mappingID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mapping_id", mappingID),
				),
			},
			{
				Config: testAccIdentityV3MappingUpdate(mappingID),
			},
		},
	})
}

func testAccCheckIdentityV3MappingDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity v3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_mapping_v3" {
			continue
		}

		_, err := mappings.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("mapping still exists")
		}
	}

	return nil
}

func testAccIdentityV3MappingBasic(mappingID string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_mapping_v3" "mapping" {
  mapping_id = "%s"
  rules      = jsonencode([{ "local" : [{ "user" : { "name" : "{0}" } }, { "groups" : "[\"admin\",\"manager\"]" }], "remote" : [{ "type" : "uid" }] }])
}
`, mappingID)
}

func testAccIdentityV3MappingUpdate(mappingID string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_mapping_v3" "mapping" {
  mapping_id = "%s"
  rules      = jsonencode([{ "local" : [{ "user" : { "name" : "samltestid-{0}" } }], "remote" : [{ "type" : "uid" }] }])
}
`, mappingID)
}
