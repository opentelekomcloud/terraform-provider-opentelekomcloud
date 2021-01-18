package opentelekomcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRdsInstanceV3_basic(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "flavor", "rds.pg.c2.medium"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "db.0.port", "8635"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "db.0.type", "PostgreSQL"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "volume.0.size", "40"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "backup_strategy.0.keep_days", "1"),
				),
			},
			{
				Config: testAccRdsInstanceV3_update(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "volume.0.size", "100"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_ip(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_eip(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "db.0.version", "10"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "public_ips.#", "1"),
				),
			},
			{
				Config: testAccRdsInstanceV3_basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "db.0.version", "10"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "public_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_ha(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	var availabilityZone2 = os.Getenv("OS_AVAILABILITY_ZONE_2")
	if availabilityZone2 == "" {
		t.Skip("OS_AVAILABILITY_ZONE_2 is empty")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_ha(postfix, availabilityZone2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "ha_replication_mode", "semisync"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "db.0.type", "MySQL"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_optionalParams(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_optionalParams(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_backupCheck(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_backupCheck(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_templateConfigCheck(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_configTemplateBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
				),
			},
			{
				Config: testAccRdsInstanceV3_configTemplateChange(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists("opentelekomcloud_rds_instance_v3.instance", &rdsInstance),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func testAccCheckRdsInstanceV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.rdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating RDSv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_instance_v3" {
			continue
		}
		instance, _ := getRdsInstance(client, rs.Primary.ID)
		if instance != nil {
			return fmt.Errorf("RDSv3 instance still exists")
		}
	}

	return nil
}

func testAccCheckRdsInstanceV3Exists(n string, rdsInstance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		client, err := config.rdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating RDSv3 client: %s", err)
		}

		found, err := getRdsInstance(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("rdsv3 instance not found")
		}

		*rdsInstance = *found

		return nil
	}
}

func testAccRdsInstanceV3_basic(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id = "%s"
  vpc_id    = "%s"
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tag = {
    foo = "bar"
    key = "value"
  }
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_update(postfix string) string {
	return fmt.Sprintf(`
resource opentelekomcloud_networking_secgroup_v2 sg {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id = "%s"
  vpc_id    = "%s"
  volume {
    type = "COMMON"
    size = 100
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tag = {
    foo = "bar1"
    value = "key"
  }
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_eip(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id          = "%s"
  vpc_id             = "%s"
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  public_ips = [opentelekomcloud_networking_floatingip_v2.fip_1.address]
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_ha(postfix string, az2 string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s", "%s"]
  db {
    password = "MySql!120521"
    type     = "MySQL"
    version  = "5.6"
    port     = "8635"
  }
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id          = "%s"
  vpc_id             = "%s"
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  flavor = "rds.mysql.s1.large.ha"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  ha_replication_mode = "semisync"
}
`, postfix, OS_AVAILABILITY_ZONE, az2, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_optionalParams(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id          = "%s"
  vpc_id             = "%s"
  volume {
    type = "COMMON"
    size = 100
  }
  flavor = "rds.pg.c2.medium"
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_backupCheck(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id          = "%s"
  vpc_id             = "%s"
  volume {
    type = "COMMON"
    size = 100
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "10:00-11:00"
    keep_days  = 5
  }

}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_configTemplateBasic(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
	name = "pg-rds-test"
	values = {
		max_connections = "10"
		autocommit = "OFF"
	}
	datastore {
		type = "postgresql"
		version = "10"
	}
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg2" {
	name = "pg-rds-test-2"
	values = {
		max_connections = "10"
		autocommit = "OFF"
	}
	datastore {
		type = "postgresql"
		version = "10"
	}
}

resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id          = "%s"
  vpc_id             = "%s"
  volume {
    type = "COMMON"
    size = 40
  }
  flavor         = "rds.pg.c2.medium"
  param_group_id = opentelekomcloud_rds_parametergroup_v3.pg.id
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_configTemplateChange(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
	name = "pg-rds-test"
	values = {
		max_connections = "10"
		autocommit = "OFF"
	}
	datastore {
		type = "postgresql"
		version = "10"
	}
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg2" {
	name = "pg-rds-test-2"
	values = {
		max_connections = "10"
		autocommit = "OFF"
	}
	datastore {
		type = "postgresql"
		version = "10"
	}
}

resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id          = "%s"
  vpc_id             = "%s"
  volume {
    type = "COMMON"
    size = 40
  }
  flavor         = "rds.pg.c2.medium"
  param_group_id = opentelekomcloud_rds_parametergroup_v3.pg2.id
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}
