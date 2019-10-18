package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/tags"
)

func resourceBMSTagsV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceBMSTagsV2Create,
		Read:   resourceBMSTagsV2Read,
		Delete: resourceBMSTagsV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceTagsV2(d *schema.ResourceData) []string {
	rawTAGS := d.Get("tags").(*schema.Set)
	tags := make([]string, rawTAGS.Len())
	for i, raw := range rawTAGS.List() {
		tags[i] = raw.(string)
	}
	return tags
}

func resourceBMSTagsV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bmsClient, err := config.bmsClient(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
	}

	createOpts := tags.CreateOpts{
		Tag: resourceTagsV2(d),
	}

	_, err = tags.Create(bmsClient, d.Get("server_id").(string), createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Tags: %s", err)
	}
	d.SetId(d.Get("server_id").(string))

	log.Printf("[INFO] Server ID: %s", d.Id())

	return resourceBMSTagsV2Read(d, meta)

}

func resourceBMSTagsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bmsClient, err := config.bmsClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
	}

	n, err := tags.Get(bmsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud tags: %s", err)
	}

	d.Set("tags", n.Tags)
	d.Set("region", GetRegion(d, config))
	d.Set("server_id", d.Id())

	return nil
}

func resourceBMSTagsV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bmsClient, err := config.bmsClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
	}

	err = tags.Delete(bmsClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud tags: %s", err)
	}

	d.SetId("")
	return nil
}
