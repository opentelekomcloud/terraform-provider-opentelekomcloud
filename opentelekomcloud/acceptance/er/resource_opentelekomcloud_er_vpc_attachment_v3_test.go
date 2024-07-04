package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/vpc"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getVpcAttachmentResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ErV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Enterprise Router client: %s", err)
	}

	return vpc.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccVpcAttachmentsV3_basic(t *testing.T) {
	var (
		obj        vpc.VpcAttachmentsResp
		rName      = "opentelekomcloud_er_vpc_attachment_v3.test"
		name       = fmt.Sprintf("er-acc-vpc-api%s", acctest.RandString(5))
		updateName = fmt.Sprintf("er-acc-vpc-api%s", acctest.RandString(5))
		bgpAsNum   = acctest.RandIntRange(64512, 65534)
	)

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getVpcAttachmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAttachment_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "opentelekomcloud_vpc_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "opentelekomcloud_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttr(rName, "auto_create_vpc_routes", "true"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckOutput("er_route_count", "3"),
				),
			},
			{
				Config: testVpcAttachment_basic_update(updateName, bgpAsNum),
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
				ImportStateIdFunc: testAccVpcAttachmentImportStateFunc(rName),
			},
		},
	})
}

func testAccVpcAttachmentImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, attachmentId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER attachment is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		attachmentId = rs.Primary.ID
		if instanceId == "" || attachmentId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, attachmentId)
		}
		return fmt.Sprintf("%s/%s", instanceId, attachmentId), nil
	}
}

func testVpcAttachment_base(name string, bgpAsNum int) string {
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
`, name, bgpAsNum)
}

func testVpcAttachment_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.test.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.test.id

  name                   = "%[2]s"
  description            = "Create by acc test"
  auto_create_vpc_routes = true
}

data "opentelekomcloud_vpc_route_table_v1" "test" {
  depends_on = [
    opentelekomcloud_er_vpc_attachment_v3.test
  ]

  vpc_id = opentelekomcloud_vpc_v1.test.id
  name   = "rtb-%[2]s"
}

output "er_route_count" {
  value = length([for route in data.opentelekomcloud_vpc_route_table_v1.test.route : route.type == "er"])
}
`, testVpcAttachment_base(name, bgpAsNum), name)
}

func testVpcAttachment_basic_update(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  vpc_id      = opentelekomcloud_vpc_v1.test.id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.test.id

  name                   = "%[2]s"
  auto_create_vpc_routes = true
}
`, testVpcAttachment_base(name, bgpAsNum), name)
}
