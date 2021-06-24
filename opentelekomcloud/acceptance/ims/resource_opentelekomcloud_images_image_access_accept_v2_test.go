package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccImagesImageAccessAcceptV2_basic(t *testing.T) {
	var member members.Member

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessAcceptV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessAcceptV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists("openstack_images_image_access_accept_v2.image_access_accept_1", &member),
					resource.TestCheckResourceAttrPtr(
						"openstack_images_image_access_accept_v2.image_access_accept_1", "status", &member.Status),
					resource.TestCheckResourceAttr(
						"openstack_images_image_access_accept_v2.image_access_accept_1", "status", "accepted"),
				),
			},
			{
				Config: testAccImagesImageAccessAcceptV2Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists("openstack_images_image_access_accept_v2.image_access_accept_1", &member),
					resource.TestCheckResourceAttrPtr(
						"openstack_images_image_access_accept_v2.image_access_accept_1", "status", &member.Status),
					resource.TestCheckResourceAttr(
						"openstack_images_image_access_accept_v2.image_access_accept_1", "status", "rejected"),
				),
			},
		},
	})
}

func testAccCheckImagesImageAccessAcceptV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ImageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenStack IMSv2: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_images_image_access_accept_v2" {
			continue
		}

		imageID, memberID, err := resourceImagesImageAccessV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = members.Get(client, imageID, memberID).Extract()
		if err == nil {
			return fmt.Errorf("image membership still exists")
		}
	}

	return nil
}

const testAccImagesImageAccessAcceptV2 = `
data "openstack_identity_auth_scope_v3" "scope" {
  name = "scope"
}
resource "openstack_images_image_v2" "image_1" {
  name   = "CirrOS-tf_1"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  container_format = "bare"
  disk_format = "qcow2"
  visibility = "shared"
  timeouts {
    create = "10m"
  }
}
resource "openstack_images_image_access_v2" "image_access_1" {
  image_id  = "${openstack_images_image_v2.image_1.id}"
  member_id = "${data.openstack_identity_auth_scope_v3.scope.project_id}"
}
`

func testAccImagesImageAccessAcceptV2Basic() string {
	return fmt.Sprintf(`
%s
resource "openstack_images_image_access_accept_v2" "image_access_accept_1" {
  image_id  = "${openstack_images_image_access_v2.image_access_1.image_id}"
  status    = "accepted"
}
`, testAccImagesImageAccessAcceptV2)
}

func testAccImagesImageAccessAcceptV2Update() string {
	return fmt.Sprintf(`
%s
resource "openstack_images_image_access_accept_v2" "image_access_accept_1" {
  image_id  = "${openstack_images_image_access_v2.image_access_1.image_id}"
  status    = "rejected"
}
`, testAccImagesImageAccessAcceptV2)
}
