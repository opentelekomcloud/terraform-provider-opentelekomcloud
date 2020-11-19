package addons

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
)

type Addon struct {
	// API type, fixed value Addon
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion" required:"true"`
	// Metadata of an Addon
	Metadata MetaData `json:"metadata" required:"true"`
	// Specifications of an Addon
	Spec Spec `json:"spec" required:"true"`
	// Status of an Addon
	Status Status `json:"status"`
}

// Metadata required to create an addon
type MetaData struct {
	// Addon unique name
	Name string `json:"name"`
	// Addon unique Id
	Id string `json:"uid"`
	// Addon tag, key/value pair format
	Labels map[string]string `json:"lables"`
	// Addon annotation, key/value pair format
	Annotations map[string]string `json:"annotaions"`
}

// Specifications to create an addon
type Spec struct {
	// For the addon version.
	Version string `json:"version" required:"true"`
	// Cluster ID.
	ClusterID string `json:"clusterID" required:"true"`
	// Addon Template Name.
	AddonTemplateName string `json:"addonTemplateName" required:"true"`
	// Addon Template Type.
	AddonTemplateType string `json:"addonTemplateType" required:"true"`
	// Addon Template Labels.
	AddonTemplateLables []string `json:"addonTemplateLables,omitempty"`
	// Addon Description.
	Description string `json:"description" required:"true"`
	// Addon Parameters
	Values Values `json:"values" required:"true"`
}

type Status struct {
	// The state of the addon
	Status string `json:"status"`
	// Reasons for the addon to become current
	Reason string `json:"reason"`
	// Error Message
	Message string `json:"message"`
	// The target versions of the addon
	TargetVersions []string `json:"targetVersions"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts an Addon.
func (r commonResult) Extract() (*Addon, error) {
	var s Addon
	err := r.ExtractInto(&s)
	return &s, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as an Addon.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as an Addon.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as an Addon.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

type ListTemplateResult struct {
	golangsdk.Result
}

type SupportVersion struct {
	// Cluster type that supports the add-on template
	ClusterType string `json:"clusterType"`
	// Cluster versions that support the add-on template,
	// the parameter value is a regular expression
	ClusterVersion []string `json:"clusterVersion"`
}

type Version struct {
	// Add-on version
	Version string `json:"version"`
	// Add-on installation parameters
	Input map[string]interface{} `json:"input"`
	// Whether the add-on version is a stable release
	Stable bool `json:"stable"`
	// Cluster versions that support the add-on template
	SupportVersions []SupportVersion `json:"supportVersions"`
	// Creation time of the add-on instance
	CreationTimestamp string `json:"creationTimestamp"`
	// Time when the add-on instance was updated
	UpdateTimestamp string `json:"updateTimestamp"`
}

type AddonSpec struct {
	// Template type (helm or static).
	Type string `json:"type" required:"true"`
	// Whether the add-on is installed by default
	Require bool `json:"require" required:"true"`
	// Group to which the template belongs
	Labels []string `json:"labels" required:"true"`
	// URL of the logo image
	LogoURL string `json:"logoURL" required:"true"`
	// URL of the readme file
	ReadmeURL string `json:"readmeURL" required:"true"`
	// Template description
	Description string `json:"description" required:"true"`
	// Template version details
	Versions []Version `json:"versions" required:"true"`
}

type AddonTemplate struct {
	// API type, fixed value Addon
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion" required:"true"`
	// Metadata of an Addon
	Metadata MetaData `json:"metadata" required:"true"`
	// Specifications of an Addon
	Spec AddonSpec `json:"spec" required:"true"`
}

type AddonTemplateList struct {
	// API type, fixed value Addon
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion" required:"true"`
	// Add-on template list
	Items []AddonTemplate `json:"items" required:"true"`
}

// Extract is a function that accepts a result and extracts an Addon.
func (r ListTemplateResult) Extract() (*AddonTemplateList, error) {
	var s AddonTemplateList
	err := r.ExtractInto(&s)
	return &s, err
}
