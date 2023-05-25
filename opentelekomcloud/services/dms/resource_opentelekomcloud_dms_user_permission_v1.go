package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/permissions"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsUsersPermissionV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsUserPermissionV1Create,
		ReadContext:   resourceDmsUserPermissionV1Read,
		DeleteContext: resourceDmsUserPermissionV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"access_policy": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"all", "pub", "sub",
							}, false),
						},
						"owner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"topic_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getPolicy(d *schema.ResourceData) []permissions.CreatePolicy {
	var refinedPolicies []permissions.CreatePolicy
	policiesRaw := d.Get("policies").([]interface{})
	for _, policyRaw := range policiesRaw {
		policy := policyRaw.(map[string]interface{})
		refinedPolicy := permissions.CreatePolicy{
			UserName:     policy["username"].(string),
			AccessPolicy: policy["access_policy"].(string),
		}
		refinedPolicies = append(refinedPolicies, refinedPolicy)
	}
	return refinedPolicies
}

func resourceDmsUserPermissionV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV11Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := permissions.CreateOpts{
		Name:     d.Get("topic_name").(string),
		Policies: getPolicy(d),
	}

	instanceId := d.Get("instance_id").(string)

	err = permissions.Create(client, instanceId, []permissions.CreateOpts{
		createOpts,
	})
	if err != nil {
		return fmterr.Errorf("error assigning OpenTelekomCloud DMSv1 permissions: %w", err)
	}

	// Store the topic name == ID
	d.SetId(createOpts.Name)

	return resourceDmsUserPermissionV1Read(ctx, d, meta)
}

func resourceDmsUserPermissionV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV11Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	listPolicies, err := getPermissionsFromList(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS permission")
	}

	mErr := multierror.Append(
		d.Set("policies", flattenPolicies(listPolicies.Policies)),
		d.Set("topic_name", listPolicies.Name),
		d.Set("topic_type", listPolicies.TopicType),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDmsUserPermissionV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV11Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	topicName := d.Get("topic_name").(string)
	deleteOpts := permissions.CreateOpts{
		Name:     topicName,
		Policies: []permissions.CreatePolicy{},
	}

	_ = permissions.Create(client, d.Get("instance_id").(string), []permissions.CreateOpts{
		deleteOpts,
	})

	// commented because of 500 internal error that is received sometimes on permission deletion
	// if err != nil {
	// 	return fmterr.Errorf("error deleting OpenTelekomCloud DMSv1 permissions: %w", err)
	// }

	log.Printf("[DEBUG] DMS permissions for topic %s deactivated.", topicName)
	d.SetId("")
	return nil
}

func getPermissionsFromList(client *golangsdk.ServiceClient, d *schema.ResourceData) (permissions.Permissions, error) {
	var policies permissions.Permissions
	v, err := permissions.List(client, d.Get("instance_id").(string), d.Get("topic_name").(string))
	if err != nil {
		return policies, err
	}

	return *v, nil
}

func flattenPolicies(rawPolicies []permissions.Policy) []map[string]interface{} {
	var policies []map[string]interface{}
	for _, rawPolicy := range rawPolicies {
		v := map[string]interface{}{
			"username":      rawPolicy.UserName,
			"access_policy": rawPolicy.AccessPolicy,
			"owner":         rawPolicy.Owner,
		}
		policies = append(policies, v)
	}
	return policies
}
