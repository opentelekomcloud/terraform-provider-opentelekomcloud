package lts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	message_template "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/message-template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationTemplateV2Create,
		UpdateContext: resourceNotificationTemplateV2Update,
		ReadContext:   resourceNotificationTemplateV2Read,
		DeleteContext: resourceNotificationTemplateV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"language": {
				Type:     schema.TypeString,
				Required: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Elem:     notificationTemplateSchema(),
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func notificationTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"sub_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func resourceNotificationTemplateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	domainID := config.GetDomainID()
	if domainID == "" {
		return diag.Errorf(errMissingDomainID)
	}

	messageTemplate, err := message_template.Create(client, message_template.CreateOpts{
		DomainId:    domainID,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Source:      d.Get("source").(string),
		Language:    d.Get("language").(string),
		Templates:   buildNotificationTemplate(d.Get("templates")),
	})
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud LTS v2 notification template: %s", err)
	}

	d.SetId(messageTemplate.Name)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNotificationTemplateV2Read(clientCtx, d, meta)
}

func buildNotificationTemplate(rawParams interface{}) []message_template.Templates {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]message_template.Templates, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = message_template.Templates{
					Type:    raw["sub_type"].(string),
					Content: raw["content"].(string),
				}
			}
		}
		return rst
	}
	return nil
}

func resourceNotificationTemplateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	domainID := config.GetDomainID()
	if domainID == "" {
		return diag.Errorf(errMissingDomainID)
	}

	requestResp, err := message_template.List(client, domainID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud LTS v2 notification template")
	}
	var tempResult message_template.MessageTemplateResponse
	for _, t := range requestResp {
		if t.Name == d.Id() {
			tempResult = t
			break
		}
	}
	if tempResult.Name == "" {
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", tempResult.Name),
		d.Set("description", tempResult.Description),
		d.Set("source", tempResult.Source),
		d.Set("language", tempResult.Language),
		d.Set("templates", flattenNotificationTemplateBody(tempResult.Templates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNotificationTemplateBody(resp []message_template.TemplateResponse) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"sub_type": v.Type,
			"content":  v.Content,
		})
	}
	return rst
}

func resourceNotificationTemplateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	domainID := config.GetDomainID()
	if domainID == "" {
		return diag.Errorf(errMissingDomainID)
	}

	_, err = message_template.Update(client, message_template.CreateOpts{
		DomainId:    domainID,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Source:      d.Get("source").(string),
		Language:    d.Get("language").(string),
		Templates:   buildNotificationTemplate(d.Get("templates")),
	})
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud LTS v2 notification template: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNotificationTemplateV2Read(clientCtx, d, meta)
}

func resourceNotificationTemplateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	domainID := config.GetDomainID()
	if domainID == "" {
		return diag.Errorf(errMissingDomainID)
	}

	err = message_template.Delete(client, message_template.DeleteOpts{
		DomainId:      domainID,
		TemplateNames: []string{d.Id()},
	})
	if err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud LTS v2 notification template: %s", err)
	}

	return nil
}
