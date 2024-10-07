package dms

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/lifecycle"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/smart_connect"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsSmartConnectV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsSmartConnectV2Create,
		ReadContext:   resourceDmsSmartConnectV2Read,
		DeleteContext: resourceDmsSmartConnectV2Delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsSmartConnectV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceId := d.Get("instance_id").(string)

	nodeCount := strconv.Itoa(d.Get("node_count").(int))

	enableOpts := smart_connect.EnableOpts{
		InstanceId:    instanceId,
		SpecCode:      d.Get("storage_spec_code").(string),
		NodeCount:     nodeCount,
		Specification: d.Get("bandwidth").(string),
	}

	enableResp, err := smart_connect.Enable(client, enableOpts)
	if err != nil {
		return diag.Errorf("error enabling smart connect for Kafka instance: %s", err)
	}

	d.SetId(enableResp.ConnectorId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"EXTENDING"},
		Target:     []string{"RUNNING"},
		Refresh:    instancesV2StateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", instanceId, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsSmartConnectV2Read(clientCtx, d, meta)
}

func resourceDmsSmartConnectV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	getResp, err := lifecycle.Get(client, d.Get("instance_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS Kafka instance")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("instance_id", getResp.InstanceID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsSmartConnectV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceId := d.Get("instance_id").(string)

	_, err = smart_connect.Disable(client, instanceId)

	if err != nil {
		return diag.Errorf("error disabling DMSv2 smart connect: %v", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CONNECTOR_DELETING"},
		Target:       []string{"RUNNING"},
		Refresh:      instancesV2StateRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", instanceId, err)
	}

	d.SetId("")
	log.Printf("[DEBUG] DMS Kafka instance %s smar connect has been disabled", d.Id())
	return nil
}
