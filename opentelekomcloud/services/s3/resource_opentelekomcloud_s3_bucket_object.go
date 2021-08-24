package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	// "github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceS3BucketObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceS3BucketObjectPut,
		ReadContext:   resourceS3BucketObjectRead,
		UpdateContext: resourceS3BucketObjectPut,
		DeleteContext: resourceS3BucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"acl": {
				Type:         schema.TypeString,
				Default:      "private",
				Optional:     true,
				ValidateFunc: ValidateS3BucketObjectAclType,
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},

			"server_side_encryption": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateS3BucketObjectServerSideEncryption,
				Computed:     true,
			},

			"sse_kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"etag": {
				Type: schema.TypeString,
				// This will conflict with SSE-C and SSE-KMS encryption and multi-part upload
				// if/when it's actually implemented. The Etag then won't match raw-file MD5.
				// See http://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html
				Optional: true,
				Computed: true,
				// ConflictsWith: []string{"kms_key_id", "server_side_encryption"},
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"website_redirect": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceS3BucketObjectPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	s3conn, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
	}

	var body io.ReadSeeker

	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmterr.Errorf("error expanding homedir in source (%s): %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmterr.Errorf("error opening S3 bucket object source (%s): %s", source, err)
		}

		body = file
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else {
		return fmterr.Errorf("must specify \"source\" or \"content\" field")
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String(d.Get("acl").(string)),
		Body:   body,
	}

	if v, ok := d.GetOk("cache_control"); ok {
		putInput.CacheControl = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_type"); ok {
		putInput.ContentType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		putInput.ContentEncoding = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_language"); ok {
		putInput.ContentLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		putInput.ContentDisposition = aws.String(v.(string))
	}

	if v, ok := d.GetOk("server_side_encryption"); ok {
		putInput.ServerSideEncryption = aws.String(v.(string))
	}

	if v, ok := d.GetOk("website_redirect"); ok {
		putInput.WebsiteRedirectLocation = aws.String(v.(string))
	}

	if v, ok := d.GetOk("sse_kms_key_id"); ok {
		putInput.SSEKMSKeyId = aws.String(v.(string))
	}

	resp, err := s3conn.PutObject(putInput)
	if err != nil {
		return fmterr.Errorf("error putting object in S3 bucket (%s): %s", bucket, err)
	}

	d.SetId(key)
	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	mErr := multierror.Append(
		d.Set("etag", strings.Trim(*resp.ETag, `"`)),
		d.Set("version_id", resp.VersionId),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return resourceS3BucketObjectRead(ctx, d, meta)
}

func resourceS3BucketObjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	s3conn, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
	}

	// restricted := false //meta.(*AWSClient).IsGovCloud() || meta.(*AWSClient).IsChinaCloud()

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	resp, err := s3conn.HeadObject(
		&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		// If S3 returns a 404 Request Failure, mark the object as destroyed
		if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == 404 {
			d.SetId("")
			log.Printf("[WARN] Error Reading Object (%s), object not found (HTTP status 404)", key)
			return nil
		}
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Reading S3 Bucket Object meta: %s", resp)

	mErr := multierror.Append(
		d.Set("cache_control", resp.CacheControl),
		d.Set("content_disposition", resp.ContentDisposition),
		d.Set("content_encoding", resp.ContentEncoding),
		d.Set("content_language", resp.ContentLanguage),
		d.Set("content_type", resp.ContentType),
		d.Set("version_id", resp.VersionId),
		d.Set("server_side_encryption", resp.ServerSideEncryption),
		d.Set("website_redirect", resp.WebsiteRedirectLocation),
		d.Set("sse_kms_key_id", resp.SSEKMSKeyId),
		d.Set("etag", strings.Trim(*resp.ETag, `"`)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceS3BucketObjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	s3conn, err := config.S3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if _, ok := d.GetOk("version_id"); ok {
		// Bucket is versioned, we need to delete all versions
		vInput := s3.ListObjectVersionsInput{
			Bucket: aws.String(bucket),
			Prefix: aws.String(key),
		}
		out, err := s3conn.ListObjectVersions(&vInput)
		if err != nil {
			return fmterr.Errorf("failed listing S3 object versions: %s", err)
		}

		for _, v := range out.Versions {
			input := s3.DeleteObjectInput{
				Bucket:    aws.String(bucket),
				Key:       aws.String(key),
				VersionId: v.VersionId,
			}
			_, err := s3conn.DeleteObject(&input)
			if err != nil {
				return fmterr.Errorf("error deleting S3 object version of %s:\n %s:\n %s",
					key, v, err)
			}
		}
	} else {
		// Just delete the object
		input := s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
		_, err := s3conn.DeleteObject(&input)
		if err != nil {
			return fmterr.Errorf("error deleting S3 bucket object: %s  Bucket: %q Object: %q", err, bucket, key)
		}
	}

	return nil
}

func ValidateS3BucketObjectAclType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	cannedAcls := map[string]bool{
		s3.ObjectCannedACLPrivate:                true,
		s3.ObjectCannedACLPublicRead:             true,
		s3.ObjectCannedACLPublicReadWrite:        true,
		s3.ObjectCannedACLAuthenticatedRead:      true,
		s3.ObjectCannedACLAwsExecRead:            true,
		s3.ObjectCannedACLBucketOwnerRead:        true,
		s3.ObjectCannedACLBucketOwnerFullControl: true,
	}

	sentenceJoin := func(m map[string]bool) string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, fmt.Sprintf("%q", k))
		}
		sort.Strings(keys)

		length := len(keys)
		words := make([]string, length)
		copy(words, keys)

		words[length-1] = fmt.Sprintf("or %s", words[length-1])
		return strings.Join(words, ", ")
	}

	if _, ok := cannedAcls[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid canned ACL type %q. Valid types are either %s",
			k, value, sentenceJoin(cannedAcls)))
	}
	return ws, errors
}

func validateS3BucketObjectServerSideEncryption(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	serverSideEncryption := map[string]bool{
		s3.ServerSideEncryptionAes256: true,
		s3.ServerSideEncryptionAwsKms: true,
	}

	if _, ok := serverSideEncryption[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Server Side Encryption value %q. Valid values are %q and %q",
			k, value, s3.ServerSideEncryptionAes256, s3.ServerSideEncryptionAwsKms))
	}
	return
}
