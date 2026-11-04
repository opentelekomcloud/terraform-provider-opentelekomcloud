package dms

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQInstancesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQInstancesV2Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exact_match_name": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_failure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: rocketMQInstanceDataSourceSchema(),
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func rocketMQInstanceDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":                         {Type: schema.TypeString, Computed: true},
		"name":                       {Type: schema.TypeString, Computed: true},
		"description":                {Type: schema.TypeString, Computed: true},
		"engine":                     {Type: schema.TypeString, Computed: true},
		"engine_version":             {Type: schema.TypeString, Computed: true},
		"status":                     {Type: schema.TypeString, Computed: true},
		"type":                       {Type: schema.TypeString, Computed: true},
		"specification":              {Type: schema.TypeString, Computed: true},
		"flavor_id":                  {Type: schema.TypeString, Computed: true},
		"arch_type":                  {Type: schema.TypeString, Computed: true},
		"broker_num":                 {Type: schema.TypeInt, Computed: true},
		"storage_space":              {Type: schema.TypeInt, Computed: true},
		"used_storage_space":         {Type: schema.TypeInt, Computed: true},
		"storage_spec_code":          {Type: schema.TypeString, Computed: true},
		"disk_encrypted":             {Type: schema.TypeBool, Computed: true},
		"vpc_id":                     {Type: schema.TypeString, Computed: true},
		"vpc_name":                   {Type: schema.TypeString, Computed: true},
		"subnet_id":                  {Type: schema.TypeString, Computed: true},
		"subnet_cidr":                {Type: schema.TypeString, Computed: true},
		"security_group_id":          {Type: schema.TypeString, Computed: true},
		"security_group_name":        {Type: schema.TypeString, Computed: true},
		"available_zones":            {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"available_zone_names":       {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"maintain_begin":             {Type: schema.TypeString, Computed: true},
		"maintain_end":               {Type: schema.TypeString, Computed: true},
		"ssl_enable":                 {Type: schema.TypeBool, Computed: true},
		"tls_mode":                   {Type: schema.TypeString, Computed: true},
		"enable_acl":                 {Type: schema.TypeBool, Computed: true},
		"dns_enable":                 {Type: schema.TypeBool, Computed: true},
		"enable_publicip":            {Type: schema.TypeBool, Computed: true},
		"publicip_id":                {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"publicip_address":           {Type: schema.TypeString, Computed: true},
		"namesrv_address":            {Type: schema.TypeString, Computed: true},
		"namesrv_domain_name":        {Type: schema.TypeString, Computed: true},
		"broker_address":             {Type: schema.TypeString, Computed: true},
		"public_namesrv_address":     {Type: schema.TypeString, Computed: true},
		"public_namesrv_domain_name": {Type: schema.TypeString, Computed: true},
		"public_broker_address":      {Type: schema.TypeString, Computed: true},
		"grpc_address":               {Type: schema.TypeString, Computed: true},
		"grpc_domain_name":           {Type: schema.TypeString, Computed: true},
		"public_grpc_address":        {Type: schema.TypeString, Computed: true},
		"public_grpc_domain_name":    {Type: schema.TypeString, Computed: true},
		"cross_vpc_accesses": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"advertised_ip": {Type: schema.TypeString, Computed: true},
					"listener_ip":   {Type: schema.TypeString, Computed: true},
					"port":          {Type: schema.TypeInt, Computed: true},
					"port_id":       {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"resource_spec_code":    {Type: schema.TypeString, Computed: true},
		"enterprise_project_id": {Type: schema.TypeString, Computed: true},
		"tags":                  {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"created_at":            {Type: schema.TypeString, Computed: true},
	}
}

func dataSourceDmsRocketMQInstancesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	opts := instances.ListOpts{
		Engine:              "reliability",
		Name:                d.Get("name").(string),
		InstanceID:          d.Get("instance_id").(string),
		Status:              d.Get("status").(string),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		// The API maximum.
		Limit: 50,
	}
	if d.Get("exact_match_name").(bool) {
		opts.ExactMatchName = strconv.FormatBool(true)
	}
	if d.Get("include_failure").(bool) {
		opts.IncludeFailure = strconv.FormatBool(true)
	}

	var all []instances.Instance
	for {
		resp, err := instances.List(client, opts)
		if err != nil {
			return diag.Errorf("error getting DMS RocketMQ instance list: %s", err)
		}
		all = append(all, resp.Instances...)
		if len(resp.Instances) == 0 || len(all) >= resp.InstanceNum {
			break
		}
		opts.Offset += len(resp.Instances)
	}
	log.Printf("[DEBUG] DMS RocketMQ instance list: %#v", all)

	ids := make([]string, len(all))
	result := make([]map[string]interface{}, len(all))
	for i, v := range all {
		crossVpcAccess, err := flattenCrossVpcInfo(v.CrossVpcInfo)
		if err != nil {
			return diag.Errorf("error parsing the cross-VPC information of DMS RocketMQ instance (%s): %s", v.InstanceID, err)
		}
		ids[i] = v.InstanceID
		result[i] = map[string]interface{}{
			"id":                         v.InstanceID,
			"name":                       v.Name,
			"description":                v.Description,
			"engine":                     v.Engine,
			"engine_version":             v.EngineVersion,
			"status":                     v.Status,
			"type":                       v.Type,
			"specification":              v.Specification,
			"flavor_id":                  v.ProductID,
			"arch_type":                  v.ArchType,
			"broker_num":                 v.BrokerNum,
			"storage_space":              v.TotalStorageSpace,
			"used_storage_space":         v.UsedStorageSpace,
			"storage_spec_code":          v.StorageSpecCode,
			"disk_encrypted":             v.DiskEncrypted,
			"vpc_id":                     v.VpcID,
			"vpc_name":                   v.VpcName,
			"subnet_id":                  v.SubnetID,
			"subnet_cidr":                v.SubnetCIDR,
			"security_group_id":          v.SecurityGroupID,
			"security_group_name":        v.SecurityGroupName,
			"available_zones":            v.AvailableZones,
			"available_zone_names":       v.AvailableZoneNames,
			"maintain_begin":             v.MaintainBegin,
			"maintain_end":               v.MaintainEnd,
			"ssl_enable":                 v.SslEnable,
			"tls_mode":                   v.TlsMode,
			"enable_acl":                 v.EnableACL,
			"dns_enable":                 v.DNSEnable,
			"enable_publicip":            v.EnablePublicIP,
			"publicip_id":                splitRocketMQPublicIPs(v.PublicIpID),
			"publicip_address":           v.PublicIpAddress,
			"namesrv_address":            v.NameSrvAddress,
			"namesrv_domain_name":        v.NameSrvDomainName,
			"broker_address":             v.BrokerAddress,
			"public_namesrv_address":     v.PublicNameSrvAddress,
			"public_namesrv_domain_name": v.PublicNameSrvDomainName,
			"public_broker_address":      v.PublicBrokerAddress,
			"grpc_address":               v.GrpcAddress,
			"grpc_domain_name":           v.GrpcDomainName,
			"public_grpc_address":        v.PublicGrpcAddress,
			"public_grpc_domain_name":    v.PublicGrpcDomainName,
			"cross_vpc_accesses":         crossVpcAccess,
			"resource_spec_code":         v.ResourceSpecCode,
			"enterprise_project_id":      v.EnterpriseProjectID,
			"tags":                       common.TagsToMap(v.Tags),
			"created_at":                 v.CreatedAt,
		}
	}
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
