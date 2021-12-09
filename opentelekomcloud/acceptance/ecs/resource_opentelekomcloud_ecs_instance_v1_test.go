package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceInstanceV1Name = "opentelekomcloud_ecs_instance_v1.instance_1"

func TestAccEcsV1InstanceBasic(t *testing.T) {
	var instance cloudservers.CloudServer
	qts := serverQuotas(10+4, "s2.medium.1")
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "availability_zone", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "auto_recovery", "true"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccEcsV1InstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "availability_zone", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "auto_recovery", "false"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccEcsV1Instance_import(t *testing.T) {
	t.Parallel()
	qts := serverQuotas(10+4, "s2.medium.1")
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceBasic,
			},
			{
				ResourceName:      resourceInstanceV1Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"delete_disks_on_termination",
				},
			},
		},
	})
}

func TestAccEcsV1InstanceDiskTypeValidation(t *testing.T) {
	t.Parallel()
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
	qts := serverQuotas(4, "s2.medium.1")
	t.Parallel()
	quotas.BookMany(t, qts)

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
	qts := serverQuotas(10+4, "s2.medium.1")
	t.Parallel()
	quotas.BookMany(t, qts)

	if env.OS_KMS_ID == "" {
		t.Skip("OS_KMS_ID is not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceDataVolumeEncryption,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "data_disks.0.kms_id", env.OS_KMS_ID),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "system_disk_kms_id", env.OS_KMS_ID),
				),
			},
		},
	})
}

func testAccCheckEcsV1InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %w", err)
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
		client, err := config.ComputeV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %w", err)
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
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  data_disks {
    size = 10
    type = "SAS"
  }

  password          = "Password@123"
  availability_zone = "%s"
  auto_recovery     = true
  delete_disks_on_termination = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccEcsV1InstanceUpdate = fmt.Sprintf(`
%s

%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_updated"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  data_disks {
    size = 10
    type = "SAS"
  }

  password                    = "Password@123"
  security_groups             = [data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id]
  availability_zone           = "%s"
  auto_recovery               = false
  delete_disks_on_termination = true

  tags = {
    muh = "value-update"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccEcsV1InstanceInvalidTypeForAZ = fmt.Sprintf(`
%s

%s


resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet)

var testAccEcsV1InstanceInvalidType = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet)

var testAccEcsV1InstanceInvalidDataDisk = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  data_disks {
    size = 10
    type = "invalid"
  }
  delete_disks_on_termination = true

  password          = "Password@123"
  availability_zone = "eu-de-03"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceImage, common.DataSourceSubnet)

var testAccEcsV1InstanceInvalidVPC = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = "abs"

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet)

var testAccEcsV1InstanceComputedVPC = fmt.Sprintf(`
%s

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
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
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
`, common.DataSourceImage)

var testAccEcsV1InstanceDataVolumeEncryption = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  password           = "Password@123"
  availability_zone  = "%s"
  auto_recovery      = true
  system_disk_kms_id = "%[4]s"

  data_disks {
    size   = 10
    type   = "SAS"
    kms_id = "%[4]s"
  }
  delete_disks_on_termination = true
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OS_KMS_ID)
