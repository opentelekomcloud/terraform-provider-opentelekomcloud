package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceBandwidthName = "opentelekomcloud_vpc_bandwidth_v2.band_test"

func TestBandwidthV2_basic(t *testing.T) {
	var b bandwidths.BandWidth

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
					resource.TestCheckResourceAttr(resourceBandwidthName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttrSet(resourceBandwidthName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceBandwidthName, "public_border_group"),
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
			{
				ResourceName:      resourceBandwidthName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckBandwidthExists(name string, bandwidth *bandwidths.BandWidth) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.VpcV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating VPC v1 client: %s", err)
		}
		found, err := testFindBandwidthV1(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found == nil {
			return fmt.Errorf("bandwidth not found")
		}
		*bandwidth = *found
		return nil
	}
}

func testCheckBandwidthV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	vpcClient, err := config.VpcV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating VPC v1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_bandwidth_v2" {
			continue
		}

		found, err := testFindBandwidthV1(vpcClient, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error checking bandwidth deletion: %w", err)
		}
		if found != nil {
			return fmt.Errorf("bandwidth still exists")
		}
	}

	return nil
}

func testFindBandwidthV1(client *golangsdk.ServiceClient, id string) (*bandwidths.BandWidth, error) {
	allBandwidths, err := bandwidths.List(client, bandwidths.ListOpts{})
	if err != nil {
		return nil, err
	}
	for i := range allBandwidths {
		if allBandwidths[i].ID == id {
			return &allBandwidths[i], nil
		}
	}
	return nil, nil
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
