package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

type ContextKey string

const (
	clientCreationFail    = "error creating OpenTelekomCloud identity v3 client: %w"
	clientV30CreationFail = "error creating OpenTelekomCloud identity v3.0 client: %w"
	keyClient             = ContextKey("client")
	keyClientV3           = "iam-v3-client"
	keyClientV30          = "iam-v30-client"
)

func ctxWithClient(parent context.Context, client *golangsdk.ServiceClient) context.Context {
	return context.WithValue(parent, keyClient, client)
}

func clientFromCtx(ctx context.Context, d *schema.ResourceData, meta interface{}) (*golangsdk.ServiceClient, error) {
	client := ctx.Value(keyClient)
	if client != nil {
		return client.(*golangsdk.ServiceClient), nil
	}
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf(clientCreationFail, err)
	}
	return client.(*golangsdk.ServiceClient), nil
}

func identityExtClient(d *schema.ResourceData, meta interface{}) (*golangsdk.ServiceClient, error) {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating identity v3-ext client: %w", err)
	}
	client.Endpoint = strings.Replace(client.Endpoint, "/v3/", "/v3-ext/", 1)
	return client, nil
}

func buildStatementsSet(role *policies.Policy) ([]interface{}, error) {
	statements := make([]interface{}, len(role.Policy.Statement))
	for i, statement := range role.Policy.Statement {
		var condition string
		if len(statement.Condition) > 0 {
			jsonOutput, err := json.Marshal(statement.Condition)
			if err != nil {
				return nil, err
			}
			condition = string(jsonOutput)
		}
		resource, err := flattenStatementResource(statement.Resource)
		if err != nil {
			return nil, fmt.Errorf("statement %d: %w", i, err)
		}
		statements[i] = map[string]interface{}{
			"effect":    statement.Effect,
			"action":    statement.Action,
			"resource":  resource,
			"condition": condition,
		}
	}

	return statements, nil
}

// flattenStatementResource converts the policy `Resource` into the `resource` field.
// `Resource` is either a list of resource paths or an object like `{"uri": ["/iam/agencies/<id>"]}`,
// in which case the URIs are returned as-is.
func flattenStatementResource(raw interface{}) ([]interface{}, error) {
	switch r := raw.(type) {
	case nil:
		return nil, nil
	case []interface{}, []string:
		return toInterfaceSlice(r)
	case map[string]interface{}:
		uris, ok := r["uri"]
		if !ok || len(r) != 1 {
			return nil, fmt.Errorf("unsupported Resource object: %v", r)
		}
		return toInterfaceSlice(uris)
	case map[string][]string:
		converted := make(map[string]interface{}, len(r))
		for k, v := range r {
			converted[k] = v
		}
		return flattenStatementResource(converted)
	default:
		return nil, fmt.Errorf("unsupported Resource type %T", r)
	}
}

func toInterfaceSlice(raw interface{}) ([]interface{}, error) {
	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			if _, ok := item.(string); !ok {
				return nil, fmt.Errorf("unsupported Resource item type %T", item)
			}
		}
		return v, nil
	case []string:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported Resource list type %T", v)
	}
}

func findKeyByValue(m map[string]string, searchValue string) (string, bool) {
	for key, value := range m {
		if value == searchValue {
			return key, true
		}
	}
	return "", false
}
