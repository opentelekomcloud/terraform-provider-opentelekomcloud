package v3

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/structs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
)

func TestSetLoadBalancerFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceLoadBalancerV3().Schema, nil)
	lb := &loadbalancers.LoadBalancer{
		ID:                      "lb-id",
		Name:                    "lb",
		ProvisioningStatus:      "ACTIVE",
		OperatingStatus:         "ONLINE",
		ProjectID:               "project-id",
		Provider:                "vlb",
		Guaranteed:              true,
		IpV6VipAddress:          "2001:db8::10",
		IpV6VipSubnetID:         "ipv6-subnet-id",
		IpV6VipPortID:           "ipv6-port-id",
		L4ScaleFlavorID:         "l4-scale",
		L7ScaleFlavorID:         "l7-scale",
		EnterpriseProjectID:     "enterprise-project-id",
		ProtectionStatus:        "protection",
		ProtectionReason:        "managed",
		Pools:                   []structs.ResourceRef{{ID: "pool-id"}},
		Listeners:               []structs.ResourceRef{{ID: "listener-id"}},
		Eips:                    []loadbalancers.EipInfo{{EipID: "eip-id", EipAddress: "192.0.2.10", IpVersion: 4}},
		GlobalEips:              []loadbalancers.GlobalEipInfo{{GlobalEipID: "global-eip-id", GlobalEipAddress: "192.0.2.20", IpVersion: 4}},
		Autoscaling:             loadbalancers.AutoscalingRef{Enable: true, MinL7FlavorID: "min-l7"},
		CustomQosLimit:          loadbalancers.CustomQosLimit{L4: loadbalancers.QosLimit{Connection: 100, CPS: 10}},
		ProxyProtocolExtensions: []loadbalancers.ProxyProtocolExtension{{Extension: loadbalancers.Extension{EpID: "endpoint-id"}}},
	}

	if diagnostics := setLoadBalancerFields(d, nil, lb); diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics)
	}

	for key, expected := range map[string]interface{}{
		"provisioning_status": "ACTIVE",
		"operating_status":    "ONLINE",
		"project_id":          "project-id",
		"guaranteed":          true,
		"ipv6_vip_subnet_id":  "ipv6-subnet-id",
		"protection_status":   "protection",
		"pools.#":             1,
		"listeners.#":         1,
		"eips.#":              1,
		"global_eips.#":       1,
	} {
		if actual := d.Get(key); actual != expected {
			t.Errorf("%s: expected %#v, got %#v", key, expected, actual)
		}
	}
}

func TestDataSourceLoadBalancerV3ModernFilters(t *testing.T) {
	resource := DataSourceLoadBalancerV3()
	for _, key := range []string{
		"description", "provisioning_status", "operating_status", "guaranteed",
		"ipv6_vip_address", "eip_id", "public_ip_id", "l4_scale_flavor",
		"enterprise_project_id", "ip_version", "deletion_protection",
		"elb_subnet_type", "protection_status", "member_device_id", "member_address",
	} {
		if field, ok := resource.Schema[key]; !ok || !field.Optional {
			t.Errorf("expected %q to be an optional filter", key)
		}
	}
}
