package opentelekomcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/ptrrecords"
)

func resourceDNSPtrRecordV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSPtrRecordV2Create,
		Read:   resourceDNSPtrRecordV2Read,
		Update: resourceDNSPtrRecordV2Update,
		Delete: resourceDNSPtrRecordV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floatingip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(300, 2147483647),
			},
			"tags": tagsSchema(),
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSPtrRecordV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.dnsV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	tagMap := d.Get("tags").(map[string]interface{})
	var tagList []ptrrecords.Tag
	for k, v := range tagMap {
		tag := ptrrecords.Tag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	createOpts := ptrrecords.CreateOpts{
		PtrName:     d.Get("name").(string),
		Description: d.Get("description").(string),
		TTL:         d.Get("ttl").(int),
		Tags:        tagList,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	fipID := d.Get("floatingip_id").(string)
	ptr, err := ptrrecords.Create(client, region, fipID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become available", ptr.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING_CREATE"},
		Refresh:    waitForDNSPtrRecord(client, ptr.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error waiting for PTR record (%s) to become ACTIVE for creation: %s", ptr.ID, err)
	}
	d.SetId(ptr.ID)

	log.Printf("[DEBUG] Created OpenTelekomCloud DNS PTR record %s: %#v", ptr.ID, ptr)
	return resourceDNSPtrRecordV2Read(d, meta)
}

func resourceDNSPtrRecordV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	ptr, err := ptrrecords.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Unable to delete ptr_record")
	}

	log.Printf("[DEBUG] Retrieved PTR record %s: %#v", d.Id(), ptr)

	// Obtain relevant info from parsing the ID
	fipID, err := parseDNSV2PtrRecordID(d.Id())
	if err != nil {
		return err
	}
	mErr := multierror.Append(nil,
		d.Set("name", ptr.PtrName),
		d.Set("description", ptr.Description),
		d.Set("floatingip_id", fipID),
		d.Set("ttl", ptr.TTL),
		d.Set("address", ptr.Address),
	)

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	// save tags
	resourceTags, err := tags.Get(client, "DNS-ptr_record", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error fetching OpenTelekomCloud DNS ptr record tags: %s", err)
	}

	tagMap := tagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmt.Errorf("error saving tags for OpenTelekomCloud DNS ptr record %s: %s", d.Id(), err)
	}

	return nil
}

func resourceDNSPtrRecordV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.dnsV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	tagMap := d.Get("tags").(map[string]interface{})
	var tagList []tags.ResourceTag
	for k, v := range tagMap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	createOpts := ptrrecords.CreateOpts{
		PtrName:     d.Get("name").(string),
		Description: d.Get("description").(string),
		TTL:         d.Get("ttl").(int),
	}

	log.Printf("[DEBUG] Update Options: %#v", createOpts)
	fipID := d.Get("floatingip_id").(string)
	ptr, err := ptrrecords.Create(client, region, fipID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating OpenTelekomCloud DNS PTR record: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		if err := UpdateResourceTags(client, d, "DNS-ptr_record", d.Id()); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become available", ptr.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING_CREATE"},
		Refresh:    waitForDNSPtrRecord(client, ptr.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error waiting for PTR record (%s) to become ACTIVE for update: %s", ptr.ID, err)
	}

	log.Printf("[DEBUG] Updated OpenTelekomCloud DNS PTR record %s: %#v", ptr.ID, ptr)
	return resourceDNSPtrRecordV2Read(d, meta)

}

func resourceDNSPtrRecordV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	err = ptrrecords.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to be deleted", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING_DELETE", "ERROR"},
		Refresh:    waitForDNSPtrRecord(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for PTR record (%s) to become DELETED for deletion: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForDNSPtrRecord(dnsClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ptr, err := ptrrecords.Get(dnsClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return ptr, "DELETED", nil
			}
			return nil, "", err
		}

		return ptr, ptr.Status, nil
	}
}

// PTR record ID, which is in {region}:{floatingip_id} format
func parseDNSV2PtrRecordID(id string) (string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 {
		return "", fmt.Errorf("unable to determine DNS PTR record ID from raw ID: %s", id)
	}

	fipID := idParts[1]
	return fipID, nil
}
