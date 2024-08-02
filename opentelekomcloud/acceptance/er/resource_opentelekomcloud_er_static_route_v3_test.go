package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/route"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/routes"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getStaticRouteFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ErV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return route.Get(client, state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccStaticRoute_basic(t *testing.T) {
	var (
		obj routes.Route

		sourceSelfResName = "opentelekomcloud_er_static_route_v3.source_self"
		destSelfResName   = "opentelekomcloud_er_static_route_v3.destination_self"
		crossVpcResName   = "opentelekomcloud_er_static_route_v3.cross_vpc"
		name              = fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
		bgpAsNum          = acctest.RandIntRange(64512, 65534)

		sourceSelfRes = common.InitResourceCheck(sourceSelfResName, &obj, getStaticRouteFunc)
		destSelfRes   = common.InitResourceCheck(destSelfResName, &obj, getStaticRouteFunc)
		crossVpcRes   = common.InitResourceCheck(crossVpcResName, &obj, getStaticRouteFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      sourceSelfRes.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRoute_basic_step1(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					sourceSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "route_table_id",
						"opentelekomcloud_er_route_table_v3.source", "id"),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "destination",
						"opentelekomcloud_vpc_v1.source", "cidr"),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.source", "id"),
					resource.TestCheckResourceAttrSet(sourceSelfResName, "type"),
					resource.TestCheckResourceAttrSet(sourceSelfResName, "status"),
					resource.TestCheckResourceAttrSet(sourceSelfResName, "created_at"),
					resource.TestCheckResourceAttrSet(sourceSelfResName, "updated_at"),
					destSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(destSelfResName, "route_table_id",
						"opentelekomcloud_er_route_table_v3.destination", "id"),
					resource.TestCheckResourceAttrPair(destSelfResName, "destination",
						"opentelekomcloud_vpc_v1.destination", "cidr"),
					resource.TestCheckResourceAttrPair(destSelfResName, "attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.destination", "id"),
					crossVpcRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(crossVpcResName, "route_table_id",
						"opentelekomcloud_er_route_table_v3.source", "id"),
					resource.TestCheckResourceAttrPair(crossVpcResName, "destination",
						"opentelekomcloud_vpc_v1.destination", "cidr"),
					resource.TestCheckResourceAttrPair(crossVpcResName, "attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.source", "id"),
				),
			},
			{
				Config: testAccStaticRoute_basic_step2(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					sourceSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.destination", "id"),
					destSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(destSelfResName, "attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.source", "id"),
					crossVpcRes.CheckResourceExists(),
					resource.TestCheckResourceAttr(crossVpcResName, "is_blackhole", "true"),
				),
			},
			{
				ResourceName:      sourceSelfResName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteImportStateFunc(sourceSelfResName),
			},
			{
				ResourceName:      destSelfResName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteImportStateFunc(destSelfResName),
			},
			{
				ResourceName:      crossVpcResName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteImportStateFunc(crossVpcResName),
			},
		},
	})
}

func testAccStaticRouteImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var routeTableId, staticRouteId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER static route is not found in the tfstate", rsName)
		}
		routeTableId = rs.Primary.Attributes["route_table_id"]
		staticRouteId = rs.Primary.ID
		if routeTableId == "" || staticRouteId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<route_table_id>/<id>', but got '%s/%s'",
				routeTableId, staticRouteId)
		}
		return fmt.Sprintf("%s/%s", routeTableId, staticRouteId), nil
	}
}

func testAccStaticRoute_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
variable "base_vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "source" {
  name = "%[1]s_source"
  cidr = cidrsubnet(var.base_vpc_cidr, 2, 1)
}

resource "opentelekomcloud_vpc_v1" "destination" {
  name = "%[1]s_destination"
  cidr = cidrsubnet(var.base_vpc_cidr, 2, 2)
}

resource "opentelekomcloud_vpc_subnet_v1" "source" {
  vpc_id = opentelekomcloud_vpc_v1.source.id

  name       = "%[1]s_source"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.source.cidr, 2, 1)
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.source.cidr, 2, 1), 1)
}

resource "opentelekomcloud_vpc_subnet_v1" "destination" {
  vpc_id = opentelekomcloud_vpc_v1.destination.id

  name       = "%[1]s_destination"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.destination.cidr, 2, 1)
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.destination.cidr, 2, 1), 1)
}

resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]
  name               = "%[1]s"
  asn                = %[2]d
}

resource "opentelekomcloud_er_route_table_v3" "source" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = "%[1]s_source"
}

resource "opentelekomcloud_er_route_table_v3" "destination" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = "%[1]s_destination"
}

resource "opentelekomcloud_er_vpc_attachment_v3" "source" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.source.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.source.id
  name        = "%[1]s_source"
}

resource "opentelekomcloud_er_vpc_attachment_v3" "destination" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.destination.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.destination.id
  name        = "%[1]s_destination"
}
`, name, bgpAsNum)
}

func testAccStaticRoute_basic_step1(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_static_route_v3" "source_self" {
  route_table_id = opentelekomcloud_er_route_table_v3.source.id
  destination    = opentelekomcloud_vpc_v1.source.cidr
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.source.id
}

resource "opentelekomcloud_er_static_route_v3" "destination_self" {
  route_table_id = opentelekomcloud_er_route_table_v3.destination.id
  destination    = opentelekomcloud_vpc_v1.destination.cidr
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.destination.id
}

resource "opentelekomcloud_er_static_route_v3" "cross_vpc" {
  route_table_id = opentelekomcloud_er_route_table_v3.source.id
  destination    = opentelekomcloud_vpc_v1.destination.cidr
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.source.id
}
`, testAccStaticRoute_base(name, bgpAsNum))
}

func testAccStaticRoute_basic_step2(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_static_route_v3" "source_self" {
  route_table_id = opentelekomcloud_er_route_table_v3.source.id
  destination    = opentelekomcloud_vpc_v1.source.cidr
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.destination.id
}

resource "opentelekomcloud_er_static_route_v3" "destination_self" {
  route_table_id = opentelekomcloud_er_route_table_v3.destination.id
  destination    = opentelekomcloud_vpc_v1.destination.cidr
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.source.id
}

resource "opentelekomcloud_er_static_route_v3" "cross_vpc" {
  route_table_id = opentelekomcloud_er_route_table_v3.source.id
  destination    = opentelekomcloud_vpc_v1.destination.cidr
  is_blackhole   = true
}
`, testAccStaticRoute_base(name, bgpAsNum))
}
