package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/domains"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDomainV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDomainV1Create,
		ReadContext:   resourceWafDomainV1Read,
		UpdateContext: resourceWafDomainV1Update,
		DeleteContext: resourceWafDomainV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
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
			"tls": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cipher": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cipher_default", "cipher_1", "cipher_2", "cipher_3",
				}, false),
			},
			"proxy": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"sip_header_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"", "default", "cloudflare", "akamai", "custom",
				}, true),
			},
			"sip_header_list": {
				Type:     schema.TypeList,
				Optional: true,
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
		cProtocol, err := common.FirstOneSet(server, "client_protocol", "front_protocol")
		if err != nil {
			return nil, err
		}
		sProtocol, err := common.FirstOneSet(server, "server_protocol", "back_protocol")
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(server["port"].(string))
		if err != nil {
			return nil, fmt.Errorf("invalid WAF domain server port: %s", err)
		}
		serverOpt := domains.ServerOpts{
			ClientProtocol: cProtocol.(string),
			ServerProtocol: sProtocol.(string),
			Address:        server["address"].(string),
			Port:           port,
		}
		serverOpts = append(serverOpts, serverOpt)
	}

	log.Printf("[DEBUG] getAllServers: %#v", serverOpts)
	return serverOpts, nil
}

func resourceWafDomainV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	sipHeaderList := d.Get("sip_header_list").([]interface{})
	headers := make([]string, len(sipHeaderList))
	for i, header := range sipHeaderList {
		headers[i] = header.(string)
	}

	proxy := d.Get("proxy").(bool)
	servers, err := getAllServers(d)
	if err != nil {
		return fmterr.Errorf("error parsing servers: %w", err)
	}
	createOpts := domains.CreateOpts{
		HostName:      d.Get("hostname").(string),
		CertificateId: d.Get("certificate_id").(string),
		Server:        servers,
		Proxy:         &proxy,
		TLS:           d.Get("tls").(string),
		Cipher:        d.Get("cipher").(string),
		SipHeaderName: d.Get("sip_header_name").(string),
		SipHeaderList: headers,
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	domain, err := domains.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Domain: %w", err)
	}

	d.SetId(domain.Id)

	if p, ok := d.GetOk("policy_id"); ok {
		if err := assignDomainPolicy(client, d.Id(), p.(string)); err != nil {
			return err
		}

		if err := policies.Delete(client, domain.PolicyID).ExtractErr(); err != nil {
			return fmterr.Errorf("error removing automatically created policy: %w", err)
		}
	}

	return resourceWafDomainV1Read(ctx, d, meta)
}

func resourceWafDomainV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}
	n, err := domains.Get(client, d.Id()).Extract()

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud Waf Domain")
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
		d.Set("cipher", n.Cipher),
		d.Set("tls", n.TLS),
		d.Set("policy_id", n.PolicyID),
		d.Set("protect_status", n.ProtectStatus),
		d.Set("access_status", n.AccessStatus),
		d.Set("protocol", n.Protocol),
	)

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
	mErr = multierror.Append(mErr,
		d.Set("server", servers),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting WAF fields: %w", err)
	}

	return nil
}

func resourceWafDomainV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}
	var updateOpts domains.UpdateOpts

	if d.HasChange("certificate_id") {
		updateOpts.CertificateId = d.Get("certificate_id").(string)
	}
	if d.HasChange("server") {
		servers, err := getAllServers(d)
		if err != nil {
			return fmterr.Errorf("error parsing servers: %w", err)
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
		sipHeaderList := d.Get("sip_header_list").([]interface{})
		headers := make([]string, len(sipHeaderList))
		for i, header := range sipHeaderList {
			headers[i] = header.(string)
		}
		updateOpts.SipHeaderList = headers
	}
	if d.HasChange("cipher") {
		cipher := d.Get("cipher").(string)
		updateOpts.Cipher = cipher
	}
	if d.HasChange("tls") {
		tls := d.Get("tls").(string)
		updateOpts.TLS = tls
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = domains.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF Domain: %w", err)
	}

	if d.HasChange("policy_id") {
		if err := assignDomainPolicy(client, d.Id(), d.Get("policy_id").(string)); err != nil {
			return err
		}
	}

	return resourceWafDomainV1Read(ctx, d, meta)
}

func assignDomainPolicy(client *golangsdk.ServiceClient, id, policyID string) diag.Diagnostics {
	if policyID == "" {
		return fmterr.Errorf("can't assign to empty policy")
	}
	opts := policies.UpdateHostsOpts{
		Hosts: []string{id},
	}
	if _, err := policies.UpdateHosts(client, policyID, opts).Extract(); err != nil {
		return fmterr.Errorf("error assigning OpenTelekomCloud WAF Policy to domain: %w", err)
	}
	return nil
}

func resourceWafDomainV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	if err := domains.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Domain: %w", err)
	}

	d.SetId("")
	return nil
}
