package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func TestAccCTSTrackerV1_basic(t *testing.T) {
	var tracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("opentelekomcloud_cts_tracker_v1.tracker_v1", &tracker),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "bucket_name", bucketName),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "file_prefix_name", "yO8Q"),
				),
			},
			{
				Config: testAccCTSTrackerV1_update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("opentelekomcloud_cts_tracker_v1.tracker_v1", &tracker),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cts_tracker_v1.tracker_v1", "file_prefix_name", "yO8Q1"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_timeout(t *testing.T) {
	var tracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1_timeout(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("opentelekomcloud_cts_tracker_v1.tracker_v1", &tracker),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ctsClient, err := config.ctsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cts_tracker_v1" {
			continue
		}

		_, err := tracker.List(ctsClient, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("cts tracker still exists.")
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}

	return nil
}

func testAccCheckCTSTrackerV1Exists(n string, trackers *tracker.Tracker) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ctsClient, err := config.ctsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating cts client: %s", err)
		}

		trackerList, err := tracker.List(ctsClient, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := trackerList[0]
		if found.TrackerName != rs.Primary.ID {
			return fmt.Errorf("cts tracker not found")
		}

		*trackers = found

		return nil
	}
}

func testAccCTSTrackerV1_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "%s"
  acl = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  = "topic_check"
  display_name    = "The display name of topic_check"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"
  is_support_smn = true
  topic_id = "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations = ["login"]
  need_notify_user_list = ["user1"]
}
`, bucketName)
}

func testAccCTSTrackerV1_update(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "%s"
  acl = "public-read"
  force_destroy = true
}
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  = "topic_check1"
  display_name    = "The display name of topic_check"
}
resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q1"
  is_support_smn = true
  topic_id = "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations = ["login"]
  need_notify_user_list = ["user1"]
}
`, bucketName)
}

func testAccCTSTrackerV1_timeout(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "%s"
  acl = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name		  = "topic_check-1"
  display_name    = "The display name of topic_check"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = "${opentelekomcloud_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"
  is_support_smn = true
  topic_id = "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations = ["login"]
  need_notify_user_list = ["user1"]

timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, bucketName)
}
