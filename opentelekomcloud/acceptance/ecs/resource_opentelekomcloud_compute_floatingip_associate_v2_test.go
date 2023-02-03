package acceptance

import (
	"fmt"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	vpc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/vpc"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ecs"
)

const (
	resourceNetworkingFloatingIpName = "opentelekomcloud_networking_floatingip_v2.fip_1"
	resourceFloatingIpAssociateName  = "opentelekomcloud_compute_floatingip_associate_v2.fip_1"
)

func TestAccComputeV2FloatingIPAssociate_basic(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP
	t.Parallel()
	qts := simpleServerWithIPQuotas(1)
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceInstanceV2Name, &instance),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNetworkingFloatingIpName, &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_importBasic(t *testing.T) {
	t.Parallel()
	qts := simpleServerWithIPQuotas(1)
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateBasic,
			},
			{
				ResourceName:      resourceFloatingIpAssociateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_fixedIP(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP
	t.Parallel()
	qts := simpleServerWithIPQuotas(1)
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateFixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceInstanceV2Name, &instance),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNetworkingFloatingIpName, &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_attachToFirstNetwork(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP
	t.Parallel()
	qts := simpleServerWithIPQuotas(1)
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateAttachToFirstNetwork,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceInstanceV2Name, &instance),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNetworkingFloatingIpName, &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_attachNew(t *testing.T) {
	var instance servers.Server
	var fip1 floatingips.FloatingIP
	var fip2 floatingips.FloatingIP
	resourceNwFloatingIp2Name := "opentelekomcloud_networking_floatingip_v2.fip_2"
	t.Parallel()
	qts := simpleServerWithIPQuotas(2)
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateAttachNew1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceInstanceV2Name, &instance),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNetworkingFloatingIpName, &fip1),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNwFloatingIp2Name, &fip2),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip1, &instance, 1),
				),
			},
			{
				Config: testAccComputeV2FloatingIPAssociateAttachNew2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceInstanceV2Name, &instance),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNetworkingFloatingIpName, &fip1),
					vpc.TestAccCheckNetworkingV2FloatingIPExists(resourceNwFloatingIp2Name, &fip2),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip2, &instance, 1),
				),
			},
		},
	})
}

func testAccCheckComputeV2FloatingIPAssociateDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_floatingip_associate_v2" {
			continue
		}

		floatingIP, instanceId, _, err := ecs.ParseComputeFloatingIPAssociateId(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := servers.Get(client, instanceId).Extract()
		if err != nil {
			// If the error is a 404, then the instance does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// But if the instance still exists, then walk through its known addresses
		// and see if there's a floating IP.
		for _, networkAddresses := range instance.Addresses {
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" || address["OS-EXT-IPS:type"] == "fixed" {
					return fmt.Errorf("floating IP %s is still attached to instance %s", floatingIP, instanceId)
				}
			}
		}
	}

	return nil
}

func testAccCheckComputeV2FloatingIPAssociateAssociated(fip *floatingips.FloatingIP, instance *servers.Server, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ComputeV2 client: %w", err)
		}

		newInstance, err := servers.Get(client, instance.ID).Extract()
		if err != nil {
			return err
		}

		// Walk through the instance's addresses and find the match
		i := 0
		for _, networkAddresses := range newInstance.Addresses {
			i += 1
			if i != n {
				continue
			}
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if (address["OS-EXT-IPS:type"] == "floating" && address["addr"] == fip.FloatingIP) ||
					(address["OS-EXT-IPS:type"] == "fixed" && address["addr"] == fip.FixedIP) {
					return nil
				}
			}
		}
		return fmt.Errorf("floating IP %s was not attached to instance %s", fip.FloatingIP, instance.ID)
	}
}

var testAccComputeV2FloatingIPAssociateBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = "Standard_Debian_10_latest"
  flavor_name	  = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
}
`, common.DataSourceSubnet, getFlavorName())

var testAccComputeV2FloatingIPAssociateFixedIP = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = "Standard_Debian_10_latest"
  flavor_name	  = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  fixed_ip    = opentelekomcloud_compute_instance_v2.instance_1.access_ip_v4
}
`, common.DataSourceSubnet, getFlavorName())

var testAccComputeV2FloatingIPAssociateAttachToFirstNetwork = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = "Standard_Debian_10_latest"
  flavor_name	  = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  fixed_ip    = opentelekomcloud_compute_instance_v2.instance_1.network.0.fixed_ip_v4
}
`, common.DataSourceSubnet, getFlavorName())

var testAccComputeV2FloatingIPAssociateAttachNew1 = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = "Standard_Debian_10_latest"
  flavor_name	  = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
}
`, common.DataSourceSubnet, getFlavorName())

var testAccComputeV2FloatingIPAssociateAttachNew2 = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
  image_name      = "Standard_Debian_10_latest"
  flavor_name	  = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_2.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
}
`, common.DataSourceSubnet, getFlavorName())
