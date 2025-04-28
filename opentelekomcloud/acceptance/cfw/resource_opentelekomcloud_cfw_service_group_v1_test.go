package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/servicegroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const serviceGroupResourceName = "opentelekomcloud_cfw_service_group_v1.group_1"

func getServiceGroupFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CfwV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CFW client: %s", err)
	}
	return servicegroup.GetServiceGroup(client, state.Primary.ID)
}

func TestAccCFWServiceGroupV1_basic(t *testing.T) {
	var group servicegroup.ServiceSetDetailResponseDto
	rc := common.InitResourceCheck(
		serviceGroupResourceName,
		&group,
		getServiceGroupFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCFWServiceGroupV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(serviceGroupResourceName, "name", "test-acc-tf-service-group"),
					resource.TestCheckResourceAttrSet(serviceGroupResourceName, "id"),
				),
			},
			{
				Config: testAccCFWServiceGroupV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(serviceGroupResourceName, "name", "test-acc-tf-service-group-updated"),
				),
			},
			{
				ResourceName:      serviceGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"object_id",
				},
			},
		},
	})
}

var testAccCFWServiceGroupV1Basic = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_service_group_v1" "group_1" {
  object_id    = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name         = "test-acc-tf-service-group"
}
`

var testAccCFWServiceGroupV1Update = `
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}

resource "opentelekomcloud_cfw_service_group_v1" "group_1" {
  object_id    = opentelekomcloud_cfw_firewall_v1.firewall_1.protect_objects.0.object_id
  name         = "test-acc-tf-service-group-updated"
}
`
