package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/vaults"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceVaultName = "opentelekomcloud_cbr_vault_v3.vault"

func getVaultResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CbrV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR client: %s", err)
	}
	return vaults.Get(client, state.Primary.ID)
}

func TestAccCBRVaultV3_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicVolumes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "2"),
					resource.TestCheckResourceAttrSet(resourceVaultName, "backup_policy_id"),
					resource.TestCheckResourceAttr(resourceVaultName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCBRVaultV3Tags,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "tags.new_test", "new_test2"),
					resource.TestCheckResourceAttr(resourceVaultName, "tags.john", "doe"),
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
			{
				ResourceName:            resourceVaultName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing.0.period_type"},
			},
		},
	})
}

func TestAccCBRVaultV3_unAssign(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicVolumes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "2"),
				),
			},
			{
				Config: testAccCBRVaultV3BasicSingleVolume,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "1"),
				),
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
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
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

func TestAccCBRVaultV3_SfsTurbo(t *testing.T) {
	shareName := tools.RandomString("sfs-turbo-", 3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultv3SFSTurboShare(shareName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "1"),
					resource.TestCheckResourceAttr(resourceVaultName, "resource.0.name", shareName),
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.size", "1000"),
				),
			},
			{
				Config: testAccCBRVaultV3TurboNoResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "0"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_extraInfoExclude(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicExtraInfo,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceVaultName, "backup_policy_id"),
				),
			},
			{
				Config: testAccCBRVaultV3BasicExtraInfoUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "1"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_extraInfoInclude_OnlySwissCloud(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BasicIncludeVolumes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "1"),
				),
			},
			{
				Config: testAccCBRVaultV3BasicIncludeVolumesUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "resource.#", "1"),
				),
			},
		},
	})
}

func TestAccCBRVaultV3_bind_rules(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRVaultV3BindRules,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "auto_bind", "true"),
					resource.TestCheckResourceAttr(resourceVaultName, "bind_rules.#", "1"),
				),
			},
		},
	})
}

func TestAccVault_locked(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.Volume, Count: 2},
				{Q: quotas.VolumeSize, Count: 20},
				{Q: quotas.CBRPolicy, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVault_locked_step("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.object_type", "server"),
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.protect_type", "backup"),
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.size", "100"),
					resource.TestCheckResourceAttr(resourceVaultName, "locked", "true"),
				),
			},
			{
				Config: testAccVault_locked_step("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.object_type", "server"),
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.protect_type", "backup"),
					resource.TestCheckResourceAttr(resourceVaultName, "billing.0.size", "100"),
					resource.TestCheckResourceAttr(resourceVaultName, "locked", "false"),
				),
				ExpectError: regexp.MustCompile("vault not support to modify locked attribute from true to false."),
			},
		},
	})
}

func testAccCheckCBRVaultV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CbrV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cbr_vault_v3" {
			continue
		}

		_, err := vaults.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CBRv3 vault still exists")
		}
	}

	return nil
}

func TestAccVault_bindPolicies(t *testing.T) {
	if os.Getenv("OS_DEST_VAULT_ID") == "" {
		t.Skip("OS_DEST_VAULT_ID is not set; skipping OpenTelekomCloud CBR bind policies test.")
	}
	var (
		vault interface{}

		randName     = tools.RandomString("cbr-pol-", 3)
		mainRcName   = "opentelekomcloud_cbr_vault_v3.test"
		legacyRcName = "opentelekomcloud_cbr_vault_v3.legacy"

		mainRc   = common.InitResourceCheck(mainRcName, &vault, getVaultResourceFunc)
		legacyRc = common.InitResourceCheck(legacyRcName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckReplication(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      mainRc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_bindPolicies_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "backup_policy_id",
						"opentelekomcloud_cbr_policy_v3.backup.0", "id"),
				),
			},
			{
				Config: testAccVault_bindPolicies_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "backup_policy_id",
						"opentelekomcloud_cbr_policy_v3.backup.1", "id"),
				),
			},
			{
				ResourceName:      mainRcName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"billing",
				},
			},
			{
				ResourceName:      legacyRcName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"billing",
					"backup_policy_id",
				},
			},
		},
	})
}

func testAccVault_bindPolicies_base(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cbr_policy_v3" "backup" {
  count = 2

  name            = format("%[1]s_%%d", count.index)
  operation_type  = "backup"

  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }

  trigger_pattern = [
    "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"
  ]
}

resource "opentelekomcloud_cbr_policy_v3" "replication" {
  count = 2

  name                   = format("%[1]s_%%d", count.index)
  operation_type         = "replication"
  destination_region     = "%[2]s"
  destination_project_id = "%[3]s"

  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }

  trigger_pattern = [
    "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"
  ]
}
`, name, env.OS_DEST_REGION, env.OS_DEST_PROJECT_ID)
}

func testAccVault_bindPolicies_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cbr_vault_v3" "legacy" {
  name             = "%[2]s_legacy"
  backup_policy_id = opentelekomcloud_cbr_policy_v3.backup[0].id

  billing {
    size          = 100
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }
}

resource "opentelekomcloud_cbr_vault_v3" "test" {
  name            = "%[2]s"

  billing {
    size          = 100
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }

  policy {
    id = opentelekomcloud_cbr_policy_v3.backup[0].id
  }
  policy {
    id                   = opentelekomcloud_cbr_policy_v3.replication[0].id
    destination_vault_id = "%[3]s"
  }
}
`, testAccVault_bindPolicies_base(name), name, env.OS_DEST_VAULT_ID)
}

func testAccVault_bindPolicies_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cbr_vault_v3" "legacy" {
  name             = "%[2]s_legacy"
  backup_policy_id = opentelekomcloud_cbr_policy_v3.backup[1].id

  billing {
    size          = 100
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }
}

resource "opentelekomcloud_cbr_vault_v3" "test" {
  name            = "%[2]s"

  billing {
    size          = 100
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }

  policy {
    id = opentelekomcloud_cbr_policy_v3.backup[1].id
  }
  policy {
    id                   = opentelekomcloud_cbr_policy_v3.replication[1].id
    destination_vault_id = "%[3]s"
  }
}
`, testAccVault_bindPolicies_base(name), name, env.OS_DEST_VAULT_ID)
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

  tags = {
    foo = "bar"
    key = "value"
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

	testAccCBRVaultV3BasicSingleVolume = `
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

	testAccCBRVaultV3Tags = `
resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  billing {
    size          = 100
    object_type   = "disk"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }

  tags = {
    foo      = "bar"
    john     = "doe"
    new_test = "new_test2"
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
	testAccCBRVaultV3BasicExtraInfo = `
resource "opentelekomcloud_cbr_policy_v3" "default_policy" {
  name           = "cbr-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=DAILY;INTERVAL=1;BYHOUR=23;BYMINUTE=00"
  ]
  operation_definition {
    max_backups = 5
    timezone    = "UTC+01:00"
  }

  enabled = "true"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for default backup policy"

  backup_policy_id = opentelekomcloud_cbr_policy_v3.default_policy.id

  auto_bind   = true
  auto_expand = true

  billing {
    size          = 10000
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }

}
`
)

var testAccCBRVaultV3BasicExtraInfoUpdate = fmt.Sprintf(`
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "%s"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "c0b36460-7aa6-44d2-990d-cc300f3a7e43"
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  data_disks {
    type = "SAS"
    size = "10"
  }
  data_disks {
    type = "SAS"
    size = "10"
  }

  password                    = "Password@123"
  availability_zone           = "%s"
  auto_recovery               = true
  delete_disks_on_termination = true
}

resource "opentelekomcloud_cbr_policy_v3" "default_policy" {
  name           = "cbr-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=DAILY;INTERVAL=1;BYHOUR=23;BYMINUTE=00"
  ]
  operation_definition {
    max_backups = 5
    timezone    = "UTC+01:00"
  }

  enabled = "true"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for default backup policy"

  billing {
    size          = 10000
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }

  resource {
    id   = opentelekomcloud_ecs_instance_v1.instance_1.id
    type = "OS::Nova::Server"

    exclude_volumes = [
      opentelekomcloud_ecs_instance_v1.instance_1.volumes_attached.1.id
    ]

  }
}
`, env.OsSubnetName, env.OS_AVAILABILITY_ZONE)

var testAccCBRVaultV3BasicIncludeVolumes = fmt.Sprintf(`
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "%s"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "fc11b59c-46e3-4b3f-84be-dd6bf9aef1b8"
  flavor   = "s3.xlarge.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  system_disk_type = "SSD"
  data_disks {
    type = "SSD"
    size = "10"
  }
  data_disks {
    type = "SSD"
    size = "10"
  }

  password                    = "Password@123"
  availability_zone           = "%s"
  auto_recovery               = true
  delete_disks_on_termination = true
}

resource "opentelekomcloud_cbr_policy_v3" "default_policy" {
  name           = "cbr-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=DAILY;INTERVAL=1;BYHOUR=23;BYMINUTE=00"
  ]
  operation_definition {
    max_backups = 5
    timezone    = "UTC+01:00"
  }

  enabled = "true"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for default backup policy"

  billing {
    size          = 10000
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }

  resource {
    id   = opentelekomcloud_ecs_instance_v1.instance_1.id
    type = "OS::Nova::Server"

    include_volumes = [
      opentelekomcloud_ecs_instance_v1.instance_1.volumes_attached.0.id
    ]
  }
}
`, env.OsSubnetName, env.OS_AVAILABILITY_ZONE)

var testAccCBRVaultV3BasicIncludeVolumesUpdate = fmt.Sprintf(`
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "%s"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "fc11b59c-46e3-4b3f-84be-dd6bf9aef1b8"
  flavor   = "s3.xlarge.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  system_disk_type = "SSD"
  data_disks {
    type = "SSD"
    size = "10"
  }
  data_disks {
    type = "SSD"
    size = "10"
  }

  password                    = "Password@123"
  availability_zone           = "%s"
  auto_recovery               = true
  delete_disks_on_termination = true
}

resource "opentelekomcloud_cbr_policy_v3" "default_policy" {
  name           = "cbr-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=DAILY;INTERVAL=1;BYHOUR=23;BYMINUTE=00"
  ]
  operation_definition {
    max_backups = 5
    timezone    = "UTC+01:00"
  }

  enabled = "true"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for default backup policy"

  billing {
    size          = 10000
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }

  resource {
    id   = opentelekomcloud_ecs_instance_v1.instance_1.id
    type = "OS::Nova::Server"

    exclude_volumes = [
      opentelekomcloud_ecs_instance_v1.instance_1.volumes_attached.0.id
    ]

  }
}
`, env.OsSubnetName, env.OS_AVAILABILITY_ZONE)

const testAccCBRVaultV3BindRules = `
resource "opentelekomcloud_cbr_policy_v3" "default_policy" {
  name           = "cbr-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=DAILY;INTERVAL=1;BYHOUR=23;BYMINUTE=00"
  ]
  operation_definition {
    max_backups = 5
    timezone    = "UTC+01:00"
  }

  enabled = "true"
}

resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for default backup policy"

  billing {
    size          = 10
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }
  auto_bind = true

  bind_rules {
    key   = "foo"
    value = "bar"
  }
}
`

func testAccCBRVaultv3SFSTurboShare(shareName string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name        = "%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone = "%s"
}
resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for terraform provider test"

  billing {
    size          = 1000
    object_type   = "turbo"
    protect_type  = "backup"
    charging_mode = "post_paid"
    period_type   = "month"
    period_num    = 2
  }

  resource {
    id   = opentelekomcloud_sfs_turbo_share_v1.sfs-turbo.id
    type = "OS::Sfs::Turbo"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, shareName, env.OS_AVAILABILITY_ZONE)
}

var testAccCBRVaultV3TurboNoResource = `
resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name = "cbr-vault-test"

  description = "CBR vault for terraform provider test"

  billing {
    size          = 1000
    object_type   = "turbo"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }
}
`

func testAccVault_locked_step(locked string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cbr_vault_v3" "vault" {
  name   = "cbr-vault-test"
  locked = %s

  billing {
    size          = 100
    object_type   = "server"
    protect_type  = "backup"
    charging_mode = "post_paid"
  }
}
`, locked)
}
