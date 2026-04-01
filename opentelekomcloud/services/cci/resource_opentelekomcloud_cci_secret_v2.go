package cci

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/secret"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCISecretV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciSecretV2Create,
		ReadContext:   resourceCciSecretV2Read,
		UpdateContext: resourceCciSecretV2Update,
		DeleteContext: resourceCciSecretV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCciSecretV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"string_data": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"immutable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceCciSecretV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	createOpts := secret.CreateOpts{
		Namespace: d.Get("namespace").(string),
		Metadata: &secret.ObjectMeta{
			Name: d.Get("name").(string),
		},
		StringData: expandStringMap(d.Get("string_data")),
		Data:       expandStringMap(d.Get("data")),
		Type:       d.Get("type").(string),
	}

	resp, err := secret.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating CCI v2 Secret: %s", err)
	}

	ns := resp.Metadata.Namespace
	name := resp.Metadata.Name
	if ns == "" || name == "" {
		return diag.Errorf("unable to find namespace or CCI v2 Secret name from API response")
	}
	d.SetId(ns + "/" + name)

	return resourceCciSecretV2Read(ctx, d, meta)
}

func resourceCciSecretV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	resp, err := secret.Get(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 Secret")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("namespace", resp.Metadata.Namespace),
		d.Set("name", resp.Metadata.Name),
		d.Set("api_version", resp.APIVersion),
		d.Set("kind", resp.Kind),
		d.Set("annotations", resp.Metadata.Annotations),
		d.Set("labels", resp.Metadata.Labels),
		d.Set("creation_timestamp", resp.Metadata.CreationTimestamp),
		d.Set("resource_version", resp.Metadata.ResourceVersion),
		d.Set("uid", resp.Metadata.UID),
		d.Set("data", resp.Data),
		d.Set("type", resp.Type),
		d.Set("immutable", resp.Immutable),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCciSecretV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	updateOpts := secret.UpdateOpts{
		Metadata: &secret.ObjectMeta{
			Name: name,
		},
		StringData: expandStringMap(d.Get("string_data")),
		Data:       expandStringMap(d.Get("data")),
		Type:       d.Get("type").(string),
	}

	_, err = secret.Update(client, ns, name, updateOpts)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Secret: %s", err)
	}

	return resourceCciSecretV2Read(ctx, d, meta)
}

func resourceCciSecretV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	_, err = secret.Delete(client, secret.DeleteOpts{
		Namespace: ns,
		Name:      name,
	})
	if err != nil {
		return diag.Errorf("error deleting CCI v2 Secret: %s", err)
	}

	return nil
}

func resourceCciSecretV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("namespace", parts[0]),
		d.Set("name", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
