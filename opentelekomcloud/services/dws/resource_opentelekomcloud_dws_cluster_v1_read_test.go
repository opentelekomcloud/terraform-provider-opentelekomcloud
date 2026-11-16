package dws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testDwsClusterV1ReadConfig(t *testing.T, status int, body string) *cfg.Config {
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

func TestDwsClusterV1ReadRemovesDeletedCluster(t *testing.T) {
	r := ResourceDcsInstanceV1()
	d := r.TestResourceData()
	d.SetId("deleted-cluster")

	diags := r.ReadContext(context.Background(), d, testDwsClusterV1ReadConfig(t, http.StatusNotFound,
		`{"error_code":"DWS.0001","error_msg":"The cluster does not exist."}`))
	if diags.HasError() {
		t.Fatalf("expected no error for a deleted cluster, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the cluster to be removed from state, ID is still %q", d.Id())
	}
}

func TestDwsClusterV1ReadKeepsClusterOnOtherErrors(t *testing.T) {
	r := ResourceDcsInstanceV1()
	d := r.TestResourceData()
	d.SetId("cluster-id")

	diags := r.ReadContext(context.Background(), d, testDwsClusterV1ReadConfig(t, http.StatusInternalServerError,
		`{"error_code":"DWS.0000","error_msg":"Internal error"}`))
	if !diags.HasError() {
		t.Fatal("expected an error when the cluster cannot be read")
	}
	if d.Id() != "cluster-id" {
		t.Fatalf("expected the cluster to stay in state, ID is %q", d.Id())
	}
}
