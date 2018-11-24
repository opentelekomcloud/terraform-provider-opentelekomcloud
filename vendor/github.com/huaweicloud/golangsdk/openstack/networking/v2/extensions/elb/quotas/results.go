package quotas

import (
	"github.com/huaweicloud/golangsdk"
)

type ResourceQuota struct {
	Type  string `json:"type"`
	Used  int    `json:"used"`
	Quota int    `json:"quota"`
	Max   int    `json:"max"`
	Min   int    `json:"min"`
}

type resources struct {
	Resources []ResourceQuota `json:"resources"`
}

type quotaInfo struct {
	Quotas resources `json:"quotas"`
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() ([]ResourceQuota, error) {
	s := &quotaInfo{}
	err := r.ExtractInto(s)
	if err == nil {
		return s.Quotas.Resources, nil
	}
	return nil, err
}
