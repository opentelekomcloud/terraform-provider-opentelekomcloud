package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/metadata"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/protocols"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
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

			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"xaccount_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"domain_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metadata": {
							Type:      schema.TypeString,
							Required:  true,
							StateFunc: common.GetHashOrEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityProtocolV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := protocols.CreateOpts{
		MappingID: d.Get("mapping_id").(string),
	}
	_, err = protocols.Create(client, provider(d), protocol(d), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error registering protocol: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", provider(d), protocol(d)))

	if d.Get("metadata.#") != 0 {
		if err := uploadMetadata(d, meta); err != nil {
			return fmterr.Errorf("error uploading metadata: %w", err)
		}
	}

	clientCtx := ctxWithClient(ctx, client)
	return resourceIdentityProtocolV3Read(clientCtx, d, meta)
}

func resourceIdentityProtocolV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProviderAndProtocol(d); err != nil {
		return diag.FromErr(err)
	}
	protocol, err := protocols.Get(client, provider(d), protocol(d)).Extract()
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

	if d.Get("metadata.#") != 0 {
		if err := setMetadata(d, meta); err != nil {
			return fmterr.Errorf("error downloading metadata: %w", err)
		}
	}

	return nil
}

func resourceIdentityProtocolV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("mapping_id") {
		opts := protocols.UpdateOpts{
			MappingID: d.Get("mapping_id").(string),
		}
		_, err = protocols.Update(client, provider(d), protocol(d), opts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating protocol: %w", err)
		}
	}

	clientCtx := ctxWithClient(ctx, client)
	return resourceIdentityProtocolV3Read(clientCtx, d, meta)
}

func resourceIdentityProtocolV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromCtx(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	err = protocols.Delete(client, provider(d), protocol(d)).ExtractErr()
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

func uploadMetadata(d *schema.ResourceData, meta interface{}) error {
	client, err := identityExtClient(d, meta)
	if err != nil {
		return err
	}
	opts := metadata.ImportOpts{
		XAccountType: d.Get("metadata.0.xaccount_type").(string),
		DomainID:     d.Get("metadata.0.domain_id").(string),
		Metadata:     d.Get("metadata.0.metadata").(string),
	}
	err = metadata.Import(client, provider(d), protocol(d), opts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error importing metadata file: %w", err)
	}
	return nil
}

func setMetadata(d *schema.ResourceData, meta interface{}) error {
	client, err := identityExtClient(d, meta)
	if err != nil {
		return err
	}
	md, err := metadata.Get(client, provider(d), protocol(d)).Extract()
	if err != nil {
		return err
	}
	value := []map[string]interface{}{{
		"xaccount_type": md.XAccountType,
		"domain_id":     md.DomainID,
		"metadata":      common.GetHashOrEmpty(md.Data),
	}}
	return d.Set("metadata", value)
}

func provider(d *schema.ResourceData) string {
	return d.Get("provider_id").(string)
}

func protocol(d *schema.ResourceData) string {
	return d.Get("protocol").(string)
}
