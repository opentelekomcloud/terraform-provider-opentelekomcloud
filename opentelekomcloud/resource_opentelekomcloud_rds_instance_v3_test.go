package opentelekomcloud

import (
	"fmt"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"testing"

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
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "volume.0.size", "100"),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "backup_strategy.0.keep_days", "1"),
				),
			},
			{
				Config: testAccRdsInstanceV3_update(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "flavor", "rds.pg.c2.large"),
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
		_, err = getRdsInstance(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("relational Database still exists")
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
    foo = "bar"
    key = "value"
  }
}
`, postfix, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_eip(val string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource opentelekomcloud_networking_secgroup_v2 sg {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type = "PostgreSQL"
    version = "9.5"
    port = "8635"
  }
  name = "tf_rds_instance_%s"
  security_group_id  = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id = "%s"
  vpc_id = "%s"
  volume {
    type = "COMMON"
    size = 100
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days = 1
  }
  tag = {
    foo = "bar"
    key = "value"
  }
  public_ips = [opentelekomcloud_networking_floatingip_v2.fip_1.address]
}
`, OS_AVAILABILITY_ZONE, val, OS_NETWORK_ID, OS_VPC_ID)
}
