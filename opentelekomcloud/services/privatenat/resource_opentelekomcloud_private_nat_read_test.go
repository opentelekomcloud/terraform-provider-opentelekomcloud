package privatenat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestPrivateNatReadRemovesMissingResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error_code":"NAT.0001","error_msg":"resource not found"}`))
	}))
	defer server.Close()

	config := &cfg.Config{
		Region: "eu-de",
		HwClient: &golangsdk.ProviderClient{
			HTTPClient: *server.Client(),
			EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
				return server.URL + "/", nil
			},
		},
	}

	resources := map[string]*schema.Resource{
		"gateway":    ResourcePrivateNatGatewayV3(),
		"snat_rule":  ResourcePrivateNatSnatRuleV3(),
		"dnat_rule":  ResourcePrivateNatDnatRuleV3(),
		"transit_ip": ResourcePrivateNatTransitIpV3(),
	}
	for name, r := range resources {
		t.Run(name, func(t *testing.T) {
			d := r.TestResourceData()
			d.SetId("deleted-out-of-band")

			diags := r.ReadContext(context.Background(), d, config)
			if diags.HasError() {
				t.Fatalf("expected no error for a 404, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the resource to be removed from state, ID is still %q", d.Id())
			}
		})
	}
}
