package er

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/association"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErAssociationV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssociationV3Create,
		ReadContext:   resourceAssociationV3Read,
		DeleteContext: resourceAssociationV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAssociationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
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

func resourceAssociationV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		instanceId   = d.Get("instance_id").(string)
		routeTableId = d.Get("route_table_id").(string)

		opts = association.CreateOpts{
			RouterID:     instanceId,
			RouteTableID: routeTableId,
			AttachmentID: d.Get("attachment_id").(string),
		}
	)

	resp, err := association.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating the association to the route table: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      associationStatusRefreshFunc(client, instanceId, routeTableId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceAssociationV3Read(clientCtx, d, meta)
}

func QueryAssociationById(client *golangsdk.ServiceClient, instanceId, routeTableId,
	associationId string) (*association.Association, error) {
	resp, err := association.List(client, association.ListOpts{
		RouterId:     instanceId,
		RouteTableId: routeTableId,
	})
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"ID": associationId,
	}
	result, err := common.FilterSliceWithField(resp.Associations, filter)
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the association (%s) does not exist", associationId)),
			},
		}
	}

	log.Printf("[DEBUG] The result filtered by resource ID (%s) is: %#v", associationId, result)
	associationResp, ok := result[0].(association.Association)
	if !ok {
		return nil, fmt.Errorf("the element type of filter result is incorrect, want 'associations.Association', but got '%T'", result[0])
	}

	return &associationResp, nil
}

func associationStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId, associationId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := QueryAssociationById(client, instanceId, routeTableId, associationId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}

		if common.StrSliceContains([]string{"failed"}, resp.State) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.State)
		}
		if common.StrSliceContains(targets, resp.State) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceAssociationV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		associationId = d.Id()
	)

	resp, err := QueryAssociationById(client, instanceId, routeTableId, associationId)
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "ER association")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("route_table_id", resp.RouteTableID),
		d.Set("attachment_id", resp.AttachmentID),
		d.Set("attachment_type", resp.ResourceType),
		d.Set("status", resp.State),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving association (%s) fields: %s", associationId, mErr)
	}
	return nil
}

func resourceAssociationV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		associationId = d.Id()

		opts = association.DeleteOpts{
			RouterID:     instanceId,
			RouteTableID: routeTableId,
			AttachmentID: d.Get("attachment_id").(string),
		}
	)
	err = association.Delete(client, opts)
	if err != nil {
		return diag.Errorf("error deleting association (%s): %s", associationId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      associationStatusRefreshFunc(client, instanceId, routeTableId, associationId, nil),
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

func resourceAssociationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<route_table_id>/<association_id>', but got '%s'", d.Id())
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("route_table_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
