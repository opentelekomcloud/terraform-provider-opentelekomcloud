package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	ac "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/access-config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getHostAccessConfigResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v3 client: %s", err)
	}

	requestResp, err := ac.List(client, ac.ListOpts{})
	if err != nil {
		return nil, err
	}
	if len(requestResp.Result) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var accessResult *ac.AccessConfigInfo
	for _, acc := range requestResp.Result {
		if acc.ID == state.Primary.ID {
			accessResult = &acc
		}
	}
	if accessResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return accessResult, nil
}

func TestAccHostAccessConfigV3_basic(t *testing.T) {
	var (
		access ac.AccessConfigInfo
		rName  = "opentelekomcloud_lts_host_access_v3.basic"
		name   = fmt.Sprintf("lts_access%s", acctest.RandString(3))
		rc     = common.InitResourceCheck(rName, &access, getHostAccessConfigResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostAccessConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.#", "2"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.#", "2"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "0"),
					resource.TestCheckResourceAttr(rName, "access_type", "AGENT"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "opentelekomcloud_lts_group_v2.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "opentelekomcloud_lts_stream_v2.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
				),
			},
			{
				Config: testHostAccessConfig_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.#", "1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var/log/*/*.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/log/*/a.log"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value-updated"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "host_group_ids.0", "opentelekomcloud_lts_host_group_v3.test", "id"),
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

func TestAccHostAccessConfigV3_windows(t *testing.T) {
	var (
		access ac.AccessConfigInfo
		rName  = "opentelekomcloud_lts_host_access_v3.windows"
		name   = fmt.Sprintf("lts_group%s", acctest.RandString(3))
		rc     = common.InitResourceCheck(rName, &access, getHostAccessConfigResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostAccessConfig_windows_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "D:\\data\\log\\*"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "D:\\data\\log\\a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.windows_log_info.0.time_offset", "7"),
					resource.TestCheckResourceAttr(rName, "access_config.0.windows_log_info.0.time_offset_unit", "day"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "access_type", "AGENT"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "opentelekomcloud_lts_group_v2.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "opentelekomcloud_lts_stream_v2.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
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

func testHostAccessConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "test" {
  name = "%[1]s"
  type = "linux"
}

resource "opentelekomcloud_lts_host_access_v3" "basic" {
  name          = "%[1]s"
  log_group_id  = opentelekomcloud_lts_group_v2.test.id
  log_stream_id = opentelekomcloud_lts_stream_v2.test.id

  access_config {
    paths       = ["/var/temp", "/var/log/*"]
    black_paths = ["/var/temp", "/var/log/*/a.log"]

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name)
}

func testHostAccessConfig_basic_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "test" {
  name = "%[1]s"
  type = "linux"
}

resource "opentelekomcloud_lts_host_access_v3" "basic" {
  name           = "%[1]s"
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.test.id
  host_group_ids = [opentelekomcloud_lts_host_group_v3.test.id]

  access_config {
    paths       = ["/var/log/*/*.log"]
    black_paths = ["/var/log/*/a.log"]

    multi_log_format {
      mode  = "time"
      value = "YYYY-MM-DD hh:mm:ss"
    }
  }

  tags = {
    key   = "value-updated"
    owner = "terraform"
  }
}
`, name)
}

func testHostAccessConfig_windows_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_host_group_v3" "windows" {
  name = "%[1]s"
  type = "windows"
}

resource "opentelekomcloud_lts_host_access_v3" "windows" {
  name           = "%[1]s"
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.test.id
  host_group_ids = [opentelekomcloud_lts_host_group_v3.windows.id]

  access_config {
    paths       = ["D:\\data\\log\\*"]
    black_paths = ["D:\\data\\log\\a.log"]

    windows_log_info {
      categories       = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name)
}
