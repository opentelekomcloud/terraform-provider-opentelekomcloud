package shared

import (
	"net"
	"os"
	"sync"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

var sharedSubnet *subnets.Subnet
var subnetOnce sync.Once
var SubnetNet *net.IPNet

func getSharedSubnet(t *testing.T) *subnets.Subnet {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Shared subnet can only be used in acceptance tests")
	}

	subnetOnce.Do(func() {
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
		th.AssertNoErr(t, err)

		subnetSlice, err := subnets.List(client, subnets.ListOpts{Name: env.OsSubnetName})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, len(subnetSlice))
		sharedSubnet = &subnetSlice[0]

		_, SubnetNet, _ = net.ParseCIDR(sharedSubnet.CIDR)
	})
	return sharedSubnet
}
