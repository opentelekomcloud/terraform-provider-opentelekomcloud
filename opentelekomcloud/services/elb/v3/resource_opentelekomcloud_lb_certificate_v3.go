package v3

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/certificates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCertificateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateV3Create,
		ReadContext:   resourceCertificateV3Read,
		UpdateContext: resourceCertificateV3Update,
		DeleteContext: resourceCertificateV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: common.SuppressStrippedNewLines,
			},
			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: common.SuppressStrippedNewLines,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"server", "client",
				}, false),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := certificates.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Domain:      d.Get("domain").(string),
		PrivateKey:  d.Get("private_key").(string),
		Certificate: d.Get("certificate").(string),
		Type:        d.Get("type").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	cert, err := certificates.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating LoadBalancer Certificate: %w", err)
	}

	// If all has been successful, set the ID on the resource
	d.SetId(cert.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceCertificateV3Read(clientCtx, d, meta)
}

func resourceCertificateV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	cert, err := certificates.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "certificateV3")
	}
	log.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), cert)

	mErr := multierror.Append(nil,
		d.Set("name", cert.Name),
		d.Set("description", cert.Description),
		d.Set("domain", cert.Domain),
		d.Set("certificate", cert.Certificate),
		d.Set("private_key", cert.PrivateKey),
		d.Set("type", cert.Type),
		d.Set("created_at", cert.CreatedAt),
		d.Set("updated_at", cert.UpdatedAt),
		d.Set("expire_time", cert.ExpireTime),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCertificateV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	var updateOpts certificates.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("domain") {
		updateOpts.Domain = d.Get("domain").(string)
	}
	if d.HasChange("private_key") {
		updateOpts.PrivateKey = d.Get("private_key").(string)
	}
	if d.HasChange("certificate") {
		updateOpts.Certificate = d.Get("certificate").(string)
	}

	log.Printf("[DEBUG] Updating certificate %s with options: %#v", d.Id(), updateOpts)

	_, err = certificates.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating the certificate: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceCertificateV3Read(clientCtx, d, meta)
}

func resourceCertificateV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	log.Printf("[DEBUG] Deleting certificate: %s", d.Id())
	if err := certificates.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting the certificate: %w", err)
	}

	return nil
}
