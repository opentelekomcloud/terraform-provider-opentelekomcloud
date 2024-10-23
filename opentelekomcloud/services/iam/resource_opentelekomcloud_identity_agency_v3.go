package iam

import (
	"bytes"
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
	newAgency "github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/agency"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/agency"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/domains"
	sdkprojects "github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"
	sdkroles "github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/roles"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityAgencyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityAgencyV3Create,
		ReadContext:   resourceIdentityAgencyV3Read,
		UpdateContext: resourceIdentityAgencyV3Update,
		DeleteContext: resourceIdentityAgencyV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delegated_domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_role": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 25,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"all_projects": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceIdentityAgencyProRoleHash,
			},
			"domain_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 25,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIdentityAgencyProRoleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["project"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["project"].(string)))
	}

	if m["all_projects"] != nil {
		buf.WriteString(fmt.Sprintf("%t-", m["all_projects"].(bool)))
	}

	r := m["roles"].(*schema.Set).List()
	s := make([]string, len(r))
	for i, item := range r {
		s[i] = item.(string)
	}
	buf.WriteString(strings.Join(s, "-"))

	return hashcode.String(buf.String())
}

func agencyClient(d *schema.ResourceData, config *cfg.Config) (*golangsdk.ServiceClient, error) {
	c, err := config.IdentityV3Client(config.GetRegion(d))

	if err != nil {
		return nil, err
	}

	c.Endpoint = strings.Replace(c.Endpoint, "v3", "v3.0", 1)
	return c, nil
}

func listProjectsOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	old := client.Endpoint
	defer func() { client.Endpoint = old }()
	client.Endpoint = strings.Replace(old, "v3.0", "v3", 1)

	opts := sdkprojects.ListOpts{
		DomainID: domainID,
	}
	allPages, err := sdkprojects.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("list projects failed, err=%s", err)
	}

	all, err := sdkprojects.ExtractProjects(allPages)
	if err != nil {
		return nil, fmt.Errorf("extract projects failed, err=%s", err)
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		r[item.Name] = item.ID
	}
	log.Printf("[TRACE] projects = %#v\n", r)
	return r, nil
}

func listRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	old := client.Endpoint
	defer func() { client.Endpoint = old }()
	client.Endpoint = strings.Replace(old, "v3.0", "v3", 1)

	opts := sdkroles.ListOpts{
		DomainID: domainID,
	}
	allPages, err := sdkroles.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("list roles failed, err=%s", err)
	}

	all, err := sdkroles.ExtractRoles(allPages)
	if err != nil {
		return nil, fmt.Errorf("extract roles failed, err=%s", err)
	}
	if len(all) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		dn, ok := item.Extra["display_name"].(string)
		if ok {
			r[dn] = item.ID
		} else {
			log.Printf("[DEBUG] Can not retrieve role:%#v", item)
		}
	}
	log.Printf("[TRACE] list roles = %#v, len=%d\n", r, len(r))
	return r, nil
}

func allRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	roles, err := listRolesOfDomain("", client)
	if err != nil {
		return nil, fmt.Errorf("error listing global roles, err=%s", err)
	}

	customRoles, err := listRolesOfDomain(domainID, client)
	if err != nil {
		return nil, fmt.Errorf("error listing domain's custom roles, err=%s", err)
	}

	if roles == nil {
		return customRoles, nil
	}

	if customRoles == nil {
		return roles, nil
	}

	for k, v := range customRoles {
		roles[k] = v
	}
	return roles, nil
}

func getDomainID(config *cfg.Config, client *golangsdk.ServiceClient) (string, error) {
	if config.DomainID != "" {
		return config.DomainID, nil
	}

	name := config.DomainName
	if name == "" {
		return "", fmt.Errorf("the required domain name was missed")
	}

	old := client.Endpoint
	defer func() { client.Endpoint = old }()
	client.Endpoint = strings.Replace(old, "v3.0", "v3", 1) + "auth/"

	opts := domains.ListOpts{
		Name: name,
	}
	allPages, err := domains.List(client, &opts).AllPages()
	if err != nil {
		return "", fmt.Errorf("list domains failed, err=%s", err)
	}

	all, err := domains.ExtractDomains(allPages)
	if err != nil {
		return "", fmt.Errorf("extract domains failed, err=%s", err)
	}

	count := len(all)
	switch count {
	case 0:
		err := &golangsdk.ErrResourceNotFound{}
		err.ResourceType = "iam"
		err.Name = name
		return "", err
	case 1:
		return all[0].ID, nil
	default:
		err := &golangsdk.ErrMultipleResourcesFound{}
		err.ResourceType = "iam"
		err.Name = name
		err.Count = count
		return "", err
	}
}

func changeToPRPair(prs *schema.Set) (r map[string]bool) {
	r = make(map[string]bool)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})

		pn := pr["project"].(string)
		pa := pr["all_projects"].(bool)
		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			if pa {
				r["all_projects"+"|"+role.(string)] = true
			} else {
				r[pn+"|"+role.(string)] = true
			}
		}
	}
	return
}

func diffChangeOfProjectRole(old, newv *schema.Set) (delete, add []string) {
	delete = make([]string, 0)
	add = make([]string, 0)

	oldprs := changeToPRPair(old)
	newprs := changeToPRPair(newv)

	for k := range oldprs {
		if _, ok := newprs[k]; !ok {
			delete = append(delete, k)
		}
	}

	for k := range newprs {
		if _, ok := oldprs[k]; !ok {
			add = append(add, k)
		}
	}
	return
}

func resourceIdentityAgencyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prs := d.Get("project_role").(*schema.Set)
	drs := d.Get("domain_roles").(*schema.Set)
	if prs.Len() == 0 && drs.Len() == 0 {
		return fmterr.Errorf("one or both of project_role and domain_roles must be input")
	}

	config := meta.(*cfg.Config)
	client, err := agencyClient(d, config)
	if err != nil {
		return fmterr.Errorf("error creating client: %s", err)
	}

	domainID, err := getDomainID(config, client)
	if err != nil {
		return fmterr.Errorf("error getting the domain id, err=%s", err)
	}

	opts := agency.CreateOpts{
		Name:            d.Get("name").(string),
		DomainID:        domainID,
		DelegatedDomain: d.Get("delegated_domain_name").(string),
		Description:     d.Get("description").(string),
	}
	log.Printf("[DEBUG] Create Identity-Agency Options: %#v", opts)
	a, err := agency.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating Identity-Agency: %s", err)
	}

	d.SetId(a.ID)

	roles, err := allRolesOfDomain(domainID, client)
	if err != nil {
		return fmterr.Errorf("error querying the roles, err=%s", err)
	}

	projects, err := listProjectsOfDomain(domainID, client)
	if err != nil {
		return fmterr.Errorf("error querying the projects, err=%s", err)
	}

	agencyID := a.ID
	for _, v := range prs.List() {
		var pn, pid string
		pr := v.(map[string]interface{})
		rs := pr["roles"].(*schema.Set)
		pa := pr["all_projects"].(bool)

		pn = pr["project"].(string)
		pid, ok := projects[pn]
		if !ok && !pa {
			return fmterr.Errorf("the project(%s) does not exist", pn)
		}

		for _, role := range rs.List() {
			r := role.(string)
			rid, ok := roles[r]
			if !ok {
				return fmterr.Errorf("the role(%s) does not exist", r)
			}

			if pa {
				domainId, err := getDomainID(config, client)
				if err != nil {
					return fmterr.Errorf("error getting the domain id, err=%s", err)
				}
				err = newAgency.GrantAgencyAllProjects(client, domainId, agencyID, rid)
				if err != nil {
					return fmterr.Errorf("error attaching role(%s) to agency(%s) for all projects, err=%s",
						rid, agencyID, err)
				}
			} else {
				err = agency.AttachRoleByProject(client, agencyID, pid, rid).ExtractErr()
				if err != nil {
					return fmterr.Errorf("error attaching role(%s) by project{%s} to agency(%s), err=%s",
						rid, pid, agencyID, err)
				}
			}
		}
	}

	for _, role := range drs.List() {
		r := role.(string)
		rid, ok := roles[r]
		if !ok {
			return fmterr.Errorf("the role(%s) does not exist", r)
		}

		err = agency.AttachRoleByDomain(client, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmterr.Errorf("error attaching role(%s) by domain{%s} to agency(%s), err=%s",
				rid, domainID, agencyID, err)
		}
	}

	return resourceIdentityAgencyV3Read(ctx, d, meta)
}

func resourceIdentityAgencyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := agencyClient(d, config)
	if err != nil {
		return fmterr.Errorf("error creating client: %s", err)
	}

	a, err := agency.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Identity-Agency")
	}
	log.Printf("[DEBUG] Retrieved Identity-Agency %s: %#v", d.Id(), a)

	mErr := multierror.Append(
		d.Set("name", a.Name),
		d.Set("delegated_domain_name", a.DelegatedDomainName),
		d.Set("description", a.Description),
		d.Set("duration", a.Duration),
		d.Set("expire_time", a.ExpireTime),
		d.Set("create_time", a.CreateTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	projects, err := listProjectsOfDomain(a.DomainID, client)
	if err != nil {
		return fmterr.Errorf("error querying the projects, err=%s", err)
	}
	agencyID := d.Id()
	prs := schema.Set{F: resourceIdentityAgencyProRoleHash}
	for pn, pid := range projects {
		roles, err := agency.ListRolesAttachedOnProject(client, agencyID, pid).ExtractRoles()
		if err != nil && !common.IsResourceNotFound(err) {
			return fmterr.Errorf("error querying the roles attached on project(%s), err=%s", pn, err)
		}
		if len(roles) == 0 {
			continue
		}
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.Extra["display_name"])
		}
		prs.Add(map[string]interface{}{
			"project": pn,
			"roles":   &v,
		})
	}

	domainId, err := getDomainID(config, client)
	if err != nil {
		return fmterr.Errorf("error quering domain id, err=%s", err)
	}

	allRoles, err := newAgency.ListAgencyAllProjects(client, domainId, agencyID)
	if err != nil {
		return fmterr.Errorf("error querying the roles given all_projects scope, err=%s", err)
	}
	if len(allRoles.Roles) > 0 {
		domainRoles, err := allRolesOfDomain(domainId, client)
		if err != nil {
			return fmterr.Errorf("error querying domain roles, err=%s", err)
		}
		rl := schema.Set{F: schema.HashString}
		for _, v := range allRoles.Roles {
			displayName, _ := findKeyByValue(domainRoles, v.Id)
			rl.Add(displayName)
		}
		prs.Add(map[string]interface{}{
			"all_projects": true,
			"roles":        &rl,
		})
	}

	err = d.Set("project_role", &prs)
	if err != nil {
		log.Printf("[ERROR]Set project_role failed, err=%s", err)
	}

	roles, err := agency.ListRolesAttachedOnDomain(client, agencyID, a.DomainID).ExtractRoles()
	if err != nil && !common.IsResourceNotFound(err) {
		return fmterr.Errorf("error querying the roles attached on domain, err=%s", err)
	}
	if len(roles) != 0 {
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.Extra["display_name"])
		}
		err = d.Set("domain_roles", &v)
		if err != nil {
			log.Printf("[ERROR]Set domain_roles failed, err=%s", err)
		}
	}

	return nil
}

func resourceIdentityAgencyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := agencyClient(d, config)
	if err != nil {
		return fmterr.Errorf("error creating client: %s", err)
	}

	aID := d.Id()

	if d.HasChange("delegated_domain_name") || d.HasChange("description") {
		updateOpts := agency.UpdateOpts{
			DelegatedDomain: d.Get("delegated_domain_name").(string),
			Description:     d.Get("description").(string),
		}
		log.Printf("[DEBUG] Updating Identity-Agency %s with options: %#v", aID, updateOpts)
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err := agency.Update(client, aID, updateOpts).Extract()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmterr.Errorf("error updating Identity-Agency %s: %s", aID, err)
		}
	}

	domainID := ""
	var roles map[string]string
	if d.HasChange("project_role") || d.HasChange("domain_roles") {
		domainID, err = getDomainID(config, client)
		if err != nil {
			return fmterr.Errorf("error getting the domain id, err=%s", err)
		}

		roles, err = allRolesOfDomain(domainID, client)
		if err != nil {
			return fmterr.Errorf("error querying the roles, err=%s", err)
		}
	}

	if d.HasChange("project_role") {
		projects, err := listProjectsOfDomain(domainID, client)
		if err != nil {
			return fmterr.Errorf("error querying the projects, err=%s", err)
		}

		o, n := d.GetChange("project_role")
		deleteprs, addprs := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))
		for _, v := range deleteprs {
			pr := strings.Split(v, "|")
			rid, ok := roles[pr[1]]
			if !ok {
				return fmterr.Errorf("the role(%s) does not exist", pr[1])
			}
			if pr[0] == "all_projects" {
				domainID, err = getDomainID(config, client)
				if err != nil {
					return fmterr.Errorf("error getting the domain id, err=%s", err)
				}
				err = newAgency.RemoveAgencyAllProjects(client, domainID, aID, rid)
				if err != nil && !common.IsResourceNotFound(err) {
					return fmterr.Errorf("error detaching role(%s) for all projects from agency(%s), err=%s",
						rid, aID, err)
				}
			} else {
				pid, ok := projects[pr[0]]
				if !ok {
					return fmterr.Errorf("the project(%s) does not exist", pr[0])
				}

				err = agency.DetachRoleByProject(client, aID, pid, rid).ExtractErr()
				if err != nil && !common.IsResourceNotFound(err) {
					return fmterr.Errorf("error detaching role(%s) by project{%s} from agency(%s), err=%s",
						rid, pid, aID, err)
				}
			}
		}

		for _, v := range addprs {
			pr := strings.Split(v, "|")
			rid, ok := roles[pr[1]]
			if !ok {
				return fmterr.Errorf("the role(%s) does not exist", pr[1])
			}
			if pr[0] == "all_projects" {
				domainID, err = getDomainID(config, client)
				if err != nil {
					return fmterr.Errorf("error getting the domain id, err=%s", err)
				}
				err = newAgency.GrantAgencyAllProjects(client, domainID, aID, rid)
				if err != nil && !common.IsResourceNotFound(err) {
					return fmterr.Errorf("error attaching role(%s) for all projects from agency(%s), err=%s",
						rid, aID, err)
				}
			} else {
				pid, ok := projects[pr[0]]
				if !ok {
					return fmterr.Errorf("the project(%s) does not exist", pr[0])
				}

				err = agency.AttachRoleByProject(client, aID, pid, rid).ExtractErr()
				if err != nil {
					return fmterr.Errorf("error attaching role(%s) by project{%s} to agency(%s), err=%s",
						rid, pid, aID, err)
				}
			}
		}
	}

	if d.HasChange("domain_roles") {
		o, n := d.GetChange("domain_roles")
		oldr := o.(*schema.Set)
		newr := n.(*schema.Set)

		for _, r := range oldr.Difference(newr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmterr.Errorf("the role(%s) does not exist", r.(string))
			}

			err = agency.DetachRoleByDomain(client, aID, domainID, rid).ExtractErr()
			if err != nil && !common.IsResourceNotFound(err) {
				return fmterr.Errorf("error detaching role(%s) by domain{%s} from agency(%s), err=%s",
					rid, domainID, aID, err)
			}
		}

		for _, r := range newr.Difference(oldr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmterr.Errorf("the role(%s) does not exist", r.(string))
			}

			err = agency.AttachRoleByDomain(client, aID, domainID, rid).ExtractErr()
			if err != nil {
				return fmterr.Errorf("error attaching role(%s) by domain{%s} to agency(%s), err=%s",
					rid, domainID, aID, err)
			}
		}
	}
	return resourceIdentityAgencyV3Read(ctx, d, meta)
}

func resourceIdentityAgencyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := agencyClient(d, config)
	if err != nil {
		return fmterr.Errorf("error creating client: %s", err)
	}

	rID := d.Id()
	log.Printf("[DEBUG] Deleting Identity-Agency %s", rID)

	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := agency.Delete(client, rID).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if common.IsResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable Identity-Agency: %s", rID)
			return nil
		}
		return fmterr.Errorf("error deleting Identity-Agency %s: %s", rID, err)
	}

	return nil
}
