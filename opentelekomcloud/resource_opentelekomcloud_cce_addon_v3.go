package opentelekomcloud

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/addons"
)

func resourceCCEAddonV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEAddonV3Create,
		Read:   resourceCCEAddonV3Read,
		Update: resourceCCEAddonV3Update,
		Delete: resourceCCEAddonV3Delete,

		Schema: map[string]*schema.Schema{
			"template_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
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
							Type:     schema.TypeMap,
							Required: true,
						},
						"custom": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceCCEAddonV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3AddonClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	basic, custom, err := getAddonValues(d)
	if err != nil {
		return fmt.Errorf("error getting values for CCE addon: %s", err)
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
			},
		},
	}, clusterID).Extract()

	if err != nil {
		errMsg := logHttpError(err)
		addonSpec, aErr := getAddonTemplateSpec(client, clusterID, templateName)
		if aErr == nil {
			errMsg = fmt.Errorf("\nAddon template spec: %s\n%s", addonSpec, errMsg)
		}
		return fmt.Errorf("error creating CCE addon instance: %s", errMsg)
	}

	d.SetId(addon.Metadata.Id)

	return resourceCCEAddonV3Read(d, meta)
}
func resourceCCEAddonV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3AddonClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", logHttpError(err))
	}

	clusterID := d.Get("cluster_id").(string)
	addon, err := addons.Get(client, d.Id(), clusterID).Extract()
	if err != nil {
		return fmt.Errorf("error reading CCE addon instance: %s", logHttpError(err))
	}

	mErr := multierror.Append(nil,
		d.Set("name", addon.Metadata.Name),
		d.Set("template_version", addon.Spec.Version),
		d.Set("template_name", addon.Spec.AddonTemplateName),
		d.Set("description", addon.Spec.Description),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting addon attributes: %s", err)
	}

	return nil
}

func getAddonValues(d *schema.ResourceData) (basic, custom map[string]interface{}, err error) {
	valLength := d.Get("values.#").(int)
	if valLength == 0 {
		err = fmt.Errorf("no values are set for CCE addon")
		return
	}
	basic = d.Get("values.0.basic").(map[string]interface{})
	custom = d.Get("values.0.custom").(map[string]interface{})
	return
}

func resourceCCEAddonV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3AddonClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", err)
	}
	clusterID := d.Get("cluster_id").(string)
	basic, custom, err := getAddonValues(d)
	if err != nil {
		return fmt.Errorf("error getting values for CCE addon: %s", err)
	}

	templateVersion := d.Get("template_version").(string)
	templateName := d.Get("template_name").(string)

	_, err = addons.Update(client, d.Id(), clusterID, addons.UpdateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.UpdateMetadata{
			Annotations: addons.UpdateAnnotations{
				AddonUpdateType: "upgrade",
			},
		},
		Spec: addons.RequestSpec{
			Version:           templateVersion,
			ClusterID:         clusterID,
			AddonTemplateName: templateName,
			Values: addons.Values{
				Basic:    basic,
				Advanced: custom,
			},
		},
	}).Extract()

	if err != nil {
		errMsg := logHttpError(err)
		addonSpec, aErr := getAddonTemplateSpec(client, clusterID, templateName)
		if aErr == nil {
			errMsg = fmt.Errorf("\nSomething got wrong installing CCE addon\nAddon template spec:\n%s\nError: %s", addonSpec, errMsg)
		}
		return fmt.Errorf("error updating CCE addon instance: %s", errMsg)
	}

	return resourceCCEAddonV3Read(d, meta)
}

func resourceCCEAddonV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3AddonClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", err)
	}
	clusterID := d.Get("cluster_id").(string)
	err = addons.Delete(client, d.Id(), clusterID).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting addon: %s", err)
	}
	d.SetId("")
	return nil
}

func getAddonTemplateSpec(client *golangsdk.ServiceClient, clusterID, templateName string) (string, error) {
	templates, err := addons.ListTemplates(client, clusterID, addons.ListOpts{Name: templateName}).Extract()
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
