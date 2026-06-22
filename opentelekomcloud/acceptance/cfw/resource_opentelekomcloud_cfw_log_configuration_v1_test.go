package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/logs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const logConfigurationResourceName = "opentelekomcloud_cfw_log_configuration_v1.log_config_1"

func getLogConfigurationFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	queryOpts := logs.QueryParameters{
		FwInstanceID: state.Primary.ID,
	}
	return logs.GetLogConfig(client, queryOpts)
}

func TestAccCFWLogConfigurationV1_basic(t *testing.T) {
	var config logs.LogConfig
	rc := common.InitResourceCheck(
		logConfigurationResourceName,
		&config,
		getLogConfigurationFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWLogConfigurationV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(logConfigurationResourceName, "lts_enable", "1"),
					resource.TestCheckResourceAttrSet(logConfigurationResourceName, "lts_log_group_id"),
				),
			},
			{
				Config: testAccCFWLogConfigurationV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(logConfigurationResourceName, "lts_enable", "0"),
					resource.TestCheckResourceAttrSet(logConfigurationResourceName, "lts_log_group_id"),
				),
			},
		},
	})
}

var testAccCFWLogConfigurationV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_lts_group_v2" "group_1" {
  group_name  = "test-acc-tf-group"
  ttl_in_days = 1
}

resource "opentelekomcloud_cfw_log_configuration_v1" "log_config_1" {
  fw_instance_id   = opentelekomcloud_cfw_firewall_v1.firewall_1.id
  lts_enable       = 1
  lts_log_group_id = opentelekomcloud_lts_group_v2.group_1.id
}
`

var testAccCFWLogConfigurationV1Update = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_lts_group_v2" "group_1" {
  group_name  = "test-acc-tf-group"
  ttl_in_days = 1
}

resource "opentelekomcloud_cfw_log_configuration_v1" "log_config_1" {
  fw_instance_id   = opentelekomcloud_cfw_firewall_v1.firewall_1.id
  lts_enable       = 0
  lts_log_group_id = opentelekomcloud_lts_group_v2.group_1.id
}
`
