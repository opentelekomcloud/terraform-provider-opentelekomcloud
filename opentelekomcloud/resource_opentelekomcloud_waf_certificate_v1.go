package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/certificates"
)

func resourceWafCertificateV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafCertificateV1Create,
		Read:   resourceWafCertificateV1Read,
		Update: resourceWafCertificateV1Update,
		Delete: resourceWafCertificateV1Delete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWafCertificateV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
		Key:     d.Get("key").(string),
	}

	certificate, err := certificates.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Certificate: %s", err)
	}

	log.Printf("[DEBUG] Waf certificate created: %#v", certificate)
	d.SetId(certificate.Id)

	return resourceWafCertificateV1Read(d, meta)
}

func resourceWafCertificateV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	n, err := certificates.Get(wafClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf Certificate: %s", err)
	}

	d.SetId(n.Id)
	d.Set("name", n.Name)

	return nil
}

func resourceWafCertificateV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts certificates.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	_, err = certificates.Update(wafClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud WAF Certificate: %s", err)
	}
	return resourceWafCertificateV1Read(d, meta)
}

func resourceWafCertificateV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	err = certificates.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF Certificate: %s", err)
	}

	d.SetId("")
	return nil
}
