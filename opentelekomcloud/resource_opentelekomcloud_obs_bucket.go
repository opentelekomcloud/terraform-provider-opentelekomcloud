package opentelekomcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
)

func resourceObsBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketCreate,
		Read:   resourceObsBucketRead,
		Update: resourceObsBucketUpdate,
		Delete: resourceObsBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STANDARD",
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD", "WARM", "COLD",
				}, true),
			},

			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					"private", "public-read", "public-read-write", "log-delivery-write",
				}, true),
			},

			"versioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
							Default:  "logs/",
						},
					},
				},
			},

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"transition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"WARM", "COLD",
										}, true),
									},
								},
							},
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
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
										Required: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"WARM", "COLD",
										}, true),
									},
								},
							},
						},
					},
				},
			},

			"website": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
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
							ValidateFunc: validateJsonString,
							StateFunc: func(v interface{}) string {
								jsonString, _ := normalizeJsonString(v)
								return jsonString
							},
						},
					},
				},
			},

			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_headers": {
							Type:     schema.TypeList,
							Optional: true,
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
							Default:  100,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"bucket_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObsBucketCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	class := d.Get("storage_class").(string)
	opts := &obs.CreateBucketInput{
		Bucket:       bucket,
		ACL:          obs.AclType(acl),
		StorageClass: obs.StorageClassType(class),
	}
	opts.Location = d.Get("region").(string)
	log.Printf("[DEBUG] OBS bucket create opts: %#v", opts)

	_, err = client.CreateBucket(opts)
	if err != nil {
		return getObsError("Error creating bucket", bucket, err)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceObsBucketUpdate(d, meta)
}

func resourceObsBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	log.Printf("[DEBUG] Update OBS bucket %s", d.Id())
	if d.HasChange("acl") && !d.IsNewResource() {
		if err := resourceObsBucketAclUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("storage_class") && !d.IsNewResource() {
		if err := resourceObsBucketClassUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		if err := resourceObsBucketTagsUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("versioning") {
		if err := resourceObsBucketVersioningUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("logging") {
		if err := resourceObsBucketLoggingUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceObsBucketLifecycleUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("website") {
		if err := resourceObsBucketWebsiteUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("cors_rule") {
		if err := resourceObsBucketCorsUpdate(obsClient, d); err != nil {
			return err
		}
	}

	return resourceObsBucketRead(d, meta)
}

func resourceObsBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	log.Printf("[DEBUG] Read OBS bucket: %s", d.Id())
	_, err = obsClient.HeadBucket(d.Id())
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			log.Printf("[WARN] OBS bucket(%s) not found", d.Id())
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("error reading OBS bucket %s: %s", d.Id(), err)
		}
	}

	mErr := &multierror.Error{}

	// for import case
	if _, ok := d.GetOk("bucket"); !ok {
		mErr = multierror.Append(mErr, d.Set("bucket", d.Id()))
	}

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("bucket_domain_name", bucketDomainName(d.Get("bucket").(string), region)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting OBS bucket fields: %s", err)
	}

	// Read storage class
	if err := setObsBucketStorageClass(obsClient, d); err != nil {
		return err
	}

	// Read the versioning
	if err := setObsBucketVersioning(obsClient, d); err != nil {
		return err
	}
	// Read the logging configuration
	if err := setObsBucketLogging(obsClient, d); err != nil {
		return err
	}

	// Read the Lifecycle configuration
	if err := setObsBucketLifecycleConfiguration(obsClient, d); err != nil {
		return err
	}

	// Read the website configuration
	if err := setObsBucketWebsiteConfiguration(obsClient, d); err != nil {
		return err
	}

	// Read the CORS rules
	if err := setObsBucketCorsRules(obsClient, d); err != nil {
		return err
	}

	// Read the tags
	if err := setObsBucketTags(obsClient, d); err != nil {
		return err
	}

	return nil
}

func resourceObsBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OBS client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] deleting OBS Bucket: %s", bucket)
	_, err = obsClient.DeleteBucket(bucket)
	if err != nil {
		obsError, ok := err.(obs.ObsError)
		if ok && obsError.Code == "BucketNotEmpty" {
			log.Printf("[WARN] OBS bucket: %s is not empty", bucket)
			if d.Get("force_destroy").(bool) {
				err = deleteAllBucketObjects(obsClient, bucket)
				if err == nil {
					log.Printf("[WARN] all objects of %s have been deleted, and try again", bucket)
					return resourceObsBucketDelete(d, meta)
				}
			}
			return err
		}
		return fmt.Errorf("error deleting OBS bucket: %s %s", bucket, err)
	}
	return nil
}

func resourceObsBucketTagsUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	tagMap := d.Get("tags").(map[string]interface{})
	var tagList []obs.Tag
	for k, v := range tagMap {
		tag := obs.Tag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	req := &obs.SetBucketTaggingInput{}
	req.Bucket = bucket
	req.Tags = tagList
	log.Printf("[DEBUG] set tags of OBS bucket %s: %#v", bucket, req)

	_, err := obsClient.SetBucketTagging(req)
	if err != nil {
		return getObsError("error updating tags of OBS bucket", bucket, err)
	}
	return nil
}

func resourceObsBucketAclUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)

	i := &obs.SetBucketAclInput{
		Bucket: bucket,
		ACL:    obs.AclType(acl),
	}
	log.Printf("[DEBUG] set ACL of OBS bucket %s: %#v", bucket, i)

	_, err := obsClient.SetBucketAcl(i)
	if err != nil {
		return getObsError("Error updating acl of OBS bucket", bucket, err)
	}

	// acl policy can not be retrieved by obsClient.GetBucketAcl method
	err = d.Set("acl", acl)
	return err
}

func resourceObsBucketClassUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	class := d.Get("storage_class").(string)

	input := &obs.SetBucketStoragePolicyInput{}
	input.Bucket = bucket
	input.StorageClass = obs.StorageClassType(class)
	log.Printf("[DEBUG] set storage class of OBS bucket %s: %#v", bucket, input)

	_, err := obsClient.SetBucketStoragePolicy(input)
	if err != nil {
		return getObsError("Error updating storage class of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketVersioningUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	version := d.Get("versioning").(bool)

	input := &obs.SetBucketVersioningInput{}
	input.Bucket = bucket
	if version {
		input.Status = obs.VersioningStatusEnabled
	} else {
		input.Status = obs.VersioningStatusSuspended
	}
	log.Printf("[DEBUG] set versioning of OBS bucket %s: %#v", bucket, input)

	_, err := obsClient.SetBucketVersioning(input)
	if err != nil {
		return getObsError("Error setting versining status of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketLoggingUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawLogging := d.Get("logging").(*schema.Set).List()
	loggingStatus := &obs.SetBucketLoggingConfigurationInput{}
	loggingStatus.Bucket = bucket

	if len(rawLogging) > 0 {
		c := rawLogging[0].(map[string]interface{})
		if val := c["target_bucket"].(string); val != "" {
			loggingStatus.TargetBucket = val
		}

		if val := c["target_prefix"].(string); val != "" {
			loggingStatus.TargetPrefix = val
		}
	}
	log.Printf("[DEBUG] set logging of OBS bucket %s: %#v", bucket, loggingStatus)

	_, err := obsClient.SetBucketLoggingConfiguration(loggingStatus)
	if err != nil {
		return getObsError("Error setting logging configuration of OBS bucket", bucket, err)
	}

	return nil
}

func mapToRule(src map[string]interface{}) (rule obs.LifecycleRule) {
	// rule ID
	rule.ID = src["name"].(string)

	// Enabled
	if val, ok := src["enabled"].(bool); ok && val {
		rule.Status = obs.RuleStatusEnabled
	} else {
		rule.Status = obs.RuleStatusDisabled
	}

	// Prefix
	rule.Prefix = src["prefix"].(string)

	// Expiration
	expiration := src["expiration"].(*schema.Set).List()
	if len(expiration) > 0 {
		raw := expiration[0].(map[string]interface{})
		exp := &rule.Expiration

		if val, ok := raw["days"].(int); ok && val > 0 {
			exp.Days = val
		}
	}

	// Transition
	transitions := src["transition"].([]interface{})
	list := make([]obs.Transition, len(transitions))
	for j, tran := range transitions {
		raw := tran.(map[string]interface{})

		if val, ok := raw["days"].(int); ok && val > 0 {
			list[j].Days = val
		}
		if val, ok := raw["storage_class"].(string); ok {
			list[j].StorageClass = obs.StorageClassType(val)
		}
	}
	rule.Transitions = list

	// NoncurrentVersionExpiration
	ncExpiration := src["noncurrent_version_expiration"].(*schema.Set).List()
	if len(ncExpiration) > 0 {
		raw := ncExpiration[0].(map[string]interface{})
		ncExp := &rule.NoncurrentVersionExpiration

		if val, ok := raw["days"].(int); ok && val > 0 {
			ncExp.NoncurrentDays = val
		}
	}

	// NoncurrentVersionTransition
	ncTransitions := src["noncurrent_version_transition"].([]interface{})
	ncList := make([]obs.NoncurrentVersionTransition, len(ncTransitions))
	for j, ncTran := range ncTransitions {
		raw := ncTran.(map[string]interface{})

		if val, ok := raw["days"].(int); ok && val > 0 {
			ncList[j].NoncurrentDays = val
		}
		if val, ok := raw["storage_class"].(string); ok {
			ncList[j].StorageClass = obs.StorageClassType(val)
		}
	}
	rule.NoncurrentVersionTransitions = ncList
	return
}

func resourceObsBucketLifecycleUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	lifecycleRules := d.Get("lifecycle_rule").([]interface{})

	if len(lifecycleRules) == 0 {
		log.Printf("[DEBUG] remove all lifecycle rules of bucket %s", bucket)
		_, err := obsClient.DeleteBucketLifecycleConfiguration(bucket)
		if err != nil {
			return getObsError("Error deleting lifecycle rules of OBS bucket", bucket, err)
		}
		return nil
	}

	rules := make([]obs.LifecycleRule, len(lifecycleRules))
	for i, lifecycleRule := range lifecycleRules {
		ruleMap := lifecycleRule.(map[string]interface{})
		rules[i] = mapToRule(ruleMap)
	}

	opts := &obs.SetBucketLifecycleConfigurationInput{}
	opts.Bucket = bucket
	opts.LifecycleRules = rules
	log.Printf("[DEBUG] set lifecycle configurations of OBS bucket %s: %#v", bucket, opts)

	_, err := obsClient.SetBucketLifecycleConfiguration(opts)
	if err != nil {
		return getObsError("error setting lifecycle rules of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketWebsiteUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	websiteList := d.Get("website").([]interface{})

	switch len(websiteList) {
	case 0:
		return resourceObsBucketWebsiteDelete(obsClient, d)
	case 1:
		var website map[string]interface{}
		if websiteList[0] != nil {
			website = websiteList[0].(map[string]interface{})
		} else {
			website = make(map[string]interface{})
		}
		return resourceObsBucketWebsitePut(obsClient, d, website)
	default:
		return fmt.Errorf("cannot specify more than one website")
	}
}

// asStringSlice returns `value` and `ok`
func asStringSlice(value interface{}) ([]string, bool) {
	sliceVal, ok := value.([]interface{})
	if !ok {
		return nil, false
	}
	result := make([]string, len(sliceVal))
	for i, val := range sliceVal {
		result[i], ok = val.(string)
		if !ok {
			return nil, false
		}
	}
	return result, true
}

func resourceObsBucketCorsUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawCors := d.Get("cors_rule").([]interface{})

	if len(rawCors) == 0 {
		// Delete CORS
		log.Printf("[DEBUG] delete CORS rules of OBS bucket: %s", bucket)
		_, err := obsClient.DeleteBucketCors(bucket)
		if err != nil {
			return getObsError("error deleting CORS rules of OBS bucket", bucket, err)
		}
		return nil
	}

	// set CORS
	rules := make([]obs.CorsRule, len(rawCors))
	for i, cors := range rawCors {
		corsMap := cors.(map[string]interface{})
		rule := obs.CorsRule{}
		for k, v := range corsMap {
			if k == "max_age_seconds" {
				rule.MaxAgeSeconds = v.(int)
				continue
			}
			value, ok := asStringSlice(v)
			if !ok {
				continue
			}
			switch k {
			case "allowed_headers":
				rule.AllowedHeader = value
			case "allowed_methods":
				rule.AllowedMethod = value
			case "allowed_origins":
				rule.AllowedOrigin = value
			case "expose_headers":
				rule.ExposeHeader = value
			}
		}
		log.Printf("[DEBUG] set CORS of OBS bucket %s: %#v", bucket, rule)
		rules[i] = rule
	}

	corsInput := &obs.SetBucketCorsInput{}
	corsInput.Bucket = bucket
	corsInput.CorsRules = rules
	log.Printf("[DEBUG] OBS bucket: %s, put CORS: %#v", bucket, corsInput)

	_, err := obsClient.SetBucketCors(corsInput)
	if err != nil {
		return getObsError("Error setting CORS rules of OBS bucket", bucket, err)
	}
	return nil
}

func resourceObsBucketWebsitePut(obsClient *obs.ObsClient, d *schema.ResourceData, website map[string]interface{}) error {
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

	websiteConfiguration := &obs.SetBucketWebsiteConfigurationInput{}
	websiteConfiguration.Bucket = bucket

	if indexDocument != "" {
		websiteConfiguration.IndexDocument = obs.IndexDocument{
			Suffix: indexDocument,
		}
	}

	if errorDocument != "" {
		websiteConfiguration.ErrorDocument = obs.ErrorDocument{
			Key: errorDocument,
		}
	}

	if redirectAllRequestsTo != "" {
		redirect, err := url.Parse(redirectAllRequestsTo)
		if err == nil && redirect.Scheme != "" {
			var redirectHostBuf bytes.Buffer
			redirectHostBuf.WriteString(redirect.Host)
			if redirect.Path != "" {
				redirectHostBuf.WriteString(redirect.Path)
			}
			websiteConfiguration.RedirectAllRequestsTo = obs.RedirectAllRequestsTo{
				HostName: redirectHostBuf.String(),
				Protocol: obs.ProtocolType(redirect.Scheme),
			}
		} else {
			websiteConfiguration.RedirectAllRequestsTo = obs.RedirectAllRequestsTo{
				HostName: redirectAllRequestsTo,
			}
		}
	}

	if routingRules != "" {
		var unmarshalledRules []obs.RoutingRule
		if err := json.Unmarshal([]byte(routingRules), &unmarshalledRules); err != nil {
			return err
		}
		websiteConfiguration.RoutingRules = unmarshalledRules
	}

	log.Printf("[DEBUG] set website configuration of OBS bucket %s: %#v", bucket, websiteConfiguration)
	_, err := obsClient.SetBucketWebsiteConfiguration(websiteConfiguration)
	if err != nil {
		return getObsError("Error updating website configuration of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketWebsiteDelete(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] delete website configuration of OBS bucket %s", bucket)
	_, err := obsClient.DeleteBucketWebsiteConfiguration(bucket)
	if err != nil {
		return getObsError("error deleting website configuration of OBS bucket", bucket, err)
	}

	return nil
}

func setObsBucketStorageClass(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketStoragePolicy(bucket)
	if err != nil {
		log.Printf("[WARN] Error getting storage class of OBS bucket %s: %s", bucket, err)
		return nil
	} else {
		class := output.StorageClass
		err = d.Set("storage_class", normalizeStorageClass(class))
	}

	return err
}

func setObsBucketVersioning(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketVersioning(bucket)
	if err != nil {
		return getObsError("error getting versioning status of OBS bucket", bucket, err)
	}

	enabled := output.Status == obs.VersioningStatusEnabled
	err = d.Set("versioning", enabled)

	return err
}

func setObsBucketLogging(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketLoggingConfiguration(bucket)
	if err != nil {
		return getObsError("Error getting logging configuration of OBS bucket", bucket, err)
	}

	var lcList []map[string]interface{}
	logging := make(map[string]interface{})

	if output.TargetBucket != "" {
		logging["target_bucket"] = output.TargetBucket
		if output.TargetPrefix != "" {
			logging["target_prefix"] = output.TargetPrefix
		}
		lcList = append(lcList, logging)
	}
	log.Printf("[DEBUG] saving logging configuration of OBS bucket: %s: %#v", bucket, lcList)

	if err := d.Set("logging", lcList); err != nil {
		return fmt.Errorf("error saving logging configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func ruleToMap(src obs.LifecycleRule) map[string]interface{} {
	rule := make(map[string]interface{})
	rule["name"] = src.ID

	// Enabled
	rule["enabled"] = src.Status == obs.RuleStatusEnabled

	if src.Prefix != "" {
		rule["prefix"] = src.Prefix
	}

	// expiration
	if days := src.Expiration.Days; days > 0 {
		expiration := make(map[string]interface{})
		expiration["days"] = days
		rule["expiration"] = schema.NewSet(expirationHash, []interface{}{expiration})
	}
	// transition
	if len(src.Transitions) > 0 {
		transitions := make([]interface{}, len(src.Transitions))
		for i, v := range src.Transitions {
			transition := make(map[string]interface{})
			transition["days"] = v.Days
			transition["storage_class"] = normalizeStorageClass(string(v.StorageClass))
			transitions[i] = transition
		}
		rule["transition"] = transitions
	}

	// noncurrent_version_expiration
	if days := src.NoncurrentVersionExpiration.NoncurrentDays; days > 0 {
		expiration := make(map[string]interface{})
		expiration["days"] = days
		rule["noncurrent_version_expiration"] = schema.NewSet(expirationHash, []interface{}{expiration})
	}

	// noncurrent_version_transition
	if len(src.NoncurrentVersionTransitions) > 0 {
		transitions := make([]interface{}, len(src.NoncurrentVersionTransitions))
		for i, v := range src.NoncurrentVersionTransitions {
			transition := make(map[string]interface{})
			transition["days"] = v.NoncurrentDays
			transition["storage_class"] = normalizeStorageClass(string(v.StorageClass))
			transitions[i] = transition
		}
		rule["noncurrent_version_transition"] = transitions
	}
	return rule
}

func setObsBucketLifecycleConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketLifecycleConfiguration(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchLifecycleConfiguration" {
				err = d.Set("lifecycle_rule", nil)
				return err
			}
			return fmt.Errorf("Error getting lifecycle configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	rawRules := output.LifecycleRules
	log.Printf("[DEBUG] getting lifecycle configuration of OBS bucket: %s, lifecycle: %#v", bucket, rawRules)

	rules := make([]map[string]interface{}, len(rawRules))
	for i, lifecycleRule := range rawRules {
		rules[i] = ruleToMap(lifecycleRule)
	}

	log.Printf("[DEBUG] saving lifecycle configuration of OBS bucket: %s, lifecycle: %#v", bucket, rules)
	if err := d.Set("lifecycle_rule", rules); err != nil {
		return fmt.Errorf("error saving lifecycle configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func handleWebsite(src *obs.GetBucketWebsiteConfigurationOutput) (map[string]interface{}, error) {
	website := make(map[string]interface{})
	website["index_document"] = src.IndexDocument.Suffix
	website["error_document"] = src.ErrorDocument.Key

	// redirect_all_requests_to
	v := src.RedirectAllRequestsTo
	if string(v.Protocol) == "" {
		website["redirect_all_requests_to"] = v.HostName
	} else {
		var host string
		var path string
		parsedHostName, err := url.Parse(v.HostName)
		if err == nil {
			host = parsedHostName.Host
			path = parsedHostName.Path
		} else {
			host = v.HostName
			path = ""
		}

		website["redirect_all_requests_to"] = (&url.URL{
			Host:   host,
			Path:   path,
			Scheme: string(v.Protocol),
		}).String()
	}

	// routing_rules
	rawRules := src.RoutingRules
	if len(rawRules) > 0 {
		rr, err := normalizeWebsiteRoutingRules(rawRules)
		if err != nil {
			return nil, fmt.Errorf("error while marshaling website routing rules: %s", err)
		}
		website["routing_rules"] = rr
	}
	return website, nil
}

func setObsBucketWebsiteConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	rawWebsite, err := obsClient.GetBucketWebsiteConfiguration(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchWebsiteConfiguration" {
				err = d.Set("website", nil)
				return err
			} else {
				return fmt.Errorf("Error getting website configuration of OBS bucket %s: %s,\n Reason: %s",
					bucket, obsError.Code, obsError.Message)
			}
		} else {
			return err
		}
	}

	log.Printf("[DEBUG] getting website configuration of OBS bucket: %s, output: %#v", bucket, rawWebsite.BucketWebsiteConfiguration)

	website, err := handleWebsite(rawWebsite)
	if err != nil {
		return err
	}

	websites := []map[string]interface{}{website}
	log.Printf("[DEBUG] saving website configuration of OBS bucket: %s, website: %#v", bucket, websites)
	if err := d.Set("website", websites); err != nil {
		return fmt.Errorf("error saving website configuration of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func setObsBucketCorsRules(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketCors(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchCORSConfiguration" {
				err = d.Set("cors_rule", nil)
				return err
			} else {
				return fmt.Errorf("Error getting CORS configuration of OBS bucket %s: %s,\n Reason: %s",
					bucket, obsError.Code, obsError.Message)
			}
		} else {
			return err
		}
	}

	corsRules := output.CorsRules
	log.Printf("[DEBUG] getting CORS rules of OBS bucket: %s, CORS: %#v", bucket, corsRules)

	rules := make([]map[string]interface{}, len(corsRules))
	for i, ruleObject := range corsRules {
		rule := make(map[string]interface{})
		rule["allowed_origins"] = ruleObject.AllowedOrigin
		rule["allowed_methods"] = ruleObject.AllowedMethod
		rule["max_age_seconds"] = ruleObject.MaxAgeSeconds
		if ruleObject.AllowedHeader != nil {
			rule["allowed_headers"] = ruleObject.AllowedHeader
		}
		if ruleObject.ExposeHeader != nil {
			rule["expose_headers"] = ruleObject.ExposeHeader
		}
		rules[i] = rule
	}

	log.Printf("[DEBUG] saving CORS rules of OBS bucket: %s, CORS: %#v", bucket, rules)
	if err := d.Set("cors_rule", rules); err != nil {
		return fmt.Errorf("error saving CORS rules of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketTags(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketTagging(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchTagSet" {
				err = d.Set("tags", nil)
				return err
			} else {
				return fmt.Errorf("error getting tags of OBS bucket %s: %s,\n Reason: %s",
					bucket, obsError.Code, obsError.Message)
			}
		} else {
			return err
		}
	}

	tagmap := make(map[string]string)
	for _, tag := range output.Tags {
		tagmap[tag.Key] = tag.Value
	}
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("error saving tags of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func deleteAllBucketObjects(obsClient *obs.ObsClient, bucket string) error {
	listOpts := &obs.ListObjectsInput{
		Bucket: bucket,
	}
	// list all objects
	resp, err := obsClient.ListObjects(listOpts)
	if err != nil {
		return getObsError("error listing objects of OBS bucket", bucket, err)
	}

	objects := make([]obs.ObjectToDelete, len(resp.Contents))
	for i, content := range resp.Contents {
		objects[i].Key = content.Key
	}

	deleteOpts := &obs.DeleteObjectsInput{
		Bucket:  bucket,
		Objects: objects,
	}
	log.Printf("[DEBUG] objects of %s will be deleted: %v", bucket, objects)
	output, err := obsClient.DeleteObjects(deleteOpts)
	if err != nil {
		return getObsError("error deleting all objects of OBS bucket", bucket, err)
	} else {
		if len(output.Errors) > 0 {
			return fmt.Errorf("error some objects are still exist in %s: %#v", bucket, output.Errors)
		}
	}
	return nil
}

func getObsError(action string, bucket string, err error) error {
	if obsError, ok := err.(obs.ObsError); ok {
		return fmt.Errorf("%s %s: %s,\n Reason: %s", action, bucket, obsError.Code, obsError.Message)
	}
	return err
}

// normalize format of storage class
func normalizeStorageClass(class string) string {
	var ret = class

	if class == "STANDARD_IA" {
		ret = "WARM"
	} else if class == "GLACIER" {
		ret = "COLD"
	}
	return ret
}

func normalizeWebsiteRoutingRules(w []obs.RoutingRule) (string, error) {
	// transform []obs.RoutingRule to []WebsiteRoutingRule
	wRules := make([]WebsiteRoutingRule, len(w))
	for i, rawRule := range w {
		rule := WebsiteRoutingRule{
			Condition: Condition{
				KeyPrefixEquals:             rawRule.Condition.KeyPrefixEquals,
				HttpErrorCodeReturnedEquals: rawRule.Condition.HttpErrorCodeReturnedEquals,
			},
			Redirect: Redirect{
				Protocol:             string(rawRule.Redirect.Protocol),
				HostName:             rawRule.Redirect.HostName,
				HttpRedirectCode:     rawRule.Redirect.HttpRedirectCode,
				ReplaceKeyWith:       rawRule.Redirect.ReplaceKeyWith,
				ReplaceKeyPrefixWith: rawRule.Redirect.ReplaceKeyPrefixWith,
			},
		}
		wRules[i] = rule
	}

	// normalize
	withNulls, err := json.Marshal(wRules)
	if err != nil {
		return "", err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(withNulls, &rules); err != nil {
		return "", err
	}

	var cleanRules []map[string]interface{}
	for _, rule := range rules {
		cleanRules = append(cleanRules, removeNil(rule))
	}

	withoutNulls, err := json.Marshal(cleanRules)
	if err != nil {
		return "", err
	}

	return string(withoutNulls), nil
}

type Condition struct {
	KeyPrefixEquals             string `json:"KeyPrefixEquals,omitempty"`
	HttpErrorCodeReturnedEquals string `json:"HttpErrorCodeReturnedEquals,omitempty"`
}

type Redirect struct {
	Protocol             string `json:"Protocol,omitempty"`
	HostName             string `json:"HostName,omitempty"`
	ReplaceKeyPrefixWith string `json:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKeyWith       string `json:"ReplaceKeyWith,omitempty"`
	HttpRedirectCode     string `json:"HttpRedirectCode,omitempty"`
}
type WebsiteRoutingRule struct {
	Condition Condition `json:"Condition,omitempty"`
	Redirect  Redirect  `json:"Redirect"`
}
