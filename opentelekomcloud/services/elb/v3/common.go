package v3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

type ContextKey string

const (
	keyClient       = ContextKey("lbv3-client")
	ErrCreateClient = "error creating ELBv3 client: %w"
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
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateClient, err)
	}
	return client.(*golangsdk.ServiceClient), nil
}
