package credentials

import "github.com/opentelekomcloud/gophertelekomcloud"

type Status string

const parentElement = "credential"

type ListOptsBuilder interface {
	ToCredentialListQuery() (string, error)
}

type ListOpts struct {
	UserID string `json:"user_id,omitempty"`
}

func (opts ListOpts) ToCredentialListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) (l ListResult) {
	q, err := opts.ToCredentialListQuery()
	if err != nil {
		l.Err = err
		return
	}
	_, l.Err = client.Get(listURL(client)+q, &l.Body, nil)
	return
}

func Get(client *golangsdk.ServiceClient, credentialID string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, credentialID), &r.Body, nil)
	return
}

type CreateOptsBuilder interface {
	ToCredentialCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	UserID      string `json:"user_id"`
	Description string `json:"description"`
}

func (opts CreateOpts) ToCredentialCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, parentElement)
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToCredentialUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts UpdateOpts) ToCredentialUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, parentElement)
}

func Update(client *golangsdk.ServiceClient, credentialID string, opts UpdateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCredentialUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, credentialID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, credentialID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, credentialID), nil)
	return
}
