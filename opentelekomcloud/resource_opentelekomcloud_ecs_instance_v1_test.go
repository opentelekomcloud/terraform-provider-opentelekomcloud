package opentelekomcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservers"
)

func TestAccEcsV1Instance_basic(t *testing.T) {
	var instance cloudservers.CloudServer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists("opentelekomcloud_ecs_instance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ecs_instance_v1.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ecs_instance_v1.instance_1", "auto_recovery", "true"),
				),
			},
			{
				Config: testAccEcsV1Instance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists("opentelekomcloud_ecs_instance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ecs_instance_v1.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_ecs_instance_v1.instance_1", "auto_recovery", "false"),
				),
			},
		},
	})
}

func TestAccEcsV1Instance_diskTypeValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEcsV1Instance_invalidTypeForAZ,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ is not supported in`),
			},
			{
				Config:      testAccEcsV1Instance_invalidType,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ doesn't exist`),
			},
			{
				Config:      testAccEcsV1Instance_invalidDataDisk,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`volume type .+ doesn't exist`),
			},
		},
	})
}

func TestAccEcsV1Instance_VPCValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEcsV1Instance_invalidVPC,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can't find VPC`),
			},
			{
				Config: testAccEcsV1Instance_computedVPC,
			},
		},
	})
}

func testAccCheckEcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ecs_instance_v1" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
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

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
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

var testAccEcsV1Instance_basic = fmt.Sprintf(`
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

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, OS_IMAGE_ID, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)

var testAccEcsV1Instance_update = fmt.Sprintf(`
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

  password                    = "Password@123"
  security_groups             = ["default", opentelekomcloud_compute_secgroup_v2.secgroup_1.name]
  availability_zone           = "%s"
  auto_recovery               = false
  delete_disks_on_termination = true

  tags = {
    foo  = "bar1"
    key1 = "value"
  }
}
`, OS_IMAGE_ID, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)

var testAccEcsV1Instance_invalidTypeForAZ = fmt.Sprintf(`
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
`, OS_IMAGE_ID, OS_VPC_ID, OS_NETWORK_ID)

var testAccEcsV1Instance_invalidType = fmt.Sprintf(`
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
`, OS_IMAGE_ID, OS_VPC_ID, OS_NETWORK_ID)

var testAccEcsV1Instance_invalidDataDisk = fmt.Sprintf(`
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
`, OS_IMAGE_ID, OS_VPC_ID, OS_NETWORK_ID)

var testAccEcsV1Instance_invalidVPC = fmt.Sprintf(`
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
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccEcsV1Instance_computedVPC = fmt.Sprintf(`
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
`, OS_IMAGE_ID)
