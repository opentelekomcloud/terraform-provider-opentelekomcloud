package mrs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/mrs/v1/job"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceMRSJobV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMRSJobV1Create,
		ReadContext:   resourceMRSJobV1Read,
		DeleteContext: resourceMRSJobV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"jar_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arguments": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"input": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_log": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hive_script_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_protected": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"job_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func JobStateRefreshFunc(client *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobGet, err := job.Get(client, jobID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return jobGet, "DELETED", nil
			}
			return nil, "", err
		}
		log.Printf("[DEBUG] JobStateRefreshFunc: %#v", jobGet)
		return jobGet, stateValue(jobGet.JobState), nil
	}
}

func resourceMRSJobV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud MRS client: %s", err)
	}

	createOpts := &job.CreateOpts{
		JobType:        d.Get("job_type").(int),
		JobName:        d.Get("job_name").(string),
		ClusterId:      d.Get("cluster_id").(string),
		JarPath:        d.Get("jar_path").(string),
		Arguments:      d.Get("arguments").(string),
		Input:          d.Get("input").(string),
		Output:         d.Get("output").(string),
		JobLog:         d.Get("job_log").(string),
		HiveScriptPath: d.Get("hive_script_path").(string),
		IsProtected:    d.Get("is_protected").(bool),
		IsPublic:       d.Get("is_public").(bool),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	jobCreate, err := job.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Job: %s", err)
	}

	d.SetId(jobCreate.Id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Starting", "Running"},
		Target:     []string{"Completed"},
		Refresh:    JobStateRefreshFunc(client, jobCreate.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"error waiting for job (%s) to become ready: %s ",
			jobCreate.Id, err)
	}

	return resourceMRSJobV1Read(ctx, d, meta)
}

func resourceMRSJobV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud  MRS client: %s", err)
	}

	jobGet, err := job.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "MRS Job")
	}
	log.Printf("[DEBUG] Retrieved MRS Job %s: %#v", d.Id(), jobGet)

	d.SetId(jobGet.Id)
	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("job_type", jobGet.JobType),
		d.Set("job_name", jobGet.JobName),
		d.Set("cluster_id", jobGet.ClusterId),
		d.Set("jar_path", jobGet.JarPath),
		d.Set("arguments", jobGet.Arguments),
		d.Set("input", jobGet.Input),
		d.Set("output", jobGet.Output),
		d.Set("job_log", jobGet.JobLog),
		d.Set("hive_script_path", jobGet.HiveScriptPath),
		d.Set("job_state", stateValue(jobGet.JobState)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func stateValue(state int) string {
	switch state {
	case -1:
		return "Terminated"
	case 1:
		return "Starting"
	case 2:
		return "Running"
	case 3:
		return "Completed"
	case 4:
		return "Abnormal"
	case 5:
		return "Error"
	default:
		return "Starting"
	}
}

func resourceMRSJobV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud client: %s", err)
	}

	rId := d.Id()
	log.Printf("[DEBUG] Deleting MRS Job %s", rId)

	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := job.Delete(client, rId)
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if common.IsResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable MRS Job: %s", rId)
			return nil
		}
		return fmterr.Errorf("error deleting MRS Job %s: %s", rId, err)
	}
	return nil
}
