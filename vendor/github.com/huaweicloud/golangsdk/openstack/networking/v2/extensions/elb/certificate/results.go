package certificate

import (
	"github.com/huaweicloud/golangsdk"
)

/*
type CreateResponse struct {
	TenantId    string `json:"tenant_id"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}
*/
type Certificate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Certificate, error) {
	s := &Certificate{}
	return s, r.ExtractInto(s)
}

type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}

type getResponse struct {
	Certificates []Certificate `json:"certificates"`
	InstanceNum  string        `json:"instance_num"`
}

// GetResult represents the result of a get operation.
type GetResult struct {
	ID string
	golangsdk.Result
}

func (r GetResult) Extract() (*Certificate, error) {
	s := &getResponse{}
	err := r.ExtractInto(s)
	if err == nil {
		for _, c := range s.Certificates {
			if c.ID == r.ID {
				return &c, nil
			}
		}
		return nil, golangsdk.ErrDefault404{}
	}
	return nil, err
}

type UpdateResult struct {
	commonResult
}
