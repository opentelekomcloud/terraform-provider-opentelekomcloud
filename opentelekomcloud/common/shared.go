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
data "opentelekomcloud_subnet_v1" "shared_subnet"  {
  name = "%s"
}
`, env.OsSubnetName)
