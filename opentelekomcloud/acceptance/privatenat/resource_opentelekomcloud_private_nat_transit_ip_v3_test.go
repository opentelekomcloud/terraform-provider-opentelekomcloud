package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/transitip"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const transitIpResourceName = "opentelekomcloud_private_nat_transit_ip_v3.transit_ip_1"

func getTransitIp(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 Client: %s", err)
	}
	getResp, err := transitip.Get(client, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching Transit IP: %s", err)
	}
	return getResp.TransitIp, nil
}

func TestAccPrivateTransitIpV3_basic(t *testing.T) {
	var gateway transitip.TransitIP
	rc := common.InitResourceCheck(
		transitIpResourceName,
		&gateway,
		getTransitIp,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitIpV3Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(transitIpResourceName, "tags.kuh", "muh"),
				),
			},
			{
				ResourceName:      transitIpResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccPrivateTransitIpV3Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  virsubnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  tags = {
    kuh = "muh"
  }
}
`, common.DataSourceSubnet)
