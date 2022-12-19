package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceInstanceName = "opentelekomcloud_dds_instance_v3.instance"

func TestAccDDSV3Instance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3ConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceInstanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", "dds-instance"),
					resource.TestCheckResourceAttr(resourceInstanceName, "mode", "ReplicaSet"),
					resource.TestCheckResourceAttr(resourceInstanceName, "ssl", "true"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_minConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3ConfigMinConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceInstanceName),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_single(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3ConfigSingle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceInstanceName),
					resource.TestCheckResourceAttr(resourceInstanceName, "name", "dds-instance"),
					resource.TestCheckResourceAttr(resourceInstanceName, "mode", "Single"),
				),
			},
		},
	})
}

func TestAccDDSInstanceV3_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3ConfigBasic,
			},
			{
				ResourceName:      resourceInstanceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"flavor",
					"password",
					"availability_zone",
				},
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DDSv3 client: %w", err)
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
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DdsV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DDSv3 client: %w", err)
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

var TestAccDDSInstanceV3ConfigBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dds_instance_v3" "instance" {
  name              = "dds-instance"
  availability_zone = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "1"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var TestAccDDSInstanceV3ConfigMinConfig = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dds_instance_v3" "instance" {
  name              = "dds-instance"
  availability_zone = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 1
    size      = 20
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var TestAccDDSInstanceV3ConfigSingle = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dds_instance_v3" "instance" {
  name              = "dds-instance"
  availability_zone = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  password          = "5ecuredPa55w0rd@"
  mode              = "Single"
  flavor {
    type      = "single"
    num       = 1
    size      = 20
    spec_code = "dds.mongodb.s2.medium.4.single"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
