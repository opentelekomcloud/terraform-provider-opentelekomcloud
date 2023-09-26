package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservers"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceInstanceV1Name = "opentelekomcloud_ecs_instance_v1.instance_1"

func TestAccEcsV1InstanceBasic(t *testing.T) {
	var instance cloudservers.CloudServer
	qts := serverQuotas(10+4, getFlavorName())
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

func TestAccEcsV1InstanceIp(t *testing.T) {
	var instance cloudservers.CloudServer
	qts := serverQuotas(10+4, getFlavorName())
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "availability_zone", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "auto_recovery", "false"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet(resourceInstanceV1Name, "nics.0.port_id"),
				),
			},
		},
	})
}

func TestAccEcsV1InstanceDeleted(t *testing.T) {
	var instance cloudservers.CloudServer
	qts := serverQuotas(10+4, getFlavorName())
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceBasic,
				Check:  testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
			},
			{
				PreConfig: func() {
					testAccEcsV1InstanceDeleted(t, instance.ID)
				},
				Config: testAccEcsV1InstanceBasic,
			},
		},
	})
}

func testAccEcsV1InstanceDeleted(t *testing.T, id string) {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV1Client(env.OS_REGION_NAME)
	th.AssertNoErr(t, err)

	serverRequests := []cloudservers.Server{{Id: id}}

	deleteOpts := cloudservers.DeleteOpts{
		Servers:      serverRequests,
		DeleteVolume: true,
	}

	jobResponse, err := cloudservers.Delete(client, deleteOpts).ExtractJobResponse()
	th.AssertNoErr(t, err)

	th.AssertNoErr(t, cloudservers.WaitForJobSuccess(client, 120, jobResponse.JobID))
}

func TestAccEcsV1Instance_import(t *testing.T) {
	t.Parallel()
	qts := serverQuotas(10+4, getFlavorName())
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
					"data_disks",
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
	qts := serverQuotas(4, getFlavorName())
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
	qts := serverQuotas(10+4, getFlavorName())
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

func TestAccEcsV1InstanceVolumeAttach(t *testing.T) {
	var instance cloudservers.CloudServer
	qts := serverQuotas(10+4, getFlavorName())
	t.Parallel()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1InstanceAttachVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "availability_zone", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "auto_recovery", "true"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "volumes_attached.#", "0"),
				),
			},
			{
				Config: testAccEcsV1InstanceAttachVolumeRepeat,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists(resourceInstanceV1Name, &instance),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "availability_zone", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "auto_recovery", "true"),
					resource.TestCheckResourceAttr(resourceInstanceV1Name, "volumes_attached.#", "1"),
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
  flavor   = "%s"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  data_disks {
    size = 10
    type = "SAS"
  }

  password                    = "Password@123"
  availability_zone           = "%s"
  auto_recovery               = true
  delete_disks_on_termination = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceImage, common.DataSourceSubnet, getFlavorName(), env.OS_AVAILABILITY_ZONE)

var testAccEcsV1InstanceUpdate = fmt.Sprintf(`
%s

%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_updated"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, getFlavorName(), env.OS_AVAILABILITY_ZONE)

var testAccEcsV1InstanceInvalidTypeForAZ = fmt.Sprintf(`
%s

%s


resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
`, common.DataSourceImage, common.DataSourceSubnet, getFlavorName())

var testAccEcsV1InstanceInvalidType = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
`, common.DataSourceImage, common.DataSourceSubnet, getFlavorName())

var testAccEcsV1InstanceInvalidDataDisk = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
`, common.DataSourceImage, common.DataSourceSubnet, getFlavorName())

var testAccEcsV1InstanceInvalidVPC = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
`, common.DataSourceImage, common.DataSourceSubnet, getFlavorName())

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
  flavor   = "%s"
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
`, common.DataSourceImage, getFlavorName())

var testAccEcsV1InstanceDataVolumeEncryption = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
`, common.DataSourceImage, common.DataSourceSubnet, getFlavorName(), env.OS_AVAILABILITY_ZONE, env.OS_KMS_ID)

var testAccEcsV1InstanceIp = fmt.Sprintf(`
%s

%s

%s

resource "opentelekomcloud_networking_floatingip_v2" "this" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
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
}


resource "opentelekomcloud_networking_floatingip_associate_v2" "this" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.this.address
  port_id     = opentelekomcloud_ecs_instance_v1.instance_1.nics.0.port_id
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, getFlavorName(), env.OS_AVAILABILITY_ZONE)

var testAccEcsV1InstanceAttachVolume = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_blockstorage_volume_v2" "myvol" {
  name              = "myvol"
  availability_zone = "%s"
  size              = 10
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  availability_zone = "%s"

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

}

resource "opentelekomcloud_compute_volume_attach_v2" "attached" {
  instance_id = opentelekomcloud_ecs_instance_v1.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.myvol.id
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, getFlavorName(), env.OS_AVAILABILITY_ZONE)

var testAccEcsV1InstanceAttachVolumeRepeat = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_blockstorage_volume_v2" "myvol" {
  name              = "myvol"
  availability_zone = "%s"
  size              = 10
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "%s"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  availability_zone = "%s"

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

}

resource "opentelekomcloud_compute_volume_attach_v2" "attached" {
  instance_id = opentelekomcloud_ecs_instance_v1.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.myvol.id
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, getFlavorName(), env.OS_AVAILABILITY_ZONE)
