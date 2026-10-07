package apigw

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/cert"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/domain"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIDomainCertificateAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCertificateAssociateV2Create,
		ReadContext:   resourceDomainCertificateAssociateV2Read,
		DeleteContext: resourceDomainCertificateAssociateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bound": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDomainCertificateAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	groupId := d.Get("group_id").(string)
	domainId := d.Get("domain_id").(string)
	certificateId := d.Get("certificate_id").(string)
	opts := cert.BindOpts{
		InstanceID:     gatewayId,
		GroupID:        groupId,
		DomainID:       domainId,
		CertificateIDs: []string{certificateId},
	}
	err = cert.Bind(client, opts)
	if err != nil {
		return diag.Errorf("error binding OpenTelekomCloud apigw cert to domain: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", gatewayId, groupId, domainId, certificateId))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDomainCertificateAssociateV2Read(clientCtx, d, meta)
}

func resourceDomainCertificateAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId, groupId, domainId, certificateId, err := ParseDomainCertificateAssociateId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	exists := true
	respCert, err := cert.Get(client, certificateId)
	if err != nil {
		exists = false
	}
	respDom, err := domain.GetCertificate(client, domain.CertificateOpts{
		GatewayID:     gatewayId,
		GroupID:       groupId,
		DomainID:      domainId,
		CertificateID: certificateId,
	})
	if err != nil {
		exists = false
	}
	if respCert.CommonName != respDom.CommonName {
		fmt.Printf("Domain in Certificate: %s\nDomain Name in associated cert of Domain:%s", respCert.CommonName, respDom.CommonName)
		exists = false
	}

	if !exists {
		d.SetId("")
	}

	// Set the attributes pulled from the composed resource ID
	mErr := multierror.Append(
		d.Set("gateway_id", gatewayId),
		d.Set("group_id", groupId),
		d.Set("domain_id", domainId),
		d.Set("certificate_id", certificateId),
		d.Set("bound", exists),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDomainCertificateAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId, groupId, domainId, certificateId, err := ParseDomainCertificateAssociateId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	opts := cert.BindOpts{
		InstanceID:     gatewayId,
		GroupID:        groupId,
		DomainID:       domainId,
		CertificateIDs: []string{certificateId},
	}
	err = cert.Unbind(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func ParseDomainCertificateAssociateId(id string) (string, string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 4 {
		return "", "", "", "", fmt.Errorf("unable to determine association ID")
	}

	gatewayId := idParts[0]
	groupId := idParts[1]
	domainId := idParts[2]
	certificateId := idParts[3]

	return gatewayId, groupId, domainId, certificateId, nil
}
