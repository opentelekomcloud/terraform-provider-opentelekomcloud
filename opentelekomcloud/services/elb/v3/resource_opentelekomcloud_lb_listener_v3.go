package v3

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/listeners"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceListenerV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerV3Create,
		ReadContext:   resourceListenerV3Read,
		UpdateContext: resourceListenerV3Update,
		DeleteContext: resourceListenerV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"admin_state_up": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				ValidateFunc: common.ValidateTrueOnly,
			},
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "HTTP", "UDP", "HTTPS",
				}, false),
			},
			"protocol_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsPortNumber,
			},
			"sni_container_refs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tls-1-0", "tls-1-1", "tls-1-2", "tls-1-2-strict",
				}, false),
			},
			"member_retry_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keep_alive_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 4000),
			},
			"client_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300),
			},
			"member_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300),
			},
			"insert_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forward_elb_ip": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"forwarded_port": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"forwarded_for_port": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"forwarded_host": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateTags,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getInsertHeaders(d *schema.ResourceData) *listeners.InsertHeaders {
	if d.Get("insert_headers.#").(int) == 0 {
		return nil
	}
	insertHeaders := d.Get("insert_headers.0").(map[string]interface{})
	forwardELBIp := insertHeaders["forward_elb_ip"].(bool)
	forwardedPort := insertHeaders["forwarded_port"].(bool)
	forwardedForPort := insertHeaders["forwarded_for_port"].(bool)
	forwardedHost := insertHeaders["forwarded_host"].(bool)
	return &listeners.InsertHeaders{
		ForwardedELBIP:   &forwardELBIp,
		ForwardedPort:    &forwardedPort,
		ForwardedForPort: &forwardedForPort,
		ForwardedHost:    &forwardedHost,
	}
}

func resourceListenerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	protocol := listeners.Protocol(d.Get("protocol").(string))
	opts := listeners.CreateOpts{
		AdminStateUp:           &adminStateUp,
		CAContainerRef:         d.Get("client_ca_tls_container_ref").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		Description:            d.Get("description").(string),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		SniContainerRefs:       common.ExpandToStringSlice(d.Get("sni_container_refs").(*schema.Set).List()),
		Tags:                   common.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		TlsCiphersPolicy:       d.Get("tls_ciphers_policy").(string),
		KeepAliveTimeout:       d.Get("keep_alive_timeout").(int),
		ClientTimeout:          d.Get("client_timeout").(int),
		MemberTimeout:          d.Get("member_timeout").(int),
		InsertHeaders:          getInsertHeaders(d),
	}
	switch protocol {
	case listeners.ProtocolHTTPS:
		http2Enable := d.Get("http2_enable").(bool)
		opts.Http2Enable = &http2Enable
		fallthrough
	case listeners.ProtocolHTTP:
		memberRetryEnable := d.Get("member_retry_enable").(bool)
		opts.EnableMemberRetry = &memberRetryEnable
	}

	log.Printf("[DEBUG] Create Options: %#v", opts)
	lb, err := listeners.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating LoadBalancerV3: %w", err)
	}

	d.SetId(lb.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceListenerV3Read(clientCtx, d, meta)
}

func resourceListenerV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	listener, err := listeners.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "listenerV3")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	return setLBListenerFields(d, listener)
}

func setLBListenerFields(d *schema.ResourceData, listener *listeners.Listener) diag.Diagnostics {
	insertHeaders := []map[string]interface{}{
		{
			"forward_elb_ip":     listener.InsertHeaders.ForwardedELBIP,
			"forwarded_port":     listener.InsertHeaders.ForwardedPort,
			"forwarded_for_port": listener.InsertHeaders.ForwardedForPort,
			"forwarded_host":     listener.InsertHeaders.ForwardedHost,
		},
	}
	mErr := multierror.Append(
		d.Set("admin_state_up", listener.AdminStateUp),
		d.Set("client_ca_tls_container_ref", listener.CAContainerRef),
		d.Set("default_pool_id", listener.DefaultPoolID),
		d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef),
		d.Set("description", listener.Description),
		d.Set("http2_enable", listener.Http2Enable),
		d.Set("name", listener.Name),
		d.Set("protocol", listener.Protocol),
		d.Set("insert_headers", insertHeaders),
		d.Set("protocol_port", listener.ProtocolPort),
		d.Set("sni_container_refs", listener.SniContainerRefs),
		d.Set("tls_ciphers_policy", listener.TlsCiphersPolicy),
		d.Set("keep_alive_timeout", listener.KeepAliveTimeout),
		d.Set("client_timeout", listener.ClientTimeout),
		d.Set("member_timeout", listener.MemberTimeout),
		d.Set("loadbalancer_id", listener.Loadbalancers[0].ID),
		d.Set("created_at", listener.CreatedAt),
		d.Set("updated_at", listener.UpdatedAt),
		d.Set("tags", common.TagsToMap(listener.Tags)),
	)

	switch listeners.Protocol(listener.Protocol) {
	case listeners.ProtocolHTTPS, listeners.ProtocolHTTP:
		mErr = multierror.Append(mErr,
			d.Set("member_retry_enable", listener.EnableMemberRetry),
		)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceListenerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var updateOpts listeners.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}
	if d.HasChange("client_ca_tls_container_ref") {
		clientCaTlsContainer := d.Get("client_ca_tls_container_ref").(string)
		updateOpts.CAContainerRef = &clientCaTlsContainer
	}
	if d.HasChange("default_pool_id") {
		updateOpts.DefaultPoolID = d.Get("default_pool_id").(string)
	}
	if d.HasChange("default_tls_container_ref") {
		defaultTlsContainerRef := d.Get("default_tls_container_ref").(string)
		updateOpts.DefaultTlsContainerRef = &defaultTlsContainerRef
	}
	if d.HasChange("http2_enable") {
		http2Enable := d.Get("http2_enable").(bool)
		updateOpts.Http2Enable = &http2Enable
	}
	if d.HasChange("tls_ciphers_policy") {
		tlsCiphersPolicy := d.Get("tls_ciphers_policy").(string)
		updateOpts.TlsCiphersPolicy = &tlsCiphersPolicy
	}
	if d.HasChange("sni_container_refs") {
		sniContainerRefs := common.ExpandToStringSlice(d.Get("sni_container_refs").(*schema.Set).List())
		updateOpts.SniContainerRefs = &sniContainerRefs
	}
	if d.HasChange("insert_headers") {
		updateOpts.InsertHeaders = getInsertHeaders(d)
	}
	if d.HasChange("member_retry_enable") {
		memberRetryEnable := d.Get("member_retry_enable").(bool)
		updateOpts.EnableMemberRetry = &memberRetryEnable
	}
	if d.HasChange("keep_alive_timeout") {
		updateOpts.KeepAliveTimeout = d.Get("keep_alive_timeout").(int)
	}
	if d.HasChange("client_timeout") {
		updateOpts.ClientTimeout = d.Get("client_timeout").(int)
	}
	if d.HasChange("member_timeout") {
		updateOpts.MemberTimeout = d.Get("member_timeout").(int)
	}

	log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)
	_, err = listeners.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("unable to update ListenerV3 %s: %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceListenerV3Read(clientCtx, d, meta)
}

func resourceListenerV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting listener: %s", d.Id())
	if err := listeners.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("unable to delete ListenerV3 %s: %s", d.Id(), err)
	}

	return nil
}
