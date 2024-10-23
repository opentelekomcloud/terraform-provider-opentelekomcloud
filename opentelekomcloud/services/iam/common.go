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
		statements[i] = map[string]interface{}{
			"effect":    statement.Effect,
			"action":    statement.Action,
			"resource":  statement.Resource,
			"condition": condition,
		}
	}

	return statements, nil
}

func findKeyByValue(m map[string]string, searchValue string) (string, bool) {
	for key, value := range m {
		if value == searchValue {
			return key, true
		}
	}
	return "", false
}
