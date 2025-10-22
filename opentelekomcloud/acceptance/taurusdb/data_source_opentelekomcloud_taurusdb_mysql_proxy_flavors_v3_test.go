package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceProxyFlavorsName = "data.opentelekomcloud_taurusdb_mysql_proxy_flavors_v3.test"

func TestAccDataSourceTaurusDBMysqlProxyFlavors_basic(t *testing.T) {
	name := "tf_taurusdb_proxy_flavors_" + acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBMysqlProxyFlavorsBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.#"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceProxyFlavorsName, "flavor_groups.0.flavors.0.az_status.%"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBMysqlProxyFlavorsBase(name string) string {
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
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}

func testAccDataSourceTaurusDBMysqlProxyFlavorsBasic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_taurusdb_mysql_proxy_flavors_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
}
`, testAccDataSourceTaurusDBMysqlProxyFlavorsBase(name))
}
