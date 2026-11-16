package as

import (
	"context"
	"encoding/json"
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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/notifications"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const asNotificationTestURN = "urn:smn:eu-de:project:test-topic"
const asNotificationTestJSON = `{"topic_urn":"urn:smn:eu-de:project:test-topic","topic_name":"test-topic","topic_scene":["SCALING_DOWN","SCALING_UP"]}`

type asNotificationTestClientContext struct {
	context.Context
	client *golangsdk.ServiceClient
}

func (c asNotificationTestClientContext) Value(key interface{}) interface{} {
	if key == keyClientV1 {
		return c.client
	}
	return c.Context.Value(key)
}

func asNotificationTestContext(t *testing.T, handler http.HandlerFunc) context.Context {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return asNotificationTestClientContext{
		Context: context.Background(),
		client: &golangsdk.ServiceClient{
			ProviderClient: &golangsdk.ProviderClient{TokenID: "test-token"},
			Endpoint:       server.URL + "/",
		},
	}
}

func asNotificationTestData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, ResourceASNotificationV1().Schema, map[string]interface{}{
		"scaling_group_id": "group",
		"topic_urn":        asNotificationTestURN,
		"topic_scene":      []interface{}{"SCALING_UP", "SCALING_DOWN"},
	})
	return d
}

func checkASNotificationRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path || r.URL.RawQuery != "" {
		t.Errorf("unexpected request %s %s, expected %s %s", r.Method, r.URL, method, path)
	}
	if r.Header.Get("X-Auth-Token") != "test-token" {
		t.Error("missing SDK authentication header")
	}
}

func TestASNotificationV1Schema(t *testing.T) {
	res, ds := ResourceASNotificationV1(), DataSourceASNotificationsV1()
	provider := &schema.Provider{
		ResourcesMap:   map[string]*schema.Resource{"opentelekomcloud_as_notification_v1": res},
		DataSourcesMap: map[string]*schema.Resource{"opentelekomcloud_as_notifications_v1": ds},
	}
	if err := provider.InternalValidate(); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"scaling_group_id", "topic_urn", "topic_scene"} {
		if !res.Schema[field].Required || !res.Schema[field].ForceNew {
			t.Errorf("%s must be required and replace the association when changed", field)
		}
	}
	if res.UpdateContext != nil || res.Importer == nil || res.Timeouts.Create == nil || res.Timeouts.Delete == nil {
		t.Fatal("expected importer/create-delete timeouts without an undocumented update operation")
	}
	scene := res.Schema["topic_scene"]
	if scene.Type != schema.TypeSet || scene.MinItems != 1 || scene.MaxItems != 5 {
		t.Fatal("scene schema must require a nonempty unordered set of at most five events")
	}
	for _, value := range []string{"SCALING_UP", "SCALING_UP_FAIL", "SCALING_DOWN", "SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL"} {
		if _, errs := scene.Elem.(*schema.Schema).ValidateFunc(value, "topic_scene"); len(errs) != 0 {
			t.Errorf("documented event %s rejected: %v", value, errs)
		}
	}
	for _, value := range []string{"", "scaling_up", "OTHER_EVENT"} {
		if _, errs := scene.Elem.(*schema.Schema).ValidateFunc(value, "topic_scene"); len(errs) == 0 {
			t.Errorf("invalid event %q accepted", value)
		}
	}
	for _, value := range []string{"", " ", "group/other", "group?query", "group#fragment"} {
		if _, errs := res.Schema["scaling_group_id"].ValidateFunc(value, "scaling_group_id"); len(errs) == 0 {
			t.Errorf("invalid identifier %q accepted", value)
		}
	}
}

func TestASNotificationV1CreateReadDelete(t *testing.T) {
	var puts, gets, deletes atomic.Int32
	ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			checkASNotificationRequest(t, r, http.MethodPut, "/scaling_notification/group")
			puts.Add(1)
			var body notifications.CreateOpts
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.TopicURN != asNotificationTestURN || len(body.TopicScene) != 2 ||
				!asNotificationTestScenes("SCALING_DOWN", "SCALING_UP").
					Equal(schema.NewSet(schema.HashString, []interface{}{body.TopicScene[0], body.TopicScene[1]})) {
				t.Errorf("wrong create request: %#v", body)
			}
			_, _ = io.WriteString(w, asNotificationTestJSON)
		case http.MethodGet:
			checkASNotificationRequest(t, r, http.MethodGet, "/scaling_notification/group")
			count := gets.Add(1)
			switch {
			case count == 1 || deletes.Load() != 0:
				_, _ = io.WriteString(w, `{"topics":[]}`)
			case count == 2:
				_, _ = io.WriteString(w, `{"topics":[{"topic_urn":"`+asNotificationTestURN+`","topic_scene":["SCALING_UP"]}]}`)
			default:
				_, _ = io.WriteString(w, `{"topics":[{"topic_urn":"unmanaged-topic","topic_scene":["SCALING_UP"],"topic_name":"unmanaged"},`+asNotificationTestJSON+`]}`)
			}
		case http.MethodDelete:
			checkASNotificationRequest(t, r, http.MethodDelete, "/scaling_notification/group/"+asNotificationTestURN)
			deletes.Add(1)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected mutation: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	res := ResourceASNotificationV1()
	d := asNotificationTestData(t)
	config := &cfg.Config{Region: "eu-de"}
	if diags := res.CreateContext(ctx, d, config); diags.HasError() {
		t.Fatal(diags)
	}
	if d.Id() != "group/"+asNotificationTestURN || d.Get("topic_name") != "test-topic" ||
		d.Get("region") != "eu-de" || d.Get("topic_scene").(*schema.Set).Len() != 2 {
		t.Fatalf("incorrect resource state: %#v", d.State())
	}
	if puts.Load() != 1 || gets.Load() < 4 {
		t.Fatalf("create did not wait for eventual visibility: PUT=%d GET=%d", puts.Load(), gets.Load())
	}
	if diags := res.DeleteContext(ctx, d, config); diags.HasError() || d.Id() != "" {
		t.Fatalf("delete failed: %v, id=%s", diags, d.Id())
	}
	if deletes.Load() != 1 {
		t.Fatalf("deleted %d associations, expected only the managed topic", deletes.Load())
	}
}

func TestASNotificationV1Import(t *testing.T) {
	res := ResourceASNotificationV1()
	var requests atomic.Int32
	ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		checkASNotificationRequest(t, r, http.MethodGet, "/scaling_notification/group")
		requests.Add(1)
		_, _ = io.WriteString(w, `{"topics":[`+asNotificationTestJSON+`]}`)
	})
	d := res.Data(nil)
	d.SetId("group/" + asNotificationTestURN)
	config := &cfg.Config{Region: "eu-de"}
	imported, err := res.Importer.StateContext(ctx, d, config)
	if err != nil || len(imported) != 1 || imported[0] != d {
		t.Fatalf("import failed: %v, %#v", err, imported)
	}
	if requests.Load() != 0 || d.Get("scaling_group_id") != "group" || d.Get("topic_urn") != asNotificationTestURN {
		t.Fatal("composite import failed to restore identity without mutating requests")
	}
	if diags := res.ReadContext(ctx, imported[0], config); diags.HasError() {
		t.Fatal(diags)
	}
	if requests.Load() != 1 || d.Get("topic_name") != "test-topic" || d.Get("topic_scene").(*schema.Set).Len() != 2 {
		t.Fatal("import refresh did not restore all SDK response fields")
	}
	for _, id := range []string{"", "group", "/topic", "group/", "group/topic/extra", " group/topic", "group/topic?x"} {
		d := res.Data(nil)
		d.SetId(id)
		if _, err := res.Importer.StateContext(ctx, d, config); err == nil {
			t.Errorf("invalid import %q accepted", id)
		}
	}
}

func TestASNotificationV1ReadErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		absent bool
	}{
		{"empty_group", http.StatusOK, `{"topics":[]}`, true},
		{"other_topic_only", http.StatusOK, `{"topics":[{"topic_urn":"other"}]}`, true},
		{"missing_group", http.StatusNotFound, `{"error_msg":"group missing"}`, true},
		{"missing_group_code", http.StatusBadRequest, `{"error":{"code":"AS.2007","message":"group missing"}}`, true},
		{"missing_group_flat_code", http.StatusBadRequest, `{"error_code":"AS.2007"}`, true},
		{"different_bad_request", http.StatusBadRequest, `{"error":{"code":"AS.2002"}}`, false},
		{"malformed_bad_request", http.StatusBadRequest, `{"error":`, false},
		{"wrong_status_missing_code", http.StatusForbidden, `{"error":{"code":"AS.2007"}}`, false},
		{"forbidden", http.StatusForbidden, `{"error_msg":"forbidden"}`, false},
		{"server_error", http.StatusInternalServerError, `{"error_msg":"server error"}`, false},
		{"malformed_json", http.StatusOK, `{"topics":`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				checkASNotificationRequest(t, r, http.MethodGet, "/scaling_notification/group")
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			})
			d := asNotificationTestData(t)
			d.SetId("group/" + asNotificationTestURN)
			diags := ResourceASNotificationV1().ReadContext(ctx, d, &cfg.Config{Region: "eu-de"})
			if tc.absent {
				if diags.HasError() || d.Id() != "" {
					t.Fatalf("missing notification must clear ID: %v, %s", diags, d.Id())
				}
			} else if !diags.HasError() || d.Id() == "" || !strings.Contains(fmt.Sprint(diags), "error reading AS notification") {
				t.Fatalf("read failure must preserve ID and report error: %v, %s", diags, d.Id())
			}
			if d.Get("topic_urn") != asNotificationTestURN || d.Get("topic_scene").(*schema.Set).Len() != 2 {
				t.Fatal("read error overwrote configured attributes")
			}
		})
	}
}

func TestASNotificationV1MutationErrors(t *testing.T) {
	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		for _, status := range []int{http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound, http.StatusConflict} {
			t.Run(fmt.Sprintf("%s/%d", method, status), func(t *testing.T) {
				var requests atomic.Int32
				ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
					path := "/scaling_notification/group"
					if method == http.MethodDelete {
						path += "/" + asNotificationTestURN
					}
					checkASNotificationRequest(t, r, method, path)
					requests.Add(1)
					w.WriteHeader(status)
					_, _ = io.WriteString(w, `{"error_msg":"mutation rejected"}`)
				})
				res := ResourceASNotificationV1()
				d := asNotificationTestData(t)
				var run func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics = res.CreateContext
				if method == http.MethodDelete {
					run = res.DeleteContext
					d.SetId("group/" + asNotificationTestURN)
				}
				diags := run(ctx, d, &cfg.Config{Region: "eu-de"})
				if method == http.MethodDelete && status == http.StatusNotFound {
					if diags.HasError() || d.Id() != "" {
						t.Fatalf("already deleted notification: %v, %s", diags, d.Id())
					}
				} else {
					if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), "mutation rejected") {
						t.Fatalf("missing mutation error: %v", diags)
					}
					if (method == http.MethodPut && d.Id() != "") || (method == http.MethodDelete && d.Id() == "") {
						t.Fatal("failed mutation changed resource identity")
					}
				}
				if requests.Load() != 1 {
					t.Fatalf("failed mutation must not continue with refresh: %d", requests.Load())
				}
			})
		}
	}
}

func TestASNotificationV1WaitErrors(t *testing.T) {
	for _, deleting := range []bool{false, true} {
		t.Run(fmt.Sprintf("deleting=%t", deleting), func(t *testing.T) {
			ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					w.WriteHeader(http.StatusForbidden)
					_, _ = io.WriteString(w, `{"error_msg":"poll rejected"}`)
				case http.MethodDelete:
					w.WriteHeader(http.StatusNoContent)
				default:
					_, _ = io.WriteString(w, asNotificationTestJSON)
				}
			})
			d := asNotificationTestData(t)
			res := ResourceASNotificationV1()
			var run func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics = res.CreateContext
			if deleting {
				d.SetId("group/" + asNotificationTestURN)
				run = res.DeleteContext
			}
			diags := run(ctx, d, &cfg.Config{Region: "eu-de"})
			if !diags.HasError() || d.Id() != "group/"+asNotificationTestURN ||
				!strings.Contains(fmt.Sprint(diags), "poll rejected") {
				t.Fatalf("polling failure must surface and retain identity: %v, %s", diags, d.Id())
			}
		})
	}
}

func TestASNotificationV1DeleteWaitsForAbsence(t *testing.T) {
	var gets atomic.Int32
	ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			checkASNotificationRequest(t, r, http.MethodDelete, "/scaling_notification/group/"+asNotificationTestURN)
			w.WriteHeader(http.StatusNoContent)
		} else {
			checkASNotificationRequest(t, r, http.MethodGet, "/scaling_notification/group")
			if gets.Add(1) == 1 {
				_, _ = io.WriteString(w, `{"topics":[`+asNotificationTestJSON+`]}`)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	})
	d := asNotificationTestData(t)
	d.SetId("group/" + asNotificationTestURN)
	if diags := ResourceASNotificationV1().DeleteContext(ctx, d, &cfg.Config{}); diags.HasError() || d.Id() != "" {
		t.Fatalf("delete wait failed: %v, %s", diags, d.Id())
	}
	if gets.Load() != 2 {
		t.Fatalf("delete did not wait until group disappeared: %d", gets.Load())
	}
}

func TestASNotificationsV1DataSource(t *testing.T) {
	var empty atomic.Bool
	ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		checkASNotificationRequest(t, r, http.MethodGet, "/scaling_notification/group")
		if empty.Load() {
			_, _ = io.WriteString(w, `{"topics":[]}`)
		} else {
			_, _ = io.WriteString(w, `{"topics":[`+asNotificationTestJSON+`,{"topic_urn":"aaa-topic","topic_name":"first","topic_scene":["SCALING_GROUP_ABNORMAL","SCALING_UP_FAIL","SCALING_DOWN_FAIL"]}]}`)
		}
	})
	ds := DataSourceASNotificationsV1()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"scaling_group_id": "group"})
	config := &cfg.Config{Region: "eu-de"}
	for i := 0; i < 2; i++ {
		if diags := ds.ReadContext(ctx, d, config); diags.HasError() {
			t.Fatal(diags)
		}
		if d.Id() != "group" || d.Get("region") != "eu-de" {
			t.Fatal("data source ID must remain stable for the group")
		}
		items := d.Get("notifications").([]interface{})
		if len(items) != 2 || items[0].(map[string]interface{})["topic_urn"] != "aaa-topic" {
			t.Fatalf("notifications must retain all topics sorted by URN: %#v", items)
		}
		item := items[1].(map[string]interface{})
		scenes := item["topic_scene"].(*schema.Set)
		if item["topic_name"] != "test-topic" || scenes.Len() != 2 ||
			!scenes.Contains("SCALING_UP") || !scenes.Contains("SCALING_DOWN") {
			t.Fatalf("incomplete notification flattening: %#v", item)
		}
		firstScenes := items[0].(map[string]interface{})["topic_scene"].(*schema.Set)
		if firstScenes.Len() != 3 || !firstScenes.Contains("SCALING_UP_FAIL") ||
			!firstScenes.Contains("SCALING_DOWN_FAIL") || !firstScenes.Contains("SCALING_GROUP_ABNORMAL") {
			t.Fatalf("failure events missing from flattened state: %#v", items[0])
		}
	}
	empty.Store(true)
	if diags := ds.ReadContext(ctx, d, config); diags.HasError() ||
		d.Id() != "group" || len(d.Get("notifications").([]interface{})) != 0 {
		t.Fatalf("empty list must clear stale data without losing group ID: %v, %#v", diags, d.State())
	}
}

func asNotificationTestScenes(scenes ...string) *schema.Set {
	values := make([]interface{}, len(scenes))
	for i, scene := range scenes {
		values[i] = scene
	}
	return schema.NewSet(schema.HashString, values)
}

func TestASNotificationsV1DataSourceErrors(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusForbidden, http.StatusInternalServerError, http.StatusOK} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				checkASNotificationRequest(t, r, http.MethodGet, "/scaling_notification/group")
				w.WriteHeader(status)
				_, _ = io.WriteString(w, `{"topics":`)
			})
			ds := DataSourceASNotificationsV1()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
				"scaling_group_id": "group",
				"notifications": []interface{}{map[string]interface{}{
					"topic_urn": asNotificationTestURN, "topic_name": "previous",
				}},
			})
			d.SetId("existing-id")
			if diags := ds.ReadContext(ctx, d, &cfg.Config{}); !diags.HasError() ||
				d.Id() != "existing-id" || len(d.Get("notifications").([]interface{})) != 1 {
				t.Fatalf("list failure must preserve state and report diagnostics: %v, %#v", diags, d.State())
			}
		})
	}
}

func TestASNotificationV1DeletedParentGroup(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusForbidden} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				checkASNotificationRequest(t, r, http.MethodDelete, "/scaling_notification/group/"+asNotificationTestURN)
				w.WriteHeader(status)
				_, _ = io.WriteString(w, `{"error":{"code":"AS.2007","message":"group missing"}}`)
			})
			d := asNotificationTestData(t)
			d.SetId("group/" + asNotificationTestURN)
			diags := ResourceASNotificationV1().DeleteContext(ctx, d, &cfg.Config{})
			if status == http.StatusBadRequest {
				if diags.HasError() || d.Id() != "" {
					t.Fatalf("deleted parent group should permit removal: %v, %s", diags, d.Id())
				}
			} else if !diags.HasError() || d.Id() == "" {
				t.Fatalf("wrong status must not hide an error: %v, %s", diags, d.Id())
			}
		})
	}
}

func TestASNotificationsV1ConfiguredClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkASNotificationRequest(t, r, http.MethodGet, "/autoscaling-api/v1/project/scaling_notification/group")
		_, _ = io.WriteString(w, `{"topics":[`+asNotificationTestJSON+`]}`)
	}))
	defer server.Close()
	config := &cfg.Config{
		Region: "eu-de",
		HwClient: &golangsdk.ProviderClient{
			TokenID: "test-token",
			EndpointLocator: func(opts golangsdk.EndpointOpts) (string, error) {
				if opts.Type != "asv1" || opts.Region != "eu-de" {
					t.Errorf("unexpected endpoint lookup: %#v", opts)
				}
				return server.URL + "/autoscaling-api/v1/project/", nil
			},
		},
	}
	ds := DataSourceASNotificationsV1()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"scaling_group_id": "group"})
	if diags := ds.ReadContext(context.Background(), d, config); diags.HasError() || d.Id() != "group" {
		t.Fatalf("configured provider client failed: %v", diags)
	}
}

func TestASNotificationV1WaitCancellation(t *testing.T) {
	for _, deleting := range []bool{false, true} {
		t.Run(fmt.Sprint(deleting), func(t *testing.T) {
			parent, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctx := asNotificationTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					if deleting {
						_, _ = io.WriteString(w, `{"topics":[`+asNotificationTestJSON+`]}`)
					} else {
						_, _ = io.WriteString(w, `{"topics":[]}`)
					}
					return
				}
				cancel()
				if deleting {
					w.WriteHeader(http.StatusNoContent)
				} else {
					_, _ = io.WriteString(w, asNotificationTestJSON)
				}
			})
			clientCtx := ctx.(asNotificationTestClientContext)
			clientCtx.Context = parent
			res := ResourceASNotificationV1()
			d := asNotificationTestData(t)
			var run func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics = res.CreateContext
			if deleting {
				d.SetId("group/" + asNotificationTestURN)
				run = res.DeleteContext
			}
			diags := run(clientCtx, d, &cfg.Config{})
			if !diags.HasError() || d.Id() != "group/"+asNotificationTestURN {
				t.Fatalf("cancelled wait must report failure and retain identity: %v, %s", diags, d.Id())
			}
		})
	}
}
