package ecs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/keypairs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeKeypairV2Create,
		ReadContext:   resourceComputeKeypairV2Read,
		DeleteContext: resourceComputeKeypairV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: common.ImportAsManaged,
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
				Computed: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceComputeKeypairV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	opts := KeyPairCreateOpts{
		keypairs.CreateOpts{
			Name:      d.Get("name").(string),
			PublicKey: d.Get("public_key").(string),
		},
		common.MapValueSpecs(d),
	}

	shared := d.Get("shared").(bool)
	if !shared {
		log.Printf("[DEBUG] Create Options: %#v", opts)
		key, err := keypairs.Create(client, opts).Extract()
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud keypair: %s", err)
		}
		if opts.CreateOpts.PublicKey == "" {
			if err := d.Set("private_key", key.PrivateKey); err != nil {
				return fmterr.Errorf("error saving private key: %s", err)
			}
		}
	} else {
		log.Printf("[DEBUG] Using non-managed key pair, skipping creation")
	}
	d.SetId(opts.Name)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceComputeKeypairV2Read(clientCtx, d, meta)
}

func resourceComputeKeypairV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	kp, err := keypairs.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "keypair")
	}

	mErr := multierror.Append(
		d.Set("name", kp.Name),
		d.Set("public_key", kp.PublicKey),
		d.Set("region", config.GetRegion(d)),
		d.Set("private_key", d.Get("private_key").(string)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceComputeKeypairV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	shared := d.Get("shared").(bool)
	if !shared {
		if err := keypairs.Delete(client, d.Id()).ExtractErr(); err != nil {
			return fmterr.Errorf("error deleting OpenTelekomCloud keypair: %s", err)
		}
	} else {
		log.Printf("[DEBUG] Using non-managed key pair, skipping deletion")
	}

	d.SetId("")
	return nil
}

func useSharedKeypair(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() != "" { // skip if not new resource
		return nil
	}

	if _, ok := d.GetOk("shared"); ok { // skip if shared is already set
		return nil
	}

	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmt.Errorf(errCreateV2Client, err)
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
	return false, nil
}
