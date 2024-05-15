package fgs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/events"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceFgsEventV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionEventV2Create,
		ReadContext:   resourceFunctionEventV2Read,
		UpdateContext: resourceFunctionEventV2Update,
		DeleteContext: resourceFunctionEventV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("function_urn", "id"),
		},

		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFunctionEventV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	eventOpts := events.CreateOpts{
		FuncUrn: d.Get("function_urn").(string),
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
	}

	eventResp, err := events.Create(fgsClient, eventOpts)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph function event: %s", err)
	}

	d.SetId(eventResp.Id)

	clientCtx := common.CtxWithClient(ctx, fgsClient, fgsClientV2)
	return resourceFunctionEventV2Read(clientCtx, d, meta)
}

func resourceFunctionEventV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	getResp, err := events.Get(fgsClient, d.Get("function_urn").(string), d.Id())
	if err != nil {
		return diag.Errorf("error retrieving function event (%s): %s", d.Id(), err)
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("name", getResp.Name),
		d.Set("content", getResp.Content),
		d.Set("updated_at", getResp.LastModified),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving function event fields: %s", err)
	}
	return nil
}

func resourceFunctionEventV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	_, err = events.Update(fgsClient, events.UpdateOpts{
		EventId: d.Id(),
		FuncUrn: d.Get("function_urn").(string),
		Content: d.Get("content").(string),
	})
	if err != nil {
		return diag.Errorf("error updating FunctionGraph function event (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, fgsClient, fgsClientV2)
	return resourceFunctionEventV2Read(clientCtx, d, meta)
}

func resourceFunctionEventV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = events.Delete(fgsClient, d.Get("function_urn").(string), d.Id())
	if err != nil {
		return diag.Errorf("error deleting FunctionGraph function event: %s", err)
	}
	return nil
}
