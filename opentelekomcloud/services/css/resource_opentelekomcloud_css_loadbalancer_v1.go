package css

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/css/v1/load_balancer"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const timeoutSeconds = 300

// ResourceCssLoadBalancerV1 defines the Terraform resource for managing
// CSS Load Balancer configuration.
func ResourceCssLoadBalancerV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssLoadBalancerEnable,
		ReadContext:   resourceCssLoadBalancerRead,
		UpdateContext: resourceCssLoadBalancerUpdate,
		DeleteContext: resourceCssLoadBalancerDisable,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCssLoadBalancerV1Import,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"elb_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"server_cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ca_cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceCssLoadBalancerEnable(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CssV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CSS v1 client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	elbID := d.Get("elb_id").(string)
	agency := d.Get("agency").(string)

	d.SetId(clusterID)

	if _, err := load_balancer.EnableLoadBalancer(client, clusterID, load_balancer.EnableLoadBalancerOpts{
		ElbId:  elbID,
		Agency: agency,
	}); err != nil {
		return diag.Errorf("failed to enable load balancer: %s", err)
	}

	if opts := extractListenerOpts(d); opts != nil {
		if _, err := load_balancer.ConfigureListener(client, clusterID, *opts); err != nil {
			return diag.Errorf("failed to configure listener: %s", err)
		}
		if err := load_balancer.WaitForListenerStatus(client, clusterID, timeoutSeconds); err != nil {
			return diag.Errorf("error waiting for listener to be active: %s", err)
		}
	}

	return resourceCssLoadBalancerRead(ctx, d, meta)
}

func resourceCssLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CssV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CSS v1 client: %s", err)
	}

	details, err := load_balancer.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("failed to get load balancer details: %s", err)
	}

	listenerList := make([]interface{}, 0)
	if details.Listener.Id != "" {
		listenerList = append(listenerList, map[string]interface{}{
			"protocol":       details.Listener.Protocol,
			"protocol_port":  details.Listener.ProtocolPort,
			"server_cert_id": details.ServerCertId,
			"ca_cert_id":     details.CacertId,
		})
	}

	mErr := multierror.Append(
		d.Set("elb_id", details.LoadBalancer.Id),
		d.Set("agency", details.Agency),
		d.Set("cluster_id", d.Id()),
		d.Set("listener", listenerList),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCssLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CssV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CSS v1 client: %s", err)
	}

	if d.HasChange("listener") || d.HasChange("elb_id") {
		if err := updateListener(d, client); err != nil {
			return diag.Errorf("failed to update listener: %s", err)
		}
	}

	return resourceCssLoadBalancerRead(ctx, d, meta)
}

func resourceCssLoadBalancerDisable(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CssV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CSS v1 client: %s", err)
	}

	clusterID := d.Id()

	if _, err := load_balancer.DisableLoadBalancer(client, clusterID); err != nil {
		return diag.Errorf("failed to disable load balancer: %s", err)
	}

	d.SetId("")
	return nil
}

func extractListenerOpts(d *schema.ResourceData) *load_balancer.ConfigureListenerOpts {
	list := d.Get("listener").([]interface{})
	if len(list) == 0 {
		return nil
	}

	m := list[0].(map[string]interface{})
	opts := load_balancer.ConfigureListenerOpts{
		Protocol:     m["protocol"].(string),
		ProtocolPort: m["protocol_port"].(int),
	}

	if strings.ToUpper(opts.Protocol) == "HTTPS" {
		if v, ok := m["server_cert_id"].(string); ok && v != "" {
			opts.ServerCertId = v
		}
		if v, ok := m["ca_cert_id"].(string); ok && v != "" {
			opts.CaCertId = v
		}
		if v, ok := m["type"].(string); ok && v != "" {
			opts.Type = v
		}
	}

	return &opts
}

func updateListener(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	clusterID := d.Id()
	elbID := d.Get("elb_id").(string)
	agency := d.Get("agency").(string)

	oldList, newList := d.GetChange("listener")
	oldLen := len(oldList.([]interface{}))
	newLen := len(newList.([]interface{}))

	switch {
	case oldLen > 0 && newLen == 0:
		// Listener removed: disable and re-enable ELB
		if _, err := load_balancer.DisableLoadBalancer(client, clusterID); err != nil {
			return err
		}
		_, err := load_balancer.EnableLoadBalancer(client, clusterID, load_balancer.EnableLoadBalancerOpts{
			ElbId:  elbID,
			Agency: agency,
		})
		return err

	case newLen > 0:
		opts := extractListenerOpts(d)
		if opts != nil {
			if _, err := load_balancer.ConfigureListener(client, clusterID, *opts); err != nil {
				return err
			}
			if err := load_balancer.WaitForListenerStatus(client, clusterID, timeoutSeconds); err != nil {
				return err
			}
		}
		if d.HasChange("elb_id") {
			if _, err := load_balancer.EnableLoadBalancer(client, clusterID, load_balancer.EnableLoadBalancerOpts{
				ElbId:  elbID,
				Agency: agency,
			}); err != nil {
				return err
			}
		}
		return nil

	default:
		// Only ELB ID changed
		if d.HasChange("elb_id") {
			_, err := load_balancer.EnableLoadBalancer(client, clusterID, load_balancer.EnableLoadBalancerOpts{
				ElbId:  elbID,
				Agency: agency,
			})
			return err
		}
		return nil
	}
}

func resourceCssLoadBalancerV1Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if diagRead := resourceCssLoadBalancerRead(ctx, d, meta); diagRead.HasError() {
		return nil, fmt.Errorf("error reading opentelekomcloud_css_loadbalancer_v1 %s: %s", d.Id(), diagRead[0].Summary)
	}

	config := meta.(*cfg.Config)
	client, err := config.CssV1Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	clusterID := d.Id()

	details, err := load_balancer.Get(client, clusterID)
	if err != nil {
		return nil, err
	}

	// If LB is disabled, import not possible
	if !details.Enabled {
		return nil, fmt.Errorf("load balancer is disabled for cluster %s; cannot import resource", clusterID)
	}

	listener := make([]interface{}, 0)
	if details.Listener.Id != "" {
		listener = append(listener, map[string]interface{}{
			"protocol":       details.Listener.Protocol,
			"protocol_port":  details.Listener.ProtocolPort,
			"server_cert_id": details.ServerCertId,
			"ca_cert_id":     details.CacertId,
		})
	}

	mErr := multierror.Append(nil,
		d.Set("elb_id", details.LoadBalancer.Id),
		d.Set("agency", details.Agency),
		d.Set("cluster_id", clusterID),
		d.Set("listener", listener),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("error setting addon attributes: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
