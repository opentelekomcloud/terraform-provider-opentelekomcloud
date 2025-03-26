package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/transfers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getLtsTransferResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v2 client: %s", err)
	}
	requestResp, err := transfers.List(client, transfers.ListTransfersOpts{})
	if err != nil {
		return nil, err
	}
	if len(requestResp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var transferResult *transfers.Transfer
	for _, transfer := range requestResp {
		if transfer.LogTransferId == state.Primary.ID {
			transferResult = &transfer
		}
	}
	if transferResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return transferResult, nil
}

func TestAccLogTankTransferV2_basic(t *testing.T) {
	var (
		transfer     transfers.Transfer
		resourceName = "opentelekomcloud_logtank_transfer_v2.transfer"
		rc           = common.InitResourceCheck(resourceName, &transfer, getLtsTransferResourceFunc)
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankTransferV2_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "period", "12"),
					resource.TestCheckResourceAttr(
						resourceName, "switch_on", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "obs_bucket_name", "tf-test-bucket-lts"),
					resource.TestCheckResourceAttr(
						resourceName, "period_unit", "hour"),
					resource.TestCheckResourceAttr(
						resourceName, "storage_format", "RAW"),
				),
			},
			{
				Config: testAccLogTankTransferV2_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "period", "30"),
					resource.TestCheckResourceAttr(
						resourceName, "switch_on", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "obs_bucket_name", "tf-test-bucket-lts"),
					resource.TestCheckResourceAttr(
						resourceName, "period_unit", "min"),
					resource.TestCheckResourceAttr(
						resourceName, "storage_format", "JSON"),
					resource.TestCheckResourceAttr(
						resourceName, "prefix_name", "prefix"),
					resource.TestCheckResourceAttr(
						resourceName, "dir_prefix_name", "dir"),
				),
			},
		},
	})
}

func TestAccLogTankTransferV2_encryptedBucket(t *testing.T) {
	var (
		transfer     transfers.Transfer
		resourceName = "opentelekomcloud_logtank_transfer_v2.transfer"
		rc           = common.InitResourceCheck(resourceName, &transfer, getLtsTransferResourceFunc)
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankTransferV2_encrypted(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "period", "30"),
					resource.TestCheckResourceAttr(
						resourceName, "switch_on", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "obs_bucket_name", "tf-test-bucket-lts-encrypted"),
					resource.TestCheckResourceAttr(
						resourceName, "period_unit", "min"),
					resource.TestCheckResourceAttr(
						resourceName, "storage_format", "JSON"),
				),
			},
		},
	})
}

const testAccLogTankTransferV2_basic = `
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-lts"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
  group_name  = "testacc_group"
  ttl_in_days = 7
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic"
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic-2" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic-2"
}

resource "opentelekomcloud_logtank_transfer_v2" "transfer" {
  log_group_id    = opentelekomcloud_logtank_group_v2.testacc_group.id
  log_stream_ids  = [opentelekomcloud_logtank_topic_v2.testacc_topic.id, opentelekomcloud_logtank_topic_v2.testacc_topic-2.id]
  obs_bucket_name = opentelekomcloud_obs_bucket.bucket.bucket
  storage_format  = "RAW"
  period          = 12
  period_unit     = "hour"
}
`

const testAccLogTankTransferV2_update = `
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-lts"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
  group_name  = "testacc_group"
  ttl_in_days = 7
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic"
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic-2" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic-2"
}

resource "opentelekomcloud_logtank_transfer_v2" "transfer" {
  log_group_id    = opentelekomcloud_logtank_group_v2.testacc_group.id
  log_stream_ids  = [opentelekomcloud_logtank_topic_v2.testacc_topic.id, opentelekomcloud_logtank_topic_v2.testacc_topic-2.id]
  obs_bucket_name = opentelekomcloud_obs_bucket.bucket.bucket
  storage_format  = "JSON"
  switch_on       = false
  period          = 30
  period_unit     = "min"
  prefix_name     = "prefix"
  dir_prefix_name = "dir"
}
`

func testAccLogTankTransferV2_encrypted() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-lts-encrypted"
  storage_class = "STANDARD"
  acl           = "private"
  server_side_encryption {
    algorithm  = "kms"
    kms_key_id = "%s"
  }
}

resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
  group_name  = "testacc_group-encr"
  ttl_in_days = 7
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic"
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic-2" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic-2"
}

resource "opentelekomcloud_logtank_transfer_v2" "transfer" {
  log_group_id    = opentelekomcloud_logtank_group_v2.testacc_group.id
  log_stream_ids  = [opentelekomcloud_logtank_topic_v2.testacc_topic.id, opentelekomcloud_logtank_topic_v2.testacc_topic-2.id]
  obs_bucket_name = opentelekomcloud_obs_bucket.bucket.bucket
  storage_format  = "JSON"
  switch_on       = false
  period          = 30
  period_unit     = "min"
  prefix_name     = "prefix"
  // dir_prefix_name = "dir"
}
`, env.OS_KMS_ID)
}
