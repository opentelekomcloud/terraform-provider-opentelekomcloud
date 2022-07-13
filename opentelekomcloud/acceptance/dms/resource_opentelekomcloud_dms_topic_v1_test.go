package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/topics"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceTopicName = "opentelekomcloud_dms_topic_v1.topic_1"
const resourceInstName = "opentelekomcloud_dms_instance_v1.instance_1"

func TestAccDmsTopicsV1_basic(t *testing.T) {
	var topic topics.Topic
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	var topicName = fmt.Sprintf("topic_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1TopicBasic(instanceName, topicName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1TopicExists(resourceTopicName, resourceInstName, topic),
				),
			},
		},
	})
}

func testAccCheckDmsV1TopicExists(n string, i string, topic topics.Topic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		rsb, ok := s.RootModule().Resources[i]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rsb.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DmsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
		}

		v, err := topics.Get(client, rsb.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting OpenTelekomCloud DMSv1 topics (%s): %w", rs.Primary.ID, err)
		}

		if v.Topics[0].Name != rs.Primary.ID {
			return fmt.Errorf("DMS topic not found")
		}
		topic = *v
		return nil
	}
}

func testAccDmsV1TopicBasic(instanceName string, topicName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "opentelekomcloud_dms_instance_v1" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
}

resource "opentelekomcloud_dms_topic_v1" "topic_1" {
  instance_id = resource.opentelekomcloud_dms_instance_v1.instance_1.id
  name = "%s"
  partition = 10
  replication = 2
  sync_replication = true
  retention_time = 80
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName, topicName)
}
