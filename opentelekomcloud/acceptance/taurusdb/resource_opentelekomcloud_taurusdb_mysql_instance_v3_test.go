package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const instanceV3ResourceName = "opentelekomcloud_taurusdb_mysql_instance_v3.instance"

func TestAccTaurusDBMySqlInstanceV3_basic(t *testing.T) {
	var instance instance.TaurusDBInstance

	name := "tf_taurusdb_instance" + acctest.RandString(3)
	updateName := "tf_taurusdb_update" + acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaurusDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMySqlInstanceV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBInstanceExists(instanceV3ResourceName, &instance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", name),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "gaussdb.mysql.xlarge.x86.8"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "availability_zone_mode", "multi"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "master_availability_zone", "eu-de-01"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "read_replicas", "1"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "port", "3306"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.start_time", "03:00-04:00"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.keep_days", "7"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "datastore.0.engine", "gaussdb-mysql"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "datastore.0.version", "8.0"),
					resource.TestCheckResourceAttrSet(instanceV3ResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(instanceV3ResourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(instanceV3ResourceName, "security_group_id"),
				),
			},
			{
				Config: testAccTaurusDBMySqlInstanceV3Update(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBInstanceExists(instanceV3ResourceName, &instance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", updateName),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "gaussdb.mysql.2xlarge.x86.8"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
			{
				ResourceName:      instanceV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccCheckTaurusDBInstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.TaurusDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_taurusdb_mysql_instance_v3" {
			continue
		}

		v, err := instance.Get(client, rs.Primary.ID)
		if err == nil && v.Id == rs.Primary.ID {
			return fmt.Errorf("TaurusDB instance <%s> still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaurusDBInstanceExists(n string, taurusInstance *instance.TaurusDBInstance) resource.TestCheckFunc {
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

		found, err := instance.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("TaurusDB instance <%s> not found", rs.Primary.ID)
		}
		*taurusInstance = *found

		return nil
	}
}

func testAccTaurusDBMySqlInstanceV3Basic(postfix string) string {
	return fmt.Sprintf(`
%s
%s
resource "opentelekomcloud_taurusdb_mysql_instance_v3" "instance" {
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
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, postfix)
}

func testAccTaurusDBMySqlInstanceV3Update(postfix string) string {
	return fmt.Sprintf(`
%s
%s
resource "opentelekomcloud_taurusdb_mysql_instance_v3" "instance" {
  name                     = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id        = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor                   = "gaussdb.mysql.2xlarge.x86.8"
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
    keep_days  = 8
  }
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, postfix)
}
