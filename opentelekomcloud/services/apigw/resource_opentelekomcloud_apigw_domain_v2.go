package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/domain"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIDomainV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainV2Create,
		ReadContext:   resourceDomainV2Read,
		UpdateContext: resourceDomainV2Update,
		DeleteContext: resourceDomainV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"min_ssl_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TLSv1.1", "TLSv1.2",
				}, false),
			},
			"http_redirect_to_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDomainV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	redirectToHTTPS := d.Get("http_redirect_to_https").(bool)
	opts := domain.CreateOpts{
		GatewayID:             d.Get("gateway_id").(string),
		GroupID:               d.Get("group_id").(string),
		UrlDomain:             d.Get("name").(string),
		MinSslVersion:         d.Get("min_ssl_version").(string),
		IsHttpRedirectToHttps: &redirectToHTTPS,
	}

	resp, err := domain.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating APIGW domain: %s", err)
	}
	d.SetId(resp.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDomainV2Read(clientCtx, d, meta)
}

func resourceDomainV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	groupId := d.Get("group_id").(string)
	resp, err := GetDomain(client, gatewayId, groupId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "APIGW domain")
	}

	mErr := multierror.Append(nil,
		d.Set("gateway_id", gatewayId),
		d.Set("group_id", groupId),
		d.Set("min_ssl_version", resp.MinSslVersion),
		d.Set("status", resp.CnameStatus),
	)
	// The group query does not return the domain name or is_http_redirect_to_https
	// (url_domains[].name is empty), so both are preserved from the configuration
	// in state and are intentionally not set here.
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving APIGW domain (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceDomainV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if d.HasChanges("min_ssl_version", "http_redirect_to_https") {
		redirectToHTTPS := d.Get("http_redirect_to_https").(bool)
		opts := domain.UpdateOpts{
			GatewayID:             d.Get("gateway_id").(string),
			GroupID:               d.Get("group_id").(string),
			DomainID:              d.Id(),
			MinSslVersion:         d.Get("min_ssl_version").(string),
			IsHttpRedirectToHttps: &redirectToHTTPS,
		}
		_, err = domain.Update(client, opts)
		if err != nil {
			return diag.Errorf("error updating APIGW domain (%s): %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDomainV2Read(clientCtx, d, meta)
}

func resourceDomainV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := domain.DeleteOpts{
		GatewayID: d.Get("gateway_id").(string),
		GroupID:   d.Get("group_id").(string),
		DomainID:  d.Id(),
	}
	err = domain.Delete(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting APIGW domain (%s): %s", d.Id(), err))
	}
	return nil
}

func resourceDomainV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <gateway_id>/<group_id>/<domain_id>")
	}

	d.SetId(parts[2])
	return []*schema.ResourceData{d}, multierror.Append(
		d.Set("gateway_id", parts[0]),
		d.Set("group_id", parts[1]),
	).ErrorOrNil()
}

func GetDomain(client *golangsdk.ServiceClient, gatewayId, groupId, domainId string) (*group.UrlDomains, error) {
	resp, err := group.Get(client, gatewayId, groupId)
	if err != nil {
		return nil, err
	}
	for i := range resp.UrlDomains {
		if resp.UrlDomains[i].DomainId == domainId {
			return &resp.UrlDomains[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}
