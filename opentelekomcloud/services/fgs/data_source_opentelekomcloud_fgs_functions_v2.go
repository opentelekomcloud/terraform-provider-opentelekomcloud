package fgs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/function"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceFunctionsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"package_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"functions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_filename": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted_user_data": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_agency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"initializer_handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initializer_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_stream_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"functiongraph_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func filterFunctions(config *cfg.Config, d *schema.ResourceData, functions []function.FuncGraph) []function.FuncGraph {
	result := functions

	if nameRaw, ok := d.GetOk("name"); ok && len(result) > 0 {
		name := nameRaw.(string)
		filtered := make([]function.FuncGraph, 0, len(result))
		for _, fn := range result {
			if fn.FuncName == name {
				filtered = append(filtered, fn)
			}
		}
		result = filtered
	}

	if urnRaw, ok := d.GetOk("urn"); ok && len(result) > 0 {
		urn := urnRaw.(string)
		filtered := make([]function.FuncGraph, 0, len(result))
		for _, fn := range result {
			if fn.FuncURN == urn {
				filtered = append(filtered, fn)
			}
		}
		result = filtered
	}

	if runtimeRaw, ok := d.GetOk("runtime"); ok && len(result) > 0 {
		runtime := runtimeRaw.(string)
		filtered := make([]function.FuncGraph, 0, len(result))
		for _, fn := range result {
			if fn.Runtime == runtime {
				filtered = append(filtered, fn)
			}
		}
		result = filtered
	}

	if epsID := config.GetEnterpriseProjectID(d); epsID != "" && len(result) > 0 {
		filtered := make([]function.FuncGraph, 0, len(result))
		for _, fn := range result {
			if fn.EnterpriseProjectId == epsID {
				filtered = append(filtered, fn)
			}
		}
		result = filtered
	}

	return result
}

func flattenFunctions(fns []function.FuncGraph) []map[string]interface{} {
	if len(fns) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(fns))
	for _, fn := range fns {
		result = append(result, map[string]interface{}{
			"name":                  fn.FuncName,
			"urn":                   fn.FuncURN,
			"package":               fn.Package,
			"runtime":               fn.Runtime,
			"timeout":               fn.Timeout,
			"handler":               fn.Handler,
			"memory_size":           fn.MemorySize,
			"code_type":             fn.CodeType,
			"code_url":              fn.CodeURL,
			"code_filename":         fn.CodeFilename,
			"user_data":             fn.UserData,
			"encrypted_user_data":   fn.EncryptedUserData,
			"version":               fn.Version,
			"app_agency":            fn.AppXrole,
			"description":           fn.Description,
			"vpc_id":                fn.FuncVpc.VpcID,
			"network_id":            fn.FuncVpc.SubnetID,
			"max_instance_num":      fn.StrategyConfig.ConcurrentNum,
			"initializer_handler":   fn.InitHandler,
			"initializer_timeout":   fn.InitTimeout,
			"enterprise_project_id": fn.EnterpriseProjectId,
			"log_group_id":          fn.LogGroupID,
			"log_stream_id":         fn.LogStreamID,
			"functiongraph_version": fn.Type,
		})
	}

	return result
}

func dataSourceFunctionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	functions, err := function.List(client, function.ListOpts{
		MaxItems:    "100",
		PackageName: d.Get("package_name").(string),
	})
	if err != nil {
		return diag.Errorf("error querying functions: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("functions", flattenFunctions(filterFunctions(config, d, functions.Functions))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
