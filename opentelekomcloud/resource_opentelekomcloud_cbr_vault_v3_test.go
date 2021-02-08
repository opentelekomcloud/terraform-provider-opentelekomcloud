package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCBRVaultV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccAsConfigPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRVaultV3_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "1"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.0.name", "cbr-test-volume"),
				),
			},
			{
				Config: testCBRVaultV3_updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "0"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_cbr_vault_v3.vault", "backup_policy_id"),
				),
			},
			{
				Config: testCBRVaultV3_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "1"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.0.name", "cbr-test-volume"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "backup_policy_id", ""),
				),
			},
		},
	})
}

var (
	testCBRVaultV3_basic = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume" {
  name = "cbr-test-volume"
  size = 10

  volume_type = "SSD"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for terraform provider test"

  billing {
    size          = 100
    object_type   = "disk"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2

    extra_info = {}
  }

  resource {
    id   = opentelekomcloud_blockstorage_volume_v2.volume.id
    type = "OS::Cinder::Volume"

    extra_info = {}
  }
}
`
	testCBRVaultV3_updated = `
resource "opentelekomcloud_cbr_policy_v3" "policy" {
  name           = "some-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"
  ]
  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }

  enabled = "false"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for terraform provider test"

  backup_policy_id = opentelekomcloud_cbr_policy_v3.policy.id

  billing {
    size          = 100
    object_type   = "disk"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2

    extra_info = {}
  }
}
`
)
