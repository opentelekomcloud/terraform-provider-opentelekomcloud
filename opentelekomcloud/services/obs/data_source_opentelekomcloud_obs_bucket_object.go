package obs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/s3"
)

func DataSourceObsBucketObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObsBucketObjectRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cache_control": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_disposition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"website_redirect_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceObsBucketObjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	input := obs.GetObjectInput{
		GetObjectMetadataInput: obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    key,
		},
	}
	if v, ok := d.GetOk("version_id"); ok {
		input.VersionId = v.(string)
	}

	versionText := ""
	uniqueId := bucket + "/" + key
	if v, ok := d.GetOk("version_id"); ok {
		versionText = fmt.Sprintf(" of version %q", v.(string))
		uniqueId += "@" + v.(string)
	}

	log.Printf("[DEBUG] Reading OBS object: %v", input)
	out, err := client.GetObject(&input)
	if err != nil {
		return fmterr.Errorf("failed getting OBS object: %s Bucket: %q Object: %q", err, bucket, key)
	}
	if out.DeleteMarker {
		return fmterr.Errorf("requested OBS object %q%s has been deleted",
			bucket+key, versionText)
	}

	log.Printf("[DEBUG] Received OBS object: %v", out)

	d.SetId(uniqueId)

	mErr := multierror.Append(
		d.Set("cache_control", out.CacheControl),
		d.Set("content_disposition", out.ContentDisposition),
		d.Set("content_encoding", out.ContentEncoding),
		d.Set("content_language", out.ContentLanguage),
		d.Set("content_length", out.ContentLength),
		d.Set("content_type", out.ContentType),
		d.Set("etag", strings.Trim(out.ETag, `"`)),
		d.Set("expiration", out.Expiration),
		d.Set("expires", out.Expires),
		d.Set("last_modified", out.LastModified.Format(time.RFC1123)),
		d.Set("metadata", out.Metadata),
		d.Set("version_id", out.VersionId),
		d.Set("website_redirect_location", out.WebsiteRedirectLocation),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(mErr)
	}

	if s3.IsContentTypeAllowed(&out.ContentType) {
		input := &obs.GetObjectInput{
			GetObjectMetadataInput: obs.GetObjectMetadataInput{
				Bucket: bucket,
				Key:    key,
			},
		}
		input.VersionId = out.VersionId
		out, err := client.GetObject(input)
		if err != nil {
			return fmterr.Errorf("failed getting OBS object: %s", err)
		}

		buf := new(bytes.Buffer)
		bytesRead, err := buf.ReadFrom(out.Body)
		if err != nil {
			return fmterr.Errorf("failed reading content of OBS object (%s): %s",
				uniqueId, err)
		}
		log.Printf("[INFO] Saving %d bytes from OBS object %s", bytesRead, uniqueId)
		_ = d.Set("body", buf.String())
	} else {
		contentType := "<EMPTY>"
		if out.ContentType != "" {
			contentType = out.ContentType
		}

		log.Printf("[INFO] Ignoring body of OBS object %s with Content-Type %q",
			uniqueId, contentType)
	}

	return nil
}
