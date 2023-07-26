package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const wafdInstanceResourceName = "opentelekomcloud_waf_dedicated_instance_v1.wafd_1"

func TestAccWafDedicatedInstanceV1_basic(t *testing.T) {
	var inst instances.Instance
	var instanceName = fmt.Sprintf("wafd_instance_%s", acctest.RandString(5))
	arch := "x86"
	flavor := "s2.large.2"
	if env.OS_REGION_NAME == "eu-ch2" {
		arch = "x86_64"
		flavor = "s3.large.2"
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedInstanceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstanceV1_basic(instanceName, arch, flavor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedInstanceV1Exists(
						wafdInstanceResourceName, &inst),
					resource.TestCheckResourceAttr(wafdInstanceResourceName, "name", instanceName),
					resource.TestCheckResourceAttr(wafdInstanceResourceName, "specification", "waf.instance.professional"),
					resource.TestCheckResourceAttr(wafdInstanceResourceName, "security_group.#", "1"),
				),
			},
			{
				Config: testAccWafDedicatedInstanceV1_update(instanceName, arch, flavor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedInstanceV1Exists(
						wafdInstanceResourceName, &inst),
					resource.TestCheckResourceAttr(wafdInstanceResourceName, "name", instanceName+"-updated"),
					resource.TestCheckResourceAttr(wafdInstanceResourceName, "specification", "waf.instance.professional"),
					resource.TestCheckResourceAttr(wafdInstanceResourceName, "security_group.#", "1"),
				),
			},
			{
				ResourceName:            wafdInstanceResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"specification"},
			},
		},
	})
}

func testAccCheckWafDedicatedInstanceV1Destroy(s *terraform.State) error {
	var client *golangsdk.ServiceClient
	var err error
	config := common.TestAccProvider.Meta().(*cfg.Config)
	if env.OS_REGION_NAME != "eu-ch2" {
		client, err = config.WafDedicatedV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Waf dedicated client: %s", err)
		}
	} else {
		client, err = config.WafDedicatedSwissV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Waf dedicated client: %s", err)
		}
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_dedicated_instance_v1" {
			continue
		}
		_, err = instances.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("waf dedicated instance (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckWafDedicatedInstanceV1Exists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var client *golangsdk.ServiceClient
		var err error
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		if env.OS_REGION_NAME != "eu-ch2" {
			client, err = config.WafDedicatedV1Client(env.OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("error creating OpenTelekomCloud Waf dedicated client: %s", err)
			}
		} else {
			client, err = config.WafDedicatedSwissV1Client(env.OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("error creating OpenTelekomCloud Waf dedicated client: %s", err)
			}
		}

		var found *instances.Instance
		found, err = instances.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		*instance = *found

		return nil
	}
}

func testAccWafDedicatedInstanceV1_basic(instanceName, arch, flavor string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_waf_dedicated_instance_v1" "wafd_1" {
    name              = "%[4]s"
    availability_zone = "%[3]s"
    specification     = "waf.instance.professional"
    flavor            = "%[6]s"
    architecture      = "%[5]s"
    vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
    subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

    security_group = [
      data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
    ]
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName, arch, flavor)
}

func testAccWafDedicatedInstanceV1_update(instanceName, arch, flavor string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_waf_dedicated_instance_v1" "wafd_1" {
    name              = "%[4]s-updated"
    availability_zone = "%[3]s"
    specification     = "waf.instance.professional"
    flavor            = "%[6]s"
    architecture      = "%[5]s"
    vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
    subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

    security_group = [
      data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
    ]
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, instanceName, arch, flavor)
}
