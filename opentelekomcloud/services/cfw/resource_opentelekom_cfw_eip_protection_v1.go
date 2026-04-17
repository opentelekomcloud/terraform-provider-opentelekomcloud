package cfw

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	cfweip "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/eip"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwEipProtectionV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWEipProtectionV1Create,
		ReadContext:   resourceCFWEipProtectionV1Read,
		DeleteContext: resourceCFWEipProtectionV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"firewall_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"eip_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"public_ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCFWEipProtectionV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	eipId := d.Get("eip_id").(string)
	objectId := d.Get("object_id").(string)
	status := d.Get("status").(int)

	err = waitForIPSync(client, 300, 5, objectId, eipId)
	if err != nil {
		return fmterr.Errorf("error finding IP with ID '%s' in list: %s", eipId, err)
	}

	createOpts := cfweip.ChangeEIPProtectionOpts{
		ObjectID: objectId,
		Status:   status,
		IPInfos: []cfweip.IPInfo{
			{
				ID:         eipId,
				PublicIP:   d.Get("public_ip").(string),
				PublicIPv6: d.Get("public_ipv6").(string),
			},
		},
	}

	firewallId := d.Get("firewall_id").(string)

	ChangeEIPProtectionResults, err := cfweip.ChangeEIPProtection(client, firewallId, createOpts)
	if err != nil {
		return fmterr.Errorf("error enabling/disabling EIP protection: %w", err)
	}
	if len(ChangeEIPProtectionResults.FailEIPIDList) != 0 {
		return fmterr.Errorf("error enabling/disabling EIP protection ID: %s", eipId)
	}
	log.Printf("[DEBUG] EIP '%s' protection status set to: %d", eipId, status)

	d.SetId(eipId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWEipProtectionV1Read(clientCtx, d, meta)
}

func resourceCFWEipProtectionV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	eipList, err := cfweip.List(client, cfweip.ListOpts{
		ObjectID: d.Get("object_id").(string),
	})
	if err != nil {
		return fmterr.Errorf("error fetching EIP List: %w", err)
	}

	for _, ip := range eipList {
		if ip.ID == d.Id() {
			mErr := multierror.Append(nil,
				d.Set("public_ip", ip.PublicIP),
				d.Set("public_ipv6", ip.PublicIPV6),
				d.Set("status", ip.Status),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	return fmterr.Errorf("unable to find EIP in EIP list or EIP does not exist: %w", err)
}

func resourceCFWEipProtectionV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}

func waitForIPSync(client *golangsdk.ServiceClient, waitTime int, interval time.Duration, objectId string, eipID string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	queryOpts := cfweip.ListOpts{
		ObjectID: objectId,
	}
	return golangsdk.WaitFor(waitTime, func() (bool, error) {
		eipResourceList, err := cfweip.List(client, queryOpts)
		if err != nil {
			return false, err
		}

		for _, eipResource := range eipResourceList {
			if eipResource.ID == eipID {
				return true, nil
			}
		}

		time.Sleep(interval * time.Second)
		return false, nil
	})
}
