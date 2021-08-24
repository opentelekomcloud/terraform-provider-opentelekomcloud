package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v1/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const trackerResource = "opentelekomcloud_cts_tracker_v1.tracker_v1"

func TestAccCTSTrackerV1_basic(t *testing.T) {
	var ctsTracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1Basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerResource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerResource, "file_prefix_name", "yO8Q"),
				),
			},
			{
				Config: testAccCTSTrackerV1Update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerResource, "file_prefix_name", "yO8Q1"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_timeout(t *testing.T) {
	var track tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1Timeout(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &track, env.OS_TENANT_NAME),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_KeyOperations(t *testing.T) {
	var track tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1AllOperations(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &track, env.OS_TENANT_NAME),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_schemaProjectName(t *testing.T) {
	var ctsTracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	var projectName2 = cfg.ProjectName(os.Getenv("OS_PROJECT_NAME_2"))
	if projectName2 == "" {
		t.Skip("OS_PROJECT_NAME_2 is empty")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1ProjectName(bucketName, projectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(trackerResource, &ctsTracker, projectName2),
					resource.TestCheckResourceAttr(trackerResource, "project_name", string(projectName2)),
				),
			},
		},
	})
	env.OS_TENANT_NAME = env.GetTenantName()
}

func testAccCheckCTSTrackerV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	ctsClient, err := config.CtsV1Client(env.OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cts_tracker_v1" {
			continue
		}

		list, err := tracker.List(ctsClient, tracker.ListOpts{
			TrackerName:    rs.Primary.ID,
			BucketName:     rs.Primary.Attributes["bucket_name"],
			FilePrefixName: rs.Primary.Attributes["file_prefix_name"],
			Status:         rs.Primary.Attributes["status"],
		})
		if err != nil {
			return fmt.Errorf("failed to retrieve CTS list: %s", err)
		}
		if len(list) != 0 {
			return fmt.Errorf("failed to delete CTS tracker")
		}
	}

	return nil
}

func testAccCheckCTSTrackerV1Exists(n string, trackers *tracker.Tracker, projectName cfg.ProjectName) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CtsV1Client(projectName)
		if err != nil {
			return fmt.Errorf("error creating cts client: %s", err)
		}

		trackerList, err := tracker.List(client, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := trackerList[0]
		if found.TrackerName != rs.Primary.ID {
			return fmt.Errorf("CTS tracker not found")
		}

		*trackers = found

		return nil
	}
}

func testAccCTSTrackerV1Basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_check"
  display_name = "The display name of topic_check"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name               = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name          = "yO8Q"
  is_support_smn            = true
  topic_id                  = opentelekomcloud_smn_topic_v2.topic_1.id
  is_send_all_key_operation = false
  operations                = ["login"]
  need_notify_user_list     = ["user1"]
}
`, bucketName)
}

func testAccCTSTrackerV1Update(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_check1"
  display_name = "The display name of topic_check"
}
resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name               = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name          = "yO8Q1"
  is_support_smn            = true
  topic_id                  = opentelekomcloud_smn_topic_v2.topic_1.id
  is_send_all_key_operation = false
  operations                = ["login"]
  need_notify_user_list     = ["user1"]
}
`, bucketName)
}

func testAccCTSTrackerV1Timeout(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_check-1"
  display_name = "The display name of topic_check"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name               = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name          = "yO8Q"
  is_support_smn            = true
  topic_id                  = opentelekomcloud_smn_topic_v2.topic_1.id
  is_send_all_key_operation = false
  operations                = ["login"]
  need_notify_user_list     = ["user1"]

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, bucketName)
}

func testAccCTSTrackerV1ProjectName(bucketName string, projectName cfg.ProjectName) string {
	return fmt.Sprintf(`
locals {
  project_name = "%s"
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q"
  project_name     = local.project_name
  is_support_smn   = false
}
`, projectName, bucketName)
}

func testAccCTSTrackerV1AllOperations(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_check"
  display_name = "The display name of topic_check"
}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name               = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name          = "yO8Q"
  is_support_smn            = true
  topic_id                  = opentelekomcloud_smn_topic_v2.topic_1.id
  is_send_all_key_operation = true
}
`, bucketName)
}
