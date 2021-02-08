package opentelekomcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/recordsets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/zones"
)

func resourceDNSRecordSetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordSetV2Create,
		Read:   resourceDNSRecordSetV2Read,
		Update: resourceDNSRecordSetV2Update,
		Delete: resourceDNSRecordSetV2Delete,
		Importer: &schema.ResourceImporter{
			State: importAsManaged,
		},

		CustomizeDiff: useSharedRecordSet,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"records": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  300,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),

			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func getRecordSetCreateOpts(d schemaOrDiff) RecordSetCreateOpts {
	recordsRaw := d.Get("records").(*schema.Set).List()
	records := make([]string, len(recordsRaw))
	for i, record := range recordsRaw {
		records[i] = record.(string)
	}

	return RecordSetCreateOpts{
		recordsets.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Records:     records,
			TTL:         d.Get("ttl").(int),
			Type:        d.Get("type").(string),
		},
		MapValueSpecs(d),
	}
}

func resourceDNSRecordSetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)

	shared := d.Get("shared").(bool)
	if shared {
		log.Printf("[DEBUG] Using non-managed DNS record set, skipping creation")
		id, _ := getExistingRecordSetID(d, meta)
		d.SetId(fmt.Sprintf("%s/%s", zoneID, id))
		return resourceDNSRecordSetV2Read(d, meta)
	}

	createOpts := getRecordSetCreateOpts(d)

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	recordSet, err := recordsets.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to become available", recordSet.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordSet.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for record set (%s) to become ACTIVE for creation: %s",
			recordSet.ID, err)
	}

	id := fmt.Sprintf("%s/%s", zoneID, recordSet.ID)
	d.SetId(id)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		resourceType, err := getDNSRecordSetResourceType(dnsClient, zoneID)
		if err != nil {
			return fmt.Errorf("error getting resource type of DNS record set %s: %s", recordSet.ID, err)
		}

		tagList := expandResourceTags(tagRaw)
		if tagErr := tags.Create(dnsClient, resourceType, recordSet.ID, tagList).ExtractErr(); tagErr != nil {
			return fmt.Errorf("error setting tags of DNS record set %s: %s", recordSet.ID, tagErr)
		}
	}

	log.Printf("[DEBUG] Created OpenTelekomCloud DNS record set %s: %#v", recordSet.ID, recordSet)
	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	n, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "record_set")
	}

	log.Printf("[DEBUG] Retrieved  record set %s: %#v", recordsetID, n)

	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("ttl", n.TTL),
		d.Set("type", n.Type),
		d.Set("records", n.Records),
		d.Set("region", GetRegion(d, config)),
		d.Set("zone_id", zoneID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf(
			"error saving records to state for OpenTelekomCloud DNS record set (%s): %s", d.Id(), err)
	}

	// save tags
	resourceType, err := getDNSRecordSetResourceType(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("error getting resource type of DNS record set %s: %s", recordsetID, err)
	}
	resourceTags, err := tags.Get(dnsClient, resourceType, recordsetID).Extract()
	if err != nil {
		return fmt.Errorf("error fetching OpenTelekomCloud DNS record set tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("error saving tags for OpenTelekomCloud DNS record set %s: %s", recordsetID, err)
	}

	return nil
}

func resourceDNSRecordSetV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	var updateOpts recordsets.UpdateOpts
	if d.HasChange("ttl") {
		updateOpts.TTL = d.Get("ttl").(int)
	}

	// `records` is required attribute for update request
	recordsRaw := d.Get("records").(*schema.Set).List()
	records := make([]string, len(recordsRaw))
	for i, recordRaw := range recordsRaw {
		records[i] = recordRaw.(string)
	}
	updateOpts.Records = records

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating  record set %s with options: %#v", recordsetID, updateOpts)

	_, err = recordsets.Update(dnsClient, zoneID, recordsetID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating OpenTelekomCloud DNS  record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to update", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for record set (%s) to become ACTIVE for updation: %s",
			recordsetID, err)
	}

	// update tags
	resourceType, err := getDNSRecordSetResourceType(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("error getting resource type of DNS record set %s: %s", d.Id(), err)
	}

	tagErr := UpdateResourceTags(dnsClient, d, resourceType, recordsetID)
	if tagErr != nil {
		return fmt.Errorf("error updating tags of DNS record set %s: %s", d.Id(), tagErr)
	}

	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
	}

	shared := d.Get("shared").(bool)
	if shared {
		log.Printf("[DEBUG] Using non-managed DNS record set, skipping deletion")
		d.SetId("")
		return nil
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	err = recordsets.Delete(dnsClient, zoneID, recordsetID).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to be deleted", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for record set (%s) to become DELETED for deletion: %s",
			recordsetID, err)
	}

	d.SetId("")
	return nil
}

func waitForDNSRecordSet(dnsClient *golangsdk.ServiceClient, zoneID, recordsetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		recordset, err := recordsets.Get(dnsClient, zoneID, recordsetId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return recordset, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud DNS record set (%s) current status: %s", recordset.ID, recordset.Status)
		return recordset, parseStatus(recordset.Status), nil
	}
}

func parseDNSV2RecordSetID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", fmt.Errorf("unable to determine DNS record set ID from raw ID: %s", id)
	}

	zoneID := idParts[0]
	recordsetID := idParts[1]

	return zoneID, recordsetID, nil
}

// get resource type of DNS record set from zone_id
func getDNSRecordSetResourceType(client *golangsdk.ServiceClient, zoneID string) (string, error) {
	zone, err := zones.Get(client, zoneID).Extract()
	if err != nil {
		return "", err
	}

	zoneType := zone.ZoneType
	if zoneType == "public" {
		return "DNS-public_recordset", nil
	} else if zoneType == "private" {
		return "DNS-private_recordset", nil
	}
	return "", fmt.Errorf("invalid zone type: %s", zoneType)
}

func getExistingRecordSetID(d schemaOrDiff, meta interface{}) (id string, err error) {
	config := meta.(*Config)
	client, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		err = fmt.Errorf("error creating OpenTelekomCloud DNS client: %s", err)
		return
	}

	createOpts := getRecordSetCreateOpts(d)

	zoneID := d.Get("zone_id").(string)
	if zoneID == "" {
		return
	}

	listOpts := recordsets.ListOpts{
		Name:   createOpts.Name,
		TTL:    createOpts.TTL,
		Type:   createOpts.Type,
		ZoneID: zoneID,
	}

	allPages, err := recordsets.ListByZone(client, zoneID, listOpts).AllPages()
	if err != nil {
		return "", fmt.Errorf("error listing record sets: %s", err)
	}
	sets, err := recordsets.ExtractRecordSets(allPages)
	if err != nil {
		return "", fmt.Errorf("error extracting record sets: %s", err)
	}
	if len(sets) != 0 {
		id = sets[0].ID
	}
	return
}

func useSharedRecordSet(d *schema.ResourceDiff, meta interface{}) (err error) {
	if d.Id() != "" { // skip if not new resource
		return
	}

	if _, ok := d.GetOk("shared"); ok { // skip if shared is already set
		return
	}

	id, err := getExistingRecordSetID(d, meta)
	if id == "" {
		_ = d.SetNew("shared", false)
		return
	}

	_ = d.SetNew("shared", true)
	return
}
