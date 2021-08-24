package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccEcsV1InstanceBasic(t *testing.T) {
	var instance cloudservers.CloudServer
	resourceName := "opentelekomcloud_ecs_instance_v1.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", env.OsAvailabilityZone),
					resource.TestCheckResourceAttr(resourceName, "auto_recovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccEcsV1InstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", env.OsAvailabilityZone),
					resource.TestCheckResourceAttr(resourceName, "auto_recovery", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccEcsV1InstanceDiskTypeValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEcsV1InstanceInvalidTypeForAZ,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ is not supported in`),
			},
			{
				Config:      testAccEcsV1InstanceInvalidType,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ doesn't exist`),
			},
			{
				Config:      testAccEcsV1InstanceInvalidDataDisk,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ doesn't exist`),
			},
		},
	})
}

func TestAccEcsV1InstanceVPCValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEcsV1InstanceInvalidVPC,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find VPC`),
			},
			{
				Config: testAccEcsV1InstanceComputedVPC,
			},
		},
	})
}

func TestAccEcsV1InstanceEncryption(t *testing.T) {
	var instance cloudservers.CloudServer
	resourceName := "opentelekomcloud_ecs_instance_v1.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceDataVolumeEncryption,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.kms_id", env.OsKmsID),
				),
			},
		},
	})
}

func testAccCheckEcsV1InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV1Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ecs_instance_v1" {
			continue
		}

		server, err := cloudservers.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckEcsV1InstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV1Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %s", err)
		}

		found, err := cloudservers.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}
		*instance = *found

		return nil
	}
}

var testAccEcsV1InstanceBasic = fmt.Sprintf(`
resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  data_disks {
    size = 10
    type = "SAS"
  }

  password          = "Password@123"
  availability_zone = "%s"
  auto_recovery     = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, env.OsImageID, env.OsRouterID, env.OsNetworkID, env.OsAvailabilityZone)

var testAccEcsV1InstanceUpdate = fmt.Sprintf(`
resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_ecs"
  description = "a security group"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_updated"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  data_disks {
    size = 10
    type = "SAS"
  }

  password                    = "Password@123"
  security_groups             = [opentelekomcloud_compute_secgroup_v2.secgroup_1.id]
  availability_zone           = "%s"
  auto_recovery               = false
  delete_disks_on_termination = true

  tags = {
    muh = "value-update"
  }
}
`, env.OsImageID, env.OsRouterID, env.OsNetworkID, env.OsAvailabilityZone)

var testAccEcsV1InstanceInvalidTypeForAZ = fmt.Sprintf(`
resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  system_disk_type = "uh-l1"

  password          = "Password@123"
  availability_zone = "eu-de-03"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, env.OsImageID, env.OsRouterID, env.OsNetworkID)

var testAccEcsV1InstanceInvalidType = fmt.Sprintf(`
resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  system_disk_type = "asdfasd"

  password          = "Password@123"
  availability_zone = "eu-de-03"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, env.OsImageID, env.OsRouterID, env.OsNetworkID)

var testAccEcsV1InstanceInvalidDataDisk = fmt.Sprintf(`
resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  data_disks {
    size = 10
    type = "invalid"
  }

  password          = "Password@123"
  availability_zone = "eu-de-03"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, env.OsImageID, env.OsRouterID, env.OsNetworkID)

var testAccEcsV1InstanceInvalidVPC = fmt.Sprintf(`
resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "abs"

  nics {
    network_id = "%s"
  }

  system_disk_type = "SSD"

  password          = "Password@123"
  availability_zone = "eu-de-03"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, env.OsImageID, env.OsNetworkID)

var testAccEcsV1InstanceComputedVPC = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  cidr = "192.168.0.0/16"
  name = "vpc-ecs-test"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0)
  gateway_ip = cidrhost(opentelekomcloud_vpc_v1.vpc.cidr, 1)
  name       = "subnet-ecs-test"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = opentelekomcloud_vpc_v1.vpc.id

  nics {
    network_id = opentelekomcloud_vpc_subnet_v1.subnet.id
  }

  system_disk_type = "SSD"

  password          = "Password@123"
  availability_zone = "eu-de-03"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, env.OsImageID)

var testAccEcsV1InstanceDataVolumeEncryption = fmt.Sprintf(`
resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "s2.medium.1"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  password          = "Password@123"
  availability_zone = "%s"
  auto_recovery     = true

  data_disks {
    size   = 10
    type   = "SAS"
    kms_id = "%s"
  }
}
`, env.OsImageID, env.OsRouterID, env.OsNetworkID, env.OsAvailabilityZone, env.OsKmsID)
