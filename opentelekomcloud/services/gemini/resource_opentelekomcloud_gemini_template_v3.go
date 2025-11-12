package gemini

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceGeminiDBParameterTemplateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiParameterTemplateCreate,
		ReadContext:   resourceGeminiParameterTemplateRead,
		UpdateContext: resourceGeminiParameterTemplateUpdate,
		DeleteContext: resourceGeminiParameterTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Elem:     templateParametersSchema(),
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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

func templateParametersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"need_restart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceGeminiParameterTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	opts := template.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      buildParametersMap(d),
		DataStore: template.DataStoreOpt{
			Type:    d.Get("instance_type").(string),
			Version: d.Get("engine_version").(string),
		},
	}

	result, err := template.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating GeminiDB parameter template: %s", err)
	}

	d.SetId(result.Id)

	return resourceGeminiParameterTemplateRead(ctx, d, meta)
}

func resourceGeminiParameterTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	opts := template.UpdateOpts{
		ConfigId:    d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      buildParametersMap(d),
	}

	err = template.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating GeminiDB parameter template: %s", err)
	}

	return resourceGeminiParameterTemplateRead(ctx, d, meta)
}

func buildParametersMap(d *schema.ResourceData) map[string]string {
	rawParameters := d.Get("parameters").(*schema.Set)
	if rawParameters.Len() == 0 {
		return nil
	}

	rst := make(map[string]string)
	for _, v := range rawParameters.List() {
		if raw, ok := v.(map[string]interface{}); ok {
			rst[raw["name"].(string)] = raw["value"].(string)
		}
	}
	return rst
}

func resourceGeminiParameterTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	var mErr *multierror.Error

	result, err := template.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB parameter template")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", result.Name),
		d.Set("description", result.Description),
		d.Set("engine_version", result.DataStoreVersionName),
		d.Set("instance_type", result.DataStoreName),
		d.Set("created_at", result.Created),
		d.Set("updated_at", result.Updated),
		d.Set("parameters", flattenParameters(d, result.ConfigurationParameters)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenParameters(d *schema.ResourceData, params []template.InstanceParameterResult) []interface{} {
	if len(params) == 0 {
		return nil
	}

	paramsMap := buildParamsMap(d)
	rst := make([]interface{}, 0, len(params))

	for _, param := range params {
		if !paramsMap[param.Name] {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"name":         param.Name,
			"value":        param.Value,
			"need_restart": param.RestartRequired,
			"readonly":     param.Readonly,
			"value_range":  param.ValueRange,
			"data_type":    param.Type,
			"description":  param.Description,
		})
	}
	return rst
}

func buildParamsMap(d *schema.ResourceData) map[string]bool {
	params := d.Get("parameters").(*schema.Set).List()
	paramsMap := make(map[string]bool)
	for _, param := range params {
		if v, ok := param.(map[string]interface{}); ok {
			paramsMap[v["name"].(string)] = true
		}
	}
	return paramsMap
}

func resourceGeminiParameterTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	err = template.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting GeminiDB parameter template")
	}

	return nil
}
