package cfw

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	group "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/dns"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwDomainNameGroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWDomainNameGroupV1Create,
		ReadContext:   resourceCFWDomainNameGroupV1Read,
		UpdateContext: resourceCFWDomainNameGroupV1Update,
		DeleteContext: resourceCFWDomainNameGroupV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("firewall_id", "object_id", "name"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"firewall_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_names": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"domain_set_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ref_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"config_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCFWDomainNameGroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := group.CreateDomainNameGroupOpts{
		FwInstanceID:  d.Get("firewall_id").(string),
		ObjectID:      d.Get("object_id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		DomainNames:   getDomainNames(d),
		DomainSetType: d.Get("domain_set_type").(int),
	}

	domainNameGroup, err := group.CreateDomainNameGroup(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW domain name group from result: %w", err)
	}
	log.Printf("[DEBUG] Create CFW domain name group %s: %#v", domainNameGroup.Id, domainNameGroup)

	d.SetId(domainNameGroup.Id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWDomainNameGroupV1Read(clientCtx, d, meta)
}

func resourceCFWDomainNameGroupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	firewallId := d.Get("firewall_id").(string)
	objectId := d.Get("object_id").(string)
	name := d.Get("name").(string)

	domainNameGroup, err := group.GetDomainNameGroup(client, name, firewallId, objectId)
	if err != nil {
		return fmterr.Errorf("error fetching CFW Domain Name Group: %w", err)
	}
	d.SetId(domainNameGroup.SetID)

	domainNames, err := group.ListDomainNames(client, d.Id(), firewallId)
	if err != nil {
		return fmterr.Errorf("error fetching domain names within CFW Domain Name Group ID '%s': %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved domain name group %s: %#v", d.Id(), domainNameGroup)

	mErr := multierror.Append(nil,
		d.Set("id", d.Id()),
		d.Set("name", domainNameGroup.Name),
		d.Set("description", domainNameGroup.Description),
		d.Set("domain_names", setDomainNames(domainNames)),
		d.Set("domain_set_type", domainNameGroup.DomainSetType),
		d.Set("ref_count", domainNameGroup.RefCount),
		d.Set("config_status", domainNameGroup.ConfigStatus),
		d.Set("rules", setRules(domainNameGroup.Rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWDomainNameGroupV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	if d.HasChange("domain_names") {
		if err := updateDomainNames(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	firewallId := d.Get("firewall_id").(string)
	updateOpts := group.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	updatedDomainNameGroup, err := group.UpdateDomainNameGroup(client, d.Id(), firewallId, updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CFW Domain Name group %s: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated CFW Domain Name group '%s': %#v", updatedDomainNameGroup.Id, updatedDomainNameGroup)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWDomainNameGroupV1Read(clientCtx, d, meta)
}

func resourceCFWDomainNameGroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW Domain Name Group %s", d.Id())

	firewallId := d.Get("firewall_id").(string)
	err = group.DeleteDomainNameGroup(client, d.Id(), firewallId)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW Domain Name Group: %s", err)
	}

	d.SetId("")
	return nil
}

func getDomainNames(d *schema.ResourceData) []group.DomainSetInfoDto {
	var domainNamesList []group.DomainSetInfoDto
	domainNamesListRaw := d.Get("domain_names").([]interface{})
	for _, v := range domainNamesListRaw {
		domainRaw := v.(map[string]interface{})
		domain := group.DomainSetInfoDto{
			DomainName:  domainRaw["domain_name"].(string),
			Description: domainRaw["description"].(string),
		}
		domainNamesList = append(domainNamesList, domain)
	}
	return domainNamesList
}

func setDomainNames(domainNamesInResp []group.DomainInfo) []map[string]interface{} {
	var domainNames []map[string]interface{}
	for _, domainNameInResp := range domainNamesInResp {
		domainName := map[string]interface{}{
			"domain_name":       domainNameInResp.DomainName,
			"description":       domainNameInResp.Description,
			"domain_address_id": domainNameInResp.DomainAddressID,
		}
		domainNames = append(domainNames, domainName)
	}
	return domainNames
}

func setRules(rulesInResp []group.UseRuleVO) []map[string]interface{} {
	var rules []map[string]interface{}
	for _, ruleInResp := range rulesInResp {
		rule := map[string]interface{}{
			"id":   ruleInResp.ID,
			"name": ruleInResp.Name,
		}
		rules = append(rules, rule)
	}
	return rules
}

func updateDomainNames(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	old, new := d.GetChange("domain_names")
	firewallId := d.Get("firewall_id").(string)
	oldList := old.([]interface{})
	newList := new.([]interface{})

	oldDomains := make(map[string]bool)
	for _, v := range oldList {
		domain := v.(map[string]interface{})
		oldDomains[domain["domain_name"].(string)] = true
	}

	newDomains := make(map[string]bool)
	for _, v := range newList {
		domain := v.(map[string]interface{})
		newDomains[domain["domain_name"].(string)] = true
	}

	var domainsToAdd []group.DomainSetInfoDto
	for _, v := range newList {
		domain := v.(map[string]interface{})
		name := domain["domain_name"].(string)
		if !oldDomains[name] {
			domainName := group.DomainSetInfoDto{
				DomainName:  name,
				Description: domain["description"].(string),
			}
			domainsToAdd = append(domainsToAdd, domainName)
		}
	}
	if len(domainsToAdd) > 0 {
		if _, err := group.AddDomainNames(client, d.Id(), group.AddDomainNameListOpts{
			FwInstanceID: firewallId,
			ObjectID:     d.Get("object_id").(string),
			DomainNames:  domainsToAdd,
		}); err != nil {
			return fmt.Errorf("error adding domains \n%#v\n: %w", domainsToAdd, err)
		}
	}

	currentDomains, err := group.ListDomainNames(client, d.Id(), firewallId)
	if err != nil {
		return fmt.Errorf("error listing domains: %w", err)
	}

	var idsToDelete []string
	for _, dom := range currentDomains {
		if !newDomains[dom.DomainName] {
			idsToDelete = append(idsToDelete, dom.DomainAddressID)
		}
	}

	if len(idsToDelete) > 0 {
		if err := group.DeleteDomainNames(client, d.Id(), firewallId, group.DeleteDomainNameListOpts{
			ObjectID:         d.Get("object_id").(string),
			DomainAddressIDs: idsToDelete,
		}); err != nil {
			return fmt.Errorf("error deleting domains: %w", err)
		}
	}
	return nil
}
