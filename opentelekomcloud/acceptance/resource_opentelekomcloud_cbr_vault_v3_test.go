package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCBRVaultV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRVaultV3_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "1"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.0.name", "cbr-test-volume"),
					resource.TestCheckResourceAttrSet("opentelekomcloud_cbr_vault_v3.vault", "backup_policy_id"),
				),
			},
			{
				Config: testCBRVaultV3_noResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "0"),
					resource.TestCheckNoResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "backup_policy_id"),
				),
			},
			{
				Config: testCBRVaultV3_noResourceResize,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "billing.0.size", "120"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_unassign(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRVaultV3_basic,
			},
			{
				Config: testCBRVaultV3_unassign,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "backup_policy_id", ""),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "0"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_instance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccFlavorPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRVaultV3_basicInstance,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "1"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.0.name", "tf-crb-test-instance"),
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "billing.0.size", "100"),
				),
			},
			{
				Config: testCBRVaultV3_noResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_cbr_vault_v3.vault", "resource.#", "0"),
				),
			},
		},
	})
}

var (
	testCBRVaultV3_basicInstance = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "instance" {
  name = "tf-crb-test-instance"

  image_id    = "%s"
  flavor_name = "%s"

  network {
    uuid = "%s"
  }
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for terraform provider test"

  billing {
    size          = 100
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }

  resource {
    id   = opentelekomcloud_compute_instance_v2.instance.id
    type = "OS::Nova::Server"
  }
}
`, OS_IMAGE_ID, OS_FLAVOR_NAME, OS_NETWORK_ID)
)

const (
	testCBRVaultV3_basic = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume" {
  name = "cbr-test-volume"
  size = 10

  volume_type = "SSD"
}

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
  }

  resource {
    id   = opentelekomcloud_blockstorage_volume_v2.volume.id
    type = "OS::Cinder::Volume"
  }
}
`

	testCBRVaultV3_unassign = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume" {
  name = "cbr-test-volume"
  size = 10

  volume_type = "SSD"
}

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

  billing {
    size          = 100
    object_type   = "disk"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }

  resource = []
}
`

	testCBRVaultV3_noResource = `
resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for terraform provider test"

  billing {
    size          = 100
    object_type   = "disk"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }
}
`
	testCBRVaultV3_noResourceResize = `
resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test-2"

  description = "CBR vault for terraform provider test"

  billing {
    size          = 120
    object_type   = "disk"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }
}
`
)
