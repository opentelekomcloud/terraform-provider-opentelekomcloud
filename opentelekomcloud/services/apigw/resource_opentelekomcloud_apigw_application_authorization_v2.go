package apigw

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	appauth "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app_auth"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIAppAuthV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppAuthV2Create,
		ReadContext:   resourceAppAuthV2Read,
		UpdateContext: resourceAppAuthV2Update,
		DeleteContext: resourceAppAuthV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppAuthV2ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppAuthV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var (
		appId  = d.Get("application_id").(string)
		envId  = d.Get("env_id").(string)
		apiIds = common.ExpandToStringListBySet(d.Get("api_ids").(*schema.Set))
	)
	err = createAppAuthForApis(ctx, client, d, apiIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", envId, appId))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAppAuthV2Read(clientCtx, d, meta)
}

func flattenAuthorizedApis(apiInfos []appauth.ApiAuth) []string {
	result := make([]string, len(apiInfos))
	for i, val := range apiInfos {
		result[i] = val.ApiID
	}
	return result
}

func resourceAppAuthV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	appId := d.Get("application_id").(string)
	opts := appauth.ListBoundOpts{
		GatewayID: gatewayId,
		AppID:     appId,
	}
	resp, err := appauth.ListAPIBound(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying OpenTelekomCloud APIGW authorized APIs from application (%s) under dedicated instance (%s)",
			appId, gatewayId))
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Application Authorization")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("api_ids", flattenAuthorizedApis(resp)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW authorization fields for specified application (%s): %s", appId, err)
	}
	return nil
}

func resourceAppAuthV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	oldRaw, newRaw := d.GetChange("api_ids")
	addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	if rmSet.Len() > 0 {
		apiIds := common.ExpandToStringListBySet(rmSet)
		err := deleteAppAuthFromApis(ctx, client, d, apiIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		apiIds := common.ExpandToStringListBySet(addSet)
		err = createAppAuthForApis(ctx, client, d, apiIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAppAuthV2Read(clientCtx, d, meta)
}

func resourceAppAuthV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	apiIds := d.Get("api_ids").(*schema.Set)
	err = deleteAppAuthFromApis(ctx, client, d, common.ExpandToStringListBySet(apiIds))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func createAppAuthForApis(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, apiIds []string) error {
	gatewayId := d.Get("gateway_id").(string)
	appId := d.Get("application_id").(string)
	envId := d.Get("env_id").(string)
	opts := appauth.CreateAuthOpts{
		GatewayID: gatewayId,
		AppIDs:    []string{appId},
		EnvID:     envId,
		ApiIDs:    apiIds,
	}

	_, err := appauth.Create(client, opts)
	if err != nil {
		return fmt.Errorf("error authorizing OpenTelekomCloud APIGW APIs to application (%s) under dedicated instance (%s): %s", appId, gatewayId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    authApisStateRefreshFunc(client, gatewayId, appId, envId, apiIds),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for OpenTelekomCloud APIGW API authorize operations completed: %s", err)
	}

	return nil
}

func deleteAppAuthFromApis(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, apiIds []string) error {
	gatewayId := d.Get("gateway_id").(string)
	appId := d.Get("application_id").(string)
	envId := d.Get("env_id").(string)
	opts := appauth.ListBoundOpts{
		GatewayID: gatewayId,
		AppID:     appId,
		EnvID:     envId,
	}
	Err404 := fmt.Sprintf("[DEBUG] All APIs have been unauthorized form application (%s) under dedicated instance (%s)", appId, gatewayId)

	resp, err := appauth.ListAPIBound(client, opts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Println(Err404)
			return nil
		}
		return fmt.Errorf("error querying authorized APIs for application (%s) under dedicated instance (%s)", appId, gatewayId)
	}
	if len(resp) < 1 {
		log.Println(Err404)
		return nil
	}

	for _, val := range resp {
		if !common.StrSliceContains(apiIds, val.ApiID) {
			continue
		}

		authId := val.ID
		err = appauth.Delete(client, gatewayId, authId)
		if err != nil {
			return fmt.Errorf("error unauthorizing OpenTelekomCloud APIGW APIs form application (%s) under dedicated instance (%s): %s", appId, gatewayId, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    deauthApisStateRefreshFunc(client, gatewayId, appId, envId, apiIds),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for OpenTelekomCloud APIGW API unauthorize operations completed: %s", err)
	}
	return nil
}

func authApisStateRefreshFunc(client *golangsdk.ServiceClient, gatewayId, appId, envId string, apiIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := appauth.ListUnboundOpts{
			GatewayID: gatewayId,
			AppID:     appId,
			EnvID:     envId,
		}
		resp, err := appauth.ListAPIUnBound(client, opts)
		if err != nil {
			return resp, "", err
		}

		for _, val := range resp {
			if common.StrSliceContains(apiIds, val.ID) {
				return resp, "PENDING", nil
			}
		}
		return resp, "COMPLETED", nil
	}
}

func deauthApisStateRefreshFunc(client *golangsdk.ServiceClient, gatewayId, appId, envId string, apiIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := appauth.ListBoundOpts{
			GatewayID: gatewayId,
			AppID:     appId,
			EnvID:     envId,
		}
		resp, err := appauth.ListAPIBound(client, opts)
		if err != nil {
			return resp, "", err
		}

		for _, val := range resp {
			if common.StrSliceContains(apiIds, val.ApiID) {
				return resp, "PENDING", nil
			}
		}
		return resp, "COMPLETED", nil
	}
}

func resourceAppAuthV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<id>' (the format of resource ID is "+
			"'<env_id>/<application_id>'), but got '%s'", importedId)
	}

	d.SetId(fmt.Sprintf("%s/%s", parts[1], parts[2]))
	mErr := multierror.Append(nil,
		d.Set("gateway_id", parts[0]),
		d.Set("env_id", parts[1]),
		d.Set("application_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d},
			fmt.Errorf("error saving OpenTelekomCloud APIGW application authorization resource fields during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
