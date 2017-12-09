package alarmrule

import (
	"log"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToAlarmRuleCreateMap() (map[string]interface{}, error)
}

type dimensionOpts struct {
	Name  string `json:"name" required:"true"`
	Value string `json:"value" required:"true"`
}

type metricOpts struct {
	Namespace  string          `json:"namespace" required:"true"`
	MetricName string          `json:"metric_name" required:"true"`
	Dimensions []dimensionOpts `json:"dimensions" required:"true"`
}

type conditionOpts struct {
	Period             int    `json:"period" required:"true"`
	Filter             string `json:"filter" required:"true"`
	ComparisonOperator string `json:"comparison_operator" required:"true"`
	Value              int    `json:"value" required:"true"`
	Unit               string `json:"unit,omitempty"`
	Count              int    `json:"count" required:"true"`
}

type actionOpts struct {
	Type             string   `json:"type" required:"true"`
	NotificationList []string `json:"notification_list" required:"true"`
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	AlarmName               string        `json:"alarm_name" required:"true"`
	AlarmDescription        string        `json:"alarm_description,omitempty"`
	Metric                  metricOpts    `json:"metric" required:"true"`
	Condition               conditionOpts `json:"condition" required:"true"`
	AlarmActions            []actionOpts  `json:"alarm_actions,omitempty"`
	InsufficientdataActions []actionOpts  `json:"insufficientdata_actions,omitempty"`
	OkActions               []actionOpts  `json:"ok_actions,omitempty"`
	AlarmEnabled            bool          `json:"alarm_enabled,omitempty"`
	AlarmActionEnabled      bool          `json:"alarm_action_enabled,omitempty"`
}

func (opts CreateOpts) ToAlarmRuleCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type realActionOpts struct {
	Type             string   `json:"type" required:"true"`
	NotificationList []string `json:"notificationList" required:"true"`
}

type createOpts struct {
	AlarmName               string           `json:"alarm_name" required:"true"`
	AlarmDescription        string           `json:"alarm_description,omitempty"`
	Metric                  metricOpts       `json:"metric" required:"true"`
	Condition               conditionOpts    `json:"condition" required:"true"`
	AlarmActions            []realActionOpts `json:"alarm_actions,omitempty"`
	InsufficientdataActions []realActionOpts `json:"insufficientdata_actions,omitempty"`
	OkActions               []realActionOpts `json:"ok_actions,omitempty"`
	AlarmEnabled            bool             `json:"alarm_enabled,omitempty"`
	AlarmActionEnabled      bool             `json:"alarm_action_enabled,omitempty"`
}

// ToAlarmRuleCreateMap casts a CreateOpts struct to a map.
func (opts createOpts) ToAlarmRuleCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func copyActionOpts(src []actionOpts) []realActionOpts {
	if len(src) == 0 {
		return nil
	}

	dest := make([]realActionOpts, len(src), len(src))
	for i, s := range src {
		d := &dest[i]
		d.Type = s.Type
		d.NotificationList = s.NotificationList
	}
	log.Printf("[DEBUG] copyActionOpts:: src = %#v, dest = %#v", src, dest)
	return dest
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	opt := opts.(CreateOpts)
	opts1 := createOpts{
		AlarmName:               opt.AlarmName,
		AlarmDescription:        opt.AlarmDescription,
		Metric:                  opt.Metric,
		Condition:               opt.Condition,
		AlarmActions:            copyActionOpts(opt.AlarmActions),
		InsufficientdataActions: copyActionOpts(opt.InsufficientdataActions),
		OkActions:               copyActionOpts(opt.OkActions),
		AlarmEnabled:            opt.AlarmEnabled,
		AlarmActionEnabled:      opt.AlarmActionEnabled,
	}
	b, err := opts1.ToAlarmRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create AlarmRule url:%q, body=%#v, opt=%#v", rootURL(c), b, opts1)
	reqOpt := &gophercloud.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToAlarmRuleUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Optional. Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`
	// Optional. Human-readable description for the Loadbalancer.
	Description string `json:"description,omitempty"`

	BandWidth int `json:"bandwidth,omitempty"`
	// Optional. The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp bool `json:"admin_state_up,omitempty"`
}

// ToAlarmRuleUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToAlarmRuleUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

/*
// Update is an operation which modifies the attributes of the specified AlarmRule.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) (r elb.JobResult) {
	b, err := opts.ToAlarmRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular AlarmRule based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) (r elb.JobResult) {
	reqOpt := &gophercloud.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.DeleteAndGetResponse(resourceURL(c, id), &r.Body, reqOpt)
	return
}
*/
