package tms

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	rt "github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/resource-tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/tms"
)

func getResourceTagsFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.TmsV2Client()
	if err != nil {
		return nil, fmt.Errorf("error creating TMS v2 client: %s", err)
	}

	var (
		projectId       = state.Primary.Attributes["project_id"]
		resourcesLen, _ = strconv.Atoi(state.Primary.Attributes["resources.#"])
		tagsLen, _      = strconv.Atoi(state.Primary.Attributes["tags.%"])
		tagsConfigured  = false
	)

	for i := 0; i < resourcesLen; i++ {
		resourceId := state.Primary.Attributes[fmt.Sprintf("resources.%d.resource_id", i)]
		opts := rt.ListOpts{
			ResourceType: state.Primary.Attributes[fmt.Sprintf("resources.%d.resource_type", i)],
			ProjectId:    projectId,
		}
		resp, err := rt.List(client, resourceId, opts)
		if err != nil {
			return nil, fmt.Errorf("error query resource (%s) tags: %s", resourceId, err)
		}
		actualTags := tms.FlattenTagsToMap(resp)
		if len(actualTags) != tagsLen {
			return nil, fmt.Errorf("some tags were not set successfully")
		}
		if len(actualTags) > 0 {
			tagsConfigured = true
		}
	}
	if !tagsConfigured {
		return nil, golangsdk.ErrDefault404{}
	}
	return tagsConfigured, nil
}

func TestAccResourceTags_basic(t *testing.T) {
	var (
		tagsConfigured bool

		rName    = "opentelekomcloud_tms_resource_tags_v1.test"
		basicCfg = testAccResourceTags_base()
		rc       = common.InitResourceCheck(rName, &tagsConfigured, getResourceTagsFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTags_basic(basicCfg),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "resources.0.resource_type", "disk"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.resource_id", "opentelekomcloud_evs_volume_v3.volume", "id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccResourceTags_basic_update(basicCfg),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "resources.0.resource_type", "disk"),
					resource.TestCheckResourceAttrPair(rName, "resources.0.resource_id", "opentelekomcloud_evs_volume_v3.volume", "id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(rName, "tags.creator", "terraform"),
				),
			},
		},
	})
}

func testAccResourceTags_base() string {
	var (
		name = fmt.Sprintf("rt_%s", acctest.RandString(3))
	)

	return fmt.Sprintf(`
data "opentelekomcloud_identity_project_v3" "project_1" {}

resource "opentelekomcloud_evs_volume_v3" "volume" {
  name              = "%s"
  description       = "tms test volume"
  availability_zone = "eu-de-01"
  volume_type       = "SSD"
  size              = 20

  lifecycle {
    ignore_changes = [
      tags
    ]
  }

}
`, name)
}

func testAccResourceTags_basic(basicCfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_tms_resource_tags_v1" "test" {
  project_id = data.opentelekomcloud_identity_project_v3.project_1.id

  resources {
    resource_type = "disk"
    resource_id   = opentelekomcloud_evs_volume_v3.volume.id
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, basicCfg)
}

func testAccResourceTags_basic_update(basicCfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_tms_resource_tags_v1" "test" {
  project_id = data.opentelekomcloud_identity_project_v3.project_1.id

  resources {
    resource_type = "disk"
    resource_id   = opentelekomcloud_evs_volume_v3.volume.id
  }

  tags = {
    foo     = "baar"
    creator = "terraform"
  }
}
`, basicCfg)
}
