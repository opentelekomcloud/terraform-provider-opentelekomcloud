package waf

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestWafDedicatedPreciseProtectionRuleV1Contents(t *testing.T) {
	for _, tc := range []struct {
		name          string
		reference     bool
		configured    bool
		create        bool
		imported      bool
		tableUpdate   bool
		editContents  bool
		editReference bool
		apiContents   string
		wantNew       bool
	}{
		{name: "reference_create", reference: true, create: true, apiContents: "/admin"},
		{name: "reference_refresh", reference: true, apiContents: "/admin"},
		{name: "reference_table_update", reference: true, tableUpdate: true, apiContents: "/admin"},
		{name: "reference_import", reference: true, imported: true, apiContents: "/admin"},
		{name: "reference_with_explicit_contents", reference: true, configured: true, apiContents: "/admin"},
		{name: "reference_id_edit", reference: true, editReference: true, apiContents: "/admin", wantNew: true},
		{name: "inline_unchanged", configured: true, apiContents: "/admin"},
		{name: "inline_remote_drift", configured: true, apiContents: "/changed", wantNew: true},
		{name: "inline_config_edit", configured: true, editContents: true, apiContents: "/admin", wantNew: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := ResourceWafDedicatedPreciseProtectionRuleV1()
			condition := map[string]interface{}{
				"category": "url", "logic_operation": "contain",
			}
			if tc.reference {
				condition["logic_operation"] = "contain_any"
				condition["value_list_id"] = "reference"
			}
			if tc.configured {
				condition["contents"] = []interface{}{"/admin"}
			}
			config := map[string]interface{}{
				"policy_id": "policy", "time": false, "priority": 50,
				"action":     []interface{}{map[string]interface{}{"category": "block"}},
				"conditions": []interface{}{condition},
			}
			rawState := config
			if tc.imported {
				rawState = map[string]interface{}{}
			}
			d := schema.TestResourceDataRaw(t, r.Schema, rawState)
			d.SetId("rule")
			if tc.imported {
				d.SetId("policy/rule")
				imported, err := r.Importer.StateContext(context.Background(), d, nil)
				if err != nil {
					t.Fatal(err)
				}
				if len(imported) != 1 || imported[0].Id() != "rule" || imported[0].Get("policy_id") != "policy" {
					t.Fatalf("unexpected imported state: %#v", imported)
				}
				d = imported[0]
			}

			var reads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				contents := tc.apiContents
				switch req.Method {
				case http.MethodGet:
					if !strings.HasSuffix(req.URL.Path, "/waf/policy/policy/custom/rule") {
						t.Errorf("unexpected GET path: %s", req.URL.Path)
					}
					if reads.Add(1) > 1 && tc.tableUpdate {
						contents = "/changed"
					}
				case http.MethodPost:
					if !tc.create || !strings.HasSuffix(req.URL.Path, "/waf/policy/policy/custom") {
						t.Errorf("unexpected POST path: %s", req.URL.Path)
					}
					var body struct {
						Time       *bool                    `json:"time"`
						Conditions []map[string]interface{} `json:"conditions"`
					}
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Error(err)
						http.Error(w, "invalid test create payload", http.StatusBadRequest)
						return
					}
					if body.Time == nil || *body.Time || len(body.Conditions) != 1 {
						t.Errorf("unexpected create options: %#v", body)
					} else {
						if _, exists := body.Conditions[0]["contents"]; exists {
							t.Error("unconfigured contents must not be sent during create")
						}
						if body.Conditions[0]["value_list_id"] != "reference" {
							t.Errorf("reference ID not preserved: %#v", body.Conditions[0])
						}
					}
				default:
					t.Errorf("unexpected request method: %s", req.Method)
				}
				responseCondition := map[string]interface{}{
					"category": "url", "logic_operation": "contain",
					"contents": []string{contents},
				}
				if tc.reference {
					responseCondition["logic_operation"] = "contain_any"
					responseCondition["value_list_id"] = "reference"
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"id": "rule", "policyid": "policy", "priority": 50,
					"status": 1, "timestamp": 123, "start": 0, "terminal": 0,
					"action":     map[string]interface{}{"category": "block"},
					"conditions": []interface{}{responseCondition},
				}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			meta := &cfg.Config{Region: "eu-de", HwClient: &golangsdk.ProviderClient{
				EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
					return server.URL + "/", nil
				},
			}}
			var diags diag.Diagnostics
			if tc.create {
				d.SetId("")
				diags = r.CreateContext(context.Background(), d, meta)
			} else {
				diags = r.ReadContext(context.Background(), d, meta)
			}
			if diags.HasError() {
				t.Fatal(diags)
			}
			expectedContents := tc.apiContents
			if tc.tableUpdate {
				if diags := r.ReadContext(context.Background(), d, meta); diags.HasError() {
					t.Fatal(diags)
				}
				expectedContents = "/changed"
			}
			if got := d.State().Attributes["conditions.0.contents.0"]; got != expectedContents {
				t.Fatalf("API contents missing from state: got %q, want %q", got, expectedContents)
			}
			if tc.editContents {
				condition["contents"] = []interface{}{"/edited"}
			}
			if tc.editReference {
				condition["value_list_id"] = "different-reference"
			}
			diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), meta)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantNew {
				if diff == nil || !diff.RequiresNew() {
					t.Fatalf("explicit input changes must require replacement: %#v", diff)
				}
			} else if diff != nil && !diff.Empty() {
				t.Fatalf("expected an unchanged plan: %#v", diff)
			}
		})
	}
}
