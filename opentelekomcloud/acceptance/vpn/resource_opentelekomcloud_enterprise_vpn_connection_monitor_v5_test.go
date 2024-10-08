package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cm "github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/connection-monitoring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getConnectionMonitorResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.EvpnV5Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud EVPN v5 client: %s", err)
	}
	return cm.Get(client, state.Primary.ID)
}

func TestAccConnectionHealthCheck_basic(t *testing.T) {
	var m cm.Monitor
	name := fmt.Sprintf("evpn_acc_cm_%s", acctest.RandString(5))
	rName := "opentelekomcloud_enterprise_vpn_connection_monitor_v5.cm_1"
	ipAddress := "172.16.1.4"
	rc := common.InitResourceCheck(
		rName,
		&m,
		getConnectionMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEvpnConnectionMonitor_basic(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "connection_id",
						"opentelekomcloud_enterprise_vpn_connection_v5.conn", "id"),
					resource.TestCheckResourceAttrSet(rName, "destination_ip"),
					resource.TestCheckResourceAttrSet(rName, "source_ip"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testEvpnConnectionMonitor_basic(name, ip string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_enterprise_vpn_connection_monitor_v5" "cm_1" {
  connection_id = opentelekomcloud_enterprise_vpn_connection_v5.conn.id
}
`, testEvpnConnection_basic(name, ip))
}
