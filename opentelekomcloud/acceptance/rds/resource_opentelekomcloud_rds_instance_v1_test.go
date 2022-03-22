package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v1/instances"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v1/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const instanceV1ResourceName = "opentelekomcloud_rds_instance_v1.instance"

func TestAccRDSV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRDSV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSInstanceV1ConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists(instanceV1ResourceName, &instance),
					resource.TestCheckResourceAttr(instanceV1ResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(instanceV1ResourceName, "region", "eu-de"),
					resource.TestCheckResourceAttr(instanceV1ResourceName, "availabilityzone", "eu-de-01"),
					testAccCheckRDSV1InstanceTag(&instance, "foo", "bar"),
					testAccCheckRDSV1InstanceTag(&instance, "key", "value"),
				),
			},
			{
				Config: testAccSInstanceV1ConfigUpdatetag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists(instanceV1ResourceName, &instance),
					testAccCheckRDSV1InstanceTag(&instance, "foo2", "bar2"),
					testAccCheckRDSV1InstanceTag(&instance, "key", "value2"),
				),
			},
			{
				Config: testAccSInstanceV1ConfigNoTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists(instanceV1ResourceName, &instance),
					testAccCheckRDSV1InstanceNoTag(&instance),
				),
			},
			{
				Config: testAccSInstanceV1ConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists(instanceV1ResourceName, &instance),
					testAccCheckRDSV1InstanceTag(&instance, "foo", "bar"),
					testAccCheckRDSV1InstanceTag(&instance, "key", "value"),
				),
			},
		},
	})
}

func testAccCheckRDSV1InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	rdsClient, err := config.RdsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud rds: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_instance_v1" {
			continue
		}

		_, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("instance still exists. ")
		}
	}

	return nil
}

func testAccCheckRDSV1InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set. ")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		rdsClient, err := config.RdsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud rds client: %s ", err)
		}

		found, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found. ")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckRDSV1InstanceTag(
	instance *instances.Instance, k, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := common.TestAccProvider.Meta().(*cfg.Config)
		tagClient, err := config.RdsTagV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud rds client: %s ", err)
		}

		taglist, err := tags.Get(tagClient, instance.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting tags for intance: %w", err)
		}
		for _, val := range taglist.Tags {
			if k != val.Key {
				continue
			}

			if v == val.Value {
				return nil
			}

			return fmt.Errorf("bad value for %s: %s", k, val.Value)
		}

		return fmt.Errorf("tag not found: %s", k)
	}
}

func testAccCheckRDSV1InstanceNoTag(
	instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := common.TestAccProvider.Meta().(*cfg.Config)
		tagClient, err := config.RdsTagV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud rds client: %s ", err)
		}

		taglist, err := tags.Get(tagClient, instance.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting tags: %w", err)
		}

		if taglist.Tags == nil {
			return nil
		}
		if len(taglist.Tags) == 0 {
			return nil
		}

		return fmt.Errorf("expected no tags, but found %v", taglist.Tags)
	}
}

var testAccSInstanceV1ConfigBasic = fmt.Sprintf(`
%s
%s

data "opentelekomcloud_rds_flavors_v1" "flavor" {
  region            = "eu-de"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.medium.ha"
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type    = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = data.opentelekomcloud_rds_flavors_v1.flavor.id
  volume {
    type = "COMMON"
    size = 100
  }
  region           = "eu-de"
  availabilityzone = "eu-de-01"
  vpc              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  nics {
    subnetid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  securitygroup {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  dbport = "8635"
  backupstrategy {
    starttime = "01:00:00"
    keepdays  = 1
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable          = true
    replicationmode = "async"
  }
  tag = {
    foo = "bar"
    key = "value"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)

var testAccSInstanceV1ConfigUpdatetag = fmt.Sprintf(`
%s
%s

data "opentelekomcloud_rds_flavors_v1" "flavor" {
  region            = "eu-de"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.medium.ha"
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type    = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = data.opentelekomcloud_rds_flavors_v1.flavor.id
  volume {
    type = "COMMON"
    size = 100
  }
  region           = "eu-de"
  availabilityzone = "eu-de-01"
  vpc              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  nics {
    subnetid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  securitygroup {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  dbport = "8635"
  backupstrategy {
    starttime = "01:00:00"
    keepdays  = 1
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable          = true
    replicationmode = "async"
  }
  tag = {
    foo2 = "bar2"
    key  = "value2"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)

var testAccSInstanceV1ConfigNoTags = fmt.Sprintf(`
%s
%s

data "opentelekomcloud_rds_flavors_v1" "flavor" {
  region            = "eu-de"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.medium.ha"
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type    = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = data.opentelekomcloud_rds_flavors_v1.flavor.id
  volume {
    type = "COMMON"
    size = 100
  }
  region           = "eu-de"
  availabilityzone = "eu-de-01"
  vpc              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  nics {
    subnetid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  securitygroup {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  dbport = "8635"
  backupstrategy {
    starttime = "01:00:00"
    keepdays  = 1
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable          = true
    replicationmode = "async"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)
