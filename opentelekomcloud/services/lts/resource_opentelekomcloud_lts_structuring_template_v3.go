package lts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	cloud_structuring "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/cloud-structuring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLtsStructuringTemplateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsStructuringTemplateV3Create,
		UpdateContext: resourceLtsStructuringTemplateV3Update,
		ReadContext:   resourceLtsStructuringTemplateV3Read,
		DeleteContext: resourceLtsStructuringTemplateV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLtsStructuringTemplateV3ImportState,
		},
		CustomizeDiff: validateLtsStructuringTemplateV3,

		Schema: map[string]*schema.Schema{
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"built_in", "custom"}, false),
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"demo_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     ltsStructuringTemplateFieldSchema(),
			},
			"tag_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     ltsStructuringTemplateFieldSchema(),
			},
			"quick_analysis": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"demo_log": {
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

func ltsStructuringTemplateFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"is_analysis": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func validateLtsStructuringTemplateV3(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	templateType := d.Get("template_type").(string)
	_, templateIDSet := d.GetOk("template_id")

	if templateType == "custom" && !templateIDSet {
		return fmt.Errorf("template_id must be specified when template_type is custom")
	}
	if templateType == "built_in" && templateIDSet {
		return fmt.Errorf("template_id can only be specified when template_type is custom")
	}
	return nil
}

func resourceLtsStructuringTemplateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = cloud_structuring.Create(client, buildLtsStructuringTemplateOpts(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v3 structuring template: %s", err)
	}

	return resourceLtsStructuringTemplateV3Read(ctx, d, meta)
}

func buildLtsStructuringTemplateOpts(d *schema.ResourceData) cloud_structuring.CreateOpts {
	opts := cloud_structuring.CreateOpts{
		LogGroupId:  d.Get("log_group_id").(string),
		LogStreamId: d.Get("log_stream_id").(string),
		// The LTS API requires template_id to be present even for built-in
		// templates. A non-nil pointer preserves the required empty string in
		// the JSON payload instead of encoding the field as null or omitting it.
		TemplateId:    pointerto.String(d.Get("template_id").(string)),
		Name:          d.Get("template_name").(string),
		Type:          d.Get("template_type").(string),
		DemoFields:    expandLtsStructuringTemplateFields(d.Get("demo_fields")),
		TagFields:     expandLtsStructuringTemplateFields(d.Get("tag_fields")),
		QuickAnalysis: pointerto.Bool(d.Get("quick_analysis").(bool)),
	}
	return opts
}

func expandLtsStructuringTemplateFields(raw interface{}) []cloud_structuring.Field {
	fieldsRaw, ok := raw.([]interface{})
	if !ok || len(fieldsRaw) == 0 {
		return nil
	}

	fields := make([]cloud_structuring.Field, 0, len(fieldsRaw))
	for _, fieldRaw := range fieldsRaw {
		field, ok := fieldRaw.(map[string]interface{})
		if !ok {
			continue
		}
		fields = append(fields, cloud_structuring.Field{
			Name:       field["field_name"].(string),
			IsAnalysis: pointerto.Bool(field["is_analysis"].(bool)),
		})
	}
	return fields
}

func resourceLtsStructuringTemplateV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	template, err := getLtsStructuringTemplate(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud LTS v3 structuring template")
	}
	if template.ID == "" {
		d.SetId("")
		return nil
	}

	d.SetId(template.ID)
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("log_group_id", template.LogGroupId),
		d.Set("log_stream_id", template.LogStreamId),
		d.Set("template_name", template.Name),
		d.Set("demo_log", template.DemoLog),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getLtsStructuringTemplate(client *golangsdk.ServiceClient, d *schema.ResourceData) (*cloud_structuring.StructuringResponse, error) {
	template, err := cloud_structuring.Get(
		client,
		d.Get("log_group_id").(string),
		d.Get("log_stream_id").(string),
	)
	// The API returns an empty, JSON-encoded string when no structuring rule
	// exists. The SDK reports that response as io.EOF, so normalize it to 404.
	if errors.Is(err, io.EOF) {
		return nil, golangsdk.ErrDefault404{}
	}
	return template, err
}

func resourceLtsStructuringTemplateV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = cloud_structuring.Update(client, buildLtsStructuringTemplateOpts(d))
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud LTS v3 structuring template: %s", err)
	}

	return resourceLtsStructuringTemplateV3Read(ctx, d, meta)
}

func resourceLtsStructuringTemplateV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = cloud_structuring.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v3 structuring template")
	}
	return nil
}

func resourceLtsStructuringTemplateV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid ID format, want '<log_group_id>/<log_stream_id>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("log_group_id", parts[0]),
		d.Set("log_stream_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
