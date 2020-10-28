package opentelekomcloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/secgroups"
)

func resourceComputeSecGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSecGroupV2Create,
		Read:   resourceComputeSecGroupV2Read,
		Update: resourceComputeSecGroupV2Update,
		Delete: resourceComputeSecGroupV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							StateFunc: func(v interface{}) string {
								return strings.ToLower(v.(string))
							},
						},
						"from_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"self": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: false,
						},
					},
				},
				Set: secgroupRuleV2Hash,
			},
		},
	}
}

func resourceComputeSecGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	// Before creating the security group, make sure all rules are valid.
	if err := checkSecGroupV2RulesForErrors(d); err != nil {
		return err
	}

	// If all rules are valid, proceed with creating the security gruop.
	createOpts := secgroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	sg, err := secgroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud security group: %s", err)
	}

	d.SetId(sg.ID)

	// Now that the security group has been created, iterate through each rule and create it
	createRuleOptsList := resourceSecGroupRulesV2(d)
	for _, createRuleOpts := range createRuleOptsList {
		_, err := secgroups.CreateRule(computeClient, createRuleOpts).Extract()
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud security group rule: %s", err)
		}
	}

	return resourceComputeSecGroupV2Read(d, meta)
}

func resourceComputeSecGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	sg, err := secgroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "security group")
	}
	me := multierror.Append(nil,
		d.Set("name", sg.Name),
		d.Set("description", sg.Description),
		d.Set("region", GetRegion(d, config)),
	)

	rulesMap, err := rulesToMap(computeClient, d, sg.Rules)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] rulesToMap(sg.Rules): %+v", rulesMap)
	if err := d.Set("rule", rulesMap); err != nil {
		return fmt.Errorf("[DEBUG] Error saving rule to state for OpenTelekomCloud server (%s): %s", d.Id(), err)
	}

	return me.ErrorOrNil()
}

func resourceComputeSecGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	updateOpts := secgroups.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] Updating Security Group (%s) with options: %+v", d.Id(), updateOpts)

	_, err = secgroups.Update(computeClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating OpenTelekomCloud security group (%s): %s", d.Id(), err)
	}

	if d.HasChange("rule") {
		oldSGRaw, newSGRaw := d.GetChange("rule")
		oldSGRSet, newSGRSet := oldSGRaw.(*schema.Set), newSGRaw.(*schema.Set)
		secGroupRulesToAdd := newSGRSet.Difference(oldSGRSet)
		secGroupRulesToRemove := oldSGRSet.Difference(newSGRSet)

		log.Printf("[DEBUG] Security group rules to add: %v", secGroupRulesToAdd)
		log.Printf("[DEBUG] Security groups rules to remove: %v", secGroupRulesToRemove)

		for _, rawRule := range secGroupRulesToAdd.List() {
			createRuleOpts := resourceSecGroupRuleCreateOptsV2(d, rawRule)
			rule, err := secgroups.CreateRule(computeClient, createRuleOpts).Extract()
			if err != nil {
				return fmt.Errorf("error adding rule to OpenTelekomCloud security group (%s): %s", d.Id(), err)
			}
			log.Printf("[DEBUG] Added rule (%s) to OpenTelekomCloud security group (%s) ", rule.ID, d.Id())
		}

		for _, r := range secGroupRulesToRemove.List() {
			rule := resourceSecGroupRuleV2(d, r)
			err := secgroups.DeleteRule(computeClient, rule.ID).ExtractErr()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}
				return fmt.Errorf("error removing rule (%s) from OpenTelekomCloud security group (%s)", rule.ID, d.Id())
			} else {
				log.Printf("[DEBUG] Removed rule (%s) from OpenTelekomCloud security group (%s): %s", rule.ID, d.Id(), err)
			}
		}
	}

	return resourceComputeSecGroupV2Read(d, meta)
}

func resourceComputeSecGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    SecGroupV2StateRefreshFunc(computeClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud security group: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceSecGroupRulesV2(d *schema.ResourceData) []secgroups.CreateRuleOpts {
	rawRules := d.Get("rule").(*schema.Set).List()
	createRuleOptsList := make([]secgroups.CreateRuleOpts, len(rawRules))
	for i, rawRule := range rawRules {
		createRuleOptsList[i] = resourceSecGroupRuleCreateOptsV2(d, rawRule)
	}
	return createRuleOptsList
}

func resourceSecGroupRuleCreateOptsV2(d *schema.ResourceData, rawRule interface{}) secgroups.CreateRuleOpts {
	rawRuleMap := rawRule.(map[string]interface{})
	groupId := rawRuleMap["from_group_id"].(string)
	if rawRuleMap["self"].(bool) {
		groupId = d.Id()
	}
	return secgroups.CreateRuleOpts{
		ParentGroupID: d.Id(),
		FromPort:      rawRuleMap["from_port"].(int),
		ToPort:        rawRuleMap["to_port"].(int),
		IPProtocol:    rawRuleMap["ip_protocol"].(string),
		CIDR:          rawRuleMap["cidr"].(string),
		FromGroupID:   groupId,
	}
}

func checkSecGroupV2RulesForErrors(d *schema.ResourceData) error {
	rawRules := d.Get("rule").(*schema.Set).List()
	for _, rawRule := range rawRules {
		rawRuleMap := rawRule.(map[string]interface{})

		// only one of cidr, from_group_id, or self can be set
		cidr := rawRuleMap["cidr"].(string)
		groupId := rawRuleMap["from_group_id"].(string)
		self := rawRuleMap["self"].(bool)
		errorMessage := fmt.Errorf("only one of cidr, from_group_id, or self can be set")

		// if cidr is set, from_group_id and self cannot be set
		if cidr != "" {
			if groupId != "" || self {
				return errorMessage
			}
		}

		// if from_group_id is set, cidr and self cannot be set
		if groupId != "" {
			if cidr != "" || self {
				return errorMessage
			}
		}

		// if self is set, cidr and from_group_id cannot be set
		if self {
			if cidr != "" || groupId != "" {
				return errorMessage
			}
		}
	}

	return nil
}

func resourceSecGroupRuleV2(d *schema.ResourceData, rawRule interface{}) secgroups.Rule {
	rawRuleMap := rawRule.(map[string]interface{})
	return secgroups.Rule{
		ID:            rawRuleMap["id"].(string),
		ParentGroupID: d.Id(),
		FromPort:      rawRuleMap["from_port"].(int),
		ToPort:        rawRuleMap["to_port"].(int),
		IPProtocol:    rawRuleMap["ip_protocol"].(string),
		IPRange:       secgroups.IPRange{CIDR: rawRuleMap["cidr"].(string)},
	}
}

func rulesToMap(computeClient *golangsdk.ServiceClient, d *schema.ResourceData, sgrs []secgroups.Rule) ([]map[string]interface{}, error) {
	sgrMap := make([]map[string]interface{}, len(sgrs))
	for i, sgr := range sgrs {
		groupId := ""
		self := false
		if sgr.Group.Name != "" {
			if sgr.Group.Name == d.Get("name").(string) {
				self = true
			} else {
				// Since Nova only returns the secgroup Name (and not the ID) for the group attribute,
				// we need to look up all security groups and match the name.
				// Nevermind that Nova wants the ID when setting the Group *and* that multiple groups
				// with the same name can exist...
				allPages, err := secgroups.List(computeClient).AllPages()
				if err != nil {
					return nil, err
				}
				securityGroups, err := secgroups.ExtractSecurityGroups(allPages)
				if err != nil {
					return nil, err
				}

				for _, sg := range securityGroups {
					if sg.Name == sgr.Group.Name {
						groupId = sg.ID
					}
				}
			}
		}

		sgrMap[i] = map[string]interface{}{
			"id":            sgr.ID,
			"from_port":     sgr.FromPort,
			"to_port":       sgr.ToPort,
			"ip_protocol":   sgr.IPProtocol,
			"cidr":          sgr.IPRange.CIDR,
			"self":          self,
			"from_group_id": groupId,
		}
	}
	return sgrMap, nil
}

func secgroupRuleV2Hash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["from_port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["to_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m["ip_protocol"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["cidr"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["from_group_id"].(string)))
	buf.WriteString(fmt.Sprintf("%t-", m["self"].(bool)))

	return hashcode.String(buf.String())
}

func SecGroupV2StateRefreshFunc(computeClient *golangsdk.ServiceClient, secGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete Security Group %s.\n", secGroupId)

		s, err := secgroups.Get(computeClient, secGroupId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Neutron Security Group %s", secGroupId)
				return s, "DELETED", nil
			}
			return s, "ACTIVE", err
		}

		err = secgroups.Delete(computeClient, secGroupId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Security Group %s", secGroupId)
				return s, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return s, "ACTIVE", nil
			}
			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] Security Group %s still active.\n", secGroupId)
		return s, "ACTIVE", nil
	}
}
