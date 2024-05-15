package apigw

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/key"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPISignatureAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSignatureAssociateV2Create,
		ReadContext:   resourceSignatureAssociateV2Read,
		UpdateContext: resourceSignatureAssociateV2Update,
		DeleteContext: resourceSignatureAssociateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSignatureAssociateV2ImportState,
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
			"signature_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"publish_ids": {
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

func resourceSignatureAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	signatureId := d.Get("signature_id").(string)
	publishIds := common.ExpandToStringListBySet(d.Get("publish_ids").(*schema.Set))
	opts := key.BindOpts{
		GatewayID:  gatewayId,
		SignID:     signatureId,
		PublishIds: publishIds,
	}
	_, err = key.BindKey(client, opts)
	if err != nil {
		return diag.Errorf("error binding OpenTelekomCloud apigw signature to the APIs: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", gatewayId, signatureId))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    signatureBindingRefreshFunc(client, gatewayId, signatureId, publishIds),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the OpenTelekomCloud apigw signature (%s) to bind: %s", gatewayId, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceSignatureAssociateV2Read(clientCtx, d, meta)
}

func resourceSignatureAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	resp, err := key.ListAPIBoundKeys(client, key.ListBoundOpts{
		GatewayID: d.Get("gateway_id").(string),
		SignID:    d.Get("signature_id").(string),
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Signature association")
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	return diag.FromErr(d.Set("publish_ids", flattenApiPublishIdsForSignature(resp)))
}

func resourceSignatureAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	signId := d.Get("signature_id").(string)
	oldRaw, newRaw := d.GetChange("publish_ids")

	addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))

	if rmSet.Len() > 0 {
		err = unbindSignatureFromApis(ctx, client, d, common.ExpandToStringListBySet(rmSet))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		// If the target (published) API already has a signature, this update will replace the signature.
		publishIds := common.ExpandToStringListBySet(addSet)
		opts := key.BindOpts{
			GatewayID:  gatewayId,
			SignID:     signId,
			PublishIds: publishIds,
		}
		_, err = key.BindKey(client, opts)
		if err != nil {
			return diag.Errorf("error binding OpenTelekomCloud apigw signature to the APIs: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING"},
			Target:     []string{"COMPLETED"},
			Refresh:    signatureBindingRefreshFunc(client, gatewayId, signId, publishIds),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the OpenTelekomCloud apigw signature (%s) to bind: %s", gatewayId, err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceSignatureAssociateV2Read(clientCtx, d, meta)
}

func resourceSignatureAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = unbindSignatureFromApis(ctx, client, d, common.ExpandToStringListBySet(d.Get("publish_ids").(*schema.Set)))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSignatureAssociateV2ImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<signature_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("gateway_id", parts[0]),
		d.Set("signature_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func signatureUnbindingRefreshFunc(client *golangsdk.ServiceClient, gatewayId, signId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := buildSignBindApiListOpts(gatewayId, signId)
		resp, err := key.ListAPIBoundKeys(client, opts)
		if err != nil {
			return resp, "", err
		}
		bindPublishIds := flattenApiPublishIdsForSignature(resp)
		if common.IsSliceContainsAnyAnotherSliceElement(bindPublishIds, publishIds, false, true) {
			return resp, "PENDING", nil
		}
		return resp, "COMPLETED", nil
	}
}

func unbindSignatureFromApis(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, rmList []string) error {
	gatewayId := d.Get("gateway_id").(string)
	signId := d.Get("signature_id").(string)

	opts := buildSignBindApiListOpts(gatewayId, signId)
	resp, err := key.ListAPIBoundKeys(client, opts)
	if err != nil {
		return fmt.Errorf("error getting binding OpenTelekomCloud apigw APIs based on signature (%s): %s", signId, err)
	}

	for _, val := range resp {
		if common.StrSliceContains(rmList, val.PublishID) {
			err = key.UnbindKey(client, gatewayId, val.ID)
			if err != nil {
				return fmt.Errorf("an error occurred during unbind signature: %s", err)
			}
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: signatureUnbindingRefreshFunc(client, gatewayId, signId, rmList),
		Timeout: d.Timeout(schema.TimeoutDelete),
		// In most cases, the unbind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the unbind operation completed: %s", err)
	}
	return nil
}

func buildSignBindApiListOpts(gatewayId, signatureId string) key.ListBoundOpts {
	return key.ListBoundOpts{
		GatewayID: gatewayId,
		SignID:    signatureId,
	}
}

func flattenApiPublishIdsForSignature(apiList []key.BindSignResp) []string {
	if len(apiList) < 1 {
		return nil
	}

	result := make([]string, len(apiList))
	for i, val := range apiList {
		result[i] = val.PublishID
	}
	return result
}

func signatureBindingRefreshFunc(client *golangsdk.ServiceClient, gatewayId, signatureId string, publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := key.ListUnboundKeys(client, key.ListUnbindOpts{
			GatewayID: gatewayId,
			SignID:    signatureId,
		})
		if err != nil {
			return resp, "ERROR", err
		}
		var result = make([]interface{}, 0)
		for _, val := range resp {
			result = append(result, val.PublishID)
		}

		if common.IsSliceContainsAnyAnotherSliceElement(common.ExpandToStringList(result), publishIds, false, true) {
			return result, "PENDING", nil
		}
		return result, "COMPLETED", nil
	}
}
