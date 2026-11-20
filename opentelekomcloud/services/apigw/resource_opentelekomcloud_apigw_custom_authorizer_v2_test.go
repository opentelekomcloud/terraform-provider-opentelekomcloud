package apigw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/authorizer"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const testFunctionUrn = "urn:fss:eu-de:0123456789abcdef0123456789abcdef:function:default:apigw-auth_func.v1"

func TestParseFunctionUrn(t *testing.T) {
	cases := []struct {
		name            string
		input           string
		expectedUrn     string
		expectedVersion string
	}{
		{"without version", testFunctionUrn, testFunctionUrn, ""},
		{"with latest version", testFunctionUrn + ":latest", testFunctionUrn, "latest"},
		{"with custom version", testFunctionUrn + ":v1.2-beta", testFunctionUrn, "v1.2-beta"},
		{"with alias", testFunctionUrn + ":!prod", testFunctionUrn, "!prod"},
		{"empty", "", "", ""},
		{"not a function urn", "foo:bar", "foo:bar", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			urn, version := parseFunctionUrn(c.input)
			if urn != c.expectedUrn || version != c.expectedVersion {
				t.Fatalf("expected (%q, %q), got (%q, %q)", c.expectedUrn, c.expectedVersion, urn, version)
			}
		})
	}
}

func TestResolveFunctionVersion(t *testing.T) {
	cases := []struct {
		name              string
		configuredVersion string
		isConfigured      bool
		urnVersion        string
		currentVersion    string
		expected          string
	}{
		{"configured version has priority", "v1", true, "latest", "v2", "v1"},
		{"version suffix of the URN", "", false, "latest", "v2", "latest"},
		{"current version as fallback", "", false, "", "latest", "latest"},
		{"empty configured version is ignored", "", true, "", "latest", "latest"},
		{"nothing specified", "", false, "", "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := resolveFunctionVersion(c.configuredVersion, c.isConfigured, c.urnVersion, c.currentVersion)
			if actual != c.expected {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}

func TestSuppressCustomAuthorizerFunctionUrnDiff(t *testing.T) {
	cases := []struct {
		name            string
		old             string
		new             string
		functionVersion string
		expected        bool
	}{
		{"state without version, config with the same version", testFunctionUrn, testFunctionUrn + ":latest", "latest", true},
		{"state with version, config without version", testFunctionUrn + ":latest", testFunctionUrn, "latest", true},
		{"version suffix differs from function_version", testFunctionUrn, testFunctionUrn + ":v2", "latest", false},
		{"both URNs with different versions", testFunctionUrn + ":latest", testFunctionUrn + ":v2", "latest", false},
		{"different functions with common prefix", testFunctionUrn, testFunctionUrn + "_v2:latest", "latest", false},
		{"different functions", testFunctionUrn, testFunctionUrn + "x", "latest", false},
		{"new resource", "", testFunctionUrn + ":latest", "latest", false},
	}

	resourceSchema := ResourceAPICustomAuthorizerV2().Schema
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
				"function_version": c.functionVersion,
			})
			actual := suppressCustomAuthorizerFunctionUrnDiff("function_urn", c.old, c.new, d)
			if actual != c.expected {
				t.Fatalf("expected %v, got %v", c.expected, actual)
			}
		})
	}
}

func TestFlattenCustomAuthorizerFunctionUrn(t *testing.T) {
	resp := &authorizer.AuthorizerResp{
		FunctionUrn: testFunctionUrn,
		Version:     "latest",
	}
	cases := []struct {
		name     string
		current  string
		expected string
	}{
		{"import", "", testFunctionUrn},
		{"without version", testFunctionUrn, testFunctionUrn},
		{"with the same version", testFunctionUrn + ":latest", testFunctionUrn + ":latest"},
		{"with a different version", testFunctionUrn + ":v2", testFunctionUrn},
		{"different function", testFunctionUrn + "x:latest", testFunctionUrn},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := flattenCustomAuthorizerFunctionUrn(c.current, resp)
			if actual != c.expected {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}

func TestBuildCustomAuthorizerOptsFunctionVersion(t *testing.T) {
	cases := []struct {
		name            string
		raw             map[string]interface{}
		expectedUrn     string
		expectedVersion string
		expectedAlias   string
	}{
		{
			name: "separate version",
			raw: map[string]interface{}{
				"function_urn":     testFunctionUrn,
				"function_version": "latest",
			},
			expectedUrn:     testFunctionUrn,
			expectedVersion: "latest",
		},
		{
			name: "version suffix in the URN",
			raw: map[string]interface{}{
				"function_urn": testFunctionUrn + ":latest",
			},
			expectedUrn:     testFunctionUrn,
			expectedVersion: "latest",
		},
		{
			name: "alias URN",
			raw: map[string]interface{}{
				"function_urn":       testFunctionUrn,
				"function_alias_uri": testFunctionUrn + ":!prod",
			},
			expectedUrn:   testFunctionUrn,
			expectedAlias: testFunctionUrn + ":!prod",
		},
	}

	resourceSchema := ResourceAPICustomAuthorizerV2().Schema
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.raw["gateway_id"] = "gateway"
			c.raw["name"] = "authorizer"
			d := schema.TestResourceDataRaw(t, resourceSchema, c.raw)
			opts, err := buildCustomAuthorizerOpts(d)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if opts.FunctionUrn != c.expectedUrn || opts.Version != c.expectedVersion || opts.AliasUrn != c.expectedAlias {
				t.Fatalf("expected (%q, %q, %q), got (%q, %q, %q)", c.expectedUrn, c.expectedVersion, c.expectedAlias,
					opts.FunctionUrn, opts.Version, opts.AliasUrn)
			}
		})
	}
}

func testAPIGWAuthorizerConfig(t *testing.T, handler http.HandlerFunc) *cfg.Config {
	t.Helper()
	server := httptest.NewServer(handler)
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

func TestCustomAuthorizerCreateAndReadWithVersionSuffix(t *testing.T) {
	var createPayload map[string]interface{}

	conf := testAPIGWAuthorizerConfig(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/apigw/instances/gw-1/authorizers":
			if err := json.NewDecoder(r.Body).Decode(&createPayload); err != nil {
				t.Errorf("failed to decode create body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"id": "auth-1",
				"name": "my_authorizer",
				"type": "FRONTEND",
				"authorizer_type": "FUNC",
				"authorizer_uri": "` + testFunctionUrn + `",
				"authorizer_version": "latest",
				"ttl": 60,
				"need_body": false
			}`))
		case r.Method == http.MethodGet && r.URL.Path == "/apigw/instances/gw-1/authorizers/auth-1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "auth-1",
				"name": "my_authorizer",
				"type": "FRONTEND",
				"authorizer_type": "FUNC",
				"authorizer_uri": "` + testFunctionUrn + `",
				"authorizer_version": "latest",
				"ttl": 60,
				"need_body": false
			}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})

	r := ResourceAPICustomAuthorizerV2()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"gateway_id":   "gw-1",
		"name":         "my_authorizer",
		"function_urn": testFunctionUrn + ":latest",
		"type":         "FRONTEND",
		"ttl":          60,
	})

	diags := r.CreateContext(context.Background(), d, conf)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if createPayload["authorizer_uri"] != testFunctionUrn {
		t.Fatalf("expected authorizer_uri %q, got %v", testFunctionUrn, createPayload["authorizer_uri"])
	}
	if createPayload["authorizer_version"] != "latest" {
		t.Fatalf("expected authorizer_version %q, got %v", "latest", createPayload["authorizer_version"])
	}
	if got := d.Get("function_version").(string); got != "latest" {
		t.Fatalf("expected state function_version %q, got %q", "latest", got)
	}
	if got := d.Get("function_urn").(string); got != testFunctionUrn+":latest" {
		t.Fatalf("expected state function_urn %q, got %q", testFunctionUrn+":latest", got)
	}
}

func TestCustomAuthorizerImportFiltersByExactName(t *testing.T) {
	conf := testAPIGWAuthorizerConfig(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/apigw/instances/gw-1/authorizers" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"authorizer_list": [
					{"id": "auth-other", "name": "other_authorizer"},
					{"id": "auth-target", "name": "target_authorizer"}
				]
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	r := ResourceAPICustomAuthorizerV2()
	d := r.TestResourceData()
	d.SetId("gw-1/target_authorizer")

	states, err := resourceCustomAuthorizerImportState(context.Background(), d, conf)
	if err != nil {
		t.Fatalf("unexpected import error: %s", err)
	}
	if len(states) != 1 || states[0].Id() != "auth-target" {
		t.Fatalf("expected imported ID %q, got %q", "auth-target", d.Id())
	}
	if got := d.Get("gateway_id").(string); got != "gw-1" {
		t.Fatalf("expected gateway_id %q, got %q", "gw-1", got)
	}

	dMissing := r.TestResourceData()
	dMissing.SetId("gw-1/missing_authorizer")
	if _, err := resourceCustomAuthorizerImportState(context.Background(), dMissing, conf); err == nil {
		t.Fatal("expected error when authorizer name is not found")
	}
}
