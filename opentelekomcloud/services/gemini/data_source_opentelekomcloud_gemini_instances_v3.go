package gemini

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceGeminiDBInstancesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBInstancesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"datastore": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage_engine": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"backup_strategy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"keep_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
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
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"support_reduce": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGeminiDBInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	var mErr *multierror.Error

	listOpts := instance.ListGeminiDBOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	allInstances, err := instance.ListGeminiDB(client, listOpts)
	if err != nil {
		return diag.Errorf("error getting GaussDB Cassandra Instances list: %s", err)
	}

	var instancesToSet []map[string]interface{}
	var instancesIds []string

	for _, instanceInAll := range allInstances {
		instanceToSet := map[string]interface{}{
			"id":                    instanceInAll.Id,
			"name":                  instanceInAll.Name,
			"status":                instanceInAll.Status,
			"vpc_id":                instanceInAll.VpcId,
			"subnet_id":             instanceInAll.SubnetId,
			"security_group_id":     instanceInAll.SecurityGroupId,
			"enterprise_project_id": instanceInAll.EnterpriseProjectId,
			"mode":                  instanceInAll.Mode,
			"db_user_name":          instanceInAll.DbUserName,
		}

		if dbPort, err := strconv.Atoi(instanceInAll.Port); err == nil {
			instanceToSet["port"] = dbPort
		}

		// set data store
		dbList := make([]map[string]interface{}, 0, 1)
		db := map[string]interface{}{
			"engine":         instanceInAll.DataStore.Type,
			"version":        instanceInAll.DataStore.Version,
			"storage_engine": instanceInAll.Engine,
		}
		dbList = append(dbList, db)
		instanceToSet["datastore"] = dbList

		specCode := ""
		wrongFlavor := "Inconsistent Flavor"
		ipsList := []string{}
		azList := []string{}
		nodesList := make([]map[string]interface{}, 0, 1)
		for _, group := range instanceInAll.Groups {
			for _, Node := range group.Nodes {
				node := map[string]interface{}{
					"id":                Node.Id,
					"name":              Node.Name,
					"status":            Node.Status,
					"private_ip":        Node.PrivateIp,
					"support_reduce":    Node.SupportReduce,
					"availability_zone": Node.AvailabilityZone,
				}
				if specCode == "" {
					specCode = Node.SpecCode
				} else if specCode != Node.SpecCode && specCode != wrongFlavor {
					specCode = wrongFlavor
				}
				nodesList = append(nodesList, node)
				azList = append(azList, Node.AvailabilityZone)
				// Only return Node private ips which doesn't support reduce
				if !Node.SupportReduce {
					ipsList = append(ipsList, Node.PrivateIp)
				}
			}
			if volSize, err := strconv.Atoi(group.Volume.Size); err == nil {
				instanceToSet["volume_size"] = volSize
			}
			if specCode != "" {
				instanceToSet["flavor"] = specCode
				instanceToSet["datastore"] = dbList
			}
		}
		instanceToSet["nodes"] = nodesList
		instanceToSet["private_ips"] = ipsList

		instanceID := instanceInAll.Id
		instancesIds = append(instancesIds, instanceID)

		instanceToSet["availability_zone"] = strings.Join(azList, ",")
		instanceToSet["node_num"] = len(nodesList)

		// set backup_strategy
		backupStrategyList := make([]map[string]interface{}, 0, 1)
		backupStrategy := map[string]interface{}{
			"start_time": instanceInAll.BackupStrategy.StartTime,
			"keep_days":  instanceInAll.BackupStrategy.KeepDays,
		}
		backupStrategyList = append(backupStrategyList, backupStrategy)
		instanceToSet["backup_strategy"] = backupStrategyList

		// save geminidb tags
		if resourceTags, err := tags.Get(client, "instances", instanceID).Extract(); err == nil {
			tagmap := common.TagsToMap(resourceTags)
			instanceToSet["tags"] = tagmap
		} else {
			log.Printf("[WARN] error fetching tags of geminidb (%s): %s", instanceID, err)
		}

		instancesToSet = append(instancesToSet, instanceToSet)
	}

	d.SetId(hashcode.Strings(instancesIds))
	mErr = multierror.Append(mErr,
		d.Set("instances", instancesToSet),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
