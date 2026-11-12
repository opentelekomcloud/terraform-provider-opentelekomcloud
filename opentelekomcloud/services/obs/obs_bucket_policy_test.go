package obs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const obsBucketPolicyTestBucket = "policy-test-bucket"
const obsBucketPolicyTestJSON = `{"Statement":[{"Effect":"Allow","Principal":{"ID":["*"]},"Action":["ListBucket"],"Resource":["policy-test-bucket"]}]}`

func TestObsBucketPolicyImport(t *testing.T) {
	var requests atomic.Int32
	config := newObsBucketPolicyTestConfig(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		checkObsBucketPolicyTestRequest(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, obsBucketPolicyTestJSON)
	})
	res := ResourceObsBucketPolicy()
	if res.Importer == nil || res.Importer.StateContext == nil {
		t.Fatal("bucket policy must support context-aware import")
	}
	d := res.Data(nil)
	d.SetId(obsBucketPolicyTestBucket)
	imported, err := res.Importer.StateContext(context.Background(), d, config)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(imported) != 1 || imported[0] != d {
		t.Fatalf("unexpected imported state: %#v", imported)
	}
	if requests.Load() != 0 || d.Get("bucket") != "" || d.Get("policy") != "" {
		t.Fatal("passthrough import must start with empty attributes and make no API requests")
	}
	if diags := res.ReadContext(context.Background(), imported[0], config); diags.HasError() {
		t.Fatalf("import refresh failed: %v", diags)
	}
	if requests.Load() != 1 {
		t.Fatalf("import refresh made %d requests, expected one GET", requests.Load())
	}
	if d.Id() != obsBucketPolicyTestBucket || d.Get("bucket") != obsBucketPolicyTestBucket ||
		d.Get("policy") != obsBucketPolicyTestJSON {
		t.Fatalf("import did not restore bucket and raw policy JSON: %s %#v", d.Id(), d.State())
	}
	equivalent := `{
		"Statement": [{"Resource": ["policy-test-bucket"], "Action": ["ListBucket"],
		"Principal": {"ID": ["*"]}, "Effect": "Allow"}]
	}`
	if !res.Schema["policy"].DiffSuppressFunc("policy", obsBucketPolicyTestJSON, equivalent, d) {
		t.Fatal("imported policy must preserve equivalent-JSON diff suppression")
	}
}

func TestObsBucketPolicyReadRemoteChange(t *testing.T) {
	config := newObsBucketPolicyTestConfig(t, func(w http.ResponseWriter, r *http.Request) {
		checkObsBucketPolicyTestRequest(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, obsBucketPolicyTestJSON)
	})
	res := ResourceObsBucketPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"bucket": "stale-state-bucket",
		"policy": `{"Statement":[]}`,
	})
	d.SetId(obsBucketPolicyTestBucket)
	if diags := res.ReadContext(context.Background(), d, config); diags.HasError() {
		t.Fatalf("refresh failed: %v", diags)
	}
	if d.Get("bucket") != obsBucketPolicyTestBucket || d.Get("policy") != obsBucketPolicyTestJSON {
		t.Fatalf("refresh did not reconcile remote bucket/policy: %#v", d.State())
	}
}

func TestObsBucketPolicyReadErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		code    string
		missing bool
	}{
		{"missing_policy", http.StatusNotFound, "NoSuchBucketPolicy", true},
		{"missing_bucket", http.StatusNotFound, "NoSuchBucket", true},
		{"access_denied", http.StatusForbidden, "AccessDenied", false},
		{"unexpected_not_found", http.StatusNotFound, "UnexpectedError", false},
		{"wrong_status_policy", http.StatusBadRequest, "NoSuchBucketPolicy", false},
		{"wrong_status_bucket", http.StatusForbidden, "NoSuchBucket", false},
		{"server_error", http.StatusInternalServerError, "InternalError", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests atomic.Int32
			config := newObsBucketPolicyTestConfig(t, func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				checkObsBucketPolicyTestRequest(t, r, http.MethodGet)
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprintf(w, `<Error><Code>%s</Code><Message>policy read rejected</Message><RequestId>policy-test-request</RequestId></Error>`, tc.code)
			})
			res := ResourceObsBucketPolicy()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"bucket": obsBucketPolicyTestBucket,
				"policy": obsBucketPolicyTestJSON,
			})
			d.SetId(obsBucketPolicyTestBucket)
			diags := res.ReadContext(context.Background(), d, config)
			if requests.Load() == 0 {
				t.Fatal("refresh did not request remote policy")
			}
			if tc.missing {
				if diags.HasError() || d.Id() != "" {
					t.Fatalf("documented missing policy/bucket must clear ID: %s %v", d.Id(), diags)
				}
			} else {
				if !diags.HasError() || d.Id() != obsBucketPolicyTestBucket {
					t.Fatalf("unexpected error must preserve ID and return diagnostics: %s %v", d.Id(), diags)
				}
				text := fmt.Sprint(diags)
				for _, expected := range []string{obsBucketPolicyTestBucket, tc.code, "policy read rejected"} {
					if !strings.Contains(text, expected) {
						t.Errorf("diagnostics missing %q: %s", expected, text)
					}
				}
			}
			if d.Get("bucket") != obsBucketPolicyTestBucket || d.Get("policy") != obsBucketPolicyTestJSON {
				t.Fatal("failed read must not overwrite existing bucket/policy attributes")
			}
		})
	}
}

func TestObsBucketPolicyImportedLifecycle(t *testing.T) {
	var requests atomic.Int32
	config := newObsBucketPolicyTestConfig(t, func(w http.ResponseWriter, r *http.Request) {
		switch requests.Add(1) {
		case 1:
			checkObsBucketPolicyTestRequest(t, r, http.MethodGet)
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, obsBucketPolicyTestJSON)
		case 2:
			checkObsBucketPolicyTestRequest(t, r, http.MethodPut)
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read update request: %v", err)
			}
			if string(body) != obsBucketPolicyTestJSON {
				t.Errorf("unexpected update policy: %s", body)
			}
			w.WriteHeader(http.StatusOK)
		case 3:
			checkObsBucketPolicyTestRequest(t, r, http.MethodDelete)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected API request: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	res := ResourceObsBucketPolicy()
	d := res.Data(nil)
	d.SetId(obsBucketPolicyTestBucket)
	imported, err := res.Importer.StateContext(context.Background(), d, config)
	if err != nil || len(imported) != 1 {
		t.Fatalf("import failed: %v", err)
	}
	for _, operation := range []struct {
		name string
		run  func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
	}{
		{"read", res.ReadContext},
		{"update", res.UpdateContext},
		{"delete", res.DeleteContext},
	} {
		if diags := operation.run(context.Background(), imported[0], config); diags.HasError() {
			t.Fatalf("imported policy %s failed: %v", operation.name, diags)
		}
	}
	if requests.Load() != 3 {
		t.Fatalf("unexpected imported lifecycle request count: %d", requests.Load())
	}
}

func TestObsBucketPolicySchemaCompatibility(t *testing.T) {
	res := ResourceObsBucketPolicy()
	bucket, policy := res.Schema["bucket"], res.Schema["policy"]
	if len(res.Schema) != 2 || bucket == nil || policy == nil {
		t.Fatal("existing bucket and policy attributes must remain unchanged")
	}
	if bucket.Type != schema.TypeString || !bucket.Required || !bucket.ForceNew ||
		policy.Type != schema.TypeString || !policy.Required || policy.ForceNew {
		t.Fatal("existing attribute types and lifecycle flags must remain unchanged")
	}
	if policy.ValidateFunc == nil || policy.DiffSuppressFunc == nil {
		t.Fatal("existing JSON validation and policy diff suppression must remain present")
	}
	if _, errs := policy.ValidateFunc("not-json", "policy"); len(errs) == 0 {
		t.Fatal("policy must still reject invalid JSON")
	}
}

func newObsBucketPolicyTestConfig(t *testing.T, handler http.HandlerFunc) *cfg.Config {
	t.Helper()
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("HTTPS_PROXY", "")
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &cfg.Config{
		Region:    "eu-de",
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		HwClient: &golangsdk.ProviderClient{
			EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
				return server.URL, nil
			},
		},
	}
}

func checkObsBucketPolicyTestRequest(t *testing.T, r *http.Request, method string) {
	t.Helper()
	if r.Method != method || r.URL.Path != "/"+obsBucketPolicyTestBucket || r.URL.RawQuery != "policy" {
		t.Errorf("unexpected request: %s %s, expected %s /%s?policy", r.Method, r.URL, method, obsBucketPolicyTestBucket)
	}
}
