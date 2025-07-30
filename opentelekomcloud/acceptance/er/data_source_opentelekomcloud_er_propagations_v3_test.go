package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourcePropagations_basic(t *testing.T) {
	var (
		name       = fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
		baseConfig = testAccDataSourcePropagations_base(name)

		all = "data.opentelekomcloud_er_propagations_v3.test"
		dc  = common.InitDataSourceCheck(all)

		byAttachmentId   = "data.opentelekomcloud_er_propagations_v3.filter_by_attachment_id"
		dcByAttachmentId = common.InitDataSourceCheck(byAttachmentId)

		byAttachmentType   = "data.opentelekomcloud_er_propagations_v3.filter_by_attachment_type"
		dcByAttachmentType = common.InitDataSourceCheck(byAttachmentType)

		byStatus   = "data.opentelekomcloud_er_propagations_v3.filter_by_status"
		dcByStatus = common.InitDataSourceCheck(byStatus)

		byNotFoundInstanceId   = "data.opentelekomcloud_er_propagations_v3.instance_id_not_found"
		dcByNotFoundInstanceId = common.InitDataSourceCheck(byNotFoundInstanceId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePropagations_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "propagations.#"),
					resource.TestCheckResourceAttrSet(all, "propagations.0.resource_id"),
					dcByAttachmentId.CheckResourceExists(),
					resource.TestCheckOutput("is_attachment_id_filter_useful", "true"),
					dcByAttachmentType.CheckResourceExists(),
					resource.TestCheckOutput("is_attachment_type_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
			{
				Config: testAccDataSourcePropagations_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dcByNotFoundInstanceId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byNotFoundInstanceId, "propagations.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourcePropagations_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_propagation_v3" "test" {
  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  attachment_id  = opentelekomcloud_er_vpc_attachment_v3.test.id
}
`, testAccPropagation_base(name))
}

func testAccDataSourcePropagations_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_er_propagations_v3" "test" {
  depends_on = [
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
}

locals {
  attachment_id = data.opentelekomcloud_er_propagations_v3.test.propagations[0].attachment_id
}

data "opentelekomcloud_er_propagations_v3" "filter_by_attachment_id" {
  depends_on = [
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  attachment_id  = local.attachment_id
}

locals {
  attachment_id_filter_result = [
    for v in data.opentelekomcloud_er_propagations_v3.filter_by_attachment_id.propagations[*].attachment_id :
    v == local.attachment_id
  ]
}

output "is_attachment_id_filter_useful" {
  value = alltrue(local.attachment_id_filter_result) && length(local.attachment_id_filter_result) > 0
}

locals {
  attachment_type = data.opentelekomcloud_er_propagations_v3.test.propagations[0].attachment_type
}

data "opentelekomcloud_er_propagations_v3" "filter_by_attachment_type" {
  depends_on = [
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id     = opentelekomcloud_er_instance_v3.test.id
  route_table_id  = opentelekomcloud_er_route_table_v3.test.id
  attachment_type = local.attachment_type
}

locals {
  attachment_type_filter_result = [
    for v in data.opentelekomcloud_er_propagations_v3.filter_by_attachment_type.propagations[*].attachment_type :
    v == local.attachment_type
  ]
}

output "is_attachment_type_filter_useful" {
  value = alltrue(local.attachment_type_filter_result) && length(local.attachment_type_filter_result) > 0
}

locals {
  status = data.opentelekomcloud_er_propagations_v3.test.propagations[0].status
}

data "opentelekomcloud_er_propagations_v3" "filter_by_status" {
  depends_on = [
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id    = opentelekomcloud_er_instance_v3.test.id
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
  status         = local.status
}

locals {
  status_filter_result = [
    for v in data.opentelekomcloud_er_propagations_v3.filter_by_status.propagations[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, baseConfig)
}

func testAccDataSourcePropagations_basic_step2(baseConfig string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_er_propagations_v3" "instance_id_not_found" {
  depends_on = [
    opentelekomcloud_er_propagation_v3.test,
  ]

  instance_id    = "%[2]s"
  route_table_id = opentelekomcloud_er_route_table_v3.test.id
}
`, baseConfig, randUUID)
}
