package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/logs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getDdsLtsLogResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS v3 client: %s", err)
	}

	list, err := logs.List(client, logs.ListOpts{})
	if err != nil {
		return nil, err
	}
	if len(list.InstanceLtsConfigs) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	var listResult logs.InstanceLtsConfigs
	for _, lr := range list.InstanceLtsConfigs {
		if lr.Instance.ID == state.Primary.ID {
			listResult = lr
			break
		}
	}
	if len(listResult.LtsConfigs) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return listResult, nil
}

func TestAccDdsLtsLog_basic(t *testing.T) {
	var instance logs.InstanceLtsConfigs
	resourceName := "opentelekomcloud_dds_lts_log_v3.log"
	rName := fmt.Sprintf("lts_dds%s", acctest.RandString(3))

	rc := common.InitResourceCheck(
		resourceName,
		&instance,
		getDdsLtsLogResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDdsLtsLog_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "opentelekomcloud_dds_instance_v3.instance", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id", "opentelekomcloud_lts_group_v2.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id", "opentelekomcloud_lts_stream_v2.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "audit_log"),
				),
			},
			{
				Config: testAccDdsLtsLog_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "opentelekomcloud_dds_instance_v3.instance", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id", "opentelekomcloud_lts_group_v2.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id", "opentelekomcloud_lts_stream_v2.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "audit_log"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDdsLtsLog_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 1
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[2]s"
}

resource "opentelekomcloud_dds_lts_log_v3" "log" {
  instance_id   = opentelekomcloud_dds_instance_v3.instance.id
  log_type      = "audit_log"
  lts_group_id  = opentelekomcloud_lts_group_v2.test.id
  lts_stream_id = opentelekomcloud_lts_stream_v2.test.id
}`, TestAccDDSInstanceV3ConfigBasic, rName)
}

func testAccDdsLtsLog_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[2]s-update"
  ttl_in_days = 1
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[2]s-update"
}

resource "opentelekomcloud_dds_lts_log_v3" "log" {
  instance_id   = opentelekomcloud_dds_instance_v3.instance.id
  log_type      = "audit_log"
  lts_group_id  = opentelekomcloud_lts_group_v2.test.id
  lts_stream_id = opentelekomcloud_lts_stream_v2.test.id
}`, TestAccDDSInstanceV3ConfigBasic, rName)
}
