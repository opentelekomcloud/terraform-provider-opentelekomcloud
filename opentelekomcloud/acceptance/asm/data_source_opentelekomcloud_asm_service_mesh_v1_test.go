package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const asmServiceMeshDataSourceName = "data.opentelekomcloud_asm_service_mesh_v1.mesh_1"

func TestAccPrivateAsmServiceMeshV3DS_basic(t *testing.T) {
	meshId := os.Getenv("OS_SERVICE_MESH_ID")
	if meshId == "" {
		t.Skip("OS_SERVICE_MESH_ID is missing but is required for ASM test.")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateAsmServiceMeshV3DS_Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAsmServiceMeshV1DataSourceID(asmServiceMeshDataSourceName),
					resource.TestCheckResourceAttr(asmServiceMeshDataSourceName, "service_meshes.0.id", meshId),
				),
			},
		},
	})
}

func testAccCheckAsmServiceMeshV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Private NAT Gateway data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Private NAT Gateway data source ID not set")
		}

		return nil
	}
}

var testAccPrivateAsmServiceMeshV3DS_Basic = `
data "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {}
`

func TestAccPrivateAsmServiceMeshV3DS_id(t *testing.T) {
	meshId := os.Getenv("OS_SERVICE_MESH_ID")
	if meshId == "" {
		t.Skip("OS_SERVICE_MESH_ID is missing but is required for ASM test.")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateAsmServiceMeshV3DS_id(meshId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(asmServiceMeshDataSourceName, "service_meshes.0.id", meshId),
				),
			},
		},
	})
}

func testAccPrivateAsmServiceMeshV3DS_id(meshId string) string {
	return fmt.Sprintf(`
data "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {
  id = "%s"
}
`, meshId)
}
