package cfw

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/acl"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwAclRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWAclRuleV1Create,
		ReadContext:   resourceCFWAclRuleV1Read,
		UpdateContext: resourceCFWAclRuleV1Update,
		DeleteContext: resourceCFWAclRuleV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("object_id", "name"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1, 2,
				}),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sequence": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dest_rule_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsUUID,
						},
						"top": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 1,
							}),
						},
						"bottom": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 1,
							}),
						},
					},
				},
			},
			"address_type": {
				Type: schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"action_type": {
				Type: schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"status": {
				Type: schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"applications": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"applications_json_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"long_connect_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"long_connect_time_hour": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"long_connect_time_minute": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"long_connect_time_second": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"long_connect_enable": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1, 2,
				}),
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     ruleAddressSchema(),
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     ruleAddressSchema(),
			},
			"service": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1),
						},
						"protocol": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.IntInSlice([]int{
								-1, 1, 6, 17, 58,
							}),
						},
						"protocols": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     schema.TypeInt,
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dest_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_set_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_set_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_service": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.IntInSlice([]int{
											-1, 1, 6, 17, 58,
										}),
									},
									"source_port": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dest_port": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"predefined_group": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     schema.TypeString,
						},
						"service_group": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     schema.TypeString,
						},
						"service_group_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"protocols": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     schema.TypeInt,
									},
									"service_set_type": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntBetween(0, 1),
									},
									"set_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"service_set_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 3),
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_open_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ruleAddressSchema() *schema.Resource {
	ruleAddress := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 7),
			},
			"address_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_set_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_address_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_list_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_type": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 1, 2,
							}),
						},
					},
				},
			},
			"domain_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_set_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"address_set_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"predefined_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"address_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	return &ruleAddress
}

func resourceCFWAclRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := acl.CreateACLRuleOpts{
		ObjectID: d.Get("object_id").(string),
		Type:     d.Get("type").(*int),
		Rules: []acl.Rule{
			{
				Name:                   d.Get("name").(string),
				Sequence:               resourceGetSequence(d),
				AddressType:            d.Get("address_type").(*int),
				ActionType:             d.Get("action_type").(*int),
				Status:                 d.Get("status").(*int),
				Applications:           d.Get("applications").([]string),
				ApplicationsJsonString: d.Get("applications_json_string").(string),
				LongConnectTime:        d.Get("long_connect_time").(int64),
				LongConnectTimeHour:    d.Get("long_connect_time_hour").(int64),
				LongConnectTimeMinute:  d.Get("long_connect_time_minute").(int64),
				LongConnectTimeSecond:  d.Get("long_connect_time_second").(int64),
				LongConnectEnable:      d.Get("long_connect_enable").(*int),
				Description:            d.Get("description").(string),
				Direction:              d.Get("direction").(*int),
				Source:                 resourceGetRuleAddress(d, "source"),
				Destination:            resourceGetRuleAddress(d, "destination"),
				Service:                resourceGetService(d),
			},
		},
	}

	ruleIdList, err := acl.CreateACLRule(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW ACL rules from result: %w", err)
	}
	log.Printf("[DEBUG] Create CFW AC: Rules: %#v", ruleIdList)

	d.SetId(ruleIdList[0].ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceCFWAclRuleV1Read(clientCtx, d, meta)
}

func resourceCFWAclRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	objectId := d.Get("object_id").(string)
	ruleName := d.Get("name").(string)
	rule, err := acl.GetACLRule(client, objectId, ruleName)
	if err != nil {
		return fmterr.Errorf("error fetching CFW ACL rule: %w", err)
	}

	log.Printf("[DEBUG] Retrieved ACL rule %s: %#v", rule.RuleId, rule)

	mErr := multierror.Append(nil,
		d.Set("id", d.Id()),
		d.Set("type", rule.Type),
		d.Set("name", rule.Name),
		d.Set("address_type", rule.AddressType),
		d.Set("action_type", rule.ActionType),
		d.Set("status", rule.Status),
		d.Set("long_connect_time", rule.LongConnectTime),
		d.Set("long_connect_time_hour", rule.LongConnectTimeHour),
		d.Set("long_connect_time_minute", rule.LongConnectTimeMinute),
		d.Set("long_connect_time_second", rule.LongConnectTimeSecond),
		d.Set("long_connect_enable", rule.LongConnectEnable),
		d.Set("description", rule.Description),
		d.Set("direction", rule.Direction),
		d.Set("source", setRuleAddress(rule.Source)),
		d.Set("destination", setRuleAddress(rule.Destination)),
		d.Set("service", setService(rule.Service)),
		d.Set("created_date", rule.CreatedDate),
		d.Set("last_open_time", rule.LastOpenTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWAclRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	source := resourceGetRuleAddress(d, "source")
	destination := resourceGetRuleAddress(d, "destination")
	service := resourceGetService(d)

	updateOpts := acl.UpdateACLRuleOpts{
		Name:                   d.Get("name").(string),
		AddressType:            d.Get("address_type").(*int),
		ActionType:             d.Get("action_type").(*int),
		Status:                 d.Get("status").(*int),
		Applications:           d.Get("applications").([]string),
		ApplicationsJsonString: d.Get("applications_json_string").(string),
		LongConnectTime:        d.Get("long_connect_time").(int64),
		LongConnectTimeHour:    d.Get("long_connect_time_hour").(int64),
		LongConnectTimeMinute:  d.Get("long_connect_time_minute").(int64),
		LongConnectTimeSecond:  d.Get("long_connect_time_second").(int64),
		LongConnectEnable:      d.Get("long_connect_enable").(*int),
		Description:            d.Get("description").(string),
		Direction:              d.Get("direction").(*int),
		Source:                 &source,
		Destination:            &destination,
		Service:                &service,
		Type:                   d.Get("type").(*int),
	}

	rule, err := acl.UpdateACLRule(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CFW ACL rule %s: %w", rule.ID, err)
	}

	log.Printf("[DEBUG] Updated ACL rule '%s' ID: %s", rule.Name, rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWAclRuleV1Read(clientCtx, d, meta)
}

func resourceCFWAclRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW ACL rule %s", d.Id())

	err = acl.DeleteACLRule(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW ACL rule: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceGetSequence(d *schema.ResourceData) acl.OrderRuleAclDto {
	sequenceRaw := d.Get("sequence").([]interface{})[0].(map[string]interface{})
	sequence := acl.OrderRuleAclDto{
		DestRuleId: sequenceRaw["dest_rule_id"].(string),
		Top:        sequenceRaw["top"].(*int),
		Bottom:     sequenceRaw["bottom"].(*int),
	}
	return sequence
}

func getRegionList(regionListRaw []interface{}) []acl.IpRegionDto {
	var regionList []acl.IpRegionDto
	for _, v := range regionListRaw {
		regionRaw := v.(map[string]interface{})
		region := acl.IpRegionDto{
			RegionID:   regionRaw["region_id"].(string),
			RegionType: regionRaw["region_type"].(*int),
		}
		regionList = append(regionList, region)
	}
	return regionList
}

func resourceGetRuleAddress(d *schema.ResourceData, argName string) acl.RuleAddressDtoRequest {
	ruleAddressDtoRaw := d.Get(argName).([]interface{})[0].(map[string]interface{})
	regionListRaw := ruleAddressDtoRaw["region_list"].([]interface{})
	ruleAddressDto := acl.RuleAddressDtoRequest{
		Type:              ruleAddressDtoRaw["type"].(*int),
		AddressType:       ruleAddressDtoRaw["address_type"].(*int),
		Address:           ruleAddressDtoRaw["address"].(string),
		AddressSetID:      ruleAddressDtoRaw["address_set_id"].(string),
		AddressSetName:    ruleAddressDtoRaw["address_set_name"].(string),
		DomainAddressName: ruleAddressDtoRaw["domain_address_name"].(string),
		RegionListJson:    ruleAddressDtoRaw["region_list_json"].(string),
		RegionList:        getRegionList(regionListRaw),
		DomainSetID:       ruleAddressDtoRaw["domain_set_id"].(string),
		DomainSetName:     ruleAddressDtoRaw["domain_set_name"].(string),
		IPAddresses:       ruleAddressDtoRaw["ip_address"].([]string),
		AddressSetType:    ruleAddressDtoRaw["address_set_type"].(*int),
		PredefinedGroup:   ruleAddressDtoRaw["predefined_group"].([]string),
		AddressGroup:      ruleAddressDtoRaw["address_group"].([]string),
	}
	return ruleAddressDto
}

func getCustomService(customServiceRaw []interface{}) []acl.ServiceItem {
	var customService []acl.ServiceItem
	for _, v := range customServiceRaw {
		serviceRaw := v.(map[string]interface{})
		service := acl.ServiceItem{
			Protocol:    serviceRaw["protocol"].(int),
			SourcePort:  serviceRaw["source_port"].(string),
			DestPort:    serviceRaw["dest_port"].(string),
			Description: serviceRaw["description"].(string),
			Name:        serviceRaw["name"].(string),
		}
		customService = append(customService, service)
	}
	return customService
}

func getServiceGroupNames(serviceGroupNamesRaw []interface{}) []acl.ServiceGroupVO {
	var serviceGroupNames []acl.ServiceGroupVO
	for _, v := range serviceGroupNamesRaw {
		serviceGroupNameRaw := v.(map[string]interface{})
		serviceGroupName := acl.ServiceGroupVO{
			Name:           serviceGroupNameRaw["name"].(string),
			Protocols:      serviceGroupNameRaw["protocols"].([]int),
			ServiceSetType: serviceGroupNameRaw["service_set_type"].(*int),
			SetID:          serviceGroupNameRaw["set_id"].(string),
		}
		serviceGroupNames = append(serviceGroupNames, serviceGroupName)
	}
	return serviceGroupNames
}

func resourceGetService(d *schema.ResourceData) acl.RuleServiceDto {
	ruleServiceRaw := d.Get("service").([]interface{})[0].(map[string]interface{})
	customServiceRaw := ruleServiceRaw["custom_service"].([]interface{})
	serviceGroupNamesRaw := ruleServiceRaw["service_group_names"].([]interface{})

	ruleService := acl.RuleServiceDto{
		Type:              ruleServiceRaw["type"].(*int),
		Protocol:          ruleServiceRaw["protocol"].(int),
		Protocols:         ruleServiceRaw["protocols"].([]int),
		SourcePort:        ruleServiceRaw["source_port"].(string),
		DestPort:          ruleServiceRaw["dest_port"].(string),
		ServiceSetID:      ruleServiceRaw["service_set_id"].(string),
		ServiceSetName:    ruleServiceRaw["service_set_name"].(string),
		CustomService:     getCustomService(customServiceRaw),
		PredefinedGroup:   ruleServiceRaw["predefined_group"].([]string),
		ServiceGroup:      ruleServiceRaw["service_group"].([]string),
		ServiceGroupNames: getServiceGroupNames(serviceGroupNamesRaw),
		ServiceSetType:    ruleServiceRaw["service_set_type"].(*int),
	}

	return ruleService
}

func setRegionList(regionListInResp []acl.IpRegionDtoResponse) []map[string]interface{} {
	var regionList []map[string]interface{}
	for _, regionInResp := range regionListInResp {
		region := map[string]interface{}{
			"region_id":   regionInResp.RegionID,
			"region_type": regionInResp.RegionType,
		}
		regionList = append(regionList, region)
	}
	return regionList
}

func setRuleAddress(ruleAddressInResp acl.RuleAddressDtoResponse) []map[string]interface{} {
	var ruleAddressList []map[string]interface{}
	ruleAddress := map[string]interface{}{
		"type":                ruleAddressInResp.Type,
		"address_type":        ruleAddressInResp.AddressType,
		"address":             ruleAddressInResp.Address,
		"address_set_id":      ruleAddressInResp.AddressSetID,
		"address_set_name":    ruleAddressInResp.AddressSetName,
		"domain_address_name": ruleAddressInResp.DomainAddressName,
		"region_list_json":    ruleAddressInResp.RegionListJson,
		"region_list":         setRegionList(ruleAddressInResp.RegionList),
		"domain_set_id":       ruleAddressInResp.DomainSetID,
		"domain_set_name":     ruleAddressInResp.DomainSetName,
		"ip_address":          ruleAddressInResp.IPAddresses,
		"address_set_type":    ruleAddressInResp.AddressSetType,
		"address_group":       ruleAddressInResp.AddressGroup,
	}
	ruleAddressList = append(ruleAddressList, ruleAddress)
	return ruleAddressList
}

func setService(ruleServiceInResp acl.RuleServiceDtoResponse) []map[string]interface{} {
	var ruleServiceList []map[string]interface{}
	ruleService := map[string]interface{}{
		"type":                ruleServiceInResp.Type,
		"protocol":            ruleServiceInResp.Protocol,
		"protocols":           ruleServiceInResp.Protocols,
		"source_port":         ruleServiceInResp.SourcePort,
		"dest_port":           ruleServiceInResp.DestPort,
		"service_set_id":      ruleServiceInResp.ServiceSetID,
		"service_set_name":    ruleServiceInResp.ServiceSetName,
		"custom_service":      setCustomService(ruleServiceInResp.CustomService),
		"service_group":       ruleServiceInResp.ServiceGroup,
		"service_group_names": setServiceGroupNames(ruleServiceInResp.ServiceGroupNames),
		"service_set_type":    ruleServiceInResp.ServiceSetType,
	}
	ruleServiceList = append(ruleServiceList, ruleService)
	return ruleServiceList
}

func setCustomService(customServiceListInResp []acl.ServiceItemResponse) []map[string]interface{} {
	var customServiceList []map[string]interface{}
	for _, customServiceInResp := range customServiceListInResp {
		customService := map[string]interface{}{
			"protocol":    customServiceInResp.Protocol,
			"source_port": customServiceInResp.SourcePort,
			"dest_port":   customServiceInResp.DestPort,
			"description": customServiceInResp.Description,
			"name":        customServiceInResp.Name,
		}
		customServiceList = append(customServiceList, customService)
	}
	return customServiceList
}

func setServiceGroupNames(serviceGroupNamesListinResp []acl.ServiceGroupVOResponse) []map[string]interface{} {
	var serviceGroupNamesList []map[string]interface{}
	for _, serviceGroupsNamesinResp := range serviceGroupNamesListinResp {
		customService := map[string]interface{}{
			"name":             serviceGroupsNamesinResp.Name,
			"protocols":        serviceGroupsNamesinResp.Protocols,
			"service_set_type": serviceGroupsNamesinResp.ServiceSetType,
			"set_id":           serviceGroupsNamesinResp.SetID,
		}
		serviceGroupNamesList = append(serviceGroupNamesList, customService)
	}
	return serviceGroupNamesList
}
