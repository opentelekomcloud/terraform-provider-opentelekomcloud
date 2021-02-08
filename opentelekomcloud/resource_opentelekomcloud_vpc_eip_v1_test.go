package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"
)

func TestAccVpcV1EIP_basic(t *testing.T) {
	var eip eips.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("opentelekomcloud_vpc_eip_v1.eip_1", &eip),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_eip_v1.eip_1", "bandwidth.0.name", "acc-band"),
				),
			},
			{
				Config: testAccVpcV1EIP_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("opentelekomcloud_vpc_eip_v1.eip_1", &eip),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_eip_v1.eip_1", "bandwidth.0.name", "acc-band-update"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_eip_v1.eip_1", "bandwidth.0.size", "25"),
				),
			},
		},
	})
}

func TestAccVpcV1EIP_timeout(t *testing.T) {
	var eip eips.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_timeouts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("opentelekomcloud_vpc_eip_v1.eip_1", &eip),
				),
			},
		},
	})
}

func testAccCheckVpcV1EIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("eror creating EIP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_eip_v1" {
			continue
		}

		_, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EIP still exists")
		}
	}

	return nil
}

func testAccCheckVpcV1EIPExists(n string, eip *eips.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating networkingV1 client: %s", err)
		}

		found, err := eips.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EIP not found")
		}

		eip = &found

		return nil
	}
}

const testAccVpcV1EIP_basic = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "acc-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
  tags = {
    this = "is"
    our  = "tags"
  }
}
`

const testAccVpcV1EIP_update = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "acc-band-update"
    size = 25
    share_type = "PER"
    charge_mode = "traffic"
  }
  tags = {
    muh = "value1"
    kuh = "value2"
  }
}
`

const testAccVpcV1EIP_timeouts = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "acc-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
