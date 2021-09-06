package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccSFSTurboShareV1DataSource_basic(t *testing.T) {
	name := tools.RandomString("turbo-", 5)
	dsName := "data.opentelekomcloud_sfs_turbo_share_v1.share"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccAzPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboShareV1Basic(name),
			},
			{
				Config: testAccSFSTurboShareV1DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(dsName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(dsName, "size", "500"),
				),
			},
		},
	})
}

func testAccSFSTurboShareV1DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_sfs_turbo_share_v1" "share" {
  name = "%s"
}
`, testAccSFSTurboShareV1Basic(name), name)
}
