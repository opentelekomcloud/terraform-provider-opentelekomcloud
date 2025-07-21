package smn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/templates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSmnMessageTemplatesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSmnMessageTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Elem:     smnMessageTemplateMessageTemplateSchema(),
				Computed: true,
			},
		},
	}
}

func smnMessageTemplateMessageTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceSmnMessageTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	var mErr *multierror.Error

	listSmnMessageTemplate, err := templates.List(client, templates.ListOpts{
		Name:     d.Get("name").(string),
		Protocol: d.Get("protocol").(string),
	})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud Smn Message Template")
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	if listSmnMessageTemplate == nil || listSmnMessageTemplate.MessageTemplates == nil {
		mErr = multierror.Append(mErr, d.Set("templates", []interface{}{}))
		return diag.FromErr(mErr.ErrorOrNil())
	}

	templateID, hasTemplateID := d.GetOk("template_id")

	rst := make([]interface{}, 0, len(listSmnMessageTemplate.MessageTemplates))
	for _, v := range listSmnMessageTemplate.MessageTemplates {
		if hasTemplateID && templateID.(string) != v.MessageTemplateID {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":         v.MessageTemplateID,
			"name":       v.Name,
			"protocol":   v.Protocol,
			"tag_names":  v.TagNames,
			"created_at": v.CreateTime,
			"updated_at": v.UpdateTime,
		})
	}

	mErr = multierror.Append(mErr, d.Set("templates", rst))

	return diag.FromErr(mErr.ErrorOrNil())
}
