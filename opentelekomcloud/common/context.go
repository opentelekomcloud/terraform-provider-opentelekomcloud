package common

import (
	"context"
	"fmt"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

type ClientCtxKey string

func CtxWithClient(parent context.Context, client *golangsdk.ServiceClient, keyClient ClientCtxKey) context.Context {
	return context.WithValue(parent, keyClient, client)
}

type ContextCreatingFunction func() (*golangsdk.ServiceClient, error)

func ClientFromCtx(ctx context.Context, keyClient string, fallbackFunc ContextCreatingFunction) (*golangsdk.ServiceClient, error) {
	client := ctx.Value(keyClient)
	if client != nil {
		return client.(*golangsdk.ServiceClient), nil
	}
	if fallbackFunc == nil {
		return nil, fmt.Errorf("no client generation function is provided as a fallback")
	}
	client, err := fallbackFunc()
	if err != nil {
		return nil, err
	}
	return client.(*golangsdk.ServiceClient), nil
}
