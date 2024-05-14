package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/key"
)

func getSignatureFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opts := key.ListOpts{
		GatewayID:   state.Primary.Attributes["gateway_id"],
		SignatureID: state.Primary.ID,
	}
	signs, err := key.List(client, opts)
	if err != nil {
		return nil, err
	}
	if len(signs) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return signs[0], nil
}

func TestAccSignature_basic(t *testing.T) {
	var signature key.SignKeyResp
	rName1 := "opentelekomcloud_apigw_signature_v2.with_key"
	rName2 := "opentelekomcloud_apigw_signature_v2.without_key"
	rName3 := "opentelekomcloud_apigw_signature_v2.hmac_with_key"
	rName4 := "opentelekomcloud_apigw_signature_v2.hmac_without_key"
	rName5 := "opentelekomcloud_apigw_signature_v2.aes_with_key"
	rName6 := "opentelekomcloud_apigw_signature_v2.aes_without_key"
	name := "apigw_acc_key"

	signKey := acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
	revSignKey := common.Reverse(signKey)
	signSecret := acctest.RandStringFromCharSet(16, acctest.CharSetAlphaNum)
	revSignSecret := common.Reverse(signSecret)
	rc1 := common.InitResourceCheck(rName1, &signature, getSignatureFunc)
	rc2 := common.InitResourceCheck(rName2, &signature, getSignatureFunc)
	rc3 := common.InitResourceCheck(rName3, &signature, getSignatureFunc)
	rc4 := common.InitResourceCheck(rName4, &signature, getSignatureFunc)
	rc5 := common.InitResourceCheck(rName5, &signature, getSignatureFunc)
	rc6 := common.InitResourceCheck(rName6, &signature, getSignatureFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignature_basic(name, signKey, signSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key"),
					resource.TestCheckResourceAttr(rName1, "type", "basic"),
					resource.TestCheckResourceAttr(rName1, "key", signKey),
					resource.TestCheckResourceAttr(rName1, "secret", signSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key"),
					resource.TestCheckResourceAttr(rName2, "type", "basic"),
					resource.TestCheckResourceAttrSet(rName2, "key"),
					resource.TestCheckResourceAttrSet(rName2, "secret"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "name", name+"_hmac_with_key"),
					resource.TestCheckResourceAttr(rName3, "type", "hmac"),
					resource.TestCheckResourceAttrSet(rName3, "key"),
					resource.TestCheckResourceAttrSet(rName3, "secret"),
					rc4.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName4, "name", name+"_hmac_without_key"),
					resource.TestCheckResourceAttr(rName4, "type", "hmac"),
					resource.TestCheckResourceAttrSet(rName4, "key"),
					resource.TestCheckResourceAttrSet(rName4, "secret"),
					rc5.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName5, "name", name+"_aes_with_key"),
					resource.TestCheckResourceAttr(rName5, "type", "aes"),
					resource.TestCheckResourceAttr(rName5, "algorithm", "aes-128-cfb"),
					resource.TestCheckResourceAttr(rName5, "key", signKey),
					resource.TestCheckResourceAttr(rName5, "secret", signSecret),
					rc6.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName6, "name", name+"_aes_without_key"),
					resource.TestCheckResourceAttr(rName6, "type", "aes"),
					resource.TestCheckResourceAttrSet(rName6, "key"),
					resource.TestCheckResourceAttrSet(rName6, "secret"),
				),
			},
			{
				Config: testAccSignature_basic_update(name, revSignKey, revSignSecret),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_with_key_update"),
					resource.TestCheckResourceAttr(rName1, "key", revSignKey),
					resource.TestCheckResourceAttr(rName1, "secret", revSignSecret),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_without_key_update"),
					resource.TestCheckResourceAttr(rName2, "key", revSignKey),
					resource.TestCheckResourceAttr(rName2, "secret", revSignSecret),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "name", name+"_hmac_with_key_update"),
					resource.TestCheckResourceAttr(rName3, "key", revSignKey),
					resource.TestCheckResourceAttr(rName3, "secret", revSignSecret),
					rc4.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName4, "name", name+"_hmac_without_key_update"),
					resource.TestCheckResourceAttr(rName4, "key", revSignKey),
					resource.TestCheckResourceAttr(rName4, "secret", revSignSecret),
					rc5.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName5, "name", name+"_aes_with_key_update"),
					resource.TestCheckResourceAttr(rName5, "key", revSignKey),
					resource.TestCheckResourceAttr(rName5, "secret", revSignSecret),
					rc6.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName6, "name", name+"_aes_without_key_update"),
					resource.TestCheckResourceAttr(rName6, "key", revSignKey+signKey),
					resource.TestCheckResourceAttr(rName6, "secret", revSignSecret),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName3),
			},
			{
				ResourceName:      rName4,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName4),
			},
			{
				ResourceName:      rName5,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName5),
			},
			{
				ResourceName:      rName6,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureImportStateFunc(rName6),
			},
		},
	})
}

func testAccSignatureImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
	}
}

func testAccSignature_base(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_apigw_gateway_v2" "gateway" {
  name                            = "%s"
  spec_id                         = "BASIC"
  vpc_id                          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id               = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones              = ["eu-de-01", "eu-de-02"]
  description                     = "test gateway 2"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
  maintain_begin                  = "02:00:00"
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}

func testAccSignature_basic(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_signature_v2" "with_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_with_key"
  type       = "basic"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "without_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_without_key"
  type       = "basic"
}

resource "opentelekomcloud_apigw_signature_v2" "hmac_with_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_hmac_with_key"
  type       = "hmac"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "hmac_without_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_hmac_without_key"
  type       = "hmac"
}

resource "opentelekomcloud_apigw_signature_v2" "aes_with_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_aes_with_key"
  type       = "aes"
  algorithm  = "aes-128-cfb"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "aes_without_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_aes_without_key"
  type       = "aes"
  algorithm  = "aes-256-cfb"
}
`, testAccSignature_base(name), name, signKey, signSecret)
}

func testAccSignature_basic_update(name, signKey, signSecret string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_signature_v2" "with_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_with_key_update"
  type       = "basic"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "without_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_without_key_update"
  type       = "basic"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "hmac_with_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_hmac_with_key_update"
  type       = "hmac"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "hmac_without_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_hmac_without_key_update"
  type       = "hmac"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "aes_with_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_aes_with_key_update"
  type       = "aes"
  algorithm  = "aes-128-cfb"
  key        = "%[3]s"
  secret     = "%[4]s"
}

resource "opentelekomcloud_apigw_signature_v2" "aes_without_key" {
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name       = "%[2]s_aes_without_key_update"
  type       = "aes"
  algorithm  = "aes-256-cfb"
  key        = format("%%s%%s", "%[3]s", strrev("%[3]s")) # the length of the 256 signature key is 32.
  secret     = "%[4]s"
}
`, testAccSignature_base(name), name, signKey, signSecret)
}
