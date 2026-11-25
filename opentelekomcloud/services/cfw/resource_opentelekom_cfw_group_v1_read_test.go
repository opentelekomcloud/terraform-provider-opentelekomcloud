package cfw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testCFWGroupReadConfig(t *testing.T, status int, body string) *cfg.Config {
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

func TestCFWGroupV1ReadRemovesDeletedGroup(t *testing.T) {
	for name, r := range map[string]*schema.Resource{
		"address group": ResourceCfwAddressGroupV1(),
		"service group": ResourceCfwServiceGroupV1(),
	} {
		t.Run(name, func(t *testing.T) {
			d := r.TestResourceData()
			d.SetId("deleted-group")

			diags := r.ReadContext(context.Background(), d, testCFWGroupReadConfig(t, http.StatusBadRequest,
				`{"error_code":"CFW.00200005","error_msg":"Object not found."}`))
			if diags.HasError() {
				t.Fatalf("expected no error for a deleted group, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the group to be removed from state, ID is still %q", d.Id())
			}
		})
	}
}

func TestCFWGroupV1ReadKeepsGroupOnOtherErrors(t *testing.T) {
	for name, r := range map[string]*schema.Resource{
		"address group": ResourceCfwAddressGroupV1(),
		"service group": ResourceCfwServiceGroupV1(),
	} {
		t.Run(name, func(t *testing.T) {
			d := r.TestResourceData()
			d.SetId("group-id")

			diags := r.ReadContext(context.Background(), d, testCFWGroupReadConfig(t, http.StatusBadRequest,
				`{"error_code":"CFW.00200001","error_msg":"The rule direction cannot be empty."}`))
			if !diags.HasError() {
				t.Fatal("expected an error for a 400 with another error code")
			}
			if d.Id() != "group-id" {
				t.Fatalf("expected the group to stay in state, ID is %q", d.Id())
			}
		})
	}
}
