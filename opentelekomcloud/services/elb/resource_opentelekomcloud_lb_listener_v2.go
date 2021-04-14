package elb

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/listeners"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceListenerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerV2Create,
		Read:   resourceListenerV2Read,
		Update: resourceListenerV2Update,
		Delete: resourceListenerV2Delete,

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
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "TERMINATED_HTTPS"}, false),
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true, // could be updated due to docs, but gopher doesn't define it
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			/*"connection_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			}, */
			// new feature 2020 to support https2
			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// new feature 2020 to handle Client certificates
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sni_container_refs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// new feature 2020 to give a choice of the http standard on https termination
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tls-1-0", "tls-1-1", "tls-1-2", "tls-1-2-strict"}, false),
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceListenerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	http2Enable := d.Get("http2_enable").(bool) // would prefer a fix in the gopher...
	adminStateUp := d.Get("admin_state_up").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.(*schema.Set).List() {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}
	createOpts := listeners.CreateOpts{
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		TenantID:               d.Get("tenant_id").(string),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		Http2Enable:            &http2Enable,
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		CAContainerRef:         d.Get("client_ca_tls_container_ref").(string),
		SniContainerRefs:       sniContainerRefs,
		TlsCiphersPolicy:       d.Get("tls_ciphers_policy").(string),
		AdminStateUp:           &adminStateUp,
	}

	/*if v, ok := d.GetOk("connection_limit"); ok {
		connectionLimit := v.(int)
		createOpts.ConnLimit = &connectionLimit
	} */

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	lbID := createOpts.LoadbalancerID
	timeout := d.Timeout(schema.TimeoutCreate)
	if err := waitForLBV2LoadBalancer(client, lbID, "ACTIVE", nil, timeout); err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create listener")
	var listener *listeners.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = listeners.Create(client, createOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error creating listener: %s", err)
	}

	// Wait for LoadBalancer to become active again before continuing
	if err := waitForLBV2LoadBalancer(client, lbID, "ACTIVE", nil, timeout); err != nil {
		return err
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "listeners", listener.ID, tagList).ExtractErr(); err != nil {
			return fmt.Errorf("error setting tags of LoadBalancer: %s", err)
		}
	}

	d.SetId(listener.ID)

	return resourceListenerV2Read(d, meta)
}

func resourceListenerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	listener, err := listeners.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "listener")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("protocol", listener.Protocol),
		d.Set("protocol_port", listener.ProtocolPort),
		d.Set("tenant_id", listener.TenantID),
		d.Set("name", listener.Name),
		d.Set("default_pool_id", listener.DefaultPoolID),
		d.Set("description", listener.Description),
		d.Set("http2_enable", listener.Http2Enable),
		d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef),
		d.Set("client_ca_tls_container_ref", listener.CAContainerRef),
		// d.Set("connection_limit", listener.ConnLimit),
		d.Set("sni_container_refs", listener.SniContainerRefs),
		d.Set("tls_ciphers_policy", listener.TlsCiphersPolicy),
		d.Set("admin_state_up", listener.AdminStateUp),
	)

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	// save tags
	resourceTags, err := tags.Get(client, "listeners", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error fetching OpenTelekomCloud LB Listener tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmt.Errorf("error saving tags for OpenTelekomCloud LB Listener: %s", err)
	}

	return nil
}

func resourceListenerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	var updateOpts listeners.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	/*if d.HasChange("connection_limit") {
		connLimit := d.Get("connection_limit").(int)
		updateOpts.ConnLimit = &connLimit
	} */
	if d.HasChange("http2_enable") {
		http2Enable := d.Get("http2_enable").(bool)
		updateOpts.Http2Enable = &http2Enable
	}
	if d.HasChange("default_tls_container_ref") {
		updateOpts.DefaultTlsContainerRef = d.Get("default_tls_container_ref").(string)
	}
	if d.HasChange("client_ca_tls_container_ref") {
		updateOpts.CAContainerRef = d.Get("client_ca_tls_container_ref").(string)
	}
	if d.HasChange("sni_container_refs") {
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_container_refs"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		updateOpts.SniContainerRefs = sniContainerRefs
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}
	if d.HasChange("tls_ciphers_policy") {
		updateOpts.TlsCiphersPolicy = d.Get("tls_ciphers_policy").(string)
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutUpdate)
	if err := waitForLBV2LoadBalancer(client, lbID, "ACTIVE", nil, timeout); err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = listeners.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error updating listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	if err := waitForLBV2LoadBalancer(client, lbID, "ACTIVE", nil, timeout); err != nil {
		return err
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "listeners", d.Id()); err != nil {
			return fmt.Errorf("error updating tags of LoadBalancer Listener %s: %s", d.Id(), err)
		}
	}

	return resourceListenerV2Read(d, meta)

}

func resourceListenerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	if err := waitForLBV2LoadBalancer(client, lbID, "ACTIVE", nil, timeout); err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting listener %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = listeners.Delete(client, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	if err := waitForLBV2LoadBalancer(client, lbID, "ACTIVE", nil, timeout); err != nil {
		return err
	}

	// Wait for Listener to delete
	if err := waitForLBV2Listener(client, d.Id(), "DELETED", nil, timeout); err != nil {
		return err
	}

	return nil
}
