package acceptance

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	var notification, peer notifications.Notification
	name := common.RandomAccResourceName()
	rc := common.InitResourceCheck(asNotificationResourceName, &notification, getASNotificationFunc)
	peerRC := common.InitResourceCheck("opentelekomcloud_as_notification_v1.peer", &peer, getASNotificationFunc)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccASNotificationPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASNotificationConfig(name, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					peerRC.CheckResourceExists(),
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
					rc.CheckResourceExists(),
					peerRC.CheckResourceExists(),
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

func getASNotificationFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Auto Scaling v1 client: %w", err)
	}
	items, err := notifications.List(client, state.Primary.Attributes["scaling_group_id"])
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.TopicURN == state.Primary.Attributes["topic_urn"] {
			return &item, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func TestGetASNotificationFunc(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		found  bool
	}{
		{"matching_topic", http.StatusOK, `{"topics":[{"topic_urn":"other"},{"topic_urn":"test-topic","topic_name":"test","topic_scene":["SCALING_UP"]}]}`, true},
		{"empty_group", http.StatusOK, `{"topics":[]}`, false},
		{"different_topic", http.StatusOK, `{"topics":[{"topic_urn":"other"}]}`, false},
		{"missing_group", http.StatusNotFound, `{"error_msg":"fixture error"}`, false},
		{"forbidden", http.StatusForbidden, `{"error_msg":"fixture error"}`, false},
		{"bad_request", http.StatusBadRequest, `{"error":{"code":"AS.2007","message":"fixture error"}}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/scaling_notification/group" {
					t.Errorf("unexpected notification request: %s %s", r.Method, r.URL.Path)
				}
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer server.Close()
			config := &cfg.Config{HwClient: &golangsdk.ProviderClient{
				EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
					return server.URL + "/", nil
				},
			}}
			state := &terraform.ResourceState{Primary: &terraform.InstanceState{
				Attributes: map[string]string{"scaling_group_id": "group", "topic_urn": "test-topic"},
			}}
			actual, err := getASNotificationFunc(config, state)
			if tc.status != http.StatusOK {
				if actual != nil || err == nil || !strings.Contains(err.Error(), "fixture error") {
					t.Fatalf("getter must propagate the API error, not report an empty result: %#v, %v", actual, err)
				}
				return
			}
			if !tc.found {
				var notFound golangsdk.ErrDefault404
				if actual != nil || !errors.As(err, &notFound) {
					t.Fatalf("missing topic must return a typed not-found error: %#v, %v", actual, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			notification, ok := actual.(*notifications.Notification)
			if !ok || notification.TopicURN != "test-topic" || notification.TopicName != "test" ||
				len(notification.TopicScene) != 1 || notification.TopicScene[0] != "SCALING_UP" {
				t.Fatalf("getter did not return the matching SDK object: %#v", actual)
			}
		})
	}
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
