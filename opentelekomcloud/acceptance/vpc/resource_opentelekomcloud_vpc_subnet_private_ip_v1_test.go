package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/privateips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVpcSubnetPrivateIPV1Name = "opentelekomcloud_vpc_subnet_private_ip_v1.test"

func TestAccVpcSubnetPrivateIPV1_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, vpcSubnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetPrivateIPV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetPrivateIPV1Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceVpcSubnetPrivateIPV1Name, "ip_address"),
					resource.TestCheckResourceAttr(resourceVpcSubnetPrivateIPV1Name, "status", "DOWN"),
					resource.TestCheckResourceAttrSet(resourceVpcSubnetPrivateIPV1Name, "tenant_id"),
					resource.TestCheckResourceAttr(resourceVpcSubnetPrivateIPV1Name, "device_owner", ""),
				),
			},
			{
				ResourceName:      resourceVpcSubnetPrivateIPV1Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVpcSubnetPrivateIPV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.VpcV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_subnet_private_ip_v1" {
			continue
		}
		_, err := privateips.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC subnet private IP %s still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

const testAccVpcSubnetPrivateIPV1Basic = `
resource "opentelekomcloud_vpc_v1" "private_ip" {
  name = "private-ip-vpc"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "private_ip" {
  name       = "private-ip-subnet"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.private_ip.id
}

resource "opentelekomcloud_vpc_subnet_private_ip_v1" "test" {
  subnet_id = opentelekomcloud_vpc_subnet_v1.private_ip.network_id
}
`
