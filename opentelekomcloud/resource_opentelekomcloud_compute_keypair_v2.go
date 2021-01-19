package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/keypairs"
)

func resourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeKeypairV2Create,
		Read:   resourceComputeKeypairV2Read,
		Delete: resourceComputeKeypairV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: useSharedKeypair,

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
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceComputeKeypairV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	opts := KeyPairCreateOpts{
		keypairs.CreateOpts{
			Name:      d.Get("name").(string),
			PublicKey: d.Get("public_key").(string),
		},
		MapValueSpecs(d),
	}

	shared := d.Get("shared").(bool)
	if !shared {
		log.Printf("[DEBUG] Create Options: %#v", opts)
		_, err := keypairs.Create(computeClient, opts).Extract()
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud keypair: %s", err)
		}
	} else {
		log.Printf("[DEBUG] Using non-managed key pair, skipping creation")
	}
	d.SetId(opts.Name)

	return resourceComputeKeypairV2Read(d, meta)
}

func resourceComputeKeypairV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	kp, err := keypairs.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "keypair")
	}

	mErr := multierror.Append(
		d.Set("name", kp.Name),
		d.Set("public_key", kp.PublicKey),
		d.Set("region", GetRegion(d, config)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	return nil
}

func resourceComputeKeypairV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	shared := d.Get("shared").(bool)

	if !shared {
		err = keypairs.Delete(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("error deleting OpenTelekomCloud keypair: %s", err)
		}
	} else {
		log.Printf("[DEBUG] Using non-managed key pair, skipping deletion")
	}

	d.SetId("")
	return nil
}

func useSharedKeypair(d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() != "" { // skip if not new resource
		return nil
	}

	config := meta.(*Config)
	client, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}
	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)
	exists, err := keyPairExist(client, name, publicKey)
	if err != nil {
		return err
	}
	_ = d.SetNew("shared", exists)
	return nil
}

// Searches for keypair with the same name and public key
func keyPairExist(client *golangsdk.ServiceClient, name, publicKey string) (exists bool, err error) {
	kp, err := keypairs.Get(client, name).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return false, nil
		}
		return false, err
	}
	if kp.PublicKey == publicKey {
		return true, nil
	}
	return true, fmt.Errorf("key %s already exist with different public key", name)
}
