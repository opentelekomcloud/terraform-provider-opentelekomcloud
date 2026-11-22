package iam

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func testIdentityReadClient(t *testing.T, handler http.HandlerFunc) *golangsdk.ServiceClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &golangsdk.ServiceClient{
		ProviderClient: &golangsdk.ProviderClient{HTTPClient: *server.Client()},
		Endpoint:       server.URL + "/v3/",
	}
}

func testIdentityReadReply(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

// testIdentityV3Ctx stores the client the way common.ClientFromCtx reads it back, by its plain string key.
func testIdentityV3Ctx(client *golangsdk.ServiceClient) context.Context {
	return context.WithValue(context.Background(), keyClientV3, client) //nolint:staticcheck
}

const testIdentityNotFound = `{"error":{"code":404,"message":"Could not find resource.","title":"Not Found"}}`

func TestIdentityProtocolReadRemovesDeletedProtocol(t *testing.T) {
	client := testIdentityReadClient(t, testIdentityReadReply(http.StatusNotFound, testIdentityNotFound))
	ctx := ctxWithClient(context.Background(), client)

	r := ResourceIdentityProtocolV3()
	d := r.TestResourceData()
	d.SetId("idp/saml")

	diags := r.ReadContext(ctx, d, &cfg.Config{})
	if diags.HasError() {
		t.Fatalf("expected no error for a deleted protocol, got: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected the protocol to be removed from state, ID is still %q", d.Id())
	}
}

func TestIdentityProtocolReadKeepsProtocolOnOtherErrors(t *testing.T) {
	client := testIdentityReadClient(t, testIdentityReadReply(http.StatusInternalServerError, `{"error":{"code":500}}`))
	ctx := ctxWithClient(context.Background(), client)

	r := ResourceIdentityProtocolV3()
	d := r.TestResourceData()
	d.SetId("idp/saml")

	diags := r.ReadContext(ctx, d, &cfg.Config{})
	if !diags.HasError() {
		t.Fatal("expected an error when the protocol cannot be read")
	}
	if d.Id() != "idp/saml" {
		t.Fatalf("expected the protocol to stay in state, ID is %q", d.Id())
	}
}

func TestIdentityRoleAssignmentReadRemovesDeletedAssignment(t *testing.T) {
	cases := []struct {
		name        string
		allProjects bool
		handler     http.HandlerFunc
	}{
		{
			name:    "role no longer assigned",
			handler: testIdentityReadReply(http.StatusOK, `{"roles":[{"id":"other-role","name":"other"}],"links":{}}`),
		},
		{
			name:    "group deleted",
			handler: testIdentityReadReply(http.StatusNotFound, testIdentityNotFound),
		},
		{
			name:        "all projects assignment removed",
			allProjects: true,
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodHead {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				w.WriteHeader(http.StatusNotFound)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := testIdentityReadClient(t, tc.handler)
			ctx := testIdentityV3Ctx(client)

			r := ResourceIdentityRoleAssignmentV3()
			d := r.TestResourceData()
			id := "/project-id/group-id/role-id"
			if tc.allProjects {
				id = "domain-id//group-id/role-id"
				if err := d.Set("all_projects", true); err != nil {
					t.Fatal(err)
				}
			}
			d.SetId(id)

			diags := r.ReadContext(ctx, d, &cfg.Config{})
			if diags.HasError() {
				t.Fatalf("expected no error for a deleted role assignment, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the role assignment to be removed from state, ID is still %q", d.Id())
			}
		})
	}
}

func TestIdentityRoleAssignmentReadKeepsAssignedRole(t *testing.T) {
	client := testIdentityReadClient(t, testIdentityReadReply(http.StatusOK,
		`{"roles":[{"id":"other-role","name":"other"},{"id":"role-id","name":"te_admin"}],"links":{}}`))
	ctx := testIdentityV3Ctx(client)

	r := ResourceIdentityRoleAssignmentV3()
	d := r.TestResourceData()
	d.SetId("/project-id/group-id/role-id")

	diags := r.ReadContext(ctx, d, &cfg.Config{})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if d.Id() == "" || d.Get("role_id").(string) != "role-id" {
		t.Fatalf("expected the assignment to stay with role_id %q, got ID %q role_id %q", "role-id", d.Id(), d.Get("role_id"))
	}
}
