package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/log"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getLbLogResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB v3 client: %s", err)
	}
	return log.Get(client, state.Primary.ID)
}

func TestAccLtsElb_basic(t *testing.T) {
	var instance log.Log
	rName := "opentelekomcloud_lb_lts_log_v3.log"
	name := fmt.Sprintf("lb_log%s", acctest.RandString(3))

	rc := common.InitResourceCheck(
		rName,
		&instance,
		getLbLogResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLbLogV3Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"opentelekomcloud_lts_stream_v2.stream_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"opentelekomcloud_lts_group_v2.group", "id"),
					resource.TestCheckResourceAttrPair(rName, "loadbalancer_id",
						"opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1", "id"),
				),
			},
			{
				Config: testAccLbLogV3Update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"opentelekomcloud_lts_stream_v2.stream_2", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"opentelekomcloud_lts_group_v2.group", "id"),
					resource.TestCheckResourceAttrPair(rName, "loadbalancer_id",
						"opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1", "id"),
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

func testAccLbLogV3Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "%[3]s"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%[2]s"]

  public_ip {
    ip_type              = "5_bgp"
    bandwidth_name       = "lb_band"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }

  deletion_protection = false

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[3]s"
  ttl_in_days = 30

  tags = {
    created_by = "terraform"
  }
}

resource "opentelekomcloud_lts_stream_v2" "stream_1" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[3]s_first"

  tags = {
    number     = "one"
    created_by = "terraform"
  }
}

resource "opentelekomcloud_lts_stream_v2" "stream_2" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[3]s_second"

  tags = {
    number     = "two"
    created_by = "terraform"
  }
}

resource "opentelekomcloud_lb_lts_log_v3" "log" {
  log_stream_id   = opentelekomcloud_lts_stream_v2.stream_1.id
  log_group_id    = opentelekomcloud_lts_group_v2.group.id
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
}


`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name)
}

func testAccLbLogV3Update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "%[3]s"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["%[2]s"]

  public_ip {
    ip_type              = "5_bgp"
    bandwidth_name       = "lb_band"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }

  deletion_protection = false

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[3]s"
  ttl_in_days = 30

  tags = {
    created_by = "terraform"
  }
}

resource "opentelekomcloud_lts_stream_v2" "stream_1" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[3]s_first"

  tags = {
    number     = "one"
    created_by = "terraform"
  }
}

resource "opentelekomcloud_lts_stream_v2" "stream_2" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[3]s_second"

  tags = {
    number     = "two"
    created_by = "terraform"
  }
}

resource "opentelekomcloud_lb_lts_log_v3" "log" {
  log_stream_id   = opentelekomcloud_lts_stream_v2.stream_2.id
  log_group_id    = opentelekomcloud_lts_group_v2.group.id
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
}


`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name)
}
