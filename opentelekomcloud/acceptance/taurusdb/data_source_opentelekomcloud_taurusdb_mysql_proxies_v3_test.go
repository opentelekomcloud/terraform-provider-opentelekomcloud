package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceName = "data.opentelekomcloud_taurusdb_mysql_proxies_v3.test"

func TestAccDataSourceTaurusDBMysqlProxies_basic(t *testing.T) {
	name := "tf_taurusdb_proxy_" + acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBMysqlProxiesBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.delay_threshold_in_seconds"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.node_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.ram"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.elb_vip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.transaction_split"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.master_node_weight.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.master_node_weight.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.master_node_weight.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.master_node_weight.0.weight"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.readonly_nodes_weight.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.readonly_nodes_weight.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.readonly_nodes_weight.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.readonly_nodes_weight.0.weight"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.0.role"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_list.0.nodes.0.frozen_flag"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBMysqlProxiesBase(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_taurusdb_mysql_instance_v3" "instance" {
  name                     = "%s"
  password                 = "Test@12345678"
  flavor                   = "gaussdb.mysql.xlarge.x86.8"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id        = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 2
}

locals {
  sorted_nodes = [
    for node in opentelekomcloud_taurusdb_mysql_instance_v3.instance.nodes :
    node if node.type == "master"
  ]

  readonly_nodes = [
    for node in opentelekomcloud_taurusdb_mysql_instance_v3.instance.nodes :
    node if node.type == "slave"
  ]
}

resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
  flavor      = "gaussdb.proxy.large.x86.2"
  node_num    = 2
  proxy_name  = "%[3]s"
  proxy_mode  = "readwrite"

  master_node_weight {
    id     = local.sorted_nodes[0].id
    weight = 50
  }

  readonly_nodes_weight {
    id     = local.readonly_nodes[0].id
    weight = 50
  }
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}

func testAccDataSourceTaurusDBMysqlProxiesBasic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_taurusdb_mysql_proxies_v3" "test" {
  depends_on = [opentelekomcloud_taurusdb_mysql_proxy_v3.test]

  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
}
`, testAccDataSourceTaurusDBMysqlProxiesBase(name))
}
