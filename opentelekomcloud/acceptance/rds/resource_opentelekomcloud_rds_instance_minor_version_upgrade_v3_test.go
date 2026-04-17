package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccRdsUpgradingMinorVersion_basic(t *testing.T) {
	rdsId := os.Getenv("RDS_INSTANCE_ID")
	if rdsId == "" {
		t.Skip("rds instance id is required for this test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsUpgradingMinorVersion_basic(rdsId),
			},
		},
	})
}

func testAccRdsUpgradingMinorVersion_basic(instanceId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_rds_instance_minor_version_upgrade_v3" "test" {
  instance_id = "%s"
}
`, instanceId)
}
