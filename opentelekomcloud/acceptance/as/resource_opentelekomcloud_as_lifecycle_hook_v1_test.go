package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	lifecyclehooks "github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/lifecycle_hooks"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const asLifecycleHookName = "opentelekomcloud_as_lifecycle_hooks_v1.as_lifecycle_hook"

func getASLifecycleHookResourceFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Autoscaling V1 Client: %s", err)
	}
	return lifecyclehooks.Get(client, state.Primary.Attributes["scaling_group_id"], state.Primary.ID)
}

func TestAccASV1LifecycleHook_basic(t *testing.T) {

	var asLifecycleHook lifecyclehooks.LifecycleHook
	rc := common.InitResourceCheck(
		asLifecycleHookName,
		&asLifecycleHook,
		getASLifecycleHookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.ASGroup, Count: 1},
				{Q: quotas.ASConfiguration, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASV1LifecycleHookBasic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(asLifecycleHookName, "scaling_lifecycle_hook_name", "as_lifecycle_hook_v1"),
				),
			},
			{
				Config: testAccASV1LifecycleHookUpdate,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(asLifecycleHookName, "default_timeout", "4800"),
				),
			},
		},
	})
}

var testAccASV1LifecycleHookBasic = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "as_config_pol_v1"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name = "%s"
  }
}

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name       = "as_group_pol_v1"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  delete_instances         = "yes"
  delete_publicip          = true
  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
}

resource "opentelekomcloud_as_lifecycle_hooks_v1" "as_lifecycle_hook"{
  scaling_lifecycle_hook_name = "as_lifecycle_hook_v1"
  scaling_group_id    = opentelekomcloud_as_group_v1.as_group.id
  scaling_lifecycle_hook_type = "INSTANCE_TERMINATING"
  default_result = "ABANDON"
  default_timeout = 3600
  notification_topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  notification_metadata = "This is a generic message"
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

var testAccASV1LifecycleHookUpdate = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "as_config_pol_v1"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name = "%s"
  }
}

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name       = "as_group_pol_v1"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  delete_instances         = "yes"
  delete_publicip          = true
  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
}

resource "opentelekomcloud_as_lifecycle_hooks_v1" "as_lifecycle_hook"{
  scaling_lifecycle_hook_name = "as_lifecycle_hook_v1"
  scaling_group_id    = opentelekomcloud_as_group_v1.as_group.id
  scaling_lifecycle_hook_type = "INSTANCE_TERMINATING"
  default_result = "ABANDON"
  default_timeout = 4800
  notification_topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  notification_metadata = "This is a generic message"
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)
