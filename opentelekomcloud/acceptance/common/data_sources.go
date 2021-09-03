package common

import (
	"fmt"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

// DataSourceSubnet can be referred as `data.opentelekomcloud_vpc_subnet_v1.shared_subnet`
var DataSourceSubnet = fmt.Sprintf(`
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet"  {
  name = "%s"
}
`, env.OsSubnetName)

// DataSourceExtNetwork can be referred as `data.opentelekomcloud_networking_network_v2.ext_network`
var DataSourceExtNetwork = fmt.Sprintf(`
data "opentelekomcloud_networking_network_v2" "ext_network" {
  name = "%s"
}
`, env.OsExtNetworkName)

// DataSourceImage can be referred as `data.opentelekomcloud_images_image_v2.latest_image`
var DataSourceImage = fmt.Sprintf(`
data "opentelekomcloud_images_image_v2" "latest_image" {
  name        = "%s"
  most_recent = true
}
`, env.OsImageName)

// DataSourceKMSKey can be referred as `data.opentelekomcloud_kms_key_v1.default_key`
var DataSourceKMSKey = fmt.Sprintf(`
data "opentelekomcloud_kms_key_v1" "default_key" {
  key_alias = "%s"
}
`, env.OsKmsName)

// DataSourceSecGroupDefault can be referred as `data.opentelekomcloud_networking_secgroup_v2.default_secgroup`
const DataSourceSecGroupDefault = `
data "opentelekomcloud_networking_secgroup_v2" "default_secgroup" {
  name = "default"
}
`
