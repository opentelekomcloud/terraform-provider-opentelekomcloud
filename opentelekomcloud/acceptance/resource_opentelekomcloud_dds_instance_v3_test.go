package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccDDSV3Instance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists("opentelekomcloud_dds_instance_v3.instance"),
					resource.TestCheckResourceAttr("opentelekomcloud_dds_instance_v3.instance", "name", "dds-instance"),
					resource.TestCheckResourceAttr("opentelekomcloud_dds_instance_v3.instance", "mode", "ReplicaSet"),
					resource.TestCheckResourceAttr("opentelekomcloud_dds_instance_v3.instance", "ssl", "true"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_minConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3Config_minConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists("opentelekomcloud_dds_instance_v3.instance"),
				),
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*cfg.Config)
	client, err := config.DdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DDSv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dds_instance_v3" {
			continue
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		ddsInstances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}

		if ddsInstances.TotalCount > 0 {
			return fmt.Errorf("instance still exists")
		}
	}

	return nil
}

func testAccCheckDDSV3InstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		config := testAccProvider.Meta().(*cfg.Config)
		client, err := config.DdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DDSv3 client: %s ", err)
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		ddsInstances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}
		if ddsInstances.TotalCount == 0 {
			return fmt.Errorf("dds instance not found")
		}

		return nil
	}
}

var TestAccDDSInstanceV3Config_basic = fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg_acc" {
  name = "secgroup_acc"
}
resource "opentelekomcloud_dds_instance_v3" "instance" {
  name              = "dds-instance"
  availability_zone = "%s"
  region            = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = opentelekomcloud_networking_secgroup_v2.sg_acc.id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type = "replica"
    num = 1
    storage = "ULTRAHIGH"
    size = 20
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days = "1"
  }
}`, OS_AVAILABILITY_ZONE, OS_REGION_NAME, OS_VPC_ID, OS_NETWORK_ID)

var TestAccDDSInstanceV3Config_minConfig = fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg_acc" {
  name = "secgroup_acc"
}
resource "opentelekomcloud_dds_instance_v3" "instance" {
  name              = "dds-instance"
  availability_zone = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = opentelekomcloud_networking_secgroup_v2.sg_acc.id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type = "replica"
    num = 1
    size = 20
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
}`, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
