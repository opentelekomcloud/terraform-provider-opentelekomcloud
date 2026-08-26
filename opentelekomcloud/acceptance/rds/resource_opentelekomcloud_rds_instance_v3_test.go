package acceptance

import (
	"fmt"
	"net"
	"os"
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.n1.large.4"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "PostgreSQL"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.keep_days", "0"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.muh", "value-create"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.kuh", "value-create"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "lower_case_table_names", "0"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "ssl_enable", "true"),
				),
			},
			{
				Config: testAccRdsInstanceV3Update(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.n1.large.4"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "tags.muh", "value-update"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "name", "tf_rds_instance_updated_"+postfix),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "lower_case_table_names", "0"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.period", "1,2,3,4"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.n1.large.4"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.n1.large.4"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "lower_case_table_names", "0"),
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
					resource.TestCheckResourceAttr(restoredResourceName, "flavor", "rds.pg.c2.medium"),
					resource.TestCheckResourceAttr(restoredResourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(restoredResourceName, "tags.muh", "value-create"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "lower_case_table_names", "0"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.version", "15"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "public_ips.#", "1"),
				),
			},
			{
				Config: testAccRdsInstanceV3Basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.version", "15"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.type", "MySQL"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "availability_zones.#", "2"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "ssl_enable", "false"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "backup_strategy.0.period", "1,2"),
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "ssl_enable", "true"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3AutoScaling(t *testing.T) {
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
					resource.TestCheckResourceAttr(instanceV3ResourceName, "db.0.port", "8635"),
				),
			},
			{
				Config: testAccRdsInstanceV3AutoScaling(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.limit_size", "500"),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "volume.0.trigger_threshold", "10"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3RestoreToPITR_NewInstance(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3RestorePITRBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
				),
			},
			{
				Config: testAccRdsInstanceV3RestorePITRUpdate(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
				),
			},
			{
				Config: testAccRdsInstanceV3RestorePITRBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_SSLEnable(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3EnableSSL(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "ssl_enable", "true"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_PrivateIp(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3PrivateIp(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttrPair(instanceV3ResourceName, "private_ips.0", instanceV3ResourceName, "private_ip"),
					testAccCheckRdsInstanceV3PrivateIp(),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_PrivateDNS(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3PrivateDNS(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "private_domain_name", "testAccDomain"),
				),
			},
		},
	})
}

func TestAccRdsInstanceV3_Swiss(t *testing.T) {
	postfix := acctest.RandString(3)
	var rdsInstance instances.InstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_Swiss(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(instanceV3ResourceName, &rdsInstance),
					resource.TestCheckResourceAttr(instanceV3ResourceName, "flavor", "rds.pg.s3.large.2"),
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

func testAccCheckRdsInstanceV3PrivateIp() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		subnet, ok := s.RootModule().Resources["data.opentelekomcloud_vpc_subnet_v1.shared_subnet"]
		if !ok {
			return fmt.Errorf("subnet not found in state")
		}

		subnetCIDR := subnet.Primary.Attributes["cidr"]

		rs, ok := s.RootModule().Resources[instanceV3ResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", instanceV3ResourceName)
		}

		actualIP := rs.Primary.Attributes["private_ip"]

		_, network, _ := net.ParseCIDR(subnetCIDR)
		ip := make(net.IP, len(network.IP))
		copy(ip, network.IP)

		ip[len(ip)-1] += byte(10)

		if actualIP != ip.String() {
			return fmt.Errorf("expected private_ip %s (from CIDR %s + 10), got %s", ip.String(), subnetCIDR, actualIP)
		}

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
    version  = "17"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = "rds.pg.n1.large.4"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 0
    period     = "1,2,3"
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  lower_case_table_names = "0"
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
  name              = "tf_rds_instance_updated_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "17"
    port     = "8636"
  }
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  flavor = "rds.pg.n1.large.4"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 2
    period     = "1,2,3,4"
  }
  tags = {
    muh = "value-update"
  }
  lower_case_table_names = "0"
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
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = "rds.pg.n1.large.4"
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

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "ha"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%s", "unsupported") == "normal" && lookup(flavor.az_status, "%s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s", "%s"]
  db {
    password = "MySql!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  flavor = local.available_flavor.name
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  ha_replication_mode = "semisync"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, az2, postfix, env.OS_AVAILABILITY_ZONE, az2)
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
    version  = "15"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 100
  }

  time_zone = "UTC+04:00"
  flavor    = "rds.pg.n1.large.4"
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
    version  = "15"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  flavor = "rds.pg.n1.large.4"
  backup_strategy {
    start_time = "10:00-11:00"
    keep_days  = 5
    period     = "1,2"
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
    version = "15"
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
    version = "15"
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor         = "rds.pg.n1.large.4"
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
    version = "15"
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
    version = "15"
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor         = "rds.pg.n1.large.4"
  param_group_id = opentelekomcloud_rds_parametergroup_v3.pg2.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ConfigurationOverride(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "15"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
  name = "pg-rds-test-%[3]s"
  values = {
    autocommit = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "15"
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]

  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }

  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  flavor            = local.available_flavor.name
  volume {
    type = "CLOUDSSD"
    size = 40
  }

  parameters = {
    max_connections = "37",
  }
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
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "ULTRAHIGH"
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
  lower_case_table_names = "0"
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
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "ULTRAHIGH"
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
  lower_case_table_names = "0"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3ConfigurationTimeZone(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
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
    type = "CLOUDSSD"
    size = 40
  }
  flavor     = local.available_flavor.name
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
  name        = "pg_tmz-%[3]s"
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
    type = "CLOUDSSD"
    size = 40
  }
  flavor = "rds.pg.n1.large.4"
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
    type = "CLOUDSSD"
    size = 40
  }
  flavor = "rds.pg.n1.large.4"
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

func testAccRdsInstanceV3RestorePITRBasic(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "15"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = local.available_flavor.name
}

resource "opentelekomcloud_rds_instance_v3" "instance_2" {
  name              = "tf_rds_source_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }

  flavor = local.available_flavor.name
}

resource "opentelekomcloud_rds_backup_v3" "test" {
  instance_id = opentelekomcloud_rds_instance_v3.instance_2.id
  name        = "tf_rds_backup_%[3]s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3RestorePITRUpdate(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "15"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = local.available_flavor.name
  restore_point {
    instance_id = opentelekomcloud_rds_backup_v3.test.instance_id
    backup_id   = opentelekomcloud_rds_backup_v3.test.id
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance_2" {
  name              = "tf_rds_source_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }

  flavor = local.available_flavor.name
}

resource "opentelekomcloud_rds_backup_v3" "test" {
  instance_id = opentelekomcloud_rds_instance_v3.instance_2.id
  name        = "tf_rds_backup_%[3]s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3AutoScaling(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
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
    type              = "CLOUDSSD"
    size              = 100
    limit_size        = 500
    trigger_threshold = 10
  }
  flavor     = local.available_flavor.name
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
  name        = "pg_tmz-%[3]s"
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

func testAccRdsInstanceV3EnableSSL(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = "3306"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor     = local.available_flavor.name
  ssl_enable = true
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  public_ips = [opentelekomcloud_networking_floatingip_v2.fip_1.address]
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3PrivateIp(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "15"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  private_ip        = cidrhost(data.opentelekomcloud_vpc_subnet_v1.shared_subnet.cidr, 10)
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = local.available_flavor.name
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3PrivateDNS(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
}

locals {
  available_flavor = [
    for flavor in data.opentelekomcloud_rds_flavors_v3.flavor.flavors :
    flavor if lookup(flavor.az_status, "%[4]s", "unsupported") == "normal"
  ][0]
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = "3306"
  }
  security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  private_domain_name = "testAccDomain"
  volume {
    type = "CLOUDSSD"
    size = 40
  }
  flavor = local.available_flavor.name
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}

func testAccRdsInstanceV3_Swiss(postfix string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%[3]s"
  availability_zone = ["%[4]s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "15"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "ULTRAHIGH"
    size = 40
  }
  flavor = "rds.pg.s3.large.2"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}
