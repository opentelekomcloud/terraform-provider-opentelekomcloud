package er

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceErInstanceV3Create,
		UpdateContext: resourceErInstanceV3Update,
		ReadContext:   resourceErInstanceV3Read,
		DeleteContext: resourceErInstanceV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_default_propagation": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_default_association": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_accept_shared_attachments": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"default_propagation_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_association_route_table_id": {
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

func getAvailabilityZones(d *schema.ResourceData) []string {
	var zones []string
	azRaw := d.Get("availability_zones").([]interface{})
	for _, az := range azRaw {
		zones = append(zones, az.(string))
	}
	return zones
}

func resourceErInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := instance.CreateOpts{
		Name:                        d.Get("name").(string),
		Description:                 d.Get("description").(string),
		Asn:                         float64(d.Get("asn").(int)),
		EnableDefaultPropagation:    pointerto.Bool(d.Get("enable_default_propagation").(bool)),
		EnableDefaultAssociation:    pointerto.Bool(d.Get("enable_default_association").(bool)),
		AvailabilityZoneIDs:         getAvailabilityZones(d),
		AutoAcceptSharedAttachments: pointerto.Bool(d.Get("auto_accept_shared_attachments").(bool)),
	}

	createResp, err := instance.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	d.SetId(createResp.Instance.ID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: waitForErInstanceActive(client, d.Id()),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for ER instance to be active: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceErInstanceV3Read(clientCtx, d, meta)
}

func resourceErInstanceV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	getResp, err := instance.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving er instance (%s): %s", d.Id(), err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", getResp.Instance.Name),
		d.Set("description", getResp.Instance.Description),
		d.Set("status", getResp.Instance.State),
		d.Set("created_at", getResp.Instance.CreatedAt),
		d.Set("updated_at", getResp.Instance.UpdatedAt),
		d.Set("asn", getResp.Instance.Asn),
		d.Set("enable_default_propagation", getResp.Instance.EnableDefaultPropagation),
		d.Set("enable_default_association", getResp.Instance.EnableDefaultAssociation),
		d.Set("default_propagation_route_table_id", getResp.Instance.DefaultPropagationRouteTableID),
		d.Set("default_association_route_table_id", getResp.Instance.DefaultAssociationRouteTableID),
		d.Set("availability_zones", getResp.Instance.AvailabilityZoneIDs),
		d.Set("auto_accept_shared_attachments", getResp.Instance.AutoAcceptSharedAttachments),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceErInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	updateInstancehasChanges := []string{
		"name",
		"description",
		"enable_default_propagation",
		"enable_default_association",
		"auto_accept_shared_attachments",
	}

	if d.HasChanges(updateInstancehasChanges...) {
		updateOpts := instance.UpdateOpts{
			InstanceID:                  d.Id(),
			Name:                        d.Get("name").(string),
			Description:                 d.Get("description").(string),
			EnableDefaultPropagation:    pointerto.Bool(d.Get("enable_default_propagation").(bool)),
			EnableDefaultAssociation:    pointerto.Bool(d.Get("enable_default_association").(bool)),
			AutoAcceptSharedAttachments: pointerto.Bool(d.Get("auto_accept_shared_attachments").(bool)),
		}

		_, err = instance.Update(client, updateOpts)
		if err != nil {
			return diag.Errorf("error updating Instance: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceErInstanceV3Read(clientCtx, d, meta)
}

func resourceErInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = instance.Delete(client, d.Id())
	if err != nil {
		return diag.Errorf("error deleting Instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Target:  []string{"Deleted"},
		Refresh: waitForErInstanceDeletion(client, d.Id()),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for the Deletion of Instance (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func waitForErInstanceActive(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instance.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil, "", nil
			}
			return nil, "", err
		}
		return resp, resp.Instance.State, nil
	}
}

func waitForErInstanceDeletion(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instance.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The opentelekomcloud ER instance has been deleted (ID:%s).", id)
				return r, "Deleted", nil
			}
			return nil, "Error", err
		}
		switch r.Instance.State {
		case "available", "deleting":
			return r, "Deleting", nil
		default:
			err = fmt.Errorf("error deleting ER instance[%s]. "+
				"Unexpected status: %v", r.Instance.ID, r.Instance.State)
			return r, "Error", err
		}
	}
}
