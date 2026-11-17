package vpcep

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testVPCEPEndpointReadConfig(t *testing.T, status int, body string) *cfg.Config {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return &cfg.Config{
		Region: "eu-de",
		HwClient: &golangsdk.ProviderClient{
			HTTPClient: *server.Client(),
			EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
				return server.URL + "/", nil
			},
		},
	}
}

func TestVPCEPEndpointReadRemovesDeletedEndpoint(t *testing.T) {
	r := ResourceVPCEPEndpointV1()
	d := r.TestResourceData()
	d.SetId("deleted-endpoint")

	diags := r.ReadContext(context.Background(), d, testVPCEPEndpointReadConfig(t, http.StatusNotFound,
		`{"error":{"error_code":"EndPoint.2015","message":"The endpoint does not exist."}}`))
	if diags.HasError() {
		t.Fatalf("expected no error for a deleted endpoint, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the endpoint to be removed from state, ID is still %q", d.Id())
	}
}

func TestVPCEPEndpointReadKeepsEndpointOnOtherErrors(t *testing.T) {
	r := ResourceVPCEPEndpointV1()
	d := r.TestResourceData()
	d.SetId("endpoint-id")

	diags := r.ReadContext(context.Background(), d, testVPCEPEndpointReadConfig(t, http.StatusInternalServerError,
		`{"error":{"error_code":"EndPoint.0000","message":"Internal error"}}`))
	if !diags.HasError() {
		t.Fatal("expected an error when the endpoint cannot be read")
	}
	if d.Id() != "endpoint-id" {
		t.Fatalf("expected the endpoint to stay in state, ID is %q", d.Id())
	}
}
