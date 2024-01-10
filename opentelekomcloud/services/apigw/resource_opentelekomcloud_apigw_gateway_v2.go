package apigw

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/gateway"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceAPIGWv2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGatewayCreate,
		ReadContext:   resourceGatewayRead,
		UpdateContext: resourceGatewayUpdate,
		DeleteContext: resourceGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^([A-Za-z][A-Za-z-_0-9]*)$"),
						"The name can only contain letters, digits, hyphens (-) and underscore (_), and must start "+
							"with a letter."),
					validation.StringLenBetween(3, 64),
				),
			},
			"spec_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BASIC", "PROFESSIONAL", "ENTERPRISE", "PLATINUM",
				}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"The description cannot contain the angle brackets (< and >)."),
					validation.StringLenBetween(0, 255),
				),
			},
			"bandwidth_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 2000),
			},
			"ingress_bandwidth_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 2000),
			},
			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "elb",
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(02|06|10|14|18|22):00:00$`),
					"The start-time format of maintenance window is not 'xx:00:00' or "+
						"the hour is not 02, 06, 10, 14, 18 or 22."),
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_ingress_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_egress_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_egress_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpcep_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMaintainEndTime(maintainStart string) (string, error) {
	result := regexp.MustCompile("^(02|06|10|14|18|22):00:00$").FindStringSubmatch(maintainStart)
	if len(result) < 2 {
		return "", fmt.Errorf("the hour is missing")
	}
	num, err := strconv.Atoi(result[1])
	if err != nil {
		return "", fmt.Errorf("the number (%s) cannot be converted to string", result[1])
	}
	return fmt.Sprintf("%02d:00:00", (num+4)%24), nil
}

func buildInstanceAvailabilityZones(d *schema.ResourceData) ([]string, error) {
	if v, ok := d.GetOk("availability_zones"); ok {
		return common.ExpandToStringSlice(v.([]interface{})), nil
	}
	return nil, fmt.Errorf("The parameter 'availability_zones' must be specified")
}

func buildInstanceCreateOpts(d *schema.ResourceData) (gateway.CreateOpts, error) {
	result := gateway.CreateOpts{
		InstanceName:         d.Get("name").(string),
		SpecID:               d.Get("spec_id").(string),
		VpcID:                d.Get("vpc_id").(string),
		SubnetID:             d.Get("subnet_id").(string),
		SecGroupID:           d.Get("security_group_id").(string),
		Description:          d.Get("description").(string),
		BandwidthSize:        pointerto.Int(d.Get("bandwidth_size").(int)),
		LoadbalancerProvider: d.Get("loadbalancer_provider").(string),
		IngressBandwidthSize: pointerto.Int(d.Get("ingress_bandwidth_size").(int)),
	}

	azList, err := buildInstanceAvailabilityZones(d)
	if err != nil {
		return result, err
	}
	result.AvailableZoneIDs = azList

	if v, ok := d.GetOk("maintain_begin"); ok {
		startTime := v.(string)
		result.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return result, err
		}
		result.MaintainEnd = endTime
	}

	log.Printf("[DEBUG] Create options of the dedicated instance is: %#v", result)
	return result, nil
}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIGW v2 client: %s", err)
	}

	opts, err := buildInstanceCreateOpts(d)
	if err != nil {
		return diag.Errorf("error creating the dedicated instance options: %s", err)
	}
	log.Printf("[DEBUG] The CreateOpts of the dedicated instance is: %#v", opts)

	resp, err := gateway.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating the dedicated instance: %s", err)
	}
	d.SetId(resp.InstanceID)

	instanceId := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating"},
		Target:       []string{"success"},
		Refresh:      InstanceStateCreateRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the dedicated instance (%s) to become running: %s", instanceId, err)
	}

	return resourceGatewayRead(ctx, d, meta)
}

// parseInstanceAvailabilityZones is a method that used to convert the string returned by the API which contains
// brackets ([ and ]) and space into a list of strings (available_zone code) and save to state.
func parseInstanceAvailabilityZones(azStr string) []string {
	codesStr := strings.TrimLeft(azStr, "[")
	codesStr = strings.TrimRight(codesStr, "]")
	codesStr = strings.ReplaceAll(codesStr, " ", "")

	return strings.Split(codesStr, ",")
}

func resourceGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIGW v2 client: %s", err)
	}

	instanceId := d.Id()
	resp, err := gateway.Get(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting instance (%s) details form server", instanceId))
	}
	log.Printf("[DEBUG] Retrieved the dedicated instance (%s): %#v", instanceId, resp)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.InstanceName),
		d.Set("vpc_id", resp.VpcID),
		d.Set("subnet_id", resp.SubnetID),
		d.Set("security_group_id", resp.SecurityGroupID),
		d.Set("description", resp.Description),
		d.Set("bandwidth_size", resp.BandwidthSize),
		// Query doesn't return 'ingress_bandwidth_size'
		d.Set("loadbalancer_provider", resp.LoadbalancerProvider),
		d.Set("availability_zones", parseInstanceAvailabilityZones(resp.AvailableZoneIDs)),
		d.Set("maintain_begin", resp.MaintainBegin),
		d.Set("maintain_end", resp.MaintainEnd),
		d.Set("supported_features", resp.SupportedFeatures),
		d.Set("status", resp.Status),
		d.Set("spec_id", resp.Spec),
		d.Set("project_id", resp.ProjectID),
		d.Set("vpc_ingress_address", resp.IngressIp),
		d.Set("public_egress_address", resp.NatEipAddress),
		d.Set("vpcep_service_name", resp.EndpointService.ServiceName),
	)

	if len(resp.PublicIps) > 0 {
		mErr = multierror.Append(mErr, d.Set("public_egress_address", resp.PublicIps[0].IpAddress))
	}
	if len(resp.NodeIps.Shubao) > 0 {
		mErr = multierror.Append(mErr, d.Set("private_egress_addresses", resp.NodeIps.Shubao))
	}
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving resource fields of the dedicated instance: %s", mErr)
	}

	return nil
}

func buildInstanceUpdateOpts(d *schema.ResourceData) (gateway.UpdateOpts, error) {
	result := gateway.UpdateOpts{}
	if d.HasChange("name") {
		result.InstanceName = d.Get("name").(string)
	}
	if d.HasChange("description") {
		result.Description = d.Get("description").(string)
	}
	if d.HasChange("security_group_id") {
		result.SecGroupID = d.Get("security_group_id").(string)
	}
	if d.HasChange("maintain_begin") {
		startTime := d.Get("maintain_begin").(string)
		result.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return result, err
		}
		result.MaintainEnd = endTime
	}

	log.Printf("[DEBUG] Update options of the dedicated instance is: %#v", result)
	return result, nil
}

func updateApigInstanceEgressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldVal, newVal := d.GetChange("bandwidth_size")
	// Disable the egress access.
	if newVal.(int) == 0 {
		err := gateway.DisableEIP(client, d.Id())
		if err != nil {
			return fmt.Errorf("unable to disable egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
		}
		return nil
	}
	// Enable the egress access.
	if oldVal.(int) == 0 {
		size := d.Get("bandwidth_size").(int)
		opts := gateway.EipOpts{
			BandwidthSize: strconv.Itoa(size),
			ID:            d.Id(),
		}
		err := gateway.EnableEIP(client, opts)
		if err != nil {
			return fmt.Errorf("unable to enable egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
		}
	}
	// Update the egress nat.
	size := d.Get("bandwidth_size").(int)
	opts := gateway.EipOpts{
		BandwidthSize: strconv.Itoa(size),
		ID:            d.Id(),
	}
	err := gateway.UpdateEIP(client, opts)
	if err != nil {
		return fmt.Errorf("unable to update egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	// Update egress access
	if d.HasChange("bandwidth_size") {
		if err = updateApigInstanceEgressAccess(d, client); err != nil {
			return diag.Errorf("update egress access failed: %s", err)
		}
	}

	// Update instance name, maintain window, description, security group ID and vpcep service name.
	updateOpts, err := buildInstanceUpdateOpts(d)
	if err != nil {
		return diag.Errorf("unable to get the update options of the dedicated instance: %s", err)
	}
	if updateOpts != (gateway.UpdateOpts{}) {
		updateOpts.ID = d.Id()
		_, err = gateway.Update(client, updateOpts)
		if err != nil {
			return diag.Errorf("error updating the dedicated instance: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Updating"},
			Target:       []string{"Running"},
			Refresh:      InstanceStateRefreshFunc(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        20 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGatewayRead(ctx, d, meta)
}

func resourceGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	if err = gateway.Delete(client, d.Id()); err != nil {
		return diag.Errorf("error deleting the dedicated instance (%s): %s", d.Id(), err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func InstanceStateCreateRefreshFunc(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := gateway.QueryProgress(client, instanceId)
		if err != nil {
			return resp, "", err
		}

		if common.StrSliceContains([]string{"failed"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status)
		}

		if resp.Status == "success" {
			return resp, resp.Status, nil
		}
		return resp, "creating", nil
	}
}

func InstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := gateway.Get(client, instanceId)
		if err != nil {
			return resp, "", err
		}

		if common.StrSliceContains([]string{"InitingFailed", "RegisterFailed", "InstallFailed",
			"UpdateFailed", "RollbackFailed", "UnRegisterFailed", "RestartFail"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status)
		}

		if resp.Status == "Running" {
			return resp, resp.Status, nil
		}
		return resp, "Updating", nil
	}
}
