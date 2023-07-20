package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/metadata"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/protocols"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/providers"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{protocolSAML, protocolOIDC}, false),
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
			"access_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"metadata"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"program", "program_console"}, false),
						},
						"provider_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"signing_key": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := common.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"authorization_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"scopes": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"response_type": {
							Type: schema.TypeString,
							// Computed: true,
							Optional: true,
							Default:  "id_token",
						},
						"response_mode": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"fragment", "form_post"}, false),
						},
					},
				},
			},
			"metadata": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"access_config"},
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

	if d.Get("metadata.#") != 0 && protocol(d) == protocolSAML {
		if err := uploadMetadata(d, meta); err != nil {
			return fmterr.Errorf("error uploading metadata: %w", err)
		}
	}

	if ac, ok := d.GetOk("access_config"); ok && protocol(d) == protocolOIDC {
		// Create access config for oidc provider.
		config := meta.(*cfg.Config)
		clientAdmin, err := config.IdentityV30Client()
		if err != nil {
			return fmterr.Errorf(clientCreationFail, err)
		}
		accessConfigArr := ac.([]interface{})
		accessConfig := accessConfigArr[0].(map[string]interface{})

		accessType := accessConfig["access_type"].(string)
		createAccessTypeOpts := providers.CreateOIDCOpts{
			AccessMode: accessType,
			IdpUrl:     accessConfig["provider_url"].(string),
			ClientId:   accessConfig["client_id"].(string),
			SigningKey: accessConfig["signing_key"].(string),
			IdpIp:      provider(d),
		}

		if accessType == "program_console" {
			scopes := common.ExpandToStringSlice(accessConfig["scopes"].([]interface{}))
			createAccessTypeOpts.Scope = strings.Join(scopes, scopeSpilt)
			createAccessTypeOpts.AuthEndpoint = accessConfig["authorization_endpoint"].(string)
			createAccessTypeOpts.ResponseType = accessConfig["response_type"].(string)
			createAccessTypeOpts.ResponseMode = accessConfig["response_mode"].(string)
		}
		log.Printf("[DEBUG] Create access type of provider: %#v", opts)

		_, err = providers.CreateOIDC(clientAdmin, createAccessTypeOpts)
		if err != nil {
			return fmterr.Errorf("Error creating the provider access config: %s", err)
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
	prot, err := protocols.Get(client, provider(d), protocol(d)).Extract()
	if err != nil {
		return fmterr.Errorf("error reading protocol: %w", err)
	}

	mErr := multierror.Append(
		d.Set("mapping_id", prot.MappingID),
		d.Set("links", prot.Links),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting protocol attributes: %w", err)
	}

	if d.Get("metadata.#") != 0 {
		if err := setMetadata(d, meta); err != nil {
			return fmterr.Errorf("error downloading metadata: %w", err)
		}
	}

	if protocol(d) == protocolOIDC {
		config := meta.(*cfg.Config)
		clientOIDC, err := config.IdentityV30Client()
		if err != nil {
			return fmterr.Errorf(clientCreationFail, err)
		}
		accessType, err := providers.GetOIDC(clientOIDC, d.Id())
		if err == nil {
			scopes := strings.Split(accessType.Scope, scopeSpilt)
			accessTypeConfig := []interface{}{
				map[string]interface{}{
					"access_type":            accessType.AccessMode,
					"provider_url":           accessType.IdpUrl,
					"client_id":              accessType.ClientId,
					"signing_key":            accessType.SigningKey,
					"scopes":                 scopes,
					"response_mode":          accessType.ResponseMode,
					"authorization_endpoint": accessType.AuthEndpoint,
					"response_type":          accessType.ResponseType,
				},
			}

			mErr = multierror.Append(
				mErr,
				d.Set("access_config", accessTypeConfig),
			)
		}
	}
	if err = mErr.ErrorOrNil(); err != nil {
		log.Printf("[ERROR] Error setting identity protocol attributes %s: %s", d.Id(), err)
		return fmterr.Errorf("Error setting identity protocol attributes: %s", err)
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

	if d.HasChange("access_config") && d.Get("protocol") == protocolOIDC {
		err = updateAccessConfig(d, meta)
		if err != nil {
			return fmterr.Errorf("error updating access_config: %w", err)
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
