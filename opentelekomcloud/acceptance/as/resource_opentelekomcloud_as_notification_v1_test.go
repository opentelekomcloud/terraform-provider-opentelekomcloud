package acceptance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/notifications"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const asNotificationResourceName = "opentelekomcloud_as_notification_v1.test"

func TestASNotificationV1Configs(t *testing.T) {
	for name, config := range map[string]string{
		"all_events":  testAccASNotificationConfig("test-notifications", false),
		"one_event":   testAccASNotificationConfig("test-notifications", true),
		"data_source": testAccASNotificationsDataSourceConfig("test-notifications"),
	} {
		t.Run(name, func(t *testing.T) {
			if _, diags := hclsyntax.ParseConfig([]byte(config), name+".tf", hcl.InitialPos); diags.HasErrors() {
				t.Fatal(diags.Error())
			}
		})
	}
}

func TestAccASNotificationV1_basic(t *testing.T) {
	name := common.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccASNotificationPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASNotificationsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASNotificationConfig(name, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASNotificationExists(asNotificationResourceName),
					resource.TestCheckResourceAttr(asNotificationResourceName, "topic_scene.#", "5"),
					resource.TestCheckResourceAttr(asNotificationResourceName, "topic_name", name+"-a"),
					resource.TestCheckResourceAttrPair(asNotificationResourceName, "scaling_group_id",
						"opentelekomcloud_as_group_v1.test", "id"),
					resource.TestCheckResourceAttrPair(asNotificationResourceName, "topic_urn",
						"opentelekomcloud_smn_topic_v2.test.0", "id"),
				),
			},
			{
				ResourceName:      asNotificationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccASNotificationConfig(name, false),
				PlanOnly: true,
			},
			{
				Config: testAccASNotificationConfig(name, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASNotificationExists(asNotificationResourceName),
					testAccCheckASNotificationExists("opentelekomcloud_as_notification_v1.peer"),
					resource.TestCheckResourceAttr(asNotificationResourceName, "topic_scene.#", "1"),
					resource.TestCheckTypeSetElemAttr(asNotificationResourceName, "topic_scene.*", "SCALING_DOWN"),
				),
			},
			{
				ResourceName:      asNotificationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccASNotificationConfig(name, true),
				PlanOnly: true,
			},
		},
	})
}

func testAccASNotificationPreCheck(t *testing.T) {
	t.Helper()
	common.TestAccPreCheck(t)
	if os.Getenv("OS_SMN_SUBSCRIPTION_ENDPOINT") == "" {
		t.Skip("OS_SMN_SUBSCRIPTION_ENDPOINT must contain an email subscription endpoint")
	}
	quotas.BookMany(t, quotas.MultipleQuotas{
		{Q: quotas.ASGroup, Count: 1},
		{Q: quotas.ASConfiguration, Count: 1},
	})
}

func testAccCheckASNotificationExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("AS notification %s not found in state", name)
		}
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return err
		}
		items, err := notifications.List(client, rs.Primary.Attributes["scaling_group_id"])
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.TopicURN == rs.Primary.Attributes["topic_urn"] {
				return nil
			}
		}
		return fmt.Errorf("AS notification %s not found remotely", rs.Primary.ID)
	}
}

func testAccCheckASNotificationsDestroy(state *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return err
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_notification_v1" {
			continue
		}
		items, err := notifications.List(client, rs.Primary.Attributes["scaling_group_id"])
		var notFound golangsdk.ErrDefault404
		if errors.As(err, &notFound) {
			continue
		}
		var badRequest golangsdk.ErrDefault400
		if errors.As(err, &badRequest) {
			var body struct {
				ErrorCode string `json:"error_code"`
				Error     struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if json.Unmarshal(badRequest.Body, &body) == nil &&
				(body.Error.Code == "AS.2007" || body.ErrorCode == "AS.2007") {
				continue
			}
		}
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.TopicURN == rs.Primary.Attributes["topic_urn"] {
				return fmt.Errorf("AS notification %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccASNotificationConfig(name string, reduced bool) string {
	scenes := `["SCALING_UP", "SCALING_UP_FAIL", "SCALING_DOWN", "SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL"]`
	if reduced {
		scenes = `["SCALING_DOWN"]`
	}
	return fmt.Sprintf(`
%[1]s
%[2]s
%[3]s

resource "opentelekomcloud_smn_topic_v2" "test" {
  count = 2
  name  = count.index == 0 ? "%[4]s-a" : "%[4]s-b"
}

resource "opentelekomcloud_smn_subscription_v2" "test" {
  count     = 2
  topic_urn = opentelekomcloud_smn_topic_v2.test[count.index].id
  protocol  = "email"
  endpoint  = %[5]q
}

resource "opentelekomcloud_as_configuration_v1" "test" {
  scaling_configuration_name = "%[4]s"
  instance_config {
    image    = data.opentelekomcloud_images_image_v2.latest_image.id
    key_name = %[6]q
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}

resource "opentelekomcloud_as_group_v1" "test" {
  scaling_group_name       = "%[4]s"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.test.id
  delete_instances         = "yes"
  delete_publicip          = true
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
}

resource "opentelekomcloud_as_notification_v1" "test" {
  depends_on       = [opentelekomcloud_smn_subscription_v2.test]
  scaling_group_id = opentelekomcloud_as_group_v1.test.id
  topic_urn        = opentelekomcloud_smn_topic_v2.test[0].id
  topic_scene      = %[7]s
}

resource "opentelekomcloud_as_notification_v1" "peer" {
  depends_on       = [opentelekomcloud_smn_subscription_v2.test]
  scaling_group_id = opentelekomcloud_as_group_v1.test.id
  topic_urn        = opentelekomcloud_smn_topic_v2.test[1].id
  topic_scene      = ["SCALING_UP"]
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet,
		name, os.Getenv("OS_SMN_SUBSCRIPTION_ENDPOINT"), env.OS_KEYPAIR_NAME, scenes)
}
