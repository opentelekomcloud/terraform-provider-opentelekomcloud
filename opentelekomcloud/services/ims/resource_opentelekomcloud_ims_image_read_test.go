package ims

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testImsReadConfig(t *testing.T, status int, body string) *cfg.Config {
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

func TestImsImageReadRemovesDeletedImage(t *testing.T) {
	resources := map[string]func() *schema.Resource{
		"image":      ResourceImsImageV2,
		"data image": ResourceImsDataImageV2,
	}

	for name, resource := range resources {
		t.Run(name+" not listed", func(t *testing.T) {
			r := resource()
			d := r.TestResourceData()
			d.SetId("deleted-image")

			diags := r.ReadContext(context.Background(), d, testImsReadConfig(t, http.StatusOK, `{"images":[]}`))
			if diags.HasError() {
				t.Fatalf("expected no error for an image that is no longer listed, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the image to be removed from state, ID is still %q", d.Id())
			}
		})

		t.Run(name+" list error", func(t *testing.T) {
			r := resource()
			d := r.TestResourceData()
			d.SetId("image-id")

			diags := r.ReadContext(context.Background(), d, testImsReadConfig(t, http.StatusForbidden, `{"message":"Forbidden"}`))
			if !diags.HasError() {
				t.Fatal("expected an error when the image list fails")
			}
			if d.Id() != "image-id" {
				t.Fatalf("expected the image to stay in state, ID is %q", d.Id())
			}
		})
	}
}
