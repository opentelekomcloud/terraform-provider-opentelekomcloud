package cc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

// defaultPlaneName is the only plane name currently supported by the central network policy API.
const defaultPlaneName = "default-plane"

func ResourceCcCentralNetworkPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcCentralNetworkPolicyV3Create,
		ReadContext:   resourceCcCentralNetworkPolicyV3Read,
		DeleteContext: resourceCcCentralNetworkPolicyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCcCentralNetworkPolicyV3ImportState,
		},

		Schema: map[string]*schema.Schema{
			"central_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"er_instances": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     centralNetworkPolicyErInstanceSchema(),
			},
			"planes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     centralNetworkPolicyPlaneSchema(),
			},
			"document_template_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_applied": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func centralNetworkPolicyPlaneSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"associate_er_tables": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     centralNetworkPolicyAssociateErTableSchema(),
			},
			"exclude_er_connections": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     centralNetworkPolicyExcludeErConnectionsSchema(),
			},
		},
	}
}

func centralNetworkPolicyAssociateErTableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_router_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func centralNetworkPolicyExcludeErConnectionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude_er_instances": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     centralNetworkPolicyErInstanceSchema(),
			},
		},
	}
}

func centralNetworkPolicyErInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCcCentralNetworkPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := policy.CreateOpts{
		CentralNetworkId: d.Get("central_network_id").(string),
		DefaultPlane:     defaultPlaneName,
		Planes:           buildCentralNetworkPolicyPlanes(d.Get("planes").([]interface{})),
		ErInstances:      buildCentralNetworkPolicyErInstances(d.Get("er_instances").([]interface{})),
	}

	created, err := policy.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud CC central network policy: %s", err)
	}

	d.SetId(created.ID)

	clientCtx := common.CtxWithClient(ctx, client, ccClientV3)
	return resourceCcCentralNetworkPolicyV3Read(clientCtx, d, meta)
}

func buildCentralNetworkPolicyPlanes(rawPlanes []interface{}) []policy.PlaneDocument {
	// The API always expects at least one plane carrying the default plane name.
	if len(rawPlanes) == 0 {
		return []policy.PlaneDocument{{Name: defaultPlaneName}}
	}

	planes := make([]policy.PlaneDocument, len(rawPlanes))
	for i, v := range rawPlanes {
		raw := v.(map[string]interface{})
		planes[i] = policy.PlaneDocument{
			Name:                 defaultPlaneName,
			AssociateErTables:    buildCentralNetworkPolicyAssociateErTables(raw["associate_er_tables"].([]interface{})),
			ExcludeErConnections: buildCentralNetworkPolicyExcludeErConnections(raw["exclude_er_connections"].([]interface{})),
		}
	}
	return planes
}

func buildCentralNetworkPolicyAssociateErTables(rawTables []interface{}) []policy.AssociateErTable {
	if len(rawTables) == 0 {
		return nil
	}
	tables := make([]policy.AssociateErTable, len(rawTables))
	for i, v := range rawTables {
		raw := v.(map[string]interface{})
		tables[i] = policy.AssociateErTable{
			ProjectId:               raw["project_id"].(string),
			RegionId:                raw["region_id"].(string),
			EnterpriseRouterId:      raw["enterprise_router_id"].(string),
			EnterpriseRouterTableId: raw["enterprise_router_table_id"].(string),
		}
	}
	return tables
}

func buildCentralNetworkPolicyExcludeErConnections(rawConnections []interface{}) [][]policy.AssociateErInstance {
	if len(rawConnections) == 0 {
		return nil
	}
	connections := make([][]policy.AssociateErInstance, len(rawConnections))
	for i, v := range rawConnections {
		raw := v.(map[string]interface{})
		connections[i] = buildCentralNetworkPolicyErInstances(raw["exclude_er_instances"].([]interface{}))
	}
	return connections
}

func buildCentralNetworkPolicyErInstances(rawInstances []interface{}) []policy.AssociateErInstance {
	if len(rawInstances) == 0 {
		return nil
	}
	instances := make([]policy.AssociateErInstance, len(rawInstances))
	for i, v := range rawInstances {
		raw := v.(map[string]interface{})
		instances[i] = policy.AssociateErInstance{
			ProjectId:          raw["project_id"].(string),
			RegionId:           raw["region_id"].(string),
			EnterpriseRouterId: raw["enterprise_router_id"].(string),
		}
	}
	return instances
}

func resourceCcCentralNetworkPolicyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	pol, err := getCentralNetworkPolicy(client, d.Get("central_network_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud CC central network policy")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("central_network_id", pol.CentralNetworkId),
		d.Set("document_template_version", pol.DocumentTemplateVersion),
		d.Set("is_applied", pol.IsApplied),
		d.Set("version", pol.Version),
		d.Set("state", pol.State),
		d.Set("er_instances", flattenCentralNetworkPolicyErInstances(pol.Document.ErInstances)),
		d.Set("planes", flattenCentralNetworkPolicyPlanes(pol.Document.Planes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// getCentralNetworkPolicy looks up a single policy by ID. The API only exposes a list endpoint,
// so we filter the result client-side and surface a 404 when the policy is gone.
func getCentralNetworkPolicy(client *golangsdk.ServiceClient, centralNetworkId, policyId string) (*policy.CentralNetworkPolicy, error) {
	resp, err := policy.List(client, policy.ListOpts{
		CentralNetworkId: centralNetworkId,
		ID:               []string{policyId},
	})
	if err != nil {
		return nil, err
	}
	for i := range resp.CentralNetworkPolicies {
		if resp.CentralNetworkPolicies[i].ID == policyId {
			return &resp.CentralNetworkPolicies[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func flattenCentralNetworkPolicyPlanes(planes []policy.PlaneDocument) []interface{} {
	if len(planes) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(planes))
	for _, plane := range planes {
		rst = append(rst, map[string]interface{}{
			"associate_er_tables":    flattenCentralNetworkPolicyAssociateErTables(plane.AssociateErTables),
			"exclude_er_connections": flattenCentralNetworkPolicyExcludeErConnections(plane.ExcludeErConnections),
		})
	}
	return rst
}

func flattenCentralNetworkPolicyAssociateErTables(tables []policy.AssociateErTable) []interface{} {
	if len(tables) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(tables))
	for _, table := range tables {
		rst = append(rst, map[string]interface{}{
			"project_id":                 table.ProjectId,
			"region_id":                  table.RegionId,
			"enterprise_router_id":       table.EnterpriseRouterId,
			"enterprise_router_table_id": table.EnterpriseRouterTableId,
		})
	}
	return rst
}

func flattenCentralNetworkPolicyExcludeErConnections(connections [][]policy.AssociateErInstance) []interface{} {
	if len(connections) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(connections))
	for _, connection := range connections {
		rst = append(rst, map[string]interface{}{
			"exclude_er_instances": flattenCentralNetworkPolicyErInstances(connection),
		})
	}
	return rst
}

func flattenCentralNetworkPolicyErInstances(instances []policy.AssociateErInstance) []interface{} {
	if len(instances) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(instances))
	for _, instance := range instances {
		rst = append(rst, map[string]interface{}{
			"project_id":           instance.ProjectId,
			"region_id":            instance.RegionId,
			"enterprise_router_id": instance.EnterpriseRouterId,
		})
	}
	return rst
}

func resourceCcCentralNetworkPolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if err = policy.Delete(client, d.Get("central_network_id").(string), d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud CC central network policy")
	}

	return nil
}

func resourceCcCentralNetworkPolicyV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <central_network_id>/<id>")
	}

	if err := d.Set("central_network_id", parts[0]); err != nil {
		return nil, err
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
