package apigw

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/key"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPISignatureV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGWSignatureV2Create,
		ReadContext:   resourceAPIGWSignatureV2Read,
		UpdateContext: resourceAPIGWSignatureV2Update,
		DeleteContext: resourceAPIGWSignatureV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGWSignatureV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
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
		},
	}
}

func resourceAPIGWSignatureV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := key.CreateOpts{
		GatewayID:     d.Get("gateway_id").(string),
		Name:          d.Get("name").(string),
		SignType:      d.Get("type").(string),
		SignKey:       d.Get("key").(string),
		SignSecret:    d.Get("secret").(string),
		SignAlgorithm: d.Get("algorithm").(string),
	}
	sign, err := key.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating the signature key: %s", err)
	}
	d.SetId(sign.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWSignatureV2Read(clientCtx, d, meta)
}

func resourceAPIGWSignatureV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	sign, err := key.List(client, key.ListOpts{
		GatewayID:   d.Get("gateway_id").(string),
		SignatureID: d.Id(),
	})
	log.Printf("[DEBUG] The signature is: %#v", sign)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Signature")
	}
	if len(sign) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Signature")
	}

	signature := sign[0]
	mErr := multierror.Append(nil,
		d.Set("name", signature.Name),
		d.Set("region", config.GetRegion(d)),
		d.Set("type", signature.SignType),
		d.Set("key", signature.SignKey),
		d.Set("secret", signature.SignSecret),
		d.Set("algorithm", signature.SignAlgorithm),
		d.Set("created_at", signature.CreateTime),
		d.Set("updated_at", signature.UpdateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving signature (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceAPIGWSignatureV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := key.UpdateOpts{
		GatewayID:     d.Get("gateway_id").(string),
		SignID:        d.Id(),
		Name:          d.Get("name").(string),
		SignType:      d.Get("type").(string),
		SignKey:       d.Get("key").(string),
		SignSecret:    d.Get("secret").(string),
		SignAlgorithm: d.Get("algorithm").(string),
	}
	_, err = key.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating the signature (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWSignatureV2Read(clientCtx, d, meta)
}

func resourceAPIGWSignatureV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var (
		instanceId  = d.Get("gateway_id").(string)
		signatureId = d.Id()
	)
	err = key.Delete(client, instanceId, signatureId)
	if err != nil {
		return diag.Errorf("error deleting the signature (%s): %s", signatureId, err)
	}
	return nil
}

func resourceAPIGWSignatureV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	err := d.Set("gateway_id", parts[0])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving gateway ID field: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
