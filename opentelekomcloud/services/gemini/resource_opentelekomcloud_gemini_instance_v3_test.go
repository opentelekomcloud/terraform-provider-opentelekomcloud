package gemini

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/instance"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func geminiRetentionTestData(t *testing.T, keepDays cty.Value, withStrategy bool) *schema.ResourceData {
	t.Helper()
	r := ResourceGeminiDBInstanceV3()
	raw := map[string]interface{}{
		"name": "test-gemini", "password": "TestOnlyPassword9!",
		"availability_zone": "eu-de-01,eu-de-02,eu-de-03",
		"vpc_id":            "test-vpc", "subnet_id": "test-subnet", "security_group_id": "test-security-group",
		"flavor": "geminidb.cassandra.xlarge.8", "volume_size": 100,
	}
	config := map[string]cty.Value{}
	for key, value := range raw {
		if number, ok := value.(int); ok {
			config[key] = cty.NumberIntVal(int64(number))
		} else {
			config[key] = cty.StringVal(value.(string))
		}
	}
	if withStrategy {
		strategy := map[string]interface{}{"start_time": "02:00-03:00"}
		if keepDays.IsKnown() && !keepDays.IsNull() {
			value, _ := keepDays.AsBigFloat().Int64()
			strategy["keep_days"] = int(value)
		}
		raw["backup_strategy"] = []interface{}{strategy}
		config["backup_strategy"] = cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"start_time": cty.StringVal("02:00-03:00"), "keep_days": keepDays,
			}),
		})
	}
	data := schema.TestResourceDataRaw(t, r.Schema, raw)
	data.SetId("fixture")
	d, err := schema.InternalMap(r.Schema).Data(data.State(), &terraform.InstanceDiff{RawConfig: cty.ObjectVal(config)})
	if err != nil {
		t.Fatal(err)
	}
	d.SetId("")
	return d
}

func geminiRetentionTestConfig(t *testing.T, handler http.HandlerFunc) *cfg.Config {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &cfg.Config{
		Region: "eu-de",
		HwClient: &golangsdk.ProviderClient{
			TokenID: "test-token", ProjectID: "test-project", HTTPClient: *server.Client(),
			EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
				return server.URL + "/v3/test-project/", nil
			},
		},
	}
}

func TestGeminiDBBackupRetentionRequest(t *testing.T) {
	for _, tc := range []struct {
		name         string
		days         cty.Value
		withStrategy bool
		wantBackup   string
		responseDays int
	}{
		{"configured", cty.NumberIntVal(14), true, `{"start_time":"02:00-03:00","keep_days":14}`, 14},
		{"disabled", cty.NumberIntVal(0), true, `{"start_time":"02:00-03:00","keep_days":0}`, 0},
		{"maximum", cty.NumberIntVal(35), true, `{"start_time":"02:00-03:00","keep_days":35}`, 35},
		{"default_retention", cty.NullVal(cty.Number), true, `{"start_time":"02:00-03:00"}`, 7},
		{"default_strategy", cty.NullVal(cty.Number), false, "", 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := geminiRetentionTestData(t, tc.days, tc.withStrategy)
			config := geminiRetentionTestConfig(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/v3/test-project/instances" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.Header.Get("X-Auth-Token") != "test-token" {
					t.Error("missing SDK authentication header")
				}
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				backup, present := body["backup_strategy"]
				if present != tc.withStrategy {
					t.Errorf("backup_strategy presence: got %t, want %t", present, tc.withStrategy)
				}
				if present {
					var actual, expected interface{}
					if err := json.Unmarshal(backup, &actual); err != nil {
						t.Error(err)
					}
					if err := json.Unmarshal([]byte(tc.wantBackup), &expected); err != nil {
						t.Error(err)
						return
					}
					if !reflect.DeepEqual(actual, expected) {
						t.Errorf("backup request: got %s, want %s", backup, tc.wantBackup)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = fmt.Fprintf(w, `{"id":"test-instance","backup_strategy":{"start_time":"02:00-03:00","keep_days":%d}}`, tc.responseDays)
			})
			client, err := config.GeminiDBV3Client("eu-de")
			if err != nil {
				t.Fatal(err)
			}
			opts := instance.CreateOpts{
				Name: d.Get("name").(string), Password: d.Get("password").(string),
				Region: "eu-de", AvailabilityZone: d.Get("availability_zone").(string),
				VpcId: d.Get("vpc_id").(string), SubnetId: d.Get("subnet_id").(string),
				SecurityGroupId: d.Get("security_group_id").(string), Mode: "Cluster",
				DataStore:      resourceGeminiDBDataStore(d, geminiDBDefaults(d)),
				Flavor:         resourceGeminiDBFlavor(d),
				BackupStrategy: resourceGeminiDBBackupStrategy(d),
			}
			result, err := instance.Create(client, opts)
			if err != nil {
				t.Fatal(err)
			}
			if result.Id != "test-instance" || result.BackupStrategy.KeepDays != tc.responseDays {
				t.Fatalf("unexpected create response: %#v", result)
			}
		})
	}
}

func TestGeminiDBCreateRetainsIDAfterNumericBackupResponse(t *testing.T) {
	d := geminiRetentionTestData(t, cty.NumberIntVal(14), true)
	config := geminiRetentionTestConfig(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v3/test-project/instances" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `{"id":"created-instance","backup_strategy":{"start_time":"02:00-03:00","keep_days":14}}`)
	})
	// Stop only the readiness wait, after the real create request decodes and saves its ID.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	diagnostics := ResourceGeminiDBInstanceV3().CreateContext(ctx, d, config)
	if d.Id() != "created-instance" {
		t.Fatalf("numeric retention response lost the created instance ID: %q, %v", d.Id(), diagnostics)
	}
	if !diagnostics.HasError() || !strings.Contains(fmt.Sprint(diagnostics), "error waiting for instance") {
		t.Fatalf("expected the cancelled readiness wait, got %v", diagnostics)
	}
}

func TestGeminiDBCreateBackupResponseErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"api_error", http.StatusBadRequest, `{"error_msg":"invalid request"}`},
		{"invalid_retention", http.StatusAccepted, `{"id":"test-instance","backup_strategy":{"keep_days":true}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := geminiRetentionTestData(t, cty.NumberIntVal(14), true)
			config := geminiRetentionTestConfig(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprint(w, tc.body)
			})
			diagnostics := ResourceGeminiDBInstanceV3().CreateContext(context.Background(), d, config)
			if !diagnostics.HasError() || !strings.Contains(fmt.Sprint(diagnostics), "error creating GeminiDB instance") {
				t.Fatalf("create error must propagate: %v", diagnostics)
			}
			if d.Id() != "" {
				t.Fatalf("failed create unexpectedly set an ID: %s", d.Id())
			}
		})
	}
}

func TestGeminiDBBackupRetentionReadAndImport(t *testing.T) {
	for _, days := range []int{0, 7, 14, 35} {
		t.Run(fmt.Sprint(days), func(t *testing.T) {
			r := ResourceGeminiDBInstanceV3()
			d := r.TestResourceData()
			d.SetId("test-instance")
			config := geminiRetentionTestConfig(t, func(w http.ResponseWriter, request *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch request.URL.Path {
				case "/v3/test-project/instances":
					if request.Method != http.MethodGet || request.URL.Query().Get("id") != "test-instance" {
						t.Errorf("unexpected instance lookup: %s %s", request.Method, request.URL)
					}
					_, _ = fmt.Fprintf(w, `{"instances":[{"id":"test-instance","name":"test-gemini","status":"normal","backup_strategy":{"start_time":"02:00-03:00","keep_days":%d}}]}`, days)
				case "/v3/test-project/instances/test-instance/tags":
					_, _ = fmt.Fprint(w, `{"tags":[]}`)
				default:
					t.Errorf("unexpected request: %s %s", request.Method, request.URL)
					w.WriteHeader(http.StatusNotFound)
				}
			})
			imported, err := r.Importer.StateContext(context.Background(), d, config)
			if err != nil || len(imported) != 1 {
				t.Fatalf("unexpected import: %#v, %v", imported, err)
			}
			diagnostics := r.ReadContext(context.Background(), imported[0], config)
			if diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if d.Id() != "test-instance" || d.Get("backup_strategy.0.keep_days") != days {
				t.Fatalf("integer retention was not saved to state: %s, %#v", d.Id(), d.Get("backup_strategy"))
			}
		})
	}
}
