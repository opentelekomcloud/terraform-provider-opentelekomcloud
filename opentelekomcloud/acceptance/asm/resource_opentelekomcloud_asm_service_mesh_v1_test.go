package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/asm/v1/servicemesh"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const asmServiceMeshResourceName = "opentelekomcloud_asm_service_mesh_v1.mesh_1"

func getPrivateAsmServiceMesh(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.AsmV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 Client: %s", err)
	}
	return servicemesh.Get(client, state.Primary.ID)
}

func TestAccPrivateAsmServiceMeshV3_basic(t *testing.T) {
	clusterId := os.Getenv("OS_CLUSTER_ID")
	nodeId := os.Getenv("OS_NODE_ID")
	if clusterId == "" || nodeId == "" {
		t.Skip("OS_CLUSTER_ID or OS_NODE_ID is missing but is required for ASM test.")
	}
	var mesh servicemesh.ServiceMesh
	rc := common.InitResourceCheck(
		asmServiceMeshResourceName,
		&mesh,
		getPrivateAsmServiceMesh,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateAsmServiceMeshV3Basic(clusterId, nodeId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(asmServiceMeshResourceName, "name", "test-acc-asm-service-mesh"),
					resource.TestCheckResourceAttr(asmServiceMeshResourceName, "cluster_ids.0", clusterId),
				),
			},
			{
				ResourceName:      asmServiceMeshResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"clusters",
				},
			},
		},
	})
}

func testAccPrivateAsmServiceMeshV3Basic(clusterId, nodeId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {
  name    = "test-acc-asm-service-mesh"
  type    = "InCluster"
  version = "1.18.7-r5"

  clusters {
    cluster_id         = "%s"
    installation_nodes = ["%s"]
  }
}
`, clusterId, nodeId)
}
