package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/propagation"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/er"
)

func getPropagationResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ErV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Enterprise Router client: %s", err)
	}

	return er.QueryPropagationById(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccPropagation_basic(t *testing.T) {
	var (
		obj propagation.Propagation

		rName = "opentelekomcloud_er_propagation_v3.test"
		name  = fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
	)

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getPropagationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPropagation_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"opentelekomcloud_er_instance_v3.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"opentelekomcloud_er_route_table_v3.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpc"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPropagationImportStateFunc(rName),
			},
		},
	})
}

func testAccPropagationImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId, propagationId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER propagation is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		routeTableId = rs.Primary.Attributes["route_table_id"]
		propagationId = rs.Primary.ID
		if instanceId == "" || routeTableId == "" || propagationId == "" {
			return "", fmt.Errorf("some import IDs are missing, want "+
				"'<instance_id>/<route_table_id>/<id>', but got '%s/%s/%s'",
				instanceId, routeTableId, propagationId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, routeTableId, propagationId), nil
	}
}

func testAccPropagation_base(name string) string {
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

  name = "%[1]s"
  asn  = %[2]d
}

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.test.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.test.id

  name                   = "%[1]s"
  auto_create_vpc_routes = true
}

resource "opentelekomcloud_er_route_table_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id

  name = "%[1]s"
}
`, name, bgpAsNum)
}

func testAccPropagation_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_propagation_v3" "test" {
  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.test.id
}
`, testAccPropagation_base(name))
}
