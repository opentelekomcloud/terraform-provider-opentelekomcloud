package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/smart_connect"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getDmsKafkav2SmartConnectTaskResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return smart_connect.GetTask(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccDmsKafkav2SmartConnectTask_basic(t *testing.T) {
	var obj interface{}
	rName := fmt.Sprintf("dms-acc-api%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_dms_smart_connect_task_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkav2SmartConnectTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckOBS(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkav2SmartConnectTask_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "destination_type", "OBS_SINK"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.consumer_strategy", "latest"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.destination_file_type", "TEXT"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.record_delimiter", ";"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.deliver_time_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.obs_bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.partition_format", "yyyy/MM/dd/HH/mm"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"destination_task.0.access_key", "destination_task.0.secret_key",
				},
				ImportStateIdFunc: testKafkav2SmartConnectTaskResourceImportState(resourceName),
			},
		},
	})
}

func TestAccDmsKafkav2SmartConnectTask_KafkaToKafka(t *testing.T) {
	var obj interface{}
	rName := fmt.Sprintf("dms-acc-api%s", acctest.RandString(5))
	resourceName := "opentelekomcloud_dms_smart_connect_task_v2.test"

	rc := common.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkav2SmartConnectTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkav2SmartConnectTask_kafkaToKafka(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "task_name"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "KAFKA_REPLICATOR_SOURCE"),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "source_task.0.peer_instance_id",
						"opentelekomcloud_dms_dedicated_instance_v2.test2", "id"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.direction", "push"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.replication_factor", "3"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.task_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.rename_topic_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.provenance_header_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.sync_consumer_offsets_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.consumer_strategy", "latest"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.compression_type", "snappy"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKafkav2SmartConnectTaskResourceImportState(resourceName),
			},
		},
	})
}

func testDmsKafkav2SmartConnectTask_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_dms_smart_connect_v2" "test" {
  instance_id       = opentelekomcloud_dms_dedicated_instance_v2.test.id
  storage_spec_code = "dms.physical.storage.high.v2"
  node_count        = 2
  bandwidth         = "100MB"
}

resource "opentelekomcloud_dms_topic_v1" "test" {
  instance_id    = opentelekomcloud_dms_dedicated_instance_v2.test.id
  name           = "%[2]s"
  partition      = 10
  retention_time = 36
}

resource "opentelekomcloud_dms_smart_connect_task_v2" "test" {
  depends_on = [opentelekomcloud_dms_smart_connect_v2.test, opentelekomcloud_dms_topic_v1.test]

  instance_id      = opentelekomcloud_dms_dedicated_instance_v2.test.id
  task_name        = "%[2]s"
  destination_type = "OBS_SINK"
  topics           = [opentelekomcloud_dms_topic_v1.test.name]

  destination_task {
    consumer_strategy     = "latest"
    destination_file_type = "TEXT"
    access_key            = "%[3]s"
    secret_key            = "%[4]s"
    obs_bucket_name       = opentelekomcloud_obs_bucket.test.bucket
    partition_format      = "yyyy/MM/dd/HH/mm"
    record_delimiter      = ";"
    deliver_time_interval = 300
  }
}`, testAccKafkaInstance_newFormat(rName), rName, env.OS_ACCESS_KEY, env.OS_SECRET_KEY)
}

func testDmsKafkav2SmartConnectTask_kafkaToKafka(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_dms_smart_connect_task_v2" "test" {
  depends_on = [
    opentelekomcloud_dms_smart_connect_v2.test1,
    opentelekomcloud_dms_smart_connect_v2.test2,
    opentelekomcloud_dms_topic_v1.test
  ]

  instance_id = opentelekomcloud_dms_dedicated_instance_v2.test1.id
  task_name   = "%[2]s"
  topics      = [opentelekomcloud_dms_topic_v1.test.name]
  source_type = "KAFKA_REPLICATOR_SOURCE"

  source_task {
    peer_instance_id              = opentelekomcloud_dms_dedicated_instance_v2.test2.id
    direction                     = "push"
    replication_factor            = 3
    task_num                      = 2
    provenance_header_enabled     = true
    sync_consumer_offsets_enabled = true
    rename_topic_enabled          = true
    consumer_strategy             = "latest"
    compression_type              = "snappy"
  }
}`, testAccKafkav2SmartConnectTaskKafkaToKafKaBase(rName), rName)
}

func testAccKafkav2SmartConnectTaskKafkaToKafKaBase(rName string) string {
	kafka1 := testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase1(rName)
	kafka2 := testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase2(rName)
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_networking_secgroup_v2" "test" {
  name        = "secgroup_dms"
  description = "My neutron security group"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "rule" {
  security_group_id = opentelekomcloud_networking_secgroup_v2.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  port_range_min    = "9092"
  port_range_max    = "9092"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "out_v4_all" {
  security_group_id = opentelekomcloud_networking_secgroup_v2.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  remote_ip_prefix  = "0.0.0.0/0"
}

data "opentelekomcloud_dms_az_v1" "test" {}

data "opentelekomcloud_dms_flavor_v2" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.opentelekomcloud_dms_flavor_v2.test.flavors[0]
}

%s

%s

resource "opentelekomcloud_dms_topic_v1" "test" {
  instance_id = opentelekomcloud_dms_dedicated_instance_v2.test1.id
  name        = "%s"
  partition   = 10
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, kafka1, kafka2, rName)
}

func testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase1(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_dedicated_instance_v2" "test1" {
  name              = "%[1]s-1"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  available_zones = [
    data.opentelekomcloud_dms_az_v1.test.id
  ]
  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3
  arch_type      = "X86"
}

resource "opentelekomcloud_dms_smart_connect_v2" "test1" {
  instance_id       = opentelekomcloud_dms_dedicated_instance_v2.test1.id
  storage_spec_code = "dms.physical.storage.high.v2"
  node_count        = 2
  bandwidth         = "100MB"
}`, rName)
}
func testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase2(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_dedicated_instance_v2" "test2" {
  name              = "%[1]s-2"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  available_zones = [
    data.opentelekomcloud_dms_az_v1.test.id
  ]
  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3
  arch_type      = "X86"
}

resource "opentelekomcloud_dms_smart_connect_v2" "test2" {
  instance_id       = opentelekomcloud_dms_dedicated_instance_v2.test2.id
  storage_spec_code = "dms.physical.storage.high.v2"
  node_count        = 2
  bandwidth         = "100MB"
}`, rName)
}

func testKafkav2SmartConnectTaskResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}

func testAccKafkaInstance_newFormat(rName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "test" {}

data "opentelekomcloud_dms_flavor_v2" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.opentelekomcloud_dms_flavor_v2.test.flavors[0]
}

resource "opentelekomcloud_dms_dedicated_instance_v2" "test" {
  name              = "%s"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id

  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  available_zones = [
    data.opentelekomcloud_dms_az_v1.test.id
  ]
  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3
  arch_type      = "X86"

  ssl_enable         = true
  access_user        = "user"
  password           = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = ""
  }
  cross_vpc_accesses {
    advertised_ip = "www.terraform-test.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.53"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, rName)
}
