package swr

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSwrPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"date_rule", "tag_rule",
							}, false),
						},
						"params": {
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tag_selector": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"label", "regexp",
										}, false),
									},
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	namespace := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	opts := policy.CreateOpts{
		Algorithm: d.Get("algorithm").(string),
		Rules:     getRules(d),
	}

	id, err := policy.Create(client, namespace, repository, opts)
	if err != nil {
		return fmterr.Errorf("error creating retention policy: %w", err)
	}
	d.SetId(strconv.Itoa(id))

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	namespace := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	retPolicy, err := policy.Get(client, namespace, repository, d.Id())
	if err != nil {
		return fmterr.Errorf("error reading retention policy: %w", err)
	}

	mErr := multierror.Append(
		d.Set("algorithm", retPolicy.Algorithm),
		d.Set("rules", setRules(retPolicy.Rules)),
		d.Set("scope", retPolicy.Scope),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting resource fields: %w", err)
	}

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	namespace := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	opts := policy.UpdateOpts{
		Algorithm: d.Get("algorithm").(string),
		Rules:     getRules(d),
	}
	err = policy.Update(client, namespace, repository, d.Id(), opts)
	if err != nil {
		return fmterr.Errorf("error updating retention policy: %w", err)
	}

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	namespace := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	err = policy.Delete(client, namespace, repository, d.Id())
	if err != nil {
		fmterr.Errorf("error deleting retention policy: %w", err)
	}

	return nil
}

func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		err := fmt.Errorf("invalid format specified for SWR retention policy import: format must be <organization>/<repository>/<policy_id>")
		return nil, err
	}
	org := parts[0]
	repo := parts[1]
	policyId := parts[2]
	d.SetId(policyId)
	if err := d.Set("organization", org); err != nil {
		return nil, err
	}
	if err := d.Set("repository", repo); err != nil {
		return nil, err
	}
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

func getRules(d *schema.ResourceData) []policy.Rule {
	rulesRaw := d.Get("rules").([]interface{})

	rules := make([]policy.Rule, 0, len(rulesRaw))
	for _, r := range rulesRaw {
		ruleRaw := r.(map[string]interface{})
		rule := policy.Rule{
			Template:     ruleRaw["template"].(string),
			Params:       common.ConvertToMapString(ruleRaw["params"]),
			TagSelectors: getTagSelectors(ruleRaw["tag_selector"].([]interface{})),
		}
		rules = append(rules, rule)
	}
	return rules
}

func getTagSelectors(tagSelectorsInput []interface{}) []policy.TagSelector {
	result := make([]policy.TagSelector, 0, len(tagSelectorsInput))
	for _, val := range tagSelectorsInput {
		tagSelectorInput := val.(map[string]interface{})
		tagSelector := policy.TagSelector{
			Kind:    tagSelectorInput["kind"].(string),
			Pattern: tagSelectorInput["pattern"].(string),
		}
		result = append(result, tagSelector)
	}
	return result
}

func setRules(rulesInResp []policy.Rule) []map[string]interface{} {
	var rules []map[string]interface{}
	for _, ruleInResp := range rulesInResp {
		rule := map[string]interface{}{
			"template":     ruleInResp.Template,
			"params":       ruleInResp.Params,
			"tag_selector": setTagSelectors(ruleInResp.TagSelectors),
		}
		rules = append(rules, rule)
	}
	return rules
}

func setTagSelectors(tsInResp []policy.TagSelector) []map[string]interface{} {
	var tags []map[string]interface{}
	for _, tagInResp := range tsInResp {
		tag := map[string]interface{}{
			"kind":    tagInResp.Kind,
			"pattern": tagInResp.Pattern,
		}
		tags = append(tags, tag)
	}
	return tags
}
