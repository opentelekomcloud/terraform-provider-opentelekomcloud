package lts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSNewGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsGroupV2Create,
		ReadContext:   resourceLtsGroupV2Read,
		UpdateContext: resourceLtsGroupV2Update,
		DeleteContext: resourceLtsGroupV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_alias": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
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

func resourceLtsGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	logGroupId, err := groups.Create(client, groups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTLInDays:    d.Get("ttl_in_days").(int),
		Tags:         ltsTags(d),
		Alias:        d.Get("group_alias").(string),
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v2 log group: %s", err)
	}

	d.SetId(logGroupId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceLtsGroupV2Read(clientCtx, d, meta)
}

func resourceLtsGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	requestResp, err := groups.List(client)
	if err != nil {
		return diag.FromErr(err)
	}
	var groupResult groups.LogGroup
	for _, gr := range requestResp {
		if gr.LogGroupId == d.Id() {
			groupResult = gr
			break
		}
	}
	if groupResult.LogGroupId == "" {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 log group by its ID (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("group_name", groupResult.LogGroupName),
		d.Set("enterprise_project_id", groupResult.Tag["_sys_enterprise_project_id"]),
		d.Set("tags", ignoreSysEpsTag(groupResult.Tag)),
		d.Set("ttl_in_days", groupResult.TTLInDays),
		d.Set("created_at", common.FormatTimeStampRFC3339(groupResult.CreationTime/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLtsGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if d.HasChange("ttl_in_days") {
		_, err = groups.Update(client, groups.UpdateLogGroupOpts{
			LogGroupId: d.Id(),
			TTLInDays:  int32(d.Get("ttl_in_days").(int)),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		tagErr := updateTags(d, meta, "groups", d.Id())
		if tagErr != nil {
			return diag.Errorf("unable to update tags for OpenTelekomCloud LTS v2 log group: %s, err:%s", d.Id(), tagErr)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceLtsGroupV2Read(clientCtx, d, meta)
}

func resourceLtsGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = groups.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v2 log group")
	}
	return nil
}
