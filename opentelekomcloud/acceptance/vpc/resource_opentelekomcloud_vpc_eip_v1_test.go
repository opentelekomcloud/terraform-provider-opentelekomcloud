package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/publicips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCEIPName = "opentelekomcloud_vpc_eip_v1.eip_1"

func TestAccVpcV1EIP_basic(t *testing.T) {
	var eip publicips.PublicIP
	t.Parallel()
	quotas.BookOne(t, quotas.FloatingIP)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIPBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceVPCEIPName, &eip),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "bandwidth.0.name", "acc-band"),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "publicip.0.name", "my_ip"),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "publicip.0.ip_version", "4"),
					resource.TestCheckResourceAttrSet(resourceVPCEIPName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceVPCEIPName, "public_border_group"),
				),
			},
			{
				Config: testAccVpcV1EIPUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceVPCEIPName, &eip),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "bandwidth.0.name", "acc-band-update"),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "publicip.0.name", "my_ip"),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "bandwidth.0.size", "25"),
				),
			},
		},
	})
}

func TestAccVpcV1EIP_UnAssing(t *testing.T) {
	var eip publicips.PublicIP
	t.Parallel()
	quotas.BookOne(t, quotas.FloatingIP)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceVPCEIPName, &eip),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "publicip.#", "1"),
				),
			},
			{
				Config: testAccVpcV1EIPUnassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceVPCEIPName, &eip),
					resource.TestCheckResourceAttr(resourceVPCEIPName, "publicip.0.port_id", ""),
				),
			},
		},
	})
}

func TestAccVpcV1EIP_timeout(t *testing.T) {
	var eip publicips.PublicIP
	t.Parallel()
	quotas.BookOne(t, quotas.FloatingIP)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIPTimeouts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceVPCEIPName, &eip),
				),
			},
		},
	})
}

func TestAccVpcV1EIP_shared(t *testing.T) {
	var eip publicips.PublicIP
	t.Parallel()
	quotas.BookOne(t, quotas.FloatingIP)
	const resourceVPCEIPSharedName = "opentelekomcloud_vpc_eip_v1.shared_1"
	name := fmt.Sprintf("eip-shared-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIPShare(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceVPCEIPSharedName, &eip),
					resource.TestCheckResourceAttr(resourceVPCEIPSharedName, "bandwidth.0.size", "50"),
					resource.TestCheckResourceAttr(resourceVPCEIPSharedName, "publicip.0.type", "5_bgp"),
				),
			},
		},
	})
}

func testAccCheckVpcV1EIPDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.VpcV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating VPC v1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_eip_v1" {
			continue
		}

		_, err := publicips.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("EIP still exists")
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return fmt.Errorf("error checking EIP deletion: %w", err)
		}
	}

	return nil
}

func testAccCheckVpcV1EIPExists(n string, eip *publicips.PublicIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.VpcV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating VPC v1 client: %s", err)
		}

		found, err := publicips.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EIP not found")
		}

		*eip = *found

		return nil
	}
}

const testAccVpcV1EIPBasic = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
    name = "my_ip"
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

const testAccVpcV1EIPUpdate = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
    name = "my_ip"
  }
  bandwidth {
    name        = "acc-band-update"
    size        = 25
    share_type  = "PER"
    charge_mode = "traffic"
  }
  tags = {
    muh = "value1"
    kuh = "value2"
  }
}
`

const testAccVpcV1EIPTimeouts = `
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

var testAccVpcV1EIP = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type    = "5_bgp"
    port_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.vip_port_id
  }

  bandwidth {
    name       = "test-bandwidth-acc"
    size       = 100
    share_type = "PER"
  }

}
`, common.DataSourceSubnet)

var testAccVpcV1EIPUnassign = fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type    = "5_bgp"
    port_id = null
  }

  unbind_port = true

  bandwidth {
    name       = "test-bandwidth-acc"
    size       = 1000
    share_type = "PER"
  }

}
`, common.DataSourceSubnet)

func testAccVpcV1EIPShare(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_bandwidth_v2" "test" {
  name = "%s"
  size = 50
}

resource "opentelekomcloud_vpc_eip_v1" "shared_1" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = opentelekomcloud_vpc_bandwidth_v2.test.id
  }
}
`, rName)
}
