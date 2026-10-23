package v3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func memberTestData(t *testing.T, resource *schema.Resource, values map[string]string, config map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	d, err := schema.InternalMap(resource.Schema).Data(
		&terraform.InstanceState{Attributes: values},
		&terraform.InstanceDiff{RawConfig: cty.ObjectVal(config)},
	)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func memberTestContext(t *testing.T, handler http.HandlerFunc) context.Context {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := &golangsdk.ServiceClient{
		ProviderClient: &golangsdk.ProviderClient{},
		Endpoint:       server.URL + "/",
	}
	return context.WithValue(context.Background(), keyClient, client)
}

func TestLBMemberV3Filters(t *testing.T) {
	values := map[string]string{
		"name": "member", "address": "192.0.2.10", "operating_status": "ONLINE",
		"ip_version": "v4", "member_type": "ip", "instance_id": "instance",
		"availability_zone": "eu-de-01", "subnet_id": "subnet",
		"weight": "0", "protocol_port": "8080", "admin_state_up": "false",
	}
	config := map[string]cty.Value{}
	for key, value := range values {
		config[key] = cty.StringVal(value)
	}
	config["weight"] = cty.NumberIntVal(0)
	config["protocol_port"] = cty.NumberIntVal(8080)
	config["admin_state_up"] = cty.False
	d := memberTestData(t, DataSourceLBMemberV3(), values, config)
	opts := buildLBMemberV3ListOpts(d)
	if !reflect.DeepEqual(opts.Weight, []int{0}) || opts.AdminStateUp == nil || *opts.AdminStateUp {
		t.Fatalf("zero/false filters lost: %#v", opts)
	}
	for key, actual := range map[string][]string{
		"name": opts.Name, "address": opts.Address, "operating_status": opts.OperatingStatus,
		"ip_version": opts.IPVersion, "member_type": opts.MemberType, "instance_id": opts.InstanceID,
		"availability_zone": opts.AvailabilityZone, "subnet_id": opts.SubnetCIDRID,
	} {
		if !reflect.DeepEqual(actual, []string{values[key]}) {
			t.Errorf("%s: got %#v", key, actual)
		}
	}
	if !reflect.DeepEqual(opts.ProtocolPort, []int{8080}) || opts.SubnetID != nil || opts.EnterpriseProjectID != nil {
		t.Fatalf("unexpected options: %#v", opts)
	}
	unconfigured := memberTestData(t, DataSourceLBMemberV3(), values, map[string]cty.Value{"pool_id": cty.StringVal("pool")})
	if actual := buildLBMemberV3ListOpts(unconfigured); !reflect.DeepEqual(actual, members.ListOpts{}) {
		t.Fatalf("computed fields must not become filters: %#v", actual)
	}
}

func TestLBMemberV3FlattenAndClear(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceLBMemberV3().Schema, nil)
	member := &members.Member{
		ID: "member", Address: "192.0.2.10", ProtocolPort: 80, Weight: 0,
		AvailabilityZone: "eu-de-01", MemberType: "ip", InstanceID: "instance",
		CreatedAt: "created", UpdatedAt: "updated", AdminStateUp: false,
		Status: []members.MemberStatus{{
			ListenerID: "listener", OperatingStatus: "OFFLINE", CreatedAt: "created", UpdatedAt: "updated",
			Reason: &members.MemberHealthCheckFailedReason{ReasonCode: "RESPONSE_MISMATCH", ExpectedResponse: "200", HealthcheckResponse: "500"},
		}},
	}
	if diagnostics := setLBMemberV3Fields(d, member); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	for key, expected := range map[string]interface{}{
		"member_id": "member", "weight": 0, "admin_state_up": false, "availability_zone": "eu-de-01",
		"member_type": "ip", "instance_id": "instance", "created_at": "created", "updated_at": "updated",
		"status.0.reason.0.reason_code":       "RESPONSE_MISMATCH",
		"status.0.reason.0.expected_response": "200", "status.0.reason.0.healthcheck_response": "500",
	} {
		if actual := d.Get(key); actual != expected {
			t.Errorf("%s: got %#v, want %#v", key, actual, expected)
		}
	}
	member.Status[0].Reason = nil
	if diagnostics := setLBMemberV3Fields(d, member); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if d.Get("status.0.reason.#") != 0 {
		t.Fatal("stale health-check reason retained")
	}
	member.Status = nil
	if diagnostics := setLBMemberV3Fields(d, member); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if d.Get("status.#") != 0 {
		t.Fatal("stale status retained")
	}
}

func TestLBMemberV3DataSourceHTTP(t *testing.T) {
	for _, test := range []struct {
		name, body string
		count      int
		wantError  bool
	}{
		{"one", `{"members":[{"id":"member","weight":0,"admin_state_up":false}],"page_info":{}}`, 2, false},
		{"none", `{"members":[],"page_info":{}}`, 2, true},
		{"multiple", `{"members":[{"id":"a"},{"id":"b"}],"page_info":{}}`, 2, true},
		{"failure", `{"error":"unavailable"}`, 1, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			ctx := memberTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != "/pools/pool/members" || r.Method != http.MethodGet {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.URL.Query().Get("weight") != "0" || r.URL.Query().Get("admin_state_up") != "false" {
					t.Errorf("zero/false filter lost: %s", r.URL)
				}
				if test.name == "failure" {
					w.WriteHeader(http.StatusBadRequest)
				}
				_, _ = w.Write([]byte(test.body))
			})
			d := memberTestData(t, DataSourceLBMemberV3(),
				map[string]string{"pool_id": "pool", "weight": "0", "admin_state_up": "false", "member_id": "stale"},
				map[string]cty.Value{"pool_id": cty.StringVal("pool"), "weight": cty.NumberIntVal(0), "admin_state_up": cty.False})
			diagnostics := dataSourceLBMemberV3Read(ctx, d, &cfg.Config{})
			if diagnostics.HasError() != test.wantError || calls != test.count {
				t.Fatalf("unexpected result: calls=%d, diagnostics=%#v", calls, diagnostics)
			}
			if !test.wantError && d.Id() != "pool/member" {
				t.Fatalf("unexpected ID %q", d.Id())
			}
		})
	}
}

func TestLBMemberV3LifecycleHTTP(t *testing.T) {
	calls := 0
	ctx := memberTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.Method {
		case http.MethodPost:
			var body struct {
				Member members.CreateOpts `json:"member"`
			}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.Member.AdminStateUp == nil || *body.Member.AdminStateUp ||
				body.Member.Weight == nil || *body.Member.Weight != 0 || body.Member.AvailabilityZone != "eu-de-01" {
				t.Errorf("unexpected create body: %#v", body)
			}
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			var body struct {
				Member members.UpdateOpts `json:"member"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.Member.Weight == nil || *body.Member.Weight != 0 || body.Member.Name == nil ||
				*body.Member.Name != "" || body.Member.AvailabilityZone != "eu-de-02" || body.Member.AdminStateUp != nil {
				t.Errorf("unexpected update body: %#v", body)
			}
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_, _ = w.Write([]byte(`{"member":{"id":"member","address":"192.0.2.10","protocol_port":80,"weight":0,"availability_zone":"eu-de-01"}}`))
	})
	d := memberTestData(t, ResourceLBMemberV3(), map[string]string{
		"pool_id": "pool", "address": "192.0.2.10", "protocol_port": "80", "weight": "0",
		"availability_zone": "eu-de-01", "admin_state_up": "false",
	}, map[string]cty.Value{"admin_state_up": cty.False, "weight": cty.NumberIntVal(0)})
	if diagnostics := resourceLBMemberV3Create(ctx, d, &cfg.Config{}); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if d.Id() != "pool/member" {
		t.Fatalf("unexpected ID %q", d.Id())
	}
	oldState := d.State()
	oldState.Attributes["name"] = "old"
	oldState.Attributes["weight"] = "1"
	updated, err := schema.InternalMap(ResourceLBMemberV3().Schema).Data(oldState, &terraform.InstanceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"name":              {Old: "old", New: ""},
			"weight":            {Old: "1", New: "0"},
			"availability_zone": {Old: "eu-de-01", New: "eu-de-02"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if diagnostics := resourceLBMemberV3Update(ctx, updated, &cfg.Config{}); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if diagnostics := resourceLBMemberV3Delete(ctx, d, &cfg.Config{}); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if calls != 5 {
		t.Fatalf("expected create/read/update/read/delete, got %d calls", calls)
	}
}

func TestLBMemberV3GetAndNotFound(t *testing.T) {
	for _, test := range []struct {
		name      string
		resource  bool
		code      int
		wantError bool
	}{
		{"data_get", false, http.StatusOK, false},
		{"data_not_found", false, http.StatusNotFound, true},
		{"resource_not_found", true, http.StatusNotFound, false},
		{"resource_error", true, http.StatusForbidden, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := memberTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/pools/pool/members/member" {
					t.Errorf("unexpected path %s", r.URL.Path)
				}
				w.WriteHeader(test.code)
				_, _ = w.Write([]byte(`{"member":{"id":"member"}}`))
			})
			resource := DataSourceLBMemberV3()
			if test.resource {
				resource = ResourceLBMemberV3()
			}
			d := memberTestData(t, resource,
				map[string]string{"pool_id": "pool", "member_id": "member"},
				map[string]cty.Value{"member_id": cty.StringVal("member")})
			d.SetId("pool/member")
			read := dataSourceLBMemberV3Read
			if test.resource {
				read = resourceLBMemberV3Read
			}
			diagnostics := read(ctx, d, &cfg.Config{})
			if diagnostics.HasError() != test.wantError {
				t.Fatalf("unexpected diagnostics: %#v", diagnostics)
			}
			if test.resource && test.code == http.StatusNotFound && d.Id() != "" {
				t.Fatal("missing resource was not removed from state")
			}
		})
	}
}

func TestLBMemberV3CompatibilityAndImport(t *testing.T) {
	resource := ResourceLBMemberV3()
	for _, key := range []string{
		"pool_id", "address", "name", "project_id", "protocol_port", "subnet_id",
		"weight", "member_id", "operating_status", "ip_version",
	} {
		if resource.Schema[key] == nil {
			t.Errorf("existing attribute %s removed", key)
		}
	}
	if !resource.Schema["admin_state_up"].ForceNew || resource.Schema["availability_zone"].ForceNew {
		t.Fatal("incorrect mutability for new attributes")
	}
	for _, key := range []string{"device_owner", "device_id", "loadbalancer_id", "loadbalancers", "enterprise_project_id"} {
		if resource.Schema[key] != nil || DataSourceLBMemberV3().Schema[key] != nil {
			t.Errorf("unsupported field %s exposed", key)
		}
	}
	d := resource.Data(nil)
	d.SetId("pool/member")
	imported, err := resource.Importer.StateContext(context.Background(), d, nil)
	if err != nil || len(imported) != 1 || poolID(d) != "pool" || memberID(d) != "member" {
		t.Fatalf("composite import changed: %#v, %v", imported, err)
	}
	for _, code := range []int{http.StatusNotFound, http.StatusForbidden} {
		ctx := memberTestContext(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
		})
		if diagnostics := resourceLBMemberV3Delete(ctx, d, &cfg.Config{}); diagnostics.HasError() != (code == http.StatusForbidden) {
			t.Fatalf("unexpected delete diagnostics: %#v", diagnostics)
		}
	}
}
