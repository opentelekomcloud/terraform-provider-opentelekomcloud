package vpcep

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpcep/v1/services"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVPCEPServiceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEPServiceCreate,
		ReadContext:   resourceVPCEPServiceRead,
		UpdateContext: resourceVPCEPServiceUpdate,
		DeleteContext: resourceVPCEPServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"port_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vip_port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 16),
					validation.StringMatch(
						regexp.MustCompile(`\w-`),
						"The value contains a maximum of 16 characters, including letters, digits, underscores (_), and hyphens (-).",
					),
				),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"approval_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"interface", "gateway"},
					true,
				),
			},
			"server_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"VM", "VIP", "LB"}, true,
				),
				DiffSuppressFunc: common.SuppressCaseInsensitive,
			},
			"port": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 0xffff),
						},
						"server_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 0xffff),
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"TCP", "UDP"}, false,
							),
							Default: "TCP",
						},
					},
				},
			},
			"tcp_proxy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"close", "toa_open", "proxy_open", "open"}, true,
				),
				DiffSuppressFunc: common.SuppressCaseInsensitive,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

const keyClient = "vpcep-client"

func resourceVPCEPServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	client, err := config.VpcEpV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrClientCreate, err)
	}

	approvalEnabled := d.Get("approval_enabled").(bool)
	opts := &services.CreateOpts{
		PortID:          d.Get("port_id").(string),
		PoolID:          d.Get("pool_id").(string),
		VIPPortID:       d.Get("vip_port_id").(string),
		ServiceName:     d.Get("name").(string),
		RouterID:        d.Get("vpc_id").(string),
		ApprovalEnabled: &approvalEnabled,
		ServiceType:     services.ServiceType(d.Get("service_type").(string)),
		ServerType:      services.ServerType(d.Get("server_type").(string)),
		Ports:           getPorts(d),
		TCPProxy:        d.Get("tcp_proxy").(string),
		Tags:            common.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	svc, err := services.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating VPC EP service: %w", err)
	}

	d.SetId(svc.ID)

	err = services.WaitForServiceStatus(
		client, d.Id(), services.StatusAvailable,
		timeoutSeconds(d, schema.TimeoutCreate),
	)
	if err != nil {
		return fmterr.Errorf("error waiting for VPC EP service to become available: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceVPCEPServiceRead(clientCtx, d, meta)
}

func resourceVPCEPServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.VpcEpV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrClientCreate, err)
	}

	svc, err := services.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("error reading VPC EP service: %w", err)
	}

	mErr := multierror.Append(
		d.Set("port_id", svc.PortID),
		d.Set("pool_id", svc.PoolID),
		d.Set("vip_port_id", svc.VIPPortID),
		d.Set("name", onlyServiceName(svc.ServiceName)),
		d.Set("vpc_id", svc.RouterID),
		d.Set("approval_enabled", svc.ApprovalEnabled),
		d.Set("service_type", svc.ServiceType),
		d.Set("server_type", svc.ServerType),
		d.Set("port", portsSlice(svc.Ports)),
		d.Set("tags", common.TagsToMap(svc.Tags)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting VPC EP service attributes: %w", err)
	}

	return nil
}

func resourceVPCEPServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	client, err := config.VpcEpV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrClientCreate, err)
	}

	opts := &services.UpdateOpts{}
	if d.HasChange("port_id") {
		opts.PortID = d.Get("port_id").(string)
	}
	if d.HasChange("vip_port_id") {
		opts.VIPPortID = d.Get("vip_port_id").(string)
	}
	if d.HasChange("name") {
		opts.ServiceName = d.Get("name").(string)
	}
	if d.HasChange("port") {
		opts.Ports = getPorts(d)
	}
	if d.HasChange("approval_enabled") {
		enabled := d.Get("approval_enabled").(bool)
		opts.ApprovalEnabled = &enabled
	}

	_, err = services.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating VPC EP service: %w", err)
	}

	err = services.WaitForServiceStatus(
		client, d.Id(), services.StatusAvailable,
		timeoutSeconds(d, schema.TimeoutUpdate),
	)
	if err != nil {
		return fmterr.Errorf("error waiting for VPC EP service to become available: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceVPCEPServiceRead(clientCtx, d, meta)
}

func resourceVPCEPServiceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	client, err := config.VpcEpV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrClientCreate, err)
	}

	err = services.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil
		}
		return fmterr.Errorf("error reading VPC EP service: %w", err)
	}

	err = services.WaitForServiceStatus(
		client, d.Id(), services.StatusDeleted,
		timeoutSeconds(d, schema.TimeoutDelete),
	)
	if err != nil {
		return fmterr.Errorf("error waiting for VPC EP service to become deleted: %w", err)
	}

	return nil
}

func timeoutSeconds(d *schema.ResourceData, key string) int {
	t := d.Timeout(key)
	return int(t.Seconds())
}

var svcNameRe = regexp.MustCompile(`[\w-]+\.([\w-]+)\.[\w-]+`)

// Get serviceName from regionName.serviceName.serviceId
func onlyServiceName(in string) string {
	matches := svcNameRe.FindStringSubmatch(in)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func getPorts(d *schema.ResourceData) []services.PortMapping {
	portSet := d.Get("port").(*schema.Set)
	pMapping := make([]services.PortMapping, portSet.Len())
	for i, p := range portSet.List() {
		port := p.(map[string]interface{})
		pMapping[i] = services.PortMapping{
			ClientPort: port["client_port"].(int),
			ServerPort: port["server_port"].(int),
			Protocol:   port["protocol"].(string),
		}
	}
	return pMapping
}

func portsSlice(pts []services.PortMapping) []interface{} {
	ports := make([]interface{}, len(pts))
	for i, p := range pts {
		ports[i] = map[string]interface{}{
			"client_port": p.ClientPort,
			"server_port": p.ServerPort,
			"protocol":    p.Protocol,
		}
	}
	return ports
}
