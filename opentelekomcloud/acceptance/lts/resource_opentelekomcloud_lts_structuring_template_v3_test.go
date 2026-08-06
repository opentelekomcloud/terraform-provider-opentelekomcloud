package acceptance

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	cloud_structuring "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/cloud-structuring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	ltsStructuringTemplateIDEnv   = "OS_LTS_STRUCTURING_TEMPLATE_ID"
	ltsStructuringTemplateNameEnv = "OS_LTS_STRUCTURING_TEMPLATE_NAME"
)

func getLtsStructuringTemplateResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v2 client: %s", err)
	}

	template, err := cloud_structuring.Get(
		client,
		state.Primary.Attributes["log_group_id"],
		state.Primary.Attributes["log_stream_id"],
	)
	if errors.Is(err, io.EOF) {
		return nil, golangsdk.ErrDefault404{}
	}
	if err != nil {
		return nil, err
	}
	if template.ID == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return template, nil
}

func TestAccLtsStructuringTemplateV3_basic(t *testing.T) {
	var (
		template cloud_structuring.StructuringResponse
		name     = common.RandomAccResourceName()
		rName    = "opentelekomcloud_lts_structuring_template_v3.test"
		rc       = common.InitResourceCheck(rName, &template, getLtsStructuringTemplateResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLtsStructuringTemplateV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"opentelekomcloud_lts_group_v2.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"opentelekomcloud_lts_stream_v2.stream", "id"),
					resource.TestCheckResourceAttr(rName, "template_type", "built_in"),
					resource.TestCheckResourceAttr(rName, "template_name", "CTS"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_fields.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				Config: testLtsStructuringTemplateV3BasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "template_type", "built_in"),
					resource.TestCheckResourceAttr(rName, "template_name", "GAUSSDB_MYSQL_SLOW"),
					resource.TestCheckResourceAttr(rName, "quick_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_fields.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testLtsStructuringTemplateV3ImportState(rName),
				ImportStateVerifyIgnore: ltsStructuringTemplateImportIgnore,
			},
		},
	})
}

func TestAccLtsStructuringTemplateV3_custom(t *testing.T) {
	templateID := os.Getenv(ltsStructuringTemplateIDEnv)
	templateName := os.Getenv(ltsStructuringTemplateNameEnv)
	if templateID == "" || templateName == "" {
		t.Skipf("%s and %s must be set for this acceptance test",
			ltsStructuringTemplateIDEnv, ltsStructuringTemplateNameEnv)
	}

	var (
		template cloud_structuring.StructuringResponse
		name     = common.RandomAccResourceName()
		rName    = "opentelekomcloud_lts_structuring_template_v3.test"
		rc       = common.InitResourceCheck(rName, &template, getLtsStructuringTemplateResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLtsStructuringTemplateV3Custom(name, templateID, templateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "template_type", "custom"),
					resource.TestCheckResourceAttr(rName, "template_id", templateID),
					resource.TestCheckResourceAttr(rName, "template_name", templateName),
					resource.TestCheckResourceAttr(rName, "quick_analysis", "true"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				Config: testLtsStructuringTemplateV3CustomUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "template_type", "built_in"),
					resource.TestCheckResourceAttr(rName, "template_name", "GAUSSDB_MYSQL_SLOW"),
					resource.TestCheckResourceAttr(rName, "quick_analysis", "false"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testLtsStructuringTemplateV3ImportState(rName),
				ImportStateVerifyIgnore: ltsStructuringTemplateImportIgnore,
			},
		},
	})
}

var ltsStructuringTemplateImportIgnore = []string{
	"template_type",
	"template_id",
	"demo_fields",
	"tag_fields",
	"quick_analysis",
}

func testLtsStructuringTemplateV3Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_lts_structuring_template_v3" "test" {
  log_group_id  = opentelekomcloud_lts_group_v2.test.id
  log_stream_id = opentelekomcloud_lts_stream_v2.stream.id
  template_type = "built_in"
  template_name = "CTS"

  demo_fields {
    field_name  = "event_type"
    is_analysis = true
  }

  demo_fields {
    field_name  = "resource_type"
    is_analysis = false
  }

  tag_fields {
    field_name  = "hostIP"
    is_analysis = true
  }
}
`, testAccLtsV2Stream_basic(name))
}

func testLtsStructuringTemplateV3BasicUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_lts_structuring_template_v3" "test" {
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.stream.id
  template_type  = "built_in"
  template_name  = "GAUSSDB_MYSQL_SLOW"
  quick_analysis = true

  demo_fields {
    field_name  = "query_time"
    is_analysis = true
  }

  demo_fields {
    field_name  = "rows_examined"
    is_analysis = false
  }

  tag_fields {
    field_name  = "hostName"
    is_analysis = false
  }
}
`, testAccLtsV2Stream_basic(name))
}

func testLtsStructuringTemplateV3Custom(name, templateID, templateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_lts_structuring_template_v3" "test" {
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.stream.id
  template_type  = "custom"
  template_id    = "%[2]s"
  template_name  = "%[3]s"
  quick_analysis = true
}
`, testAccLtsV2Stream_basic(name), templateID, templateName)
}

func testLtsStructuringTemplateV3CustomUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_lts_structuring_template_v3" "test" {
  log_group_id   = opentelekomcloud_lts_group_v2.test.id
  log_stream_id  = opentelekomcloud_lts_stream_v2.stream.id
  template_type  = "built_in"
  template_name  = "GAUSSDB_MYSQL_SLOW"
  quick_analysis = false

  demo_fields {
    field_name  = "query_time"
    is_analysis = true
  }

  demo_fields {
    field_name  = "rows_examined"
    is_analysis = true
  }

  tag_fields {
    field_name  = "hostName"
    is_analysis = true
  }
}
`, testAccLtsV2Stream_basic(name))
}

func testLtsStructuringTemplateV3ImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		groupID := rs.Primary.Attributes["log_group_id"]
		streamID := rs.Primary.Attributes["log_stream_id"]
		if groupID == "" || streamID == "" {
			return "", fmt.Errorf("log_group_id and log_stream_id must be set to build the import ID")
		}
		return fmt.Sprintf("%s/%s", groupID, streamID), nil
	}
}
