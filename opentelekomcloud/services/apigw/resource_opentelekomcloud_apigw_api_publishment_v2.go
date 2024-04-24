package apigw

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/api"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIApiPublishmentV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceApiPublishmentV2Create,
		ReadContext:   ResourceApiPublishmentV2Read,
		UpdateContext: ResourceApiPublishmentV2Update,
		DeleteContext: ResourceApiPublishmentV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApiPublishmentV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publish_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The publish ID of the API in current environment.",
			},
			"environment_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the current version was published.",
			},
			"history": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceApiPublishmentV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	apiId := d.Get("api_id").(string)
	envId := d.Get("environment_id").(string)

	if versionId, hasVer := d.GetOk("version_id"); hasVer {
		history, err := GetVersionHistory(client, gatewayId, envId, apiId)
		if err != nil {
			return diag.Errorf("error finding the publish versions of the OpenTelekomCloud APIGW API (%s) in the environment (%s): %s",
				apiId, envId, err)
		}
		if ver, ok := isPublished(history, versionId.(string)); !ok {
			return diag.Errorf("the version (%s) has not published", versionId.(string))
		} else if desc, ok := d.GetOk("description"); ok && desc != ver.Description {
			return diag.Errorf("the description is no correct, expected '%s', actual '%s', please check your description "+
				"input or API version", ver.Description, desc)
		}
		if _, ok := isLatestVersion(history, versionId.(string)); !ok {
			_, err := apis.SwitchVersion(client, apis.VersionApiOpts{
				GatewayID: gatewayId,
				ApiID:     apiId,
				VersionID: versionId.(string),
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		err = publishApiToSpecifiedEnv(client, gatewayId, envId, apiId, d.Get("description").(string))
		if err != nil {
			return diag.Errorf("error publishing OpenTelekomCloud APIGW API: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", gatewayId, envId, apiId))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return ResourceApiPublishmentV2Read(clientCtx, d, meta)
}

func ResourceApiPublishmentV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := GetVersionHistory(client, instanceId, envId, apiId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error getting the publish versions of the OpenTelekomCloud APIGW API (%s) in the environment (%s)", apiId, envId))
	}

	publishInfo, err := getCertainPublishInfo(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr := multierror.Append(nil,
		d.Set("environment_id", publishInfo.EnvID),
		d.Set("api_id", publishInfo.ApiID),
		d.Set("description", publishInfo.Description),
		d.Set("published_at", publishInfo.PublishTime),
		d.Set("environment_name", publishInfo.EnvName),
		d.Set("region", config.GetRegion(d)),
		setApiPublishHistory(d, resp),
	)

	if publishId, err := getPublishIdByEnvId(client, instanceId, publishInfo.ApiID, publishInfo.EnvID); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("publish_id", publishId))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW API publishment fields: %s", err)
	}
	return nil
}

func ResourceApiPublishmentV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	versionId := d.Get("version_id").(string)
	gatewayId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if versionId == "" {
		if err = publishApiToSpecifiedEnv(client, gatewayId, envId, apiId, d.Get("description").(string)); err != nil {
			return diag.Errorf("error publishing OpenTelekomCloud APIGW API: %s", err)
		}
	} else {
		if !d.HasChange("version_id") && d.HasChange("description") {
			return diag.Errorf("only for new OpenTelekomCloud APIGW API publishment, the description can be updated")
		}
		description := d.Get("description").(string)

		history, err := GetVersionHistory(client, gatewayId, envId, apiId)
		if err != nil {
			return diag.Errorf("error getting version histories of the OpenTelekomCloud APIGW API (%s): %s", apiId, err)
		}
		if ver, ok := isPublished(history, versionId); !ok {
			return diag.Errorf("this version (%s) has not published", versionId)
		} else if description != "" && ver.Description != description {
			return diag.Errorf("this description is not belongs to version (%s)", versionId)
		}

		if _, err := apis.SwitchVersion(client, apis.VersionApiOpts{
			GatewayID: gatewayId,
			ApiID:     apiId,
			VersionID: versionId,
		}); err != nil {
			return diag.Errorf("%s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return ResourceApiPublishmentV2Read(clientCtx, d, meta)
}

func ResourceApiPublishmentV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = offlineApiFromSpecifiedEnv(client, instanceId, envId, apiId)
	if err != nil {
		return diag.Errorf("error offlining API: %s", err)
	}
	d.SetId("")

	return nil
}

func isPublished(versionList []apis.VersionResp, versionId string) (*apis.VersionResp, bool) {
	for _, v := range versionList {
		if v.VersionID == versionId {
			return &v, true
		}
	}
	return nil, false
}

func isLatestVersion(versionList []apis.VersionResp, versionId string) (*apis.VersionResp, bool) {
	if versionList[0].VersionID == versionId {
		return &versionList[0], true
	}
	return nil, false
}

func publishApiToSpecifiedEnv(client *golangsdk.ServiceClient, gatewayId, envId, apiId, description string) error {
	opts := apis.ManageOpts{
		GatewayID:   gatewayId,
		Action:      "online",
		EnvID:       envId,
		ApiID:       apiId,
		Description: description,
	}
	_, err := apis.ManageApi(client, opts)
	if err != nil {
		return err
	}
	return nil
}

func offlineApiFromSpecifiedEnv(client *golangsdk.ServiceClient, gatewayId, envId, apiId string) error {
	opts := apis.ManageOpts{
		GatewayID: gatewayId,
		Action:    "offline",
		EnvID:     envId,
		ApiID:     apiId,
	}
	_, err := apis.ManageApi(client, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetVersionHistory(client *golangsdk.ServiceClient, gatewayId, envId, apiId string) ([]apis.VersionResp, error) {
	history, err := apis.GetHistory(client, gatewayId, apiId, apis.ListHistoryOpts{
		EnvID: envId,
	})
	if err != nil {
		return nil, err
	}

	return history, nil
}

func getCertainPublishInfo(resp []apis.VersionResp) (*apis.VersionResp, error) {
	if len(resp) == 0 {
		return nil, fmt.Errorf("the API does not have any published information")
	}
	for _, ver := range resp {
		if ver.Status == 1 {
			return &ver, nil
		}
	}
	return nil, fmt.Errorf("unable to find any publish information for the API")
}

func setApiPublishHistory(d *schema.ResourceData, resp []apis.VersionResp) error {
	result := make([]map[string]interface{}, len(resp))
	for i, ver := range resp {
		result[i] = map[string]interface{}{
			"version_id":  ver.VersionID,
			"description": ver.Description,
		}
	}
	return d.Set("history", result)
}

func flattenResourceId(id string) (gatewayId, envId, apiId string, err error) {
	ids := strings.Split(id, "/")
	if len(ids) != 3 {
		err = fmt.Errorf("invalid ID format, want '<gateway_id>/<environment_id>/<api_id>', but '%s'", id)
		return
	}
	gatewayId = ids[0]
	envId = ids[1]
	apiId = ids[2]
	return
}

func getPublishIdByEnvId(client *golangsdk.ServiceClient, instanceId, apiId, envId string) (string, error) {
	resp, err := apis.Get(client, instanceId, apiId)
	if err != nil {
		return "", err
	}
	var (
		publishIds = strings.Split(resp.PublishID, "|")
		envIds     = strings.Split(resp.RunEnvId, "|")
	)
	log.Printf("[DEBUG] The list of publish IDs is: %#v", publishIds)
	log.Printf("[DEBUG] The list of environment IDs is: %#v", envIds)

	for i, val := range envIds {
		if val == envId {
			if len(publishIds) < i {
				return "", fmt.Errorf("the length of publish ID list is not correct, expected '%d', actual '%d'",
					len(envIds), len(publishIds))
			}
			return publishIds[i], nil
		}
	}
	return "", fmt.Errorf("the OpenTelekomCloud APIGW API is not published in this environment (%s)", envId)
}

func resourceApiPublishmentV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	mErr := multierror.Append(nil,
		d.Set("gateway_id", instanceId),
		d.Set("environment_id", envId),
		d.Set("api_id", apiId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
