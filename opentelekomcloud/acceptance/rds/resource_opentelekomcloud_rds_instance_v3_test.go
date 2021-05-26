package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/rds"
)

const resourceName = "opentelekomcloud_rds_instance_v3.instance"

func TestAccRdsInstanceV3Basic(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c2.medium"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(resourceName, "db.0.type", "PostgreSQL"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
					resource.TestCheckResourceAttr(resourceName, "tags.kuh", "value-create"),
				),
			},
			{
				Config: testAccRdsInstanceV3Update(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "muh.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3ElasticIP(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3ElasticIP(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(resourceName, "db.0.version", "10"),
					resource.TestCheckResourceAttr(resourceName, "public_ips.#", "1"),
				),
			},
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "db.0.version", "10"),
					resource.TestCheckResourceAttr(resourceName, "public_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3HA(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	var availabilityZone2 = os.Getenv("OS_AVAILABILITY_ZONE_2")
	if availabilityZone2 == "" {
		t.Skip("OS_AVAILABILITY_ZONE_2 is empty")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3HA(postfix, availabilityZone2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "semisync"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "db.0.type", "MySQL"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3OptionalParams(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3OptionalParams(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3Backup(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3Backup(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3TemplateConfig(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3ConfigTemplateBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
			{
				Config: testAccRdsInstanceV3ConfigTemplateUpdate(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &rdsInstance),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3InvalidDBVersion(t *testing.T) {
	postfix := acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccRdsInstanceV3InvalidDBVersion(postfix),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find version.+`),
			},
		},
	})
}

func testAccCheckRdsInstanceV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.RdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating RDSv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_instance_v3" {
			continue
		}
		instance, _ := rds.GetRdsInstance(client, rs.Primary.ID)
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
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.RdsV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating RDSv3 client: %s", err)
		}

		found, err := rds.GetRdsInstance(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("RDSv3 instance not found")
		}

		*rdsInstance = *found

		return nil
	}
}

func testAccRdsInstanceV3Basic(postfix string) string {
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
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3Update(postfix string) string {
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
  tags = {
    muh = "value-update"
  }
}
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3ElasticIP(postfix string) string {
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
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3HA(postfix string, az2 string) string {
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
`, postfix, env.OS_AVAILABILITY_ZONE, az2, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3OptionalParams(postfix string) string {
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
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3Backup(postfix string) string {
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
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3ConfigTemplateBasic(postfix string) string {
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
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3ConfigTemplateUpdate(postfix string) string {
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
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}

func testAccRdsInstanceV3InvalidDBVersion(postfix string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!12052"
    type     = "PostgreSQL"
    version  = "5.6"
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
}
`, postfix, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID, env.OS_VPC_ID)
}
