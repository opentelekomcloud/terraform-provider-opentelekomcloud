package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const resourceBandwidthAssociateName = "opentelekomcloud_vpc_bandwidth_associate_v2.associate"

func TestBandwidthAssociateV2_basic(t *testing.T) {
	var b bandwidths.Bandwidth

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
	var b bandwidths.Bandwidth

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
