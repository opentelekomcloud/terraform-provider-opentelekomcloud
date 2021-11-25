package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceBandwidthName = "opentelekomcloud_vpc_bandwidth_v2.band_test"

func TestBandwidthV2_basic(t *testing.T) {
	var b bandwidths.Bandwidth

	t.Parallel()
	quotas.BookOne(t, quotas.SharedBandwidth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCheckBandwidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testBandwidthV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testCheckBandwidthExists(resourceBandwidthName, &b),
					resource.TestCheckResourceAttr(resourceBandwidthName, "size", "100"),
					resource.TestCheckResourceAttr(resourceBandwidthName, "status", "NORMAL"),
				),
			},
			{
				Config: testBandwidthV2Updated,
				Check: resource.ComposeTestCheckFunc(
					testCheckBandwidthExists(resourceBandwidthName, &b),
					resource.TestCheckResourceAttr(resourceBandwidthName, "size", "50"),
					resource.TestCheckResourceAttr(resourceBandwidthName, "status", "NORMAL"),
				),
			},
		},
	})
}

func TestBandwidthV2_import(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.SharedBandwidth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCheckBandwidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testBandwidthV2Basic,
			},
			{
				ResourceName:      resourceBandwidthName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckBandwidthExists(name string, bandwidth *bandwidths.Bandwidth) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating networkingV2 client: %s", err)
		}
		found, err := bandwidths.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("bandwidth not found")
		}
		*bandwidth = *found
		return nil
	}
}

func testCheckBandwidthV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("eror creating NetworkingV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_bandwidth_v2" {
			continue
		}

		_, err := bandwidths.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("bandwidth still exists")
		}
	}

	return nil
}

const testBandwidthV2Basic = `
resource "opentelekomcloud_vpc_bandwidth_v2" "band_test" {
  name = "shared-test"
  size = 100
}
`

const testBandwidthV2Updated = `
resource "opentelekomcloud_vpc_bandwidth_v2" "band_test" {
  name = "shared-test"
  size = 50
}
`
