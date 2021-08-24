package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ImportAsManaged(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("shared", false)
	return []*schema.ResourceData{d}, nil
}

// ImportByPath can be used to import resource by complex ID
// (e.g. identity protocol by `<provider>/<identity>` or CCE addon by `<cluster_id>/<addon_id>`)
//
// Usage in schema:
//   StateContext: common.ImportByPath("provider", "protocol"),
func ImportByPath(attributes ...string) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		id := d.Id()
		parts := strings.SplitN(id, "/", len(attributes))

		if len(parts) != len(attributes) {
			attrsDescription := strings.Join(attributes, "/")
			return nil, fmt.Errorf("resource ID should have format %s, but is %s", attrsDescription, id)
		}

		for i, attr := range attributes {
			_ = d.Set(attr, parts[i])
		}
		return schema.ImportStatePassthroughContext(ctx, d, meta)
	}
}
