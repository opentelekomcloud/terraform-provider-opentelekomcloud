package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/permissions"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	testPermissionInstanceID = "instance-id"
	testPermissionTopic      = "topic-1"
)

var (
	testPermissionProducer = permissions.CreatePolicy{UserName: "producer", AccessPolicy: "pub"}
	testPermissionConsumer = permissions.CreatePolicy{UserName: "consumer", AccessPolicy: "sub"}
)

// testPermissionAPI is a DMS API that keeps the policies of one topic and every list it was sent.
type testPermissionAPI struct {
	policies []permissions.CreatePolicy
	requests [][]permissions.CreatePolicy
	status   int
}

func (api *testPermissionAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	topicsPath := fmt.Sprintf("/instances/%s/topics/", testPermissionInstanceID)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == topicsPath+"accesspolicy":
		var body struct {
			Topics []permissions.CreateOpts `json:"topics"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Topics) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		api.requests = append(api.requests, body.Topics[0].Policies)
		if api.status != 0 {
			w.WriteHeader(api.status)
			return
		}
		api.policies = body.Topics[0].Policies
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodGet && r.URL.Path == topicsPath+testPermissionTopic+"/accesspolicy":
		topic := permissions.Permissions{Name: testPermissionTopic}
		for _, policy := range api.policies {
			topic.Policies = append(topic.Policies, permissions.Policy{
				UserName:     policy.UserName,
				AccessPolicy: policy.AccessPolicy,
			})
		}
		_ = json.NewEncoder(w).Encode(topic)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func testPermissionConfig(t *testing.T, api *testPermissionAPI) *cfg.Config {
	t.Helper()
	server := httptest.NewServer(api)
	t.Cleanup(server.Close)

	return &cfg.Config{
		Region: "eu-de",
		HwClient: &golangsdk.ProviderClient{
			HTTPClient: *server.Client(),
			EndpointLocator: func(golangsdk.EndpointOpts) (string, error) {
				return server.URL + "/", nil
			},
		},
	}
}

func testPermissionResourceConfig(topic string, policies []permissions.CreatePolicy) *terraform.ResourceConfig {
	rawPolicies := make([]interface{}, len(policies))
	for i, policy := range policies {
		rawPolicies[i] = map[string]interface{}{
			"username":      policy.UserName,
			"access_policy": policy.AccessPolicy,
		}
	}
	return terraform.NewResourceConfigRaw(map[string]interface{}{
		"instance_id": testPermissionInstanceID,
		"topic_name":  topic,
		"policies":    rawPolicies,
	})
}

// testPermissionApply plans and applies the policies and returns the state and the plan.
func testPermissionApply(
	t *testing.T, r *schema.Resource, state *terraform.InstanceState, policies []permissions.CreatePolicy, config *cfg.Config,
) (*terraform.InstanceState, *terraform.InstanceDiff, error) {
	t.Helper()
	diff, err := r.Diff(context.Background(), state, testPermissionResourceConfig(testPermissionTopic, policies), config)
	if err != nil {
		t.Fatalf("error planning: %s", err)
	}
	if diff == nil {
		t.Fatal("expected a plan with changes")
	}
	newState, diags := r.Apply(context.Background(), state, diff, config)
	if diags.HasError() {
		return newState, diff, fmt.Errorf("%v", diags)
	}
	return newState, diff, nil
}

func TestDmsUserPermissionV1PolicyChangeIsOneRequest(t *testing.T) {
	cases := map[string][]permissions.CreatePolicy{
		"add a user": {
			testPermissionProducer, testPermissionConsumer, {UserName: "consumer-2", AccessPolicy: "sub"},
		},
		"remove a user":      {testPermissionProducer},
		"change a user":      {testPermissionProducer, {UserName: "consumer", AccessPolicy: "all"}},
		"replace every user": {{UserName: "consumer-2", AccessPolicy: "all"}},
	}

	for name, changed := range cases {
		t.Run(name, func(t *testing.T) {
			api := &testPermissionAPI{}
			config := testPermissionConfig(t, api)
			r := ResourceDmsUsersPermissionV1()

			before := []permissions.CreatePolicy{testPermissionProducer, testPermissionConsumer}
			state, _, err := testPermissionApply(t, r, nil, before, config)
			if err != nil {
				t.Fatalf("error creating: %s", err)
			}
			api.requests = nil

			state, diff, err := testPermissionApply(t, r, state, changed, config)
			if err != nil {
				t.Fatalf("error updating: %s", err)
			}
			if diff.RequiresNew() {
				t.Error("expected the policies to be updated in place")
			}
			if want := [][]permissions.CreatePolicy{changed}; !reflect.DeepEqual(api.requests, want) {
				t.Errorf("expected the requests %v, got %v", want, api.requests)
			}
			if state.ID != testPermissionTopic {
				t.Errorf("expected the ID %q, got %q", testPermissionTopic, state.ID)
			}

			after, err := r.Diff(context.Background(), state, testPermissionResourceConfig(testPermissionTopic, changed), config)
			if err != nil {
				t.Fatalf("error planning: %s", err)
			}
			if after != nil && !after.Empty() {
				t.Errorf("expected no changes after the update, got %v", after)
			}
		})
	}
}

func TestDmsUserPermissionV1UpdateErrorIsReturned(t *testing.T) {
	api := &testPermissionAPI{}
	config := testPermissionConfig(t, api)
	r := ResourceDmsUsersPermissionV1()

	state, _, err := testPermissionApply(t, r, nil, []permissions.CreatePolicy{testPermissionProducer}, config)
	if err != nil {
		t.Fatalf("error creating: %s", err)
	}

	api.status = http.StatusBadRequest
	changed := []permissions.CreatePolicy{testPermissionProducer, testPermissionConsumer}
	state, _, err = testPermissionApply(t, r, state, changed, config)
	if err == nil {
		t.Fatal("expected the error of the API")
	}
	if state == nil || state.ID != testPermissionTopic {
		t.Fatalf("expected the resource to stay in state, got %v", state)
	}
	if want := []permissions.CreatePolicy{testPermissionProducer}; !reflect.DeepEqual(api.policies, want) {
		t.Errorf("expected the topic to keep the policies %v, got %v", want, api.policies)
	}
}

func TestDmsUserPermissionV1TopicChangeRecreates(t *testing.T) {
	config := testPermissionConfig(t, &testPermissionAPI{})
	r := ResourceDmsUsersPermissionV1()
	policies := []permissions.CreatePolicy{testPermissionProducer}

	state, _, err := testPermissionApply(t, r, nil, policies, config)
	if err != nil {
		t.Fatalf("error creating: %s", err)
	}
	diff, err := r.Diff(context.Background(), state, testPermissionResourceConfig("topic-2", policies), config)
	if err != nil {
		t.Fatalf("error planning: %s", err)
	}
	if !diff.RequiresNew() {
		t.Error("expected another topic to need a new resource")
	}
}
