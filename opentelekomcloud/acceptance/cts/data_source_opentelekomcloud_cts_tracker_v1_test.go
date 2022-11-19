package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCTSTrackerV1DataSource_basic(t *testing.T) {
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1DataSource_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1DataSourceID("data.opentelekomcloud_cts_tracker_v1.d_tracker"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cts_tracker_v1.d_tracker", "bucket_name", bucketName),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cts_tracker_v1.d_tracker", "status", "enabled"),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find cts tracker data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("tracker data source not set ")
		}

		return nil
	}
}

func testAccCTSTrackerV1DataSource_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q"
  is_lts_enabled   = true
}

data "opentelekomcloud_cts_tracker_v1" "d_tracker" {
  tracker_name = opentelekomcloud_cts_tracker_v1.tracker_v1.tracker_name
}
`, bucketName)
}
