package dws

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dws/v1/cluster"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsClusterV1Create,
		ReadContext:   resourceDwsClusterV1Read,
		UpdateContext: resourceDwsClusterV1Update,
		DeleteContext: resourceDwsClusterV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(4, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[\-_A-Za-z0-9]+$`),
						"Only letters, digits, underscores (_), and hyphens (-) are allowed.",
					),
				),
			},
			"user_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateName,
			},
			"user_pwd": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"number_of_node": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 256),
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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(8000, 30000),
				ForceNew:     true,
				Computed:     true,
			},
			"number_of_cn": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 20),
				Default:      3,
				ForceNew:     true,
			},
			"keep_last_manual_snapshot": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_bind_type": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"auto_assign", "not_use", "bind_existing",
							}, false),
						},
						"eip_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jdbc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"public_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_connect_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jdbc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"recent_event": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDwsClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DwsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := cluster.CreateClusterOpts{
		NodeType:         d.Get("node_type").(string),
		Name:             d.Get("name").(string),
		NumberOfNode:     d.Get("number_of_node").(int),
		SubnetId:         d.Get("network_id").(string),
		SecurityGroupId:  d.Get("security_group_id").(string),
		VpcId:            d.Get("vpc_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		Port:             d.Get("port").(int),
		UserName:         d.Get("user_name").(string),
		UserPwd:          d.Get("user_pwd").(string),
		NumberOfCn:       d.Get("number_of_cn").(int),
	}

	if _, ok := d.GetOk("public_ip.0"); ok {
		createOpts.PublicIp = cluster.PublicIp{
			PublicBindType: d.Get("public_ip.0.public_bind_type").(string),
			EipId:          d.Get("public_ip.0.eip_id").(string),
		}
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	clusterID, err := cluster.CreateCluster(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating DWS cluster: %w", err)
	}
	log.Printf("[INFO] cluster ID: %s", clusterID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"AVAILABLE"},
		Refresh:    dwsClusterV1StateRefreshFunc(client, clusterID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", clusterID, err)
	}

	d.SetId(clusterID)

	return resourceDwsClusterV1Read(ctx, d, meta)
}

func resourceDwsClusterV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DwsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	v, err := cluster.ListClusterDetails(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] DWS cluster %s: %+v", d.Id(), v)

	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("network_id", v.SubnetId),
		d.Set("node_type", v.NodeType),
		d.Set("number_of_node", v.NumberOfNode),
		d.Set("security_group_id", v.SecurityGroupId),
		d.Set("user_name", v.UserName),
		d.Set("vpc_id", v.VpcId),
		d.Set("availability_zone", v.AvailabilityZone),
		d.Set("port", v.Port),
		d.Set("created", v.Created),
		d.Set("recent_event", v.RecentEvent),
		d.Set("status", v.Status),
		d.Set("sub_status", v.SubStatus),
		d.Set("task_status", v.TaskStatus),
		d.Set("updated", v.Updated),
		d.Set("version", v.Version),
		d.Set("private_ip", v.PrivateIp),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	if v.PublicIp.EipId != "" {
		value := []interface{}{map[string]string{
			"eip_id":           v.PublicIp.EipId,
			"public_bind_type": v.PublicIp.PublicBindType,
		}}
		if err := d.Set("public_ip", value); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(v.Endpoints) > 0 {
		private := make([]interface{}, 0, len(v.Endpoints))
		for _, endpoint := range v.Endpoints {
			transformed := map[string]interface{}{
				"connect_info": endpoint.ConnectInfo,
				"jdbc_url":     endpoint.JdbcUrl,
			}
			private = append(private, transformed)
		}
		if err := d.Set("endpoints", private); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(v.PublicEndpoints) > 0 {
		public := make([]interface{}, 0, len(v.PublicEndpoints))
		for _, endpoint := range v.PublicEndpoints {
			transformed := map[string]interface{}{
				"public_connect_info": endpoint.PublicConnectInfo,
				"jdbc_url":            endpoint.JdbcUrl,
			}
			public = append(public, transformed)
		}
		if err := d.Set("public_endpoints", public); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceDwsClusterV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DwsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	// extend cluster
	if d.HasChange("number_of_node") {
		oldValue, newValue := d.GetChange("number_of_node")
		num := newValue.(int) - oldValue.(int)
		err = cluster.ResizeCluster(client, cluster.ResizeClusterOpts{
			ClusterId: d.Id(),
			Count:     num,
		})
		if err != nil {
			return fmterr.Errorf("Extend DWS cluster failed, cluster_id: %s, error: %s", d.Id(), err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"DONE"},
			Refresh:      dwsClusterV1StateRefreshFuncUpdate(client, d.Id(), true),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 20 * d.Timeout(schema.TimeoutUpdate),
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for cluster (%s) to update: %w", d.Id(), err)
		}
	}

	// change pwd
	if d.HasChange("user_pwd") {
		newValue := d.Get("user_pwd")

		err = cluster.ResetPassword(client, cluster.ResetPasswordOpts{
			ClusterId:   d.Id(),
			NewPassword: newValue.(string),
		})
		if err != nil {
			return fmterr.Errorf("reset password of DWS cluster failed. cluster_id: %s, error: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"DONE"},
			Refresh:      dwsClusterV1StateRefreshFuncUpdate(client, d.Id(), false),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 20 * d.Timeout(schema.TimeoutUpdate),
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for cluster (%s) to update: %w", d.Id(), err)
		}
	}

	return resourceDwsClusterV1Read(ctx, d, meta)
}

func resourceDwsClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DwsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	_, err = cluster.ListClusterDetails(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DWS instance")
	}
	var keepSnapshots = new(int)
	*keepSnapshots = d.Get("keep_last_manual_snapshot").(int)
	err = cluster.DeleteCluster(client, cluster.DeleteClusterOpts{
		ClusterId:              d.Id(),
		KeepLastManualSnapshot: keepSnapshots,
	})
	if err != nil {
		return fmterr.Errorf("error deleting DWS instance: %w", err)
	}

	log.Printf("[DEBUG] Waiting for cluster (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE"},
		Target:     []string{"DELETED"},
		Refresh:    dwsClusterV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for cluster (%s) to delete: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] DWS instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func dwsClusterV1StateRefreshFunc(client *golangsdk.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := cluster.ListClusterDetails(client, clusterID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}

func dwsClusterV1StateRefreshFuncUpdate(client *golangsdk.ServiceClient, clusterID string, isExtendTask bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := cluster.ListClusterDetails(client, clusterID)
		if err != nil {
			return nil, "FAILED", err
		}
		if resp.FailedReasons.ErrorMsg != "" && resp.FailedReasons.ErrorCode != "" {
			return nil, "FAILED", fmt.Errorf("error_code: %s, error_msg: %s", resp.FailedReasons.ErrorCode,
				resp.FailedReasons.ErrorMsg)
		}

		cState, cErr := parseClusterStatus(resp, isExtendTask)
		if cErr != nil {
			return nil, "FAILED", cErr
		}
		if cState {
			return resp, "DONE", nil
		}
		return resp, "PENDING", nil
	}
}

// when extend=true: if TaskStatus = RESIZE_FAILURE ,return error; else just check cluster is no task running
func parseClusterStatus(detail *cluster.ClusterDetail, extend bool) (bool, error) {
	if len(detail.ActionProgress) > 0 {
		return false, nil
	}

	if detail.Status != "AVAILABLE" {
		return false, nil
	}

	if extend && detail.TaskStatus == "RESIZE_FAILURE" {
		return false, fmt.Errorf("RESIZE_FAILURE")
	}

	if detail.TaskStatus != "" {
		return false, nil
	}

	if detail.SubStatus != "NORMAL" {
		return false, nil
	}

	return true, nil
}
