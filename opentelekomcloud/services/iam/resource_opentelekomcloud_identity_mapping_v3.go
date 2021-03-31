package iam

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/mappings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

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
				Type:     schema.TypeList,
				Required: true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"local": {
							Type:     schema.TypeList,
							Required: true,
							Elem: schema.Resource{
								Schema: map[string]*schema.Schema{
									"user": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"group": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"groups": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"remote": {
							Type:     schema.TypeList,
							Required: true,
							Elem: schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"not_any_of": {
										Type:     schema.TypeSet,
										Optional: true,
									},
									"any_one_of": {
										Type:     schema.TypeSet,
										Optional: true,
									},
									"regex": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceIdentityMappingV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}

	rulesRaw := d.Get("rules").([]map[string]interface{})
	var rulesList []mappings.RuleOpts
	for _, ruleRaw := range rulesRaw {
		var localRulesList []mappings.LocalRuleOpts
		localRuleRaw := ruleRaw["local"].([]map[string]interface{})
		for _, local := range localRuleRaw {
			var user *mappings.UserOpts
			var group *mappings.GroupOpts
			var groups string
			if v, found := local["user"]; found {
				user = &mappings.UserOpts{
					Name: v.(map[string]interface{})["name"].(string),
				}
			}
			if v, found := local["group"]; found {
				group = &mappings.GroupOpts{
					Name: v.(map[string]interface{})["name"].(string),
				}
			}
			if v, found := local["groups"]; found {
				groups = v.(string)
			}
			localRule := mappings.LocalRuleOpts{
				User:   user,
				Group:  group,
				Groups: groups,
			}
			localRulesList = append(localRulesList, localRule)
		}

		var remoteRulesList []mappings.RemoteRuleOpts
		remoteRuleRaw := ruleRaw["remote"].([]map[string]interface{})
		for _, remote := range remoteRuleRaw {
			var notAnyOf []string
			var anyOneOf []string
			var regex bool

			if v, found := remote["not_any_of"]; found {
				notAnyOf = v.([]string)
			}
			if v, found := remote["any_one_of"]; found {
				anyOneOf = v.([]string)
			}
			if v, found := remote["regex"]; found {
				regex = v.(bool)
			}

			remoteRule := mappings.RemoteRuleOpts{
				Type:     remote["type"].(string),
				NotAnyOf: notAnyOf,
				AnyOneOf: anyOneOf,
				Regex:    &regex,
			}
			remoteRulesList = append(remoteRulesList, remoteRule)
		}
		rule := mappings.RuleOpts{
			Local:  localRulesList,
			Remote: remoteRulesList,
		}
		rulesList = append(rulesList, rule)
	}
	createOpts := mappings.CreateOpts{
		Rules: rulesList,
	}
	mappingID := d.Get("mapping_id").(string)
	mapping, err := mappings.Create(client, mappingID, createOpts).Extract()
	if err != nil {
		return err
	}

	d.SetId(mapping.ID)

	return resourceIdentityMappingV3Read(d, meta)
}

func resourceIdentityMappingV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}

	return nil
}

func resourceIdentityMappingV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}

	return resourceIdentityMappingV3Read(d, meta)
}

func resourceIdentityMappingV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}

	if err := mappings.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud identity mapping: %s", err)
	}

	return nil
}
