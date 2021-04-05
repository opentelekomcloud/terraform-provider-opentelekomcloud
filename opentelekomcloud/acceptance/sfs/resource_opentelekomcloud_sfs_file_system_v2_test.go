package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs/v2/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSFSFileSystemV2_basic(t *testing.T) {
	var share shares.Share
	resourceName := "opentelekomcloud_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-test1"),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "access_to", env.OS_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "access_type", "cert"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccSFSFileSystemV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-test2"),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "2"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "access_to", env.OS_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "access_type", "cert"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_timeout(t *testing.T) {
	var share shares.Share
	resourceName := "opentelekomcloud_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_clean(t *testing.T) {
	var share shares.Share
	resourceName := "opentelekomcloud_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_clean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
				),
			},
			{
				Config: testAccSFSFileSystemV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists("opentelekomcloud_sfs_file_system_v2.sfs_1", &share),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "access_to", env.OS_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "access_type", "cert"),
				),
			},
			{
				Config: testAccSFSFileSystemV2_clean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
				),
			},
		},
	})
}

func testAccCheckSFSFileSystemV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SfsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SFSv2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sfs_file_system_v2" {
			continue
		}

		_, err := shares.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("share file still exists")
		}
	}

	return nil
}

func testAccCheckSFSFileSystemV2Exists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.SfsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud sfs client: %s", err)
		}

		found, err := shares.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("share file not found")
		}

		*share = *found

		return nil
	}
}

var testAccSFSFileSystemV2_basic = fmt.Sprintf(`
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  share_proto       = "NFS"
  size              = 1
  name              = "sfs-test1"
  availability_zone = "eu-de-01"
  access_to         = "%s"
  access_type       = "cert"
  access_level      = "rw"
  description       = "sfs_c2c_test-file"

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, env.OS_VPC_ID)

var testAccSFSFileSystemV2_update = fmt.Sprintf(`
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  share_proto       = "NFS"
  size              = 2
  name              = "sfs-test2"
  availability_zone = "eu-de-01"
  access_to         = "%s"
  access_type       = "cert"
  access_level      = "rw"
  description       = "sfs_c2c_test-file"

  tags = {
    muh = "value-update"
  }
}
`, env.OS_VPC_ID)

var testAccSFSFileSystemV2_timeout = fmt.Sprintf(`
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  share_proto  = "NFS"
  size         = 1
  name         = "sfs-test1"
  access_to    = "%s"
  access_type  = "cert"
  access_level = "rw"
  description  = "sfs_c2c_test-file"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}`, env.OS_VPC_ID)

var testAccSFSFileSystemV2_clean = fmt.Sprintf(`
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  share_proto       = "NFS"
  size              = 1
  name              = "sfs-test1"
  availability_zone = "eu-de-01"
}
`)
