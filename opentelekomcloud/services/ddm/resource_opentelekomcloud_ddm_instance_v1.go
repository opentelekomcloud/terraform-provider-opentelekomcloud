package ddm

import (
	"context"
	"fmt"
	"log"
	"unicode"

	// "reflect"
	// "sort"
	// "strconv"
	"regexp"
	// "strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	// "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	ddmv1instances "github.com/opentelekomcloud/gophertelekomcloud/openstack/ddm/v1/instances"
	ddmv2instances "github.com/opentelekomcloud/gophertelekomcloud/openstack/ddm/v2/instances"

	// ddmv3instances "github.com/opentelekomcloud/gophertelekomcloud/openstack/ddm/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDdmInstanceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdmInstanceV1Create,
		ReadContext:   resourceDdmInstanceV1Read,
		UpdateContext: resourceDdmInstanceV1Update,
		DeleteContext: resourceDdmInstanceV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: isValidateName,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"engine_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"param_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"time_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: isValidUTCOffset,
			},
			"username": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				ValidateFunc: isValidUsername,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"purge_rds_on_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
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
		},
	}
}

func resourceDdmInstanceV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DdmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	instanceDetails := ddmv1instances.CreateInstanceDetail{
		Name:              d.Get("name").(string),
		FlavorId:          d.Get("flavor_id").(string),
		NodeNum:           d.Get("node_num").(int),
		EngineId:          d.Get("engine_id").(string),
		AvailableZones:    resourceDDMAvailabilityZones(d),
		VpcId:             d.Get("vpc_id").(string),
		SecurityGroupId:   d.Get("security_group_id").(string),
		SubnetId:          d.Get("subnet_id").(string),
		ParamGroupId:      d.Get("param_group_id").(string),
		TimeZone:          d.Get("time_zone").(string),
		AdminUserName:     d.Get("username").(string),
		AdminUserPassword: d.Get("password").(string),
	}
	createOpts := ddmv1instances.CreateOpts{
		Instance: instanceDetails,
	}

	ddmInstance, err := ddmv1instances.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting instance from result: %w", err)
	}
	log.Printf("[DEBUG] Create instance %s: %#v", ddmInstance.Id, ddmInstance)

	d.SetId(ddmInstance.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATE", "CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceDdmInstanceV1Read(clientCtx, d, meta)
}

func resourceDdmInstanceV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DdmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	instance, err := ddmv1instances.QueryInstanceDetails(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching DDM instance: %w", err)
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", d.Id(), instance)

	mErr := multierror.Append(nil,
		d.Set("region", d.Get("region").(string)),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("username", instance.AdminUserName),
		d.Set("availability_zone", []string{instance.AvailableZone}),
		d.Set("node_num", instance.NodeCount),
		d.Set("access_ip", instance.AccessIp),
		d.Set("access_port", instance.AccessPort),
		d.Set("node_status", instance.NodeStatus),
		d.Set("created_at", instance.Created),
		d.Set("updated_at", instance.Updated),
		d.Set("time_zone", d.Get("time_zone").(string)),
	)

	var nodesList []map[string]interface{}
	for _, nodeObj := range instance.Nodes {
		node := make(map[string]interface{})
		node["ip"] = nodeObj.IP
		node["port"] = nodeObj.Port
		node["status"] = nodeObj.Status
		nodesList = append(nodesList, node)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("nodes", nodesList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDdmInstanceV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	clientV1, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DdmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	clientV2, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DdmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if d.HasChange("name") {
		_, newNameRaw := d.GetChange("name")
		newName := newNameRaw.(string)
		_, err = ddmv1instances.Rename(clientV1, d.Id(), newName)
		if err != nil {
			return fmterr.Errorf("error renaming DDM instance: %w", err)
		}
	}

	if d.HasChange("node_num") {
		err = resourceDDMScaling(clientV2, d)
		if err != nil {
			return fmterr.Errorf("error in DDM instance scaling: %w", err)
		}
	}
	clientCtx := common.CtxWithClient(ctx, clientV1, keyClientV1)
	return resourceDdmInstanceV1Read(clientCtx, d, meta)
}

func resourceDdmInstanceV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.DdmV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting Instance %s", d.Id())

	deleteRdsData := d.Get("purge_rds_on_delete").(bool)
	_, err = ddmv1instances.Delete(client, d.Id(), deleteRdsData)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud RDSv3 instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"RUNNING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to get deleted: %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func isValidateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	// Check length between 4 and 64
	if len(value) > 64 || len(value) < 4 {
		errors = append(errors, fmt.Errorf("%q must contain more than 4 and less than 64 characters", k))
	}

	// Check if contains invalid character
	pattern := `^[\-A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, and hyphens allowed in %q", k))
	}

	// Check if doesn't start with a letter
	if !unicode.IsLetter(rune(value[0])) {
		errors = append(errors, fmt.Errorf("%q must start with a letter", k))
	}

	return
}

func isValidUTCOffset(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	// Regular expression pattern for matching UTC offsets from +12:00 to -12:00
	pattern := `^UTC(?:(\+12:00|\-12:00)|([+-](0[0-9]|1[01]):([0-5][0-9])))$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)
	if !re.MatchString(value) {
		errors = append(errors, fmt.Errorf("only valid utc offsets allowed in %q", k))
	}

	return
}

// Checks if the admin username is valid.
func isValidUsername(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	// Check length between 1 and 32
	if len(value) > 32 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must contain more than 1 and less than 32 characters", k))
	}

	// Check if contains invalid character
	pattern := `^[\_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, and underscores allowed in %q", k))
	}

	// Check if doesn't start with a letter
	if !unicode.IsLetter(rune(value[0])) {
		errors = append(errors, fmt.Errorf("%q must start with a letter", k))
	}

	return
}

func resourceDDMAvailabilityZones(d *schema.ResourceData) []string {
	azRaw := d.Get("availability_zone").([]interface{})
	zones := make([]string, 0)
	for _, v := range azRaw {
		zones = append(zones, v.(string))
	}
	return zones
}

func resourceDDMScaling(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldNodeNumRaw, newNodeNumRaw := d.GetChange("node_num")
	oldNodeNum := oldNodeNumRaw.(int)
	newNodeNum := newNodeNumRaw.(int)
	if oldNodeNum < newNodeNum {
		log.Printf("[DEBUG] Scaling up Instance %s", d.Id())
		scaleOutOpts := ddmv2instances.ScaleOutOpts{
			FlavorId:   d.Get("flavor_id").(string),
			NodeNumber: newNodeNum - oldNodeNum,
		}
		_, err := ddmv2instances.ScaleOut(client, d.Id(), scaleOutOpts)
		if err != nil {
			return fmt.Errorf("error scaling up DDM instance: %w", err)
		}
	} else {
		log.Printf("[DEBUG] Scaling down Instance %s", d.Id())
		if oldNodeNum-newNodeNum < 1 {
			return fmt.Errorf("error scaling down DDM instance: %s\n num_nodes needs to be 1 or greater", d.Id())
		}
		scaleInOpts := ddmv2instances.ScaleInOpts{
			NodeNumber: oldNodeNum - newNodeNum,
		}
		_, err := ddmv2instances.ScaleIn(client, d.Id(), scaleInOpts)
		if err != nil {
			return fmt.Errorf("error scaling up DDM instance: %w", err)
		}
	}
	return nil
}
