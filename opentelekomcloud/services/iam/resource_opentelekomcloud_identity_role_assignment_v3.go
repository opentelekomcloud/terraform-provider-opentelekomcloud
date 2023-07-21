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
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"all_projects": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"project_id"},
				Optional:      true,
				ForceNew:      true,
			},
		},
	}
}

func resourceIdentityRoleAssignmentV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	domainID := d.Get("domain_id").(string)
	groupID := d.Get("group_id").(string)
	projectID := d.Get("project_id").(string)
	roleID := d.Get("role_id").(string)

	if !d.Get("all_projects").(bool) {
		opts := roles.AssignOpts{
			DomainID:  domainID,
			GroupID:   groupID,
			ProjectID: projectID,
		}

		err = roles.Assign(identityClient, roleID, opts).ExtractErr()
		if err != nil {
			return fmterr.Errorf("error assigning role: %s", err)
		}
	} else {
		err = roles.UpdateGroupAllProjects(identityClient, domainID, groupID, roleID)
		if err != nil {
			return fmterr.Errorf("error assigning role to all projects: %s", err)
		}
	}

	d.SetId(buildRoleAssignmentID(domainID, projectID, groupID, roleID))

	clientCtx := common.CtxWithClient(ctx, identityClient, keyClientV3)
	return resourceIdentityRoleAssignmentV3Read(clientCtx, d, meta)
}

func resourceIdentityRoleAssignmentV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	var mErr *multierror.Error
	var roleAssignmentId string

	domainID, projectID, groupID, _ := ExtractRoleAssignmentID(d.Id())

	if !d.Get("all_projects").(bool) {
		roleAssignmentId, err = getRoleAssignment(identityClient, d)
	} else {
		err = getAllProjectsRoleAssignment(identityClient, d)
		if err != nil {
			diag.FromErr(err)
		}
		_, _, _, roleAssignmentId = ExtractRoleAssignmentID(d.Id())
	}
	if err != nil {
		return fmterr.Errorf("error getting role assignment: %s", err)
	}

	log.Printf("[DEBUG] Retrieved OpenStack role assignment ID: %s", roleAssignmentId)
	mErr = multierror.Append(
		d.Set("domain_id", domainID),
		d.Set("project_id", projectID),
		d.Set("group_id", groupID),
		d.Set("role_id", roleAssignmentId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIdentityRoleAssignmentV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	domainID, projectID, groupID, roleID := ExtractRoleAssignmentID(d.Id())
	if !d.Get("all_projects").(bool) {
		opts := roles.UnassignOpts{
			DomainID:  domainID,
			GroupID:   groupID,
			ProjectID: projectID,
		}
		err = roles.Unassign(identityClient, roleID, opts).ExtractErr()
		if err != nil {
			return fmterr.Errorf("error unassigning role: %s", err)
		}
	} else {
		err = roles.RemoveGroupAllProjects(identityClient, domainID, groupID, roleID)
		if err != nil {
			return fmterr.Errorf("error assigning role to all projects: %s", err)
		}
	}

	return nil
}

func getRoleAssignment(identityClient *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	domainID, projectID, groupID, roleID := ExtractRoleAssignmentID(d.Id())

	opts := roles.ListAssignmentsOpts{
		GroupID:        groupID,
		ScopeDomainID:  domainID,
		ScopeProjectID: projectID,
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

	return assignment.ID, err
}

func getAllProjectsRoleAssignment(identityClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	domainID, _, groupID, roleID := ExtractRoleAssignmentID(d.Id())

	err := roles.CheckGroupAllProjects(identityClient, domainID, groupID, roleID)

	return err
}

// Role assignments have no ID in OpenStack. Build an ID out of the IDs that make up the role assignment
func buildRoleAssignmentID(domainID, projectID, groupID, roleID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", domainID, projectID, groupID, roleID)
}

func ExtractRoleAssignmentID(roleAssignmentID string) (string, string, string, string) {
	split := strings.Split(roleAssignmentID, "/")
	if len(split) != 4 {
		// Otherwise previous role assignments in tfstate file won't be backwards compatible
		return split[0], split[1], split[2], split[4]
	}
	return split[0], split[1], split[2], split[3]
}
