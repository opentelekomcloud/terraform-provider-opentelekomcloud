package opentelekomcloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/domains"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/policies"
)

func resourceWafDomainV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafDomainV1Create,
		Read:   resourceWafDomainV1Read,
		Update: resourceWafDomainV1Update,
		Delete: resourceWafDomainV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"front_protocol": {
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use `client_protocol` instead",
						},
						"client_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"back_protocol": {
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use `server_protocol` instead",
						},
						"server_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"proxy": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"sip_header_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"", "default", "cloudflare", "akamai", "custom"}, true),
			},
			"sip_header_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"txt_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getAllServers(d *schema.ResourceData) ([]domains.ServerOpts, error) {
	var serverOpts []domains.ServerOpts

	servers := d.Get("server").([]interface{})
	for _, v := range servers {
		server := v.(map[string]interface{})
		cProtocol, err := firstOneSet(server, "client_protocol", "front_protocol")
		if err != nil {
			return nil, err
		}
		sProtocol, err := firstOneSet(server, "server_protocol", "back_protocol")
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(server["port"].(string))
		if err != nil {
			return nil, fmt.Errorf("invalid WAF domain server port: %s", err)
		}
		v := domains.ServerOpts{
			ClientProtocol: cProtocol.(string),
			ServerProtocol: sProtocol.(string),
			Address:        server["address"].(string),
			Port:           port,
		}
		serverOpts = append(serverOpts, v)
	}

	log.Printf("[DEBUG] getAllServers: %#v", serverOpts)
	return serverOpts, nil
}

func resourceWafDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	var hosts []string
	if v, ok := d.GetOk("policy_id"); ok {
		policyId := v.(string)
		policy, err := policies.Get(wafClient, policyId).Extract()
		if err != nil {
			return fmt.Errorf("error retrieving OpenTelekomCloud Waf Policy %s: %s", policyId, err)
		}
		hosts = append(hosts, policy.Hosts...)
	}

	v := d.Get("sip_header_list").([]interface{})
	headers := make([]string, len(v))
	for i, v := range v {
		headers[i] = v.(string)
	}

	proxy := d.Get("proxy").(bool)
	servers, err := getAllServers(d)
	if err != nil {
		return fmt.Errorf("error parsing servers: %s", err)
	}
	createOpts := domains.CreateOpts{
		HostName:      d.Get("hostname").(string),
		CertificateId: d.Get("certificate_id").(string),
		Server:        servers,
		Proxy:         &proxy,
		SipHeaderName: d.Get("sip_header_name").(string),
		SipHeaderList: headers,
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	domain, err := domains.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomcomCloud WAF Domain: %s", err)
	}

	log.Printf("[DEBUG] Waf domain created: %#v", domain)
	d.SetId(domain.Id)

	if v, ok := d.GetOk("policy_id"); ok {
		var updateHostsOpts policies.UpdateHostsOpts
		policyId := v.(string)
		hosts = append(hosts, d.Id())
		updateHostsOpts.Hosts = hosts
		log.Printf("[DEBUG] Waf policy update Hosts: %#v", hosts)

		_, err = policies.UpdateHosts(wafClient, policyId, updateHostsOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating OpenTelekomCloud WAF Policy Hosts: %s", err)
		}
	}

	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	n, err := domains.Get(wafClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error retrieving OpenTelekomCloud Waf Domain: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("hostname", n.HostName),
		d.Set("certificate_id", n.CertificateId),
		d.Set("proxy", n.Proxy),
		d.Set("sip_header_name", n.SipHeaderName),
		d.Set("sip_header_list", n.SipHeaderList),
		d.Set("access_code", n.AccessCode),
		d.Set("cname", n.Cname),
		d.Set("txt_code", n.TxtCode),
		d.Set("sub_domain", n.SubDomain),
	)
	if n.PolicyID != "" {
		mErr = multierror.Append(mErr, d.Set("policy_id", n.PolicyID))
	}
	mErr = multierror.Append(mErr,
		d.Set("protect_status", n.ProtectStatus),
		d.Set("access_status", n.AccessStatus),
		d.Set("protocol", n.Protocol),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting WAF fields: %s", err)
	}

	servers := make([]map[string]interface{}, len(n.Server))
	for i, server := range n.Server {
		servers[i] = make(map[string]interface{})
		servers[i]["front_protocol"] = server.ClientProtocol
		servers[i]["client_protocol"] = server.ClientProtocol
		servers[i]["back_protocol"] = server.ServerProtocol
		servers[i]["server_protocol"] = server.ServerProtocol
		servers[i]["address"] = server.Address
		servers[i]["port"] = strconv.Itoa(server.Port)
	}
	return d.Set("server", servers)
}

func resourceWafDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts domains.UpdateOpts

	if d.HasChange("certificate_id") {
		updateOpts.CertificateId = d.Get("certificate_id").(string)
	}
	if d.HasChange("server") {
		servers, err := getAllServers(d)
		if err != nil {
			return fmt.Errorf("error parsing servers: %s", err)
		}
		updateOpts.Server = servers
	}
	if d.HasChange("proxy") {
		proxy := d.Get("proxy").(bool)
		updateOpts.Proxy = &proxy
	}
	if d.HasChange("sip_header_name") {
		updateOpts.SipHeaderName = d.Get("sip_header_name").(string)
	}
	if d.HasChange("sip_header_list") {
		v := d.Get("sip_header_list").([]interface{})
		headers := make([]string, len(v))
		for i, v := range v {
			headers[i] = v.(string)
		}
		updateOpts.SipHeaderList = headers
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = domains.Update(wafClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating OpenTelekomCloud WAF Domain: %s", err)
	}
	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	err = domains.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud WAF Domain: %s", err)
	}

	d.SetId("")
	return nil
}
