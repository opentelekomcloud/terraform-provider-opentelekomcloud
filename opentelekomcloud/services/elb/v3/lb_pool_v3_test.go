package v3

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/structs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/pools"
)

func TestLBPoolV3SchemaCoverage(t *testing.T) {
	resourceSchema := ResourceLBPoolV3().Schema
	for _, key := range []string{
		"slow_start", "healthmonitor_id", "listener_ids", "loadbalancer_ids",
		"member_ids", "protection_status", "protection_reason", "created_at", "updated_at",
	} {
		if _, ok := resourceSchema[key]; !ok {
			t.Fatalf("resource schema is missing %q", key)
		}
	}

	dataSourceSchema := DataSourceLBPoolV3().Schema
	for _, key := range []string{
		"id", "healthmonitor_id", "member_address", "member_device_id",
		"member_instance_id", "protection_status",
	} {
		if _, ok := dataSourceSchema[key]; !ok {
			t.Fatalf("data source schema is missing %q", key)
		}
	}
	for _, unsupported := range []string{"admin_state_up", "enterprise_project_id", "az_affinity"} {
		if _, ok := dataSourceSchema[unsupported]; ok {
			t.Fatalf("unsupported field %q must not be exposed", unsupported)
		}
	}
}

func TestBuildLBPoolV3ListOpts(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSourceLBPoolV3().Schema, map[string]interface{}{
		"description":                "pool description",
		"healthmonitor_id":           "monitor-id",
		"lb_algorithm":               "ROUND_ROBIN",
		"protocol":                   "HTTP",
		"name":                       "pool-name",
		"loadbalancer_id":            "loadbalancer-id",
		"ip_version":                 "v4",
		"member_address":             "192.0.2.10",
		"member_device_id":           "device-id",
		"member_deletion_protection": false,
		"listener_id":                "listener-id",
		"member_instance_id":         "instance-id",
		"vpc_id":                     "vpc-id",
		"type":                       "instance",
		"protection_status":          "nonProtection",
	})

	memberDeletionProtection := false
	opts := buildLBPoolV3ListOpts(d, &memberDeletionProtection)
	if !reflect.DeepEqual(opts.Description, []string{"pool description"}) ||
		!reflect.DeepEqual(opts.HealthMonitorID, []string{"monitor-id"}) ||
		!reflect.DeepEqual(opts.LBMethod, []string{"ROUND_ROBIN"}) ||
		!reflect.DeepEqual(opts.Protocol, []string{"HTTP"}) ||
		!reflect.DeepEqual(opts.ListenerID, []string{"listener-id"}) ||
		!reflect.DeepEqual(opts.MemberInstanceID, []string{"instance-id"}) {
		t.Fatalf("unexpected list options: %#v", opts)
	}
	if opts.MemberDeletionProtectionEnable == nil || *opts.MemberDeletionProtectionEnable {
		t.Fatalf("false member deletion protection filter was not preserved: %#v", opts)
	}
}

func TestSetLBPoolV3Fields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSourceLBPoolV3().Schema, nil)
	pool := &pools.Pool{
		ID:            "pool-id",
		Name:          "pool-name",
		LBMethod:      "ROUND_ROBIN",
		Protocol:      "HTTP",
		Loadbalancers: []structs.ResourceRef{{ID: "loadbalancer-id"}},
		Listeners:     []structs.ResourceRef{{ID: "listener-id"}},
		Members:       []structs.ResourceRef{{ID: "member-id"}},
		SlowStart:     &pools.SlowStart{Enable: true, Duration: 30},
		CreatedAt:     "2026-09-22T09:00:00Z",
		UpdatedAt:     "2026-09-22T09:01:00Z",
	}

	if diags := setLBPoolV3Fields(d, pool); diags.HasError() {
		t.Fatalf("failed to flatten pool: %#v", diags)
	}
	if d.Id() != "pool-id" ||
		d.Get("slow_start.0.duration").(int) != 30 ||
		d.Get("loadbalancer_ids.#").(int) != 1 ||
		d.Get("listener_ids.#").(int) != 1 ||
		d.Get("member_ids.#").(int) != 1 {
		t.Fatalf("unexpected flattened state: %#v", d.State())
	}
}
