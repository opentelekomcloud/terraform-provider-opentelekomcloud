package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/protocols"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityProtocolV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProtocolV3Create,
		ReadContext:   resourceIdentityProtocolV3Read,
		UpdateContext: resourceIdentityProtocolV3Update,
		DeleteContext: resourceIdentityProtocolV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("provider_id", "protocol"),
		},

		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProtocolV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	clientCtx := ctxWithClient(ctx, client)

	provider := d.Get("provider_id").(string)
	protocol := d.Get("protocol").(string)

	opts := protocols.CreateOpts{
		MappingID: d.Get("mapping_id").(string),
	}
	_, err = protocols.Create(client, provider, protocol, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error registering protocol: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", provider, protocol))

	return resourceIdentityProtocolV3Read(clientCtx, d, meta)
}

func resourceIdentityProtocolV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProviderAndProtocol(d); err != nil {
		return diag.FromErr(err)
	}
	protocol, err := protocols.Get(client,
		d.Get("provider_id").(string),
		d.Get("protocol").(string),
	).Extract()
	if err != nil {
		return fmterr.Errorf("error reading protocol: %w", err)
	}
	mErr := multierror.Append(
		d.Set("mapping_id", protocol.MappingID),
		d.Set("links", protocol.Links),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting protocol attributes: %w", err)
	}

	return nil
}

func resourceIdentityProtocolV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	clientCtx := ctxWithClient(ctx, client)

	if d.HasChange("mapping_id") {
		provider := d.Get("provider_id").(string)
		protocol := d.Get("protocol").(string)

		opts := protocols.UpdateOpts{
			MappingID: d.Get("mapping_id").(string),
		}
		_, err = protocols.Update(client, provider, protocol, opts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating protocol: %w", err)
		}
	}

	return resourceIdentityProtocolV3Read(clientCtx, d, meta)
}

func resourceIdentityProtocolV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	err = protocols.Delete(client, d.Get("provider_id").(string), d.Get("protocol").(string)).ExtractErr()
	return diag.FromErr(err)
}

func setProviderAndProtocol(d *schema.ResourceData) error {
	parts := strings.SplitN(d.Id(), "/", 2)
	mErr := multierror.Append(
		d.Set("provider_id", parts[0]),
		d.Set("protocol", parts[1]),
	)
	return mErr.ErrorOrNil()
}
