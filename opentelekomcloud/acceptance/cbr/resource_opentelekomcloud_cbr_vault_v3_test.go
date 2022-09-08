package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/vaults"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceVaultName = "opentelekomcloud_cbr_vault_v3.vault"

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
    type = "SATA"
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
