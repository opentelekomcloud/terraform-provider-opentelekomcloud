package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"

	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/dataarts/v1.1/cluster"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceMigrationsClusterName = "opentelekomcloud_dataarts_migrations_cluster_v1.cluster"

func getClusterFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DataArtsMigrationsV1Client(cfg.ProjectName(env.OS_PROJECT_ID))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud Dataarts Migrations client: %s", err)
	}
	return apis.Get(c, state.Primary.ID)
}

func TestAccDAMigrationsCluster_basic(t *testing.T) {
	var (
		api   apis.ClusterQuery
		rName = resourceMigrationsClusterName
		name  = fmt.Sprintf("da_mig_cluster_acc_api_%s", acctest.RandString(5))
		vpcID = env.OS_VPC_ID
	)

	rc := common.InitResourceCheck(
		rName,
		&api,
		getClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDAMigrationsCluster_basic(vpcID, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "language", "eng"),
					resource.TestCheckResourceAttr(rName, "email", "anyemail@example.com"),
					resource.TestCheckResourceAttr(rName, "phone_number", "123456789"),
					resource.TestCheckResourceAttr(rName, "is_schedule_boot_off", "false"),
					resource.TestCheckResourceAttr(rName, "is_auto_off", "false"),
				),
			},
			// {
			// 	ResourceName:      rName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateIdFunc: testAccApiResourceImportStateFunc(),
			// },
		},
	})
}

func testAccApiResourceImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "opentelekomcloud_apigw_api_v2.api"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{gateway_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccDAMigrationsCluster_basic(vpcID, name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dataarts_migrations_cluster_v1" "cluster" {
  language                     = "en"
  auto_remind                  = false
  phone_number                 = "123456789"
  email                        = "anyemail@example.com"

  #schedule_boot_time     = "10/10/2024"
  is_schedule_boot_off     = false
  #schedule_off_time     = "10/11/2024"
  vpc_id = "%s"
  name = "%s"
  #is_auto_off     = false

  datastore {
    type      = "cdm"
    version = "2.10.0.100"
  }

  #extended_properties {
	#workspace = "default"
    #resource = "oneone"
    #trial = "no"
  #}

  #sys_tags {
	#  value = "onetag"
	 # key = "keyforvalue"
  #}

  instances {
	  availability_zone = "eu-de-03"
	  #flavor = "5ddb1071-c5d7-40e0-a874-8a032e81a697"
	  flavor_id = "a79fd5ae-1833-448a-88e8-3ea2b913e1f6"
	  type = "cdm"
	  nics {
	      security_group = "12fcdc62-9bff-4df0-ac14-5258050d004b"
	      net = "03c5f385-17a0-4312-9e10-9b084edc18a1"
      }
  }
}
`, vpcID, name)
}
