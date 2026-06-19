package cc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/central_network"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCcCentralNetworkV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcCentralNetworkV3Create,
		ReadContext:   resourceCcCentralNetworkV3Read,
		UpdateContext: resourceCcCentralNetworkV3Update,
		DeleteContext: resourceCcCentralNetworkV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_plane_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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

func resourceCcCentralNetworkV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := central_network.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
	}

	created, err := central_network.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud CC central network: %s", err)
	}

	d.SetId(created.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForCentralNetworkActive(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for CC central network (%s) to become available: %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, ccClientV3)
	return resourceCcCentralNetworkV3Read(clientCtx, d, meta)
}

func resourceCcCentralNetworkV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	cn, err := central_network.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud CC central network")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", cn.Name),
		d.Set("description", cn.Description),
		d.Set("enterprise_project_id", cn.EnterpriseProjectId),
		d.Set("state", cn.State),
		d.Set("default_plane_id", cn.DefaultPlaneId),
		d.Set("domain_id", cn.DomainId),
		d.Set("created_at", cn.CreatedAt),
		d.Set("updated_at", cn.UpdatedAt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCcCentralNetworkV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := central_network.UpdateOpts{
			CentralNetworkId: d.Id(),
			Name:             d.Get("name").(string),
			Description:      pointerto.String(d.Get("description").(string)),
		}
		if _, err = central_network.Update(client, updateOpts); err != nil {
			return diag.Errorf("error updating OpenTelekomCloud CC central network (%s): %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, ccClientV3)
	return resourceCcCentralNetworkV3Read(clientCtx, d, meta)
}

func resourceCcCentralNetworkV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if err = central_network.Delete(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud CC central network")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForCentralNetworkDeleted(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for CC central network (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func waitForCentralNetworkActive(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cn, err := central_network.Get(client, id)
		if err != nil {
			return nil, "", err
		}
		switch cn.State {
		case "AVAILABLE":
			return cn, "COMPLETED", nil
		case "FAILED", "DELETED":
			return cn, cn.State, fmt.Errorf("unexpected status (%s) for CC central network (%s)", cn.State, id)
		default:
			return cn, "PENDING", nil
		}
	}
}

func waitForCentralNetworkDeleted(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cn, err := central_network.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] OpenTelekomCloud CC central network (%s) has been deleted", id)
				return "deleted", "COMPLETED", nil
			}
			return nil, "", err
		}
		if cn.State == "DELETED" {
			return cn, "COMPLETED", nil
		}
		if cn.State == "FAILED" {
			return cn, cn.State, fmt.Errorf("unexpected status (%s) while deleting CC central network (%s)", cn.State, id)
		}
		return cn, "PENDING", nil
	}
}
