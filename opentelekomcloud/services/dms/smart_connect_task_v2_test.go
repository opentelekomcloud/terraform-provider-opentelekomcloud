package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/smart_connect"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestSmartConnectTaskV2Schema(t *testing.T) {
	r := ResourceDmsSmartConnectTaskV2()
	p := &schema.Provider{ResourcesMap: map[string]*schema.Resource{"opentelekomcloud_dms_smart_connect_task_v2": r}}
	if err := p.InternalValidate(); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{
		"instance_id", "task_name", "topics", "topics_regex", "start_later", "source_type",
		"destination_type", "source_task", "destination_task",
	} {
		if r.Schema[field] == nil || !r.Schema[field].ForceNew {
			t.Errorf("existing immutable attribute lost: %s", field)
		}
	}
	aliases := map[string]string{
		"current_cluster_name": "current_instance_alias",
		"cluster_name":         "peer_instance_alias",
		"instance_id":          "peer_instance_id",
		"bootstrap_servers":    "peer_instance_address",
	}
	for _, block := range []struct {
		name     string
		request  reflect.Type
		response reflect.Type
	}{
		{"source_task", reflect.TypeOf(smart_connect.SmartConnectTaskSourceConfig{}), reflect.TypeOf(smart_connect.SmartConnectTaskRespSourceConfig{})},
		{"destination_task", reflect.TypeOf(smart_connect.SmartConnectTaskSinkConfig{}), reflect.TypeOf(smart_connect.SmartConnectTaskRespSinkConfig{})},
	} {
		fields := r.Schema[block.name].Elem.(*schema.Resource).Schema
		for i := 0; i < block.request.NumField(); i++ {
			name := strings.Split(block.request.Field(i).Tag.Get("json"), ",")[0]
			if alias, ok := aliases[name]; ok && block.name == "source_task" {
				name = alias
			}
			if fields[name] == nil || !fields[name].Optional || !fields[name].ForceNew {
				t.Errorf("missing immutable request field %s.%s", block.name, name)
			}
		}
		for i := 0; i < block.response.NumField(); i++ {
			name := block.response.Field(i).Tag.Get("json")
			if alias, ok := aliases[name]; ok && block.name == "source_task" {
				name = alias
			}
			if fields[name] == nil {
				t.Errorf("missing response field %s.%s", block.name, name)
			}
		}
		if !fields["redis_password"].Sensitive || fields["redis_password"].Computed {
			t.Errorf("%s Redis password must be sensitive and configuration-only", block.name)
		}
	}
	for _, name := range []string{"target_db", "obs_part_size", "flush_size", "timezone", "schema_generator_class", "partitioner_class", "value_converter", "key_converter", "kv_delimiter"} {
		field := r.Schema["destination_task"].Elem.(*schema.Resource).Schema[name]
		if !field.Computed || field.Optional {
			t.Errorf("response-only destination field %s is configurable", name)
		}
	}
}

type smartConnectTestClientContext struct {
	context.Context
	client *golangsdk.ServiceClient
}

func (c smartConnectTestClientContext) Value(key interface{}) interface{} {
	if key == dmsClientV2 {
		return c.client
	}
	return c.Context.Value(key)
}

func smartConnectTestContext(t *testing.T, handler http.HandlerFunc) context.Context {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := &golangsdk.ServiceClient{
		ProviderClient: &golangsdk.ProviderClient{},
		Endpoint:       server.URL + "/v2/project/",
	}
	return smartConnectTestClientContext{Context: context.Background(), client: client}
}

func TestSmartConnectTaskV2Lifecycle(t *testing.T) {
	for _, test := range []struct {
		name     string
		config   map[string]interface{}
		request  string
		response string
		check    map[string]interface{}
	}{
		{
			name: "redis",
			config: map[string]interface{}{
				"source_type": "REDIS_REPLICATOR_SOURCE", "destination_type": "REDIS_REPLICATOR_SINK",
				"source_task": []interface{}{map[string]interface{}{
					"redis_address": "192.0.2.1:6379,192.0.2.2:6379", "redis_type": "cluster",
					"dcs_instance_id": "source-dcs", "redis_password": "source-test-password",
					"sync_mode": "CUSTOM_OFFSET", "full_sync_wait_ms": 10000, "full_sync_max_retry": 10, "ratelimit": -1,
				}},
				"destination_task": []interface{}{map[string]interface{}{
					"redis_address": "192.0.2.3:6379", "redis_type": "standalone",
					"dcs_instance_id": "sink-dcs", "redis_password": "sink-test-password",
				}},
			},
			request: `{
				"task_name":"test", "topics":"topic", "start_later":false,
				"source_type":"REDIS_REPLICATOR_SOURCE", "sink_type":"REDIS_REPLICATOR_SINK",
				"source_task":{"redis_address":"192.0.2.1:6379,192.0.2.2:6379","redis_type":"cluster",
					"dcs_instance_id":"source-dcs","redis_password":"source-test-password",
					"sync_mode":"CUSTOM_OFFSET","full_sync_wait_ms":10000,"full_sync_max_retry":10,"ratelimit":-1},
				"sink_task":{"redis_address":"192.0.2.3:6379","redis_type":"standalone",
					"dcs_instance_id":"sink-dcs","redis_password":"sink-test-password"}
			}`,
			response: `{
				"source_type":"REDIS_REPLICATOR_SOURCE", "sink_type":"REDIS_REPLICATOR_SINK",
				"source_task":{"redis_address":"192.0.2.1:6379,192.0.2.2:6379","redis_type":"cluster",
					"dcs_instance_id":"source-dcs","sync_mode":"CUSTOM_OFFSET",
					"full_sync_wait_ms":10000,"full_sync_max_retry":10,"ratelimit":-1},
				"sink_task":{"redis_address":"192.0.2.3:6379","redis_type":"standalone","dcs_instance_id":"sink-dcs","target_db":0}
			}`,
			check: map[string]interface{}{
				"source_task.0.redis_address": "192.0.2.1:6379,192.0.2.2:6379",
				"source_task.0.redis_type":    "cluster", "source_task.0.dcs_instance_id": "source-dcs",
				"source_task.0.redis_password": "source-test-password", "source_task.0.sync_mode": "CUSTOM_OFFSET",
				"source_task.0.full_sync_wait_ms": 10000, "source_task.0.full_sync_max_retry": 10, "source_task.0.ratelimit": -1,
				"source_task.0.peer_instance_address.#": 0,
				"destination_task.0.redis_address":      "192.0.2.3:6379", "destination_task.0.redis_type": "standalone",
				"destination_task.0.dcs_instance_id": "sink-dcs", "destination_task.0.redis_password": "sink-test-password",
				"destination_task.0.target_db": 0,
			},
		},
		{
			name: "kafka",
			config: map[string]interface{}{
				"source_type": "KAFKA_REPLICATOR_SOURCE",
				"source_task": []interface{}{map[string]interface{}{
					"current_instance_alias": "A", "peer_instance_alias": "B", "peer_instance_id": "peer",
					"peer_instance_address": []interface{}{"192.0.2.1:9092"}, "security_protocol": "SASL_SSL",
					"user_name": "test-user", "password": "kafka-test-password", "sasl_mechanism": "PLAIN",
					"direction": "push", "replication_factor": 3, "task_num": 2, "consumer_strategy": "latest",
					"compression_type": "snappy", "topics_mapping": []interface{}{"topic:target"},
					"sync_consumer_offsets_enabled": false, "rename_topic_enabled": false, "provenance_header_enabled": true,
				}},
			},
			request: `{
				"task_name":"test","topics":"topic","start_later":false,"source_type":"KAFKA_REPLICATOR_SOURCE",
				"source_task":{"current_cluster_name":"A","cluster_name":"B","instance_id":"peer",
					"bootstrap_servers":"192.0.2.1:9092","security_protocol":"SASL_SSL","user_name":"test-user",
					"password":"kafka-test-password","sasl_mechanism":"PLAIN","direction":"push",
					"replication_factor":3,"task_num":2,"consumer_strategy":"latest","compression_type":"snappy",
					"topics_mapping":"topic:target","sync_consumer_offsets_enabled":false,
					"rename_topic_enabled":false,"provenance_header_enabled":true}
			}`,
			response: `{
				"source_type":"KAFKA_REPLICATOR_SOURCE","sink_type":"NONE","sink_task":null,
				"source_task":{"current_cluster_name":"A","cluster_name":"B","instance_id":"peer",
					"bootstrap_servers":"192.0.2.1:9092","security_protocol":"SASL_SSL","user_name":"test-user",
					"sasl_mechanism":"PLAIN","direction":"push","replication_factor":3,"task_num":2,
					"consumer_strategy":"latest","compression_type":"snappy","topics_mapping":"topic:target",
					"sync_consumer_offsets_enabled":false,"rename_topic_enabled":false,"provenance_header_enabled":true}
			}`,
			check: map[string]interface{}{
				"source_task.0.current_instance_alias": "A", "source_task.0.peer_instance_alias": "B",
				"source_task.0.peer_instance_id": "peer", "source_task.0.peer_instance_address.#": 1,
				"source_task.0.security_protocol": "SASL_SSL", "source_task.0.user_name": "test-user",
				"source_task.0.password": "kafka-test-password", "source_task.0.sasl_mechanism": "PLAIN",
				"source_task.0.direction": "push", "source_task.0.replication_factor": 3, "source_task.0.task_num": 2,
				"source_task.0.consumer_strategy": "latest", "source_task.0.compression_type": "snappy",
				"source_task.0.topics_mapping.#": 1, "source_task.0.sync_consumer_offsets_enabled": false,
				"source_task.0.rename_topic_enabled": false, "source_task.0.provenance_header_enabled": true,
				"destination_task.#": 0,
			},
		},
		{
			name: "obs",
			config: map[string]interface{}{
				"destination_type": "OBS_SINK",
				"destination_task": []interface{}{map[string]interface{}{
					"access_key": "test-ak", "secret_key": "test-sk", "consumer_strategy": "earliest",
					"destination_file_type": "TEXT", "deliver_time_interval": 300, "obs_bucket_name": "bucket",
					"obs_path": "path", "partition_format": "yyyy/MM", "record_delimiter": ";", "store_keys": false,
				}},
			},
			request: `{
				"task_name":"test","topics":"topic","start_later":false,"sink_type":"OBS_SINK",
				"sink_task":{"access_key":"test-ak","secret_key":"test-sk","consumer_strategy":"earliest",
					"destination_file_type":"TEXT","deliver_time_interval":300,"obs_bucket_name":"bucket",
					"obs_path":"path","partition_format":"yyyy/MM","record_delimiter":";","store_keys":false}
			}`,
			response: `{
				"source_type":"NONE","source_task":null,"sink_type":"OBS_SINK",
				"sink_task":{"consumer_strategy":"earliest","destination_file_type":"TEXT","deliver_time_interval":300,
					"obs_bucket_name":"bucket","obs_path":"path","partition_format":"yyyy/MM","record_delimiter":";",
					"store_keys":false,"obs_part_size":5242880,"flush_size":1000000,"timezone":"UTC",
					"schema_generator_class":"schema","partitioner_class":"partitioner","value_converter":"value",
					"key_converter":"key","kv_delimiter":":"}
			}`,
			check: map[string]interface{}{
				"destination_task.0.access_key": "test-ak", "destination_task.0.secret_key": "test-sk",
				"destination_task.0.consumer_strategy": "earliest", "destination_task.0.destination_file_type": "TEXT",
				"destination_task.0.deliver_time_interval": 300, "destination_task.0.obs_bucket_name": "bucket",
				"destination_task.0.obs_path": "path", "destination_task.0.partition_format": "yyyy/MM",
				"destination_task.0.record_delimiter": ";", "destination_task.0.store_keys": false,
				"destination_task.0.obs_part_size": 5242880, "destination_task.0.flush_size": 1000000,
				"destination_task.0.timezone": "UTC", "destination_task.0.schema_generator_class": "schema",
				"destination_task.0.partitioner_class": "partitioner", "destination_task.0.value_converter": "value",
				"destination_task.0.key_converter": "key", "destination_task.0.kv_delimiter": ":", "source_task.#": 0,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.config["instance_id"] = "instance"
			test.config["task_name"] = "test"
			test.config["topics"] = []interface{}{"topic"}
			d := schema.TestResourceDataRaw(t, ResourceDmsSmartConnectTaskV2().Schema, test.config)
			var response map[string]interface{}
			if err := json.Unmarshal([]byte(test.response), &response); err != nil {
				t.Fatal(err)
			}
			response["id"], response["task_name"], response["topics"] = "task", "test", "topic"
			response["status"], response["create_time"] = "RUNNING", int64(1708427753133)
			var methods []string
			ctx := smartConnectTestContext(t, func(w http.ResponseWriter, r *http.Request) {
				methods = append(methods, r.Method)
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPost:
					if r.URL.Path != "/v2/project/instances/instance/connector/tasks" {
						t.Errorf("unexpected create path: %s", r.URL.Path)
					}
					var actual, expected map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&actual); err != nil {
						t.Error(err)
					}
					if err := json.Unmarshal([]byte(test.request), &expected); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(actual, expected) {
						t.Errorf("unexpected request: got %#v, want %#v", actual, expected)
					}
				case http.MethodGet, http.MethodDelete:
					if r.URL.Path != "/v2/project/instances/instance/connector/tasks/task" {
						t.Errorf("unexpected task path: %s", r.URL.Path)
					}
				default:
					t.Errorf("unexpected method: %s", r.Method)
				}
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			config := &cfg.Config{Region: "eu-de"}
			if diagnostics := resourceDmsV2SmartConnectTaskCreate(ctx, d, config); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if d.Id() != "task" {
				t.Fatalf("unexpected resource ID: %s", d.Id())
			}
			test.check["task_name"], test.check["created_at"] = "test", "1708427753133"
			test.check["region"], test.check["status"] = "eu-de", "RUNNING"
			for field, expected := range test.check {
				if actual := d.Get(field); actual != expected {
					t.Errorf("%s: got %#v, want %#v", field, actual, expected)
				}
			}
			if diagnostics := resourceDmsV2SmartConnectTaskDelete(ctx, d, config); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if !reflect.DeepEqual(methods, []string{"POST", "GET", "GET", "DELETE"}) {
				t.Errorf("unexpected lifecycle: %v", methods)
			}
		})
	}
}

func TestSmartConnectTaskV2FlattenAndClear(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceDmsSmartConnectTaskV2().Schema, nil)
	for _, db := range []int{-1, 0, 3} {
		sink := &smart_connect.SmartConnectTaskRespSinkConfig{TargetDB: &db}
		if err := d.Set("destination_task", flattenSinkTaskResponse(d, sink)); err != nil {
			t.Fatal(err)
		}
		if actual := d.Get("destination_task.0.target_db"); actual != db {
			t.Fatalf("database lost: got %v, want %d", actual, db)
		}
	}
	if err := d.Set("destination_task", flattenSinkTaskResponse(d, &smart_connect.SmartConnectTaskRespSinkConfig{RedisAddress: "192.0.2.1:6379"})); err != nil {
		t.Fatal(err)
	}
	if actual := d.Get("destination_task.0.target_db"); actual != 0 {
		t.Errorf("stale target database retained: %v", actual)
	}
	for _, source := range []*smart_connect.SmartConnectTaskRespSourceConfig{nil, {}} {
		if err := d.Set("source_task", flattenSourceTaskResponse(d, source)); err != nil {
			t.Fatal(err)
		}
		if d.Get("source_task.#") != 0 {
			t.Fatal("empty source retained")
		}
	}
	for _, sink := range []*smart_connect.SmartConnectTaskRespSinkConfig{nil, {}} {
		if err := d.Set("destination_task", flattenSinkTaskResponse(d, sink)); err != nil {
			t.Fatal(err)
		}
		if d.Get("destination_task.#") != 0 {
			t.Fatal("stale sink retained")
		}
	}
}

func TestSmartConnectTaskV2ReadClearsNullBlocks(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceDmsSmartConnectTaskV2().Schema, map[string]interface{}{
		"instance_id": "instance",
		"source_task": []interface{}{map[string]interface{}{
			"redis_address": "192.0.2.1:6379", "redis_password": "source-test-password",
		}},
		"destination_task": []interface{}{map[string]interface{}{
			"obs_bucket_name": "bucket", "secret_key": "test-sk",
		}},
	})
	d.SetId("task")
	ctx := smartConnectTestContext(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"task","source_task":null,"sink_task":null,"status":"RUNNING"}`)
	})
	if diags := resourceDmsV2SmartConnectTaskRead(ctx, d, &cfg.Config{Region: "eu-de"}); diags.HasError() {
		t.Fatal(diags)
	}
	if d.Get("source_task.#") != 0 || d.Get("destination_task.#") != 0 {
		t.Fatal("stale source/destination configuration retained")
	}
}

func TestSmartConnectTaskV2ErrorsAndImport(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodDelete, http.MethodPost} {
		for _, code := range []int{http.StatusNotFound, http.StatusForbidden} {
			t.Run(fmt.Sprintf("%s-%d", method, code), func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, ResourceDmsSmartConnectTaskV2().Schema, map[string]interface{}{"instance_id": "instance"})
				d.SetId("task")
				ctx := smartConnectTestContext(t, func(w http.ResponseWriter, r *http.Request) {
					if r.Method != method {
						t.Errorf("got %s, want %s", r.Method, method)
					}
					w.WriteHeader(code)
				})
				op := resourceDmsV2SmartConnectTaskRead
				switch method {
				case http.MethodPost:
					op = resourceDmsV2SmartConnectTaskCreate
				case http.MethodDelete:
					op = resourceDmsV2SmartConnectTaskDelete
				}
				diags := op(ctx, d, &cfg.Config{})
				deleted := code == http.StatusNotFound && method != http.MethodPost
				if diags.HasError() == deleted {
					t.Fatalf("unexpected error handling: %v", diags)
				}
				if deleted && d.Id() != "" {
					t.Fatal("deleted resource retained in state")
				}
			})
		}
	}
	for _, id := range []string{"instance/task", "invalid"} {
		d := schema.TestResourceDataRaw(t, ResourceDmsSmartConnectTaskV2().Schema, nil)
		d.SetId(id)
		_, err := resourceDmsV2SmartConnectTaskImportState(context.Background(), d, nil)
		if id == "invalid" {
			if err == nil {
				t.Fatal("invalid import accepted")
			}
		} else if err != nil || d.Id() != "task" || d.Get("instance_id") != "instance" {
			t.Fatalf("composite import changed: %v", err)
		}
	}
}
