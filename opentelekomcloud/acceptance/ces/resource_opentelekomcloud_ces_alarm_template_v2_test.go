package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/templates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceAlarmTemplateV2Name = "opentelekomcloud_ces_alarm_template_v2.test"

func getAlarmTemplateV2Func(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CesV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CES v2 client: %s", err)
	}

	template, err := templates.Get(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return template, nil
}

func TestAccCESAlarmTemplateV2_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("ces-template-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("ces-template-%s", acctest.RandString(5))
	rName := resourceAlarmTemplateV2Name

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getAlarmTemplateV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCESAlarmTemplateV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "It is a test template"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
				),
			},
			{
				Config: testAccCESAlarmTemplateV2BasicUpdate(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "It is an updated template"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_associate_alarm"},
			},
		},
	})
}

func TestAccCESAlarmTemplateV2_event(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("ces-template-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("ces-template-%s", acctest.RandString(5))
	rName := resourceAlarmTemplateV2Name

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getAlarmTemplateV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCESAlarmTemplateV2Event(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "It is a test event template"),
					resource.TestCheckResourceAttr(rName, "type", "2"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
				),
			},
			{
				Config: testAccCESAlarmTemplateV2EventUpdate(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "It is an updated event template"),
					resource.TestCheckResourceAttr(rName, "type", "2"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_associate_alarm"},
			},
		},
	})
}

func TestAccCESAlarmTemplateV2_multiplePolicies(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("ces-template-%s", acctest.RandString(5))
	rName := resourceAlarmTemplateV2Name

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getAlarmTemplateV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCESAlarmTemplateV2MultiplePolicies(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "policies.#", "2"),
				),
			},
			{
				Config: testAccCESAlarmTemplateV2MultiplePoliciesUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "policies.#", "3"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_associate_alarm"},
			},
		},
	})
}

func testAccCESAlarmTemplateV2Basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "%s"
  description = "It is a test template"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 1
    filter              = "average"
    comparison_operator = ">="
    value               = 80
    unit                = "%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }
}
`, name)
}

func testAccCESAlarmTemplateV2BasicUpdate(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "%s"
  description = "It is an updated template"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "mem_util"
    period              = 300
    filter              = "max"
    comparison_operator = ">"
    value               = 90
    unit                = "%%"
    count               = 5
    alarm_level         = 1
    suppress_duration   = 3600
  }
}
`, name)
}

func testAccCESAlarmTemplateV2Event(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "%s"
  type        = 2
  description = "It is a test event template"

  policies {
    namespace           = "SYS.ECS"
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    alarm_level         = 2
    suppress_duration   = 0
  }
}
`, name)
}

func testAccCESAlarmTemplateV2EventUpdate(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "%s"
  type        = 2
  description = "It is an updated event template"

  policies {
    namespace           = "SYS.ECS"
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 2
    unit                = "count"
    count               = 3
    alarm_level         = 1
    suppress_duration   = 300
  }
}
`, name)
}

func testAccCESAlarmTemplateV2MultiplePolicies(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "%s"
  description = "Test template with multiple policies"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "mem_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 85
    unit                = "%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }
}
`, name)
}

func testAccCESAlarmTemplateV2MultiplePoliciesUpdate(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "%s"
  description = "Test template with multiple policies - updated"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "mem_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 85
    unit                = "%%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "disk_util_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">="
    value               = 90
    unit                = "%%"
    count               = 3
    alarm_level         = 1
    suppress_duration   = 600
  }
}
`, name)
}
