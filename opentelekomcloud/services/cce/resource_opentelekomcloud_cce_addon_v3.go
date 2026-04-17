package cce

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/addons"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCEAddonV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCEAddonV3Create,
		ReadContext:   resourceCCEAddonV3Read,
		UpdateContext: resourceCCEAddonV3Update,
		DeleteContext: resourceCCEAddonV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceCCEAddonV3Import,
		},

		Schema: map[string]*schema.Schema{
			"template_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:         schema.TypeMap,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							ExactlyOneOf: []string{"values.0.basic", "values.0.basic_json"},
						},
						"basic_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := common.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ExactlyOneOf: []string{"values.0.basic", "values.0.basic_json"},
						},
						"custom": {
							Type:          schema.TypeMap,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							ConflictsWith: []string{"values.0.custom_json"},
						},
						"custom_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := common.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ConflictsWith: []string{"values.0.custom"},
						},
						"flavor": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: common.ValidateJsonString,
							StateFunc: func(v interface{}) string {
								jsonString, _ := common.NormalizeJsonString(v)
								return jsonString
							},
							ConflictsWith: []string{"values.0.flavor_json"},
						},
						"flavor_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := common.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ConflictsWith: []string{"values.0.flavor"},
						},
					},
				},
			},
		},
	}
}

func resourceCCEAddonV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientAddonV3, func() (*golangsdk.ServiceClient, error) {
		return config.CceV3AddonClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterID := d.Get("cluster_id").(string)
	basic, custom, flavor, err := buildAddonValues(d)
	if err != nil {
		return fmterr.Errorf("error getting values for CCE addon: %s", err)
	}

	templateName := d.Get("template_name").(string)
	addon, err := addons.Create(client, addons.CreateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.CreateMetadata{
			Annotations: addons.CreateAnnotations{
				AddonInstallType: "install",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("template_version").(string),
			ClusterID:         clusterID,
			AddonTemplateName: templateName,
			Values: addons.Values{
				Basic:    basic,
				Advanced: custom,
				Flavor:   flavor,
			},
		},
	}, clusterID)

	if err != nil {
		errMsg := logHttpError(err)
		addonSpec, aErr := getAddonTemplateSpec(client, clusterID, templateName)
		if aErr == nil {
			errMsg = fmt.Errorf("\nAddon template spec: %s\n%s", addonSpec, errMsg)
		}
		return fmterr.Errorf("error creating CCE addon instance: %s", errMsg)
	}

	d.SetId(addon.Metadata.Id)

	log.Printf("[DEBUG] Waiting for CCEAddon (%s) to become available", addon.Metadata.Id)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"installing", "abnormal"},
		Target:       []string{"running", "available", "abnormal"},
		Refresh:      waitForCCEAddonActive(client, addon.Metadata.Id, clusterID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("Error creating CCEAddon: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientAddonV3)
	return resourceCCEAddonV3Read(clientCtx, d, meta)
}

func resourceCCEAddonV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientAddonV3, func() (*golangsdk.ServiceClient, error) {
		return config.CceV3AddonClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterID := d.Get("cluster_id").(string)
	addon, err := addons.Get(client, d.Id(), clusterID)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error reading CCE addon instance: %s", logHttpError(err))
	}

	mErr := multierror.Append(nil,
		d.Set("name", addon.Metadata.Name),
		d.Set("cluster_id", addon.Spec.ClusterID),
		d.Set("template_version", addon.Spec.Version),
		d.Set("template_name", addon.Spec.AddonTemplateName),
		d.Set("description", addon.Spec.Description),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting addon attributes: %s", err)
	}

	return nil
}

func resourceCCEAddonV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientAddonV3, func() (*golangsdk.ServiceClient, error) {
		return config.CceV3AddonClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterID := d.Get("cluster_id").(string)
	addonID := d.Id()

	addon, err := addons.Get(client, addonID, clusterID)
	if err != nil {
		return diag.Errorf("error reading current CCE add-on (%s) before update: %s", addonID, logHttpError(err))
	}

	desiredBasic, desiredCustom, desiredFlavor, err := buildAddonValues(d)
	if err != nil {
		return diag.Errorf("error getting values for CCE add-on: %s", err)
	}

	updateValues := addon.Spec.Values
	if addonValuesFieldChanged(d, "basic", "basic_json") {
		updateValues.Basic = desiredBasic
	}
	if addonValuesFieldChanged(d, "custom", "custom_json") {
		updateValues.Advanced = desiredCustom
	}
	if addonValuesFieldChanged(d, "flavor", "flavor_json") {
		updateValues.Flavor = desiredFlavor
	}

	updateOpts := addons.UpdateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.UpdateMetadata{
			Annotations: addons.UpdateAnnotations{
				AddonUpdateType: "upgrade",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("template_version").(string),
			ClusterID:         clusterID,
			AddonTemplateName: d.Get("template_name").(string),
			Values:            updateValues,
		},
	}

	_, err = addons.Update(client, addonID, clusterID, updateOpts)
	if err != nil {
		return diag.Errorf("error updating CCE add-on (%s): %s", addonID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"installing", "upgrading", "abnormal"},
		Target:       []string{"running", "available", "abnormal"},
		Refresh:      waitForCCEAddonActive(client, addonID, clusterID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE add-on (%s) to become available: %s", addonID, err)
	}

	return resourceCCEAddonV3Read(ctx, d, meta)
}

func addonValuesFieldChanged(d *schema.ResourceData, plainKey, jsonKey string) bool {
	return d.HasChange(fmt.Sprintf("values.0.%s", plainKey)) ||
		d.HasChange(fmt.Sprintf("values.0.%s", jsonKey))
}

func buildAddonValues(d *schema.ResourceData) (basic, custom, flavor map[string]interface{}, err error) {
	values := d.Get("values").([]interface{})
	if len(values) == 0 || values[0] == nil {
		basic = make(map[string]interface{})
		return
	}

	valuesMap := values[0].(map[string]interface{})
	if basicRaw := valuesMap["basic"].(map[string]interface{}); len(basicRaw) != 0 {
		basic = basicRaw
	} else if basicJsonRaw := valuesMap["basic_json"].(string); basicJsonRaw != "" {
		err = json.Unmarshal([]byte(basicJsonRaw), &basic)
		if err != nil {
			err = fmt.Errorf("error unmarshalling basic json: %s", err)
			return
		}
	}

	if customRaw := valuesMap["custom"].(map[string]interface{}); len(customRaw) != 0 {
		custom = customRaw
	} else if customJsonRaw := valuesMap["custom_json"].(string); customJsonRaw != "" {
		err = json.Unmarshal([]byte(customJsonRaw), &custom)
		if err != nil {
			err = fmt.Errorf("error unmarshalling custom json: %s", err)
			return
		}
	}

	if flavorRaw := valuesMap["flavor"].(string); flavorRaw != "" {
		err = json.Unmarshal([]byte(flavorRaw), &flavor)
		if err != nil {
			err = fmt.Errorf("error unmarshalling flavor json %s", err)
			return
		}
	} else if flavorJsonRaw := valuesMap["flavor_json"].(string); flavorJsonRaw != "" {
		err = json.Unmarshal([]byte(flavorJsonRaw), &flavor)
		if err != nil {
			err = fmt.Errorf("error unmarshalling flavor json %s", err)
			return
		}
	}

	return
}

func resourceCCEAddonV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientAddonV3, func() (*golangsdk.ServiceClient, error) {
		return config.CceV3AddonClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterID := d.Get("cluster_id").(string)

	if err := addons.Delete(client, d.Id(), clusterID); err != nil {
		return fmterr.Errorf("error deleting CCE addon : %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"available"},
		Target:       []string{"deleted"},
		Refresh:      waitForCCEAddonDelete(client, d.Id(), clusterID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getAddonTemplateSpec(client *golangsdk.ServiceClient, clusterID, templateName string) (string, error) {
	templates, err := addons.ListTemplates(client, clusterID, addons.ListOpts{Name: templateName})
	if err != nil {
		return "", err
	}
	jsonTemplate, _ := json.Marshal(templates)
	return string(jsonTemplate), nil
}

func logHttpError(err error) error {
	if httpErr, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
		return fmt.Errorf("response: %s\n %s", httpErr.Error(), httpErr.Body)
	}
	return err
}

func waitForCCEAddonDelete(client *golangsdk.ServiceClient, addonID, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		addon, err := addons.Get(client, addonID, clusterID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return addon, "deleted", nil
			}
			return nil, "error", fmt.Errorf("error waiting CCE addon to become deleted: %s", err)
		}

		return addon, "available", nil
	}
}

func resourceCCEAddonV3Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for CCE Addon. Format must be <cluster id>/<addon id>")
		return nil, err
	}
	clusterID := parts[0]
	addonID := parts[1]
	d.SetId(addonID)

	config := meta.(*cfg.Config)
	client, err := config.CceV3AddonClient(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CCE client: %s", logHttpError(err))
	}

	addon, err := addons.Get(client, d.Id(), clusterID)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil, fmt.Errorf("addon wasn't found")
		}

		return nil, fmt.Errorf("error reading CCE addon instance: %s", logHttpError(err))
	}

	mErr := multierror.Append(nil,
		d.Set("name", addon.Metadata.Name),
		d.Set("cluster_id", addon.Spec.ClusterID),
		d.Set("template_version", addon.Spec.Version),
		d.Set("template_name", addon.Spec.AddonTemplateName),
		d.Set("description", addon.Spec.Description),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("error setting addon attributes: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func waitForCCEAddonActive(cceAddonClient *golangsdk.ServiceClient, id, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := addons.Get(cceAddonClient, id, clusterID)
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Status, nil
	}
}
