package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/route_table"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getRouteTableResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ErV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return route_table.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccRouteTable_basic(t *testing.T) {
	var (
		obj route_table.RouteTable

		rName      = "opentelekomcloud_er_route_table_v3.test"
		name       = fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
		updateName = fmt.Sprintf("er-acc-api%s-updated", acctest.RandString(5))
		baseConfig = testRouteTable_base(name)
	)

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getRouteTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRouteTable_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttrSet(rName, "is_default_association"),
					resource.TestCheckResourceAttrSet(rName, "is_default_propagation"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testRouteTable_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRouteTableImportStateFunc(rName),
			},
		},
	})
}

func testAccRouteTableImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER route table is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		routeTableId = rs.Primary.ID
		if instanceId == "" || routeTableId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, routeTableId)
		}
		return fmt.Sprintf("%s/%s", instanceId, routeTableId), nil
	}
}

func testRouteTable_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	return fmt.Sprintf(`

resource "opentelekomcloud_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "test" {
  vpc_id = opentelekomcloud_vpc_v1.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.test.cidr, 4, 1), 1)
}

resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]
  name               = "%[1]s"
  asn                = %[2]d

  # Enable default routes
  enable_default_propagation = true
  enable_default_association = true
}
`, name, bgpAsNum)
}

func testRouteTable_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_route_table_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = "%[2]s"
  description = "Create by acc test"
}

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.test.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.test.id
  name        = "%[2]s"
}

`, baseConfig, name)
}

func testRouteTable_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_route_table_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = "%[2]s"
}

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.test.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.test.id
  name        = "%[2]s"
}
`, baseConfig, name)
}
