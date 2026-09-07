package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/css/v1/load_balancer"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

// getCssLoadBalancerFunc loads the CSS load balancer config from the remote API
func getCssLoadBalancerFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CssV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	clusterID := os.Getenv("OS_CSS_CLUSTER_ID")
	if clusterID == "" {
		return nil, fmt.Errorf("environment variable OS_CSS_CLUSTER_ID is not set")
	}

	return load_balancer.Get(client, clusterID)
}

// TestAccCssLoadBalancerConfiguration_basic performs acceptance testing for CSS load balancer configuration
func TestAccCssLoadBalancerConfiguration_basic(t *testing.T) {
	clusterID := os.Getenv("OS_CSS_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("OS_CSS_CLUSTER_ID env var is not set")
	}

	var obj load_balancer.LoadBalancerResp
	rName := "opentelekomcloud_css_loadbalancer_v1.css_lb"
	rc := common.InitResourceCheck(
		rName,
		&obj,
		getCssLoadBalancerFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkCssLoadBalancerDisable,
		Steps: []resource.TestStep{
			{
				Config: testCssLoadBalancerV1_basic(clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", clusterID),
					resource.TestCheckResourceAttr(rName, "agency", "css_upgrade_agency"),
				),
			},
			{
				Config: testCssLoadBalancerV1_update(clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", clusterID),
					resource.TestCheckResourceAttr(rName, "agency", "css_upgrade_agency"),
				),
			},
			{
				Config: testCssLoadBalancerV1_update2(clusterID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", clusterID),
					resource.TestCheckResourceAttr(rName, "agency", "css_upgrade_agency"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// Optionally provide ImportStateIdFunc if import ID differs from default
			},
		},
	})
}

// testCssLoadBalancerV1_basic returns Terraform config string for the basic configuration
func testCssLoadBalancerV1_basic(clusterID string) string {
	return fmt.Sprintf(`
variable "elb_id" {}
variable "server_cert_id" {}

resource "opentelekomcloud_css_loadbalancer_v1" "css_lb" {
  cluster_id = "%s"
  elb_id     = var.elb_id
  agency     = "css_upgrade_agency"

  listener {
    protocol       = "HTTPS"
    protocol_port  = 443
    server_cert_id = var.server_cert_id
  }
}
`, clusterID)
}

// testCssLoadBalancerV1_update returns updated Terraform config (without listener block)
func testCssLoadBalancerV1_update(clusterID string) string {
	return fmt.Sprintf(`
variable "elb_id" {}

resource "opentelekomcloud_css_loadbalancer_v1" "css_lb" {
  cluster_id = "%s"
  elb_id     = var.elb_id
  agency     = "css_upgrade_agency"
}
`, clusterID)
}

// testCssLoadBalancerV1_update2 returns updated Terraform config (re-adding listener block)
func testCssLoadBalancerV1_update2(clusterID string) string {
	return fmt.Sprintf(`
variable "elb_id" {}
variable "server_cert_id" {}

resource "opentelekomcloud_css_loadbalancer_v1" "css_lb" {
  cluster_id = "%s"
  elb_id     = var.elb_id
  agency     = "css_upgrade_agency"

  listener {
    protocol       = "HTTPS"
    protocol_port  = 443
    server_cert_id = var.server_cert_id
  }
}
`, clusterID)
}

// checkCssLoadBalancerDisable verifies that the load balancer was disabled or removed
func checkCssLoadBalancerDisable(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)

	client, err := config.CssV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	clusterID := os.Getenv("OS_CSS_CLUSTER_ID")
	if clusterID == "" {
		return fmt.Errorf("environment variable OS_CSS_CLUSTER_ID is not set")
	}

	details, err := load_balancer.Get(client, clusterID)
	if err != nil {
		// Assuming a 404 or not found means it's destroyed
		return nil
	}

	if details.Enabled {
		return fmt.Errorf("load_balancer is still enabled for cluster %s", clusterID)
	}

	return nil
}
