package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/roles"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityRoleAssignmentV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRoleAssignmentV3Create,
		ReadContext:   resourceIdentityRoleAssignmentV3Read,
		DeleteContext: resourceIdentityRoleAssignmentV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"project_id"},
				Optional:      true,
				ForceNew:      true,
			},

			"group_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"user_id"},
				Optional:      true,
				ForceNew:      true,
			},

			"project_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"domain_id"},
				Optional:      true,
				ForceNew:      true,
			},

			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"group_id"},
				Optional:      true,
				ForceNew:      true,
			},
		},
	}
}

func resourceIdentityRoleAssignmentV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	domainID := d.Get("domain_id").(string)
	groupID := d.Get("group_id").(string)
	projectID := d.Get("project_id").(string)
	roleID := d.Get("role_id").(string)
	userID := d.Get("user_id").(string)
	opts := roles.AssignOpts{
		DomainID:  domainID,
		GroupID:   groupID,
		ProjectID: projectID,
		UserID:    userID,
	}

	err = roles.Assign(identityClient, roleID, opts).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error assigning role: %s", err)
	}

	d.SetId(buildRoleAssignmentID(domainID, projectID, groupID, userID, roleID))

	return resourceIdentityRoleAssignmentV3Read(ctx, d, meta)
}

func resourceIdentityRoleAssignmentV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	roleAssignment, err := getRoleAssignment(identityClient, d)
	if err != nil {
		return fmterr.Errorf("error getting role assignment: %s", err)
	}
	domainID, projectID, groupID, userID, _ := ExtractRoleAssignmentID(d.Id())

	log.Printf("[DEBUG] Retrieved OpenStack role assignment: %#v", roleAssignment)
	mErr := multierror.Append(
		d.Set("domain_id", domainID),
		d.Set("project_id", projectID),
		d.Set("group_id", groupID),
		d.Set("user_id", userID),
		d.Set("role_id", roleAssignment.ID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIdentityRoleAssignmentV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	domainID, projectID, groupID, userID, roleID := ExtractRoleAssignmentID(d.Id())
	opts := roles.UnassignOpts{
		DomainID:  domainID,
		GroupID:   groupID,
		ProjectID: projectID,
		UserID:    userID,
	}
	err = roles.Unassign(identityClient, roleID, opts).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error unassigning role: %s", err)
	}

	return nil
}

func getRoleAssignment(identityClient *golangsdk.ServiceClient, d *schema.ResourceData) (roles.RoleAssignment, error) {
	domainID, projectID, groupID, userID, roleID := ExtractRoleAssignmentID(d.Id())

	opts := roles.ListAssignmentsOpts{
		GroupID:        groupID,
		ScopeDomainID:  domainID,
		ScopeProjectID: projectID,
		UserID:         userID,
	}

	pager := roles.ListAssignments(identityClient, opts)
	var assignment roles.RoleAssignment

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoleAssignments(page)
		if err != nil {
			return false, err
		}

		for _, a := range assignmentList {
			if a.ID == roleID {
				assignment = a
				return false, nil
			}
		}

		return true, nil
	})

	return assignment, err
}

// Role assignments have no ID in OpenStack. Build an ID out of the IDs that make up the role assignment
func buildRoleAssignmentID(domainID, projectID, groupID, userID, roleID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", domainID, projectID, groupID, userID, roleID)
}

func ExtractRoleAssignmentID(roleAssignmentID string) (string, string, string, string, string) {
	split := strings.Split(roleAssignmentID, "/")
	return split[0], split[1], split[2], split[3], split[4]
}
