package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const userV2TestPath = "/v2/project-id/instances/instance-id/users"

func userV2TestConfig(t *testing.T, handler http.HandlerFunc) *cfg.Config {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &cfg.Config{
		Region: "eu-de",
		HwClient: &golangsdk.ProviderClient{
			EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
				return server.URL + "/v2/project-id/", nil
			},
		},
	}
}

func TestDmsUserV2SchemaCompatibility(t *testing.T) {
	r := ResourceDmsUsersV2()
	for key, typ := range map[string]schema.ValueType{
		"instance_id": schema.TypeString, "username": schema.TypeString, "password": schema.TypeString,
		"role": schema.TypeString, "default_app": schema.TypeBool, "creation_time": schema.TypeFloat,
	} {
		if r.Schema[key] == nil || r.Schema[key].Type != typ {
			t.Errorf("existing attribute %s removed or changed type", key)
		}
	}
	for _, key := range []string{"instance_id", "username"} {
		if !r.Schema[key].Required || !r.Schema[key].ForceNew {
			t.Errorf("existing %s configuration semantics changed", key)
		}
	}
	if !r.Schema["password"].Required || !r.Schema["password"].Sensitive || r.Schema["password"].ForceNew {
		t.Fatal("password must remain required, sensitive and updatable")
	}
	for _, key := range []string{"role", "default_app", "creation_time"} {
		if !r.Schema[key].Computed {
			t.Errorf("existing %s output must remain computed", key)
		}
	}
	description := r.Schema["description"]
	if description.Type != schema.TypeString || !description.Optional || !description.Computed || !description.ForceNew {
		t.Fatal("description must be a creation-only optional/computed string")
	}
	if r.Importer == nil || r.Importer.StateContext == nil {
		t.Fatal("existing importer removed")
	}
}

func TestDmsUserV2Lifecycle(t *testing.T) {
	for _, description := range []string{"Kafka application user", ""} {
		t.Run(fmt.Sprintf("description=%q", description), func(t *testing.T) {
			calls := make(chan string, 8)
			config := userV2TestConfig(t, func(w http.ResponseWriter, r *http.Request) {
				calls <- r.Method + " " + r.URL.Path
				if r.Method == http.MethodGet && r.URL.Path == userV2TestPath {
					w.Header().Set("Content-Type", "application/json")
					_, err := fmt.Fprintf(w, `{"users":[
						{"user_name":"other-user","user_desc":"not the managed user"},
						{"user_name":"test-user","user_desc":%q,"role":"guest","default_app":false,"created_time":1615431764734}
					]}`, description)
					if err != nil {
						t.Errorf("write user response: %s", err)
					}
					return
				}
				var expected map[string]interface{}
				switch {
				case r.Method == http.MethodPost && r.URL.Path == userV2TestPath:
					expected = map[string]interface{}{"user_name": "test-user", "user_passwd": "Create1!password"}
					if description != "" {
						expected["user_desc"] = description
					}
				case r.Method == http.MethodPut && r.URL.Path == userV2TestPath+"/test-user":
					expected = map[string]interface{}{"new_password": "Changed2!password"}
				case r.Method == http.MethodPut && r.URL.Path == userV2TestPath:
					expected = map[string]interface{}{"action": "delete", "users": []interface{}{"test-user"}}
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.Error(w, "unexpected request", http.StatusBadRequest)
					return
				}
				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode user request: %s", err)
					http.Error(w, "invalid body", http.StatusBadRequest)
					return
				}
				if !reflect.DeepEqual(body, expected) {
					t.Errorf("request body = %#v, want %#v", body, expected)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			values := map[string]interface{}{
				"instance_id": "instance-id", "username": "test-user", "password": "Create1!password",
			}
			if description != "" {
				values["description"] = description
			}
			d := schema.TestResourceDataRaw(t, ResourceDmsUsersV2().Schema, values)
			if diagnostics := resourceDmsUsersV2Create(context.Background(), d, config); diagnostics.HasError() {
				t.Fatalf("create failed: %v", diagnostics)
			}
			if d.Id() != "test-user" {
				t.Fatalf("ID = %q, want managed username", d.Id())
			}
			for key, expected := range map[string]interface{}{
				"username": "test-user", "description": description, "password": "Create1!password",
				"role": "guest", "default_app": false, "creation_time": float64(1615431764734),
			} {
				if actual := d.Get(key); actual != expected {
					t.Errorf("%s = %#v, want %#v", key, actual, expected)
				}
			}
			values["password"] = "Changed2!password"
			updated := schema.TestResourceDataRaw(t, ResourceDmsUsersV2().Schema, values)
			updated.SetId(d.Id())
			if !updated.HasChange("password") {
				t.Fatal("password update fixture must exercise the reset request")
			}
			if diagnostics := resourceDmsUsersV2Update(context.Background(), updated, config); diagnostics.HasError() {
				t.Fatalf("password update failed: %v", diagnostics)
			}
			if updated.Get("password") != values["password"] || updated.Get("description") != description {
				t.Fatal("read after password reset changed password or description")
			}
			if diagnostics := resourceDmsUsersV2Delete(context.Background(), updated, config); diagnostics.HasError() {
				t.Fatalf("delete failed: %v", diagnostics)
			}
			if updated.Id() != "" {
				t.Fatal("successful deletion must clear the resource ID")
			}
			for _, expected := range []string{
				"POST " + userV2TestPath, "GET " + userV2TestPath,
				"PUT " + userV2TestPath + "/test-user", "GET " + userV2TestPath, "PUT " + userV2TestPath,
			} {
				select {
				case actual := <-calls:
					if actual != expected {
						t.Errorf("request = %q, want %q", actual, expected)
					}
				default:
					t.Errorf("missing request %q", expected)
				}
			}
			if len(calls) != 0 {
				t.Fatal("unexpected additional user requests")
			}
		})
	}
}

func TestDmsUserV2ReadClearsStaleFields(t *testing.T) {
	config := userV2TestConfig(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprint(w, `{"users":[{"user_name":"test-user","role":"guest","default_app":false,"created_time":0}]}`)
		if err != nil {
			t.Errorf("write user response: %s", err)
		}
	})
	d := schema.TestResourceDataRaw(t, ResourceDmsUsersV2().Schema, map[string]interface{}{
		"instance_id": "instance-id", "username": "test-user", "password": "Create1!password",
		"description": "stale", "default_app": true, "creation_time": float64(1615431764734),
	})
	d.SetId("test-user")
	if diagnostics := resourceDmsUsersV2Read(context.Background(), d, config); diagnostics.HasError() {
		t.Fatalf("read failed: %v", diagnostics)
	}
	if d.Get("description") != "" || d.Get("default_app") != false || d.Get("creation_time") != float64(0) {
		t.Fatal("missing description and zero-value fields must clear stale state")
	}
	if d.Get("password") != "Create1!password" {
		t.Fatal("list must preserve the request-only password")
	}
}

func TestDmsUserV2Errors(t *testing.T) {
	for _, tc := range []struct {
		name        string
		operation   schema.CreateContextFunc
		status      int
		body        string
		expectError bool
		expectID    string
	}{
		{"create", resourceDmsUsersV2Create, http.StatusForbidden, `{"error_msg":"test failure"}`, true, ""},
		{"read", resourceDmsUsersV2Read, http.StatusInternalServerError, `{"error_msg":"test failure"}`, true, "test-user"},
		{"reset_password", resourceDmsUsersV2Update, http.StatusForbidden, `{"error_msg":"test failure"}`, true, "test-user"},
		{"delete", resourceDmsUsersV2Delete, http.StatusInternalServerError, `{"error_msg":"test failure"}`, true, "test-user"},
		{"missing_user", resourceDmsUsersV2Read, http.StatusOK, `{"users":[]}`, true, "test-user"},
		{"missing_instance", resourceDmsUsersV2Read, http.StatusNotFound, `{}`, false, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := userV2TestConfig(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				if _, err := fmt.Fprint(w, tc.body); err != nil {
					t.Errorf("write error response: %s", err)
				}
			})
			d := schema.TestResourceDataRaw(t, ResourceDmsUsersV2().Schema, map[string]interface{}{
				"instance_id": "instance-id", "username": "test-user", "password": "Changed2!password",
			})
			if tc.name != "create" {
				d.SetId("test-user")
			}
			diagnostics := tc.operation(context.Background(), d, config)
			if diagnostics.HasError() != tc.expectError || d.Id() != tc.expectID {
				t.Fatalf("diagnostics = %v, ID = %q", diagnostics, d.Id())
			}
			for _, diagnostic := range diagnostics {
				if diagnostic.Severity == diag.Error && strings.Contains(diagnostic.Summary, "Changed2!password") {
					t.Fatal("password leaked in diagnostics")
				}
			}
		})
	}
}
