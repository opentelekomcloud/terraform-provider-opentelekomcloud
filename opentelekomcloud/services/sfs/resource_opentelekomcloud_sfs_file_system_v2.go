package sfs

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs/v2/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSFSFileSystemV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSFileSystemV2Create,
		ReadContext:   resourceSFSFileSystemV2Read,
		UpdateContext: resourceSFSFileSystemV2Update,
		DeleteContext: resourceSFSFileSystemV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NFS",
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"access_level": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_to"},
				Deprecated:   "Use the opentelekomcloud_sfs_share_access_rule_v2 resource instead",
			},
			"access_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Default:    "cert",
				Deprecated: "Use the opentelekomcloud_sfs_share_access_rule_v2 resource instead",
			},
			"access_to": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_level"},
				Deprecated:   "Use the opentelekomcloud_sfs_share_access_rule_v2 resource instead",
			},
			"share_access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_rule_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceSFSMetadataV2(d *schema.ResourceData) map[string]string {
	meta := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		meta[key] = val.(string)
	}
	return meta
}

func resourceSFSFileSystemV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share Client: %s", err)
	}

	createOpts := shares.CreateOpts{
		ShareProto:       d.Get("share_proto").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		IsPublic:         d.Get("is_public").(bool),
		Metadata:         resourceSFSMetadataV2(d),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	share, err := shares.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    waitForSFSFileStatus(client, share.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error creating share file: %s", err)
	}

	accessLevel := d.Get("access_level").(string)
	accessTo := d.Get("access_to").(string)
	if accessLevel != "" || accessTo != "" {
		grantAccessOpts := shares.GrantAccessOpts{
			AccessLevel: d.Get("access_level").(string),
			AccessType:  d.Get("access_type").(string),
			AccessTo:    d.Get("access_to").(string),
		}

		_, err = shares.GrantAccess(client, share.ID, grantAccessOpts).ExtractAccess()
		if err != nil {
			return fmterr.Errorf("error applying access rules to share file: %s", err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "sfs", share.ID, tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of SFS File System: %s", err)
		}
	}

	d.SetId(share.ID)

	return resourceSFSFileSystemV2Read(ctx, d, meta)
}

func resourceSFSFileSystemV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share: %s", err)
	}

	share, err := shares.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Shares: %s", err)
	}
	mErr := multierror.Append(nil,
		d.Set("name", share.Name),
		d.Set("share_proto", share.ShareProto),
		d.Set("status", share.Status),
		d.Set("size", share.Size),
		d.Set("description", share.Description),
		d.Set("share_type", share.ShareType),
		d.Set("volume_type", share.VolumeType),
		d.Set("is_public", share.IsPublic),
		d.Set("availability_zone", share.AvailabilityZone),
		d.Set("region", config.GetRegion(d)),
		d.Set("export_location", share.ExportLocation),
		d.Set("host", share.Host),
	)

	// NOTE: This tries to remove system metadata.
	metadata := make(map[string]string)
	for key, val := range share.Metadata {
		if strings.HasPrefix(key, "#sfs") {
			continue
		}
		if strings.Contains(key, "enterprise_project_id") || strings.Contains(key, "share_used") {
			continue
		}
		metadata[key] = val
	}
	if err := d.Set("metadata", metadata); err != nil {
		return diag.FromErr(err)
	}

	// save tags
	resourceTags, err := tags.Get(client, "sfs", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud SFS File System tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud SFS File System: %s", err)
	}

	rules, err := shares.ListAccessRights(client, d.Id()).ExtractAccessRights()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Shares: %s", err)
	}

	attachResourceID := d.Get("access_to").(string)
	if attachResourceID != "" {
		for _, rule := range rules {
			if rule.ID == attachResourceID {
				mErr = multierror.Append(mErr,
					d.Set("share_access_id", rule.ID),
					d.Set("access_rule_status", rule.State),
					d.Set("access_to", rule.AccessTo),
					d.Set("access_type", rule.AccessType),
					d.Set("access_level", rule.AccessLevel),
				)
				break
			}
		}
	}

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceSFSFileSystemV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Share File: %s", err)
	}
	var updateOpts shares.UpdateOpts

	if d.HasChange("description") || d.HasChange("name") {
		updateOpts.DisplayName = d.Get("name").(string)
		updateOpts.DisplayDescription = d.Get("description").(string)

		_, err = shares.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud Share File: %s", err)
		}
	}
	if d.HasChange("access_to") || d.HasChange("access_level") || d.HasChange("access_type") {
		shareAccessID := d.Get("share_access_id").(string)
		if shareAccessID != "" {
			deleteAccessOpts := shares.DeleteAccessOpts{AccessID: d.Get("share_access_id").(string)}
			if err := shares.DeleteAccess(client, d.Id(), deleteAccessOpts).Err; err != nil {
				return fmterr.Errorf("error changing access rules for share file: %s", err)
			}
		}

		accessLevel := d.Get("access_level").(string)
		accessTo := d.Get("access_to").(string)
		if accessTo != "" || accessLevel != "" {
			grantAccessOpts := shares.GrantAccessOpts{
				AccessLevel: d.Get("access_level").(string),
				AccessType:  d.Get("access_type").(string),
				AccessTo:    d.Get("access_to").(string),
			}

			log.Printf("[DEBUG] Grant Access Rules: %#v", grantAccessOpts)
			_, err := shares.GrantAccess(client, d.Id(), grantAccessOpts).ExtractAccess()
			if err != nil {
				return fmterr.Errorf("error changing access rules for share file: %s", err)
			}
		}
	}

	if d.HasChange("size") {
		oldSizeRaw, newSizeRaw := d.GetChange("size")
		newSize := newSizeRaw.(int)
		if oldSizeRaw.(int) < newSize {
			expandOpts := shares.ExpandOpts{OSExtend: shares.OSExtendOpts{NewSize: newSize}}
			if err := shares.Expand(client, d.Id(), expandOpts).ExtractErr(); err != nil {
				return fmterr.Errorf("error expanding OpenTelekomCloud Share File size: %s", err)
			}
		} else {
			shrinkOpts := shares.ShrinkOpts{OSShrink: shares.OSShrinkOpts{NewSize: newSize}}
			if err := shares.Shrink(client, d.Id(), shrinkOpts).ExtractErr(); err != nil {
				return fmterr.Errorf("error shrinking OpenTelekomCloud Share File size: %s", err)
			}
		}
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "sfs", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of SFS File System %s: %s", d.Id(), err)
		}
	}

	return resourceSFSFileSystemV2Read(ctx, d, meta)
}

func resourceSFSFileSystemV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Shared File: %s", err)
	}
	err = shares.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Shared File: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSFileStatus(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Share File: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSFSFileStatus(client *golangsdk.ServiceClient, shareID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		share, err := shares.Get(client, shareID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud shared File %s", shareID)
				return share, "deleted", nil
			}
			return nil, "", err
		}
		return share, share.Status, nil
	}
}
