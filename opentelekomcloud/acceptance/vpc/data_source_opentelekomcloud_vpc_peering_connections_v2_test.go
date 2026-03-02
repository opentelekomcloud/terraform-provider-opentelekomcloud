package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const (
	dataSourceVpcPeeringsName         = "data.opentelekomcloud_vpc_peering_connections_v2.by_vpc_id"
	dataSourceVpcPeeringsNameByName   = "data.opentelekomcloud_vpc_peering_connections_v2.by_name"
	dataSourceVpcPeeringsNameByStatus = "data.opentelekomcloud_vpc_peering_connections_v2.by_status"
	dataSourceVpcPeeringsNameEmpty    = "data.opentelekomcloud_vpc_peering_connections_v2.empty"
)

func TestAccVpcPeeringConnectionsV2DataSource_basic(t *testing.T) {
	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.Router, Count: 3},
	}
	quotas.BookMany(t, qts)

	dcByVpcId := common.InitDataSourceCheck(dataSourceVpcPeeringsName)
	dcByName := common.InitDataSourceCheck(dataSourceVpcPeeringsNameByName)
	dcByStatus := common.InitDataSourceCheck(dataSourceVpcPeeringsNameByStatus)
	dcEmpty := common.InitDataSourceCheck(dataSourceVpcPeeringsNameEmpty)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcPeeringConnectionsV2Config,
				Check: resource.ComposeTestCheckFunc(
					dcByVpcId.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringsName, "peering_connections.#", "2"),

					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringsNameByName, "peering_connections.#", "1"),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringsNameByName, "peering_connections.0.name", "opentelekomcloud_peerings_ds_1"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceVpcPeeringsNameByStatus, "peering_connections.0.status"),

					dcEmpty.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceVpcPeeringsNameEmpty, "peering_connections.#", "0"),
				),
			},
		},
	})
}

const testAccDataSourceVpcPeeringConnectionsV2Config = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_ds_peerings_1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_ds_peerings_2"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_3" {
  name = "vpc_test_ds_peerings_3"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud_peerings_ds_1"
  description = "test peering 1"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_2" {
  name        = "opentelekomcloud_peerings_ds_2"
  description = "test peering 2"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_3.id
}

data "opentelekomcloud_vpc_peering_connections_v2" "by_vpc_id" {
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id

  depends_on = [
    opentelekomcloud_vpc_peering_connection_v2.peering_1,
    opentelekomcloud_vpc_peering_connection_v2.peering_2,
  ]
}

data "opentelekomcloud_vpc_peering_connections_v2" "by_name" {
  name = opentelekomcloud_vpc_peering_connection_v2.peering_1.name
}

data "opentelekomcloud_vpc_peering_connections_v2" "by_status" {
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  status = "ACTIVE"

  depends_on = [
    opentelekomcloud_vpc_peering_connection_v2.peering_1,
    opentelekomcloud_vpc_peering_connection_v2.peering_2,
  ]
}

data "opentelekomcloud_vpc_peering_connections_v2" "empty" {
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_1.id

  depends_on = [
    opentelekomcloud_vpc_peering_connection_v2.peering_1,
    opentelekomcloud_vpc_peering_connection_v2.peering_2,
  ]
}
`
