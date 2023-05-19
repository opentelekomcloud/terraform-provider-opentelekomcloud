package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/rds"
)

const instanceV3ResourceName = "opentelekomcloud_rds_instance_v3.instance"

func TestAccRdsInstanceV3Basic(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "PostgreSQL"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.muh", "value-create"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.kuh", "value-create"),
				),
			},
			{
				Config: testAccRdsInstanceV3Update(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.muh", "value-update"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8636"),
				),
			},
		},
	})
}

func TestAccRdsPostgre13V3ParamsBasic(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsPostrgre13V3ParamsBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "PostgreSQL"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.muh", "value-create"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.kuh", "value-create"),
				),
			},
			{
				Config: testAccRdsPostrgre13V3ParamsBasicUpdate(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3RestoreBackup(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse
	var restoredInstance instances.InstanceResponse

	restoredResourceName := "opentelekomcloud_rds_instance_v3.from_backup"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
				),
			},
			{
				PreConfig: func() {
					forceRdsBackup(t, &rdsInstance.Id)
				},
				Config: testAccRdsInstanceV3Restored(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					testAccCheckRdsInstanceV3Exists(restoredResourceName, &restoredInstance),
					resource.TestCheckResourceAttr(restoredResourceName, "flavor", "rds.pg.c2.large"),
					resource.TestCheckResourceAttr(restoredResourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(restoredResourceName, "tags.muh", "value-create"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3ElasticIP(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3ElasticIP(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.version", "10"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "public_ips.#", "1"),
				),
			},
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.version", "10"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "public_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3HA(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	var availabilityZone2 = os.Getenv("OS_AVAILABILITY_ZONE_2")
	if availabilityZone2 == "" {
		t.Skip("OS_AVAILABILITY_ZONE_2 is empty")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3HA(postfix, availabilityZone2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "ha_replication_mode", "semisync"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "MySQL"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3OptionalParams(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3OptionalParams(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3Backup(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3Backup(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3TemplateConfig(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3ConfigTemplateBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
			{
				Config: testAccRdsInstanceV3ConfigTemplateUpdate(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3InvalidDBVersion(t *testing.T) {
	postfix := acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccRdsInstanceV3InvalidDBVersion(postfix),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find version.+`),
			},
		},
	})
}

func TestAccRdsInstanceV3InvalidFlavor(t *testing.T) {
	postfix := acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccRdsInstanceV3InvalidFlavor(postfix),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find flavor.+`),
			},
		},
	})
}

func TestAccRdsInstanceV3_configurationParameters(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3ConfigurationOverride(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "parameters.max_connections", "37"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3TimeZoneAndSSL(t *testing.T) {
	// Test is failing on deletion because RDSv3 SSL switchover doesn't change instance `action` status / doesn't
	// return job_id but still blocks the instance from performing other actions like port_change/instance_deletion
	// https://jira.tsi-dev.otc-service.com/browse/OTCDB-3026
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3ConfigurationTimeZone(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "MySQL"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.mysql.m1.large"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "ssl_enable", "true"),
				),
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

func testAccCheckRdsInstanceV3Exists(n string, rdsInstance *instances.InstanceResponse) resource.TestCheckFunc {
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
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3Update(postfix string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8636"
  }
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
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
`, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ElasticIP(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  public_ips = [opentelekomcloud_networking_floatingip_v2.fip_1.address]
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3HA(postfix string, az2 string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s", "%s"]
  db {
    password = "MySql!120521"
    type     = "MySQL"
    version  = "5.6"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
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
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE, az2)
}

func testAccRdsInstanceV3OptionalParams(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 100
  }
  flavor = "rds.pg.c2.large"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3Backup(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 100
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "10:00-11:00"
    keep_days  = 5
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ConfigTemplateBasic(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
  name = "pg-rds-test"
  values = {
    max_connections = "1200"
    autocommit      = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "12"
  }
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg2" {
  name = "pg-rds-test-2"
  values = {
    max_connections = "1200"
    autocommit      = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "12"
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "12"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor         = "rds.pg.c2.large"
  param_group_id = opentelekomcloud_rds_parametergroup_v3.pg.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ConfigTemplateUpdate(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
  name = "pg-rds-test"
  values = {
    max_connections = "1200"
    autocommit      = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "10"
  }
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg2" {
  name = "pg-rds-test-2"
  values = {
    max_connections = "1200"
    autocommit      = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "12"
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "12"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor         = "rds.pg.c2.large"
  param_group_id = opentelekomcloud_rds_parametergroup_v3.pg2.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3InvalidDBVersion(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!12052"
    type     = "PostgreSQL"
    version  = "5.6"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ConfigurationOverride(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
  name = "pg-rds-test"
  values = {
    autocommit = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "10"
  }
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

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  flavor            = "rds.pg.c2.large"
  volume {
    type = "COMMON"
    size = 40
  }

  parameters = {
    max_connections = "37",
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3InvalidFlavor(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "bla.bla.rds"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3Restored(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

data "opentelekomcloud_rds_backup_v3" "backup" {
  instance_id = opentelekomcloud_rds_instance_v3.instance.id
  type        = "auto"
}

resource "opentelekomcloud_rds_instance_v3" "from_backup" {
  name              = "${opentelekomcloud_rds_instance_v3.instance.name}-restored"
  availability_zone = opentelekomcloud_rds_instance_v3.instance.availability_zone

  restore_point {
    instance_id = opentelekomcloud_rds_instance_v3.instance.id
    backup_id   = data.opentelekomcloud_rds_backup_v3.backup.id
  }

  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ConfigurationTimeZone(postfix string) string {
	return fmt.Sprintf(`
%s
%s
resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "MySql!112822"
    type     = "MySQL"
    version  = "8.0"
    port     = "8635"
  }
  param_group_id    = opentelekomcloud_rds_parametergroup_v3.pg_1.id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor     = "rds.mysql.m1.large"
  ssl_enable = true
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}


resource "opentelekomcloud_rds_parametergroup_v3" "pg_1" {
  name        = "pg_tmz"
  description = "time zone template"

  values = {
    time_zone = "Africa/Casablanca"
  }

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}


`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsPostrgre13V3ParamsBasic(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "13"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  parameters = {
    max_prepared_transactions = "200"
    max_connections           = "250"
    autovacuum_max_workers    = "20"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsPostrgre13V3ParamsBasicUpdate(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "13"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  parameters = {
    max_prepared_transactions = "250"
    max_connections           = "300"
    autovacuum_max_workers    = "30"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}
