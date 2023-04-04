package vpc

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceNetworkingSecGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secgroup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingSecGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating networking client: %w", err)
	}

	listOpts := groups.ListOpts{
		ID:       d.Get("secgroup_id").(string),
		Name:     d.Get("name").(string),
		TenantID: d.Get("tenant_id").(string),
	}

	pages, err := groups.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to list security groups: %s", err)
	}

	allSecGroups, err := groups.ExtractGroups(pages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve security groups: %s", err)
	}

	var filteredGroups []groups.SecGroup
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, group := range allSecGroups {
			if r.MatchString(group.Name) {
				filteredGroups = append(filteredGroups, group)
			}
		}
		if len(filteredGroups) < 1 {
			return fmterr.Errorf("no Security Group found by regex: %s", d.Get("name_regex"))
		}
		if len(filteredGroups) > 1 {
			return fmterr.Errorf("more than one Security Group with regex: %s", d.Get("name_regex"))
		}
		allSecGroups = filteredGroups
	}

	if len(allSecGroups) < 1 {
		return fmterr.Errorf("no Security Group found with name: %s", d.Get("name"))
	}

	if len(allSecGroups) > 1 {
		return fmterr.Errorf("more than one Security Group found with name: %s", d.Get("name"))
	}

	secGroup := allSecGroups[0]

	log.Printf("[DEBUG] Retrieved Security Group %s: %+v", secGroup.ID, secGroup)
	d.SetId(secGroup.ID)

	mErr := multierror.Append(
		d.Set("name", secGroup.Name),
		d.Set("description", secGroup.Description),
		d.Set("tenant_id", secGroup.TenantID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting security group fields: %w", err)
	}

	return nil
}
