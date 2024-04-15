package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/certificates"
	domains "github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/hosts"
)

const (
	// protectStatusEnable 1: protection status enabled.
	protectStatusEnable = 1
)

func ResourceWafDedicatedDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedDomainV1Create,
		ReadContext:   resourceWafDedicatedDomainV1Read,
		UpdateContext: resourceWafDedicatedDomainV1Update,
		DeleteContext: resourceWafDedicatedDomainV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 80,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
						},
						"server_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
						},
						"address": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"port": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 65535),
							Required:     true,
							ForceNew:     true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keep_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"tls": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TLS v1.0", "TLS v1.1", "TLS v1.2", "TLS v1.3",
				}, false),
			},
			"cipher": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cipher_1",
					"cipher_2",
					"cipher_3",
					"cipher_4",
					"cipher_default",
				}, false),
			},
			"pci_3ds": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"tls", "cipher"},
			},
			"pci_dss": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"tls", "cipher"},
			},
			"timeout_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"send_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"read_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alarm_page": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compliance_certification": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"traffic_identifier": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceWafDedicatedDomainV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	servers := d.Get("server").([]interface{})
	serverOpts := make([]domains.PremiumWafServer, len(servers))
	for i, v := range servers {
		s := v.(map[string]interface{})
		serverOpts[i] = domains.PremiumWafServer{
			FrontProtocol: s["client_protocol"].(string),
			BackProtocol:  s["server_protocol"].(string),
			Address:       s["address"].(string),
			Port:          s["port"].(int),
			Type:          s["type"].(string),
			VpcId:         s["vpc_id"].(string),
		}
	}

	certificateId := d.Get("certificate_id").(string)
	certificateName := ""
	if certificateId != "" {
		certificate, err := certificates.Get(client, certificateId)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud WAF dedicated certificate.")
		}
		certificateName = certificate.Name
	}
	proxy := d.Get("proxy").(bool)
	opts := domains.CreateOpts{
		CertificateId:   d.Get("certificate_id").(string),
		CertificateName: certificateName,
		Hostname:        d.Get("domain").(string),
		Proxy:           &proxy,
		PolicyId:        d.Get("policy_id").(string),
		Server:          serverOpts,
	}

	domain, err := domains.Create(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.ID)

	if d.Get("protect_status").(int) != protectStatusEnable {
		protectStatus := d.Get("protect_status").(int)
		opts := domains.ProtectUpdateOpts{
			ProtectStatus: protectStatus,
		}
		err = domains.UpdateProtectStatus(client, d.Id(), opts)
		if err != nil {
			log.Printf("[ERROR] error change the protection status of OpenTelekomCloud WAF dedicate domain %s: %s", d.Id(), err)
		}
	}

	var timeout bool
	if v, ok := d.GetOk("timeout_config"); ok && len(v.([]interface{})) > 0 {
		timeout = true
	}

	if d.HasChanges("tls", "cipher", "pci_3ds", "pci_dss") || timeout {
		if err := updateWafDedicatedDomain(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedDomainV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedDomainV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	domain, err := domains.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error obtain OpenTelekomCloud WAF dedicated domain information")
	}
	log.Printf("[DEBUG] Get the WAF dedicated domain : %#v", domain)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("domain", domain.Hostname),
		d.Set("server", buildDomainServerAttributes(domain)),
		d.Set("certificate_id", domain.CertificateId),
		d.Set("certificate_name", domain.CertificateName),
		d.Set("policy_id", domain.PolicyId),
		d.Set("proxy", domain.Proxy),
		d.Set("protect_status", domain.ProtectStatus),
		d.Set("access_status", domain.AccessStatus),
		d.Set("protocol", domain.Protocol),
		d.Set("tls", domain.Tls),
		d.Set("cipher", domain.Cipher),
		d.Set("created_at", domain.CreatedAt),
	)

	complianceCertification := make(map[string]interface{})
	if domain.Flag.Pci3ds != "" {
		pci3ds, err := strconv.ParseBool(domain.Flag.Pci3ds)
		if err != nil {
			log.Printf("[WARN] error parse bool pci 3ds, %s", err)
		}
		mErr = multierror.Append(mErr, d.Set("pci_3ds", pci3ds))
		complianceCertification["pci_3ds"] = pci3ds
	}

	if domain.Flag.PciDss != "" {
		pciDss, err := strconv.ParseBool(domain.Flag.PciDss)
		if err != nil {
			log.Printf("[WARN] error parse bool pci dss, %s", err)
		}
		mErr = multierror.Append(mErr, d.Set("pci_dss", pciDss))
		complianceCertification["pci_dss"] = pciDss
	}

	if mErr.ErrorOrNil() != nil {
		return fmterr.Errorf("error setting OpenTelekomCloud WAF domain fields: %s", err)
	}

	if err := d.Set("compliance_certification", complianceCertification); err != nil {
		return diag.FromErr(err)
	}

	trafficMark := map[string]interface{}{}
	if domain.TrafficMark != nil {
		trafficMark["ip_tag"] = strings.Join(domain.TrafficMark.Sip, ",")
		trafficMark["session_tag"] = domain.TrafficMark.Cookie
		trafficMark["user_tag"] = domain.TrafficMark.Params
	}
	if err := d.Set("traffic_identifier", trafficMark); err != nil {
		return diag.FromErr(err)
	}

	alarmPage := map[string]interface{}{}
	if domain.BlockPage != nil {
		alarmPage["template_name"] = domain.BlockPage.Template
		alarmPage["redirect_url"] = domain.BlockPage.RedirectUrl
	}
	if err := d.Set("alarm_page", alarmPage); err != nil {
		return diag.FromErr(err)
	}

	var timeoutConfig []map[string]interface{}

	if domain.TimeoutConfig != nil {
		timeoutRaw := map[string]interface{}{}
		timeoutRaw["read_timeout"] = domain.TimeoutConfig.ReadTimeout
		timeoutRaw["send_timeout"] = domain.TimeoutConfig.SendTimeout
		timeoutRaw["connect_timeout"] = domain.TimeoutConfig.ConnectionTimeout
		timeoutConfig = append(timeoutConfig, timeoutRaw)
	}
	if err := d.Set("timeout_config", timeoutConfig); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedDomainV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	if d.HasChanges("tls", "cipher", "proxy", "certificate_id", "pci_3ds", "pci_dss", "timeout_config") {
		if err := updateWafDedicatedDomain(client, d); err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud WAF dedicated domain: %s", err)
		}
	}

	if d.HasChanges("protect_status") {
		protectStatus := d.Get("protect_status").(int)
		err = domains.UpdateProtectStatus(client, d.Id(), domains.ProtectUpdateOpts{ProtectStatus: protectStatus})
		if err != nil {
			return fmterr.Errorf("[ERROR] error change the protection status of OpenTelekomCloud WAF dedicated domain %s: %s",
				d.Id(), err)
		}
	}

	if d.HasChanges("policy_id") {
		err := updateWafDedicatedPolicy(d, meta)
		if err != nil {
			return err
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedDomainV1Read(clientCtx, d, meta)
}

func updateWafDedicatedPolicy(d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Client: %s", err)
	}

	_, newId := d.GetChange("policy_id")
	policyId := newId.(string)
	updateHostsOpts := policies.UpdateHostsOpts{
		Hosts: []string{d.Id()},
	}
	log.Printf("[DEBUG] Bind OpenTelekomCloud Waf dedicated domain %s to policy %s", d.Id(), policyId)

	_, err = policies.UpdateHosts(client, policyId, updateHostsOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF Policy Hosts: %s", err)
	}
	return nil
}

func resourceWafDedicatedDomainV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	log.Printf("[DEBUG] Delete OpenTelekomCloud WAF dedicated domain (keep_policy: %v).", d.Get("keep_policy"))
	keepPolicy := d.Get("keep_policy").(bool)
	err = domains.Delete(client, d.Id(), domains.DeleteOpts{KeepPolicy: pointerto.Bool(keepPolicy)})
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF dedicated domain: %s", err)
	}

	d.SetId("")
	return nil
}

func updateWafDedicatedDomain(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateOpts := domains.UpdateOpts{
		Tls:    d.Get("tls").(string),
		Cipher: d.Get("cipher").(string),
	}

	if d.HasChange("proxy") && !d.IsNewResource() {
		updateOpts.Proxy = pointerto.Bool(d.Get("proxy").(bool))
	}

	if d.HasChange("certificate_id") && !d.IsNewResource() {
		if v, ok := d.GetOk("certificate_id"); ok {
			certificate, err := certificates.Get(client, v.(string))
			if err != nil {
				return fmt.Errorf("error retrieving OpenTelekomCloud WAF dedicated certificate: %s", err)
			}
			updateOpts.CertificateName = certificate.Name
			updateOpts.CertificateId = certificate.ID
		}
	}

	if v, ok := d.GetOk("timeout_config"); ok && len(v.([]interface{})) > 0 {
		rawArray := v.([]interface{})[0].(map[string]interface{})
		updateOpts = domains.UpdateOpts{
			TimeoutConfig: &domains.TimeoutConfigObject{
				ConnectionTimeout: rawArray["connect_timeout"].(int),
				SendTimeout:       rawArray["send_timeout"].(int),
				ReadTimeout:       rawArray["read_timeout"].(int),
			},
		}
	}

	if d.HasChanges("pci_3ds", "pci_dss") {
		flag, err := getHostFlag(d)
		if err != nil {
			return err
		}
		updateOpts.Flag = flag
	}

	log.Printf("[DEBUG] OpenTelekomCloud Waf dedicated domain update: %#v", updateOpts)
	_, err := domains.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmt.Errorf("error updating OpenTelekomCloud WAF dedicated domain: %s", err)
	}
	return nil
}

func getHostFlag(d *schema.ResourceData) (*domains.FlagObject, error) {
	pci3ds := d.Get("pci_3ds").(bool)
	pciDss := d.Get("pci_dss").(bool)
	if !pci3ds && !pciDss {
		return nil, nil
	}

	if d.Get("tls").(string) != "TLS v1.2" || d.Get("cipher").(string) != "cipher_2" {
		return nil, fmt.Errorf("pci_3ds and pci_dss must be used together with tls and cipher. " +
			"Tls must be set to TLS v1.2, and cipher must be set to cipher_2")
	}
	return &domains.FlagObject{
		Pci3ds: strconv.FormatBool(pci3ds),
		PciDss: strconv.FormatBool(pciDss),
	}, nil
}

func buildDomainServerAttributes(domain *domains.Host) []map[string]interface{} {
	servers := make([]map[string]interface{}, 0, len(domain.Server))
	for _, s := range domain.Server {
		servers = append(servers, map[string]interface{}{
			"client_protocol": s.FrontProtocol,
			"server_protocol": s.BackProtocol,
			"address":         s.Address,
			"port":            s.Port,
			"type":            s.Type,
			"vpc_id":          s.VpcId,
		})
	}
	return servers
}
