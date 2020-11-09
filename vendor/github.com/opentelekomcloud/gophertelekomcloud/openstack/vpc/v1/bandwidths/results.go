package bandwidths

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type BandWidth struct {
	// Specifies the bandwidth name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name"`

	// Specifies the bandwidth size. The value ranges from 1 Mbit/s to
	// 300 Mbit/s.
	Size int `json:"size"`

	// Specifies the bandwidth ID, which uniquely identifies the
	// bandwidth.
	ID string `json:"id"`

	// Specifies whether the bandwidth is shared or exclusive. The
	// value can be PER or WHOLE.
	ShareType string `json:"share_type"`

	// Specifies the elastic IP address of the bandwidth.  The
	// bandwidth, whose type is set to WHOLE, supports up to 20 elastic IP addresses. The
	// bandwidth, whose type is set to PER, supports only one elastic IP address.
	PublicipInfo []PublicIpinfo `json:"publicip_info"`

	// Specifies the tenant ID of the user.
	TenantId string `json:"tenant_id"`

	// Specifies the bandwidth type.
	BandwidthType string `json:"bandwidth_type"`

	// Specifies the charging mode (by traffic or by bandwidth).
	ChargeMode string `json:"charge_mode"`

	// Specifies the billing information.
	BillingInfo string `json:"billing_info"`

	// Enterprise project id
	EnterpriseProjectID string `json:"enterprise_project_id"`

	// Specifies the status of bandwidth
	Status string `json:"status"`
}

type PublicIpinfo struct {
	// Specifies the tenant ID of the user.
	PublicipId string `json:"publicip_id"`

	// Specifies the elastic IP address.
	PublicipAddress string `json:"publicip_address"`

	// Specifies the elastic IP v6 address.
	Publicipv6Address string `json:"publicipv6_address"`

	// Specifies the elastic IP version.
	IPVersion int `json:"ip_version"`

	// Specifies the elastic IP address type. The value can be
	// 5_telcom, 5_union, or 5_bgp.
	PublicipType string `json:"publicip_type"`
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*BandWidth, error) {
	var entity BandWidth
	err := r.ExtractIntoStructPtr(&entity, "bandwidth")
	return &entity, err
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*[]BandWidth, error) {
	var list []BandWidth
	err := r.ExtractIntoSlicePtr(&list, "bandwidths")
	return &list, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*BandWidth, error) {
	var entity BandWidth
	err := r.ExtractIntoStructPtr(&entity, "bandwidth")
	return &entity, err
}
func (r BandWidthPage) IsEmpty() (bool, error) {
	list, err := ExtractBandWidths(r)
	return len(list) == 0, err
}

type BandWidthPage struct {
	pagination.LinkedPageBase
}

func ExtractBandWidths(r pagination.Page) ([]BandWidth, error) {
	var s struct {
		BandWidths []BandWidth `json:"bandwidths"`
	}
	err := r.(BandWidthPage).ExtractInto(&s)
	return s.BandWidths, err
}
func (r BandWidthPage) NextPageURL() (string, error) {
	s, err := ExtractBandWidths(r)
	if err != nil {
		return "", err
	}
	return r.WrapNextPageURL(s[len(s)-1].ID)
}
