// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package opentelekomcloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func resourceRdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsInstanceV3Create,
		Read:   resourceRdsInstanceV3Read,
		Delete: resourceRdsInstanceV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"db": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backup_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"ha_replication_mode": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"param_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRdsInstanceV3UserInputParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"terraform_resource_data": d,
		"availability_zone":       d.Get("availability_zone"),
		"backup_strategy":         d.Get("backup_strategy"),
		"db":                      d.Get("db"),
		"ha_replication_mode":     d.Get("ha_replication_mode"),
		"name":                    d.Get("name"),
		"network_id":              d.Get("network_id"),
		"param_group_id":          d.Get("param_group_id"),
		"security_group_id":       d.Get("security_group_id"),
		"volume":                  d.Get("volume"),
		"vpc_id":                  d.Get("vpc_id"),
	}
}

func resourceRdsInstanceV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.sdkClient(GetRegion(d, config), "rdsv1", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}
	client.Endpoint = strings.Replace(client.Endpoint, "/rds/v1/", "/v3/", 1)

	opts := resourceRdsInstanceV3UserInputParams(d)
	opts["region"] = GetRegion(d, config)

	arrayIndex := map[string]int{
		"backup_strategy": 0,
		"db":              0,
		"volume":          0,
	}

	params, err := buildRdsInstanceV3CreateParameters(opts, arrayIndex)
	if err != nil {
		return fmt.Errorf("Error building the request body of api(create)")
	}
	r, err := sendRdsInstanceV3CreateRequest(d, params, client)
	if err != nil {
		return fmt.Errorf("Error creating RdsInstanceV3: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	obj, err := asyncWaitRdsInstanceV3Create(d, config, r, client, timeout)
	if err != nil {
		return err
	}
	id, err := navigateValue(obj, []string{"instance", "id"}, nil)
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id.(string))

	return resourceRdsInstanceV3Read(d, meta)
}

func resourceRdsInstanceV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.sdkClient(GetRegion(d, config), "rdsv1", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}
	client.Endpoint = strings.Replace(client.Endpoint, "/rds/v1/", "/v3/", 1)

	res := make(map[string]interface{})

	v, err := fetchRdsInstanceV3ByList(d, client)
	if err != nil {
		return err
	}
	res["list"] = v

	return setRdsInstanceV3Properties(d, res)
}

func resourceRdsInstanceV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.sdkClient(GetRegion(d, config), "rdsv1", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}
	client.Endpoint = strings.Replace(client.Endpoint, "/rds/v1/", "/v3/", 1)

	url, err := replaceVars(d, "instances/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	log.Printf("[DEBUG] Deleting Instance %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONBody:     nil,
		JSONResponse: &r.Body,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	if r.Err != nil {
		return fmt.Errorf("Error deleting Instance %q: %s", d.Id(), r.Err)
	}

	_, err = asyncWaitRdsInstanceV3Delete(d, config, r.Body, client, d.Timeout(schema.TimeoutDelete))
	return err
}

func buildRdsInstanceV3CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	availabilityZoneProp, err := navigateValue(opts, []string{"availability_zone"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := isEmptyValue(reflect.ValueOf(availabilityZoneProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["availability_zone"] = availabilityZoneProp
	}

	backupStrategyProp, err := expandRdsInstanceV3CreateBackupStrategy(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(backupStrategyProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["backup_strategy"] = backupStrategyProp
	}

	configurationIDProp, err := navigateValue(opts, []string{"param_group_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(configurationIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["configuration_id"] = configurationIDProp
	}

	datastoreProp, err := expandRdsInstanceV3CreateDatastore(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(datastoreProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["datastore"] = datastoreProp
	}

	diskEncryptionIDProp, err := navigateValue(opts, []string{"volume", "disk_encryption_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(diskEncryptionIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["disk_encryption_id"] = diskEncryptionIDProp
	}

	flavorRefProp, err := navigateValue(opts, []string{"db", "flavor"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(flavorRefProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["flavor_ref"] = flavorRefProp
	}

	haProp, err := expandRdsInstanceV3CreateHa(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(haProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["ha"] = haProp
	}

	nameProp, err := navigateValue(opts, []string{"name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(nameProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["name"] = nameProp
	}

	passwordProp, err := navigateValue(opts, []string{"db", "password"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(passwordProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["password"] = passwordProp
	}

	portProp, err := navigateValue(opts, []string{"db", "port"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(portProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["port"] = portProp
	}

	regionProp, err := expandRdsInstanceV3CreateRegion(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(regionProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["region"] = regionProp
	}

	securityGroupIDProp, err := navigateValue(opts, []string{"security_group_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(securityGroupIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["security_group_id"] = securityGroupIDProp
	}

	subnetIDProp, err := navigateValue(opts, []string{"network_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(subnetIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["subnet_id"] = subnetIDProp
	}

	volumeProp, err := expandRdsInstanceV3CreateVolume(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(volumeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["volume"] = volumeProp
	}

	vpcIDProp, err := navigateValue(opts, []string{"vpc_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(vpcIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		params["vpc_id"] = vpcIDProp
	}

	return params, nil
}

func expandRdsInstanceV3CreateBackupStrategy(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	keepDaysProp, err := navigateValue(d, []string{"backup_strategy", "keep_days"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := isEmptyValue(reflect.ValueOf(keepDaysProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["keep_days"] = keepDaysProp
	}

	startTimeProp, err := navigateValue(d, []string{"backup_strategy", "start_time"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(startTimeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["start_time"] = startTimeProp
	}

	return req, nil
}

func expandRdsInstanceV3CreateDatastore(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	typeProp, err := navigateValue(d, []string{"db", "type"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := isEmptyValue(reflect.ValueOf(typeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["type"] = typeProp
	}

	versionProp, err := navigateValue(d, []string{"db", "version"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(versionProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["version"] = versionProp
	}

	return req, nil
}

func expandRdsInstanceV3CreateHa(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	modeProp, err := expandRdsInstanceV3CreateHaMode(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := isEmptyValue(reflect.ValueOf(modeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["mode"] = modeProp
	}

	replicationModeProp, err := navigateValue(d, []string{"ha_replication_mode"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(replicationModeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["replication_mode"] = replicationModeProp
	}

	return req, nil
}

func expandRdsInstanceV3CreateHaMode(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"ha_replication_mode"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v1, ok := v.(string); ok && v1 != "" {
		return "ha", nil
	}
	return "", nil
}

func expandRdsInstanceV3CreateVolume(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	sizeProp, err := navigateValue(d, []string{"volume", "size"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err := isEmptyValue(reflect.ValueOf(sizeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["size"] = sizeProp
	}

	typeProp, err := navigateValue(d, []string{"volume", "type"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(typeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		req["type"] = typeProp
	}

	return req, nil
}

func sendRdsInstanceV3CreateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url := client.ServiceURL("instances")

	r := golangsdk.Result{}
	_, r.Err = client.Post(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(create): %s", r.Err)
	}
	return r.Body, nil
}

func asyncWaitRdsInstanceV3Create(d *schema.ResourceData, config *Config, result interface{},
	client *golangsdk.ServiceClient, timeout time.Duration) (interface{}, error) {

	var data = make(map[string]string)
	pathParameters := map[string][]string{
		"id": []string{"job_id"},
	}
	for key, path := range pathParameters {
		value, err := navigateValue(result, path, nil)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving async operation path parameter: %s", err)
		}
		data[key] = value.(string)
	}

	url, err := replaceVars(d, "jobs?id={id}", data)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	return waitToFinish(
		[]string{"Completed"},
		[]string{"Running"},
		timeout, 1*time.Second,
		func() (interface{}, string, error) {
			r := golangsdk.Result{}
			_, r.Err = client.Get(
				url, &r.Body,
				&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
			if r.Err != nil {
				return nil, "", nil
			}

			status, err := navigateValue(r.Body, []string{"status"}, nil)
			if err != nil {
				return nil, "", nil
			}
			return r.Body, status.(string), nil
		},
	)
}

func asyncWaitRdsInstanceV3Delete(d *schema.ResourceData, config *Config, result interface{},
	client *golangsdk.ServiceClient, timeout time.Duration) (interface{}, error) {

	var data = make(map[string]string)
	pathParameters := map[string][]string{
		"id": []string{"job_id"},
	}
	for key, path := range pathParameters {
		value, err := navigateValue(result, path, nil)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving async operation path parameter: %s", err)
		}
		data[key] = value.(string)
	}

	url, err := replaceVars(d, "jobs?id={id}", data)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	return waitToFinish(
		[]string{"Completed"},
		[]string{"Running"},
		timeout, 1*time.Second,
		func() (interface{}, string, error) {
			r := golangsdk.Result{}
			_, r.Err = client.Get(
				url, &r.Body,
				&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
			if r.Err != nil {
				return nil, "", nil
			}

			status, err := navigateValue(r.Body, []string{"status"}, nil)
			if err != nil {
				return nil, "", nil
			}
			return r.Body, status.(string), nil
		},
	)
}

func fetchRdsInstanceV3ByList(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	opts := resourceRdsInstanceV3UserInputParams(d)

	arrayIndex := map[string]int{
		"backup_strategy": 0,
		"db":              0,
		"volume":          0,
	}

	identity := make(map[string]interface{})

	if v, err := navigateValue(opts, []string{"name"}, arrayIndex); err == nil {
		identity["name"] = v
	} else {
		return nil, err
	}

	identity["id"] = d.Id()

	p := make([]string, 0, 2)

	if v, err := convertToStr(identity["name"]); err == nil {
		p = append(p, fmt.Sprintf("name=%v", v))
	} else {
		return nil, err
	}

	if v, err := convertToStr(identity["id"]); err == nil {
		p = append(p, fmt.Sprintf("id=%v", v))
	} else {
		return nil, err
	}
	queryLink := "?" + strings.Join(p, "&")

	link := client.ServiceURL("instances") + queryLink

	return findRdsInstanceV3ByList(client, link, identity)
}

func findRdsInstanceV3ByList(client *golangsdk.ServiceClient, link string, identity map[string]interface{}) (interface{}, error) {
	r, err := sendRdsInstanceV3ListRequest(client, link)
	if err != nil {
		return nil, err
	}

	for _, item := range r.([]interface{}) {
		val := item.(map[string]interface{})

		bingo := true
		for k, v := range identity {
			if val[k] != v {
				bingo = false
				break
			}
		}
		if bingo {
			return item, nil
		}
	}

	return nil, fmt.Errorf("Error finding the resource by list api")
}

func sendRdsInstanceV3ListRequest(client *golangsdk.ServiceClient, url string) (interface{}, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(
		url, &r.Body,
		&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(list) for resource(RdsInstanceV3), error: %s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"instances"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setRdsInstanceV3Properties(d *schema.ResourceData, response map[string]interface{}) error {
	opts := resourceRdsInstanceV3UserInputParams(d)

	backupStrategyProp, _ := opts["backup_strategy"]
	backupStrategyProp, err := flattenRdsInstanceV3BackupStrategy(response, nil, backupStrategyProp)
	if err != nil {
		return fmt.Errorf("Error reading Instance:backup_strategy, err: %s", err)
	}
	if err = d.Set("backup_strategy", backupStrategyProp); err != nil {
		return fmt.Errorf("Error setting Instance:backup_strategy, err: %s", err)
	}

	createdProp, err := navigateValue(response, []string{"list", "created"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:created, err: %s", err)
	}
	if err = d.Set("created", createdProp); err != nil {
		return fmt.Errorf("Error setting Instance:created, err: %s", err)
	}

	dbProp, _ := opts["db"]
	dbProp, err = flattenRdsInstanceV3Db(response, nil, dbProp)
	if err != nil {
		return fmt.Errorf("Error reading Instance:db, err: %s", err)
	}
	if err = d.Set("db", dbProp); err != nil {
		return fmt.Errorf("Error setting Instance:db, err: %s", err)
	}

	haReplicationModeProp, err := navigateValue(response, []string{"list", "ha", "replication_mode"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:ha_replication_mode, err: %s", err)
	}
	if err = d.Set("ha_replication_mode", haReplicationModeProp); err != nil {
		return fmt.Errorf("Error setting Instance:ha_replication_mode, err: %s", err)
	}

	nameProp, err := navigateValue(response, []string{"list", "name"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:name, err: %s", err)
	}
	if err = d.Set("name", nameProp); err != nil {
		return fmt.Errorf("Error setting Instance:name, err: %s", err)
	}

	networkIDProp, err := navigateValue(response, []string{"list", "subnet_id"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:network_id, err: %s", err)
	}
	if err = d.Set("network_id", networkIDProp); err != nil {
		return fmt.Errorf("Error setting Instance:network_id, err: %s", err)
	}

	nodesProp, _ := opts["nodes"]
	nodesProp, err = flattenRdsInstanceV3Nodes(response, nil, nodesProp)
	if err != nil {
		return fmt.Errorf("Error reading Instance:nodes, err: %s", err)
	}
	if err = d.Set("nodes", nodesProp); err != nil {
		return fmt.Errorf("Error setting Instance:nodes, err: %s", err)
	}

	privateIpsProp, err := navigateValue(response, []string{"list", "private_ips"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:private_ips, err: %s", err)
	}
	if err = d.Set("private_ips", privateIpsProp); err != nil {
		return fmt.Errorf("Error setting Instance:private_ips, err: %s", err)
	}

	publicIpsProp, err := navigateValue(response, []string{"list", "public_ips"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:public_ips, err: %s", err)
	}
	if err = d.Set("public_ips", publicIpsProp); err != nil {
		return fmt.Errorf("Error setting Instance:public_ips, err: %s", err)
	}

	securityGroupIDProp, err := navigateValue(response, []string{"list", "security_group_id"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:security_group_id, err: %s", err)
	}
	if err = d.Set("security_group_id", securityGroupIDProp); err != nil {
		return fmt.Errorf("Error setting Instance:security_group_id, err: %s", err)
	}

	volumeProp, _ := opts["volume"]
	volumeProp, err = flattenRdsInstanceV3Volume(response, nil, volumeProp)
	if err != nil {
		return fmt.Errorf("Error reading Instance:volume, err: %s", err)
	}
	if err = d.Set("volume", volumeProp); err != nil {
		return fmt.Errorf("Error setting Instance:volume, err: %s", err)
	}

	vpcIDProp, err := navigateValue(response, []string{"list", "vpc_id"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Instance:vpc_id, err: %s", err)
	}
	if err = d.Set("vpc_id", vpcIDProp); err != nil {
		return fmt.Errorf("Error setting Instance:vpc_id, err: %s", err)
	}

	return nil
}

func flattenRdsInstanceV3BackupStrategy(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		result = make([]interface{}, 1, 1)
	}
	if result[0] == nil {
		result[0] = make(map[string]interface{})
	}
	r := result[0].(map[string]interface{})

	keepDaysProp, err := navigateValue(d, []string{"list", "backup_strategy", "keep_days"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:keep_days, err: %s", err)
	}
	r["keep_days"] = keepDaysProp

	startTimeProp, err := navigateValue(d, []string{"list", "backup_strategy", "start_time"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:start_time, err: %s", err)
	}
	r["start_time"] = startTimeProp

	return result, nil
}

func flattenRdsInstanceV3Db(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		result = make([]interface{}, 1, 1)
	}
	if result[0] == nil {
		result[0] = make(map[string]interface{})
	}
	r := result[0].(map[string]interface{})

	flavorProp, err := navigateValue(d, []string{"list", "flavor_ref"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:flavor, err: %s", err)
	}
	r["flavor"] = flavorProp

	portProp, err := navigateValue(d, []string{"list", "port"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:port, err: %s", err)
	}
	r["port"] = portProp

	typeProp, err := navigateValue(d, []string{"list", "datastore", "type"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:type, err: %s", err)
	}
	r["type"] = typeProp

	userNameProp, err := navigateValue(d, []string{"list", "db_user_name"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:user_name, err: %s", err)
	}
	r["user_name"] = userNameProp

	versionProp, err := navigateValue(d, []string{"list", "datastore", "version"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:version, err: %s", err)
	}
	r["version"] = versionProp

	return result, nil
}

func flattenRdsInstanceV3Nodes(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		v, err := navigateValue(d, []string{"list", "nodes"}, arrayIndex)
		if err != nil {
			return nil, err
		}
		n := len(v.([]interface{}))
		result = make([]interface{}, n, n)
	}

	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	for i := 0; i < len(result); i++ {
		newArrayIndex["list.nodes"] = i
		if result[i] == nil {
			result[i] = make(map[string]interface{})
		}
		r := result[i].(map[string]interface{})

		availabilityZoneProp, err := navigateValue(d, []string{"list", "nodes", "availability_zone"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Instance:availability_zone, err: %s", err)
		}
		r["availability_zone"] = availabilityZoneProp

		idProp, err := navigateValue(d, []string{"list", "nodes", "id"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Instance:id, err: %s", err)
		}
		r["id"] = idProp

		nameProp, err := navigateValue(d, []string{"list", "nodes", "name"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Instance:name, err: %s", err)
		}
		r["name"] = nameProp

		roleProp, err := navigateValue(d, []string{"list", "nodes", "role"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Instance:role, err: %s", err)
		}
		r["role"] = roleProp

		statusProp, err := navigateValue(d, []string{"list", "nodes", "status"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Instance:status, err: %s", err)
		}
		r["status"] = statusProp
	}

	return result, nil
}

func flattenRdsInstanceV3Volume(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		result = make([]interface{}, 1, 1)
	}
	if result[0] == nil {
		result[0] = make(map[string]interface{})
	}
	r := result[0].(map[string]interface{})

	diskEncryptionIDProp, err := navigateValue(d, []string{"list", "disk_encryption_id"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:disk_encryption_id, err: %s", err)
	}
	r["disk_encryption_id"] = diskEncryptionIDProp

	sizeProp, err := navigateValue(d, []string{"list", "volume", "size"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:size, err: %s", err)
	}
	r["size"] = sizeProp

	typeProp, err := navigateValue(d, []string{"list", "volume", "type"}, arrayIndex)
	if err != nil {
		return nil, fmt.Errorf("Error reading Instance:type, err: %s", err)
	}
	r["type"] = typeProp

	return result, nil
}
