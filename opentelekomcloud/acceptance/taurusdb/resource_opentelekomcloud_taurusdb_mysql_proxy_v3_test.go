package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/proxy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const proxyResourceName = "opentelekomcloud_taurusdb_mysql_proxy_v3.test"

func TestAccTaurusDBMySQLProxyV3_basic(t *testing.T) {
	var proxyInstance proxy.ProxyInstanceResponse

	name := "tf_taurusdb_proxy_" + acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaurusDBProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMySQLProxyV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBProxyExists(proxyResourceName, &proxyInstance),
					resource.TestCheckResourceAttrSet(proxyResourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(proxyResourceName, "flavor"),
					resource.TestCheckResourceAttr(proxyResourceName, "node_num", "2"),
					resource.TestCheckResourceAttr(proxyResourceName, "proxy_name", name),
					resource.TestCheckResourceAttrSet(proxyResourceName, "address"),
					resource.TestCheckResourceAttrSet(proxyResourceName, "port"),
					resource.TestCheckResourceAttrSet(proxyResourceName, "status"),
					resource.TestCheckResourceAttr(proxyResourceName, "master_node_weight.#", "1"),
					resource.TestCheckResourceAttr(proxyResourceName, "readonly_nodes_weight.#", "1"),
					resource.TestCheckResourceAttrSet(proxyResourceName, "nodes.#"),
				),
			},
			{
				Config: testAccTaurusDBMySQLProxyV3Update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBProxyExists(proxyResourceName, &proxyInstance),
					resource.TestCheckResourceAttr(proxyResourceName, "node_num", "3"),
					resource.TestCheckResourceAttr(proxyResourceName, "master_node_weight.#", "1"),
					resource.TestCheckResourceAttr(proxyResourceName, "readonly_nodes_weight.#", "2"),
				),
			},
			{
				ResourceName:      proxyResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTaurusDBProxyImportStateFunc(proxyResourceName),
				ImportStateVerifyIgnore: []string{
					"proxy_mode", "readonly_nodes_weight",
				},
			},
		},
	})
}

func TestAccTaurusDBMySQLProxyV3_reduceNodes(t *testing.T) {
	name := "tf_taurusdb_proxy_" + acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaurusDBProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMySQLProxyV3WithNodes(name, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(proxyResourceName, "node_num", "3"),
				),
			},
			{
				Config:      testAccTaurusDBMySQLProxyV3WithNodes(name, 2),
				ExpectError: regexp.MustCompile("reducing the number of proxy nodes is not supported"),
			},
		},
	})
}

func testAccCheckTaurusDBProxyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.TaurusDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_taurusdb_mysql_proxy_v3" {
			continue
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		proxyID := rs.Primary.ID

		proxies, err := proxy.List(client, instanceID, proxy.ListOpts{})
		if err != nil {
			return nil
		}

		for _, p := range proxies {
			if p.Proxy.PoolId == proxyID {
				return fmt.Errorf("TaurusDB MySQL proxy still exists: %s", proxyID)
			}
		}
	}

	return nil
}

func testAccCheckTaurusDBProxyExists(n string, proxyInstance *proxy.ProxyInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.TaurusDBV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating TaurusDB client: %s", err)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		proxyID := rs.Primary.ID

		proxies, err := proxy.List(client, instanceID, proxy.ListOpts{})
		if err != nil {
			return fmt.Errorf("error listing proxies: %s", err)
		}

		for _, p := range proxies {
			if p.Proxy.PoolId == proxyID {
				*proxyInstance = p
				return nil
			}
		}

		return fmt.Errorf("TaurusDB MySQL proxy not found: %s", proxyID)
	}
}

func testAccTaurusDBProxyImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}

func testAccTaurusDBMySQLProxyV3Base(name string) string {
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
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}

func testAccTaurusDBMySQLProxyV3Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
  flavor      = "gaussdb.proxy.large.x86.2"
  node_num    = 2
  proxy_name  = "%s"
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
`, testAccTaurusDBMySQLProxyV3Base(name), name)
}

func testAccTaurusDBMySQLProxyV3Update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
  flavor      = "gaussdb.proxy.large.x86.2"
  node_num    = 3
  proxy_name  = "%s"
  proxy_mode  = "readwrite"

  master_node_weight {
    id     = local.sorted_nodes[0].id
    weight = 40
  }

  readonly_nodes_weight {
    id     = local.readonly_nodes[0].id
    weight = 30
  }

  readonly_nodes_weight {
    id     = local.readonly_nodes[1].id
    weight = 30
  }
}
`, testAccTaurusDBMySQLProxyV3Base(name), name)
}

func testAccTaurusDBMySQLProxyV3WithNodes(name string, nodeNum int) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
  flavor      = "gaussdb.proxy.large.x86.2"
  node_num    = %d
  proxy_name  = "%s"
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
`, testAccTaurusDBMySQLProxyV3Base(name), nodeNum, name)
}
