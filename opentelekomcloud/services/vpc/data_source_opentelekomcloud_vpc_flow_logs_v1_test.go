package vpc

import (
	"reflect"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/flow_logs"
)

func TestDataSourceVpcFlowLogsV1Schema(t *testing.T) {
	dataSource := DataSourceVpcFlowLogsV1()
	for _, field := range []string{
		"id", "name", "tenant_id", "description", "resource_type", "resource_id",
		"traffic_type", "log_group_id", "log_topic_id", "status", "limit", "marker",
	} {
		if !dataSource.Schema[field].Optional {
			t.Fatalf("%s must be an optional filter", field)
		}
	}
	if !dataSource.Schema["flow_logs"].Computed {
		t.Fatal("flow_logs must be computed")
	}
}

func TestFlattenVpcFlowLog(t *testing.T) {
	flowLog := flow_logs.FlowLog{
		ID:           "flow-log-id",
		Name:         "flow-log",
		TenantID:     "tenant-id",
		Description:  "description",
		ResourceType: "vpc",
		ResourceID:   "vpc-id",
		TrafficType:  "all",
		LogGroupID:   "group-id",
		LogTopicID:   "topic-id",
		IndexEnabled: true,
		AdminState:   true,
		Status:       "ACTIVE",
		CreatedAt:    "created",
		UpdatedAt:    "updated",
	}

	expected := map[string]interface{}{
		"id":            "flow-log-id",
		"name":          "flow-log",
		"tenant_id":     "tenant-id",
		"description":   "description",
		"resource_type": "vpc",
		"resource_id":   "vpc-id",
		"traffic_type":  "all",
		"log_group_id":  "group-id",
		"log_topic_id":  "topic-id",
		"enabled":       true,
		"status":        "ACTIVE",
		"created_at":    "created",
		"updated_at":    "updated",
	}
	if actual := flattenVpcFlowLog(flowLog); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected flattened flow log: %#v", actual)
	}
}
