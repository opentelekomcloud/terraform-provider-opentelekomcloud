package v3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/monitors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const monitorV3TestJSON = `{
	"id":"monitor","pool_id":"ignored","name":"health","project_id":"project","type":"HTTP",
	"delay":10,"timeout":5,"max_retries":3,"max_retries_down":2,"http_method":"GET",
	"url_path":"/health","domain_name":"example.com","expected_codes":"200-299",
	"admin_state_up":false,"monitor_port":8080,"pools":[{"id":"pool"}],
	"created_at":"2026-10-28T10:00:00Z","updated_at":"2026-10-28T10:01:00Z"
}`

type monitorTestClientContext struct {
	context.Context
	client *golangsdk.ServiceClient
}

func (c monitorTestClientContext) Value(key interface{}) interface{} {
	if key == keyClient {
		return c.client
	}
	return c.Context.Value(key)
}

func monitorTestContext(t *testing.T, handler http.HandlerFunc) context.Context {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return monitorTestClientContext{
		Context: context.Background(),
		client: &golangsdk.ServiceClient{
			ProviderClient: &golangsdk.ProviderClient{},
			Endpoint:       server.URL + "/v3/project/elb/",
		},
	}
}

func TestMonitorV3SchemaAndFlatten(t *testing.T) {
	r, ds := ResourceMonitorV3(), DataSourceLBMonitorV3()
	provider := &schema.Provider{
		ResourcesMap:   map[string]*schema.Resource{"opentelekomcloud_lb_monitor_v3": r},
		DataSourcesMap: map[string]*schema.Resource{"opentelekomcloud_lb_monitor_v3": ds},
	}
	if err := provider.InternalValidate(); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{
		"pool_id", "type", "admin_state_up", "delay", "timeout", "max_retries",
		"max_retries_down", "url_path", "domain_name", "http_method", "expected_codes",
		"project_id", "name", "monitor_port", "created_at", "updated_at",
	} {
		if r.Schema[field] == nil || ds.Schema[field] == nil {
			t.Errorf("missing field %s", field)
		}
	}
	for _, field := range []string{"http_version", "content_check_enable", "send_content", "receive_content", "content_format", "enterprise_project_id"} {
		if r.Schema[field] != nil || ds.Schema[field] != nil {
			t.Errorf("unsupported field exposed: %s", field)
		}
	}
	if r.Schema["admin_state_up"].Default != true || ds.Schema["admin_state_up"].Default != nil {
		t.Fatal("administrative state defaults changed or leak into filters")
	}
	if !r.Schema["pool_id"].ForceNew || !r.Schema["project_id"].ForceNew || r.Schema["type"].ForceNew {
		t.Fatal("existing replacement behavior changed")
	}
	if _, errs := r.Schema["type"].ValidateFunc("PING", "type"); len(errs) != 0 {
		t.Fatal("legacy PING protocol removed")
	}
	var monitor monitors.Monitor
	if err := json.Unmarshal([]byte(monitorV3TestJSON), &monitor); err != nil {
		t.Fatal(err)
	}
	d := schema.TestResourceDataRaw(t, r.Schema, nil)
	if diags := setMonitorV3Fields(d, &monitor); diags.HasError() {
		t.Fatal(diags)
	}
	for key, expected := range map[string]interface{}{
		"pool_id": "pool", "name": "health", "project_id": "project", "type": "HTTP",
		"admin_state_up": false, "delay": 10, "timeout": 5, "max_retries": 3, "max_retries_down": 2,
		"monitor_port": 8080, "domain_name": "example.com", "url_path": "/health",
		"http_method": "GET", "expected_codes": "200-299",
		"created_at": "2026-10-28T10:00:00Z", "updated_at": "2026-10-28T10:01:00Z",
	} {
		if actual := d.Get(key); actual != expected {
			t.Errorf("%s: got %#v, want %#v", key, actual, expected)
		}
	}
	monitor.CreatedAt, monitor.UpdatedAt = "", ""
	if diags := setMonitorV3Fields(d, &monitor); diags.HasError() {
		t.Fatal(diags)
	}
	if d.Get("created_at") != "" || d.Get("updated_at") != "" {
		t.Fatal("stale timestamps retained")
	}
	d.SetId("monitor")
	imported, err := r.Importer.StateContext(context.Background(), d, nil)
	if err != nil || len(imported) != 1 || imported[0].Id() != "monitor" {
		t.Fatalf("import ID format changed: %v", err)
	}
}

func TestMonitorV3Filters(t *testing.T) {
	values := map[string]string{
		"name": "health", "type": "HTTP", "delay": "10", "timeout": "5", "max_retries": "3",
		"max_retries_down": "2", "monitor_port": "8080", "domain_name": "example.com",
		"url_path": "/health", "expected_codes": "200-299", "admin_state_up": "false",
	}
	config := map[string]cty.Value{}
	for _, key := range []string{"name", "type", "domain_name", "url_path", "expected_codes"} {
		config[key] = cty.StringVal(values[key])
	}
	for key, value := range map[string]int64{"delay": 10, "timeout": 5, "max_retries": 3, "max_retries_down": 2, "monitor_port": 8080} {
		config[key] = cty.NumberIntVal(value)
	}
	config["admin_state_up"] = cty.False
	d := memberTestData(t, DataSourceLBMonitorV3(), values, config)
	falseValue := false
	expected := monitors.ListOpts{
		Name: "health", Type: "HTTP", Delay: 10, Timeout: 5, MaxRetries: 3,
		MaxRetriesDown: []int{2}, MonitorPort: []int{8080}, DomainName: []string{"example.com"},
		URLPath: "/health", ExpectedCodes: "200-299", AdminStateUp: &falseValue,
	}
	if actual := buildMonitorV3ListOpts(d); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("filter mappings differ: got %#v, want %#v", actual, expected)
	}
	unconfigured := memberTestData(t, DataSourceLBMonitorV3(), values, map[string]cty.Value{})
	if actual := buildMonitorV3ListOpts(unconfigured); !reflect.DeepEqual(actual, monitors.ListOpts{}) {
		t.Fatalf("computed state must not become filters: %#v", actual)
	}
}

func TestMonitorV3Lifecycle(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceMonitorV3().Schema, map[string]interface{}{
		"pool_id": "pool", "type": "HTTP", "admin_state_up": false, "delay": 10, "timeout": 5,
		"max_retries": 3, "max_retries_down": 2, "url_path": "/health", "domain_name": "example.com",
		"http_method": "GET", "expected_codes": "200-299", "project_id": "project", "name": "health", "monitor_port": 8080,
	})
	var methods []string
	ctx := monitorTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			var body map[string]map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			expected := map[string]interface{}{
				"type": "HTTP", "admin_state_up": false, "delay": float64(10), "timeout": float64(5),
				"max_retries": float64(3), "max_retries_down": float64(2), "url_path": "/health",
				"domain_name": "example.com", "http_method": "GET", "expected_codes": "200-299",
				"name": "health", "monitor_port": float64(8080),
			}
			if r.Method == http.MethodPost {
				expected["pool_id"], expected["project_id"] = "pool", "project"
				w.WriteHeader(http.StatusCreated)
			}
			if !reflect.DeepEqual(body["healthmonitor"], expected) {
				t.Errorf("request mappings changed: %#v", body)
			}
		}
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_, _ = fmt.Fprintf(w, `{"healthmonitor":%s}`, monitorV3TestJSON)
	})
	config := &cfg.Config{}
	if diags := resourceMonitorV3Create(ctx, d, config); diags.HasError() {
		t.Fatal(diags)
	}
	if d.Id() != "monitor" || d.Get("created_at") != "2026-10-28T10:00:00Z" {
		t.Fatal("monitor ID or metadata lost")
	}
	if diags := resourceMonitorV3Update(ctx, d, config); diags.HasError() {
		t.Fatal(diags)
	}
	if diags := resourceMonitorV3Delete(ctx, d, config); diags.HasError() {
		t.Fatal(diags)
	}
	if !reflect.DeepEqual(methods, []string{"POST", "GET", "PUT", "GET", "DELETE"}) {
		t.Errorf("unexpected lifecycle: %v", methods)
	}
}

func TestMonitorV3DataSourceRead(t *testing.T) {
	for _, test := range []struct {
		name       string
		id         bool
		body       string
		status     int
		errorText  string
		pagination bool
	}{
		{name: "id", id: true, body: `{"healthmonitor":` + monitorV3TestJSON + `}`},
		{name: "filters", body: `{"healthmonitors":[` + monitorV3TestJSON + `],"page_info":{}}`},
		{name: "empty", body: `{"healthmonitors":[]}`, errorText: "no results"},
		{name: "multiple", body: `{"healthmonitors":[` + monitorV3TestJSON + `,` + monitorV3TestJSON + `]}`, errorText: "more than one"},
		{name: "later-page-match", body: `{"healthmonitors":[` + monitorV3TestJSON + `],"page_info":{"next_marker":"next"}}`, errorText: "more than one", pagination: true},
		{name: "no-pool", id: true, body: `{"healthmonitor":{"id":"monitor","pools":[]}}`},
		{name: "not-found", id: true, status: http.StatusNotFound, errorText: "retrieving"},
		{name: "forbidden", status: http.StatusForbidden, errorText: "retrieving"},
		{name: "missing-id", id: true, body: `{"healthmonitor":{"name":"health"}}`, errorText: "ID"},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := map[string]cty.Value{"name": cty.StringVal("health"), "admin_state_up": cty.False}
			if test.id {
				raw["id"] = cty.StringVal("monitor")
			}
			d := memberTestData(t, DataSourceLBMonitorV3(), map[string]string{
				"id": "old-monitor", "name": "health", "admin_state_up": "false", "pool_id": "stale-pool",
			}, raw)
			if test.id {
				if err := d.Set("id", "monitor"); err != nil {
					t.Fatal(err)
				}
			}
			sawNext := false
			ctx := monitorTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				if test.id {
					if !strings.HasSuffix(r.URL.Path, "/healthmonitors/monitor") || r.URL.RawQuery != "" {
						t.Errorf("ID lookup did not bypass filters: %s", r.URL)
					}
				} else if !strings.HasSuffix(r.URL.Path, "/healthmonitors") || r.URL.Query().Get("name") != "health" || r.URL.Query().Get("admin_state_up") != "false" {
					t.Errorf("computed ID bypassed configured filters: %s", r.URL)
				}
				if test.status != 0 {
					w.WriteHeader(test.status)
					return
				}
				if test.pagination && r.URL.Query().Get("marker") == "next" {
					sawNext = true
					_, _ = fmt.Fprintf(w, `{"healthmonitors":[%s],"page_info":{}}`, strings.ReplaceAll(monitorV3TestJSON, `"monitor"`, `"second"`))
					return
				}
				_, _ = fmt.Fprint(w, test.body)
			})
			diags := dataSourceLBMonitorV3Read(ctx, d, &cfg.Config{})
			if test.errorText != "" {
				if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), test.errorText) {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
			} else {
				expectedPool := "pool"
				if test.name == "no-pool" {
					expectedPool = ""
				}
				if diags.HasError() || d.Id() != "monitor" || d.Get("pool_id") != expectedPool {
					t.Fatalf("lookup did not populate monitor or clear stale pool: %v", diags)
				}
			}
			if test.pagination && !sawNext {
				t.Fatal("uniqueness was not checked beyond the first page")
			}
		})
	}
}

func TestMonitorV3ResourceNotFound(t *testing.T) {
	for _, code := range []int{http.StatusNotFound, http.StatusForbidden} {
		for _, method := range []string{http.MethodGet, http.MethodDelete} {
			t.Run(fmt.Sprintf("%s-%d", method, code), func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, ResourceMonitorV3().Schema, nil)
				d.SetId("monitor")
				ctx := monitorTestContext(t, func(w http.ResponseWriter, r *http.Request) {
					if r.Method != method {
						t.Errorf("unexpected method %s", r.Method)
					}
					w.WriteHeader(code)
				})
				op := resourceMonitorV3Read
				if method == http.MethodDelete {
					op = resourceMonitorV3Delete
				}
				diags := op(ctx, d, &cfg.Config{})
				if code == http.StatusNotFound {
					if diags.HasError() || d.Id() != "" {
						t.Fatalf("404 did not remove resource: %v", diags)
					}
				} else if !diags.HasError() || d.Id() != "monitor" {
					t.Fatalf("API error swallowed: %v", diags)
				}
			})
		}
	}
}
