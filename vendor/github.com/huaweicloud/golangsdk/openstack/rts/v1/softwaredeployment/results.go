package softwaredeployment

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Deployment struct {
	// Specifies the stack action that triggers this deployment resource.
	Action string `json:"action"`
	//Specifies the ID of the software Deployments resource running on an instance.
	ConfigId string `json:"config_id"`
	//Specifies the creation time. The timestamp format is ISO 8601: CCYY-MM-DDThh:mm:ss±hh:mm
	CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
	//Specifies the ID of this deployment resource.
	Id string `json:"id"`
	//Specifies input data stored in the form of a key-value pair.
	InputValues map[string]interface{} `json:"input_values"`
	//Specifies output data stored in the form of a key-value pair.
	OutputValues map[string]interface{} `json:"output_values"`
	//Specifies the ID of the instance deployed by the software Deployments.
	ServerId string `json:"server_id"`
	//Specifies the current status of deployment resources. Valid values include COMPLETE, IN_PROGRESS, and FAILED.
	Status string `json:"status"`
	//Specifies the cause of the current deployment resource status.
	StatusReason string `json:"status_reason"`
	//Specifies the updated time. The timestamp format is ISO 8601: CCYY-MM-DDThh:mm:ss±hh:mm
	UpdatedTime golangsdk.JSONRFC3339NoZ `json:"updated_time"`
}

// DeploymentPage is the page returned by a pager when traversing over a
// collection of Software Deployments.
type DeploymentPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of Software Deployments has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r DeploymentPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"software_deployments_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a DeploymentPage struct is empty.
func (r DeploymentPage) IsEmpty() (bool, error) {
	is, err := ExtractDeployments(r)
	return len(is) == 0, err
}

// ExtractDeployments accepts a Page struct, specifically a DeploymentPage struct,
// and extracts the elements into a slice of Software Deployments structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractDeployments(r pagination.Page) ([]Deployment, error) {
	var s struct {
		Deployments []Deployment `json:"software_deployments"`
	}
	err := (r.(DeploymentPage)).ExtractInto(&s)
	return s.Deployments, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a Software Deployments.
func (r commonResult) Extract() (*Deployment, error) {
	var s struct {
		SoftwareDeployment *Deployment `json:"software_deployment"`
	}
	err := r.ExtractInto(&s)
	return s.SoftwareDeployment, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Software Deployments.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Software Deployments.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Software Deployments.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
