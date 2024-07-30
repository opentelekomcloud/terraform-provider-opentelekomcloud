package migrations

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/dataarts/v1.1/cluster"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

// TODO: add to the documentation a warning: "All changes in parameters will cause the cluster recreation."
func ResourceClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterV1Create,
		ReadContext:   resourceClusterV1Read,
		UpdateContext: resourceClusterV1Update,
		DeleteContext: resourceClusterV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_remind": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"phone_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"language": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"schedule_boot_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_schedule_boot_off": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				Elem:     instanceSchema(),
			},
			"datastore": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(ClusterTypeCDM),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"extended_properties": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"trial": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"schedule_off_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sys_tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"is_auto_off": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_endpoint_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"cluster_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func instanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type: schema.TypeString,
				// Optional: true,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(FlavorTypeSmall),
					string(FlavorTypeMedium),
					string(FlavorTypeLarge),
					string(FlavorTypeXLarge),
				}, false),
			},
			"flavor": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rel": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"href": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ClusterTypeCDM),
				}, false),
			},
			"nics": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"net": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"volume": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shard_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_fix_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func buildCluster(d *schema.ResourceData) apis.Cluster {

	i := d.Get("instances").(*schema.Set)
	ds := d.Get("datastore").(*schema.Set)
	ep := d.Get("extended_properties").(*schema.Set)
	st := d.Get("sys_tags").([]interface{})

	return apis.Cluster{
		ScheduleBootTime:   d.Get("schedule_boot_time").(string),
		IsScheduleBootOff:  pointerto.Bool(d.Get("is_schedule_boot_off").(bool)),
		Instances:          buildInstances(i),
		DataStore:          buildDatastore(ds),
		ExtendedProperties: buildExtendedProperties(ep),
		ScheduleOffTime:    d.Get("schedule_off_time").(string),
		VpcId:              d.Get("vpc_id").(string),
		Name:               d.Get("name").(string),
		SysTags:            buildSysTags(st),
		IsAutoOff:          d.Get("is_auto_off").(bool),
	}
}

func buildInstances(instances *schema.Set) []apis.Instance {
	if instances.Len() == 0 {
		return nil
	}
	insSlice := make([]apis.Instance, instances.Len())
	for index, instance := range instances.List() {
		i := instance.(map[string]interface{})
		insSlice[index] = apis.Instance{
			AZ:        i["availability_zone"].(string),
			Nics:      buildNics(i["nics"].(*schema.Set)),
			FlavorRef: i["flavor_id"].(string),
			Type:      i["type"].(string),
		}
	}
	return insSlice
}

func buildNics(nics *schema.Set) []apis.Nic {
	if nics.Len() == 0 {
		return nil
	}
	nicSlice := make([]apis.Nic, nics.Len())

	for index, nic := range nics.List() {
		n := nic.(map[string]interface{})
		nicSlice[index] = apis.Nic{
			SecurityGroupId: n["security_group"].(string),
			NetId:           n["net"].(string),
		}
	}

	return nicSlice
}

func buildDatastore(ds *schema.Set) *apis.Datastore {
	if ds.Len() == 1 {
		return nil
	}

	d := ds.List()[0].(map[string]interface{})
	return &apis.Datastore{
		Type:    d["type"].(string),
		Version: d["version"].(string),
	}
}

func buildExtendedProperties(ep *schema.Set) *apis.ExtendedProp {
	if ep.Len() < 1 {
		return nil
	}

	p := ep.List()[0].(map[string]interface{})

	return &apis.ExtendedProp{
		WorkSpaceId: p["workspace"].(string),
		ResourceId:  p["resource"].(string),
		Trial:       p["trial"].(string),
	}
}

func buildSysTags(st []interface{}) []tags.ResourceTag {
	if len(st) < 1 {
		return nil
	}

	ts := make([]tags.ResourceTag, len(st))

	for _, t := range st {
		tag := t.(map[string]interface{})
		ts = append(ts, tags.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		})
	}

	return ts
}

func buildApiCreateOpts(d *schema.ResourceData) (apis.CreateOpts, error) {

	opts := apis.CreateOpts{
		AutoRemind: d.Get("auto_remind").(bool),
		Email:      d.Get("email").(string),
		PhoneNum:   d.Get("phone_number").(string),
		XLang:      d.Get("language").(string),
		Cluster:    buildCluster(d),
	}

	log.Printf("[DEBUG] The Migration Cluster Creation Opts are : %+v", opts)
	return opts, nil
}

func resourceClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DataArtsMigrationsV1Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	opts, err := buildApiCreateOpts(d)
	if err != nil {
		return diag.Errorf("unable to build the OpenTelekomCloud DataArts Migrations Cluster API create opts: %s", err)
	}
	resp, err := apis.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud DataArts Migrations Cluster API: %s", err)
	}
	d.SetId(resp.Id)

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud DataArts Migrations Cluster (%s) to become available", resp.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"100"},
		Target:     []string{"200"},
		Refresh:    WaitForDAMigrationClusterActive(client, resp.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DataArts Migrations Cluster: %s", err)
	}
	d.SetId(resp.Id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceClusterV1Read(clientCtx, d, meta)
}

func resourceClusterV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DataArtsMigrationsV1Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	resp, err := apis.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud DataArts Migrations Cluster")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("public_endpoint", resp.PublicEndpoint),
		d.Set("datastore", flattenDatastore(resp.Datastore)),
		d.Set("instances", flattenInstances(d, resp.Instances)),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("vpc_id", resp.VpcId),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("is_auto_off", resp.IsAutoOff),
		d.Set("public_endpoint_domain_name", resp.PublicEndpointDomainName),
		d.Set("flavor_name", resp.FlavorName),
		d.Set("availability_zone", resp.AzName),
		d.Set("endpoint_domain_name", resp.EndpointDomainName),
		d.Set("is_schedule_boot_off", resp.IsScheduleBootOff),
		d.Set("namespace", resp.Namespace),
		d.Set("eip_id", resp.EipId),
		d.Set("links", flattenLinks(resp.Links)),
		d.Set("db_user", resp.DbUser),
		d.Set("cluster_mode", resp.ClusterMode),
		d.Set("status", analyseClusterStatus(resp.Status)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving  OpenTelekomCloud DataArts Migrations Cluster API fields: %s", err)
	}

	log.Printf("[DEBUG] The Migration Cluster READ states are : %+v", d)

	return nil
}

// resourceClusterV1Update is a function for a backward compatibility with the ForceNew param for the field 'instances'
func resourceClusterV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DataArtsMigrationsV1Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	deleteOpts := apis.DeleteOpts{}

	if lastBackups, ok := d.Get("keep_last_manual_backup").(int); ok {
		deleteOpts.KeepBackup = lastBackups

	}

	if _, err = apis.Delete(client, d.Id(), deleteOpts); err != nil {
		return diag.Errorf("unable to delete the OpenTelekomCloud DataArts Migrations API (%s): %s", d.Id(), err)
	}
	d.SetId("")

	return nil
}

// func resourceClusterV1ImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
// 	config := meta.(*cfg.Config)
// 	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
// 		return config.DataArtsMigrationsV1Client(config.GetProjectName(d))
// 	})
// 	if err != nil {
// 		return []*schema.ResourceData{d}, fmt.Errorf(errCreationV1Client, err)
// 	}
//
// 	mErr := multierror.Append(nil,
// 		d.Set("gateway_id", parts[0]),
// 		d.Set("policy_id", parts[1]),
// 	)
// 	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
// }

func WaitForDAMigrationClusterActive(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := apis.Get(cceClient, clusterId)
		if err != nil {
			return nil, "", fmt.Errorf("error waiting for OpenTelekomCloud DataArts Migrations Cluster to become active: %w", err)
		}

		return resp, resp.Status, nil
	}
}

func flattenInstances(d *schema.ResourceData, reqParams []apis.DetailedInstances) []map[string]interface{} {
	if len(reqParams) == 0 {
		return nil
	}

	var az, flavorID string
	var nics *schema.Set

	// That's a dirty hack, because we can't get these 3 fields from an API and we should reuse them.
	// If you don't set up them, terraform will recalculate a state and try to apply it again with all null fields.
	instances := d.Get("instances").(*schema.Set).List()

	if len(instances) > 0 {
		inst := instances[0].(map[string]interface{})
		az = inst["availability_zone"].(string)
		flavorID = inst["flavor_id"].(string)
		nics = inst["nics"].(*schema.Set)
	}

	result := make([]map[string]interface{}, len(reqParams))
	for i, v := range reqParams {
		param := map[string]interface{}{
			"availability_zone": az,
			"flavor_id":         flavorID,
			"nics":              flattenNics(nics),
			"flavor":            flattenFlavor(v.Flavor),
			"type":              v.Type,
			"volume":            flattenVolume(v.Volume),
			"status":            analyseStatus(v.Status),
			"vm_id":             v.Id,
			"name":              v.Name,
			"role":              v.Role,
			"group":             v.Group,
			"public_ip":         v.PublicIp,
			"manage_ip":         v.ManageIp,
			"traffic_ip":        v.TrafficIp,
			"shard_id":          v.ShardId,
			"manage_fix_ip":     v.ManageFixIp,
			"private_ip":        v.PrivateIp,
			"internal_ip":       v.InternalIp,
		}

		result[i] = param
	}
	return result
}

func flattenFlavor(f apis.Flavor) []map[string]interface{} {

	fs := make([]map[string]interface{}, 1)

	fs[0] = map[string]interface{}{
		"id": f.Id,
	}
	fs[0]["links"] = flattenLinks(f.Links)

	return fs
}

func flattenLinks(links []apis.ClusterLinks) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}

	ls := make([]map[string]interface{}, 0)
	for _, v := range links {
		link := map[string]interface{}{
			"rel":  v.Rel,
			"href": v.Href,
		}
		ls = append(ls, link)
	}

	return ls
}

func flattenVolume(v apis.Volume) []map[string]interface{} {
	fv := make([]map[string]interface{}, 1)
	fv[0] = map[string]interface{}{
		"type": v.Type,
		"size": v.Size,
	}
	return fv
}

func flattenNics(nics *schema.Set) []map[string]interface{} {
	if nics == nil || len(nics.List()) == 0 {
		return nil
	}

	nicSlice := make([]map[string]interface{}, nics.Len())

	for index, nic := range nics.List() {
		n := nic.(map[string]interface{})
		nicSlice[index] = map[string]interface{}{
			"security_group": n["security_group"].(string),
			"net":            n["net"].(string),
		}
	}

	return nicSlice
}

func analyseStatus(s string) string {
	apiType := map[string]string{
		"100": "creating",
		"200": "normal",
		"300": "failed",
		"303": "failed to be created",
		"400": "deleted",
		"800": "frozen",
	}
	if v, ok := apiType[s]; ok {
		return v
	}
	return ""
}

func flattenDatastore(d apis.Datastore) []map[string]interface{} {
	ds := make([]map[string]interface{}, 1)
	ds[0] = map[string]interface{}{
		"type":    d.Type,
		"version": d.Version,
	}

	return ds
}
func analyseClusterStatus(s string) string {
	apiType := map[string]string{
		"100": "creating",
		"200": "normal",
		"300": "failed",
		"303": "failed to be created",
		"400": "deleted",
		"800": "frozen",
		"900": "stopped",
		"910": "stopping",
		"920": "starting",
	}
	if v, ok := apiType[s]; ok {
		return v
	}

	return ""
}
