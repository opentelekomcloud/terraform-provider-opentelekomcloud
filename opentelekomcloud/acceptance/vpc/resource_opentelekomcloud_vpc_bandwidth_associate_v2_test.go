package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const resourceBandwidthAssociateName = "opentelekomcloud_vpc_bandwidth_associate_v2.associate"

func TestBandwidthAssociateV2_basic(t *testing.T) {
	var b bandwidths.BandWidth

	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.SharedBandwidth, Count: 1},
		{Q: quotas.FloatingIP, Count: 2},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCheckBandwidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testBandwidthAssociateV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testCheckBandwidthExists(resourceBandwidthAssociateName, &b),
					resource.TestCheckResourceAttr(resourceBandwidthAssociateName, "floating_ips.#", "1"),
				),
			},
			{
				Config: testBandwidthAssociateV2Updated,
				Check: resource.ComposeTestCheckFunc(
					testCheckBandwidthExists(resourceBandwidthAssociateName, &b),
					resource.TestCheckResourceAttr(resourceBandwidthAssociateName, "floating_ips.#", "1"),
				),
			},
		},
	})
}

func TestBandwidthAssociateV2_EIPv1(t *testing.T) {
	var b bandwidths.BandWidth

	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.SharedBandwidth, Count: 1},
		{Q: quotas.FloatingIP, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCheckBandwidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testBandwidthAssociateV2EipV1,
				Check: resource.ComposeTestCheckFunc(
					testCheckBandwidthExists(resourceBandwidthAssociateName, &b),
					resource.TestCheckResourceAttr(resourceBandwidthAssociateName, "floating_ips.#", "1"),
				),
				ExpectNonEmptyPlan: true, // opentelekomcloud_vpc_eip_v1 bandwidth is updated
			},
		},
	})
}

func TestBandwidthAssociateV2_import(t *testing.T) {
	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.SharedBandwidth, Count: 1},
		{Q: quotas.FloatingIP, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCheckBandwidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testBandwidthAssociateV2Basic,
			},
			{
				ResourceName:            resourceBandwidthAssociateName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_charge_mode", "backup_size"},
			},
		},
	})
}

func TestBandwidthAssociateV2_EcsIpv6PortID(t *testing.T) {
	networkID := os.Getenv("OS_IPV6_ENABLED_NETWORK_ID")
	if networkID == "" {
		t.Skip("TestBandwidthAssociateV2_EcsIpv6PortID requires OS_IPV6_ENABLED_NETWORK_ID to continue.")
	}

	var b bandwidths.BandWidth

	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.SharedBandwidth, Count: 1},
		{Q: quotas.FloatingIP, Count: 1},
		{Q: quotas.Server, Count: 1},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCheckBandwidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testBandwidthAssociateV2EcsIpv6PortID(networkID),
				Check: resource.ComposeTestCheckFunc(
					testCheckBandwidthExists(resourceBandwidthAssociateName, &b),
					resource.TestCheckResourceAttrSet("opentelekomcloud_ecs_instance_v1.instance_1", "nics.0.port_id"),
					resource.TestCheckResourceAttrPair(
						"opentelekomcloud_vpc_eip_v1.eip", "publicip.0.port_id",
						"opentelekomcloud_ecs_instance_v1.instance_1", "nics.0.port_id",
					),
				),
			},
		},
	})
}

const testBandwidthAssociateV2Basic = `
resource "opentelekomcloud_networking_floatingip_v2" "ip" {}

resource "opentelekomcloud_vpc_bandwidth_v2" "band_test" {
  name = "shared-test-associate"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  bandwidth    = opentelekomcloud_vpc_bandwidth_v2.band_test.id
  floating_ips = [opentelekomcloud_networking_floatingip_v2.ip.id]
}
`

const testBandwidthAssociateV2Updated = `
resource "opentelekomcloud_networking_floatingip_v2" "ip2" {}

resource "opentelekomcloud_vpc_bandwidth_v2" "band_test" {
  name = "shared-test-associate"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  bandwidth    = opentelekomcloud_vpc_bandwidth_v2.band_test.id
  floating_ips = [opentelekomcloud_networking_floatingip_v2.ip2.id]
}
`

const testBandwidthAssociateV2EipV1 = `
resource "opentelekomcloud_vpc_eip_v1" "eip" {
  bandwidth {
    name       = "tmp-band"
    share_type = "PER"
    size       = 10
  }
  publicip {
    type = "5_bgp"
  }
}

resource "opentelekomcloud_vpc_bandwidth_v2" "band_test" {
  name = "shared-test-associate"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  bandwidth    = opentelekomcloud_vpc_bandwidth_v2.band_test.id
  floating_ips = [opentelekomcloud_vpc_eip_v1.eip.id]
}
`

func testBandwidthAssociateV2EcsIpv6PortID(networkID string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_vpc_subnet_v1" "ipv6_subnet" {
  name = "%s"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "ecs-ipv6-bandwidth-associate"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s3.xlarge.4"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.ipv6_subnet.vpc_id

  nics {
    network_id  = data.opentelekomcloud_vpc_subnet_v1.ipv6_subnet.network_id
    ipv6_enable = true
  }

  data_disks {
    size = 10
    type = "SSD"
  }

  password                    = "Password@123"
  availability_zone           = "%s"
  delete_disks_on_termination = true
}

resource "opentelekomcloud_vpc_eip_v1" "eip" {
  publicip {
    type    = "5_bgp"
    port_id = opentelekomcloud_ecs_instance_v1.instance_1.nics.0.port_id
  }

  bandwidth {
    name       = "tmp-band"
    share_type = "PER"
    size       = 10
  }
}

resource "opentelekomcloud_vpc_bandwidth_v2" "band_test" {
  name = "shared-test-associate-ipv6-port"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  depends_on = [opentelekomcloud_vpc_eip_v1.eip]

  bandwidth    = opentelekomcloud_vpc_bandwidth_v2.band_test.id
  floating_ips = [opentelekomcloud_ecs_instance_v1.instance_1.nics.0.port_id]
}
`, common.DataSourceImage, networkID, env.OS_AVAILABILITY_ZONE)
}
