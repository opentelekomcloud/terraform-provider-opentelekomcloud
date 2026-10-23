package v3

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
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
