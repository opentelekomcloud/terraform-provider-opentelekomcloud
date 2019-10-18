package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCTSTrackerV1DataSource_basic(t *testing.T) {
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1DataSource_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1DataSourceID("data.opentelekomcloud_cts_tracker_v1.tracker_v1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cts_tracker_v1.tracker_v1", "bucket_name", bucketName),
					resource.TestCheckResourceAttr("data.opentelekomcloud_cts_tracker_v1.tracker_v1", "status", "enabled"),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cts tracker data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("tracker data source not set ")
		}

		return nil
	}
}

func testAccCTSTrackerV1DataSource_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket		= "%s"
  acl			= "public-read"
  force_destroy = true
}
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name			= "tf-test-topic"
  display_name	= "The display name of tf-test-topic"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name		= "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name  = "yO8Q"
  is_support_smn 	= true
  topic_id 			= "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations 		= ["login"]
  need_notify_user_list = ["user1"]
}

data "opentelekomcloud_cts_tracker_v1" "tracker_v1" {  
  tracker_name = "${opentelekomcloud_cts_tracker_v1.tracker_v1.id}"
}
`, bucketName)
}
