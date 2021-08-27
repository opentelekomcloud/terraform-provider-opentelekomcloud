package common

import (
	"fmt"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

var DataSourceVPC = fmt.Sprintf(`
data "opentelekomcloud_vpc_v1" "shared_vpc"  {
  name = "%s"
}
`, env.OsRouterName)

var DataSourceSubnet = fmt.Sprintf(`
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet"  {
  name = "%s"
}
`, env.OsSubnetName)

var DataSourceExtNetwork = fmt.Sprintf(`
data "opentelekomcloud_networking_network_v2" "ext_network" {
  name = "%s"
}
`, env.OsExtNetworkName)

var DataSourceImage = fmt.Sprintf(`
data "opentelekomcloud_images_image_v2" "latest_image" {
  name        = "%s"
  most_recent = true
}
`, env.OsImageName)

var DataSourceSecGroupDefault = fmt.Sprintf(`
data "opentelekomcloud_networking_secgroup_v2" "default_secgroup" {
  name = "default"
}
`)
