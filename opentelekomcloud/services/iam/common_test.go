package iam

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/policies"
)

func policyFromJSON(t *testing.T, resource string) *policies.Policy {
	t.Helper()
	body := `{"policy": {"Version": "1.1", "Statement": [{"Effect": "Allow", "Action": ["iam:agencies:assume"]`
	if resource != "" {
		body += `, "Resource": ` + resource
	}
	body += `}]}}`

	var p policies.Policy
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		t.Fatalf("error unmarshalling policy: %s", err)
	}
	return &p
}

func TestBuildStatementsSetResource(t *testing.T) {
	cases := []struct {
		name        string
		resource    string
		expResource []interface{}
		expErr      string
	}{
		{
			name:        "list",
			resource:    `["OBS:*:*:bucket:test-bucket", "OBS:*:*:object:your_object"]`,
			expResource: []interface{}{"OBS:*:*:bucket:test-bucket", "OBS:*:*:object:your_object"},
		},
		{
			name:        "uri object",
			resource:    `{"uri": ["/iam/agencies/id1", "/iam/agencies/id2"]}`,
			expResource: []interface{}{"/iam/agencies/id1", "/iam/agencies/id2"},
		},
		{
			name: "no resource",
		},
		{
			name:     "unsupported object key",
			resource: `{"foo": ["bar"]}`,
			expErr:   "unsupported Resource object",
		},
		{
			name:     "unsupported type",
			resource: `"OBS:*:*:bucket:*"`,
			expErr:   "unsupported Resource type string",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			statements, err := buildStatementsSet(policyFromJSON(t, c.resource))
			if c.expErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.expErr) {
					t.Fatalf("expected error containing %q, got %v", c.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			statement := statements[0].(map[string]interface{})
			if !reflect.DeepEqual(statement["resource"], c.expResource) {
				t.Errorf("resource: expected %v, got %v", c.expResource, statement["resource"])
			}
		})
	}
}

func TestBuildStatementsSetTypedResource(t *testing.T) {
	role := &policies.Policy{}
	role.Policy.Statement = []policies.Statement{
		{Effect: "Allow", Action: []string{"obs:bucket:ListBucket"}, Resource: []string{"OBS:*:*:bucket:*"}},
		{Effect: "Allow", Action: []string{"iam:agencies:assume"}, Resource: map[string][]string{"uri": {"/iam/agencies/id1"}}},
	}

	statements, err := buildStatementsSet(role)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got := statements[0].(map[string]interface{})["resource"]; !reflect.DeepEqual(got, []interface{}{"OBS:*:*:bucket:*"}) {
		t.Errorf("unexpected resource: %v", got)
	}
	if got := statements[1].(map[string]interface{})["resource"]; !reflect.DeepEqual(got, []interface{}{"/iam/agencies/id1"}) {
		t.Errorf("unexpected resource: %v", got)
	}
}

func TestBuildRolePolicyResource(t *testing.T) {
	cases := []struct {
		name   string
		raw    map[string]interface{}
		exp    interface{}
		expErr string
	}{
		{
			name: "resource",
			raw:  map[string]interface{}{"resource": []interface{}{"OBS:*:*:bucket:*"}},
			exp:  []string{"OBS:*:*:bucket:*"},
		},
		{
			name: "uri",
			raw:  map[string]interface{}{"resource": []interface{}{"/iam/agencies/id1", "/iam/agencies/id2"}},
			exp:  map[string][]string{"uri": {"/iam/agencies/id1", "/iam/agencies/id2"}},
		},
		{
			name: "no resource",
			raw:  map[string]interface{}{},
		},
		{
			name:   "mixed uri and path",
			raw:    map[string]interface{}{"resource": []interface{}{"OBS:*:*:bucket:*", "/iam/agencies/id1"}},
			expErr: "can't mix URIs",
		},
		{
			name: "uri with condition",
			raw: map[string]interface{}{
				"resource":  []interface{}{"/iam/agencies/id1"},
				"condition": `{"StringEquals": {"g:ProjectName": ["eu-de"]}}`,
			},
			expErr: "`condition` can't be used with URI resources",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			statement := map[string]interface{}{
				"effect": "Allow",
				"action": []interface{}{"iam:agencies:assume"},
			}
			for k, v := range c.raw {
				statement[k] = v
			}
			d := schema.TestResourceDataRaw(t, ResourceIdentityRoleV3().Schema, map[string]interface{}{
				"description":   "role",
				"display_name":  "role",
				"display_layer": "domain",
				"statement":     []interface{}{statement},
			})

			policy, err := buildRolePolicy(d)
			if c.expErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.expErr) {
					t.Fatalf("expected error containing %q, got %v", c.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(policy.Statement[0].Resource, c.exp) {
				t.Errorf("expected %#v, got %#v", c.exp, policy.Statement[0].Resource)
			}
		})
	}
}
