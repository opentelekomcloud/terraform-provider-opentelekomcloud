package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/quota"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceTaurusDBV3Quota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBV3QuotaCreateOrUpdate,
		UpdateContext: resourceTaurusDBV3QuotaCreateOrUpdate,
		ReadContext:   resourceTaurusDBV3QuotaRead,
		DeleteContext: resourceTaurusDBV3QuotaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"vcpus_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"ram_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"availability_instance_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_vcpus_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_ram_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTaurusDBV3QuotaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	opts := quota.SetQuotasOpts{
		QuotaList: []quota.SetQuota{
			{
				EnterpriseProjectId:   d.Get("enterprise_project_id").(string),
				EnterpriseProjectName: d.Get("enterprise_project_name").(string),
				InstanceQuota:         d.Get("instance_quota").(int),
				VcpusQuota:            d.Get("vcpus_quota").(int),
				RamQuota:              d.Get("ram_quota").(int),
			},
		},
	}

	_, err = quota.SetQuotas(client, opts)
	if err != nil {
		return diag.Errorf("error updating TaurusDB quota: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(d.Get("enterprise_project_id").(string))
	}

	return resourceTaurusDBV3QuotaRead(ctx, d, meta)
}

func resourceTaurusDBV3QuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	opts := quota.ListQuotasOpts{}
	resp, err := quota.ListQuotas(client, opts)
	if err != nil {
		return diag.Errorf("error reading TaurusDB quotas: %s", err)
	}

	var foundQuota *quota.Quota
	for _, q := range resp.QuotaList {
		if q.EnterpriseProjectId == d.Id() {
			foundQuota = &q
			break
		}
	}

	if foundQuota == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "TaurusDB quota")
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("enterprise_project_id", foundQuota.EnterpriseProjectId),
		d.Set("enterprise_project_name", foundQuota.EnterpriseProjectName),
		d.Set("instance_quota", foundQuota.InstanceQuota),
		d.Set("vcpus_quota", foundQuota.VcpusQuota),
		d.Set("ram_quota", foundQuota.RamQuota),
		d.Set("availability_instance_quota", foundQuota.AvailabilityInstanceQuota),
		d.Set("availability_vcpus_quota", foundQuota.AvailabilityVcpusQuota),
		d.Set("availability_ram_quota", foundQuota.AvailabilityRamQuota),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDBV3QuotaDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting TaurusDB quota resource is not supported. The TaurusDB quota resource is only " +
		"removed from the state, the TaurusDB quota remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
