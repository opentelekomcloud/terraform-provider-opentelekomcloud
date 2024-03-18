package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/replications"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const rpResourceName = "opentelekomcloud_sdrs_replication_pair_v1.pair_1"

func TestAccSdrsReplicatonPairV1_basic(t *testing.T) {
	var replication replications.Replication

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSdrsReplicatonPairV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsReplicationPairV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsReplicationPairV1Exists(rpResourceName, &replication),
					resource.TestCheckResourceAttr(rpResourceName, "name", "replication_1"),
					resource.TestCheckResourceAttr(rpResourceName, "description", "test description"),
				),
			},
			{
				Config: testAccSdrsReplicationPairV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsReplicationPairV1Exists(rpResourceName, &replication),
					resource.TestCheckResourceAttr(rpResourceName, "name", "replication_1_updated"),
				),
			},
		},
	})
}

func TestAccSdrsReplicationPairV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_sdrs_replication_pair_v1.pair_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccSdrsProtectedInstanceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsReplicationPairV1Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_target_volume",
				},
			},
		},
	})
}

func testAccCheckSdrsReplicatonPairV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SdrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sdrs_replication_pair_v1" {
			continue
		}

		_, err := replications.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS replication pair still exists")
		}
	}

	return nil
}

func testAccCheckSdrsReplicationPairV1Exists(n string, replication *replications.Replication) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.SdrsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
		}

		found, err := replications.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("SDRS replication pair not found")
		}

		*replication = *found

		return nil
	}
}

var testAccSdrsReplicationPairV1Base = fmt.Sprintf(`
%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  description              = "test description"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}

resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "eu-de-02"
  volume_type       = "SATA"
  size              = 12
}

`, common.DataSourceSubnet)

var testAccSdrsReplicationPairV1Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_sdrs_replication_pair_v1" "pair_1" {
  name                 = "replication_1"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  volume_id            = opentelekomcloud_evs_volume_v3.volume_1.id
  description          = "test description"
  delete_target_volume = true
}
`, testAccSdrsReplicationPairV1Base)

var testAccSdrsReplicationPairV1Update = fmt.Sprintf(`
%s

resource "opentelekomcloud_sdrs_replication_pair_v1" "pair_1" {
  name                 = "replication_1_updated"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  volume_id            = opentelekomcloud_evs_volume_v3.volume_1.id
  description          = "test description"
  delete_target_volume = true
}
`, testAccSdrsReplicationPairV1Base)
