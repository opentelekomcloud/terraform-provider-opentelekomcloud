package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/servergroups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccComputeV2ServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("opentelekomcloud_compute_servergroup_v2.sg_1", &sg),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroup_affinity,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("opentelekomcloud_compute_servergroup_v2.sg_1", &sg),
					testAccCheckComputeV2InstanceExists("opentelekomcloud_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceInServerGroup(&instance, &sg),
				),
			},
		},
	})
}

func testAccCheckComputeV2ServerGroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	computeClient, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_servergroup_v2" {
			continue
		}

		_, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("serverGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2ServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		computeClient, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
		}

		found, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("serverGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceInServerGroup(instance *servers.Server, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

const testAccComputeV2ServerGroup_basic = `
resource "opentelekomcloud_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}
`

var testAccComputeV2ServerGroup_affinity = fmt.Sprintf(`
resource "opentelekomcloud_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  scheduler_hints {
    group = opentelekomcloud_compute_servergroup_v2.sg_1.id
  }
}
`, env.OS_NETWORK_ID)
