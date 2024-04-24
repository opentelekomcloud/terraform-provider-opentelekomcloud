package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/env_vars"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceAPIGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"environment_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"registration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createEnvironmentVariables(client *golangsdk.ServiceClient, instanceId, groupId string,
	environmentSet *schema.Set) error {
	for _, envVar := range environmentSet.List() {
		envMap := envVar.(map[string]interface{})
		envId := envMap["environment_id"].(string)
		for _, v := range envMap["variable"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			opt := env_vars.CreateOpts{
				VariableValue: variable["value"].(string),
				VariableName:  variable["name"].(string),
				GroupID:       groupId,
				EnvID:         envId,
				GatewayID:     instanceId,
			}
			if _, err := env_vars.Create(client, opt); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeEnvironmentVariables(client *golangsdk.ServiceClient, instanceId string,
	environmentSet *schema.Set) error {
	for _, env := range environmentSet.List() {
		envMap := env.(map[string]interface{})
		for _, v := range envMap["variable"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			err := env_vars.Delete(client, instanceId, variable["id"].(string))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)

	opt := group.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		GatewayID:   instanceId,
	}

	resp, err := group.Create(client, opt)
	if err != nil {
		return diag.Errorf("error creating dedicated group: %s", err)
	}
	d.SetId(resp.ID)

	if environmentSet, ok := d.GetOk("environment"); ok {
		err = createEnvironmentVariables(client, instanceId, d.Id(), environmentSet.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceGroupRead(ctx, d, meta)
}

func queryEnvironmentVariables(client *golangsdk.ServiceClient, instanceId, groupId string) ([]env_vars.EnvVarsResp, error) {
	opt := env_vars.ListOpts{
		GroupID:   groupId,
		GatewayID: instanceId,
	}
	pages, err := env_vars.List(client, opt)
	if err != nil {
		return nil, fmt.Errorf("error getting environment variable list from server: %s", err)
	}
	return pages, nil
}

func flattenEnvironmentVariables(variables []env_vars.EnvVarsResp) []map[string]interface{} {
	if len(variables) < 1 {
		return nil
	}

	environmentMap := make(map[string]interface{})
	for _, variable := range variables {
		varMap := map[string]interface{}{
			"name":  variable.VariableName,
			"value": variable.VariableValue,
			"id":    variable.ID,
		}
		if val, ok := environmentMap[variable.EnvID]; !ok {
			environmentMap[variable.EnvID] = []map[string]interface{}{
				varMap,
			}
		} else {
			environmentMap[variable.EnvID] = append(val.([]map[string]interface{}), varMap)
		}
	}

	result := make([]map[string]interface{}, 0, len(environmentMap))
	for k, v := range environmentMap {
		envMap := map[string]interface{}{
			"variable":       v,
			"environment_id": k,
		}
		result = append(result, envMap)
	}

	return result
}

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Id()
	)

	resp, err := group.Get(client, instanceId, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("registration_time", resp.RegisterTime),
		d.Set("updated_at", resp.UpdateTime),
	)
	var variables []env_vars.EnvVarsResp
	if variables, err = queryEnvironmentVariables(client, instanceId, groupId); err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr, d.Set("environment", flattenEnvironmentVariables(variables)))

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving dedicated group fields： %s", mErr)
	}
	return nil
}

func updateEnvironmentVariables(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oldRaws, newRaws = d.GetChange("environment")
		addRaws          = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws       = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
		instanceId       = d.Get("instance_id").(string)
		groupId          = d.Id()
	)
	if err := removeEnvironmentVariables(client, instanceId, removeRaws); err != nil {
		return err
	}
	return createEnvironmentVariables(client, instanceId, groupId, addRaws)
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Id()
	)

	if d.HasChanges("name", "description") {
		opt := group.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			GroupID:     groupId,
			GatewayID:   instanceId,
		}
		_, err = group.Update(client, opt)
		if err != nil {
			return diag.Errorf("error updating dedicated group (%s): %s", groupId, err)
		}
	}

	if d.HasChange("environment") {
		if err := updateEnvironmentVariables(client, d); err != nil {
			return diag.Errorf("error updating environment variables: %s", err)
		}
	}
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = group.Delete(client, instanceId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting group from the instance (%s): %s", instanceId, err)
	}

	return nil
}

func resourceGroupResourceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
