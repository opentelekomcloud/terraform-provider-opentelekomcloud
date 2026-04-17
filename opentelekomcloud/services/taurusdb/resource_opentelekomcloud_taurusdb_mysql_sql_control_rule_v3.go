package taurusdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/sqlfilter"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceTaurusDbV3MySQLSqlControlRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDbV3MySQLSqlControlRuleCreate,
		ReadContext:   resourceTaurusDbV3MySQLSqlControlRuleRead,
		UpdateContext: resourceTaurusDbV3MySQLSqlControlRuleUpdate,
		DeleteContext: resourceTaurusDbV3MySQLSqlControlRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sql_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pattern": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTaurusDbV3MySQLSqlControlRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceID := d.Get("instance_id").(string)
	nodeID := d.Get("node_id").(string)
	sqlType := d.Get("sql_type").(string)
	pattern := d.Get("pattern").(string)
	maxConcurrency := d.Get("max_concurrency").(int)

	opts := sqlfilter.UpdateSqlFilterRulesOpts{
		SqlFilterRules: []sqlfilter.NodeSqlFilterRuleInfo{
			{
				NodeId: nodeID,
				Rules: []sqlfilter.NodeSqlFilterRule{
					{
						SqlType: sqlType,
						Patterns: []sqlfilter.NodeSqlFilterRulePattern{
							{
								Pattern:        pattern,
								MaxConcurrency: maxConcurrency,
							},
						},
					},
				},
			},
		},
	}

	jobID, err := sqlfilter.UpdateSqlFilterRules(client, instanceID, opts)
	if err != nil {
		return diag.Errorf("error creating TaurusDB MySQL SQL control rule: %s", err)
	}

	if jobID == nil || *jobID == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	if err := waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", instanceID, nodeID, sqlType, pattern))

	return resourceTaurusDbV3MySQLSqlControlRuleRead(ctx, d, meta)
}

func resourceTaurusDbV3MySQLSqlControlRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	parts := strings.SplitN(d.Id(), "/", 4)
	if len(parts) != 4 {
		return diag.Errorf("invalid id format, must be <instance_id>/<node_id>/<sql_type>/<pattern>")
	}

	instanceID := parts[0]
	nodeID := parts[1]
	sqlType := parts[2]
	pattern := parts[3]

	opts := sqlfilter.GetSqlFilterRulesOpts{
		NodeId: nodeID,
	}

	resp, err := sqlfilter.GetSqlFilterRules(client, instanceID, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving TaurusDb MySQL SQL control rule")
	}

	var maxConcurrency *int
	for _, rule := range resp.SqlFilterRules {
		if rule.SqlType == sqlType {
			for _, p := range rule.Patterns {
				if p.Pattern == pattern {
					maxConcurrency = &p.MaxConcurrency
					break
				}
			}
		}
		if maxConcurrency != nil {
			break
		}
	}

	if maxConcurrency == nil {
		return common.CheckDeletedDiag(d, fmt.Errorf("not found"), "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("instance_id", instanceID),
		d.Set("node_id", nodeID),
		d.Set("sql_type", sqlType),
		d.Set("pattern", pattern),
		d.Set("max_concurrency", *maxConcurrency),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDbV3MySQLSqlControlRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	if d.HasChange("max_concurrency") {
		instanceID := d.Get("instance_id").(string)
		nodeID := d.Get("node_id").(string)
		sqlType := d.Get("sql_type").(string)
		pattern := d.Get("pattern").(string)
		maxConcurrency := d.Get("max_concurrency").(int)

		opts := sqlfilter.UpdateSqlFilterRulesOpts{
			SqlFilterRules: []sqlfilter.NodeSqlFilterRuleInfo{
				{
					NodeId: nodeID,
					Rules: []sqlfilter.NodeSqlFilterRule{
						{
							SqlType: sqlType,
							Patterns: []sqlfilter.NodeSqlFilterRulePattern{
								{
									Pattern:        pattern,
									MaxConcurrency: maxConcurrency,
								},
							},
						},
					},
				},
			},
		}

		jobID, err := sqlfilter.UpdateSqlFilterRules(client, instanceID, opts)
		if err != nil {
			return diag.Errorf("error updating TaurusDb MySQL SQL control rule: %s", err)
		}

		if jobID == nil || *jobID == "" {
			return diag.Errorf("unable to find the job ID from the API response")
		}

		if err := waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDbV3MySQLSqlControlRuleRead(ctx, d, meta)
}

func resourceTaurusDbV3MySQLSqlControlRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceID := d.Get("instance_id").(string)
	nodeID := d.Get("node_id").(string)
	sqlType := d.Get("sql_type").(string)
	pattern := d.Get("pattern").(string)

	opts := sqlfilter.DeleteSqlFilterRulesOpts{
		SqlFilterRules: []sqlfilter.DeleteNodeSqlFilterRuleInfo{
			{
				NodeId: nodeID,
				Rules: []sqlfilter.DeleteNodeSqlFilterRule{
					{
						SqlType:  sqlType,
						Patterns: []string{pattern},
					},
				},
			},
		},
	}

	jobID, err := sqlfilter.DeleteSqlFilterRules(client, instanceID, opts)
	if err != nil {
		return diag.Errorf("error deleting TaurusDB MySQL SQL control rule: %s", err)
	}

	if jobID == nil || *jobID == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	if err := waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
