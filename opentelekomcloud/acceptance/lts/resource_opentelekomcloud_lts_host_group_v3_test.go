package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	hg "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/host-groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getHostGroupResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v3 client: %s", err)
	}

	requestResp, err := hg.List(client, hg.ListOpts{})
	if err != nil {
		return nil, err
	}
	if len(requestResp.Result) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var groupResult *hg.HostGroupResponse
	for _, group := range requestResp.Result {
		if group.ID == state.Primary.ID {
			groupResult = &group
		}
	}
	if groupResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return groupResult, nil
}

func TestAccHostGroup_basic(t *testing.T) {
	var (
		group hg.HostGroupResponse
		rName = "opentelekomcloud_lts_host_group_v3.hg"
		name  = fmt.Sprintf("lts_group%s", acctest.RandString(3))
		rc    = common.InitResourceCheck(rName, &group, getHostGroupResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestLtsPreCheckLts(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "agent_access_type", "IP"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config:            testHostGroup_basic(name),
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testHostGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key_update", "value"),
				),
			},
			{
				Config: testHostGroup_updateRm(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update_again"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.#", "0"),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "0"),
				),
			},
		},
	})
}

func testHostGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance" {
  name              = "%[3]s"
  description       = "my_desc"
  availability_zone = "%[2]s"

  image_name = "Standard_Debian_11_latest"
  flavor_id  = "s3.large.2"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  tags = {
    lts = "acc-test"
  }

  user_data = <<-EOF
#!/bin/bash
set +o history
curl https://icagent-eu-de.obs.eu-de.otc.t-systems.com/ICAgent_linux/apm_agent_install.sh > apm_agent_install.sh && REGION=eu-de bash apm_agent_install.sh -ak %[5]s -sk %[6]s -region eu-de -projectid %[4]s -accessip lts-access.eu-de.otc.t-systems.com -obsdomain obs.eu-de.otc.t-systems.com
set -o history
  EOF

  stop_before_destroy = true
}

resource "opentelekomcloud_lts_host_group_v3" "hg" {
  name     = "%[3]s"
  type     = "linux"
  host_ids = [opentelekomcloud_compute_instance_v2.instance.id]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name, env.OS_PROJECT_ID, env.OS_ACCESS_KEY, env.OS_SECRET_KEY)
}

func testHostGroup_update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance" {
  name              = "%[3]s"
  description       = "my_desc"
  availability_zone = "%[2]s"

  image_name = "Standard_Debian_11_latest"
  flavor_id  = "s3.large.2"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  tags = {
    lts = "acc-test"
  }

  user_data = <<-EOF
#!/bin/bash
set +o history
curl https://icagent-eu-de.obs.eu-de.otc.t-systems.com/ICAgent_linux/apm_agent_install.sh > apm_agent_install.sh && REGION=eu-de bash apm_agent_install.sh -ak %[5]s -sk %[6]s -region eu-de -projectid %[4]s -accessip lts-access.eu-de.otc.t-systems.com -obsdomain obs.eu-de.otc.t-systems.com
set -o history
  EOF

  stop_before_destroy = true
}

resource "opentelekomcloud_lts_host_group_v3" "hg" {
  name     = "%[3]s-update"
  type     = "linux"
  host_ids = [opentelekomcloud_compute_instance_v2.instance.id]

  tags = {
    foo        = "bar_update"
    key_update = "value"
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name, env.OS_PROJECT_ID, env.OS_ACCESS_KEY, env.OS_SECRET_KEY)
}

func testHostGroup_updateRm(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance" {
  name              = "%[3]s"
  description       = "my_desc"
  availability_zone = "%[2]s"

  image_name = "Standard_Debian_11_latest"
  flavor_id  = "s3.large.2"

  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }

  tags = {
    lts = "acc-test"
  }

  user_data = <<-EOF
#!/bin/bash
set +o history
curl https://icagent-eu-de.obs.eu-de.otc.t-systems.com/ICAgent_linux/apm_agent_install.sh > apm_agent_install.sh && REGION=eu-de bash apm_agent_install.sh -ak %[5]s -sk %[6]s -region eu-de -projectid %[4]s -accessip lts-access.eu-de.otc.t-systems.com -obsdomain obs.eu-de.otc.t-systems.com
set -o history
  EOF

  stop_before_destroy = true
}

resource "opentelekomcloud_lts_host_group_v3" "hg" {
  name = "%[3]s-update_again"
  type = "linux"
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, name, env.OS_PROJECT_ID, env.OS_ACCESS_KEY, env.OS_SECRET_KEY)
}
