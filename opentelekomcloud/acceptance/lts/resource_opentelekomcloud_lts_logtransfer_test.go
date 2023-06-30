package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/transfers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccLogTankTransferV2_basic(t *testing.T) {
	var transfer transfers.Transfer
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLogTankTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankTransferV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankTransferV2Exists(
						"opentelekomcloud_logtank_transfer_v2.transfer", &transfer),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "period", "12"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "switch_on", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "obs_bucket_name", "tf-test-bucket-lts"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "period_unit", "hour"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "storage_format", "RAW"),
				),
			},
			{
				Config: testAccLogTankTransferV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankTransferV2Exists(
						"opentelekomcloud_logtank_transfer_v2.transfer", &transfer),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "period", "30"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "switch_on", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "obs_bucket_name", "tf-test-bucket-lts"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "period_unit", "min"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "storage_format", "JSON"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "prefix_name", "prefix"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "dir_prefix_name", "dir"),
				),
			},
		},
	})
}

func TestAccLogTankTransferV2_encryptedBucket(t *testing.T) {
	var transfer transfers.Transfer
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLogTankTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankTransferV2_encrypted(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankTransferV2Exists(
						"opentelekomcloud_logtank_transfer_v2.transfer", &transfer),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "period", "30"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "switch_on", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "obs_bucket_name", "tf-test-bucket-lts-encrypted"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "period_unit", "min"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_transfer_v2.transfer", "storage_format", "JSON"),
				),
			},
		},
	})
}

func testAccCheckLogTankTransferV2Exists(n string, transfer *transfers.Transfer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		ltsclient, err := config.LtsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
		}

		allTransfers, err := transfers.ListTransfers(ltsclient, transfers.ListTransfersOpts{})
		if err != nil {
			return err
		}

		for _, transferRaw := range allTransfers {
			if transferRaw.LogTransferId == rs.Primary.ID {
				*transfer = transferRaw
				break
			}
		}

		return nil
	}
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
  dir_prefix_name = "dir"
}
`, env.OS_KMS_ID)
}
