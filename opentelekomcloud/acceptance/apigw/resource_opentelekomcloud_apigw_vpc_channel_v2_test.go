package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/channel"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwChannelName = "opentelekomcloud_apigw_vpc_channel_v2.channel"

func getChannelFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return channel.Get(client, state.Primary.Attributes["gateway_id"], state.Primary.ID)
}

func TestAccChannel_basic(t *testing.T) {
	var ch channel.ChannelResp

	name := fmt.Sprintf("apigw_acc_channel%s", acctest.RandString(5))
	updateName := fmt.Sprintf("apigw_acc_channel_up%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceApigwChannelName,
		&ch,
		getChannelFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccChannel_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "name", name),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "port", "80"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "lb_algorithm", "1"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member_type", "ecs"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "type", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.threshold_normal", "1"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.threshold_abnormal", "1"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.interval", "1"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.timeout", "1"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.path", ""),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.method", ""),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.port", "0"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.http_codes", ""),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member.#", "1"),
				),
			},
			{
				Config: testAccChannel_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "name", updateName),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "port", "8000"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "lb_algorithm", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member_type", "ecs"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "type", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.threshold_normal", "10"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.threshold_abnormal", "10"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.interval", "300"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.path", "/terraform/"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.method", "HEAD"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.http_codes", "201,202,303-404"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member.#", "2"),
				),
			},
			{
				ResourceName:      resourceApigwChannelName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(),
			},
		},
	})
}

func TestAccChannel_eipMembers(t *testing.T) {
	var ch channel.ChannelResp
	name := fmt.Sprintf("apigw_acc_channel%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceApigwChannelName,
		&ch,
		getChannelFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccChannel_eipMembers(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "name", name),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "port", "80"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "lb_algorithm", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member_type", "ip"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "type", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.threshold_normal", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.threshold_abnormal", "2"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.interval", "60"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.timeout", "10"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.path", "/"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.method", "HEAD"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "health_check.0.http_codes", "201,202,303-404"),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member.#", "1"),
				),
			},
			{
				Config: testAccChannel_eipMembers_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "name", name),
					resource.TestCheckResourceAttr(resourceApigwChannelName, "member.#", "2"),
				),
			},
			{
				ResourceName:      resourceApigwChannelName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(),
			},
		},
	})
}

func testAccChannelImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceApigwChannelName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceApigwChannelName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["gateway_id"],
				rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
	}
}

func testAccChannel_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_compute_instance_v2" "instance" {
  count = 1

  name        = format("%[2]s-%%d", count.index)
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_apigw_vpc_channel_v2" "channel" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  port         = 80
  lb_algorithm = 1
  member_type  = "ecs"
  type         = 2

  health_check {
    protocol           = "TCP"
    threshold_normal   = 1 # minimum value
    threshold_abnormal = 1 # minimum value
    interval           = 1 # minimum value
    timeout            = 1 # minimum value
  }

  dynamic "member" {
    for_each = opentelekomcloud_compute_instance_v2.instance[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), name)
}

func testAccChannel_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_compute_instance_v2" "instance" {
  count = 2

  name        = format("%[2]s-%%d", count.index)
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_apigw_vpc_channel_v2" "channel" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  port         = 8000
  lb_algorithm = 2
  member_type  = "ecs"
  type         = 2

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 10  # maximum value
    threshold_abnormal = 10  # maximum value
    interval           = 300 # maximum value
    timeout            = 30  # maximum value
    path               = "/terraform/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  dynamic "member" {
    for_each = opentelekomcloud_compute_instance_v2.instance[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), name)
}

func testAccChannel_eipMembers(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_vpc_eip_v1" "eip" {
  count = 1

  publicip {
    type = "5_bgp"
    name = "my_ip"
  }
  bandwidth {
    name        = "acc-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_apigw_vpc_channel_v2" "channel" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  port         = 80
  lb_algorithm = 2
  member_type  = "ip"
  type         = 2

  health_check {
    protocol           = "HTTP"
    threshold_normal   = 2
    threshold_abnormal = 2
    interval           = 60
    timeout            = 10
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  dynamic "member" {
    for_each = opentelekomcloud_vpc_eip_v1.eip[*].publicip.0.ip_address

    content {
      host = member.value
    }
  }
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), name)
}

func testAccChannel_eipMembers_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_vpc_eip_v1" "eip" {
  count = 2

  publicip {
    type = "5_bgp"
    name = "my_ip"
  }
  bandwidth {
    name        = "acc-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_apigw_vpc_channel_v2" "channel" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  port         = 80
  lb_algorithm = 2
  member_type  = "ip"
  type         = 2

  health_check {
    protocol           = "HTTP"
    threshold_normal   = 2
    threshold_abnormal = 2
    interval           = 60
    timeout            = 10
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  dynamic "member" {
    for_each = opentelekomcloud_vpc_eip_v1.eip[*].publicip.0.ip_address

    content {
      host = member.value
    }
  }
}
`, testAccAPIGWv2GatewayBasic("gateway-"+name), name)
}
