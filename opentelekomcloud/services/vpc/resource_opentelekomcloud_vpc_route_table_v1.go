package vpc

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/routetables"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVPCRouteTableV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteTableCreate,
		ReadContext:   resourceVpcRouteTableRead,
		UpdateContext: resourceVpcRouteTableUpdate,
		DeleteContext: resourceVpcRouteTableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(regexp.MustCompile("^[0-9a-zA-Z-_.]*$"),
						"only letters, digits, underscores (_), hyphens (-), and dot (.) are allowed"),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 255),
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"The angle brackets (< and >) are not allowed."),
				),
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"route": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: common.ValidateCIDR,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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

func resourceVpcRouteTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := routetables.CreateOpts{
		VpcID:       d.Get("vpc_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	allRouteOpts := buildRouteTableRoutes(d)
	if len(allRouteOpts) <= MaxCreateRoutes {
		createOpts.Routes = allRouteOpts
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC route table create options: %#v", createOpts)
	routeTable, err := routetables.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud VPC route table: %s", err)
	}

	d.SetId(routeTable.ID)

	if v, ok := d.GetOk("subnets"); ok {
		subnets := common.ExpandToStringSlice(v.(*schema.Set).List())
		err = associateRouteTableSubnets(client, d.Id(), subnets)
		if err != nil {
			return diag.Errorf("error associating subnets with OpenTelekomCloud VPC route table %s: %s", d.Id(), err)
		}
	}

	if len(allRouteOpts) > MaxCreateRoutes {
		updateOpts := routetables.UpdateOpts{
			Routes: &routetables.RouteAction{
				Add: allRouteOpts,
			},
		}

		log.Printf("[DEBUG] add routes to OpenTelekomCloud VPC route table %s: %#v", d.Id(), updateOpts)
		_, err = routetables.Update(client, d.Id(), updateOpts)
		if err != nil {
			return diag.Errorf("error creating OpenTelekomCloud VPC route: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcRouteTableRead(clientCtx, d, meta)
}

func resourceVpcRouteTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	routeTable, err := routetables.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud VPC route table")
	}

	if err := setVpcRouteTableV1Fields(d, routeTable, config.GetRegion(d)); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud VPC route table: %s", err)
	}

	return nil
}

func resourceVpcRouteTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	var changed bool
	var updateOpts routetables.UpdateOpts
	if d.HasChanges("name", "description") {
		changed = true
		desc := d.Get("description").(string)
		updateOpts.Description = &desc
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("route") {
		changed = true
		routesOpts := &routetables.RouteAction{}

		oldR, newR := d.GetChange("route")
		add := newR.(*schema.Set).Difference(oldR.(*schema.Set))
		del := oldR.(*schema.Set).Difference(newR.(*schema.Set))

		if delLen := del.Len(); delLen > 0 {
			delRouteOpts := make([]routetables.DeleteRouteOpts, delLen)
			for i, item := range del.List() {
				opts := item.(map[string]interface{})
				delRouteOpts[i] = routetables.DeleteRouteOpts{
					Type:        opts["type"].(string),
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
				}
			}
			routesOpts.Del = delRouteOpts
		}

		if addLen := add.Len(); addLen > 0 {
			addRouteOpts := make([]routetables.RouteOpts, addLen)
			for i, item := range add.List() {
				opts := item.(map[string]interface{})
				desc := opts["description"].(string)
				addRouteOpts[i] = routetables.RouteOpts{
					Type:        opts["type"].(string),
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
					Description: &desc,
				}
			}
			routesOpts.Add = addRouteOpts
		}
		updateOpts.Routes = routesOpts
	}

	if changed {
		log.Printf("[DEBUG] OpenTelekomCloud VPC route table update options: %#v", updateOpts)
		if _, err := routetables.Update(client, d.Id(), updateOpts); err != nil {
			return diag.Errorf("error updating OpenTelekomCloud VPC route table: %s", err)
		}
	}

	if d.HasChange("subnets") {
		oldS, newS := d.GetChange("subnets")
		associate := newS.(*schema.Set).Difference(oldS.(*schema.Set))
		disassociate := oldS.(*schema.Set).Difference(newS.(*schema.Set))

		disassociateSubnets := common.ExpandToStringSlice(disassociate.List())
		if len(disassociateSubnets) > 0 {
			err = disassociateRouteTableSubnets(client, d.Id(), disassociateSubnets)
			if err != nil {
				return diag.Errorf("error disassociating subnets with OpenTelekomCloud VPC route table %s: %s", d.Id(), err)
			}
		}

		associateSubnets := common.ExpandToStringSlice(associate.List())
		if len(associateSubnets) > 0 {
			err = associateRouteTableSubnets(client, d.Id(), associateSubnets)
			if err != nil {
				return diag.Errorf("error associating subnets with OpenTelekomCloud VPC route table %s: %s", d.Id(), err)
			}
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcRouteTableRead(clientCtx, d, meta)
}

func resourceVpcRouteTableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	if v, ok := d.GetOk("subnets"); ok {
		subnets := common.ExpandToStringSlice(v.(*schema.Set).List())
		err = disassociateRouteTableSubnets(client, d.Id(), subnets)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				d.SetId("")
				return nil
			}
			if _, ok := err.(golangsdk.ErrDefault400); !ok {
				return diag.Errorf("error disassociating subnets with OpenTelekomCloud VPC route table %s: %s", d.Id(), err)
			}
			log.Printf("[WARN] Error disassociating subnets from route table %s (may already be disassociated): %s", d.Id(), err)
		}
	}

	err = routetables.Delete(client, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error deleting OpenTelekomCloud VPC route table: %s", err)
	}

	d.SetId("")
	return nil
}

func associateRouteTableSubnets(client *golangsdk.ServiceClient, id string, subnets []string) error {
	_, err := routetables.Associate(client, id, routetables.AssociateOpts{Subnets: subnets})
	return err
}

func disassociateRouteTableSubnets(client *golangsdk.ServiceClient, id string, subnets []string) error {
	_, err := routetables.Disassociate(client, id, routetables.DisassociateOpts{Subnets: subnets})
	return err
}

func buildRouteTableRoutes(d *schema.ResourceData) []routetables.RouteOpts {
	rawRoutes := d.Get("route").(*schema.Set).List()
	routeOpts := make([]routetables.RouteOpts, len(rawRoutes))

	for i, raw := range rawRoutes {
		opts := raw.(map[string]interface{})
		routeDesc := opts["description"].(string)
		routeOpts[i] = routetables.RouteOpts{
			Type:        opts["type"].(string),
			NextHop:     opts["nexthop"].(string),
			Destination: opts["destination"].(string),
			Description: &routeDesc,
		}
	}

	return routeOpts
}

func expandRouteTableRoutes(routes []routetables.Route) []map[string]interface{} {
	r := make([]map[string]interface{}, 0, len(routes))
	for _, item := range routes {
		// ignore local rule as it can not be modified
		if item.Type == "local" {
			continue
		}
		step := map[string]interface{}{
			"destination": item.DestinationCIDR,
			"type":        item.Type,
			"nexthop":     item.NextHop,
			"description": item.Description,
		}
		r = append(r, step)
	}
	return r
}

func expandRouteTableSubnets(subnets []routetables.Subnet) []string {
	result := make([]string, len(subnets))
	for i, item := range subnets {
		result[i] = item.ID
	}
	return result
}
