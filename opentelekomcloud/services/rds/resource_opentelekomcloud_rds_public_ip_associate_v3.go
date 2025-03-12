package rds

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsPublicIpAssociateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsPublicIpAssociateV3Create,
		ReadContext:   resourceRdsPublicIpAssociateV3Read,
		UpdateContext: resourceRdsPublicIpAssociateV3Update,
		DeleteContext: resourceRdsPublicIpAssociateV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateIP,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRdsPublicIpAssociateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	d.SetId(d.Get("instance_id").(string))

	ip := d.Get("public_ip").(string)
	ipId := d.Get("public_ip_id").(string)

	jobId, err := instances.AttachEip(client, instances.AttachEipOpts{
		InstanceId: d.Get("instance_id").(string),
		PublicIp:   ip,
		PublicIpId: ipId,
		IsBind:     pointerto.Bool(true),
	})

	if err != nil {
		return fmterr.Errorf("error attaching public ip to rds instance: %w", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	if err := instances.WaitForJobCompleted(client, int(timeout.Seconds()), *jobId); err != nil {
		return diag.FromErr(err)
	}

	return resourceRdsPublicIpAssociateV3Read(ctx, d, meta)
}

func resourceRdsPublicIpAssociateV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	rdsInstance, err := GetRdsInstance(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching RDS instance: %s", err)
	}
	if rdsInstance == nil {
		d.SetId("")
		return nil
	}

	if len(rdsInstance.PublicIps) > 0 {
		if err = d.Set("public_ip", rdsInstance.PublicIps[0]); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceRdsPublicIpAssociateV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	if d.HasChange("public_ip") || d.HasChange("public_ip_id") {
		timeout := d.Timeout(schema.TimeoutUpdate)
		oldIp, _ := d.GetChange("public_ip")
		oldId, _ := d.GetChange("public_ip_id")
		// detach old ip
		jobId, err := instances.AttachEip(client, instances.AttachEipOpts{
			PublicIp:   oldIp.(string),
			PublicIpId: oldId.(string),
			InstanceId: d.Id(),
			IsBind:     pointerto.Bool(false),
		})
		if err != nil {
			return fmterr.Errorf("error detaching old ip: %w", err)
		}

		if err := instances.WaitForJobCompleted(client, int(timeout.Seconds()), *jobId); err != nil {
			return diag.FromErr(err)
		}

		// attach new ip
		jobId, err = instances.AttachEip(client, instances.AttachEipOpts{
			InstanceId: d.Id(),
			PublicIp:   d.Get("public_ip").(string),
			PublicIpId: d.Get("public_ip_id").(string),
			IsBind:     pointerto.Bool(true),
		})
		if err != nil {
			return fmterr.Errorf("error attaching new ip: %w", err)
		}

		if err := instances.WaitForJobCompleted(client, int(timeout.Seconds()), *jobId); err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceRdsPublicIpAssociateV3Read(clientCtx, d, meta)
}

func resourceRdsPublicIpAssociateV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	log.Printf("[DEBUG] Unassigning public ip for Instance %s", d.Id())

	jobId, err := instances.AttachEip(client, instances.AttachEipOpts{
		InstanceId: d.Get("instance_id").(string),
		IsBind:     pointerto.Bool(false),
	})

	if err != nil {
		return fmterr.Errorf("error detaching public ip from RDS instance: %w", err)
	}
	timeout := d.Timeout(schema.TimeoutUpdate)
	if err := instances.WaitForJobCompleted(client, int(timeout.Seconds()), *jobId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
