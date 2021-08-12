package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

var keyAlias = tools.RandomString("key_alias_", 3)

func TestAccKmsKeyV1DataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_kms_key_v1.key1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyV1DataSourceKey,
			},
			{
				Config: testAccKmsKeyV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyV1DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "key_alias", keyAlias),
				),
			},
		},
	})
}

func testAccCheckKmsKeyV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find KMS key data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kms key data source ID not set")
		}

		return nil
	}
}

var testAccKmsKeyV1DataSourceKey = fmt.Sprintf(`
resource "opentelekomcloud_kms_key_v1" "key1" {
  key_alias       = "%s"
  key_description = "test description"
  pending_days    = "7"
  is_enabled      = true
}`, keyAlias)

var testAccKmsKeyV1DataSourceBasic = fmt.Sprintf(`
%s
data "opentelekomcloud_kms_key_v1" "key1" {
  key_alias       = opentelekomcloud_kms_key_v1.key1.key_alias
  key_id          = opentelekomcloud_kms_key_v1.key1.id
  key_description = "test description"
  key_state       = "2"
}
`, testAccKmsKeyV1DataSourceKey)
