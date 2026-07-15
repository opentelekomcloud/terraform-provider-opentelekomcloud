package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccDDSPublicIpAssociateV3_basic(t *testing.T) {
	nodeId := os.Getenv("OS_DDS_NODE_ID")
	if nodeId == "" {
		t.Skip("OS_DDS_NODE_ID must be set for DDS public IP associate tests")
	}

	resName := "opentelekomcloud_dds_public_ip_associate_v3.public_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDDSPublicIpAssociateV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSPublicIpAssociateV3Basic(nodeId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSPublicIpAssociateV3Exists(resName),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
				),
			},
			{
				Config: testAccDDSPublicIpAssociateV3Update(nodeId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSPublicIpAssociateV3Exists(resName),
					resource.TestCheckResourceAttrSet(resName, "public_ip"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"public_ip_id",
				},
			},
		},
	})
}

func testAccCheckDDSPublicIpAssociateV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DDSv3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dds_public_ip_associate_v3" {
			continue
		}

		opts := instances.ListInstanceOpts{}
		ddsInstances, err := instances.List(client, opts)
		if err != nil {
			return err
		}

		for _, instance := range ddsInstances.Instances {
			for _, group := range instance.Groups {
				for _, node := range group.Nodes {
					if node.Id == rs.Primary.ID && node.PublicIP != "" {
						return fmt.Errorf("node still has public IP assigned")
					}
				}
			}
		}
	}

	return nil
}

func testAccCheckDDSPublicIpAssociateV3Exists(n string) resource.TestCheckFunc {
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

		opts := instances.ListInstanceOpts{}
		ddsInstances, err := instances.List(client, opts)
		if err != nil {
			return err
		}

		found := false
		for _, instance := range ddsInstances.Instances {
			for _, group := range instance.Groups {
				for _, node := range group.Nodes {
					if node.Id == rs.Primary.ID {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("dds node not found")
		}

		return nil
	}
}

func testAccDDSPublicIpAssociateV3Basic(nodeId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_compute_floatingip_v2" "eip_1" {}

resource "opentelekomcloud_compute_floatingip_v2" "eip_2" {}

resource "opentelekomcloud_dds_public_ip_associate_v3" "public_ip" {
  node_id      = "%s"
  public_ip    = opentelekomcloud_compute_floatingip_v2.eip_1.address
  public_ip_id = opentelekomcloud_compute_floatingip_v2.eip_1.id
}
`, nodeId)
}

func testAccDDSPublicIpAssociateV3Update(nodeId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_compute_floatingip_v2" "eip_1" {}

resource "opentelekomcloud_compute_floatingip_v2" "eip_2" {}

resource "opentelekomcloud_dds_public_ip_associate_v3" "public_ip" {
  node_id      = "%s"
  public_ip    = opentelekomcloud_compute_floatingip_v2.eip_2.address
  public_ip_id = opentelekomcloud_compute_floatingip_v2.eip_2.id
}
`, nodeId)
}
