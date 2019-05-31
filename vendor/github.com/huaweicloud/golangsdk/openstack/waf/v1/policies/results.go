package policies

import (
	"github.com/huaweicloud/golangsdk"
)

type Policy struct {
	//Policy ID
	Id string `json:"id"`
	//Policy Name
	Name string `json:"name"`
	//Protective Action
	Action Action `json:"action"`
	//Protection Switches
	Options Options `json:"options"`
	//Protection Level
	Level int `json:"level"`
	//Detection Mode
	FullDetection bool `json:"full_detection"`
	//Domain IDs
	Hosts []string `json:"hosts"`
}

type Action struct {
	//Protective Action
	Category string `json:"category" required:"true"`
}

type Options struct {
	//Whether Basic Web Protection is enabled
	WebAttack *bool `json:"webattack,omitempty"`
	//Whether General Check in Basic Web Protection is enabled
	Common *bool `json:"common,omitempty"`
	//Whether the master crawler detection switch in Basic Web Protection is enabled
	Crawler *bool `json:"crawler,omitempty"`
	//Whether the Search Engine switch in Basic Web Protection is enabled
	CrawlerEngine *bool `json:"crawler_engine,omitempty"`
	//Whether the Scanner switch in Basic Web Protection is enabled
	CrawlerScanner *bool `json:"crawler_scanner,omitempty"`
	//Whether the Script Tool switch in Basic Web Protection is enabled
	CrawlerScript *bool `json:"crawler_script,omitempty"`
	//Whether detection of other crawlers in Basic Web Protection is enabled
	CrawlerOther *bool `json:"crawler_other,omitempty"`
	//Whether webshell detection in Basic Web Protection is enabled
	WebShell *bool `json:"webshell,omitempty"`
	//Whether CC Attack Protection is enabled
	Cc *bool `json:"cc,omitempty"`
	//Whether Precise Protection is enabled
	Custom *bool `json:"custom,omitempty"`
	//Whether Blacklist and Whitelist is enabled
	WhiteblackIp *bool `json:"whiteblackip,omitempty"`
	//Whether Data Masking is enabled
	Privacy *bool `json:"privacy,omitempty"`
	//Whether False Alarm Masking is enabled
	Ignore *bool `json:"ignore,omitempty"`
	//Whether Web Tamper Protection is enabled
	AntiTamper *bool `json:"antitamper,omitempty"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a policy.
func (r commonResult) Extract() (*Policy, error) {
	var response Policy
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Policy.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Policy.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Policy.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
