package obs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceObsBucketPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketPolicyPut,
		ReadContext:   resourceObsBucketPolicyRead,
		UpdateContext: resourceObsBucketPolicyPut,
		DeleteContext: resourceObsBucketPolicyDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     common.ValidateJsonString,
				DiffSuppressFunc: common.SuppressEquivalentAwsPolicyDiffs,
			},
		},
	}
}

func resourceObsBucketPolicyPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	policy := d.Get("policy").(string)
	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] OBS bucket: %s, put policy: %s", bucket, policy)

	params := &obs.SetBucketPolicyInput{
		Bucket: bucket,
		Policy: policy,
	}

	err = resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		if _, err := client.SetBucketPolicy(params); err != nil {
			if err, ok := err.(obs.ObsError); ok {
				if err.Code == "MalformedPolicy" {
					return resource.RetryableError(err)
				}
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmterr.Errorf("error putting OBS policy: %s", err)
	}

	d.SetId(bucket)

	return nil
}

func resourceObsBucketPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	log.Printf("[DEBUG] OBS bucket policy, read for bucket: %s", d.Id())
	pol, err := client.GetBucketPolicy(d.Id())

	if err != nil {
		return fmterr.Errorf("error getting bucket policy")
	}

	if err := d.Set("policy", pol.Policy); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceObsBucketPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] OBS bucket: %s, delete policy", bucket)
	_, err = client.DeleteBucketPolicy(bucket)

	if err != nil {
		if obsErr, ok := err.(obs.ObsError); ok && obsErr.Code == "NoSuchBucket" {
			return nil
		}
		return fmterr.Errorf("error deleting OBS policy: %s", err)
	}
	return nil
}
