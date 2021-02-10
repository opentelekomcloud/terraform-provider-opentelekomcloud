package obs

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	version = "2008-10-17"
)

func ResourceObsBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketPolicyPut,
		Read:   resourceObsBucketPolicyRead,
		Update: resourceObsBucketPolicyPut,
		Delete: resourceObsBucketPolicyDelete,

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

func resourceObsBucketPolicyPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	policy := d.Get("policy").(string)
	bucket := d.Get("bucket").(string)
	d.SetId(bucket)

	log.Printf("[DEBUG] OBS bucket: %s, put policy: %s", bucket, policy)

	params := &obs.SetBucketPolicyInput{
		Bucket: bucket,
		Policy: policy,
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
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
		return fmt.Errorf("error putting OBS policy: %s", err)
	}

	return nil
}

func resourceObsBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	log.Printf("[DEBUG] OBS bucket policy, read for bucket: %s", d.Id())
	pol, err := client.GetBucketPolicy(d.Id())

	if err != nil {
		return fmt.Errorf("error getting bucket policy")
	}

	if err := d.Set("policy", pol.Policy); err != nil {
		return err
	}

	return nil
}

func resourceObsBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] OBS bucket: %s, delete policy", bucket)
	_, err = client.DeleteBucketPolicy(bucket)

	if err != nil {
		if obsErr, ok := err.(obs.ObsError); ok && obsErr.Code == "NoSuchBucket" {
			return nil
		}
		return fmt.Errorf("error deleting OBS policy: %s", err)
	}
	return nil
}
