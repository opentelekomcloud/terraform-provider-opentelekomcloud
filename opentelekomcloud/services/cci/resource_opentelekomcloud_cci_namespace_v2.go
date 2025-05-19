package cci

import (
	"context"
	"fmt"
	"time"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/namespace"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCINamespaceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciNamespaceV2Create,
		ReadContext:   resourceCciNamespaceV2Read,
		DeleteContext: resourceCciNamespaceV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"finalizers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCciNamespaceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := namespace.CreateOpts{
		Metadata: namespace.Metadata{
			Name: d.Get("name").(string),
		},
	}

	resp, err := namespace.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating CCI namespace: %s", err)
	}

	d.SetId(resp.Metadata.Name)

	err = waitForCreateNamespaceStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceCciNamespaceV2Read(ctx, d, meta)
}

func waitForCreateNamespaceStatus(ctx context.Context, client *golangsdk.ServiceClient, ns string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Active"},
		Refresh: func() (interface{}, string, error) {
			resp, err := namespace.Get(client, ns)
			if err != nil {
				return nil, "failed", err
			}
			return resp, resp.Status.Phase, nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the namespace to complete active timeout: %s", err)
	}
	return nil
}

func resourceCciNamespaceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	resp, err := namespace.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies namespace form server")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("api_version", resp.APIVersion),
		d.Set("kind", resp.Kind),
		d.Set("name", resp.Metadata.Name),
		d.Set("annotations", resp.Metadata.Annotations),
		d.Set("labels", resp.Metadata.Labels),
		d.Set("creation_timestamp", resp.Metadata.CreationTimestamp),
		d.Set("resource_version", resp.Metadata.ResourceVersion),
		d.Set("uid", resp.Metadata.UID),
		d.Set("finalizers", resp.Spec.Finalizers),
		d.Set("status", resp.Status.Phase),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCciNamespaceV2Update(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCciNamespaceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	_, err = namespace.Delete(client, namespace.DeleteOpts{
		Name: d.Id(),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForDeleteNamespaceStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForDeleteNamespaceStatus(ctx context.Context, client *golangsdk.ServiceClient, ns string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Active", "Terminating"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := namespace.Get(client, ns)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "", "DELETED", nil
				}
				return nil, "ERROR", err
			}
			return resp, resp.Status.Phase, nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the namespace to complete delete timeout: %s", err)
	}
	return nil
}
