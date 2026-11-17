package cts

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

// testCTSNotificationReadContext returns a context carrying a CTS v3 client that answers every
// request with body. The client goes through the ctx because CtsV3Client re-authenticates.
func testCTSNotificationReadContext(t *testing.T, body string) (context.Context, *cfg.Config) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	client, err := openstack.NewCTSV3(&golangsdk.ProviderClient{
		HTTPClient: *server.Client(),
		EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
			return server.URL + "/v3/project/", nil
		},
	}, golangsdk.EndpointOpts{Region: "eu-de"})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), keyClientV3, client) //nolint:staticcheck // ClientFromCtx reads a string key

	return ctx, &cfg.Config{Region: "eu-de"}
}

func TestCTSEventNotificationReadRemovesMissingNotification(t *testing.T) {
	cases := map[string]string{
		"empty list":     `{"notifications":[]}`,
		"only other one": `{"notifications":[{"notification_id":"other-id","notification_name":"keyevent-old","notification_type":"smn","status":"enabled"}]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			r := ResourceCTSEventNotificationV3()
			d := r.TestResourceData()
			d.SetId("notification-id/keyevent")

			ctx, config := testCTSNotificationReadContext(t, body)
			diags := r.ReadContext(ctx, d, config)
			if diags.HasError() {
				t.Fatalf("expected no error for a deleted notification, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the notification to be removed from state, ID is still %q", d.Id())
			}
		})
	}
}

func TestCTSEventNotificationReadPicksTheNamedNotification(t *testing.T) {
	r := ResourceCTSEventNotificationV3()
	d := r.TestResourceData()
	d.SetId("notification-id/keyevent")

	ctx, config := testCTSNotificationReadContext(t, `{"notifications":[
		{"notification_id":"other-id","notification_name":"keyevent-old","notification_type":"smn","status":"disabled"},
		{"notification_id":"notification-id","notification_name":"keyevent","notification_type":"smn","status":"enabled","operation_type":"complete"}
	]}`)
	diags := r.ReadContext(ctx, d, config)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if got := d.Get("notification_id").(string); got != "notification-id" {
		t.Fatalf("expected notification_id of the named notification, got %q", got)
	}
	if got := d.Get("status").(string); got != "enabled" {
		t.Fatalf("expected status of the named notification, got %q", got)
	}
}
