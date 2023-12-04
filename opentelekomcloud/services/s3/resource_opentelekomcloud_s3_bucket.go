package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceS3Bucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceS3BucketCreate,
		ReadContext:   resourceS3BucketRead,
		UpdateContext: resourceS3BucketUpdate,
		DeleteContext: resourceS3BucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceS3BucketImportState,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  common.ValidateName,
				ConflictsWith: []string{"bucket_prefix"},
			},
			"bucket_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(0, 63-resource.UniqueIDSuffixLength),
				ConflictsWith: []string{"bucket"},
			},
			"bucket_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Default:  "private",
				Optional: true,
			},
			"policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     common.ValidateJsonString,
				DiffSuppressFunc: common.SuppressEquivalentAwsPolicyDiffs,
			},
			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"website": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"error_document": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"redirect_all_requests_to": {
							Type: schema.TypeString,
							ConflictsWith: []string{
								"website.0.index_document",
								"website.0.error_document",
								"website.0.routing_rules",
							},
							Optional: true,
						},
						"routing_rules": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: common.ValidateJsonString,
							StateFunc: func(v interface{}) string {
								jsonString, _ := common.NormalizeJsonString(v)
								return jsonString
							},
						},
					},
				},
			},
			"hosted_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"website_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"website_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"mfa_delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"logging": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["target_bucket"]))
					buf.WriteString(fmt.Sprintf("%s-", m["target_prefix"]))
					return hashcode.String(buf.String())
				},
			},
			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateS3BucketLifecycleRuleId,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"abort_incomplete_multipart_upload_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      ExpirationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: common.ValidateRFC3339Timestamp,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateS3BucketLifecycleExpirationDays,
									},
									"expired_object_delete_marker": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      ExpirationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateS3BucketLifecycleExpirationDays,
									},
								},
							},
						},
					},
				},
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceS3BucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.S3Client(region)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud S3 client: %s", err)
	}

	// Get the bucket and acl
	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else if v, ok := d.GetOk("bucket_prefix"); ok {
		bucket = resource.PrefixedUniqueId(v.(string))
	} else {
		bucket = resource.UniqueId()
	}
	_ = d.Set("bucket", bucket)
	acl := d.Get("acl").(string)

	log.Printf("[DEBUG] S3 bucket create: %s, ACL: %s", bucket, acl)

	createOpts := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}

	log.Printf("[DEBUG] S3 bucket create: %s, using region: %s", bucket, region)

	err = resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to create new S3 bucket: %q", bucket)
		_, err := client.CreateBucket(createOpts)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "OperationAborted" {
				log.Printf("[WARN] Got an error while trying to create S3 bucket %s: %s", bucket, err)
				return resource.RetryableError(fmt.Errorf("error creating S3 bucket %s, retrying: %s", bucket, err))
			}
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmterr.Errorf("error creating S3 bucket: %s", err)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceS3BucketUpdate(ctx, d, meta)
}

func resourceS3BucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud S3 client: %s", err)
	}

	if err := setTagsS3(ctx, client, d); err != nil {
		return fmterr.Errorf("%q: %s", d.Get("bucket").(string), err)
	}

	if d.HasChange("policy") {
		if err := resourceS3BucketPolicyUpdate(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("cors_rule") {
		if err := resourceS3BucketCorsUpdate(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("website") {
		if err := resourceS3BucketWebsiteUpdate(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("versioning") {
		if d.Get("versioning.#").(int) > 0 {
			enabled := d.Get("versioning.0.enabled").(bool)

			if enabled || !d.IsNewResource() {
				if err := resourceS3BucketVersioningUpdate(ctx, client, d); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("acl") {
		if err := resourceS3BucketAclUpdate(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("logging") {
		if err := resourceS3BucketLoggingUpdate(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceAwsS3BucketLifecycleUpdate(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceS3BucketRead(ctx, d, meta)
}

func resourceS3BucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud S3 client: %s", err)
	}

	_, err = retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(d.Id()),
		})
	})
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN] S3 Bucket (%s) not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		} else {
			// some of the AWS SDK's errors can be empty strings, so let's add
			// some additional context.
			return fmterr.Errorf("error reading S3 bucket %s: %s", d.Id(), err)
		}
	}

	// In the import case, we won't have this
	if _, ok := d.GetOk("bucket"); !ok {
		_ = d.Set("bucket", d.Id())
	}
	domainName := BucketDomainName(d.Get("bucket").(string), config.GetRegion(d))
	_ = d.Set("bucket_domain_name", domainName)

	// Read the policy
	if _, ok := d.GetOk("policy"); ok {
		pol, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
			return client.GetBucketPolicy(&s3.GetBucketPolicyInput{
				Bucket: aws.String(d.Id()),
			})
		})
		log.Printf("[DEBUG] S3 bucket: %s, read policy: %v", d.Id(), pol)
		if err != nil {
			if err := d.Set("policy", ""); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if v := pol.(*s3.GetBucketPolicyOutput).Policy; v == nil {
				if err := d.Set("policy", ""); err != nil {
					return diag.FromErr(err)
				}
			} else {
				policy, err := common.NormalizeJsonString(*v)
				if err != nil {
					return fmterr.Errorf("policy contains an invalid JSON: %w", err)
				}
				_ = d.Set("policy", policy)
			}
		}
	}

	// Read the CORS
	corsResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.GetBucketCors(&s3.GetBucketCorsInput{
			Bucket: aws.String(d.Id()),
		})
	})
	cors := corsResponse.(*s3.GetBucketCorsOutput)
	if err != nil {
		// An S3 Bucket might not have CORS configuration set.
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() != "NoSuchCORSConfiguration" {
			return diag.FromErr(err)
		}
		log.Printf("[WARN] S3 bucket: %s, no CORS configuration could be found.", d.Id())
	}
	log.Printf("[DEBUG] S3 bucket: %s, read CORS: %v", d.Id(), cors)
	if cors.CORSRules != nil {
		rules := make([]map[string]interface{}, 0, len(cors.CORSRules))
		for _, ruleObject := range cors.CORSRules {
			rule := make(map[string]interface{})
			rule["allowed_headers"] = common.FlattenStringList(ruleObject.AllowedHeaders)
			rule["allowed_methods"] = common.FlattenStringList(ruleObject.AllowedMethods)
			rule["allowed_origins"] = common.FlattenStringList(ruleObject.AllowedOrigins)
			// Both the "ExposeHeaders" and "MaxAgeSeconds" might not be set.
			if ruleObject.AllowedOrigins != nil {
				rule["expose_headers"] = common.FlattenStringList(ruleObject.ExposeHeaders)
			}
			if ruleObject.MaxAgeSeconds != nil {
				rule["max_age_seconds"] = int(*ruleObject.MaxAgeSeconds)
			}
			rules = append(rules, rule)
		}
		if err := d.Set("cors_rule", rules); err != nil {
			return diag.FromErr(err)
		}
	}

	// Read the website configuration
	wsResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(d.Id()),
		})
	})
	ws := wsResponse.(*s3.GetBucketWebsiteOutput)
	var websites []map[string]interface{}
	if err == nil {
		w := make(map[string]interface{})

		if v := ws.IndexDocument; v != nil {
			w["index_document"] = *v.Suffix
		}

		if v := ws.ErrorDocument; v != nil {
			w["error_document"] = *v.Key
		}

		if v := ws.RedirectAllRequestsTo; v != nil {
			if v.Protocol == nil {
				w["redirect_all_requests_to"] = *v.HostName
			} else {
				var host string
				var path string
				parsedHostName, err := url.Parse(*v.HostName)
				if err == nil {
					host = parsedHostName.Host
					path = parsedHostName.Path
				} else {
					host = *v.HostName
					path = ""
				}

				w["redirect_all_requests_to"] = (&url.URL{
					Host:   host,
					Path:   path,
					Scheme: *v.Protocol,
				}).String()
			}
		}

		if v := ws.RoutingRules; v != nil {
			rr, err := normalizeRoutingRules(v)
			if err != nil {
				return fmterr.Errorf("error while marshaling routing rules: %s", err)
			}
			w["routing_rules"] = rr
		}

		websites = append(websites, w)
	}
	if err := d.Set("website", websites); err != nil {
		return diag.FromErr(err)
	}

	versioningResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.GetBucketVersioning(&s3.GetBucketVersioningInput{
			Bucket: aws.String(d.Id()),
		})
	})
	versioning := versioningResponse.(*s3.GetBucketVersioningOutput)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] S3 Bucket: %s, versioning: %v", d.Id(), versioning)
	if versioning != nil {
		vcl := make([]map[string]interface{}, 0, 1)
		vc := make(map[string]interface{})
		if versioning.Status != nil && *versioning.Status == s3.BucketVersioningStatusEnabled {
			vc["enabled"] = true
		} else {
			vc["enabled"] = false
		}

		if versioning.MFADelete != nil && *versioning.MFADelete == s3.MFADeleteEnabled {
			vc["mfa_delete"] = true
		} else {
			vc["mfa_delete"] = false
		}
		vcl = append(vcl, vc)
		if err := d.Set("versioning", vcl); err != nil {
			return diag.FromErr(err)
		}
	}

	// Read the logging configuration
	loggingResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.GetBucketLogging(&s3.GetBucketLoggingInput{
			Bucket: aws.String(d.Id()),
		})
	})
	logging := loggingResponse.(*s3.GetBucketLoggingOutput)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] S3 Bucket: %s, logging: %v", d.Id(), logging)
	lcl := make([]map[string]interface{}, 0, 1)
	if v := logging.LoggingEnabled; v != nil {
		lc := make(map[string]interface{})
		if *v.TargetBucket != "" {
			lc["target_bucket"] = *v.TargetBucket
		}
		if *v.TargetPrefix != "" {
			lc["target_prefix"] = *v.TargetPrefix
		}
		lcl = append(lcl, lc)
	}
	if err := d.Set("logging", lcl); err != nil {
		return diag.FromErr(err)
	}

	// Read the lifecycle configuration

	lifecycleResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(d.Id()),
		})
	})
	lifecycle := lifecycleResponse.(*s3.GetBucketLifecycleConfigurationOutput)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() != 404 {
			return diag.FromErr(err)
		}
	}
	log.Printf("[DEBUG] S3 Bucket: %s, lifecycle: %v", d.Id(), lifecycle)
	if len(lifecycle.Rules) > 0 {
		rules := make([]map[string]interface{}, 0, len(lifecycle.Rules))

		for _, lifecycleRule := range lifecycle.Rules {
			rule := make(map[string]interface{})

			// ID
			if lifecycleRule.ID != nil && *lifecycleRule.ID != "" {
				rule["id"] = *lifecycleRule.ID
			}
			filter := lifecycleRule.Filter
			if filter != nil {
				if filter.And != nil {
					// Prefix
					if filter.And.Prefix != nil && *filter.And.Prefix != "" {
						rule["prefix"] = *filter.And.Prefix
					}
				} else {
					// Prefix
					if filter.Prefix != nil && *filter.Prefix != "" {
						rule["prefix"] = *filter.Prefix
					}
				}
			} else {
				if lifecycleRule.Prefix != nil { // nolint:staticcheck
					rule["prefix"] = *lifecycleRule.Prefix
				}
			}

			// Enabled
			if lifecycleRule.Status != nil {
				if *lifecycleRule.Status == s3.ExpirationStatusEnabled {
					rule["enabled"] = true
				} else {
					rule["enabled"] = false
				}
			}

			// AbortIncompleteMultipartUploadDays
			if lifecycleRule.AbortIncompleteMultipartUpload != nil {
				if lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
					rule["abort_incomplete_multipart_upload_days"] = int(*lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation)
				}
			}

			// expiration
			if lifecycleRule.Expiration != nil {
				e := make(map[string]interface{})
				if lifecycleRule.Expiration.Date != nil {
					e["date"] = lifecycleRule.Expiration.Date.Format("2006-01-02")
				}
				if lifecycleRule.Expiration.Days != nil {
					e["days"] = int(*lifecycleRule.Expiration.Days)
				}
				if lifecycleRule.Expiration.ExpiredObjectDeleteMarker != nil {
					e["expired_object_delete_marker"] = *lifecycleRule.Expiration.ExpiredObjectDeleteMarker
				}
				rule["expiration"] = schema.NewSet(ExpirationHash, []interface{}{e})
			}
			// noncurrent_version_expiration
			if lifecycleRule.NoncurrentVersionExpiration != nil {
				e := make(map[string]interface{})
				if lifecycleRule.NoncurrentVersionExpiration.NoncurrentDays != nil {
					e["days"] = int(*lifecycleRule.NoncurrentVersionExpiration.NoncurrentDays)
				}
				rule["noncurrent_version_expiration"] = schema.NewSet(ExpirationHash, []interface{}{e})
			}
			// transition
			if len(lifecycleRule.Transitions) > 0 {
				transitions := make([]interface{}, 0, len(lifecycleRule.Transitions))
				for _, v := range lifecycleRule.Transitions {
					t := make(map[string]interface{})
					if v.Date != nil {
						t["date"] = v.Date.Format("2006-01-02")
					}
					if v.Days != nil {
						t["days"] = int(*v.Days)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["transition"] = schema.NewSet(transitionHash, transitions)
			}
			// noncurrent_version_transition
			if len(lifecycleRule.NoncurrentVersionTransitions) > 0 {
				transitions := make([]interface{}, 0, len(lifecycleRule.NoncurrentVersionTransitions))
				for _, v := range lifecycleRule.NoncurrentVersionTransitions {
					t := make(map[string]interface{})
					if v.NoncurrentDays != nil {
						t["days"] = int(*v.NoncurrentDays)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["noncurrent_version_transition"] = schema.NewSet(transitionHash, transitions)
			}

			rules = append(rules, rule)
		}

		if err := d.Set("lifecycle_rule", rules); err != nil {
			return diag.FromErr(err)
		}
	}

	// Add the region as an attribute

	locationResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return client.GetBucketLocation(
			&s3.GetBucketLocationInput{
				Bucket: aws.String(d.Id()),
			},
		)
	})
	location := locationResponse.(*s3.GetBucketLocationOutput)
	if err != nil {
		return diag.FromErr(err)
	}
	var region string
	if location.LocationConstraint != nil {
		region = *location.LocationConstraint
	}
	if err := d.Set("region", region); err != nil {
		return diag.FromErr(err)
	}

	// Add the hosted zone ID for this bucket's region as an attribute
	// UNUSED?
	/*
		hostedZoneID := HostedZoneIDForRegion(region)
		if err := d.Set("hosted_zone_id", hostedZoneID); err != nil {
			return err
		}
	*/

	// Add website_endpoint as an attribute
	websiteEndpoint, err := websiteEndpoint(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	if websiteEndpoint != nil {
		if err := d.Set("website_endpoint", websiteEndpoint.Endpoint); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("website_domain", websiteEndpoint.Domain); err != nil {
			return diag.FromErr(err)
		}
	}

	tagSet, err := getTagSetS3(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", tagsToMapS3(tagSet)); err != nil {
		return diag.FromErr(err)
	}

	// UNUSED?
	// d.Set("arn", fmt.Sprintf("arn:%s:s3:::%s", meta.(*Config).partition, d.Id()))

	return nil
}

func resourceS3BucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	s3conn, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud S3 client: %s", err)
	}

	log.Printf("[DEBUG] S3 Delete Bucket: %s", d.Id())
	_, err = s3conn.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		ec2err, ok := err.(awserr.Error)
		if ok && ec2err.Code() == "BucketNotEmpty" {
			if d.Get("force_destroy").(bool) {
				// bucket may have things delete them
				log.Printf("[DEBUG] S3 Bucket attempting to forceDestroy %+v", err)

				bucket := d.Get("bucket").(string)
				resp, err := s3conn.ListObjectVersions(
					&s3.ListObjectVersionsInput{
						Bucket: aws.String(bucket),
					},
				)

				if err != nil {
					return fmterr.Errorf("error S3 Bucket list Object Versions err: %s", err)
				}

				objectsToDelete := make([]*s3.ObjectIdentifier, 0)

				if len(resp.DeleteMarkers) != 0 {
					for _, v := range resp.DeleteMarkers {
						objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
							Key:       v.Key,
							VersionId: v.VersionId,
						})
					}
				}

				if len(resp.Versions) != 0 {
					for _, v := range resp.Versions {
						objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
							Key:       v.Key,
							VersionId: v.VersionId,
						})
					}
				}

				params := &s3.DeleteObjectsInput{
					Bucket: aws.String(bucket),
					Delete: &s3.Delete{
						Objects: objectsToDelete,
					},
				}

				_, err = s3conn.DeleteObjects(params)

				if err != nil {
					return fmterr.Errorf("error S3 Bucket force_destroy error deleting: %s", err)
				}

				// this line recourses until all objects are deleted or an error is returned
				return resourceS3BucketDelete(ctx, d, meta)
			}
		}
		return fmterr.Errorf("error deleting S3 Bucket: %s %q", err, d.Get("bucket").(string))
	}
	return nil
}

func resourceS3BucketPolicyUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	if policy != "" {
		log.Printf("[DEBUG] S3 bucket: %s, put policy: %s", bucket, policy)

		params := &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucket),
			Policy: aws.String(policy),
		}

		err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
			if _, err := s3conn.PutBucketPolicy(params); err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == "MalformedPolicy" || awsErr.Code() == "NoSuchBucket" {
						return resource.RetryableError(awsErr)
					}
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error putting S3 policy: %s", err)
		}
	} else {
		log.Printf("[DEBUG] S3 bucket: %s, delete policy: %s", bucket, policy)
		_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
			return s3conn.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{
				Bucket: aws.String(bucket),
			})
		})

		if err != nil {
			return fmt.Errorf("error deleting S3 policy: %s", err)
		}
	}

	return nil
}

func resourceS3BucketCorsUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawCors := d.Get("cors_rule").([]interface{})

	if len(rawCors) == 0 {
		// Delete CORS
		log.Printf("[DEBUG] S3 bucket: %s, delete CORS", bucket)

		_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
			return s3conn.DeleteBucketCors(&s3.DeleteBucketCorsInput{
				Bucket: aws.String(bucket),
			})
		})
		if err != nil {
			return fmt.Errorf("error deleting S3 CORS: %s", err)
		}
	} else {
		// Put CORS
		rules := make([]*s3.CORSRule, 0, len(rawCors))
		for _, cors := range rawCors {
			corsMap := cors.(map[string]interface{})
			r := &s3.CORSRule{}
			for k, v := range corsMap {
				log.Printf("[DEBUG] S3 bucket: %s, put CORS: %#v, %#v", bucket, k, v)
				if k == "max_age_seconds" {
					r.MaxAgeSeconds = aws.Int64(int64(v.(int)))
				} else {
					vMap := make([]*string, len(v.([]interface{})))
					for i, vv := range v.([]interface{}) {
						str := vv.(string)
						vMap[i] = aws.String(str)
					}
					switch k {
					case "allowed_headers":
						r.AllowedHeaders = vMap
					case "allowed_methods":
						r.AllowedMethods = vMap
					case "allowed_origins":
						r.AllowedOrigins = vMap
					case "expose_headers":
						r.ExposeHeaders = vMap
					}
				}
			}
			rules = append(rules, r)
		}
		corsInput := &s3.PutBucketCorsInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: rules,
			},
		}
		log.Printf("[DEBUG] S3 bucket: %s, put CORS: %#v", bucket, corsInput)

		_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
			return s3conn.PutBucketCors(corsInput)
		})
		if err != nil {
			return fmt.Errorf("error putting S3 CORS: %s", err)
		}
	}

	return nil
}

func resourceS3BucketWebsiteUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	ws := d.Get("website").([]interface{})
	switch len(ws) {
	case 0:
		return resourceS3BucketWebsiteDelete(ctx, s3conn, d)
	case 1:
		var w map[string]interface{}
		if ws[0] != nil {
			w = ws[0].(map[string]interface{})
		} else {
			w = make(map[string]interface{})
		}
		return resourceS3BucketWebsitePut(ctx, s3conn, d, w)
	default:
		return fmt.Errorf("cannot specify more than one website")
	}
}

func resourceS3BucketWebsitePut(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData, website map[string]interface{}) error {
	bucket := d.Get("bucket").(string)

	var indexDocument, errorDocument, redirectAllRequestsTo, routingRules string
	if v, ok := website["index_document"]; ok {
		indexDocument = v.(string)
	}
	if v, ok := website["error_document"]; ok {
		errorDocument = v.(string)
	}
	if v, ok := website["redirect_all_requests_to"]; ok {
		redirectAllRequestsTo = v.(string)
	}
	if v, ok := website["routing_rules"]; ok {
		routingRules = v.(string)
	}

	if indexDocument == "" && redirectAllRequestsTo == "" {
		return fmt.Errorf("must specify either index_document or redirect_all_requests_to")
	}

	websiteConfiguration := &s3.WebsiteConfiguration{}

	if indexDocument != "" {
		websiteConfiguration.IndexDocument = &s3.IndexDocument{Suffix: aws.String(indexDocument)}
	}

	if errorDocument != "" {
		websiteConfiguration.ErrorDocument = &s3.ErrorDocument{Key: aws.String(errorDocument)}
	}

	if redirectAllRequestsTo != "" {
		redirect, err := url.Parse(redirectAllRequestsTo)
		if err == nil && redirect.Scheme != "" {
			var redirectHostBuf bytes.Buffer
			redirectHostBuf.WriteString(redirect.Host)
			if redirect.Path != "" {
				redirectHostBuf.WriteString(redirect.Path)
			}
			websiteConfiguration.RedirectAllRequestsTo = &s3.RedirectAllRequestsTo{HostName: aws.String(redirectHostBuf.String()), Protocol: aws.String(redirect.Scheme)}
		} else {
			websiteConfiguration.RedirectAllRequestsTo = &s3.RedirectAllRequestsTo{HostName: aws.String(redirectAllRequestsTo)}
		}
	}

	if routingRules != "" {
		var unmarshaledRules []*s3.RoutingRule
		if err := json.Unmarshal([]byte(routingRules), &unmarshaledRules); err != nil {
			return err
		}
		websiteConfiguration.RoutingRules = unmarshaledRules
	}

	putInput := &s3.PutBucketWebsiteInput{
		Bucket:               aws.String(bucket),
		WebsiteConfiguration: websiteConfiguration,
	}

	log.Printf("[DEBUG] S3 put bucket website: %#v", putInput)

	_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return s3conn.PutBucketWebsite(putInput)
	})
	if err != nil {
		return fmt.Errorf("error putting S3 website: %s", err)
	}

	return nil
}

func resourceS3BucketWebsiteDelete(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	deleteInput := &s3.DeleteBucketWebsiteInput{Bucket: aws.String(bucket)}

	log.Printf("[DEBUG] S3 delete bucket website: %#v", deleteInput)

	_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return s3conn.DeleteBucketWebsite(deleteInput)
	})
	if err != nil {
		return fmt.Errorf("error deleting S3 website: %s", err)
	}

	mErr := multierror.Append(
		d.Set("website_endpoint", ""),
		d.Set("website_domain", ""),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}

func websiteEndpoint(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) (*Website, error) {
	// If the bucket doesn't have a website configuration, return an empty
	// endpoint
	if _, ok := d.GetOk("website"); !ok {
		return nil, nil
	}

	bucket := d.Get("bucket").(string)

	// Lookup the region for this bucket

	locationResponse, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return s3conn.GetBucketLocation(
			&s3.GetBucketLocationInput{
				Bucket: aws.String(bucket),
			},
		)
	})
	location := locationResponse.(*s3.GetBucketLocationOutput)
	if err != nil {
		return nil, err
	}
	var region string
	if location.LocationConstraint != nil {
		region = *location.LocationConstraint
	}

	return WebsiteEndpoint(bucket, region), nil
}

func WebsiteEndpoint(bucket string, region string) *Website {
	domain := WebsiteDomainUrl(region)
	return &Website{Endpoint: fmt.Sprintf("%s.%s", bucket, domain), Domain: domain}
}

func WebsiteDomainUrl(region string) string {
	return fmt.Sprintf("s3-website.%s.amazonaws.com", region)
}

func resourceS3BucketAclUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	acl := d.Get("acl").(string)
	bucket := d.Get("bucket").(string)

	i := &s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	log.Printf("[DEBUG] S3 put bucket ACL: %#v", i)

	_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return s3conn.PutBucketAcl(i)
	})
	if err != nil {
		return fmt.Errorf("error putting S3 ACL: %s", err)
	}

	return nil
}

func resourceS3BucketVersioningUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	v := d.Get("versioning").([]interface{})
	bucket := d.Get("bucket").(string)
	vc := &s3.VersioningConfiguration{}

	if len(v) > 0 {
		c := v[0].(map[string]interface{})

		if c["enabled"].(bool) {
			vc.Status = aws.String(s3.BucketVersioningStatusEnabled)
		} else {
			vc.Status = aws.String(s3.BucketVersioningStatusSuspended)
		}

		if c["mfa_delete"].(bool) {
			vc.MFADelete = aws.String(s3.MFADeleteEnabled)
		} else {
			vc.MFADelete = aws.String(s3.MFADeleteDisabled)
		}
	} else {
		vc.Status = aws.String(s3.BucketVersioningStatusSuspended)
	}

	i := &s3.PutBucketVersioningInput{
		Bucket:                  aws.String(bucket),
		VersioningConfiguration: vc,
	}
	log.Printf("[DEBUG] S3 put bucket versioning: %#v", i)

	_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return s3conn.PutBucketVersioning(i)
	})
	if err != nil {
		return fmt.Errorf("error putting S3 versioning: %s", err)
	}

	return nil
}

func resourceS3BucketLoggingUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	logging := d.Get("logging").(*schema.Set).List()
	bucket := d.Get("bucket").(string)
	loggingStatus := &s3.BucketLoggingStatus{}

	if len(logging) > 0 {
		c := logging[0].(map[string]interface{})

		loggingEnabled := &s3.LoggingEnabled{}
		if val, ok := c["target_bucket"]; ok {
			loggingEnabled.TargetBucket = aws.String(val.(string))
		}
		if val, ok := c["target_prefix"]; ok {
			loggingEnabled.TargetPrefix = aws.String(val.(string))
		}

		loggingStatus.LoggingEnabled = loggingEnabled
	}

	i := &s3.PutBucketLoggingInput{
		Bucket:              aws.String(bucket),
		BucketLoggingStatus: loggingStatus,
	}
	log.Printf("[DEBUG] S3 put bucket logging: %#v", i)

	_, err := retryOnAwsCode(ctx, "NoSuchBucket", func() (interface{}, error) {
		return s3conn.PutBucketLogging(i)
	})
	if err != nil {
		return fmt.Errorf("error putting S3 logging: %s", err)
	}

	return nil
}

func resourceAwsS3BucketLifecycleUpdate(ctx context.Context, s3conn *s3.S3, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)

	lifecycleRules := d.Get("lifecycle_rule").([]interface{})

	// fmt.Printf("lifecycleRules=%+v.\n", lifecycleRules)
	if len(lifecycleRules) == 0 {
		i := &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		}

		err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
			if _, err := s3conn.DeleteBucketLifecycle(i); err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error removing S3 lifecycle: %s", err)
		}
		return nil
	}

	rules := make([]*s3.LifecycleRule, 0, len(lifecycleRules))

	for i, lifecycleRule := range lifecycleRules {
		r := lifecycleRule.(map[string]interface{})

		rule := &s3.LifecycleRule{}

		// OTC only supports deprecated location for this.
		rule.SetPrefix(r["prefix"].(string))

		// ID
		if val, ok := r["id"].(string); ok && val != "" {
			rule.ID = aws.String(val)
		} else {
			rule.ID = aws.String(resource.PrefixedUniqueId("tf-s3-lifecycle-"))
		}

		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rule.Status = aws.String(s3.ExpirationStatusEnabled)
		} else {
			rule.Status = aws.String(s3.ExpirationStatusDisabled)
		}

		// AbortIncompleteMultipartUpload
		if val, ok := r["abort_incomplete_multipart_upload_days"].(int); ok && val > 0 {
			rule.AbortIncompleteMultipartUpload = &s3.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: aws.Int64(int64(val)),
			}
		}

		// Expiration
		expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expiration) > 0 {
			e := expiration[0].(map[string]interface{})
			i := &s3.LifecycleExpiration{}

			if val, ok := e["date"].(string); ok && val != "" {
				t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
				if err != nil {
					return fmt.Errorf("error Parsing Swift S3 Bucket Lifecycle Expiration Date: %s", err)
				}
				i.Date = aws.Time(t)
			} else if val, ok := e["days"].(int); ok && val > 0 {
				i.Days = aws.Int64(int64(val))
			} else if val, ok := e["expired_object_delete_marker"].(bool); ok {
				i.ExpiredObjectDeleteMarker = aws.Bool(val)
			}
			rule.Expiration = i
		}

		// NonCurrentVersionExpiration
		ncExpiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_expiration", i)).(*schema.Set).List()
		if len(ncExpiration) > 0 {
			e := ncExpiration[0].(map[string]interface{})

			if val, ok := e["days"].(int); ok && val > 0 {
				rule.NoncurrentVersionExpiration = &s3.NoncurrentVersionExpiration{
					NoncurrentDays: aws.Int64(int64(val)),
				}
			}
		}

		rules = append(rules, rule)
	}

	// fmt.Printf("Rules=%+v.\n", rules)
	i := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: rules,
		},
	}
	// fmt.Printf("PutBucketLifecycleConfigurationInput=%+v.\n", i)

	err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		if _, err := s3conn.PutBucketLifecycleConfiguration(i); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error putting S3 lifecycle: %s", err)
	}

	return nil
}

func normalizeRoutingRules(w []*s3.RoutingRule) (string, error) {
	withNulls, err := json.Marshal(w)
	if err != nil {
		return "", err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(withNulls, &rules); err != nil {
		return "", err
	}

	var cleanRules []map[string]interface{}
	for _, rule := range rules {
		cleanRules = append(cleanRules, RemoveNil(rule))
	}

	withoutNulls, err := json.Marshal(cleanRules)
	if err != nil {
		return "", err
	}

	return string(withoutNulls), nil
}

func transitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return hashcode.String(buf.String())
}

type Website struct {
	Endpoint, Domain string
}

func resourceS3BucketImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1)
	results[0] = d

	config := meta.(*cfg.Config)
	conn, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
	}
	pol, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchBucketPolicy" {
			// Bucket without policy
			return results, nil
		}
		return nil, fmt.Errorf("error importing AWS S3 bucket policy: %w", err)
	}

	policy := ResourceS3BucketPolicy()
	pData := policy.Data(nil)
	pData.SetId(d.Id())
	pData.SetType("opentelekomcloud_s3_bucket_policy")
	_ = pData.Set("bucket", d.Id())
	_ = pData.Set("policy", pol)
	results = append(results, pData)

	return results, nil
}
