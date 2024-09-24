package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectiongroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const pgResourceName = "opentelekomcloud_sdrs_protectiongroup_v1.group_1"

func getProtectionGroupResourceFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SdrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return protectiongroups.Get(client, state.Primary.ID)
}

func TestAccSdrsProtectionGroupV1_basic(t *testing.T) {
	var group protectiongroups.ServerGroupResponseInfo
	rc := common.InitResourceCheck(
		pgResourceName,
		&group,
		getProtectionGroupResourceFunc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectionGroupV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_1"),
				),
			},
			{
				Config: testAccSdrsProtectionGroupV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_updated"),
				),
			},
		},
	})
}

func TestAccSdrsProtectionGroupV1_enabling(t *testing.T) {
	var group protectiongroups.ServerGroupResponseInfo
	rc := common.InitResourceCheck(
		pgResourceName,
		&group,
		getProtectionGroupResourceFunc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectionGroupV1enableBasic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_1"),
				),
			},
			{
				Config: testAccSdrsProtectionGroupV1Enable,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_updated"),
					resource.TestCheckResourceAttr(pgResourceName, "enable", "true"),
				),
			},
			{
				Config: testAccSdrsProtectionGroupV1Disable,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(pgResourceName, "name", "group_updated"),
					resource.TestCheckResourceAttr(pgResourceName, "enable", "false"),
				),
			},
		},
	})
}

var testAccSdrsProtectionGroupV1Basic = fmt.Sprintf(`
%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  description              = "test description"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}
`, common.DataSourceSubnet)

var testAccSdrsProtectionGroupV1Update = fmt.Sprintf(`
%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_updated"
  description              = "test description"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}
`, common.DataSourceSubnet)

var testAccSdrsProtectionGroupV1enableBasic = fmt.Sprintf(`
%s

%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  description              = "test description"
  source_availability_zone = "eu-de-01"
  target_availability_zone = "eu-de-02"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  availability_zone = "%s"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  description          = "some interesting description"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true
}
`, common.DataSourceSubnet, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)

var testAccSdrsProtectionGroupV1Enable = fmt.Sprintf(`
%s

%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_updated"
  description              = "test description"
  source_availability_zone = "eu-de-01"
  target_availability_zone = "eu-de-02"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
  enable                   = true
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  availability_zone = "%s"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  description          = "some interesting description"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true
}
`, common.DataSourceSubnet, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)

var testAccSdrsProtectionGroupV1Disable = fmt.Sprintf(`
%s

%s

data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_updated"
  description              = "test description"
  source_availability_zone = "eu-de-01"
  target_availability_zone = "eu-de-02"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  dr_type                  = "migration"
  enable                   = false
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  nics {
    network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  availability_zone = "%s"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  description          = "some interesting description"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true
}
`, common.DataSourceSubnet, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)
