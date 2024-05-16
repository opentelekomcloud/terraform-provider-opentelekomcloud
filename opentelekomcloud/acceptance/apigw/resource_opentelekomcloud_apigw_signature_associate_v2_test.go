package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/key"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	resourceApigwSignatureBasic = "opentelekomcloud_apigw_signature_v2.basic"
	resourceApigwSignatureHmac  = "opentelekomcloud_apigw_signature_v2.hmac"
	resourceApigwSignatureAes   = "opentelekomcloud_apigw_signature_v2.aes"
)

func getKeySignatureFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	opts := key.ListBoundOpts{
		GatewayID: state.Primary.Attributes["gateway_id"],
		SignID:    state.Primary.Attributes["signature_id"],
	}
	resp, err := key.ListAPIBoundKeys(c, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccSignatureAssociate_basic(t *testing.T) {
	var (
		apiDetails []key.BindSignResp

		name   = fmt.Sprintf("apigw_acc_key%s", acctest.RandString(5))
		rName1 = "opentelekomcloud_apigw_signature_associate_v2.basic_bind"
		rName2 = "opentelekomcloud_apigw_signature_associate_v2.hmac_bind"
		rName3 = "opentelekomcloud_apigw_signature_associate_v2.aes_bind"

		rc1 = common.InitResourceCheck(rName1, &apiDetails, getKeySignatureFunc)
		rc2 = common.InitResourceCheck(rName2, &apiDetails, getKeySignatureFunc)
		rc3 = common.InitResourceCheck(rName3, &apiDetails, getKeySignatureFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignatureAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName1, "gateway_id", resourceGwName, "id"),
					resource.TestCheckResourceAttrPair(rName1, "signature_id", resourceApigwSignatureBasic, "id"),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName2, "gateway_id", resourceGwName, "id"),
					resource.TestCheckResourceAttrPair(rName2, "signature_id", resourceApigwSignatureHmac, "id"),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName3, "gateway_id", resourceGwName, "id"),
					resource.TestCheckResourceAttrPair(rName3, "signature_id", resourceApigwSignatureAes, "id"),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				Config: testAccSignatureAssociate_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName3),
			},
		},
	})
}

func testAccSignatureAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.Attributes["signature_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<signature_id>', but got '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["signature_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["signature_id"]), nil
	}
}

func testAccSignatureAssociate_base(name string) string {
	relatedConfig := testAccApigwApi_basic(testAccApigwApi_base(name), name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_environment_v2" "envs" {
  count = 6

  name        = "%[2]s_${count.index}"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub" {
  count = 6

  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.envs[count.index].id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_signature_v2" "basic" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_basic"
  type       = "basic"
}

resource "opentelekomcloud_apigw_signature_v2" "hmac" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_hmac"
  type       = "hmac"
}

resource "opentelekomcloud_apigw_signature_v2" "aes" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_aes"
  type       = "aes"
  algorithm  = "aes-128-cfb"
}
`, relatedConfig, name)
}

func testAccSignatureAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_signature_associate_v2" "basic_bind" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  signature_id = opentelekomcloud_apigw_signature_v2.basic.id
  publish_ids  = slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 0, 2)
}

resource "opentelekomcloud_apigw_signature_associate_v2" "hmac_bind" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  signature_id = opentelekomcloud_apigw_signature_v2.hmac.id
  publish_ids  = slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 2, 4)
}

resource "opentelekomcloud_apigw_signature_associate_v2" "aes_bind" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  signature_id = opentelekomcloud_apigw_signature_v2.aes.id
  publish_ids  = slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 4, 6)
}
`, testAccSignatureAssociate_base(name))
}

func testAccSignatureAssociate_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_signature_associate_v2" "basic_bind" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  signature_id = opentelekomcloud_apigw_signature_v2.basic.id
  publish_ids  = slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 1, 3)
}

resource "opentelekomcloud_apigw_signature_associate_v2" "hmac_bind" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  signature_id = opentelekomcloud_apigw_signature_v2.hmac.id
  publish_ids  = slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 3, 5)
}

resource "opentelekomcloud_apigw_signature_associate_v2" "aes_bind" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  signature_id = opentelekomcloud_apigw_signature_v2.aes.id
  publish_ids = setunion(slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 0, 1),
  slice(opentelekomcloud_apigw_api_publishment_v2.pub[*].publish_id, 5, 6))
}
`, testAccSignatureAssociate_base(name))
}
