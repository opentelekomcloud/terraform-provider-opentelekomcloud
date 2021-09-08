package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceVaultName = "opentelekomcloud_cbr_vault_v3.vault"

func TestAccCBRVaultV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicVolumes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "2"),
					resource.TestCheckResourceAttr(resourceVaultName, "resource.0.name", "cbr-test-volume"),
					resource.TestCheckResourceAttrSet(resourceVaultName, "backup_policy_id"),
				),
			},
			{
				Config: testAccCBRVaultV3NoResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "0"),
					resource.TestCheckNoResourceAttr(resourceVaultName, "backup_policy_id"),
				),
			},
			{
				Config: testAccCBRVaultV3NoResourceResize,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.size", "120"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_unAssign(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicVolumes,
			},
			{
				Config: testAccCBRVaultV3Unassign,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "backup_policy_id", ""),
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "0"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_instance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicInstance,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "1"),
					resource.TestCheckResourceAttr(resourceVaultName, "resource.0.name", "tf-crb-test-instance"),
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.size", "100"),
				),
			},
			{
				Config: testAccCBRVaultV3NoResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "0"),
				),
			},
		},
	})
}

var (
	testAccCBRVaultV3BasicInstance = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_compute_instance_v2" "instance" {
  name = "tf-crb-test-instance"

  image_id    = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor_name = "%s"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OsFlavorID)
)

const (
	testAccCBRVaultV3BasicVolumes = `
resource "opentelekomcloud_blockstorage_volume_v2" "volume" {
  name = "cbr-test-volume"
  size = 10

  volume_type = "SSD"
}
resource "opentelekomcloud_blockstorage_volume_v2" "volume2" {
  name = "cbr-test-volume-2"
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

  resource {
    id   = opentelekomcloud_blockstorage_volume_v2.volume2.id
    type = "OS::Cinder::Volume"
  }
}
`

	testAccCBRVaultV3Unassign = `
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

	testAccCBRVaultV3NoResource = `
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
	testAccCBRVaultV3NoResourceResize = `
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
