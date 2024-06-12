package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func getInstanceResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ErV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Enterprise Router client: %s", err)
	}

	return instance.Get(client, state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("er-acc-api%s", acctest.RandString(5))
	rName := "opentelekomcloud_er_instance_v3.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic_step1(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttr(rName, "enable_default_propagation", "true"),
					resource.TestCheckResourceAttr(rName, "enable_default_association", "false"),
					resource.TestCheckResourceAttr(rName, "auto_accept_shared_attachments", "true"),
				),
			},
			{
				Config: testInstance_basic_step2(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "description", "test-update"),
					resource.TestCheckResourceAttr(rName, "enable_default_propagation", "false"),
					resource.TestCheckResourceAttr(rName, "enable_default_association", "true"),
					resource.TestCheckResourceAttr(rName, "auto_accept_shared_attachments", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testInstance_basic_step1(name string, bgpAsNum int) string {
	return fmt.Sprintf(`




resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]

  name        = "%[1]s"
  asn         = %[2]d
  description = "test"

  enable_default_propagation     = true
  enable_default_association     = false
  auto_accept_shared_attachments = true
}
`, name, bgpAsNum)
}

func testInstance_basic_step2(name string, bgpAsNum int) string {
	return fmt.Sprintf(`


resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = ["eu-de-01", "eu-de-02"]

  name        = "%[1]s"
  asn         = %[2]d
  description = "test-update"

  enable_default_propagation     = false
  enable_default_association     = true
  auto_accept_shared_attachments = false
}
`, name, bgpAsNum)
}
