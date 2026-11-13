package ddm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testDdmReadConfig(t *testing.T, status int) *cfg.Config {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"message":"` + http.StatusText(status) + `"}`))
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

func TestDdmReadRemovesMissingResources(t *testing.T) {
	cases := []struct {
		name     string
		resource func() *schema.Resource
		setup    func(d *schema.ResourceData)
	}{
		{
			name:     "instance",
			resource: ResourceDdmInstanceV1,
			setup:    func(d *schema.ResourceData) { d.SetId("deleted-instance") },
		},
		{
			name:     "schema",
			resource: ResourceDdmSchemaV1,
			setup: func(d *schema.ResourceData) {
				d.SetId("deleted-instance/db1")
				_ = d.Set("instance_id", "deleted-instance")
				_ = d.Set("name", "db1")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+" 404", func(t *testing.T) {
			r := tc.resource()
			d := r.TestResourceData()
			tc.setup(d)

			diags := r.ReadContext(context.Background(), d, testDdmReadConfig(t, http.StatusNotFound))
			if diags.HasError() {
				t.Fatalf("expected no error for a 404, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the resource to be removed from state, ID is still %q", d.Id())
			}
		})

		t.Run(tc.name+" 403", func(t *testing.T) {
			r := tc.resource()
			d := r.TestResourceData()
			tc.setup(d)
			id := d.Id()

			diags := r.ReadContext(context.Background(), d, testDdmReadConfig(t, http.StatusForbidden))
			if !diags.HasError() {
				t.Fatal("expected an error for a 403")
			}
			if d.Id() != id {
				t.Fatalf("expected the resource to stay in state, ID is %q", d.Id())
			}
		})
	}
}
