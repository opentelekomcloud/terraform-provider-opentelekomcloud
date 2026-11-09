package er

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestErInstanceReadRemovesMissingResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error_code":"ER.00200001","error_msg":"instance not found"}`))
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

	r := ResourceErInstanceV3()
	d := r.TestResourceData()
	d.SetId("deleted-out-of-band")

	diags := r.ReadContext(context.Background(), d, config)
	if diags.HasError() {
		t.Fatalf("expected no error for a 404, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the instance to be removed from state, ID is still %q", d.Id())
	}
}
