package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	var (
		res = "data.opentelekomcloud_er_instances_v3.test"
		dc  = common.InitDataSourceCheck(res)

		byInstanceId   = "data.opentelekomcloud_er_instances_v3.filter_by_instance_id"
		dcByInstanceId = common.InitDataSourceCheck(byInstanceId)

		byName   = "data.opentelekomcloud_er_instances_v3.filter_by_name"
		dcByName = common.InitDataSourceCheck(byName)

		byStatus   = "data.opentelekomcloud_er_instances_v3.filter_by_status"
		dcByStatus = common.InitDataSourceCheck(byStatus)

		byTags   = "data.opentelekomcloud_er_instances_v3.filter_by_tags"
		dcByTags = common.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceErInstancesV3_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(res, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by parameter 'instance_id'.
					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byInstanceId, "instances.#", "1"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.id", "opentelekomcloud_er_instance_v3.test", "id"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.asn", "opentelekomcloud_er_instance_v3.test", "asn"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.name", "opentelekomcloud_er_instance_v3.test", "name"),
					resource.TestCheckResourceAttrPair(byInstanceId, "instances.0.description", "opentelekomcloud_er_instance_v3.test", "description"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.tags.key", "value"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.enable_default_propagation", "true"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.enable_default_association", "true"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.auto_accept_shared_attachments", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.default_propagation_route_table_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.default_association_route_table_id"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					// Filter by parameter 'name'.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by parameter 'status'.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by parameter 'tags'.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceErInstancesV3_base() string {
	var (
		name     = fmt.Sprintf("er-data-inst%s", acctest.RandString(5))
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]
  name               = "%[1]s"
  asn                = %[2]d
  description        = "Created by terraform"

  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, bgpAsNum)
}

func testAccDataSourceErInstancesV3_basic() string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_er_instances_v3" "test" {
  depends_on = [opentelekomcloud_er_instance_v3.test]
}

# Filter by instance ID
locals {
  instance_id = opentelekomcloud_er_instance_v3.test.id
}

data "opentelekomcloud_er_instances_v3" "filter_by_instance_id" {
  instance_id = local.instance_id
}

locals {
  instance_id_filter_result = [
    for v in data.opentelekomcloud_er_instances_v3.filter_by_instance_id.instances[*].id : v == local.instance_id
  ]
}

output "is_instance_id_filter_useful" {
  value = length(local.instance_id_filter_result) > 0 && alltrue(local.instance_id_filter_result)
}

# Filter by name
locals {
  instance_name = opentelekomcloud_er_instance_v3.test.name
}

data "opentelekomcloud_er_instances_v3" "filter_by_name" {
  depends_on = [opentelekomcloud_er_instance_v3.test]

  name = local.instance_name
}

locals {
  name_filter_result = [
    for v in data.opentelekomcloud_er_instances_v3.filter_by_name.instances[*].name : v == local.instance_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by status
locals {
  instance_status = opentelekomcloud_er_instance_v3.test.status
}

data "opentelekomcloud_er_instances_v3" "filter_by_status" {
  status = local.instance_status
}

locals {
  status_filter_result = [
    for v in data.opentelekomcloud_er_instances_v3.filter_by_status.instances[*].status : v == local.instance_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by tags
locals {
  instance_tags = opentelekomcloud_er_instance_v3.test.tags
}

data "opentelekomcloud_er_instances_v3" "filter_by_tags" {
  depends_on = [opentelekomcloud_er_instance_v3.test]

  tags = local.instance_tags
}

locals {
  tags_filter_result = [
    for v in data.opentelekomcloud_er_instances_v3.filter_by_tags.instances[*].tags : length(v) == length(local.instance_tags) &&
    length(v) == length(merge(v, local.instance_tags))
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataSourceErInstancesV3_base())
}
