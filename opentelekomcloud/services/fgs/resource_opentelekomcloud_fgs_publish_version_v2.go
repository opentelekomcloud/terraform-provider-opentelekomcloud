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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/alias"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourcePublishVersionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublishVersionV2Create,
		ReadContext:   resourcePublishVersionV2Read,
		DeleteContext: resourcePublishVersionV2Delete,

		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the URN of the function.",
			},
			"digest": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the digest of the function code.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the version of the function.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the description of the version.",
			},
			"func_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the name of the function.",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the last modified time of the version.",
			},
		},
	}
}

func parsePublishVersionResourceId(resourceId string) (funcUrn, version string) {
	parts := strings.SplitN(resourceId, "/", 2)
	if len(parts) < 2 {
		log.Printf("[ERROR] invalid ID format for FGS publish version resource: %s", resourceId)
		return
	}
	funcUrn = parts[0]
	version = parts[1]
	return
}

func resourcePublishVersionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := alias.PublishOpts{
		FuncUrn:     d.Get("function_urn").(string),
		Digest:      d.Get("digest").(string),
		Version:     d.Get("version").(string),
		Description: d.Get("description").(string),
	}

	createResp, err := alias.PublishVersion(fgsClient, createOpts)
	if err != nil {
		return diag.Errorf("error publishing OpenTelekomCloud FunctionGraph version: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", createOpts.FuncUrn, createResp.Version))

	clientCtx := common.CtxWithClient(ctx, fgsClient, fgsClientV2)
	return resourcePublishVersionV2Read(clientCtx, d, meta)
}

func resourcePublishVersionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	fgsClient, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	funcUrn, version := parsePublishVersionResourceId(d.Id())

	listOpts := alias.ListVersionOpts{
		FuncUrn: funcUrn,
	}

	listResp, err := alias.ListVersion(fgsClient, listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud FunctionGraph version")
	}

	var found bool
	for _, v := range listResp.Functions {
		if v.Version == version {
			mErr := multierror.Append(
				d.Set("function_urn", funcUrn),
				d.Set("version", v.Version),
				d.Set("digest", v.Digest),
				d.Set("description", v.Description),
				d.Set("func_name", v.FuncName),
				d.Set("last_modified", v.LastModified),
			)
			if err := mErr.ErrorOrNil(); err != nil {
				return diag.Errorf("error setting resource fields of FGS publish version (%s): %s", d.Id(), err)
			}
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return diag.Errorf("error finding FGS publish version (%s): %s", d.Id(), err)
	}

	return nil
}

func resourcePublishVersionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}
