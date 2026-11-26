package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gpfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGPFSACLName = "opentelekomcloud_gpfs_acl_v3.acl"

func getGPFSACLFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.NewGPFSClient(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GPFS client: %s", err)
	}
	output, err := client.GetFSAccessRules(state.Primary.ID)
	if err != nil {
		if serviceErr, ok := err.(gpfs.Error); ok && strings.HasPrefix(serviceErr.Status, "404") {
			return nil, golangsdk.ErrDefault404{}
		}
		return nil, err
	}
	if len(output.Statement) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return output, nil
}

func checkGPFSACLGrants(output *gpfs.GetFSAccessRulesOutput, count int, action string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if len(output.Statement) != count {
			return fmt.Errorf("expected %d live ACL statements, got %d", count, len(output.Statement))
		}
		expected := map[string]bool{
			state.RootModule().Resources["opentelekomcloud_vpc_v1.test.1"].Primary.ID: false,
		}
		if count == 2 {
			expected[state.RootModule().Resources["opentelekomcloud_vpc_v1.test.0"].Primary.ID] = false
		}
		for _, statement := range output.Statement {
			if _, ok := expected[statement.Condition.SourceVPC]; !ok {
				return fmt.Errorf("unexpected VPC grant in live ACL: %s", statement.Condition.SourceVPC)
			}
			if expected[statement.Condition.SourceVPC] {
				return fmt.Errorf("duplicate VPC grant in live ACL: %s", statement.Condition.SourceVPC)
			}
			expected[statement.Condition.SourceVPC] = true
			if string(statement.Action) != action {
				return fmt.Errorf("expected action %s, got %s", action, statement.Action)
			}
		}
		return nil
	}
}

func checkGPFSACLImport(output *gpfs.GetFSAccessRulesOutput) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected one imported ACL, got %d", len(states))
		}
		attributes := states[0].Attributes
		if attributes["statement.#"] != fmt.Sprint(len(output.Statement)) {
			return fmt.Errorf("import did not retain the complete ACL")
		}
		for _, statement := range output.Statement {
			found := false
			for i := range output.Statement {
				prefix := fmt.Sprintf("statement.%d.", i)
				if attributes[prefix+"vpc_id"] != statement.Condition.SourceVPC {
					continue
				}
				found = true
				if attributes[prefix+"action"] != string(statement.Action) ||
					attributes[prefix+"effect"] != string(statement.Effect) ||
					attributes[prefix+"sid"] != statement.ID ||
					attributes[prefix+"vpc_source_ip.#"] != fmt.Sprint(len(statement.Condition.VPCSourceIP)) {
					return fmt.Errorf("import changed statement for VPC %s", statement.Condition.SourceVPC)
				}
				for j, address := range statement.Condition.VPCSourceIP {
					if attributes[fmt.Sprintf("%svpc_source_ip.%d", prefix, j)] != address {
						return fmt.Errorf("import changed source IP for VPC %s", statement.Condition.SourceVPC)
					}
				}
				break
			}
			if !found {
				return fmt.Errorf("import omitted VPC %s", statement.Condition.SourceVPC)
			}
		}
		return nil
	}
}

func TestAccGPFSACLV3_multiVpc(t *testing.T) {
	name := fmt.Sprintf("terraform-acc-%s", acctest.RandString(5))
	var acl gpfs.GetFSAccessRulesOutput
	rc := common.InitResourceCheck(resourceGPFSACLName, &acl, getGPFSACLFunc)
	fsCheck := common.InitResourceCheck(resourceGPFSFileSystemName, &gpfs.Bucket{}, getGPFSFileSystemFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      resource.ComposeTestCheckFunc(rc.CheckResourceDestroy(), fsCheck.CheckResourceDestroy()),
		Steps: []resource.TestStep{
			{
				Config: testAccGPFSACLV3Config(name, true, "FullControl"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					fsCheck.CheckResourceExists(),
					checkGPFSACLGrants(&acl, 2, "FullControl"),
					resource.TestCheckResourceAttr(resourceGPFSACLName, "id", name),
					resource.TestCheckResourceAttr(resourceGPFSACLName, "statement.#", "2"),
				),
			},
			{
				Config:   testAccGPFSACLV3Config(name, true, "FullControl"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceGPFSACLName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"statement"},
				ImportStateCheck:        checkGPFSACLImport(&acl),
			},
			{
				Config: testAccGPFSACLV3Config(name, false, "FullControl"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					checkGPFSACLGrants(&acl, 1, "FullControl"),
					resource.TestCheckResourceAttr(resourceGPFSACLName, "id", name),
					resource.TestCheckResourceAttr(resourceGPFSACLName, "statement.#", "1"),
				),
			},
			{
				Config: testAccGPFSACLV3Config(name, false, "Read"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					checkGPFSACLGrants(&acl, 1, "Read"),
					resource.TestCheckResourceAttr(resourceGPFSACLName, "id", name),
				),
			},
			{
				Config:   testAccGPFSACLV3Config(name, false, "Read"),
				PlanOnly: true,
			},
		},
	})
}

func testAccGPFSACLV3Config(name string, both bool, action string) string {
	first := ""
	if both {
		first = gpfsACLFirstStatement
	}
	return fmt.Sprintf(gpfsACLTestConfig, name, first, action)
}

const gpfsACLFirstStatement = `
  statement {
    vpc_id = opentelekomcloud_vpc_v1.test[0].id
    sid    = "first-vpc"
  }
`

const gpfsACLTestConfig = `
resource "opentelekomcloud_gpfs_file_system_v3" "fs" {
  name = "%[1]s"
}

resource "opentelekomcloud_vpc_v1" "test" {
  count = 2
  name  = "%[1]s-${count.index}"
  cidr  = "192.168.${count.index}.0/24"
}

resource "opentelekomcloud_gpfs_acl_v3" "acl" {
  fs_name = opentelekomcloud_gpfs_file_system_v3.fs.name
  %[2]s
  statement {
    vpc_id = opentelekomcloud_vpc_v1.test[1].id
    action = "%[3]s"
  }
}
`
