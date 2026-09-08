package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cts/v3/tracker"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const trackerV3Resource = "opentelekomcloud_cts_tracker_v3.tracker_v3"

func TestAccCTSTrackerV3_basic(t *testing.T) {
	var ctsTracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV3Basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_support_validate", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_sort_by_service", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "disabled"),
					resource.TestCheckResourceAttr(trackerV3Resource, "compress_type", "json"),
				),
			},
			{
				Config: testAccCTSTrackerV3Update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q1"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_support_validate", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_sort_by_service", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "enabled"),
					resource.TestCheckResourceAttr(trackerV3Resource, "compress_type", "gzip"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV3_supportValidate(t *testing.T) {
	var ctsTracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV3SupportValidate(bucketName, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(trackerV3Resource, "compress_type", "gzip"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_support_validate", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_sort_by_service", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "enabled"),
				),
			},
			{
				// enabling `is_support_validate` alone must not reset the rest of
				// the tracker configuration, see #3541
				Config: testAccCTSTrackerV3SupportValidate(bucketName, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(trackerV3Resource, "compress_type", "gzip"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_support_validate", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_sort_by_service", "false"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "enabled"),
				),
			},
			{
				// the other half of the #3541 loop: changing an unrelated argument
				// must not reset `is_support_validate` back to `false`
				Config: testAccCTSTrackerV3SupportValidate(bucketName, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV3Exists(trackerV3Resource, &ctsTracker, env.OS_TENANT_NAME),
					resource.TestCheckResourceAttr(trackerV3Resource, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(trackerV3Resource, "file_prefix_name", "yO8Q"),
					resource.TestCheckResourceAttr(trackerV3Resource, "compress_type", "gzip"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_lts_enabled", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_support_validate", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "is_sort_by_service", "true"),
					resource.TestCheckResourceAttr(trackerV3Resource, "status", "enabled"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV3_importBasic(t *testing.T) {
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV3ImportBasic(bucketName),
			},

			{
				ResourceName:      trackerV3Resource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCTSTrackerV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	ctsClient, err := config.CtsV3Client(env.OS_TENANT_NAME)
	if err != nil {
		return fmt.Errorf("error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cts_tracker_v3" {
			continue
		}

		ctsTracker, err := tracker.List(ctsClient, "system")
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				continue
			} else {
				return fmt.Errorf("failed to retrieve CTS list: %s", err)
			}
		}

		if len(ctsTracker) != 0 {
			return fmt.Errorf("failed to delete CTS tracker")
		}
	}

	return nil
}

func testAccCheckCTSTrackerV3Exists(n string, trackers *tracker.Tracker, projectName cfg.ProjectName) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CtsV3Client(projectName)
		if err != nil {
			return fmt.Errorf("error creating cts client: %s", err)
		}

		ctsTracker, err := tracker.List(client, "system")
		if err != nil {
			return err
		}

		if len(ctsTracker) == 0 {
			return fmt.Errorf("CTS tracker not found")
		}

		if ctsTracker[0].TrackerName != rs.Primary.ID {
			return fmt.Errorf("CTS tracker not found")
		}

		trackers = &ctsTracker[0]

		return nil
	}
}

func testAccCTSTrackerV3Basic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%s"
  acl    = "public-read"
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name         = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name    = "yO8Q"
  status              = "disabled"
  compress_type       = "json"
  is_sort_by_service  = true
  is_support_validate = false
}
`, bucketName)
}

func testAccCTSTrackerV3Update(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name         = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name    = "yO8Q1"
  is_lts_enabled      = true
  status              = "enabled"
  compress_type       = "gzip"
  is_sort_by_service  = false
  is_support_validate = true
}
`, bucketName)
}

func testAccCTSTrackerV3SupportValidate(bucketName string, supportValidate, sortByService bool) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name         = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name    = "yO8Q"
  compress_type       = "gzip"
  is_lts_enabled      = true
  is_obs_created      = false
  is_sort_by_service  = %t
  is_support_validate = %t
  status              = "enabled"
}
`, bucketName, sortByService, supportValidate)
}

func testAccCTSTrackerV3ImportBasic(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q1"
  is_lts_enabled   = false
  status           = "enabled"
}
`, bucketName)
}
