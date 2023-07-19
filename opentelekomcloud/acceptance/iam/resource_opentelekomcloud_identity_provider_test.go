package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const providerNewResource = "opentelekomcloud_identity_provider.provider"

func TestAccIdentityProviderBasic(t *testing.T) {
	var nameSaml = fmt.Sprintf("test-provider-saml-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderBasic(nameSaml),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(providerNewResource, "status", "true"),
					resource.TestCheckResourceAttr(providerNewResource, "description", "test-provider"),
				),
			},
			{
				Config: testAccIdentityProviderUpdate(nameSaml),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(providerNewResource, "status", "false"),
					resource.TestCheckResourceAttr(providerNewResource, "description", "test-provider-2"),
					resource.TestCheckResourceAttrSet(providerNewResource, "metadata"),
				),
			},
		},
	})
}

func TestAccIdentityProviderOIDC(t *testing.T) {
	var nameOIDC = fmt.Sprintf("test-provider-oidc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderOIDCBasic(nameOIDC),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(providerNewResource, "status", "true"),
					resource.TestCheckResourceAttr(providerNewResource, "description", "test-provider"),
					resource.TestCheckResourceAttr(providerNewResource, "access_config.0.access_type", "program_console"),
					resource.TestCheckResourceAttr(providerNewResource, "access_config.0.provider_url", "https://accounts.example.com"),
					resource.TestCheckResourceAttr(providerNewResource, "access_config.0.client_id", "your_client_id"),
				),
			},
			{
				Config: testAccIdentityProviderOIDCUpdate(nameOIDC),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(providerNewResource, "status", "false"),
					resource.TestCheckResourceAttr(providerNewResource, "description", "test-provider-updated"),
					resource.TestCheckResourceAttr(providerNewResource, "access_config.0.access_type", "program"),
					resource.TestCheckResourceAttr(providerNewResource, "access_config.0.provider_url", "https://accounts.example.com"),
					resource.TestCheckResourceAttr(providerNewResource, "access_config.0.client_id", "your_client_id_2"),
				),
			},
		},
	})
}

func TestAccIdentityOIDCProvider_import(t *testing.T) {
	var nameOIDC = fmt.Sprintf("test-provider-oidc-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderOIDCBasic(nameOIDC),
			},
			{
				ResourceName:      providerNewResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIdentitySamlProvider_import(t *testing.T) {
	var nameSaml = fmt.Sprintf("test-provider-oidc-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderUpdate(nameSaml),
			},
			{
				ResourceName:      providerNewResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIdentityProviderDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity v3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != providerResource {
			continue
		}

		_, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("provider still exists")
		}
	}

	return nil
}

func testAccIdentityProviderBasic(providerName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_provider" "provider" {
  name        = "%s"
  protocol    = "saml"
  status      = true
  description = "test-provider"
}
  `, providerName)
}

func testAccIdentityProviderUpdate(providerName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_provider" "provider" {
  name        = "%s"
  protocol    = "saml"
  status      = false
  description = "test-provider-2"
  metadata    = %s
}
  `, providerName, Metadata)
}

func testAccIdentityProviderOIDCBasic(providerName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_provider" "provider" {
  name        = "%s"
  protocol    = "oidc"
  status      = true
  description = "test-provider"
  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "your_client_id"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
    response_type          = "id_token"
    response_mode          = "fragment"
    signing_key = jsonencode(
      {
        keys = [
          {
            alg = "RS256"
            e   = "AQAB"
            kid = "..."
            kty = "RSA"
            n   = "..."
            use = "sig"
          },
        ]
      }
    )
  }
}
  `, providerName)
}

func testAccIdentityProviderOIDCUpdate(providerName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_provider" "provider" {
  name        = "%s"
  protocol    = "oidc"
  status      = false
  description = "test-provider-updated"
  access_config {
    access_type  = "program"
    provider_url = "https://accounts.example.com"
    client_id    = "your_client_id_2"
    signing_key = jsonencode(
      {
        keys = [
          {
            kid : "d05ef20c4512645vv1...",
            n : "cws_cnjiwsbvweolwn_-vnl...",
            e : "AQAB",
            kty : "RSA",
            use : "sig",
            alg : "RS256"
          },
        ]
      }
    )
  }
}
  `, providerName)
}
