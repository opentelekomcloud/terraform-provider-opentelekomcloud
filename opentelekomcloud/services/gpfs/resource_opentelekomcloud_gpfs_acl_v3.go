package gpfs

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gpfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceGPFSACLV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGPFSACLV3Write,
		ReadContext:   resourceGPFSACLV3Read,
		UpdateContext: resourceGPFSACLV3Write,
		DeleteContext: resourceGPFSACLV3Delete,
		CustomizeDiff: validateGPFSACLStatements,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"fs_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"statement": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(gpfs.AccessRuleActionFullControl),
							ValidateFunc: validation.StringInSlice([]string{
								string(gpfs.AccessRuleActionFullControl),
								string(gpfs.AccessRuleActionRead),
							}, false),
						},
						"effect": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(gpfs.AccessRuleEffectAllow),
							ValidateFunc: validation.StringInSlice([]string{string(gpfs.AccessRuleEffectAllow)}, false),
						},
						"vpc_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"vpc_source_ip": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func validateGPFSACLStatements(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	config := d.GetRawConfig()
	if config.IsNull() || !config.IsKnown() {
		return nil
	}
	statements := config.GetAttr("statement")
	if statements.IsNull() {
		return fmt.Errorf("GPFS ACL requires at least one statement")
	}
	if !statements.IsKnown() {
		return nil
	}
	seen := make(map[string]bool)
	for _, statement := range statements.AsValueSlice() {
		if statement.IsNull() || !statement.IsKnown() {
			continue
		}
		vpc := statement.GetAttr("vpc_id")
		if vpc.IsNull() || !vpc.IsKnown() {
			continue
		}
		if seen[vpc.AsString()] {
			return fmt.Errorf("GPFS ACL must contain only one statement per VPC (%s)", vpc.AsString())
		}
		seen[vpc.AsString()] = true
	}
	if !statements.IsWhollyKnown() {
		return nil
	}
	old, planned := d.GetChange("statement")
	if len(old.([]interface{})) == 0 {
		return nil
	}
	previous := make(map[string]map[string]interface{})
	for _, item := range old.([]interface{}) {
		value := item.(map[string]interface{})
		previous[value["vpc_id"].(string)] = value
	}
	values := planned.([]interface{})
	for i, configured := range statements.AsValueSlice() {
		value := values[i].(map[string]interface{})
		prior, exists := previous[value["vpc_id"].(string)]
		// Optional/computed metadata belongs to the VPC, not its current list index.
		for _, key := range []string{"sid", "vpc_source_ip"} {
			if configured.GetAttr(key).IsNull() {
				switch {
				case exists:
					value[key] = prior[key]
				case key == "sid":
					value[key] = ""
				default:
					value[key] = []interface{}{}
				}
			}
		}
	}
	return d.SetNew("statement", values)
}

func expandGPFSACLStatements(raw []interface{}) []gpfs.AccessRule {
	statements := make([]gpfs.AccessRule, 0, len(raw))
	for _, item := range raw {
		value := item.(map[string]interface{})
		statements = append(statements, gpfs.AccessRule{
			ID:     value["sid"].(string),
			Action: gpfs.AccessRuleAction(value["action"].(string)),
			Effect: gpfs.AccessRuleEffect(value["effect"].(string)),
			Condition: gpfs.AccessRuleCondition{
				SourceVPC:   value["vpc_id"].(string),
				VPCSourceIP: common.ExpandToStringList(value["vpc_source_ip"].([]interface{})),
			},
		})
	}
	return statements
}

func flattenGPFSACLStatements(statements []gpfs.AccessRule, previous []interface{}) []interface{} {
	ordered := make([]gpfs.AccessRule, 0, len(statements))
	used := make([]bool, len(statements))
	// Keep Terraform's statement order even when the API reorders the ACL.
	for _, item := range previous {
		value := item.(map[string]interface{})
		for i, statement := range statements {
			if !used[i] && statement.Condition.SourceVPC == value["vpc_id"] {
				ordered = append(ordered, statement)
				used[i] = true
				break
			}
		}
	}
	for i, statement := range statements {
		if !used[i] {
			ordered = append(ordered, statement)
		}
	}
	result := make([]interface{}, 0, len(ordered))
	for _, statement := range ordered {
		result = append(result, map[string]interface{}{
			"sid":           statement.ID,
			"action":        string(statement.Action),
			"effect":        string(statement.Effect),
			"vpc_id":        statement.Condition.SourceVPC,
			"vpc_source_ip": statement.Condition.VPCSourceIP,
		})
	}
	return result
}

func gpfsACLMatches(actual, desired []gpfs.AccessRule) bool {
	if len(actual) != len(desired) {
		return false
	}
	used := make([]bool, len(actual))
	for _, expected := range desired {
		found := false
		for i, statement := range actual {
			if !used[i] && statement.Condition.SourceVPC == expected.Condition.SourceVPC &&
				statement.Action == expected.Action && statement.Effect == expected.Effect &&
				(expected.ID == "" || statement.ID == expected.ID) {
				used[i], found = true, true
				break
			}
		}
		if !found {
			return false
		}
	}
	// VpcSourceIp is documented but currently unsupported by the service.
	return true
}

func setGPFSACLState(d *schema.ResourceData, output *gpfs.GetFSAccessRulesOutput) diag.Diagnostics {
	seen := make(map[string]bool)
	for _, statement := range output.Statement {
		if seen[statement.Condition.SourceVPC] {
			return diag.Errorf("GPFS ACL (%s) contains multiple statements for VPC (%s); this resource supports one statement per VPC",
				d.Id(), statement.Condition.SourceVPC)
		}
		seen[statement.Condition.SourceVPC] = true
	}
	mErr := multierror.Append(nil,
		d.Set("fs_name", d.Id()),
		d.Set("statement", flattenGPFSACLStatements(output.Statement, d.Get("statement").([]interface{}))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving GPFS ACL (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceGPFSACLV3Write(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GPFS client: %s", err)
	}
	timeout := d.Timeout(schema.TimeoutUpdate)
	if d.Id() == "" {
		timeout = d.Timeout(schema.TimeoutCreate)
	}
	fsName := d.Get("fs_name").(string)
	desired := expandGPFSACLStatements(d.Get("statement").([]interface{}))
	if len(desired) == 0 {
		return diag.Errorf("GPFS ACL (%s) requires at least one statement", fsName)
	}
	if _, err = client.CreateFSAccessRules(&gpfs.CreateFSAccessRulesInput{
		FSName: fsName, Statement: desired,
	}); err != nil {
		return diag.Errorf("error configuring GPFS ACL (%s): %s", fsName, err)
	}
	d.SetId(fsName)
	var output *gpfs.GetFSAccessRulesOutput
	wait := resource.StateChangeConf{
		Pending: []string{"pending"}, Target: []string{"ready"},
		Timeout: timeout, MinTimeout: time.Second,
		Refresh: func() (interface{}, string, error) {
			var readErr error
			output, readErr = client.GetFSAccessRules(fsName)
			if isGPFSService404(readErr) {
				return nil, "pending", nil
			}
			if readErr != nil {
				return nil, "", readErr
			}
			if gpfsACLMatches(output.Statement, desired) {
				return output, "ready", nil
			}
			return output, "pending", nil
		},
	}
	if _, err = wait.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for GPFS ACL (%s) to be configured: %s", fsName, err)
	}
	return setGPFSACLState(d, output)
}

func resourceGPFSACLV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GPFS client: %s", err)
	}
	output, err := client.GetFSAccessRules(d.Id())
	if isGPFSService404(err) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error retrieving GPFS ACL (%s): %s", d.Id(), err)
	}
	if len(output.Statement) == 0 {
		d.SetId("")
		return nil
	}
	return setGPFSACLState(d, output)
}

func resourceGPFSACLV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GPFS client: %s", err)
	}
	if _, err = client.DeleteFSAccessRules(d.Id()); err != nil {
		if isGPFSService404(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error deleting GPFS ACL (%s): %s", d.Id(), err)
	}
	wait := resource.StateChangeConf{
		Pending: []string{"pending"}, Target: []string{"deleted"},
		Timeout: d.Timeout(schema.TimeoutDelete), MinTimeout: time.Second,
		Refresh: func() (interface{}, string, error) {
			output, readErr := client.GetFSAccessRules(d.Id())
			if isGPFSService404(readErr) {
				return &gpfs.GetFSAccessRulesOutput{}, "deleted", nil
			}
			if readErr != nil {
				return nil, "", readErr
			}
			if len(output.Statement) == 0 {
				return output, "deleted", nil
			}
			return output, "pending", nil
		},
	}
	if _, err = wait.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for GPFS ACL (%s) to be deleted: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
