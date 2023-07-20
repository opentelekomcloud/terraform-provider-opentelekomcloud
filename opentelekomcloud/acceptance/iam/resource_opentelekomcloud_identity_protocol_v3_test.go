package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const protocolResourceName = "opentelekomcloud_identity_protocol_v3.saml"

func TestAccIdentityV3Protocol_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProtocolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProtocolBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(protocolResourceName, "mapping_id", mapping),
					resource.TestCheckResourceAttr(protocolResourceName, "provider_id", providerName),
					resource.TestCheckResourceAttr(protocolResourceName, "links.%", "2"),
					resource.TestCheckResourceAttrSet(protocolResourceName, "links.self"),
				),
			},
		},
	})
}

func TestAccIdentityV3Protocol_OIDC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProtocolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProtocolOIDC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(protocolResourceName, "mapping_id", mapping),
					resource.TestCheckResourceAttr(protocolResourceName, "provider_id", providerName),
					resource.TestCheckResourceAttr(protocolResourceName, "links.%", "2"),
					resource.TestCheckResourceAttrSet(protocolResourceName, "links.self"),
					resource.TestCheckResourceAttr(protocolResourceName, "access_config.0.access_type", "program_console"),
					resource.TestCheckResourceAttr(protocolResourceName, "access_config.0.provider_url", "https://accounts.example.com"),
					resource.TestCheckResourceAttr(protocolResourceName, "access_config.0.client_id", "your_client_id"),
				),
			},
		},
	})
}

func TestAccIdentityV3Protocol_metadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProtocolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProtocolMetadata,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(protocolResourceName, "metadata.0.metadata"),
					resource.TestCheckResourceAttrSet(protocolResourceName, "metadata.0.domain_id"),
				),
			},
		},
	})
}

func TestAccIdentityV3Protocol_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProtocolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProtocolBasic,
			},
			{
				ResourceName:      protocolResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIdentityV3ProtocolDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	identityClient, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity v3 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_protocol_v3" {
			continue
		}

		_, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("provider still exists")
		}
	}

	return nil
}

var domainId = os.Getenv("OS_TEST_DOMAINID")

var (
	mapping                        = tools.RandomString("mapping-", 3)
	protocolName                   = "saml"
	protocolOIDC                   = "oidc"
	testAccIdentityV3ProtocolBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_identity_protocol_v3" "saml" {
  protocol    = "%s"
  provider_id = opentelekomcloud_identity_provider_v3.provider.id
  mapping_id  = opentelekomcloud_identity_mapping_v3.mapping.id
}
`, testAccIdentityV3ProviderBasic, testAccIdentityV3MappingBasic(mapping), protocolName)
	testAccIdentityV3ProtocolMetadata = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_identity_protocol_v3" "saml" {
  protocol    = "%s"
  provider_id = opentelekomcloud_identity_provider_v3.provider.id
  mapping_id  = opentelekomcloud_identity_mapping_v3.mapping.id

  metadata {
    metadata  = %s
    domain_id = "%s"
  }
}
`, testAccIdentityV3ProviderBasic, testAccIdentityV3MappingBasic(mapping), protocolName, Metadata, domainId)

	testAccIdentityV3ProtocolOIDC = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_identity_protocol_v3" "saml" {
  protocol    = "%s"
  provider_id = opentelekomcloud_identity_provider_v3.provider.id
  mapping_id  = opentelekomcloud_identity_mapping_v3.mapping.id
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
`, testAccIdentityV3ProviderBasic, testAccIdentityV3MappingBasic(mapping), protocolOIDC)
)
