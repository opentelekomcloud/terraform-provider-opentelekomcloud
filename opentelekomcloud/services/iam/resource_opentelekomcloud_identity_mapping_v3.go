package iam

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/mappings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const mappingError = "error %s identity mapping v3: %w"

func ResourceIdentityMappingV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityMappingV3Create,
		Read:   resourceIdentityMappingV3Read,
		Update: resourceIdentityMappingV3Update,
		Delete: resourceIdentityMappingV3Delete,

		Schema: map[string]*schema.Schema{
			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rules": {
				Type:         schema.TypeList,
				Required:     true,
				ValidateFunc: common.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					jsonString, _ := common.NormalizeJsonString(v)
					return jsonString
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIdentityMappingV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	rulesRaw := d.Get("rules")
	rulesBytes, err := json.Marshal(rulesRaw)
	if err != nil {
		return err
	}
	rules := make([]mappings.RuleOpts, 1)
	if err := json.Unmarshal(rulesBytes, &rules); err != nil {
		return err
	}

	createOpts := mappings.CreateOpts{
		Rules: rules,
	}
	mappingID := d.Get("mapping_id").(string)
	mapping, err := mappings.Create(client, mappingID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf(mappingError, "creating", err)
	}

	d.SetId(mapping.ID)

	return resourceIdentityMappingV3Read(d, meta)
}

func resourceIdentityMappingV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	mapping, err := mappings.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(mappingError, "reading", err)
	}

	rules, err := common.NormalizeJsonString(mapping.Rules)
	if err != nil {
		return err
	}
	if err := d.Set("rules", rules); err != nil {
		return err
	}

	if err := d.Set("links", mapping.Links); err != nil {
		return fmt.Errorf("error setting identity mapping links: %w", err)
	}

	return nil
}

func resourceIdentityMappingV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}
	changes := false
	updateOpts := mappings.UpdateOpts{}

	if d.HasChange("rules") {
		changes = true
		rulesRaw := d.Get("rules")
		rulesBytes, err := json.Marshal(rulesRaw)
		if err != nil {
			return err
		}
		rules := make([]mappings.RuleOpts, 1)
		if err := json.Unmarshal(rulesBytes, &rules); err != nil {
			return err
		}
		updateOpts.Rules = rules
	}
	if changes {
		_, err := mappings.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return err
		}
	}

	return resourceIdentityMappingV3Read(d, meta)
}

func resourceIdentityMappingV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	if err := mappings.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf(mappingError, "deleting", err)
	}

	return nil
}
