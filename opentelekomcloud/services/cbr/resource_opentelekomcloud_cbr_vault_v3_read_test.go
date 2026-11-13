package cbr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testCBRVaultReadConfig(t *testing.T, status int, body string) *cfg.Config {
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

func TestCBRVaultReadRemovesMissingVault(t *testing.T) {
	config := testCBRVaultReadConfig(t, http.StatusNotFound, `{"message":"Not Found"}`)

	r := ResourceCBRVaultV3()
	d := r.TestResourceData()
	d.SetId("deleted-out-of-band")

	diags := r.ReadContext(context.Background(), d, config)
	if diags.HasError() {
		t.Fatalf("expected no error for a 404, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the vault to be removed from state, ID is still %q", d.Id())
	}
}

func TestCBRVaultReadKeepsOtherErrors(t *testing.T) {
	config := testCBRVaultReadConfig(t, http.StatusForbidden, `{"message":"Forbidden"}`)

	r := ResourceCBRVaultV3()
	d := r.TestResourceData()
	d.SetId("vault-id")

	diags := r.ReadContext(context.Background(), d, config)
	if !diags.HasError() {
		t.Fatal("expected an error for a 403")
	}
	if d.Id() != "vault-id" {
		t.Fatalf("expected the vault to stay in state, ID is %q", d.Id())
	}
}
