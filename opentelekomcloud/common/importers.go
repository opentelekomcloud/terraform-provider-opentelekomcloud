package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ImportAsManaged(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("shared", false)
	return []*schema.ResourceData{d}, nil
}

// SetComplexID setting ID from multiple attributes by names
//
//	SetComplexID(d, "provider", "protocol")
//
// will set ID from components:
//
//	"<provider>/<protocol>"
func SetComplexID(d *schema.ResourceData, attributes ...string) error {
	parts := make([]string, len(attributes))
	for i, name := range attributes {
		v, ok := d.Get(name).(string)
		if !ok {
			return fmt.Errorf("all ID components must be strings, but %s is not", name)
		}
		parts[i] = v
	}
	d.SetId(strings.Join(parts, "/"))
	return nil
}

// SetIDComponents setting attributes from complex ID, e.g.:
//
//	SetIDComponents(d, "provider", "protocol")
//
// will set fields from ID formatted as
//
//	<provider>/<protocol>
func SetIDComponents(d *schema.ResourceData, attributes ...string) error {
	id := d.Id()
	parts := strings.SplitN(id, "/", len(attributes))

	if len(parts) != len(attributes) {
		attrsDescription := strings.Join(attributes, "/")
		return fmt.Errorf("resource ID should have format %s, but is %s", attrsDescription, id)
	}

	mErr := &multierror.Error{}
	for i, attr := range attributes {
		v := parts[i]
		if attr == "id" {
			d.SetId(v)
			continue
		}
		mErr = multierror.Append(mErr, d.Set(attr, v))
	}
	return mErr.ErrorOrNil()
}

// ImportByPath can be used to import resource by complex ID
// (e.g. identity protocol by `<provider>/<identity>` or CCE addon by `<cluster_id>/<addon_id>`)
//
// Usage in schema:
//
//	StateContext: common.ImportByPath("provider", "protocol"),
func ImportByPath(attributes ...string) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		err := SetIDComponents(d, attributes...)
		if err != nil {
			return nil, err
		}
		return schema.ImportStatePassthroughContext(ctx, d, meta)
	}
}
