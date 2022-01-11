package kms

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/kms/v1/grants"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceKmsGrantV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsGrantV1Create,
		ReadContext:   resourceKmsGrantV1Read,
		DeleteContext: resourceKmsGrantV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"grantee_principal": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operations": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"retiring_principal": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"issuing_principal": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKmsGrantV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	rawOperations := d.Get("operations").(*schema.Set).List()
	operations := make([]string, len(rawOperations))
	for i, operation := range rawOperations {
		operations[i] = operation.(string)
	}

	keyID := d.Get("key_id").(string)
	createOpts := grants.CreateOpts{
		KeyID:             keyID,
		GranteePrincipal:  d.Get("grantee_principal").(string),
		Operations:        operations,
		Name:              d.Get("name").(string),
		RetiringPrincipal: d.Get("retiring_principal").(string),
	}

	createGrant, err := grants.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud KMSv1 Grant: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", keyID, createGrant.GrantID))

	return resourceKmsGrantV1Read(ctx, d, meta)
}

func resourceKmsGrantV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	kmsID, grantID, err := ResourceKMSGrantV1ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := grants.ListOpts{
		KeyID: kmsID,
	}
	grantList, err := grants.List(client, listOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	var createdGrant *grants.Grant
	for _, grant := range grantList.Grants {
		if grant.GrantID == grantID {
			createdGrant = &grant
			break
		}
	}
	if createdGrant == nil {
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(
		d.Set("key_id", createdGrant.KeyID),
		d.Set("grantee_principal", createdGrant.GranteePrincipal),
		d.Set("operations", createdGrant.Operations),
		d.Set("name", createdGrant.Name),
		d.Set("retiring_principal", createdGrant.RetiringPrincipal),
		d.Set("creation_date", createdGrant.CreationDate),
		d.Set("issuing_principal", createdGrant.IssuingPrincipal),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKmsGrantV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	kmsID, grantID, err := ResourceKMSGrantV1ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteOpts := &grants.DeleteOpts{
		KeyID:   kmsID,
		GrantID: grantID,
	}

	if err := grants.Delete(client, deleteOpts).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
