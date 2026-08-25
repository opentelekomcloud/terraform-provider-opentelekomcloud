package acceptance

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pg_ext "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/postgres-extensions"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getPostgresExtension(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.RdsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RDSv3 Client: %s", err)
	}
	// resource ID = "instanceId/DatabaseName/ExtensionName"
	parts := strings.Split(state.Primary.ID, "/")
	resp, err := pg_ext.List(client, parts[0], pg_ext.ListOpts{
		DatabaseName: parts[1],
	})
	if err != nil {
		return nil, err
	}
	for _, ext := range resp.Extensions {
		if ext.Name == parts[2] && ext.Created {
			return ext, nil
		}
	}
	return nil, fmt.Errorf("Extension not created")
}

func TestAccRdsPostgresExtensionV3Basic(t *testing.T) {
	rdsId := os.Getenv("OS_RDS_ID")
	if rdsId == "" {
		t.Skip("OS_RDS_ID env var required for the test is missing")
	}
	resName := "opentelekomcloud_rds_postgres_extension_v3.extension"
	var ext pg_ext.Extension

	rc := common.InitResourceCheck(
		resName,
		&ext,
		getPostgresExtension,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsPostgresExtensionV3Basic(rdsId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resName, "version"),
					resource.TestCheckResourceAttr(resName, "extension_name", "hstore"),
					resource.TestCheckResourceAttr(resName, "database_name", "postgres"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRdsPostgresExtensionV3Basic(rdsId string) string {
	return fmt.Sprintf(`

resource "opentelekomcloud_rds_postgres_extension_v3" "extension" {
  instance_id    = "%s"
  database_name  = "postgres"
  extension_name = "hstore"
}
`, rdsId)
}
