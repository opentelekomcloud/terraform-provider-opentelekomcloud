package ces

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/dashboards"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDashboardV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardV2Create,
		UpdateContext: resourceDashboardV2Update,
		ReadContext:   resourceDashboardV2Read,
		DeleteContext: resourceDashboardV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
						"only letters, digits, underscores (_), and hyphens (-) are allowed",
					),
				),
			},
			"row_widget_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dashboard_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_favorite": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
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

func resourceDashboardV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	rowWidgetNum := d.Get("row_widget_num").(int)
	createOpts := dashboards.CreateOpts{
		DashboardName: d.Get("name").(string),
		RowWidgetNum:  &rowWidgetNum,
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		createOpts.EnterpriseId = v.(string)
	}

	if v, ok := d.GetOk("dashboard_id"); ok {
		createOpts.DashboardId = v.(string)
	}

	log.Printf("[DEBUG] CES Dashboard V2 Create Options: %#v", createOpts)

	dashboardId, err := dashboards.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CES Dashboard V2: %w", err)
	}

	log.Printf("[DEBUG] CES Dashboard V2 created with ID: %s", dashboardId)
	d.SetId(dashboardId)

	// is_favorite is not supported in Create API, set it via Update.
	// When copying a dashboard, the API ignores row_widget_num from the request,
	// so we also need to update it after creation.
	_, isFavoriteOk := d.GetOk("is_favorite")
	_, dashboardIdOk := d.GetOk("dashboard_id")
	if isFavoriteOk || dashboardIdOk {
		err = updateDashboardV2(client, d)
		if err != nil {
			return fmterr.Errorf("error updating CES Dashboard V2 after creation: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceDashboardV2Read(clientCtx, d, meta)
}

func resourceDashboardV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	listOpts := dashboards.ListOpts{
		DashboardId: d.Id(),
	}
	dashboardList, err := dashboards.List(client, listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CES Dashboard V2")
	}

	if len(dashboardList) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "CES Dashboard V2")
	}

	dashboard := dashboardList[0]

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", dashboard.DashboardName),
		d.Set("row_widget_num", dashboard.RowWidgetNum),
		d.Set("enterprise_project_id", dashboard.EnterpriseId),
		d.Set("is_favorite", dashboard.IsFavorite),
		d.Set("creator_name", dashboard.CreatorName),
	)

	if dashboard.CreateTime > 0 {
		mErr = multierror.Append(mErr, d.Set("created_at", time.Unix(dashboard.CreateTime/1000, 0).Format(time.RFC3339)))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDashboardV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	err = updateDashboardV2(client, d)
	if err != nil {
		return fmterr.Errorf("error updating CES Dashboard V2: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceDashboardV2Read(clientCtx, d, meta)
}

func resourceDashboardV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	dashboardId := d.Id()
	log.Printf("[DEBUG] Deleting CES Dashboard V2: %s", dashboardId)

	deleteOpts := dashboards.BatchDeleteOpts{
		DashboardIds: []string{dashboardId},
	}

	_, err = dashboards.BatchDelete(client, deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES Dashboard V2")
	}

	return nil
}

func updateDashboardV2(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	rowWidgetNum := d.Get("row_widget_num").(int)
	isFavorite := d.Get("is_favorite").(bool)
	updateOpts := dashboards.UpdateOpts{
		DashboardName: d.Get("name").(string),
		RowWidgetNum:  &rowWidgetNum,
		IsFavorite:    &isFavorite,
	}

	log.Printf("[DEBUG] CES Dashboard V2 Update Options: %#v", updateOpts)

	err := dashboards.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmt.Errorf("error updating CES Dashboard V2: %w", err)
	}

	return nil
}
