package cci

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/configmap"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCIConfigMapV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciConfigMapV2Create,
		ReadContext:   resourceCciConfigMapV2Read,
		UpdateContext: resourceCciConfigMapV2Update,
		DeleteContext: resourceCciConfigMapV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCciConfigMapV2ImportState,
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
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"binary_data": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"immutable": {
				Type:     schema.TypeBool,
				Computed: true,
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
		},
	}
}

func resourceCciConfigMapV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	createOpts := configmap.CreateOpts{
		Metadata: configmap.ObjectMeta{
			Name: d.Get("name").(string),
		},
		Data:       expandStringMap(d.Get("data")),
		BinaryData: expandStringMap(d.Get("binary_data")),
	}

	resp, err := configmap.Create(client, ns, createOpts)
	if err != nil {
		return diag.Errorf("error creating CCI v2 ConfigMap: %s", err)
	}

	respNs := resp.Metadata.Namespace
	respName := resp.Metadata.Name
	if respNs == "" || respName == "" {
		return diag.Errorf("unable to find namespace or CCI v2 ConfigMap name from API response")
	}
	d.SetId(respNs + "/" + respName)

	return resourceCciConfigMapV2Read(ctx, d, meta)
}

func resourceCciConfigMapV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	resp, err := configmap.Get(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 ConfigMap")
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
		d.Set("binary_data", resp.BinaryData),
		d.Set("immutable", resp.Immutable),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCciConfigMapV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	updateOpts := configmap.UpdateOpts{
		Metadata: &configmap.ObjectMeta{
			Name: name,
		},
		Data:       expandStringMap(d.Get("data")),
		BinaryData: expandStringMap(d.Get("binary_data")),
	}

	_, err = configmap.Update(client, ns, name, updateOpts)
	if err != nil {
		return diag.Errorf("error updating CCI v2 ConfigMap: %s", err)
	}

	return resourceCciConfigMapV2Read(ctx, d, meta)
}

func resourceCciConfigMapV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	_, err = configmap.Delete(client, ns, name, configmap.DeleteOpts{})
	if err != nil {
		return diag.Errorf("error deleting CCI v2 ConfigMap (%s/%s): %s", ns, name, err)
	}

	return nil
}

func resourceCciConfigMapV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
