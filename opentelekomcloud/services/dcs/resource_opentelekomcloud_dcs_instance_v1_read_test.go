package dcs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testDcsInstanceV1ReadConfig(t *testing.T, status int, body string) *cfg.Config {
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

func TestDcsInstanceV1ReadRemovesDeletedInstance(t *testing.T) {
	r := ResourceDcsInstanceV1()
	d := r.TestResourceData()
	d.SetId("deleted-instance")

	diags := r.ReadContext(context.Background(), d, testDcsInstanceV1ReadConfig(t, http.StatusNotFound,
		`{"error_code":"DCS.4010","error_msg":"The DCS instance does not exist."}`))
	if diags.HasError() {
		t.Fatalf("expected no error for a deleted instance, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the instance to be removed from state, ID is still %q", d.Id())
	}
}

func TestDcsInstanceV1ReadKeepsInstanceOnOtherErrors(t *testing.T) {
	r := ResourceDcsInstanceV1()
	d := r.TestResourceData()
	d.SetId("instance-id")

	diags := r.ReadContext(context.Background(), d, testDcsInstanceV1ReadConfig(t, http.StatusInternalServerError,
		`{"error_code":"DCS.5000","error_msg":"Internal error"}`))
	if !diags.HasError() {
		t.Fatal("expected an error when the instance cannot be read")
	}
	if d.Id() != "instance-id" {
		t.Fatalf("expected the instance to stay in state, ID is %q", d.Id())
	}
}
