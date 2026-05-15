package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/alarms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceAlarmRuleV2Name = "opentelekomcloud_ces_alarm_rule_v2.alarmrule_1"

func getAlarmRuleV2Func(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CesV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CES v2 client: %s", err)
	}
	listResp, err := alarms.List(c, alarms.ListOpts{
		AlarmId: state.Primary.ID,
	})
	if err != nil {
		return nil, err
	}
	if listResp == nil || len(listResp.Alarms) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return listResp.Alarms[0], nil
}

func TestCESAlarmRuleV2_basic(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = resourceAlarmRuleV2Name
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "alarm_rule_v2_1"),
					resource.TestCheckResourceAttr(rName, "namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "type", "MULTI_INSTANCE"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "false"),
				),
			},
			{
				Config: testCESAlarmRuleV2Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "alarm_rule_v2_1"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "false"),
				),
			},
			{
				ResourceName:      resourceAlarmRuleV2Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestCESAlarmRuleV2_withNotification(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = resourceAlarmRuleV2Name
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2WithNotification,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "alarm_rule_v2_1"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "alarm_actions.0.type", "notification"),
				),
			},
		},
	})
}

var testCESAlarmRuleV2Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_ces_alarm_rule_v2" "alarmrule_1" {
  name      = "alarm_rule_v2_1"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }

  notification_enabled = false
  alarm_enabled        = true
}
`, common.DataSourceSubnet)

var testCESAlarmRuleV2Update = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_ces_alarm_rule_v2" "alarmrule_1" {
  name      = "alarm_rule_v2_1"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 10
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }

  notification_enabled = false
  alarm_enabled        = false
}
`, common.DataSourceSubnet)

var testCESAlarmRuleV2WithNotification = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "alarmrule_1" {
  name      = "alarm_rule_v2_1"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.topic_1.topic_urn
    ]
  }
}
`, common.DataSourceSubnet)

func TestCESAlarmRuleV2_multiplePolicies(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2MultiplePolicies(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "type", "MULTI_INSTANCE"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "policies.#", "2"),
				),
			},
			{
				Config: testCESAlarmRuleV2MultiplePoliciesUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "false"),
					resource.TestCheckResourceAttr(rName, "policies.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2MultiplePolicies(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  name        = "ecs-%[2]s"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test.id
    }
  }

  policies {
    metric_name         = "network_incoming_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    level               = 3
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    level               = 3
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}

func testCESAlarmRuleV2MultiplePoliciesUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  name        = "ecs-%[2]s"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test.id
    }
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 1200
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    level               = 4
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 3600
    filter              = "average"
    comparison_operator = ">="
    value               = 20
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    level               = 4
  }

  notification_enabled = true
  alarm_enabled        = false

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}

func TestCESAlarmRuleV2_multipleResources(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 2},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 2},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 8},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2MultipleResources(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
				),
			},
			{
				Config: testCESAlarmRuleV2MultipleResourcesUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "resources.#", "2"),
				),
			},
		},
	})
}

func testCESAlarmRuleV2MultipleResources(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  count       = 2
  name        = "ecs-%[2]s-${count.index}"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[0].id
    }
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}

func testCESAlarmRuleV2MultipleResourcesUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  count       = 2
  name        = "ecs-%[2]s-${count.index}"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[0].id
    }
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[1].id
    }
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}

func TestCESAlarmRuleV2_multipleDimensions(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 2},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 2},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 8},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2MultipleDimensions(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "type", "MULTI_INSTANCE"),
					resource.TestCheckResourceAttr(rName, "resources.#", "2"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "false"),
				),
			},
			{
				Config: testCESAlarmRuleV2MultipleDimensionsUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2MultipleDimensions(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  count       = 2
  name        = "ecs-%[2]s-${count.index}"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[0].id
    }
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[1].id
    }
  }

  policies {
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%%%%"
    count               = 1
    suppress_duration   = 300
    level               = 2
  }

  notification_enabled = false
  alarm_enabled        = true
}
`, common.DataSourceSubnet, name)
}

func testCESAlarmRuleV2MultipleDimensionsUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  count       = 2
  name        = "ecs-%[2]s-${count.index}"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[0].id
    }
  }

  policies {
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%%%%"
    count               = 1
    suppress_duration   = 300
    level               = 2
  }

  notification_enabled = false
  alarm_enabled        = false
}
`, common.DataSourceSubnet, name)
}

func TestCESAlarmRuleV2_sysEvent(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2SysEvent(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "type", "EVENT.SYS"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "false"),
				),
			},
			{
				Config: testCESAlarmRuleV2SysEventUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "false"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2SysEvent(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[1]s"
  namespace = "SYS.ECS"
  type      = "EVENT.SYS"

  policies {
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    suppress_duration   = 0
    level               = 2
  }

  notification_enabled = false
  alarm_enabled        = true
}
`, name)
}

func testCESAlarmRuleV2SysEventUpdate(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[1]s"
  namespace = "SYS.ECS"
  type      = "EVENT.SYS"

  policies {
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 2
    unit                = "count"
    count               = 2
    suppress_duration   = 300
    level               = 3
  }

  notification_enabled = false
  alarm_enabled        = false
}
`, name)
}

func TestCESAlarmRuleV2_sysEventWithNotification(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2SysEventWithNotification(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "type", "EVENT.SYS"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_actions.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2SysEventWithNotification(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[1]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[1]s"
  namespace = "SYS.ECS"
  type      = "EVENT.SYS"

  policies {
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    suppress_duration   = 0
    level               = 2
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, name)
}

func TestCESAlarmRuleV2_allInstance(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2AllInstance(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ALL_INSTANCE"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2AllInstance(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[1]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[1]s"
  namespace = "AGT.ECS"
  type      = "ALL_INSTANCE"

  resources {
    dimensions {
      name = "instance_id"
    }

    dimensions {
      name = "mount_point"
    }
  }

  policies {
    metric_name         = "disk_usedPercent"
    period              = 1
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    count               = 1
    suppress_duration   = 0
    level               = 2
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, name)
}

func TestCESAlarmRuleV2_withOkActions(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2WithOkActions(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "alarm_actions.0.type", "notification"),
					resource.TestCheckResourceAttr(rName, "ok_actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "ok_actions.0.type", "notification"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2WithOkActions(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  name        = "ecs-%[2]s"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[2]s"
  namespace = "SYS.ECS"
  type      = "MULTI_INSTANCE"

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test.id
    }
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  ok_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}

func TestCESAlarmRuleV2_withAlarmTemplate(t *testing.T) {
	var (
		ar    alarms.Alarm
		rName = "opentelekomcloud_ces_alarm_rule_v2.test"
		name  = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&ar,
		getAlarmRuleV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 2},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 2},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 8},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleV2WithAlarmTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "type", "MULTI_INSTANCE"),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttrPair(rName, "alarm_template_id",
						"opentelekomcloud_ces_alarm_template_v2.test", "id"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
				),
			},
			{
				Config: testCESAlarmRuleV2WithAlarmTemplateUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "alarm_enabled", "false"),
					resource.TestCheckResourceAttrPair(rName, "alarm_template_id",
						"opentelekomcloud_ces_alarm_template_v2.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleV2WithAlarmTemplate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  count       = 2
  name        = "ecs-%[2]s-${count.index}"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "template-%[2]s"
  description = "Test alarm template for alarm rule"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%%%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  depends_on = [
    opentelekomcloud_compute_instance_v2.test,
    opentelekomcloud_smn_topic_v2.test
  ]
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name              = "%[2]s"
  namespace         = "SYS.ECS"
  type              = "MULTI_INSTANCE"
  alarm_template_id = opentelekomcloud_ces_alarm_template_v2.test.id

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[0].id
    }
  }

  notification_enabled = true
  alarm_enabled        = true

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  ok_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}

func testCESAlarmRuleV2WithAlarmTemplateUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  count       = 2
  name        = "ecs-%[2]s-${count.index}"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}

resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "template-%[2]s"
  description = "Test alarm template for alarm rule"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%%%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  depends_on = [
    opentelekomcloud_compute_instance_v2.test,
    opentelekomcloud_smn_topic_v2.test
  ]
}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name              = "%[2]s"
  namespace         = "SYS.ECS"
  type              = "MULTI_INSTANCE"
  alarm_template_id = opentelekomcloud_ces_alarm_template_v2.test.id

  resources {
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test[1].id
    }
  }

  notification_enabled = true
  alarm_enabled        = false

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  ok_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }
}
`, common.DataSourceSubnet, name)
}
