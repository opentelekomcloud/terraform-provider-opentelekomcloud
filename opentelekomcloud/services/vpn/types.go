package vpn

import (
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
)

// VpnIPSecPolicyCreateOpts represents the attributes used when creating a new IPSec policy.
type VpnIPSecPolicyCreateOpts struct {
	ipsecpolicies.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnServiceCreateOpts represents the attributes used when creating a new VPN service.
type VpnServiceCreateOpts struct {
	services.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnEndpointGroupCreateOpts represents the attributes used when creating a new endpoint group.
type VpnEndpointGroupCreateOpts struct {
	endpointgroups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnIKEPolicyCreateOpts represents the attributes used when creating a new IKE policy.
type VpnIKEPolicyCreateOpts struct {
	ikepolicies.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnIKEPolicyLifetimeCreateOpts represents the attributes used when creating a new lifetime for an IKE policy.
type VpnIKEPolicyLifetimeCreateOpts struct {
	ikepolicies.LifetimeCreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// VpnSiteConnectionCreateOpts represents the attributes used when creating a new IPSec site connection.
type VpnSiteConnectionCreateOpts struct {
	siteconnections.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}
