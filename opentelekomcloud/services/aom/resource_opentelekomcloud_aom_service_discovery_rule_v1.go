package aom

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/aom/v1/monitoring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const (
	keyClientV1                    = "aom-v1-client"
	errCreationV1Client            = "error creating OpenTelekomCloud AOM V1 client: %w"
	serviceDiscoveryFixedEventName = "aom_inventory_rules_event"
)

func ResourceServiceDiscoveryRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceDiscoveryRuleV1Create,
		ReadContext:   resourceServiceDiscoveryRuleV1Read,
		UpdateContext: resourceServiceDiscoveryRuleV1Update,
		DeleteContext: resourceServiceDiscoveryRuleV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"detect_log": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"is_default_rule": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"is_detect": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
						},
						"log_file_fix": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"log", "trace", "out",
								}, false),
							},
						},
						"discovery_rule": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"check_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"cmdLine", "env", "scope"}, false),
									},
									"check_mode": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"contain", "equals"}, false),
									},
									"check_content": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name_rule": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_name_rule": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     nameExtractionRuleSchema(),
									},
									"application_name_rule": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     nameExtractionRuleSchema(),
									},
								},
							},
						},
						"log_path_rule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"args": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"value": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"app_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"attr_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"cmdLine", "env",
								}, false),
							},
						},
						"data_source": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"event_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nameExtractionRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"cmdLineHash", "cmdLine", "env", "str"}, false),
			},
			"args": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func expandServiceDiscoverySpec(raw []interface{}) monitoring.AppRulesSpec {
	item := raw[0].(map[string]interface{})

	discoveryRules := make([]monitoring.DiscoveryRule, 0)
	for _, r := range item["discovery_rule"].([]interface{}) {
		rule := r.(map[string]interface{})
		discoveryRules = append(discoveryRules, monitoring.DiscoveryRule{
			CheckType:    rule["check_type"].(string),
			CheckMode:    rule["check_mode"].(string),
			CheckContent: common.ExpandToStringList(rule["check_content"].([]interface{})),
		})
	}

	nameRulesRaw := item["name_rule"].([]interface{})[0].(map[string]interface{})
	appNameRules := make([]monitoring.AppNameRule, 0)
	for _, r := range nameRulesRaw["app_name_rule"].([]interface{}) {
		rule := r.(map[string]interface{})
		appNameRules = append(appNameRules, monitoring.AppNameRule{
			NameType: rule["name_type"].(string),
			Args:     common.ExpandToStringList(rule["args"].([]interface{})),
			Value:    common.ExpandToStringList(rule["value"].([]interface{})),
		})
	}
	applicationNameRules := make([]monitoring.ApplicationNameRule, 0)
	for _, r := range nameRulesRaw["application_name_rule"].([]interface{}) {
		rule := r.(map[string]interface{})
		applicationNameRules = append(applicationNameRules, monitoring.ApplicationNameRule{
			NameType: rule["name_type"].(string),
			Args:     common.ExpandToStringList(rule["args"].([]interface{})),
			Value:    common.ExpandToStringList(rule["value"].([]interface{})),
		})
	}

	logPathRules := make([]monitoring.LogPathRule, 0)
	for _, r := range item["log_path_rule"].([]interface{}) {
		rule := r.(map[string]interface{})
		logPathRules = append(logPathRules, monitoring.LogPathRule{
			NameType: rule["name_type"].(string),
			Args:     common.ExpandToStringList(rule["args"].([]interface{})),
			Value:    common.ExpandToStringList(rule["value"].([]interface{})),
		})
	}

	return monitoring.AppRulesSpec{
		AppType:       item["app_type"].(string),
		AttrList:      common.ExpandToStringList(item["attr_list"].([]interface{})),
		DetectLog:     strconv.FormatBool(item["detect_log"].(bool)),
		IsDefaultRule: strconv.FormatBool(item["is_default_rule"].(bool)),
		IsDetect:      strconv.FormatBool(item["is_detect"].(bool)),
		Priority:      item["priority"].(int),
		LogFileFix:    common.ExpandToStringList(item["log_file_fix"].([]interface{})),
		DataSource:    item["data_source"].(string),
		DiscoveryRule: discoveryRules,
		NameRule: monitoring.NameRule{
			AppNameRule:         appNameRules,
			ApplicationNameRule: applicationNameRules,
		},
		LogPathRule: logPathRules,
	}
}

func flattenServiceDiscoverySpec(spec monitoring.AppRulesSpec) ([]map[string]interface{}, error) {
	detectLog, err := strconv.ParseBool(spec.DetectLog)
	if err != nil {
		return nil, err
	}
	isDefaultRule, err := strconv.ParseBool(spec.IsDefaultRule)
	if err != nil {
		return nil, err
	}
	isDetect, err := strconv.ParseBool(spec.IsDetect)
	if err != nil {
		return nil, err
	}

	discoveryRules := make([]map[string]interface{}, 0, len(spec.DiscoveryRule))
	for _, rule := range spec.DiscoveryRule {
		discoveryRules = append(discoveryRules, map[string]interface{}{
			"check_type":    rule.CheckType,
			"check_mode":    rule.CheckMode,
			"check_content": rule.CheckContent,
		})
	}

	appNameRules := make([]map[string]interface{}, 0, len(spec.NameRule.AppNameRule))
	for _, rule := range spec.NameRule.AppNameRule {
		appNameRules = append(appNameRules, map[string]interface{}{
			"name_type": rule.NameType,
			"args":      rule.Args,
			"value":     rule.Value,
		})
	}
	applicationNameRules := make([]map[string]interface{}, 0, len(spec.NameRule.ApplicationNameRule))
	for _, rule := range spec.NameRule.ApplicationNameRule {
		applicationNameRules = append(applicationNameRules, map[string]interface{}{
			"name_type": rule.NameType,
			"args":      rule.Args,
			"value":     rule.Value,
		})
	}

	logPathRules := make([]map[string]interface{}, 0, len(spec.LogPathRule))
	for _, rule := range spec.LogPathRule {
		logPathRules = append(logPathRules, map[string]interface{}{
			"name_type": rule.NameType,
			"args":      rule.Args,
			"value":     rule.Value,
		})
	}

	return []map[string]interface{}{{
		"detect_log":      detectLog,
		"is_default_rule": isDefaultRule,
		"is_detect":       isDetect,
		"priority":        spec.Priority,
		"log_file_fix":    spec.LogFileFix,
		"discovery_rule":  discoveryRules,
		"name_rule": []map[string]interface{}{{
			"app_name_rule":         appNameRules,
			"application_name_rule": applicationNameRules,
		}},
		"log_path_rule": logPathRules,
		"app_type":      spec.AppType,
		"attr_list":     spec.AttrList,
		"data_source":   spec.DataSource,
	}}, nil
}

func buildAppRule(d *schema.ResourceData, config *cfg.Config) monitoring.AppRule {
	enabled := d.Get("enabled").(bool)
	return monitoring.AppRule{
		Name:      d.Get("name").(string),
		Enable:    &enabled,
		EventName: serviceDiscoveryFixedEventName,
		Hostid:    []string{},
		Projectid: config.TenantID,
		Desc:      d.Get("description").(string),
		Spec:      expandServiceDiscoverySpec(d.Get("spec").([]interface{})),
	}
}

func resourceServiceDiscoveryRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AomV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	// On creation Id and CreateTime are left blank
	rule := buildAppRule(d, config)
	opts := monitoring.AddOrModifyServiceDiscoveryRuleOpts{
		AppRules: []monitoring.AppRule{rule},
	}
	resp, err := monitoring.AddOrModifyServiceDiscoveryRule(client, opts)
	if err != nil {
		return diag.Errorf("error creating AOM service discovery rule (%s): %s", rule.Name, err)
	}
	if len(resp.Results) == 0 || resp.Results[0]["id"] == "" {
		return fmterr.Errorf("no rule ID returned in AOM service discovery rule (%s) creation response", rule.Name)
	}

	d.SetId(resp.Results[0]["id"])

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceServiceDiscoveryRuleV1Read(clientCtx, d, meta)
}

func resourceServiceDiscoveryRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AomV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	rule, err := GetServiceDiscoveryRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AOM service discovery rule")
	}

	spec, err := flattenServiceDiscoverySpec(rule.Spec)
	if err != nil {
		return diag.Errorf("error parsing AOM service discovery rule (%s) spec: %s", d.Id(), err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", rule.Name),
		d.Set("enabled", rule.Enable != nil && *rule.Enable),
		d.Set("description", rule.Desc),
		d.Set("event_name", rule.EventName),
		d.Set("project_id", rule.Projectid),
		d.Set("create_time", rule.CreateTime),
		d.Set("spec", spec),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving AOM service discovery rule (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceServiceDiscoveryRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AomV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	// The modify operation uses the same endpoint as add and identifies the
	// rule by Id, requiring the CreateTime returned by the API.
	rule := buildAppRule(d, config)
	rule.Id = d.Id()
	rule.CreateTime = d.Get("create_time").(string)
	opts := monitoring.AddOrModifyServiceDiscoveryRuleOpts{
		AppRules: []monitoring.AppRule{rule},
	}
	if _, err = monitoring.AddOrModifyServiceDiscoveryRule(client, opts); err != nil {
		return diag.Errorf("error updating AOM service discovery rule (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceServiceDiscoveryRuleV1Read(clientCtx, d, meta)
}

func resourceServiceDiscoveryRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AomV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	_, err = monitoring.DeleteServiceDiscoveryRule(client, []string{d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AOM service discovery rule")
	}
	return nil
}

// GetServiceDiscoveryRule finds the service discovery rule by ID in the rule list.
func GetServiceDiscoveryRule(client *golangsdk.ServiceClient, id string) (*monitoring.AppRule, error) {
	rules, err := monitoring.ListServiceDiscoveryRules(client, monitoring.ListServiceDiscoveryRulesOpts{Id: id})
	if err != nil {
		return nil, err
	}
	for i := range rules {
		if rules[i].Id == id {
			return &rules[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}
