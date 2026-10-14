package gpfs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gpfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceGPFSAccessRuleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGPFSAccessRuleV3Create,
		ReadContext:   resourceGPFSAccessRuleV3Read,
		DeleteContext: resourceGPFSAccessRuleV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"fs_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_source_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(gpfs.AccessRuleActionFullControl),
				ValidateFunc: validation.StringInSlice([]string{
					string(gpfs.AccessRuleActionFullControl),
					string(gpfs.AccessRuleActionRead),
				}, false),
			},
		},
	}
}

func resourceGPFSAccessRuleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GPFS client: %w", err)
	}

	fsName := d.Get("fs_name").(string)
	opts := gpfs.CreateFSAccessRulesInput{
		FSName: fsName,
		Statement: []gpfs.AccessRule{{
			Action: gpfs.AccessRuleAction(d.Get("action").(string)),
			Effect: gpfs.AccessRuleEffectAllow,
			Condition: gpfs.AccessRuleCondition{
				SourceVPC:   d.Get("vpc_id").(string),
				VPCSourceIP: common.ExpandToStringList(d.Get("vpc_source_ip").([]interface{})),
			},
		}},
	}
	if _, err = client.CreateFSAccessRules(&opts); err != nil {
		return diag.Errorf("error creating GPFS access rule for file system (%s): %s", fsName, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", fsName, d.Get("vpc_id").(string)))
	return resourceGPFSAccessRuleV3Read(ctx, d, meta)
}

func resourceGPFSAccessRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GPFS client: %w", err)
	}

	fsName, vpcId, err := parseGPFSAccessRuleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := getGPFSAccessRule(client, fsName, vpcId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving GPFS access rule (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("fs_name", fsName),
		d.Set("vpc_id", rule.Condition.SourceVPC),
		d.Set("vpc_source_ip", rule.Condition.VPCSourceIP),
		d.Set("action", string(rule.Action)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving GPFS access rule (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceGPFSAccessRuleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GPFS client: %w", err)
	}

	fsName, _, err := parseGPFSAccessRuleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err = client.DeleteFSAccessRules(fsName); err != nil {
		return diag.Errorf("error deleting GPFS access rules for file system (%s): %s", fsName, err)
	}
	return nil
}

func parseGPFSAccessRuleId(id string) (fsName, vpcId string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid format of GPFS access rule ID (%s), expected <fs_name>/<vpc_id>", id)
	}
	return parts[0], parts[1], nil
}

func getGPFSAccessRule(client *gpfs.Client, fsName, vpcId string) (*gpfs.AccessRule, error) {
	output, err := client.GetFSAccessRules(fsName)
	if err != nil {
		return nil, err
	}
	for i := range output.Statement {
		if output.Statement[i].Condition.SourceVPC == vpcId {
			return &output.Statement[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}
