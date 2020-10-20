package opentelekomcloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud"
)

func TestAccRdsInstanceV3_basic(t *testing.T) {
	name := acctest.RandString(3)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(),
				),
			},
			{
				Config: testAccRdsInstanceV3_eip(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(),
					resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "public_ips.#", "1"),
				),
			},
			//{
			//	Config: testAccRdsInstanceV3_update(name),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckRdsInstanceV3Exists(),
			//		resource.TestCheckResourceAttr("opentelekomcloud_rds_instance_v3.instance", "public_ips.#", "0"),
			//	),
			//},
		},
	})
}

func testAccRdsInstanceV3_basic(val string) string {
	return fmt.Sprintf(`
resource opentelekomcloud_networking_secgroup_v2 sg {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type = "PostgreSQL"
    version = "10"
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
}
`, OS_AVAILABILITY_ZONE, val, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccRdsInstanceV3_update(val string) string {
	return fmt.Sprintf(`
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
    size = 200
  }
  flavor = "rds.pg.c2.medium"
  backup_strategy {
    start_time = "09:00-10:00"
    keep_days = 2
  }
  tag = {
    foo1 = "bar1"
    key = "value1"
  }
}
	`, OS_AVAILABILITY_ZONE, val, OS_NETWORK_ID, OS_VPC_ID)
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

func testAccCheckRdsInstanceV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.rdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_instance_v3" {
			continue
		}

		_, err = fetchRdsInstanceV3ByListOnTest(rs, client)
		if err != nil {
			if strings.Index(err.Error(), "Error finding the resource by list api") != -1 {
				return nil
			}
			return err
		}
		return fmt.Errorf("opentelekomcloud_rds_instance_v3 still exists")
	}

	return nil
}

func testAccCheckRdsInstanceV3Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.rdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["opentelekomcloud_rds_instance_v3.instance"]
		if !ok {
			return fmt.Errorf("Error checking opentelekomcloud_rds_instance_v3.instance exist, err=not found this resource")
		}

		_, err = fetchRdsInstanceV3ByListOnTest(rs, client)
		if err != nil {
			if strings.Index(err.Error(), "Error finding the resource by list api") != -1 {
				return fmt.Errorf("opentelekomcloud_rds_instance_v3 is not exist")
			}
			return fmt.Errorf("Error checking opentelekomcloud_rds_instance_v3.instance exist, err=%s", err)
		}
		return nil
	}
}

func fetchRdsInstanceV3ByListOnTest(rs *terraform.ResourceState,
	client *golangsdk.ServiceClient) (interface{}, error) {

	identity := map[string]interface{}{"id": rs.Primary.ID}

	queryLink := "?id=" + identity["id"].(string)

	link := client.ServiceURL("instances") + queryLink

	return findRdsInstanceV3ByList(client, link, identity)
}
