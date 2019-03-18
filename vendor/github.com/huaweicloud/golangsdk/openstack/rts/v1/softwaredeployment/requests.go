package softwaredeployment

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API.
type ListOpts struct {
	//Specifies the ID of the instance deployed by the software configuration.
	ServerId string `q:"server_id"`
	//Specifies the ID of this deployment resource.
	Id string
	//Specifies the ID of the software configuration resource running on an instance.
	ConfigId string
	//Specifies the current status of deployment resources. Valid values include COMPLETE, IN_PROGRESS, and FAILED.
	Status string
	// Specifies the stack action that triggers this deployment resource.
	Action string
}

// List returns collection of
// Software Deployment. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those Software Deployment that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Deployment, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := rootURL(client) + q.String()
	pages, err := pagination.NewPager(client, u, func(r pagination.PageResult) pagination.Page {
		return DeploymentPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allConfigs, err := ExtractDeployments(pages)
	if err != nil {
		return nil, err
	}

	return FilterDeployments(allConfigs, opts)
}

func FilterDeployments(deployments []Deployment, opts ListOpts) ([]Deployment, error) {

	var refinedDeployments []Deployment
	var matched bool
	m := map[string]interface{}{}

	if opts.Id != "" {
		m["Id"] = opts.Id
	}
	if opts.ServerId != "" {
		m["ServerId"] = opts.ServerId
	}
	if opts.ConfigId != "" {
		m["ConfigId"] = opts.ConfigId
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if opts.Action != "" {
		m["Action"] = opts.Action
	}

	if len(m) > 0 && len(deployments) > 0 {
		for _, deployment := range deployments {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&deployment, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedDeployments = append(refinedDeployments, deployment)
			}
		}

	} else {
		refinedDeployments = deployments
	}

	return refinedDeployments, nil
}

func getStructField(v *Deployment, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSoftwareDeploymentCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new Software Deployment. There are
// no required values.
type CreateOpts struct {
	// Specifies the stack action that triggers this deployment resource.
	Action string `json:"action,omitempty"`
	//Specifies the ID of the software configuration resource running on an instance.
	ConfigId string `json:"config_id" required:"true"`
	//Specifies input data stored in the form of a key-value pair.
	InputValues map[string]interface{} `json:"input_values,omitempty"`
	//Specifies the ID of the instance deployed by the software configuration.
	ServerId string `json:"server_id" required:"true"`
	//Specifies the ID of the authenticated tenant who can perform operations on the deployment resources.
	TenantId string `json:"stack_user_project_id,omitempty"`
	//Specifies the current status of deployment resources. Valid values include COMPLETE, IN_PROGRESS, and FAILED.
	Status string `json:"status,omitempty"`
	//Specifies the cause of the current deployment resource status.
	StatusReason string `json:"status_reason,omitempty"`
}

// ToSoftwareDeploymentCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToSoftwareDeploymentCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new Software Deployment
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSoftwareDeploymentCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular software Deployment based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

//UpdateOptsBuilder is an interface by which can be able to build the request body of software deployment.
type UpdateOptsBuilder interface {
	ToSoftwareDeploymentUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	// Specifies the stack action that triggers this deployment resource.
	Action string `json:"action,omitempty"`
	//Specifies the ID of the software configuration resource running on an instance.
	ConfigId string `json:"config_id" required:"true"`
	//Specifies input data stored in the form of a key-value pair.
	InputValues map[string]interface{} `json:"input_values,omitempty"`
	//Specifies output data stored in the form of a key-value pair.
	OutputValues map[string]interface{} `json:"output_values" required:"true"`
	//Specifies the current status of deployment resources. Valid values include COMPLETE, IN_PROGRESS, and FAILED.
	Status string `json:"status,omitempty"`
	//Specifies the cause of the current deployment resource status.
	StatusReason string `json:"status_reason,omitempty"`
}

//ToSoftwareDeploymentUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToSoftwareDeploymentUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Update is a method which can be able to update the name of software deployment.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSoftwareDeploymentUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular Software Deployment based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
