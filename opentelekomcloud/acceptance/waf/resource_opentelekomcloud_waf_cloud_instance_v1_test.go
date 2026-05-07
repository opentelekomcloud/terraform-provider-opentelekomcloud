package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/cloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const wafCloudInstanceResourceName = "opentelekomcloud_waf_cloud_instance_v1.cloud_1"
const wafCloudPostPaidDomainResourceType = "hws.resource.type.waf.payperusedomain"

func getWafCloudInstanceResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.WafDedicatedV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud WAF dedicated client: %w", err)
	}

	return getWafCloudInstanceResource(client, state.Primary.ID)
}

func TestAccWafCloudInstanceV1_basic(t *testing.T) {
	var cloudResource cloud.ResourceResponse
	rc := common.InitResourceCheck(wafCloudInstanceResourceName, &cloudResource, getWafCloudInstanceResourceFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafCloudInstanceV1Basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(wafCloudInstanceResourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(wafCloudInstanceResourceName, "website", "dt"),
					resource.TestCheckResourceAttrSet(wafCloudInstanceResourceName, "resource_spec_code"),
					resource.TestCheckResourceAttrSet(wafCloudInstanceResourceName, "status"),
				),
			},
			{
				ResourceName:            wafCloudInstanceResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"website", "charging_mode"},
			},
		},
	})
}

func getWafCloudInstanceResource(client *golangsdk.ServiceClient, id string) (*cloud.ResourceResponse, error) {
	subscription, err := cloud.Get(client)
	if err != nil {
		return nil, err
	}

	for i := range subscription.Resources {
		r := &subscription.Resources[i]
		if r.ResourceType == wafCloudPostPaidDomainResourceType && r.ID == id {
			return r, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("the cloud WAF (%s) does not exist", id)),
		},
	}
}

func testAccWafCloudInstanceV1Basic() string {
	return `
resource "opentelekomcloud_waf_cloud_instance_v1" "cloud_1" {
  charging_mode = "postPaid"
  website       = "dt"
}
`
}
