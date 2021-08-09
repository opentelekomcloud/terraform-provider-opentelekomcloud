package antiddos

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/antiddos/v1/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAntiDdosV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAntiDdosV1Create,
		ReadContext:   resourceAntiDdosV1Read,
		UpdateContext: resourceAntiDdosV1Update,
		DeleteContext: resourceAntiDdosV1Delete,

		DeprecationMessage: "AntiDDoS protection for Elastic IP is provided by default and shouldn't be created manually.",

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_l7": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"traffic_pos_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: common.ValidateAntiDdosTrafficPosID,
			},
			"http_request_pos_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: common.ValidateAntiDdosHttpRequestPosID,
			},
			"cleaning_access_pos_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: common.ValidateAntiDdosCleaningAccessPosID,
			},
			"app_type_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: common.ValidateAntiDdosAppTypeID,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAntiDdosV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	antiddosClient, err := config.AntiddosV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating AntiDdos client: %s", err)
	}

	createOpts := antiddos.CreateOpts{
		EnableL7:            d.Get("enable_l7").(bool),
		TrafficPosId:        d.Get("traffic_pos_id").(int),
		HttpRequestPosId:    d.Get("http_request_pos_id").(int),
		CleaningAccessPosId: d.Get("cleaning_access_pos_id").(int),
		AppTypeId:           d.Get("app_type_id").(int),
	}

	_, err = antiddos.Create(antiddosClient, d.Get("floating_ip_id").(string), createOpts).Extract()

	if err != nil {
		return fmterr.Errorf("error creating AntiDdos: %s", err)
	}

	d.SetId(d.Get("floating_ip_id").(string))

	log.Printf("[INFO] AntiDdos ID: %s", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"configging"},
		Target:     []string{"normal"},
		Refresh:    waitForAntiDdosActive(antiddosClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return fmterr.Errorf("error waiting for AntiDdos (%s) to become normal: %s", d.Id(), stateErr)
	}

	return resourceAntiDdosV1Read(ctx, d, meta)
}

func resourceAntiDdosV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	antiddosClient, err := config.AntiddosV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating AntiDdos client: %s", err)
	}

	n, err := antiddos.Get(antiddosClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving AntiDdos: %s", err)
	}

	mErr := multierror.Append(
		d.Set("floating_ip_id", d.Id()),
		d.Set("enable_l7", n.EnableL7),
		d.Set("app_type_id", n.AppTypeId),
		d.Set("cleaning_access_pos_id", n.CleaningAccessPosId),
		d.Set("traffic_pos_id", n.TrafficPosId),
		d.Set("http_request_pos_id", n.HttpRequestPosId),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAntiDdosV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	antiddosClient, err := config.AntiddosV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating AntiDdos client: %s", err)
	}

	var updateOpts antiddos.UpdateOpts

	updateOpts.EnableL7 = d.Get("enable_l7").(bool)
	updateOpts.AppTypeId = d.Get("app_type_id").(int)
	updateOpts.CleaningAccessPosId = d.Get("cleaning_access_pos_id").(int)
	updateOpts.TrafficPosId = d.Get("traffic_pos_id").(int)
	updateOpts.HttpRequestPosId = d.Get("http_request_pos_id").(int)

	_, err = antiddos.Update(antiddosClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating AntiDdos: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"configging"},
		Target:     []string{"normal"},
		Refresh:    waitForAntiDdosActive(antiddosClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      3 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return fmterr.Errorf("error waiting for AntiDdos to become normal: %s", stateErr)
	}

	return resourceAntiDdosV1Read(ctx, d, meta)
}

func resourceAntiDdosV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	antiddosClient, err := config.AntiddosV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating AntiDdos client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "configging"},
		Target:     []string{"notConfig"},
		Refresh:    waitForAntiDdosDelete(antiddosClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting AntiDdos: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForAntiDdosActive(antiddosClient *golangsdk.ServiceClient, antiddosId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := antiddos.Get(antiddosClient, antiddosId).Extract()
		if err != nil {
			return nil, "", err
		}

		return s, "normal", nil
	}
}

func waitForAntiDdosDelete(antiddosClient *golangsdk.ServiceClient, antiddosId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ddosstatus, err := antiddos.ListStatus(antiddosClient, antiddos.ListStatusOpts{FloatingIpId: antiddosId})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				log.Printf("[INFO] Successfully deleted AntiDdos %s", antiddosId)
				return ddosstatus, "notConfig", nil
			}
			return ddosstatus, "normal", err
		}
		r := ddosstatus[0]
		if r.Status != "configging" {
			_, err = antiddos.Delete(antiddosClient, antiddosId).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault403); ok {
					log.Printf("[INFO] Successfully deleted Antiddos %s", antiddosId)
					return r, "notConfig", nil
				}
				return r, r.Status, err
			}
		}

		return r, "normal", nil
	}
}
