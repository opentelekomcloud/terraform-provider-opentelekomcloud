package opentelekomcloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/domains"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/policies"
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
							Type:     schema.TypeString,
							Required: true,
						},
						"back_protocol": {
							Type:     schema.TypeString,
							Required: true,
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

func getAllServers(d *schema.ResourceData) []domains.ServerOpts {
	var serverOpts []domains.ServerOpts

	servers := d.Get("server").([]interface{})
	for _, v := range servers {
		server := v.(map[string]interface{})

		v := domains.ServerOpts{
			ClientProtocol: server["front_protocol"].(string),
			ServerProtocol: server["back_protocol"].(string),
			Address:        server["address"].(string),
			Port:           server["port"].(string),
		}
		serverOpts = append(serverOpts, v)
	}

	log.Printf("[DEBUG] getAllServers: %#v", serverOpts)
	return serverOpts
}

func resourceWafDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	hosts := []string{}
	if hasFilledOpt(d, "policy_id") {
		policy_id := d.Get("policy_id").(string)
		policy, err := policies.Get(wafClient, policy_id).Extract()
		if err != nil {
			return fmt.Errorf("Error retrieving OpenTelekomCloud Waf Policy %s: %s", policy_id, err)
		}
		hosts = append(hosts, policy.Hosts...)
	}

	v := d.Get("sip_header_list").([]interface{})
	headers := make([]string, len(v))
	for i, v := range v {
		headers[i] = v.(string)
	}

	proxy := d.Get("proxy").(bool)
	createOpts := domains.CreateOpts{
		HostName:      d.Get("hostname").(string),
		CertificateId: d.Get("certificate_id").(string),
		Server:        getAllServers(d),
		Proxy:         &proxy,
		SipHeaderName: d.Get("sip_header_name").(string),
		SipHeaderList: headers,
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	domain, err := domains.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Domain: %s", err)
	}

	log.Printf("[DEBUG] Waf domain created: %#v", domain)
	d.SetId(domain.Id)

	if hasFilledOpt(d, "policy_id") {
		var updateHostsOpts policies.UpdateHostsOpts
		policy_id := d.Get("policy_id").(string)
		hosts = append(hosts, d.Id())
		updateHostsOpts.Hosts = hosts
		log.Printf("[DEBUG] Waf policy update Hosts: %#v", hosts)

		_, err = policies.UpdateHosts(wafClient, policy_id, updateHostsOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF Policy Hosts: %s", err)
		}
	}

	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	n, err := domains.Get(wafClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf Domain: %s", err)
	}

	d.SetId(n.Id)
	d.Set("hostname", n.HostName)
	d.Set("certificate_id", n.CertificateId)
	d.Set("proxy", n.Proxy)
	d.Set("sip_header_name", n.SipHeaderName)
	d.Set("sip_header_list", n.SipHeaderList)
	d.Set("access_code", n.AccessCode)
	d.Set("cname", n.Cname)
	d.Set("txt_code", n.TxtCode)
	d.Set("sub_domain", n.SubDomain)
	if n.PolicyID != "" {
		d.Set("policy_id", n.PolicyID)
	}
	d.Set("protect_status", n.ProtectStatus)
	d.Set("access_status", n.AccessStatus)
	d.Set("protocol", n.Protocol)

	servers := make([]map[string]interface{}, len(n.Server))
	for i, server := range n.Server {
		servers[i] = make(map[string]interface{})
		servers[i]["front_protocol"] = server.ClientProtocol
		servers[i]["back_protocol"] = server.ServerProtocol
		servers[i]["address"] = server.Address
		servers[i]["port"] = strconv.Itoa(server.Port)
	}
	d.Set("server", servers)
	return nil
}

func resourceWafDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts domains.UpdateOpts

	if d.HasChange("certificate_id") {
		updateOpts.CertificateId = d.Get("certificate_id").(string)
	}
	if d.HasChange("server") {
		updateOpts.Server = getAllServers(d)
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
		return fmt.Errorf("Error updating OpenTelekomCloud WAF Domain: %s", err)
	}
	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	err = domains.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF Domain: %s", err)
	}

	d.SetId("")
	return nil
}
