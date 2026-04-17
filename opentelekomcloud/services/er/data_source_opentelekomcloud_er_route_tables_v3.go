package er

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/association"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/propagation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/route"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/route_table"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRouteTablesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRouteTablesV3Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     routeTableRelationshipSchemaResource(),
						},
						"propagations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     routeTableRelationshipSchemaResource(),
						},
						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_blackhole": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"attachments": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attachment_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"attachment_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"resource_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func routeTableRelationshipSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func queryRouteTableAssociations(client *golangsdk.ServiceClient, instanceId, routeTableId string) ([]map[string]interface{},
	error) {
	resp, err := association.List(client, association.ListOpts{
		RouterId:     instanceId,
		RouteTableId: routeTableId,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Associations) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(resp.Associations))
	for i, assoc := range resp.Associations {
		result[i] = map[string]interface{}{
			"attachment_id":   assoc.AttachmentID,
			"id":              assoc.ID,
			"attachment_type": assoc.ResourceType,
		}
	}
	return result, nil
}

func queryRouteTablePropagations(client *golangsdk.ServiceClient, instanceId, routeTableId string) ([]map[string]interface{},
	error) {
	resp, err := propagation.List(client, propagation.ListOpts{
		RouterId:     instanceId,
		RouteTableId: routeTableId,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Propagations) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(resp.Propagations))
	for i, prop := range resp.Propagations {
		result[i] = map[string]interface{}{
			"attachment_id":   prop.AttachmentID,
			"id":              prop.ID,
			"attachment_type": prop.ResourceType,
		}
	}
	return result, nil
}

func queryRouteTableRoutes(client *golangsdk.ServiceClient, routeTableId string) ([]map[string]interface{},
	error) {
	resp, err := route.List(client, route.ListOpts{
		RouteTableId: routeTableId,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Routes) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(resp.Routes))
	for i, rt := range resp.Routes {
		rr := map[string]interface{}{
			"destination":  rt.Destination,
			"is_blackhole": rt.IsBlackhole,
			"id":           rt.RouteId,
		}
		if len(rt.NextHops) < 1 {
			result[i] = rr
			continue
		}

		attachments := make([]map[string]interface{}, len(rt.NextHops))
		for i, attachment := range rt.NextHops {
			attachments[i] = map[string]interface{}{
				"attachment_id":   attachment.AttachmentId,
				"attachment_type": attachment.ResourceType,
				"resource_id":     attachment.ResourceId,
			}
		}
		rr["attachments"] = attachments
		result[i] = rr
	}
	return result, nil
}

func filterRouteTablesByTags(d *schema.ResourceData, all []route_table.RouteTable) ([]route_table.RouteTable, error) {
	filter := map[string]interface{}{
		"ID":   d.Get("route_table_id"),
		"Name": d.Get("name"),
	}
	filterResult, err := common.FilterSliceWithField(all, filter)
	if err != nil {
		return nil, fmt.Errorf("error filting security groups list: %s", err)
	}

	tagFilter := d.Get("tags").(map[string]interface{})
	result := make([]route_table.RouteTable, 0, len(filterResult))
	for _, val := range filterResult {
		item := val.(route_table.RouteTable)
		tagmap := common.TagsToMap(item.Tags)

		if !common.MapContains(tagmap, tagFilter) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func flattenRouteTables(client *golangsdk.ServiceClient, instanceId string,
	all []route_table.RouteTable) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, routeTable := range all {
		routeTableId := routeTable.ID

		associationList, _ := queryRouteTableAssociations(client, instanceId, routeTableId)
		propagationList, _ := queryRouteTablePropagations(client, instanceId, routeTableId)
		routeList, _ := queryRouteTableRoutes(client, routeTableId)

		result[i] = map[string]interface{}{
			"id":           routeTableId,
			"name":         routeTable.Name,
			"description":  routeTable.Description,
			"associations": associationList,
			"propagations": propagationList,
			"routes":       routeList,
			"created_at":   routeTable.CreatedAt,
			"updated_at":   routeTable.UpdatedAt,
			"tags":         common.TagsToMap(routeTable.Tags),
		}
	}
	return result
}

func dataSourceRouteTablesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := route_table.List(client, route_table.ListOpts{
		RouterId: instanceId,
	})
	if err != nil {
		return diag.Errorf("error retrieving route tables: %s", err)
	}
	filterResult, err := filterRouteTablesByTags(d, resp.RouteTables)
	if err != nil {
		return diag.Errorf("error retrieving route tables: %s", err)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("route_tables", flattenRouteTables(client, instanceId, filterResult)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving route table list field: %s", mErr)
	}
	return nil
}
