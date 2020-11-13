package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/addons"
)

func resourceCCEAddonV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEAddonV3Create,
		Read:   resourceCCEAddonV3Read,
		Delete: resourceCCEAddonV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

func getValuesValues(d *schema.ResourceData) (basic, custom map[string]interface{}, err error) {
	values := d.Get("values").([]interface{})
	if len(values) == 0 {
		err = fmt.Errorf("no values are set for CCE addon") // should be impossible, as Required: true
		return
	}
	valuesMap := values[0].(map[string]interface{})

	basicRaw, ok := valuesMap["basic"]
	if !ok {
		err = fmt.Errorf("no basic values are set for CCE addon") // should be impossible, as Required: true
		return
	}
	if customRaw, ok := valuesMap["custom"]; ok {
		custom = customRaw.(map[string]interface{})
	}
	basic = basicRaw.(map[string]interface{})
	return
}

func resourceCCEAddonV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	basic, custom, err := getValuesValues(d)
	if err != nil {
		return fmt.Errorf("error getting values for CCE addon: %s", err)
	}

	addon, err := addons.Create(client, addons.CreateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.CreateMetadata{
			Annotations: addons.Annotations{
				AddonInstallType: "install",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("template_version").(string),
			ClusterID:         clusterID,
			AddonTemplateName: d.Get("template_name").(string),
			Values: addons.Values{
				Basic:    basic,
				Advanced: custom,
			},
		},
	}, clusterID).Extract()

	if err != nil {
		return fmt.Errorf("error creating CCE addon instance: %s", err)
	}

	d.SetId(addon.Metadata.Id)

	return resourceCCEAddonV3Read(d, meta)
}
func resourceCCEAddonV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating CCE client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	addon, err := addons.Get(client, d.Id(), clusterID).Extract()
	if err != nil {
		return fmt.Errorf("error reading CCE addon instance: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", addon.Metadata.Name),
		d.Set("template_version", addon.Spec.Version),
		d.Set("template_name", addon.Spec.AddonTemplateName),
		d.Set("values", addon.Spec.Values),
		d.Set("description", addon.Spec.Description),
	)

	return mErr.ErrorOrNil()
}
func resourceCCEAddonV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cceV3Client(GetRegion(d, config))
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
