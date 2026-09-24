package fgs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/alias"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceFgsAliasV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFgsAliasV2Create,
		ReadContext:   resourceFgsAliasV2Read,
		UpdateContext: resourceFgsAliasV2Update,
		DeleteContext: resourceFgsAliasV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFgsAliasV2ImportState,
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
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"additional_version_weights": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"additional_version_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     aliasStrategySchema(),
			},
			"alias_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func aliasStrategySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"combine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     aliasStrategyRulesSchema(),
			},
		},
	}
}

func aliasStrategyRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"param": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"op": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func parseFgsAliasResourceId(resourceId string) (funcUrn, aliasName string) {
	parts := strings.SplitN(resourceId, "/", 2)
	if len(parts) < 2 {
		log.Printf("[ERROR] invalid ID format for FGS alias resource: %s", resourceId)
		return
	}
	funcUrn = parts[0]
	aliasName = parts[1]
	return
}

func expandAliasWeights(raw map[string]interface{}) map[string]int {
	if len(raw) == 0 {
		return nil
	}
	result := make(map[string]int, len(raw))
	for k, v := range raw {
		result[k] = v.(int)
	}
	return result
}

func expandAliasStrategy(raw []interface{}) map[string]alias.VectorStrategy {
	if len(raw) == 0 {
		return nil
	}
	result := make(map[string]alias.VectorStrategy, len(raw))
	for _, v := range raw {
		item := v.(map[string]interface{})
		strategy := alias.VectorStrategy{
			CombineType: item["combine_type"].(string),
		}
		if rules, ok := item["rules"].([]interface{}); ok && len(rules) > 0 {
			strategy.Rules = make([]alias.VersionStrategyRules, 0, len(rules))
			for _, r := range rules {
				rule := r.(map[string]interface{})
				strategy.Rules = append(strategy.Rules, alias.VersionStrategyRules{
					RuleType: rule["rule_type"].(string),
					Param:    rule["param"].(string),
					Op:       rule["op"].(string),
					Value:    rule["value"].(string),
				})
			}
		}
		result[item["version"].(string)] = strategy
	}
	return result
}

func flattenAliasStrategy(strategies map[string]alias.VectorStrategy) []map[string]interface{} {
	if len(strategies) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(strategies))
	for k, v := range strategies {
		item := map[string]interface{}{
			"version":      k,
			"combine_type": v.CombineType,
		}
		if len(v.Rules) > 0 {
			rules := make([]map[string]interface{}, 0, len(v.Rules))
			for _, r := range v.Rules {
				rules = append(rules, map[string]interface{}{
					"rule_type": r.RuleType,
					"param":     r.Param,
					"op":        r.Op,
					"value":     r.Value,
				})
			}
			item["rules"] = rules
		}
		result = append(result, item)
	}
	return result
}

func resourceFgsAliasV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := alias.CreateAliasOpts{
		FuncUrn:                   d.Get("function_urn").(string),
		Name:                      d.Get("name").(string),
		Version:                   d.Get("version").(string),
		Description:               d.Get("description").(string),
		AdditionalVersionWeights:  expandAliasWeights(d.Get("additional_version_weights").(map[string]interface{})),
		AdditionalVersionStrategy: expandAliasStrategy(d.Get("additional_version_strategy").([]interface{})),
	}

	_, err = alias.CreateAlias(fgsClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud FunctionGraph alias: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", createOpts.FuncUrn, createOpts.Name))

	clientCtx := common.CtxWithClient(ctx, fgsClient, fgsClientV2)
	return resourceFgsAliasV2Read(clientCtx, d, meta)
}

func resourceFgsAliasV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	funcUrn, aliasName := parseFgsAliasResourceId(d.Id())

	getResp, err := alias.GetAlias(fgsClient, funcUrn, aliasName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud FunctionGraph alias")
	}

	mErr := multierror.Append(
		d.Set("function_urn", funcUrn),
		d.Set("name", getResp.Name),
		d.Set("version", getResp.Version),
		d.Set("description", getResp.Description),
		d.Set("alias_urn", getResp.AliasUrn),
		d.Set("last_modified", getResp.LastModified),
		d.Set("additional_version_weights", getResp.AdditionalVersionWeights),
		d.Set("additional_version_strategy", flattenAliasStrategy(getResp.AdditionalVersionStrategy)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of FGS alias (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceFgsAliasV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	funcUrn, aliasName := parseFgsAliasResourceId(d.Id())

	updateOpts := alias.UpdateAliasOpts{
		FuncUrn:                   funcUrn,
		AliasName:                 aliasName,
		Version:                   d.Get("version").(string),
		Description:               d.Get("description").(string),
		AdditionalVersionWeights:  expandAliasWeights(d.Get("additional_version_weights").(map[string]interface{})),
		AdditionalVersionStrategy: expandAliasStrategy(d.Get("additional_version_strategy").([]interface{})),
	}

	_, err = alias.UpdateAlias(fgsClient, updateOpts)
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud FunctionGraph alias: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, fgsClient, fgsClientV2)
	return resourceFgsAliasV2Read(clientCtx, d, meta)
}

func resourceFgsAliasV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	funcUrn, aliasName := parseFgsAliasResourceId(d.Id())

	err = alias.Delete(fgsClient, funcUrn, aliasName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud FunctionGraph alias")
	}
	return nil
}

func resourceFgsAliasV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	funcUrn, aliasName := parseFgsAliasResourceId(d.Id())
	mErr := multierror.Append(
		d.Set("function_urn", funcUrn),
		d.Set("name", aliasName),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
