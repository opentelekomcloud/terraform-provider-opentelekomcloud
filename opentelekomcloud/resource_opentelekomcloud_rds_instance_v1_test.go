package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/tags"
)

func TestAccRDSV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRDSV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSInstanceV1Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("opentelekomcloud_rds_instance_v1.instance", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_instance_v1.instance", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_instance_v1.instance", "region", "eu-de"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_instance_v1.instance", "availabilityzone", "eu-de-01"),
					testAccCheckRDSV1InstanceTag(&instance, "foo", "bar"),
					testAccCheckRDSV1InstanceTag(&instance, "key", "value"),
				),
			},
			{
				Config: testAccSInstanceV1Config_updatetag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("opentelekomcloud_rds_instance_v1.instance", &instance),
					testAccCheckRDSV1InstanceTag(&instance, "foo2", "bar2"),
					testAccCheckRDSV1InstanceTag(&instance, "key", "value2"),
				),
			},
			{
				Config: testAccSInstanceV1Config_notags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("opentelekomcloud_rds_instance_v1.instance", &instance),
					testAccCheckRDSV1InstanceNoTag(&instance),
				),
			},
			{
				Config: testAccSInstanceV1Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("opentelekomcloud_rds_instance_v1.instance", &instance),
					testAccCheckRDSV1InstanceTag(&instance, "foo", "bar"),
					testAccCheckRDSV1InstanceTag(&instance, "key", "value"),
				),
			},
		},
	})
}

func testAccCheckRDSV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	rdsClient, err := config.rdsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud rds: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_instance_v1" {
			continue
		}

		_, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckRDSV1InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		rdsClient, err := config.rdsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud rds client: %s ", err)
		}

		found, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found. ")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckRDSV1InstanceTag(
	instance *instances.Instance, k, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		tagClient, err := config.rdsTagV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud rds client: %s ", err)
		}

		taglist, err := tags.Get(tagClient, instance.ID).Extract()
		for _, val := range taglist.Tags {
			if k != val.Key {
				continue
			}

			if v == val.Value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, val.Value)
		}

		return fmt.Errorf("Tag not found: %s", k)
	}
}

func testAccCheckRDSV1InstanceNoTag(
	instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		tagClient, err := config.rdsTagV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud rds client: %s ", err)
		}

		taglist, err := tags.Get(tagClient, instance.ID).Extract()

		if taglist.Tags == nil {
			return nil
		}
		if len(taglist.Tags) == 0 {
			return nil
		}

		return fmt.Errorf("Expected no tags, but found %v", taglist.Tags)
	}
}

var testAccSInstanceV1Config_basic = fmt.Sprintf(`
data "opentelekomcloud_rds_flavors_v1" "flavor" {
    region = "eu-de"
    datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.medium.ha"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = "${data.opentelekomcloud_rds_flavors_v1.flavor.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "eu-de"
  availabilityzone = "eu-de-01"
  vpc = "%s"
  nics {
    subnetid = "%s"
  }
  securitygroup {
    id = "${opentelekomcloud_compute_secgroup_v2.secgrp_rds.id}"
  }
  dbport = "8635"
  backupstrategy {
    starttime = "01:00:00"
    keepdays = 1
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable = true
    replicationmode = "async"
  }
  tag = {
    foo = "bar"
    key = "value"
  }
  depends_on = ["opentelekomcloud_compute_secgroup_v2.secgrp_rds"]
}`, OS_VPC_ID, OS_NETWORK_ID)

var testAccSInstanceV1Config_updatetag = fmt.Sprintf(`
data "opentelekomcloud_rds_flavors_v1" "flavor" {
    region = "eu-de"
    datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.medium.ha"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = "${data.opentelekomcloud_rds_flavors_v1.flavor.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "eu-de"
  availabilityzone = "eu-de-01"
  vpc = "%s"
  nics {
    subnetid = "%s"
  }
  securitygroup {
    id = "${opentelekomcloud_compute_secgroup_v2.secgrp_rds.id}"
  }
  dbport = "8635"
  backupstrategy {
    starttime = "01:00:00"
    keepdays = 1
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable = true
    replicationmode = "async"
  }
  tag = {
    foo2 = "bar2"
    key = "value2"
  }
  depends_on = ["opentelekomcloud_compute_secgroup_v2.secgrp_rds"]
}`, OS_VPC_ID, OS_NETWORK_ID)

var testAccSInstanceV1Config_notags = fmt.Sprintf(`
data "opentelekomcloud_rds_flavors_v1" "flavor" {
    region = "eu-de"
    datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.medium.ha"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = "${data.opentelekomcloud_rds_flavors_v1.flavor.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "eu-de"
  availabilityzone = "eu-de-01"
  vpc = "%s"
  nics {
    subnetid = "%s"
  }
  securitygroup {
    id = "${opentelekomcloud_compute_secgroup_v2.secgrp_rds.id}"
  }
  dbport = "8635"
  backupstrategy {
    starttime = "01:00:00"
    keepdays = 1
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable = true
    replicationmode = "async"
  }
  depends_on = ["opentelekomcloud_compute_secgroup_v2.secgrp_rds"]
}`, OS_VPC_ID, OS_NETWORK_ID)
