package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elbaas/listeners"
)

var ProtocolFormats = []string{"HTTP", "TCP", "HTTPS", "SSL", "UDP"}

func ValidateProtocolFormat(v interface{}, k string) (ws []string, errors []error) {
	return ValidateStringList(v, k, ProtocolFormats)
}

func resourceEListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceEListenerCreate,
		Read:   resourceEListenerRead,
		Update: resourceEListenerUpdate,
		Delete: resourceEListenerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateProtocolFormat,
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"backend_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateStringList(v, k, []string{"HTTP", "TCP", "UDP"})
				},
			},
			"backend_port": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateIntRange(v, k, 1, 65535)
				},
			},

			"lb_algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateStringList(v, k, []string{"roundrobin", "leastconn", "source"})
				},
			},

			"session_sticky": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"session_sticky_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"cookie_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateIntRange(v, k, 1, 1440)
				},
			},

			"tcp_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateIntRange(v, k, 1, 5)
				},
			},

			"tcp_draining": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tcp_draining_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateIntRange(v, k, 0, 60)
				},
			},

			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},

			"udp_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateIntRange(v, k, 1, 1440)
				},
			},

			"ssl_protocols": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateStringList(v, k, []string{"TLSv1.2", "TLSv1.2 TLSv1.1 TLSv1"})
				},
			},

			"ssl_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateStringList(v, k, []string{"Default", "Extended", "Strict"})
				},
			},
		},
	}
}

func resourceEListenerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var certificates []string
	if raw, ok := d.GetOk("certificates"); ok {
		for _, v := range raw.([]interface{}) {
			certificates = append(certificates, v.(string))
		}
	}
	createOpts := listeners.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		LoadbalancerID:      d.Get("loadbalancer_id").(string),
		Protocol:            listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:        d.Get("protocol_port").(int),
		BackendProtocol:     listeners.Protocol(d.Get("backend_protocol").(string)),
		BackendProtocolPort: d.Get("backend_port").(int),
		Algorithm:           d.Get("lb_algorithm").(string),
		SessionSticky:       d.Get("session_sticky").(bool),
		StickySessionType:   d.Get("session_sticky_type").(string),
		CookieTimeout:       d.Get("cookie_timeout").(int),
		TcpTimeout:          d.Get("tcp_timeout").(int),
		TcpDraining:         d.Get("tcp_draining").(bool),
		TcpDrainingTimeout:  d.Get("tcp_draining_timeout").(int),
		CertificateID:       d.Get("certificate_id").(string),
		Certificates:        certificates,
		UDPTimeout:          d.Get("udp_timeout").(int),
		SSLProtocols:        d.Get("ssl_protocols").(string),
		SSLCiphers:          d.Get("ssl_ciphers").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	listener, err := listeners.Create(client, createOpts).Extract()
	if err != nil {
		return err
	}
	d.SetId(listener.ID)

	log.Printf("[DEBUG] Successfully created listener %s", listener.ID)

	return resourceEListenerRead(d, meta)
}

func resourceEListenerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	listener, err := listeners.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "listener")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	d.Set("backend_port", listener.BackendProtocolPort)
	d.Set("backend_protocol", listener.BackendProtocol)
	d.Set("session_sticky_type", listener.StickySessionType)
	d.Set("description", listener.Description)
	d.Set("load_balancer_id", listener.LoadbalancerID)
	d.Set("protocol", listener.Protocol)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("cookie_timeout", listener.CookieTimeout)
	d.Set("admin_state_up", listener.AdminStateUp)
	d.Set("session_sticky", listener.SessionSticky)
	d.Set("lb_algorithm", listener.Algorithm)
	d.Set("name", listener.Name)
	d.Set("certificate_id", listener.CertificateID)
	d.Set("certificates", listener.Certificates)
	d.Set("tcp_timeout", listener.TcpTimeout)
	d.Set("udp_timeout", listener.UDPTimeout)
	d.Set("ssl_protocols", listener.SSLProtocols)
	d.Set("ssl_ciphers", listener.SSLCiphers)
	d.Set("tcp_draining", listener.TcpDraining)
	d.Set("tcp_draining_timeout", listener.TcpDrainingTimeout)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceEListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts listeners.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("protocol_port") {
		updateOpts.ProtocolPort = d.Get("protocol_port").(int)
	}
	if d.HasChange("backend_port") {
		updateOpts.BackendProtocolPort = d.Get("backend_port").(int)
	}
	if d.HasChange("lb_algorithm") {
		updateOpts.Algorithm = d.Get("lb_algorithm").(string)
	}
	if d.HasChange("tcp_timeout") {
		updateOpts.TcpTimeout = d.Get("tcp_timeout").(int)
	}
	if d.HasChange("tcp_draining") {
		updateOpts.TcpDraining = d.Get("tcp_draining").(bool)
	}
	if d.HasChange("tcp_draining_timeout") {
		updateOpts.TcpDrainingTimeout = d.Get("tcp_draining_timeout").(int)
	}
	if d.HasChange("udp_timeout") {
		updateOpts.UDPTimeout = d.Get("udp_timeout").(int)
	}
	if d.HasChange("ssl_protocols") {
		updateOpts.SSLProtocols = d.Get("ssl_protocols").(string)
	}
	if d.HasChange("ssl_ciphers") {
		updateOpts.SSLCiphers = d.Get("ssl_ciphers").(string)
	}

	_, err = listeners.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return err
	}

	return resourceEListenerRead(d, meta)

}

func resourceEListenerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting listener %s", id)

	err = listeners.Delete(client, id).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}
