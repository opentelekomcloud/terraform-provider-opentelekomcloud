package iam

import (
	"context"
	"fmt"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const clientCreationFail = "error creating OpenTelekomCloud identity v3 client: %w"

func ctxWithClient(parent context.Context, client *golangsdk.ServiceClient) context.Context {
	return context.WithValue(parent, "client", client)
}

func clientFromCtx(ctx context.Context, meta interface{}) (*golangsdk.ServiceClient, error) {
	client := ctx.Value("client")
	if client != nil {
		return client.(*golangsdk.ServiceClient), nil
	}
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf(clientCreationFail, err)
	}
	return client.(*golangsdk.ServiceClient), nil
}
