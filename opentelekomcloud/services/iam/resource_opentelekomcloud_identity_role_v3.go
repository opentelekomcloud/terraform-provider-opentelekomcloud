package iam

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRoleV3Create,
		ReadContext:   resourceIdentityRoleV3Read,
		UpdateContext: resourceIdentityRoleV3Update,
		DeleteContext: resourceIdentityRoleV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_layer": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"statement": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"catalog": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityRoleV3UserInputParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"description":   d.Get("description"),
		"display_layer": d.Get("display_layer"),
		"display_name":  d.Get("display_name"),
		"statement":     d.Get("statement"),
	}
}

func resourceIdentityRoleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	opts := resourceIdentityRoleV3UserInputParams(d)

	r, err := sendIdentityRoleV3CreateRequest(d, opts, nil, client)
	if err != nil {
		return fmterr.Errorf("error creating IdentityRoleV3: %s", err)
	}

	id, err := common.NavigateValue(r, []string{"role", "id"}, nil)
	if err != nil {
		return fmterr.Errorf("error constructing id: %s", err)
	}
	d.SetId(id.(string))

	return resourceIdentityRoleV3Read(ctx, d, meta)
}

func resourceIdentityRoleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	res := make(map[string]interface{})

	err = readIdentityRoleV3Read(ctx, d, client, res)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(setIdentityRoleV3Properties(d, res))
}

func resourceIdentityRoleV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	opts := resourceIdentityRoleV3UserInputParams(d)

	_, err = sendIdentityRoleV3UpdateRequest(d, opts, nil, client)
	if err != nil {
		return fmterr.Errorf("error updating (IdentityRoleV3: %v): %s", d.Id(), err)
	}

	return resourceIdentityRoleV3Read(ctx, d, meta)
}

func resourceIdentityRoleV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	url, err := common.ReplaceVars(d, "OS-ROLE/roles/{id}", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	url = client.ServiceURL(url)

	log.Printf("[DEBUG] Deleting Role %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      common.SuccessHTTPCodes,
		JSONBody:     nil,
		JSONResponse: &r.Body,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	if r.Err != nil {
		return fmterr.Errorf("error deleting Role %q: %s", d.Id(), r.Err)
	}

	return nil
}

func sendIdentityRoleV3CreateRequest(_ *schema.ResourceData, opts map[string]interface{},
	arrayIndex map[string]int, client *golangsdk.ServiceClient) (interface{}, error) {
	url := client.ServiceURL("OS-ROLE/roles")

	params, err := buildIdentityRoleV3CreateParameters(opts, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("error building the request body of api(create)")
	}

	r := golangsdk.Result{}
	_, r.Err = client.Post(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: common.SuccessHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmt.Errorf("error run api(create): %s", r.Err)
	}
	return r.Body, nil
}

func buildIdentityRoleV3CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	descriptionProp, err := common.NavigateValue(opts, []string{"description"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	empty, err := common.IsEmptyValue(reflect.ValueOf(descriptionProp))
	if err != nil {
		return nil, err
	}
	if !empty {
		params["description"] = descriptionProp
	}

	displayNameProp, err := common.NavigateValue(opts, []string{"display_name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	empty, err = common.IsEmptyValue(reflect.ValueOf(displayNameProp))
	if err != nil {
		return nil, err
	}
	if !empty {
		params["display_name"] = displayNameProp
	}

	policyProp, err := expandIdentityRoleV3CreatePolicy(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	empty, err = common.IsEmptyValue(reflect.ValueOf(policyProp))
	if err != nil {
		return nil, err
	}
	if !empty {
		params["policy"] = policyProp
	}

	typeProp, err := expandIdentityRoleV3CreateType(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	empty, err = common.IsEmptyValue(reflect.ValueOf(typeProp))
	if err != nil {
		return nil, err
	}
	if !empty {
		params["type"] = typeProp
	}

	params = map[string]interface{}{"role": params}

	return params, nil
}

func expandIdentityRoleV3CreatePolicy(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	statementProp, err := expandIdentityRoleV3CreatePolicyStatement(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := common.IsEmptyValue(reflect.ValueOf(statementProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["Statement"] = statementProp
	}

	req["Version"] = "1.1"

	return req, nil
}

func expandIdentityRoleV3CreatePolicyStatement(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	v, err := common.NavigateValue(d, []string{"statement"}, newArrayIndex)
	if err != nil {
		return nil, err
	}

	n := len(v.([]interface{}))
	req := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		newArrayIndex["statement"] = i
		transformed := make(map[string]interface{})

		actionProp, err := common.NavigateValue(d, []string{"statement", "action"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		e, err := common.IsEmptyValue(reflect.ValueOf(actionProp))
		if err != nil {
			return nil, err
		}
		if !e {
			transformed["Action"] = actionProp
		}

		effectProp, err := common.NavigateValue(d, []string{"statement", "effect"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		e, err = common.IsEmptyValue(reflect.ValueOf(effectProp))
		if err != nil {
			return nil, err
		}
		if !e {
			transformed["Effect"] = effectProp
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandIdentityRoleV3CreateType(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := common.NavigateValue(d, []string{"display_layer"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v == "domain" {
		return "AX", nil
	} else if v == "project" {
		return "XA", nil
	}
	return nil, fmt.Errorf("unknown display layer:%v", v)
}

func sendIdentityRoleV3UpdateRequest(d *schema.ResourceData, opts map[string]interface{},
	arrayIndex map[string]int, client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := common.ReplaceVars(d, "OS-ROLE/roles/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	params, err := buildIdentityRoleV3UpdateParameters(opts, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("error building the request body of api(update)")
	}

	r := golangsdk.Result{}
	_, r.Err = client.Patch(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: common.SuccessHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmt.Errorf("error run api(update): %s", r.Err)
	}
	return r.Body, nil
}

func buildIdentityRoleV3UpdateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	descriptionProp, err := common.NavigateValue(opts, []string{"description"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := common.IsEmptyValue(reflect.ValueOf(descriptionProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["description"] = descriptionProp
	}

	displayNameProp, err := common.NavigateValue(opts, []string{"display_name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = common.IsEmptyValue(reflect.ValueOf(displayNameProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["display_name"] = displayNameProp
	}

	policyProp, err := expandIdentityRoleV3UpdatePolicy(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = common.IsEmptyValue(reflect.ValueOf(policyProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["policy"] = policyProp
	}

	typeProp, err := expandIdentityRoleV3UpdateType(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = common.IsEmptyValue(reflect.ValueOf(typeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["type"] = typeProp
	}

	params = map[string]interface{}{"role": params}

	return params, nil
}

func expandIdentityRoleV3UpdatePolicy(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	statementProp, err := expandIdentityRoleV3UpdatePolicyStatement(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := common.IsEmptyValue(reflect.ValueOf(statementProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["Statement"] = statementProp
	}

	req["Version"] = "1.1"

	return req, nil
}

func expandIdentityRoleV3UpdatePolicyStatement(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	v, err := common.NavigateValue(d, []string{"statement"}, newArrayIndex)
	if err != nil {
		return nil, err
	}

	n := len(v.([]interface{}))
	req := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		newArrayIndex["statement"] = i
		transformed := make(map[string]interface{})

		actionProp, err := common.NavigateValue(d, []string{"statement", "action"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		e, err := common.IsEmptyValue(reflect.ValueOf(actionProp))
		if err != nil {
			return nil, err
		}
		if !e {
			transformed["Action"] = actionProp
		}

		effectProp, err := common.NavigateValue(d, []string{"statement", "effect"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		e, err = common.IsEmptyValue(reflect.ValueOf(effectProp))
		if err != nil {
			return nil, err
		}
		if !e {
			transformed["Effect"] = effectProp
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandIdentityRoleV3UpdateType(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := common.NavigateValue(d, []string{"display_layer"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v == "domain" {
		return "AX", nil
	} else if v == "project" {
		return "XA", nil
	}
	return nil, fmt.Errorf("unknown display layer:%v", v)
}

func readIdentityRoleV3Read(_ context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, result map[string]interface{}) error {
	url, err := common.ReplaceVars(d, "OS-ROLE/roles/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Get(
		url, &r.Body,
		&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return fmt.Errorf("error running api(read) for resource(IdentityRoleV3: %v), error: %s", d.Id(), r.Err)
	}

	v, err := common.NavigateValue(r.Body, []string{"role"}, nil)
	if err != nil {
		return err
	}
	result["read"] = v

	return nil
}

func setIdentityRoleV3Properties(d *schema.ResourceData, response map[string]interface{}) error {
	opts := resourceIdentityRoleV3UserInputParams(d)

	statementProp, _ := opts["statement"]
	statementProp, err := flattenIdentityRoleV3Statement(response, nil, statementProp)
	if err != nil {
		return fmt.Errorf("error reading Role:statement, err: %s", err)
	}
	if err = d.Set("statement", statementProp); err != nil {
		return fmt.Errorf("error setting Role:statement, err: %s", err)
	}

	catalogProp, err := common.NavigateValue(response, []string{"read", "catalog"}, nil)
	if err != nil {
		return fmt.Errorf("error reading Role:catalog, err: %s", err)
	}
	if err = d.Set("catalog", catalogProp); err != nil {
		return fmt.Errorf("error setting Role:catalog, err: %s", err)
	}

	descriptionProp, err := common.NavigateValue(response, []string{"read", "description"}, nil)
	if err != nil {
		return fmt.Errorf("error reading Role:description, err: %s", err)
	}
	if err = d.Set("description", descriptionProp); err != nil {
		return fmt.Errorf("error setting Role:description, err: %s", err)
	}

	displayLayerProp, _ := opts["display_layer"]
	displayLayerProp, err = flattenIdentityRoleV3DisplayLayer(response, nil, displayLayerProp)
	if err != nil {
		return fmt.Errorf("error reading Role:display_layer, err: %s", err)
	}
	if err = d.Set("display_layer", displayLayerProp); err != nil {
		return fmt.Errorf("error setting Role:display_layer, err: %s", err)
	}

	displayNameProp, err := common.NavigateValue(response, []string{"read", "display_name"}, nil)
	if err != nil {
		return fmt.Errorf("error reading Role:display_name, err: %s", err)
	}
	if err = d.Set("display_name", displayNameProp); err != nil {
		return fmt.Errorf("error setting Role:display_name, err: %s", err)
	}

	domainIDProp, err := common.NavigateValue(response, []string{"read", "domain_id"}, nil)
	if err != nil {
		return fmt.Errorf("error reading Role:domain_id, err: %s", err)
	}
	if err = d.Set("domain_id", domainIDProp); err != nil {
		return fmt.Errorf("error setting Role:domain_id, err: %s", err)
	}

	nameProp, err := common.NavigateValue(response, []string{"read", "name"}, nil)
	if err != nil {
		return fmt.Errorf("error reading Role:name, err: %s", err)
	}
	if err = d.Set("name", nameProp); err != nil {
		return fmt.Errorf("error setting Role:name, err: %s", err)
	}

	return nil
}

func flattenIdentityRoleV3Statement(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		v, err := common.NavigateValue(d, []string{"read", "policy", "Statement"}, arrayIndex)
		if err != nil {
			return nil, err
		}
		n := len(v.([]interface{}))
		result = make([]interface{}, n, n)
	}

	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	for i := 0; i < len(result); i++ {
		newArrayIndex["read.policy.Statement"] = i
		if result[i] == nil {
			result[i] = make(map[string]interface{})
		}
		r := result[i].(map[string]interface{})

		actionProp, err := common.NavigateValue(d, []string{"read", "policy", "Statement", "Action"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("error reading Role:action, err: %s", err)
		}
		r["action"] = actionProp

		effectProp, err := common.NavigateValue(d, []string{"read", "policy", "Statement", "Effect"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("error reading Role:effect, err: %s", err)
		}
		r["effect"] = effectProp
	}

	return result, nil
}

func flattenIdentityRoleV3DisplayLayer(d interface{}, arrayIndex map[string]int, _ interface{}) (interface{}, error) {
	v, err := common.NavigateValue(d, []string{"read", "type"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v == "AX" {
		return "domain", nil
	} else if v == "XA" {
		return "project", nil
	}
	return nil, fmt.Errorf("unknown display type:%v", v)
}
