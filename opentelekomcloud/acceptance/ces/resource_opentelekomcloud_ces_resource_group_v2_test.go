package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/resourcegroups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceResourceGroupV2Name = "opentelekomcloud_ces_resource_group_v2.test"

func getResourceGroupV2Func(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CesV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CES v2 client: %s", err)
	}
	group, err := resourcegroups.Get(c, state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return group, nil
}

func TestAccCESResourceGroupV2_basic(t *testing.T) {
	var (
		group resourcegroups.ResourceGroupDetail
		rName = resourceResourceGroupV2Name
		name  = acctest.RandomWithPrefix("ces-rg")
	)

	rc := common.InitResourceCheck(
		rName,
		&group,
		getResourceGroupV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroupV2_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testResourceGroupV2_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources",
				},
			},
		},
	})
}

func TestAccCESResourceGroupV2_tags(t *testing.T) {
	var (
		group resourcegroups.ResourceGroupDetail
		rName = resourceResourceGroupV2Name
		name  = acctest.RandomWithPrefix("ces-rg")
	)

	rc := common.InitResourceCheck(
		rName,
		&group,
		getResourceGroupV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroupV2_tags(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "TAG"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testResourceGroupV2_tags_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(rName, "tags.foo_update", "bar_update"),
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

func TestAccCESResourceGroupV2_eps(t *testing.T) {
	var (
		group resourcegroups.ResourceGroupDetail
		rName = resourceResourceGroupV2Name
		name  = acctest.RandomWithPrefix("ces-rg")
	)

	rc := common.InitResourceCheck(
		rName,
		&group,
		getResourceGroupV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroupV2_eps(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "EPS"),
					resource.TestCheckResourceAttr(rName, "associated_eps_ids.0", "0"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
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

func testResourceGroupV2_base(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "test" {
  name        = "ecs-%[2]s"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet, name)
}

func testResourceGroupV2_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name = "%s"

  resources {
    namespace = "SYS.ECS"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.test.id
    }
  }
}
`, testResourceGroupV2_base(name), name)
}

func testResourceGroupV2_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name = "%s-update"

  resources {
    namespace = "SYS.EVS"
    dimensions {
      name  = "disk_name"
      value = "${opentelekomcloud_compute_instance_v2.test.id}-sda"
    }
  }
}
`, testResourceGroupV2_base(name), name)
}

func testResourceGroupV2_tags(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name = "%s"
  type = "TAG"
  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name)
}

func testResourceGroupV2_tags_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name = "%s-update"
  type = "TAG"
  tags = {
    key        = "value_update"
    foo_update = "bar_update"
  }
}
`, name)
}

func testResourceGroupV2_eps(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_resource_group_v2" "test" {
  name               = "%s"
  type               = "EPS"
  associated_eps_ids = ["0"]
}
`, name)
}
