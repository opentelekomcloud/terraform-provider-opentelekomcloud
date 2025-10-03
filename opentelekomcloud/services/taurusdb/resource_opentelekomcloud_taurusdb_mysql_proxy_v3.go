package taurusdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/proxy"
)

func ResourceTaurusDbV3Proxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDbProxyCreate,
		ReadContext:   resourceTaurusDbProxyRead,
		UpdateContext: resourceTaurusDbProxyUpdate,
		DeleteContext: resourceTaurusDbProxyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDbMySQLProxyImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"proxy_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"master_node_weight": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     gaussDBMysqlProxyNodeWeightSchema(),
			},
			"readonly_nodes_weight": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     gaussDBMysqlProxyNodeWeightSchema(),
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBMysqlProxyNodeSchema(),
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func gaussDBMysqlProxyNodeWeightSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
	return &sc
}

func gaussDBMysqlProxyNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen_flag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceTaurusDbProxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceID := d.Get("instance_id").(string)

	opts := buildEnableProxyOpts(d)
	jobID, err := proxy.EnableProxy(client, instanceID, opts)
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL proxy: %s", err)
	}

	err = waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	proxyID, err := getProxyIDByStatus(client, instanceID, "ACTIVE")
	if err != nil {
		return diag.Errorf("error retrieving GaussDB MySQL proxy ID: %s", err)
	}
	d.SetId(proxyID)

	return resourceTaurusDbProxyRead(ctx, d, meta)
}

func buildEnableProxyOpts(d *schema.ResourceData) proxy.EnableProxyOpts {
	opts := proxy.EnableProxyOpts{
		FlavorRef: d.Get("flavor").(string),
		NodeNum:   d.Get("node_num").(int),
	}

	if v, ok := d.GetOk("proxy_name"); ok {
		opts.ProxyName = v.(string)
	}

	if v, ok := d.GetOk("proxy_mode"); ok {
		opts.ProxyMode = v.(string)
	}

	nodesWeight := buildNodesReadWeight(d)
	if len(nodesWeight) > 0 {
		opts.NodesReadWeight = nodesWeight
	}

	return opts
}

func buildNodesReadWeight(d *schema.ResourceData) []proxy.NodesWeight {
	var weights []proxy.NodesWeight

	// Master node weight
	if v, ok := d.GetOk("master_node_weight"); ok {
		masterList := v.([]interface{})
		if len(masterList) > 0 {
			master := masterList[0].(map[string]interface{})
			weights = append(weights, proxy.NodesWeight{
				Id:     master["id"].(string),
				Weight: master["weight"].(int),
			})
		}
	}

	// Readonly nodes weight
	if v, ok := d.GetOk("readonly_nodes_weight"); ok {
		readonlySet := v.(*schema.Set).List()
		for _, item := range readonlySet {
			node := item.(map[string]interface{})
			weights = append(weights, proxy.NodesWeight{
				Id:     node["id"].(string),
				Weight: node["weight"].(int),
			})
		}
	}

	return weights
}

func resourceTaurusDbProxyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceID := d.Get("instance_id").(string)
	proxyID := d.Id()

	proxies, err := proxy.List(client, instanceID, proxy.ListOpts{})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving TaurusDb MySQL proxy")
	}

	var proxyInfo *proxy.ProxyInstanceResponse
	for i, p := range proxies {
		if p.Proxy.PoolId == proxyID {
			proxyInfo = &proxies[i]
			break
		}
	}

	if proxyInfo == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "TaurusDb MySQL proxy not found")
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("flavor", proxyInfo.Proxy.FlavorRef),
		d.Set("node_num", proxyInfo.Proxy.NodeNum),
		d.Set("proxy_name", proxyInfo.Proxy.Name),
		d.Set("address", proxyInfo.Proxy.Address),
		d.Set("port", proxyInfo.Proxy.Port),
		d.Set("status", proxyInfo.Proxy.Status),
		d.Set("master_node_weight", flattenMasterNodeWeight(&proxyInfo.MasterNode)),
		d.Set("readonly_nodes_weight", flattenReadonlyNodesWeight(proxyInfo.ReadonlyNodes, d)),
		d.Set("nodes", flattenProxyNodes(proxyInfo.Proxy.Nodes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMasterNodeWeight(node *proxy.MysqlProxyNodeV3) []interface{} {
	if node == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":     node.Id,
			"weight": node.Weight,
		},
	}
}

func flattenReadonlyNodesWeight(nodes []proxy.MysqlProxyNodeV3, d *schema.ResourceData) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	readonlyNodesWeightRaw := d.Get("readonly_nodes_weight").(*schema.Set).List()
	if len(readonlyNodesWeightRaw) == 0 {
		return nil
	}

	configuredNodes := make(map[string]bool)
	for _, v := range readonlyNodesWeightRaw {
		configuredNodes[v.(map[string]interface{})["id"].(string)] = true
	}

	result := make([]interface{}, 0)
	for _, node := range nodes {
		if configuredNodes[node.Id] {
			result = append(result, map[string]interface{}{
				"id":     node.Id,
				"weight": node.Weight,
			})
		}
	}

	return result
}

func flattenProxyNodes(nodes []proxy.MysqlProxyNodes) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"id":          node.Id,
			"status":      node.Status,
			"name":        node.Name,
			"role":        node.Role,
			"az_code":     node.AzCode,
			"frozen_flag": node.FrozenFlag,
		})
	}

	return result
}

func resourceTaurusDbProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceID := d.Get("instance_id").(string)
	proxyID := d.Id()

	if d.HasChange("flavor") {
		err = updateProxyFlavor(ctx, client, instanceID, proxyID, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("node_num") {
		oldNum, newNum := d.GetChange("node_num")
		oldCount := oldNum.(int)
		newCount := newNum.(int)

		if newCount < oldCount {
			return diag.Errorf("reducing the number of proxy nodes is not supported. "+
				"Current nodes: %d, requested: %d. You can only increase the number of nodes.",
				oldCount, newCount)
		}

		if newCount > oldCount {
			err = enlargeProxyNodes(ctx, client, instanceID, proxyID, newCount-oldCount, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChanges("master_node_weight", "readonly_nodes_weight") {
		err = updateProxyWeight(ctx, client, instanceID, proxyID, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDbProxyRead(ctx, d, meta)
}

func updateProxyFlavor(ctx context.Context, client *golangsdk.ServiceClient, instanceID, proxyID string,
	d *schema.ResourceData) error {
	flavorRef := d.Get("flavor").(string)

	jobID, err := proxy.Resize(client, instanceID, proxyID, flavorRef)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy flavor: %s", err)
	}

	err = waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func enlargeProxyNodes(ctx context.Context, client *golangsdk.ServiceClient, instanceID, proxyID string,
	nodeNum int, d *schema.ResourceData) error {
	opts := proxy.EnlargeOpts{
		InstanceID: instanceID,
		NodeNum:    nodeNum,
		ProxyId:    proxyID,
	}

	jobID, err := proxy.Enlarge(client, opts)
	if err != nil {
		return fmt.Errorf("error enlarging GaussDB MySQL proxy nodes: %s", err)
	}

	err = waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func updateProxyWeight(ctx context.Context, client *golangsdk.ServiceClient, instanceID, proxyID string,
	d *schema.ResourceData) error {
	opts := buildUpdateWeightOpts(d)

	jobID, err := proxy.UpdateWeight(client, instanceID, proxyID, opts)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy weight: %s", err)
	}

	err = waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateWeightOpts(d *schema.ResourceData) proxy.UpdateWeightOpts {
	opts := proxy.UpdateWeightOpts{}

	// Master node weight
	if v, ok := d.GetOk("master_node_weight"); ok {
		masterList := v.([]interface{})
		if len(masterList) > 0 {
			master := masterList[0].(map[string]interface{})
			weight := master["weight"].(int)
			opts.MasterWeight = &weight
		}
	}

	// Readonly nodes weight
	if v, ok := d.GetOk("readonly_nodes_weight"); ok {
		readonlySet := v.(*schema.Set).List()
		if len(readonlySet) > 0 {
			nodes := make([]proxy.TaurusModifyProxyWeightReadonlyNode, 0, len(readonlySet))
			for _, item := range readonlySet {
				node := item.(map[string]interface{})
				nodes = append(nodes, proxy.TaurusModifyProxyWeightReadonlyNode{
					Id:     node["id"].(string),
					Weight: node["weight"].(int),
				})
			}
			opts.ReadonlyNodes = nodes
		}
	}

	return opts
}

func resourceTaurusDbProxyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceID := d.Get("instance_id").(string)
	proxyID := d.Id()

	opts := &proxy.DisableProxyOpts{
		ProxyIds: []string{proxyID},
	}

	jobID, err := proxy.DisableProxy(client, instanceID, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting GaussDB MySQL proxy")
	}

	err = waitForJobComplete(ctx, client, *jobID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForJobComplete(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      taurusDbV3MysqlDatabaseStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for job (%s) to complete: %s", jobID, err)
	}

	return nil
}

func getProxyIDByStatus(client *golangsdk.ServiceClient, instanceID, status string) (string, error) {
	proxies, err := proxy.List(client, instanceID, proxy.ListOpts{})
	if err != nil {
		return "", err
	}

	for _, p := range proxies {
		if p.Proxy.Status == status {
			return p.Proxy.PoolId, nil
		}
	}

	return "", fmt.Errorf("proxy with status %s not found", status)
}

func resourceTaurusDbMySQLProxyImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<proxy_id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
