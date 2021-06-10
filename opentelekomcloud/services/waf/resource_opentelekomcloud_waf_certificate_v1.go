package waf

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/certificates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafCertificateV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafCertificateV1Create,
		ReadContext:   resourceWafCertificateV1Read,
		UpdateContext: resourceWafCertificateV1Update,
		DeleteContext: resourceWafCertificateV1Delete,
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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateCertificateName,
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
		},
	}
}

func resourceWafCertificateV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(wafClientError, err)
	}

	createOpts := certificates.CreateOpts{
		Name:    d.Get("name").(string),
		Content: strings.TrimSpace(d.Get("content").(string)),
		Key:     strings.TrimSpace(d.Get("key").(string)),
	}

	certificate, err := certificates.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Certificate: %w", err)
	}

	log.Printf("[DEBUG] Waf certificate created: %#v", certificate)
	d.SetId(certificate.Id)

	return resourceWafCertificateV1Read(ctx, d, meta)
}

func resourceWafCertificateV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(wafClientError, err)
	}
	n, err := certificates.Get(wafClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Certificate: %w", err)
	}

	expires := time.Unix(int64(n.ExpireTime/1000), 0).UTC().Format("2006-01-02 15:04:05 MST")

	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("expires", expires),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceWafCertificateV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(wafClientError, err)
	}

	var updateOpts certificates.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	_, err = certificates.Update(wafClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF Certificate: %w", err)
	}
	return resourceWafCertificateV1Read(ctx, d, meta)
}

func resourceWafCertificateV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(wafClientError, err)
	}

	err = certificates.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Certificate: %s", err)
	}

	d.SetId("")
	return nil
}

var validNameRe = regexp.MustCompile(`^[\w-]{1,256}$`)

func validateCertificateName(v interface{}, path cty.Path) diag.Diagnostics {
	name := v.(string)
	if validNameRe.MatchString(name) {
		return nil
	}
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       fmt.Sprintf("invalid certificate name: %s", name),
		Detail:        "The maximum length is 256 characters. Only digits, letters, underscores (_), and hyphens (-) are allowed.",
		AttributePath: path,
	}}
}
