package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const rpAttachResourceName = "opentelekomcloud_sdrs_replication_attach_v1.attach_1"

func TestAccSdrsReplicatonAttachV1_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsReplicationAttachV1Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rpAttachResourceName, "device", "/dev/vdb"),
				),
			},
		},
	})
}

func TestAccSDRSReplicationAttach_Import(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsReplicationAttachV1Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rpAttachResourceName, "device", "/dev/vdb"),
					resource.TestCheckResourceAttrSet(rpAttachResourceName, "status"),
					resource.TestCheckResourceAttrPair(rpAttachResourceName, "instance_id", "opentelekomcloud_sdrs_protected_instance_v1.instance_1", "id"),
					resource.TestCheckResourceAttrPair(rpAttachResourceName, "replication_id", "opentelekomcloud_sdrs_replication_pair_v1.pair_1", "id"),
				),
			},
			{
				ResourceName:      rpAttachResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccSdrsReplicationAttachV1Basic = fmt.Sprintf(`

%s
%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s3.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  availability_zone = "eu-de-02"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true
}

resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "eu-de-02"
  volume_type       = "SATA"
  size              = 12
}

resource "opentelekomcloud_sdrs_replication_pair_v1" "pair_1" {
  name                 = "replication_1"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  volume_id            = opentelekomcloud_evs_volume_v3.volume_1.id
  delete_target_volume = true
}

resource "opentelekomcloud_sdrs_replication_attach_v1" "attach_1" {
  instance_id    = opentelekomcloud_sdrs_protected_instance_v1.instance_1.id
  replication_id = opentelekomcloud_sdrs_replication_pair_v1.pair_1.id
  device         = "/dev/vdb"
}
`, common.DataSourceSubnet, common.DataSourceImage)
