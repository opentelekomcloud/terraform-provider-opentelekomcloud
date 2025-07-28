package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceRouteTables_basic(t *testing.T) {
	var (
		all = "data.opentelekomcloud_er_route_tables_v3.test"
		dc  = common.InitDataSourceCheck(all)

		byRouteTableId   = "data.opentelekomcloud_er_route_tables_v3.filter_by_route_table_id"
		dcByRouteTableId = common.InitDataSourceCheck(byRouteTableId)

		byName   = "data.opentelekomcloud_er_route_tables_v3.filter_by_name"
		dcByName = common.InitDataSourceCheck(byName)

		byTags   = "data.opentelekomcloud_er_route_tables_v3.filter_by_tags"
		dcByTags = common.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTables_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "route_tables.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByRouteTableId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.id",
						"opentelekomcloud_er_route_table_v3.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.name",
						"opentelekomcloud_er_route_table_v3.test", "name"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.description",
						"opentelekomcloud_er_route_table_v3.test", "description"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.associations.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.associations.0.id",
						"opentelekomcloud_er_association_v3.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.associations.0.attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.associations.0.attachment_type", "vpc"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.propagations.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.propagations.0.id",
						"opentelekomcloud_er_propagation_v3.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.propagations.0.attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.propagations.0.attachment_type", "vpc"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.id",
						"opentelekomcloud_er_static_route_v3.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.destination",
						"opentelekomcloud_er_static_route_v3.test", "destination"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.0.is_blackhole", "false"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.0.attachments.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.attachments.0.attachment_id",
						"opentelekomcloud_er_vpc_attachment_v3.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.0.attachments.0.attachment_type", "vpc"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.attachments.0.resource_id",
						"opentelekomcloud_vpc_v1.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.tags.key", "value"),
					resource.TestCheckOutput("is_route_table_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRouteTables_base() string {
	var (
		name     = fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

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
  availability_zones = ["eu-de-02"]

  name = "%[2]s"
  asn  = %[3]d
}

resource "opentelekomcloud_er_route_table_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = "%[2]s"
  description = "Created by script"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.test.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.test.id
  name        = "%[2]s"
}

resource "opentelekomcloud_er_static_route_v3" "test" {
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  destination    = opentelekomcloud_vpc_v1.test.cidr
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.test.id
}

resource "opentelekomcloud_er_association_v3" "test" {
  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.test.id
}

resource "opentelekomcloud_er_propagation_v3" "test" {
  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.test.id
}
`, name, name, bgpAsNum)
}

func testAccDataSourceRouteTables_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_er_route_tables_v3" "test" {
  depends_on = [
    opentelekomcloud_er_static_route_v3.test,
    opentelekomcloud_er_association_v3.test,
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id = opentelekomcloud_er_instance_v3.test.id
}

# Filter by route table ID
locals {
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
}

data "opentelekomcloud_er_route_tables_v3" "filter_by_route_table_id" {
  depends_on = [
    opentelekomcloud_er_static_route_v3.test,
    opentelekomcloud_er_association_v3.test,
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = local.route_table_id
}

locals {
  route_table_id_filter_result = [
    for v in data.opentelekomcloud_er_route_tables_v3.filter_by_route_table_id.route_tables[*].id : v == local.route_table_id
  ]
}

output "is_route_table_id_filter_useful" {
  value = length(local.route_table_id_filter_result) > 0 && alltrue(local.route_table_id_filter_result)
}

# Filter by name
locals {
  route_table_name = opentelekomcloud_er_route_table_v3.test.name
}

data "opentelekomcloud_er_route_tables_v3" "filter_by_name" {
  depends_on = [
    opentelekomcloud_er_static_route_v3.test,
    opentelekomcloud_er_association_v3.test,
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = local.route_table_name
}

locals {
  name_filter_result = [
    for v in data.opentelekomcloud_er_route_tables_v3.filter_by_route_table_id.route_tables[*].name : v == local.route_table_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by tags
locals {
  route_table_tags = opentelekomcloud_er_route_table_v3.test.tags
}

data "opentelekomcloud_er_route_tables_v3" "filter_by_tags" {
  depends_on = [
    opentelekomcloud_er_static_route_v3.test,
    opentelekomcloud_er_association_v3.test,
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id = opentelekomcloud_er_instance_v3.test.id
  tags        = local.route_table_tags
}

locals {
  tags_filter_result = [
    for v in data.opentelekomcloud_er_route_tables_v3.filter_by_tags.route_tables[*].tags : length(v) == length(local.route_table_tags) &&
    length(v) == length(merge(v, local.route_table_tags))
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataSourceRouteTables_base())
}
