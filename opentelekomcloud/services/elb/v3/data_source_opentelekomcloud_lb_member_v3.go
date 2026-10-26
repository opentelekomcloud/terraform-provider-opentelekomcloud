package v3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLBMemberV3() *schema.Resource {
	fields := map[string]*schema.Schema{
		"pool_id": {Type: schema.TypeString, Required: true},
	}
	for _, key := range []string{
		"member_id", "name", "address", "subnet_id", "operating_status",
		"ip_version", "member_type", "instance_id", "availability_zone",
	} {
		fields[key] = &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true}
	}
	for _, key := range []string{"weight", "protocol_port"} {
		fields[key] = &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true}
	}
	fields["admin_state_up"] = &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true}
	for _, key := range []string{"project_id", "created_at", "updated_at"} {
		fields[key] = &schema.Schema{Type: schema.TypeString, Computed: true}
	}
	fields["status"] = memberStatusSchema()
	return &schema.Resource{
		ReadContext: dataSourceLBMemberV3Read,
		Schema:      fields,
	}
}

func buildLBMemberV3ListOpts(d *schema.ResourceData) members.ListOpts {
	filter := func(key string) []string {
		if !common.IsAttrSet(d, key) {
			return nil
		}
		return []string{d.Get(key).(string)}
	}
	opts := members.ListOpts{
		Name:             filter("name"),
		Address:          filter("address"),
		OperatingStatus:  filter("operating_status"),
		IPVersion:        filter("ip_version"),
		MemberType:       filter("member_type"),
		InstanceID:       filter("instance_id"),
		AvailabilityZone: filter("availability_zone"),
		SubnetCIDRID:     filter("subnet_id"),
	}
	if common.IsAttrSet(d, "weight") {
		opts.Weight = []int{d.Get("weight").(int)}
	}
	if common.IsAttrSet(d, "protocol_port") {
		opts.ProtocolPort = []int{d.Get("protocol_port").(int)}
	}
	if common.IsAttrSet(d, "admin_state_up") {
		value := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &value
	}
	return opts
}

func dataSourceLBMemberV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	var member *members.Member
	// Only a configured member ID should bypass filters, not a previous read's computed value.
	if common.IsAttrSet(d, "member_id") {
		member, err = members.Get(client, poolID(d), memberID(d))
	} else {
		var matches []members.Member
		matches, err = members.List(client, poolID(d), buildLBMemberV3ListOpts(d))
		if err == nil {
			if len(matches) == 0 {
				return common.DataSourceTooFewDiag
			}
			if len(matches) > 1 {
				return common.DataSourceTooManyDiag
			}
			member = &matches[0]
		}
	}
	if err != nil {
		return fmterr.Errorf("error retrieving LB pool member v3: %w", err)
	}
	if diagnostics := setLBMemberV3Fields(d, member); diagnostics.HasError() {
		return diagnostics
	}
	return diag.FromErr(common.SetComplexID(d, "pool_id", "member_id"))
}
