package obs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/s3"
)

func DataSourceObsBucket() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObsBucketRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versioning": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"logging": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"lifecycle_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"transition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"noncurrent_version_transition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"website": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_all_requests_to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routing_rules": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cors_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_headers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"server_side_encryption": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"event_notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						"filter_rule": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceObsBucketRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.NewObjectStorageClient(region)
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] Reading OBS bucket: %v", bucket)
	out, err := client.HeadBucket(bucket)
	if err != nil {
		return fmterr.Errorf("failed getting OBS bucket (%s): %w", bucket, err)
	}

	log.Printf("[DEBUG] Received OBS bucket: %v", out)

	d.SetId(bucket)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("bucket", bucket),
		d.Set("bucket_domain_name", s3.BucketDomainName(bucket, region)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting OBS bucket fields: %s", err)
	}

	// Read storage class
	if err := setObsBucketStorageClass(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the versioning
	if err := setObsBucketVersioning(client, d); err != nil {
		return diag.FromErr(err)
	}
	// Read the logging configuration
	if err := setObsBucketLogging(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the Lifecycle configuration
	if err := setObsBucketLifecycleConfiguration(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the website configuration
	if err := setObsBucketWebsiteConfiguration(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the CORS rules
	if err := setObsBucketCorsRules(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the tags
	if err := setObsBucketTags(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read SSE settings
	if err := setObsBucketEncryption(client, d); err != nil {
		return diag.FromErr(err)
	}

	// Read notifications settings
	if err := setObsBucketNotifications(client, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
