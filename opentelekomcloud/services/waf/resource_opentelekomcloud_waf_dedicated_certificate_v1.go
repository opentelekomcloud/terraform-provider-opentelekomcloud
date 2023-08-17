package waf

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/certificates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDedicatedCertificateV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedCertificateV1Create,
		ReadContext:   resourceWafDedicatedCertificateV1Read,
		DeleteContext: resourceWafDedicatedCertificateV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-_.]{1,64}$`), "Invalid certificate name. "+
					"Certificate name. The value can contain a maximum of 64 characters. Only digits, "+
					"letters, hyphens (-), underscores (_), and periods (.) are allowed."),
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: common.GetHashOrEmpty,
			},
			"key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: common.GetHashOrEmpty,
			},
			"expires": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWafDedicatedCertificateV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	opts := certificates.CreateOpts{
		Name:    d.Get("name").(string),
		Content: strings.TrimSpace(d.Get("content").(string)),
		Key:     strings.TrimSpace(d.Get("key").(string)),
	}

	certificate, err := certificates.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated Certificate: %w", err)
	}

	log.Printf("[DEBUG] Waf Dedicated Certificate created: %#v", certificate.ID)
	d.SetId(certificate.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedCertificateV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedCertificateV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	n, err := certificates.Get(client, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated Certificate: %w", err)
	}
	layout := "2006-01-02 15:04:05 MST"
	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("expires", time.Unix(n.ExpireAt/1000, 0).UTC().Format(layout)),
		d.Set("created_at", time.Unix(n.CreatedAt/1000, 0).UTC().Format(layout)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceWafDedicatedCertificateV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	err = certificates.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated Certificate: %s", err)
	}

	d.SetId("")
	return nil
}
