package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccTaurusDBMySqlInstanceDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_taurusdb_mysql_instance_v3.test"
	rName := "tf_taurusdb_instance" + acctest.RandString(3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMySqlInstanceDataSourceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBMySqlInstanceDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "flavor", "gaussdb.mysql.xlarge.x86.8"),
					resource.TestCheckResourceAttr(dataSourceName, "availability_zone_mode", "multi"),
					resource.TestCheckResourceAttr(dataSourceName, "master_availability_zone", "eu-de-01"),
					resource.TestCheckResourceAttr(dataSourceName, "read_replicas", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "port", "3306"),
					resource.TestCheckResourceAttr(dataSourceName, "backup_strategy.0.start_time", "03:00-04:00"),
					resource.TestCheckResourceAttr(dataSourceName, "backup_strategy.0.keep_days", "7"),
					resource.TestCheckResourceAttr(dataSourceName, "datastore.0.engine", "gaussdb-mysql"),
					resource.TestCheckResourceAttr(dataSourceName, "datastore.0.version", "8.0.22"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "db_user_name"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBMySqlInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB MySQL instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the TaurusDB MySQL instance data source ID not set")
		}

		return nil
	}
}

func testAccTaurusDBMySqlInstanceDataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_taurusdb_mysql_instance_v3" "test" {
  name                     = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id        = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor                   = "gaussdb.mysql.xlarge.x86.8"
  password                 = "Test123!@#"
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 1

  datastore {
    engine  = "gaussdb-mysql"
    version = "8.0"
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }
}

data "opentelekomcloud_taurusdb_mysql_instance_v3" "test" {
  name = opentelekomcloud_taurusdb_mysql_instance_v3.test.name

  depends_on = [
    opentelekomcloud_taurusdb_mysql_instance_v3.test,
  ]
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}
