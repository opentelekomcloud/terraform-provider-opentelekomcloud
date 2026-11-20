package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func backupWindowCreateTestConfig(days int) map[string]interface{} {
	return map[string]interface{}{
		"name": "test-gemini", "password": "TestOnlyPassword9!",
		"availability_zone": "eu-de-01,eu-de-02,eu-de-03",
		"vpc_id":            "test-vpc", "subnet_id": "test-subnet", "security_group_id": "test-group",
		"flavor": "geminidb.cassandra.xlarge.8", "volume_size": 100,
		"backup_strategy": []interface{}{
			map[string]interface{}{"start_time": "03:00-04:00", "keep_days": days},
		},
	}
}

func backupWindowCreateTestRawConfig(days int) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"backup_strategy": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"start_time": cty.StringVal("03:00-04:00"), "keep_days": cty.NumberIntVal(int64(days)),
			}),
		}),
	})
}

func TestGeminiDBBackupWindowCreatePlan(t *testing.T) {
	for _, days := range []int{0, 7} {
		t.Run(map[int]string{0: "disabled", 7: "enabled"}[days], func(t *testing.T) {
			r := ResourceGeminiDBInstanceV3()
			// The protocol shim supplies raw config even when the prior state has no ID.
			rawConfig := backupWindowCreateTestRawConfig(days)
			state := &terraform.InstanceState{RawConfig: rawConfig}
			diff, err := r.SimpleDiff(context.Background(), state, terraform.NewResourceConfigRaw(backupWindowCreateTestConfig(days)), nil)
			if err != nil {
				t.Fatal(err)
			}
			window := diff.Attributes["backup_strategy.0.start_time"]
			if window == nil || window.New != "03:00-04:00" {
				t.Fatalf("create plan dropped the required backup window: %#v", window)
			}
			diff.RawConfig = rawConfig
			d, err := schema.InternalMap(r.Schema).Data(nil, diff)
			if err != nil {
				t.Fatal(err)
			}
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests++
				if request.Method != http.MethodPost || request.URL.Path != "/v3/project/instances" {
					t.Errorf("unexpected request: %s %s", request.Method, request.URL)
				}
				var body struct {
					Backup struct {
						StartTime string `json:"start_time"`
						KeepDays  *int   `json:"keep_days"`
					} `json:"backup_strategy"`
				}
				if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body.Backup.StartTime != "03:00-04:00" || body.Backup.KeepDays == nil || *body.Backup.KeepDays != days {
					t.Errorf("invalid create backup strategy: %#v", body.Backup)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = fmt.Fprintf(w, `{"id":"test-instance","backup_strategy":{"keep_days":%d}}`, days)
			}))
			defer server.Close()
			meta := &cfg.Config{
				Region: "eu-de",
				HwClient: &golangsdk.ProviderClient{
					ProjectID: "project", HTTPClient: *server.Client(),
					EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
						return server.URL + "/v3/project/", nil
					},
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			// Cancel only the readiness wait, after Create sends the actual planned request.
			diagnostics := r.CreateContext(ctx, d, meta)
			if requests != 1 || d.Id() != "test-instance" {
				t.Fatalf("planned create did not succeed: requests %d, ID %q, diagnostics %v", requests, d.Id(), diagnostics)
			}
			if !diagnostics.HasError() || !strings.Contains(fmt.Sprint(diagnostics), "error waiting for instance") {
				t.Fatalf("expected only the cancelled readiness wait: %v", diagnostics)
			}
		})
	}
}

func TestGeminiDBBackupWindowExistingPlan(t *testing.T) {
	for _, days := range []int{0, 7} {
		t.Run(map[int]string{0: "disabled", 7: "enabled"}[days], func(t *testing.T) {
			r := ResourceGeminiDBInstanceV3()
			d := schema.TestResourceDataRaw(t, r.Schema, backupWindowCreateTestConfig(days))
			d.SetId("existing-instance")
			if err := d.Set("backup_strategy", []interface{}{
				map[string]interface{}{"start_time": "00:00-01:00", "keep_days": days},
			}); err != nil {
				t.Fatal(err)
			}
			state := d.State()
			for _, key := range []string{"datastore.#", "nodes.#", "private_ips.#"} {
				state.Attributes[key] = "0"
			}
			state.RawConfig = backupWindowCreateTestRawConfig(days)
			diff, err := r.SimpleDiff(context.Background(), state,
				terraform.NewResourceConfigRaw(backupWindowCreateTestConfig(days)), nil)
			if err != nil {
				t.Fatal(err)
			}
			if days == 0 {
				if diff != nil && !diff.Empty() {
					t.Fatalf("disabled backups must keep a no-op plan: %#v", diff)
				}
			} else if diff == nil || diff.Attributes["backup_strategy.0.start_time"] == nil ||
				diff.Attributes["backup_strategy.0.start_time"].New != "03:00-04:00" {
				t.Fatalf("enabled backup window change was suppressed: %#v", diff)
			}
		})
	}
}
