package smn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/templates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSmnMessageTemplateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmnMessageTemplateCreateV2,
		UpdateContext: resourceSmnMessageTemplateUpdateV2,
		ReadContext:   resourceSmnMessageTemplateReadV2,
		DeleteContext: resourceSmnMessageTemplateDeleteV2,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSmnMessageTemplateCreateV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := templates.CreateOpts{
		Name:     d.Get("name").(string),
		Content:  d.Get("content").(string),
		Protocol: d.Get("protocol").(string),
	}

	smnTemplate, err := templates.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating SMN message template: %s", err)
	}

	d.SetId(smnTemplate.MessageTemplateID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceSmnMessageTemplateReadV2(clientCtx, d, meta)
}

func resourceSmnMessageTemplateReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	var mErr *multierror.Error

	getMessageTemplate, err := templates.Get(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", getMessageTemplate.Name),
		d.Set("protocol", getMessageTemplate.Protocol),
		d.Set("tag_names", getMessageTemplate.TagNames),
		d.Set("content", getMessageTemplate.Content),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSmnMessageTemplateUpdateV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if d.HasChanges("content") {
		updateOpts := templates.UpdateOpts{
			TemplateID: d.Id(),
			Content:    d.Get("content").(string),
		}

		_, err = templates.Update(client, updateOpts)
		if err != nil {
			return diag.Errorf("error updating SMN message template: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceSmnMessageTemplateReadV2(clientCtx, d, meta)
}

func resourceSmnMessageTemplateDeleteV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	err = templates.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMN message template")
	}

	return nil
}
