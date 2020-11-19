package addons

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAddonCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new addon
type CreateOpts struct {
	// API type, fixed value Addon
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion" required:"true"`
	// Metadata required to create an addon
	Metadata CreateMetadata `json:"metadata" required:"true"`
	// specifications to create an addon
	Spec RequestSpec `json:"spec" required:"true"`
}

type CreateMetadata struct {
	Annotations CreateAnnotations `json:"annotations" required:"true"`
}

type CreateAnnotations struct {
	AddonInstallType string `json:"addon.install/type" required:"true"`
}

// Specifications to create an addon
type RequestSpec struct {
	// For the addon version.
	Version string `json:"version" required:"true"`
	// Cluster ID.
	ClusterID string `json:"clusterID" required:"true"`
	// Addon Template Name.
	AddonTemplateName string `json:"addonTemplateName" required:"true"`
	// Addon Parameters
	Values Values `json:"values" required:"true"`
}

type Values struct {
	Basic    map[string]interface{} `json:"basic" required:"true"`
	Advanced map[string]interface{} `json:"custom,omitempty"`
}

// ToAddonCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToAddonCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// addon.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder, clusterId string) (r CreateResult) {
	b, err := opts.ToAddonCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c, clusterId), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular addon based on its unique ID.
func Get(c *golangsdk.ServiceClient, id, clusterId string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id, clusterId), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type UpdateOptsBuilder interface {
	ToAddonUpdateMap() (map[string]interface{}, error)
}

type UpdateMetadata struct {
	// Add-on annotations in the format of key-value pairs.
	// For add-on upgrade, the value is fixed at {"addon.upgrade/type":"upgrade"}.
	Annotations UpdateAnnotations `json:"annotations" required:"true"`
	// Add-on labels in the format of key-value pairs.
	Labels map[string]string `json:"metadata,omitempty"`
}

type UpdateAnnotations struct {
	AddonUpdateType string `json:"addon.update/type" required:"true"`
}

type UpdateOpts struct {
	// API type, fixed value Addon
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion" required:"true"`
	// Metadata required to create an addon
	Metadata UpdateMetadata `json:"metadata" required:"true"`
	// specifications to create an addon
	Spec RequestSpec `json:"spec" required:"true"`
}

func (opts UpdateOpts) ToAddonUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Update(c *golangsdk.ServiceClient, id, clusterId string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToAddonUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id, clusterId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular addon based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id, clusterId string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id, clusterId), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

type ListOptsBuilder interface {
	ToAddonListQuery() (string, error)
}

type ListOpts struct {
	Name string `q:"addon_template_name"`
}

func (opts ListOpts) ToAddonListQuery() (string, error) {
	u, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func ListTemplates(c *golangsdk.ServiceClient, clusterID string, opts ListOptsBuilder) (r ListTemplateResult) {
	url := templatesURL(c, clusterID)
	if opts != nil {
		q, err := opts.ToAddonListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += q
	}
	_, r.Err = c.Get(url, r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
