package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/streams"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSStreamV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamV2Create,
		ReadContext:   resourceStreamV2Read,
		UpdateContext: resourceStreamV2Update,
		DeleteContext: resourceStreamV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStreamV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stream_alias": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStreamV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := streams.CreateOpts{
		GroupId:       d.Get("group_id").(string),
		LogStreamName: d.Get("stream_name").(string),
		Tags:          ltsTags(d),
		Alias:         d.Get("stream_alias").(string),
		// EnterpriseProjectName: "",
	}

	streamId, err := streams.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v2 log stream: %s", err)
	}
	d.SetId(streamId)

	if _, ok := d.GetOk("ttl_in_days"); ok {
		_, err = streams.Update(client, streams.UpdateLogStreamOpts{
			GroupId:   d.Get("group_id").(string),
			StreamId:  d.Id(),
			TTLInDays: d.Get("ttl_in_days").(int),
		})
		if err != nil {
			return diag.Errorf("error setting TTL for OpenTelekomCloud LTS v2 log stream %s: %s", streamId, err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceStreamV2Read(clientCtx, d, meta)
}

func resourceStreamV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	requestResp, err := streams.List(client, d.Get("group_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	var streamResult *streams.LogStream
	for _, stream := range requestResp {
		if stream.LogStreamId == d.Id() {
			streamResult = &stream
		}
	}
	if streamResult == nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 log stream by its ID (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("stream_name", streamResult.LogStreamName),
		d.Set("ttl_in_days", streamResult.TTLInDays),
		d.Set("enterprise_project_id", streamResult.Tag["_sys_enterprise_project_id"]),
		d.Set("tags", ignoreSysEpsTag(streamResult.Tag)),
		d.Set("filter_count", streamResult.FilterCount),
		d.Set("created_at", common.FormatTimeStampRFC3339(streamResult.CreationTime/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceStreamV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if d.HasChange("ttl_in_days") {
		_, err = streams.Update(client, streams.UpdateLogStreamOpts{
			GroupId:   d.Get("group_id").(string),
			StreamId:  d.Id(),
			TTLInDays: d.Get("ttl_in_days").(int),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		tagErr := updateTags(d, meta, "topics", d.Id())
		if tagErr != nil {
			return diag.Errorf("unable to update tags for OpenTelekomCloud LTS v2 log stream: %s, err:%s", d.Id(), tagErr)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceStreamV2Read(clientCtx, d, meta)
}

func resourceStreamV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	err = streams.Delete(client, streams.DeleteOpts{
		GroupId:  d.Get("group_id").(string),
		StreamId: d.Id(),
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v2 log stream")
	}
	return nil
}

func resourceStreamV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, want '<group_id>/<stream_id>', but '%s'", d.Id())
	}

	groupID := parts[0]
	streamID := parts[1]

	d.SetId(streamID)
	mErr := multierror.Append(nil,
		d.Set("group_id", groupID),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
