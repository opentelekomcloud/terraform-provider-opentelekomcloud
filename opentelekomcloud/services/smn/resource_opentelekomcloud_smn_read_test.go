package smn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const testTopicUrn = "urn:smn:eu-de:project:topic"

// testSmnReadContext returns a context carrying an SMN v2 client that talks to handler. The
// client is injected the same way ClientFromCtx looks it up, because SmnV2Client re-authenticates.
func testSmnReadContext(t *testing.T, handler http.HandlerFunc) (context.Context, *cfg.Config) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := openstack.NewSMNV2(&golangsdk.ProviderClient{
		HTTPClient: *server.Client(),
		EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
			return server.URL + "/v2/project/", nil
		},
	}, golangsdk.EndpointOpts{Region: "eu-de"})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), keyClientV2, client) //nolint:staticcheck // ClientFromCtx reads a string key

	return ctx, &cfg.Config{Region: "eu-de"}
}

// topicSubscriptions serves total subscriptions of testTopicUrn, at most 100 per page, on both the
// topic and the project-wide subscription lists.
func topicSubscriptions(total int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allSubscriptions := strings.HasSuffix(r.URL.Path, "/notifications/subscriptions")
		if !allSubscriptions && !strings.Contains(r.URL.Path, "/topics/"+testTopicUrn+"/subscriptions") {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error_code":"SMN.00010001","error_msg":"topic not found"}`))
			return
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit > 100 {
			limit = 100
		}
		items := []map[string]any{}
		for i := offset; i < total && i < offset+limit; i++ {
			items = append(items, map[string]any{
				"subscription_urn": fmt.Sprintf("%s:sub-%d", testTopicUrn, i),
				"topic_urn":        testTopicUrn,
				"protocol":         "email",
				"endpoint":         fmt.Sprintf("user%d@example.com", i),
				"status":           1,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"subscription_count": total, "subscriptions": items})
	}
}

func TestSmnSubscriptionRead(t *testing.T) {
	cases := []struct {
		name      string
		topicUrn  string
		id        string
		total     int
		wantID    string
		wantEmail string
	}{
		{
			name:      "subscription on a later page",
			topicUrn:  testTopicUrn,
			id:        testTopicUrn + ":sub-150",
			total:     180,
			wantID:    testTopicUrn + ":sub-150",
			wantEmail: "user150@example.com",
		},
		{
			name:     "subscription deleted outside terraform",
			topicUrn: testTopicUrn,
			id:       testTopicUrn + ":sub-500",
			total:    3,
			wantID:   "",
		},
		{
			name:     "topic deleted outside terraform",
			topicUrn: "urn:smn:eu-de:project:gone",
			id:       "urn:smn:eu-de:project:gone:sub-0",
			wantID:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := ResourceSubscription()
			d := r.TestResourceData()
			d.SetId(tc.id)
			_ = d.Set("topic_urn", tc.topicUrn)

			ctx, config := testSmnReadContext(t, topicSubscriptions(tc.total))
			diags := r.ReadContext(ctx, d, config)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if d.Id() != tc.wantID {
				t.Fatalf("ID = %q, want %q", d.Id(), tc.wantID)
			}
			if tc.wantEmail != "" && d.Get("endpoint").(string) != tc.wantEmail {
				t.Fatalf("endpoint = %q, want %q", d.Get("endpoint"), tc.wantEmail)
			}
		})
	}
}

func TestSmnSubscriptionReadKeepsStateOnError(t *testing.T) {
	r := ResourceSubscription()
	d := r.TestResourceData()
	d.SetId(testTopicUrn + ":sub-0")
	_ = d.Set("topic_urn", testTopicUrn)

	ctx, config := testSmnReadContext(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error_code":"SMN.00010003","error_msg":"forbidden"}`))
	})
	diags := r.ReadContext(ctx, d, config)
	if !diags.HasError() {
		t.Fatal("expected an error when the subscription list fails")
	}
	if d.Id() == "" {
		t.Fatal("expected the subscription to stay in state")
	}
}

func TestSmnMessageTemplateReadRemovesDeletedTemplate(t *testing.T) {
	r := ResourceSmnMessageTemplateV2()
	d := r.TestResourceData()
	d.SetId("deleted-template")

	ctx, config := testSmnReadContext(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error_code":"SMN.00011020","error_msg":"message template not found"}`))
	})
	diags := r.ReadContext(ctx, d, config)
	if diags.HasError() {
		t.Fatalf("expected no error for a deleted message template, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the template to be removed from state, ID is still %q", d.Id())
	}
}
