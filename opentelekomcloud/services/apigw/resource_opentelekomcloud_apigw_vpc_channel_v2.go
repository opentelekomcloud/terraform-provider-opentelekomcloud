package apigw

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/channel"
)

func ResourceAPIGWVpcChannelV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGWVpcChannelV2Create,
		ReadContext:   resourceAPIGWVpcChannelV2Read,
		UpdateContext: resourceAPIGWVpcChannelV2Update,
		DeleteContext: resourceAPIGWVpcChannelV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGWVpcChannelV2ResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"member_group": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"microservice_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"microservice_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"microservice_tags": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"member": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"is_backup": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: specialMembersOrderHash,
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold_normal": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"threshold_abnormal": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"http_codes": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_client_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"microservice": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cce_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Required: true,
									},
									"workload_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"workload_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"label_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"label_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							ExactlyOneOf: []string{"microservice.0.cse_config"},
						},
						"cse_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"service_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// If the backend server is configured through host, the hash value is calculated through hHost, and if the backend
// server is configured through ID and name, the hash value is calculated through ID.
func specialMembersOrderHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	id := m["id"].(string)
	name := m["name"].(string)
	// When configuring the backend server through host, the ID and name returned by the API both are host values.
	if id == name {
		// The values of ID and name are the same and both are host values.
		buf.WriteString(m["host"].(string))
	} else {
		buf.WriteString(id)
	}

	return hashcode.String(buf.String())
}

func buildMicroserviceTags(labels map[string]interface{}) []channel.MicroserviceTags {
	result := make([]channel.MicroserviceTags, 0, len(labels))
	for k, v := range labels {
		result = append(result, channel.MicroserviceTags{
			Key:   k,
			Value: v.(string),
		})
	}
	return result
}

func buildChannelMemberGroups(groups []interface{}) []channel.MemberGroups {
	if len(groups) < 1 {
		return nil
	}

	result := make([]channel.MemberGroups, len(groups))
	for i, val := range groups {
		group := val.(map[string]interface{})
		result[i] = channel.MemberGroups{
			Name:                group["name"].(string),
			Description:         group["description"].(string),
			Weight:              pointerto.Int(group["weight"].(int)),
			MicroserviceVersion: group["microservice_version"].(string),
			MicroservicePort:    group["microservice_port"].(int),
			MicroserviceTags:    buildMicroserviceTags(group["microservice_tags"].(map[string]interface{})),
		}
	}

	return result
}

func buildChannelMembers(members *schema.Set) []channel.Members {
	if members.Len() < 1 {
		return nil
	}

	result := make([]channel.Members, members.Len())
	for i, val := range members.List() {
		member := val.(map[string]interface{})
		result[i] = channel.Members{
			Host:            member["host"].(string),
			EcsId:           member["id"].(string),
			EcsName:         member["name"].(string),
			Weight:          pointerto.Int(member["weight"].(int)),
			IsBackup:        pointerto.Bool(member["is_backup"].(bool)),
			MemberGroupName: member["group_name"].(string),
			Status:          member["status"].(int),
			Port:            pointerto.Int(member["port"].(int)),
		}
	}

	return result
}

func buildChannelHealthCheckConfig(healthConfigs []interface{}) *channel.VpcHealthConfig {
	if len(healthConfigs) < 1 {
		return nil
	}

	healthConfig := healthConfigs[0].(map[string]interface{})
	return &channel.VpcHealthConfig{
		Protocol:           healthConfig["protocol"].(string),
		HealthyThreshold:   healthConfig["threshold_normal"].(int),
		UnhealthyThreshold: healthConfig["threshold_abnormal"].(int),
		Interval:           healthConfig["interval"].(int),
		Path:               healthConfig["path"].(string),
		Method:             healthConfig["method"].(string),
		Port:               pointerto.Int(healthConfig["port"].(int)),
		HttpCode:           healthConfig["http_codes"].(string),
		EnableClientSsl:    pointerto.Bool(healthConfig["enable_client_ssl"].(bool)),
		Status:             pointerto.Int(healthConfig["status"].(int)),
		Timeout:            healthConfig["timeout"].(int),
	}
}

func buildChannelMicroserviceConfig(microserviceConfigs []interface{}) *channel.MicroserviceConfig {
	if len(microserviceConfigs) < 1 {
		return nil
	}

	micConfig := microserviceConfigs[0].(map[string]interface{})
	if cceConfig, ok := micConfig["cce_config"]; ok {
		log.Printf("[DEBUG] The CCE configuration of the microservice is: %#v", cceConfig)
		configs := cceConfig.([]interface{})
		if len(configs) > 0 {
			details := configs[0].(map[string]interface{})
			return &channel.MicroserviceConfig{
				ServiceType: "CCE",
				CceInfo: &channel.CceInfo{
					ClusterId:    details["cluster_id"].(string),
					Namespace:    details["namespace"].(string),
					WorkloadType: details["workload_type"].(string),
					AppName:      details["workload_name"].(string),
					LabelKey:     details["label_key"].(string),
					LabelValue:   details["label_value"].(string),
				},
			}
		}
	}
	if cseConfig, ok := micConfig["cse_config"]; ok {
		log.Printf("[DEBUG] The CSE configuration of the microservice is: %#v", cseConfig)
		configs := cseConfig.([]interface{})
		if len(configs) > 0 {
			details := configs[0].(map[string]interface{})
			return &channel.MicroserviceConfig{
				ServiceType: "CSE",
				CseInfo: &channel.CseInfo{
					EngineID:  details["engine_id"].(string),
					ServiceID: details["service_id"].(string),
				},
			}
		}
	}
	return nil
}

func buildChannelCreateOpts(d *schema.ResourceData) channel.CreateOpts {
	return channel.CreateOpts{
		GatewayID:          d.Get("gateway_id").(string),
		Name:               d.Get("name").(string),
		Port:               d.Get("port").(int),
		LbAlgorithm:        d.Get("lb_algorithm").(int),
		MemberType:         d.Get("member_type").(string),
		Type:               d.Get("type").(int),
		MemberGroups:       buildChannelMemberGroups(d.Get("member_group").([]interface{})),
		Members:            buildChannelMembers(d.Get("member").(*schema.Set)),
		VpcHealthConfig:    buildChannelHealthCheckConfig(d.Get("health_check").([]interface{})),
		MicroserviceConfig: buildChannelMicroserviceConfig(d.Get("microservice").([]interface{})),
	}
}

func resourceAPIGWVpcChannelV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := buildChannelCreateOpts(d)
	v, err := channel.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud APIGW channel: %s", err)
	}
	d.SetId(v.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWVpcChannelV2Read(clientCtx, d, meta)
}

func resourceAPIGWVpcChannelV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	resp, err := channel.Get(client, gatewayId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "channel")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("port", resp.Port),
		d.Set("lb_algorithm", resp.LbAlgorithm),
		d.Set("member_type", resp.MemberType),
		d.Set("type", resp.Type),
		d.Set("member_group", flattenChannelMemberGroups(resp.MemberGroups)),
		d.Set("member", flattenChannelMembers(resp.Members)),
		d.Set("health_check", flattenHealthCheckConfig(resp.VpcHealthConfig)),
		d.Set("microservice", flattenChannelMicroserivceConfig(resp.MicroserviceConfig)),
		d.Set("created_at", resp.CreatedAt),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW channel (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func resourceAPIGWVpcChannelV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	_, err = channel.Update(client, d.Id(), buildChannelCreateOpts(d))
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud APIGW channel (%s): %s", d.Id(), err)
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWVpcChannelV2Read(clientCtx, d, meta)
}

func resourceAPIGWVpcChannelV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	if err = channel.Delete(client, gatewayId, d.Id()); err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud APIGW channel (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAPIGWVpcChannelV2ResourceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	channelId := d.Id()
	parts := strings.Split(channelId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<id>', but got '%s'",
			channelId)
	}

	d.SetId(parts[1])
	err := d.Set("gateway_id", parts[0])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving gateway ID: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func flattenMicroserviceLabels(labels []channel.MicroserviceTags) map[string]interface{} {
	result := make(map[string]interface{})
	for _, label := range labels {
		result[label.Key] = label.Value
	}
	return result
}

func flattenChannelMemberGroups(groups []channel.MemberGroupsResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(groups))
	for i, v := range groups {
		result[i] = map[string]interface{}{
			"name":                 v.Name,
			"description":          v.Description,
			"weight":               v.Weight,
			"microservice_version": v.MicroserviceVersion,
			"microservice_port":    v.MicroservicePort,
			"microservice_tagss":   flattenMicroserviceLabels(v.MicroserviceTags),
		}
	}
	return result
}

func flattenChannelMicroserivceCceConfig(cceConfig *channel.CceInfoResp) []map[string]interface{} {
	if cceConfig == nil {
		return nil
	}
	result := []map[string]interface{}{
		{
			"cluster_id":    cceConfig.ClusterId,
			"namespace":     cceConfig.Namespace,
			"workload_type": cceConfig.WorkloadType,
			"workload_name": cceConfig.AppName,
			"label_key":     cceConfig.LabelKey,
			"label_value":   cceConfig.LabelValue,
		},
	}
	return result
}

func flattenChannelMicroserivceCseConfig(cseConfig *channel.CseInfoResp) []map[string]interface{} {
	if cseConfig == nil {
		return nil
	}
	result := []map[string]interface{}{
		{
			"engine_id":  cseConfig.EngineID,
			"service_id": cseConfig.ServiceID,
		},
	}
	return result
}

func flattenHealthCheckConfig(healthConfig *channel.VpcHealthConfigResp) []map[string]interface{} {
	if healthConfig == nil {
		return nil
	}
	result := []map[string]interface{}{
		{
			"protocol":           strings.ToUpper(healthConfig.Protocol),
			"threshold_normal":   healthConfig.HealthyThreshold,
			"threshold_abnormal": healthConfig.UnhealthyThreshold,
			"interval":           healthConfig.Interval,
			"timeout":            healthConfig.Timeout,
			"path":               healthConfig.Path,
			"method":             healthConfig.Method,
			"port":               healthConfig.Port,
			"http_codes":         healthConfig.HttpCode,
			"enable_client_ssl":  healthConfig.EnableClientSsl,
			"status":             healthConfig.Status,
		},
	}
	return result
}

func flattenChannelMicroserivceConfig(microserviceConfig *channel.MicroserviceConfigResp) []map[string]interface{} {
	if microserviceConfig == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, 1)
	switch microserviceConfig.ServiceType {
	case "CCE":
		result = append(result, map[string]interface{}{
			"cce_config": flattenChannelMicroserivceCceConfig(microserviceConfig.CceInfo),
		})
	case "CSE":
		result = append(result, map[string]interface{}{
			"cse_config": flattenChannelMicroserivceCseConfig(microserviceConfig.CseInfo),
		})
	}
	return result
}

func flattenChannelMembers(members []channel.MembersResp) []map[string]interface{} {
	if len(members) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(members))
	for i, member := range members {
		result[i] = map[string]interface{}{
			"host":       member.Host,
			"id":         member.EcsId,
			"name":       member.EcsName,
			"weight":     member.Weight,
			"is_backup":  member.IsBackup,
			"group_name": member.MemberGroupName,
			"status":     member.Status,
			"port":       member.Port,
		}
	}
	return result
}
