package v3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceListenerV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListenerV3Read,

		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"keep_alive_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// output
			"insert_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forward_elb_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_port": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_for_port": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"forwarded_host": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_retry_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sni_container_refs": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"memory_retry_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"advanced_forwarding": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sni_match_algo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceListenerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if id := d.Get("id"); id != "" {
		listener, err := listeners.Get(client, id.(string)).Extract()
		if err != nil {
			return fmterr.Errorf("error finding listener by ID: %w", err)
		}
		d.SetId(listener.ID)
		return setLBListenerFields(d, listener)
	}

	opts := listeners.ListOpts{
		ProtocolPort:            common.IntSlice(d.Get("protocol_port")),
		Protocol:                toProtocolSlice(common.StrSlice(d.Get("protocol"))),
		Description:             common.StrSlice(d.Get("description")),
		DefaultTLSContainerRef:  common.StrSlice(d.Get("default_tls_container_ref")),
		ClientCATLSContainerRef: common.StrSlice(d.Get("client_ca_tls_container_ref")),
		DefaultPoolID:           common.StrSlice(d.Get("default_pool_id")),
		Name:                    common.StrSlice(d.Get("name")),
		LoadBalancerID:          common.StrSlice(d.Get("loadbalancer_id")),
		TLSCiphersPolicy:        common.StrSlice(d.Get("tls_ciphers_policy")),
		MemberAddress:           common.StrSlice(d.Get("member_address")),
		MemberDeviceID:          common.StrSlice(d.Get("member_device_id")),
		MemberTimeout:           common.IntSlice(d.Get("member_timeout")),
		ClientTimeout:           common.IntSlice(d.Get("client_timeout")),
		KeepAliveTimeout:        common.IntSlice(d.Get("keep_alive_timeout")),
	}
	pages, err := listeners.List(client, opts).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing LB listeners v3: %w", err)
	}
	listenerSlice, err := listeners.ExtractListeners(pages)
	if err != nil {
		return fmterr.Errorf("error extracting listeners: %w", err)
	}
	if len(listenerSlice) < 1 {
		return common.DataSourceTooFewDiag
	}
	if len(listenerSlice) > 1 {
		return common.DataSourceTooManyDiag
	}
	listener := listenerSlice[0]
	d.SetId(listener.ID)
	return setLBListenerFields(d, &listener)
}

func toProtocolSlice(src []string) []listeners.Protocol {
	res := make([]listeners.Protocol, len(src))
	for i, v := range src {
		res[i] = listeners.Protocol(v)
	}
	return res
}
