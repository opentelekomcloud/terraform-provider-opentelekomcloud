package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpcep/v1/endpoints"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vpcep"
)

const resourceEndpointName = "opentelekomcloud_vpcep_endpoint_v1.endpoint"

func TestEndpoint_basic(t *testing.T) {
	var ep endpoints.Endpoint
	name := tools.RandomString("tf-test-ep-", 4)
	t.Parallel()
	quotas.BookMany(t, endpointQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testEndpointBasic(name),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists(resourceEndpointName, &ep),
					resource.TestCheckResourceAttr(resourceEndpointName, "port_ip", "192.168.0.12"),
					resource.TestCheckResourceAttr(resourceEndpointName, "tags.fizz", "buzz"),
					resource.TestCheckResourceAttr(resourceEndpointName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceEndpointName, "dns_names.#", "1"),
					resource.TestCheckResourceAttr(resourceEndpointName, "service_name", name),
				),
			},
		},
	})
}

func TestEndpoint_import(t *testing.T) {
	name := tools.RandomString("tf-test-ep-", 4)
	t.Parallel()
	quotas.BookMany(t, endpointQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testEndpointImport(name),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceEndpointName,
			},
		},
	})
}

func testEndpointBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpcep_endpoint_v1" "endpoint" {
  service_id = opentelekomcloud_vpcep_service_v1.service.id
  vpc_id     = opentelekomcloud_vpcep_service_v1.service.vpc_id
  subnet_id  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  port_ip    = "192.168.0.12"
  enable_dns = true

  tags = {
    "fizz" : "buzz"
  }
}
`, testServiceBasic(name)) // without acceptance required
}

func testEndpointImport(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpcep_endpoint_v1" "endpoint" {
  service_id = opentelekomcloud_vpcep_service_v1.service.id
  vpc_id     = opentelekomcloud_vpcep_service_v1.service.vpc_id
  subnet_id  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  port_ip    = "192.168.0.14"
  enable_dns = true

  tags = {
    "fizz" : "buzz"
  }
}
`, testServiceBasic(name)) // without acceptance required
}

func checkEndpointDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.VpcEpV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(vpcep.ErrClientCreate, err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpcep_service_v1" {
			continue
		}
		svc, err := endpoints.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("error getting service state: %w", err)
		}
		return fmt.Errorf("VPC Endpoint %s still exists", svc.ServiceName)
	}
	return nil
}

func checkEndpointExists(name string, ep *endpoints.Endpoint) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.VpcEpV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(vpcep.ErrClientCreate, err)
		}
		found, err := endpoints.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting endpoint: %w", err)
		}
		*ep = *found
		return nil
	}
}
