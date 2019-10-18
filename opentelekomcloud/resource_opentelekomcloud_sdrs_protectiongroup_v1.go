package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sdrs/v1/protectiongroups"
)

func resourceSdrsProtectiongroupV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSdrsProtectiongroupV1Create,
		Read:   resourceSdrsProtectiongroupV1Read,
		Update: resourceSdrsProtectiongroupV1Update,
		Delete: resourceSdrsProtectiongroupV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSdrsProtectiongroupV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sdrsClient, err := config.sdrsV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud SDRS Client: %s", err)
	}

	createOpts := protectiongroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SourceAZ:    d.Get("source_availability_zone").(string),
		TargetAZ:    d.Get("target_availability_zone").(string),
		DomainID:    d.Get("domain_id").(string),
		SourceVpcID: d.Get("source_vpc_id").(string),
		DrType:      d.Get("dr_type").(string),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	n, err := protectiongroups.Create(sdrsClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud SDRS Protectiongroup: %s", err)
	}

	if err := protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := protectiongroups.GetJobEntity(sdrsClient, n.JobID, "server_group_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)
		return resourceSdrsProtectiongroupV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSdrsProtectiongroupV1Create.")
}

func resourceSdrsProtectiongroupV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	sdrsClient, err := config.sdrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud SDRS client: %s", err)
	}
	n, err := protectiongroups.Get(sdrsClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("source_availability_zone", n.SourceAZ)
	d.Set("target_availability_zone", n.TargetAZ)
	d.Set("domain_id", n.DomainID)
	d.Set("source_vpc_id", n.SourceVpcID)
	d.Set("dr_type", n.DrType)

	return nil
}

func resourceSdrsProtectiongroupV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := config.sdrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud SDRS Client: %s", err)
	}
	var updateOpts protectiongroups.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = protectiongroups.Update(sdrsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}
	return resourceSdrsProtectiongroupV1Read(d, meta)
}

func resourceSdrsProtectiongroupV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := config.sdrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud SDRS client: %s", err)
	}

	n, err := protectiongroups.Delete(sdrsClient, d.Id()).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud SDRS Protectiongroup: %s", err)
	}

	if err := protectiongroups.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
