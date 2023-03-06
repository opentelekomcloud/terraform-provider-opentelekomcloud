package ims

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/image/v2/images"
	ims "github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceImagesImageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesImageV2Create,
		ReadContext:   resourceImagesImageV2Read,
		UpdateContext: resourceImagesImageV2Update,
		DeleteContext: resourceImagesImageV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"container_format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceImagesImageV2ValidateContainerFormat,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_format": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     resourceImagesImageV2ValidateDiskFormat,
				DiffSuppressFunc: common.SuppressDiffAll, // UNDONE: OTC appears broken here, so hack work-around...
			},

			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_cache_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  fmt.Sprintf("%s/.terraform/image_cache", os.Getenv("HOME")),
			},

			"image_source_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"local_file_path"},
			},

			"local_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"image_source_url"},
			},

			"min_disk_gb": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntAtLeast(1),
				Default:          0,
				DiffSuppressFunc: common.SuppressMinDisk,
			},

			"min_ram_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Default:      0,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"update_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: resourceImagesImageV2ValidateVisibility,
				Default:      "private",
			},
		},
	}
}

func resourceImagesImageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	imageClient, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	createOpts := images.CreateOpts{
		Name:            d.Get("name").(string),
		ContainerFormat: d.Get("container_format").(string),
		DiskFormat:      d.Get("disk_format").(string),
		MinDisk:         d.Get("min_disk_gb").(int),
		MinRam:          d.Get("min_ram_mb").(int),
		Protected:       d.Get("protected").(bool),
		Visibility:      d.Get("visibility").(string),
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = resourceImagesImageV2BuildTags(tags)
	}

	d.Partial(true)

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newImg, err := images.Create(imageClient, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Image: %s", err)
	}

	d.SetId(newImg.Id)

	// downloading/getting image file props
	imgFilePath, err := resourceImagesImageV2File(d)
	if err != nil {
		return fmterr.Errorf("error opening file for Image: %s", err)
	}
	fileSize, fileChecksum, err := resourceImagesImageV2FileProps(imgFilePath)
	if err != nil {
		return fmterr.Errorf("error getting file props: %s", err)
	}

	// upload
	imgFile, err := os.Open(imgFilePath)
	if err != nil {
		return fmterr.Errorf("error opening file %q: %s", imgFilePath, err)
	}
	defer func() {
		_ = imgFile.Close()
	}()
	log.Printf("[WARN] Uploading image %s (%d bytes). This can be pretty long.", d.Id(), fileSize)

	res := images.Upload(imageClient, d.Id(), imgFile)
	if res != nil {
		return fmterr.Errorf("error while uploading file %q: %s", imgFilePath, res)
	}

	// wait for active
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"queued", "saving"},
		Target:     []string{"active"},
		Refresh:    resourceImagesImageV2RefreshFunc(imageClient, d.Id(), fileSize, fileChecksum),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for Image: %s", err)
	}

	d.Partial(false)

	return resourceImagesImageV2Read(ctx, d, meta)
}

func resourceImagesImageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	imageClient, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	img, err := images.Get(imageClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "image")
	}

	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	size, _ := strconv.Atoi(img.ImageSize)
	mErr := multierror.Append(
		d.Set("owner", img.Owner),
		d.Set("status", img.Status),
		d.Set("file", img.File),
		d.Set("schema", img.Schema),
		d.Set("checksum", img.Checksum),
		d.Set("size_bytes", size),
		d.Set("created_at", img.CreatedAt.Format(time.RFC3339)),
		d.Set("update_at", img.UpdatedAt.Format(time.RFC3339)),
		d.Set("container_format", img.ContainerFormat),
		d.Set("disk_format", img.DiskFormat),
		d.Set("min_disk_gb", img.MinDisk),
		d.Set("min_ram_mb", img.MinRam),
		d.Set("file", img.File),
		d.Set("name", img.Name),
		d.Set("protected", img.Protected),
		d.Set("tags", img.Tags),
		d.Set("visibility", img.Visibility),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceImagesImageV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	imageClient, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	updateOpts := make([]ims.UpdateImageOpts, 0)

	if d.HasChange("visibility") {
		updateOpts = append(updateOpts, ims.UpdateImageOpts{
			Op:    "replace",
			Path:  "/visibility",
			Value: d.Get("visibility").(string),
		})
	}

	if d.HasChange("name") {
		updateOpts = append(updateOpts, ims.UpdateImageOpts{
			Op:    "replace",
			Path:  "/name",
			Value: d.Get("name").(string),
		})
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(*schema.Set).List()
		updateOpts = append(updateOpts, ims.UpdateImageOpts{
			Op:    "replace",
			Path:  "/tags",
			Value: resourceImagesImageV2BuildTags(tags),
		})
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = images.Update(imageClient, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating image: %s", err)
	}

	return resourceImagesImageV2Read(ctx, d, meta)
}

func resourceImagesImageV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	imageClient, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	log.Printf("[DEBUG] Deleting Image %s", d.Id())
	if err := images.Delete(imageClient, images.DeleteImageOpts{
		ImageId:      d.Id(),
		DeleteBackup: true,
	}); err != nil {
		return fmterr.Errorf("error deleting Image: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceImagesImageV2ValidateVisibility(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validVisibilities := []string{
		"public",
		"private",
		"shared",
		"community",
	}

	for _, v := range validVisibilities {
		if value == v {
			return
		}
	}

	err := fmt.Errorf("%s must be one of %s", k, validVisibilities)
	errors = append(errors, err)
	return
}

var DiskFormats = [9]string{"ami", "ari", "aki", "vhd", "vmdk", "raw", "qcow2", "vdi", "iso"}

func resourceImagesImageV2ValidateDiskFormat(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range DiskFormats {
		if value == DiskFormats[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, DiskFormats))
	return
}

var ContainerFormats = [9]string{"ami", "ari", "aki", "bare", "ovf"}

func resourceImagesImageV2ValidateContainerFormat(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range ContainerFormats {
		if value == ContainerFormats[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, ContainerFormats))
	return
}

func fileMD5Checksum(f *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func resourceImagesImageV2FileProps(filename string) (int64, string, error) {
	var filesize int64
	var fileChecksum string

	file, err := os.Open(filename)
	if err != nil {
		return -1, "", fmt.Errorf("error opening file for Image: %s", err)
	}
	defer func() { _ = file.Close() }()

	fstat, err := file.Stat()
	if err != nil {
		return -1, "", fmt.Errorf("error reading image file %q: %s", file.Name(), err)
	}

	filesize = fstat.Size()
	fileChecksum, err = fileMD5Checksum(file)

	if err != nil {
		return -1, "", fmt.Errorf("error computing image file %q checksum: %s", file.Name(), err)
	}

	return filesize, fileChecksum, nil
}

func resourceImagesImageV2File(d *schema.ResourceData) (string, error) {
	if filename := d.Get("local_file_path").(string); filename != "" {
		return filename, nil
	} else if furl := d.Get("image_source_url").(string); furl != "" {
		dir := d.Get("image_cache_path").(string)
		_ = os.MkdirAll(dir, 0700)
		filename := filepath.Join(dir, fmt.Sprintf("%x.img", md5.Sum([]byte(furl))))

		if _, err := os.Stat(filename); err != nil {
			if !os.IsNotExist(err) {
				return "", fmt.Errorf("error while trying to access file %q: %s", filename, err)
			}
			log.Printf("[DEBUG] File doens't exists %s. will download from %s", filename, furl)
			file, err := os.Create(filename)
			if err != nil {
				return "", fmt.Errorf("error creating file %q: %s", filename, err)
			}
			defer func() { _ = file.Close() }()
			resp, err := http.Get(furl)
			if err != nil {
				return "", fmt.Errorf("error downloading image from %q", furl)
			}
			defer func() { _ = resp.Body.Close() }()
			if _, err = io.Copy(file, resp.Body); err != nil {
				return "", fmt.Errorf("error downloading image %q to file %q: %s", furl, filename, err)
			}
			return filename, nil
		} else {
			log.Printf("[DEBUG] File exists %s", filename)
			return filename, nil
		}
	} else {
		return "", fmt.Errorf("error in config. no file specified")
	}
}

func resourceImagesImageV2RefreshFunc(client *golangsdk.ServiceClient, id string, _ int64, _ string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		img, err := images.Get(client, id)
		if err != nil {
			return nil, "", err
		}
		log.Printf("[DEBUG] OpenTelekomCloud image status is: %s", img.Status)

		// Huawei provider doesn't have this set initially.
		/*
			if img.Checksum != checksum || int64(img.SizeBytes) != fileSize {
				return img, fmt.Sprintf("%s", img.Status), fmt.Errorf("error wrong size %v or checksum %q", img.SizeBytes, img.Checksum)
			}
		*/
		return img, img.Status, nil
	}
}

func resourceImagesImageV2BuildTags(v []interface{}) []string {
	var tags []string
	for _, tag := range v {
		tags = append(tags, tag.(string))
	}

	return tags
}
