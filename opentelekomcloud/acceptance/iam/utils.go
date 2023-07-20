package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccIdentityV3AgencyPreCheck(t *testing.T) {
	if env.OS_TENANT_NAME == "" {
		t.Skip("OS_TENANT_NAME must be set for acceptance tests")
	}
}

const Metadata = `<<EOT
<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2023-04-28T16:06:53Z"
                     cacheDuration="PT604800S"
                     entityID="https://idp.hfbk-dresden.de/idp/shibboleth">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="https://idp.hfbk-dresden.de/idp/profile/SAML2/POST/SLO"
                                     index="1" />

    </md:SPSSODescriptor>
</md:EntityDescriptor>
EOT`
