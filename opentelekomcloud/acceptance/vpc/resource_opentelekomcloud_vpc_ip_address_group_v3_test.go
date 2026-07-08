package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/addressgroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const vpcIPAddressGroupResourceName = "opentelekomcloud_vpc_ip_address_group_v3.group_1"

func getVpcIPAddressGroupV3(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 Client: %s", err)
	}
	return addressgroup.Get(client, state.Primary.ID)
}

func TestAccVpcIPAddressGroupV3_basic(t *testing.T) {
	var group addressgroup.AddressGroup
	rc := common.InitResourceCheck(
		vpcIPAddressGroupResourceName,
		&group,
		getVpcIPAddressGroupV3,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIPAddressGroupV3Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(vpcIPAddressGroupResourceName, "name", "test-acc-ip-address-group-v3"),
					resource.TestCheckResourceAttr(vpcIPAddressGroupResourceName, "description", "created"),
					resource.TestCheckResourceAttr(vpcIPAddressGroupResourceName, "ip_version", "4"),
				),
			},
			{
				Config: testAccVpcIPAddressGroupV3Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(vpcIPAddressGroupResourceName, "name", "test-acc-ip-address-group-v3"),
					resource.TestCheckResourceAttr(vpcIPAddressGroupResourceName, "description", "updated"),
					resource.TestCheckResourceAttr(vpcIPAddressGroupResourceName, "ip_version", "4"),
				),
			},
			{
				ResourceName:      vpcIPAddressGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccVpcIPAddressGroupV3Basic = `
resource "opentelekomcloud_vpc_ip_address_group_v3" "group_1" {
  name        = "test-acc-ip-address-group-v3"
  description = "created"
  ip_version  = 4
}
`

var testAccVpcIPAddressGroupV3Update = `
resource "opentelekomcloud_vpc_ip_address_group_v3" "group_1" {
  name        = "test-acc-ip-address-group-v3"
  description = "updated"
  ip_version  = 4
}
`
