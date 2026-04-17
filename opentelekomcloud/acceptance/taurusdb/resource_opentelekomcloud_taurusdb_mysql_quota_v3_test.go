package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/quota"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getResourceQuota(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.TaurusDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	opts := quota.ListQuotasOpts{}
	resp, err := quota.ListQuotas(client, opts)
	if err != nil {
		return nil, err
	}

	epID := state.Primary.Attributes["enterprise_project_id"]
	for _, q := range resp.QuotaList {
		if q.EnterpriseProjectId == epID {
			return q, nil
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccTaurusDBQuota_basic(t *testing.T) {
	var obj interface{}
	resourceName := "opentelekomcloud_taurusdb_mysql_quota_v3.test"

	rc := common.InitResourceCheck(
		resourceName,
		&obj,
		getResourceQuota,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBQuota_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", env.OS_ENTERPRISE_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_name", env.OS_ENTERPRISE_PROJECT_NAME),
					resource.TestCheckResourceAttr(resourceName, "instance_quota", "10"),
					resource.TestCheckResourceAttr(resourceName, "vcpus_quota", "20"),
					resource.TestCheckResourceAttr(resourceName, "ram_quota", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_instance_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_vcpus_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_ram_quota"),
				),
			},
			{
				Config: testAccTaurusDBQuota_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", env.OS_ENTERPRISE_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "instance_quota", "0"),
					resource.TestCheckResourceAttr(resourceName, "vcpus_quota", "50"),
					resource.TestCheckResourceAttr(resourceName, "ram_quota", "-1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTaurusDBQuota_basic() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_taurusdb_mysql_quota_v3" "test" {
  enterprise_project_id   = "%s"
  enterprise_project_name = "%s"
  instance_quota          = 10
  vcpus_quota             = 20
  ram_quota               = 30
}`, env.OS_ENTERPRISE_PROJECT_ID, env.OS_ENTERPRISE_PROJECT_NAME)
}

func testAccTaurusDBQuota_basic_update() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_taurusdb_mysql_quota_v3" "test" {
  enterprise_project_id   = "%s"
  enterprise_project_name = "%s"
  instance_quota          = 0
  vcpus_quota             = 50
}`, env.OS_ENTERPRISE_PROJECT_ID, env.OS_ENTERPRISE_PROJECT_NAME)
}
