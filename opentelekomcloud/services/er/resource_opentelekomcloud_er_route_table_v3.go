package er

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/association"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/propagation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/route_table"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErRouteTableV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouteTableV3Create,
		UpdateContext: resourceRouteTableV3Update,
		ReadContext:   resourceRouteTableV3Read,
		DeleteContext: resourceRouteTableV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRouteTableImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 255),
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
				),
			},
			"is_default_association": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default_propagation": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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

func buildRouteTableCreateOpts(d *schema.ResourceData) route_table.CreateOpts {
	return route_table.CreateOpts{
		RouterID:    d.Get("instance_id").(string),
		Name:        d.Get("name").(string),
		Description: pointerto.String(d.Get("description").(string)),
	}
}

func resourceRouteTableV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	opts := buildRouteTableCreateOpts(d)
	resp, err := route_table.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating route table: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      routeTableStatusRefreshFunc(client, instanceId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceRouteTableV3Read(clientCtx, d, meta)
}

func routeTableStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := route_table.Get(client, instanceId, routeTableId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}
		log.Printf("[DEBUG] The details of the route table (%s) is: %#v", routeTableId, resp)

		if common.StrSliceContains([]string{"failed"}, resp.State) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.State)
		}
		if common.StrSliceContains(targets, resp.State) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceRouteTableV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	routeTableId := d.Id()
	resp, err := route_table.Get(client, instanceId, routeTableId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER route table")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("is_default_association", resp.IsDefaultAssociation),
		d.Set("is_default_propagation", resp.IsDefaultPropagation),
		d.Set("status", resp.State),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving route table (%s) fields: %s", routeTableId, mErr)
	}
	return nil
}

func updateRouteTableBasicInfo(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		instanceId   = d.Get("instance_id").(string)
		routeTableId = d.Id()
	)

	opts := route_table.UpdateOpts{
		RouterID:     instanceId,
		RouteTableId: routeTableId,
		Name:         d.Get("name").(string),
		Description:  pointerto.String(d.Get("description").(string)),
	}

	_, err := route_table.Update(client, opts)
	if err != nil {
		return fmt.Errorf("error updating route table (%s): %s", routeTableId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      routeTableStatusRefreshFunc(client, instanceId, routeTableId, []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceRouteTableV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if d.HasChanges("name", "description") {
		if err = updateRouteTableBasicInfo(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRouteTableV3Read(ctx, d, meta)
}

func releaseRouteTableAssociations(client *golangsdk.ServiceClient, instanceId, routeTableId string) error {
	resp, err := association.List(client, association.ListOpts{
		RouterId:     instanceId,
		RouteTableId: routeTableId,
	})
	if err != nil {
		return fmt.Errorf("error getting association list from the specified route table (%s): %s", routeTableId, err)
	}
	for _, respAssociate := range resp.Associations {
		opts := association.DeleteOpts{
			RouterID:     instanceId,
			RouteTableID: routeTableId,
			AttachmentID: respAssociate.AttachmentID,
		}
		err := association.Delete(client, opts)
		if err != nil {
			return fmt.Errorf("error disable the association: %s", err)
		}
	}

	return nil
}

func releaseRouteTablePropagations(client *golangsdk.ServiceClient, instanceId, routeTableId string) error {
	resp, err := propagation.List(client, propagation.ListOpts{
		RouterId:     instanceId,
		RouteTableId: routeTableId,
	})
	if err != nil {
		return fmt.Errorf("error getting association list from the specified route table (%s): %s", routeTableId, err)
	}
	for _, respPropagation := range resp.Propagations {
		opts := propagation.DeleteOpts{
			RouterID:     instanceId,
			RouteTableID: routeTableId,
			AttachmentID: respPropagation.AttachmentID,
		}
		err := propagation.Delete(client, opts)
		if err != nil {
			return fmt.Errorf("error disable the propagation: %s", err)
		}
	}

	return nil
}

func resourceRouteTableV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	routeTableId := d.Id()

	err = releaseRouteTableAssociations(client, instanceId, routeTableId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = releaseRouteTablePropagations(client, instanceId, routeTableId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = route_table.Delete(client, instanceId, routeTableId)
	if err != nil {
		return diag.Errorf("error deleting route table (%s): %s", routeTableId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      routeTableStatusRefreshFunc(client, instanceId, routeTableId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRouteTableImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid format for import ID, want '<instance_id>/<route_table_id>', but '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
