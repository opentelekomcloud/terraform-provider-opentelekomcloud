package fgs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/dependency_version"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDependencyVersionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDependencyVersionV2Create,
		ReadContext:   resourceDependencyVersionV2Read,
		DeleteContext: resourceDependencyVersionV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"link",
					"file",
				},
			},
			"file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"link",
					"file",
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"download_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dependency_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDependencyVersionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := dependency_version.CreateOpts{
		Name:        d.Get("name").(string),
		Runtime:     d.Get("runtime").(string),
		Description: d.Get("description").(string),
	}

	if d.Get("link").(string) != "" {
		createOpts.DependType = "obs"
		createOpts.DependLink = d.Get("link").(string)
	} else {
		createOpts.DependType = "zip"
		createOpts.DependFile = d.Get("file").(string)
	}

	createResp, err := dependency_version.Create(fgsClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud dependency package version: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%v", createResp.DepId, createResp.Version))

	clientCtx := common.CtxWithClient(ctx, fgsClient, fgsClientV2)
	return resourceDependencyVersionV2Read(clientCtx, d, meta)
}

func ParseDependVersionResourceId(resourceId string) (dependId, versionInfo string) {
	parts := strings.Split(resourceId, "/")
	if len(parts) < 2 {
		log.Printf("[ERROR] invalid ID format for dependency package version resource, it must contain two parts: "+
			"dependency package information and version information, e.g. '<dependency name>/<version number>'. "+
			"but the ID that you provided does not meet this requirement '%s'", resourceId)
		return
	}
	dependId = parts[0]
	versionInfo = parts[1]
	return
}

func resourceDependencyVersionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	dependId, dependVersion := ParseDependVersionResourceId(d.Id())

	getResp, err := dependency_version.Get(fgsClient, dependId, dependVersion)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud FunctionGraph dependency package version")
	}

	mErr := multierror.Append(
		d.Set("runtime", getResp.Runtime),
		d.Set("name", getResp.Name),
		d.Set("description", getResp.Description),
		d.Set("download_link", getResp.Link),
		d.Set("file_name", getResp.FileName),
		d.Set("etag", getResp.Etag),
		d.Set("size", getResp.Size),
		d.Set("owner", getResp.Owner),
		d.Set("version", getResp.Version),
		d.Set("dependency_id", getResp.DepId),
		d.Set("version_id", getResp.Id),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of custom dependency package version (%s): %s",
			d.Id(), err)
	}

	return nil
}

func resourceDependencyVersionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	dependId, dependVersion := ParseDependVersionResourceId(d.Id())

	err = dependency_version.Delete(fgsClient, dependId, dependVersion)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting custom dependency package version")
	}
	return nil
}
