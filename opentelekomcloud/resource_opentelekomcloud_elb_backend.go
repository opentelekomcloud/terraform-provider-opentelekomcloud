package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elbaas/backendmember"
)

func resourceBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackendCreate,
		Read:   resourceBackendRead,
		Delete: resourceBackendDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBackendCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	addOpts := backendmember.AddOpts{
		ServerId: d.Get("server_id").(string),
		Address:  d.Get("address").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", addOpts)

	listener_id := d.Get("listener_id").(string)
	job, err := backendmember.Add(client, listener_id, addOpts).ExtractJobResponse()
	if err != nil {
		return err
	}

	log.Printf("Waiting for backend to become active")
	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutCreate)/time.Second)); err != nil {
		return err
	}

	entity, err := golangsdk.GetJobEntity(client, job.URI, "members")

	if members, ok := entity.([]interface{}); ok {
		if len(members) > 0 {
			vmember := members[0]
			if member, ok := vmember.(map[string]interface{}); ok {
				if vid, ok := member["id"]; ok {
					if id, ok := vid.(string); ok {
						d.SetId(id)
						return resourceBackendRead(d, meta)
					}
				}
			}
		}
	}
	return fmt.Errorf("Unexpected conversion error in resourceBackendCreate.")
}

func resourceBackendRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	listener_id := d.Get("listener_id").(string)
	id := d.Id()
	backend, err := backendmember.Get(client, listener_id, id).Extract()
	if err != nil {
		return CheckDeleted(d, err, "backend member")
	}

	log.Printf("[DEBUG] Retrieved backend member %s: %#v", id, backend)

	b := backend[0]
	d.Set("server_id", b.ServerID)
	d.Set("address", b.ServerAddress)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBackendDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.elbV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting backend member %s", d.Id())
	listener_id := d.Get("listener_id").(string)
	id := d.Id()
	job, err := backendmember.Remove(client, listener_id, id).ExtractJobResponse()
	if err != nil {
		return err
	}

	log.Printf("Waiting for backend member %s to delete", id)

	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutDelete)/time.Second)); err != nil {
		return err
	}

	log.Printf("Successfully deleted backend member %s", id)
	return nil
}
