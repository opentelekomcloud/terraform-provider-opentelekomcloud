package v3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/certificates"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
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

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
		},
	}
}

func resourceCertificateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	cert, err := certificates.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "certificateV3"))
	}
	log.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), cert)

	mErr := multierror.Append(nil,
		d.Set("name", cert.Name),
		d.Set("description", cert.Description),
		d.Set("domain", cert.Domain),
		d.Set("certificate", cert.Certificate),
		d.Set("private_key", cert.PrivateKey),
		d.Set("type", cert.Type),
		d.Set("created_at", cert.CreateTime),
		d.Set("updated_at", cert.UpdateTime),
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
		return diag.FromErr(err)
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

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err := certificates.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error updating certificate %s: %s", d.Id(), err)
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
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting certificate: %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	if err := certificates.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			if err := handleCertificateDeletionError(ctx, d, client, err); err != nil {
				return resource.RetryableError(err)
			}
			return nil
		}))
	}

	return nil
}

func handleCertificateDeletionError(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, err error) error {
	if common.IsResourceNotFound(err) {
		log.Printf("[INFO] deleting an unavailable certificate: %s", d.Id())
		return nil
	}

	err409, ok := err.(golangsdk.ErrDefault409)
	if !ok {
		return fmt.Errorf("error deleting certificate %s: %w", d.Id(), err)
	}
	var dep struct {
		ListenerIDs []string `json:"listener_ids"`
	}
	if err := json.Unmarshal(err409.Body, &dep); err != nil {
		return fmt.Errorf("error loading assigned listeners: %w", err)
	}

	mErr := new(multierror.Error)
	for _, listenerID := range dep.ListenerIDs {
		mErr = multierror.Append(mErr,
			unAssignCert(ctx, client, d.Id(), listenerID),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	log.Printf("[DEBUG] Retry deleting certificate %s", d.Id())

	if err := certificates.Delete(client, d.Id()).ExtractErr(); err != nil {
		return err
	}

	return nil
}

func unAssignCert(_ context.Context, client *golangsdk.ServiceClient, certID, listenerID string) error {
	listener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("failed to get listener %s: %w", listenerID, err)
	}

	otherCerts := make([]string, 0)
	for _, cert := range listener.SniContainerRefs {
		if cert != certID {
			otherCerts = append(otherCerts, cert)
		}
	}
	opts := listeners.UpdateOpts{
		SniContainerRefs: &otherCerts,
	}
	_, err = listeners.Update(client, listener.ID, opts).Extract()
	if err != nil {
		return fmt.Errorf("error unassigning certificate %s from listener %s: %w", certID, listener.ID, err)
	}
	return nil
}
