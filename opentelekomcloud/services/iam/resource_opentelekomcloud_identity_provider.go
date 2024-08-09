package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/mappings"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/metadata"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/protocols"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/providers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	protocolSAML = "saml"
	protocolOIDC = "oidc"

	scopeSpilt = " "
)

func ResourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-]{1,64}$`),
					"The maximum length is 64 characters. "+
						"Only letters, digits, underscores (_), and hyphens (-) are allowed"),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{protocolSAML, protocolOIDC}, false),
			},
			"mapping_rules": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: common.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					jsonString, _ := common.NormalizeJsonString(v)
					return jsonString
				},
			},
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"metadata": {
				Type:          schema.TypeString,
				Optional:      true,
				StateFunc:     common.GetHashOrEmpty,
				ConflictsWith: []string{"access_config"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"metadata"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"program", "program_console"}, false),
						},
						"provider_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"signing_key": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := common.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"authorization_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"scopes": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"response_type": {
							Type: schema.TypeString,
							// Computed: true,
							Optional: true,
							Default:  "id_token",
						},
						"response_mode": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"fragment", "form_post"}, false),
						},
					},
				},
			},
			"conversion_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"groups": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"remote": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"login_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	clientAdmin, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	// Create a SAML protocol provider.
	opts := providers.CreateOpts{
		ID:          d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("status").(bool),
	}
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Create identity options %s : %#v", name, opts)
	p, err := providers.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf(providerError, "creating", err)
	}
	d.SetId(p.ID)

	prot := d.Get("protocol").(string)
	// Create protocol and default mapping
	err = createProtocol(client, d)
	if err != nil {
		log.Printf("[ERROR] Error in creating provider protocol: %s,", err)
		return diag.FromErr(err)
	}

	// Import metadata, metadata only worked on saml protocol providers
	if prot == protocolSAML {
		err = importMetadata(config, d, meta)
		if err != nil {
			log.Printf("[ERROR] Error importing matedata into identity provider: %s,", err)
			return diag.FromErr(err)
		}
	} else if ac, ok := d.GetOk("access_config"); ok {
		// Create access config for oidc provider.
		accessConfigArr := ac.([]interface{})
		accessConfig := accessConfigArr[0].(map[string]interface{})

		accessType := accessConfig["access_type"].(string)
		createAccessTypeOpts := providers.CreateOIDCOpts{
			AccessMode: accessType,
			IdpUrl:     accessConfig["provider_url"].(string),
			ClientId:   accessConfig["client_id"].(string),
			SigningKey: accessConfig["signing_key"].(string),
			IdpIp:      p.ID,
		}

		if accessType == "program_console" {
			scopes := common.ExpandToStringSlice(accessConfig["scopes"].([]interface{}))
			createAccessTypeOpts.Scope = strings.Join(scopes, scopeSpilt)
			createAccessTypeOpts.AuthEndpoint = accessConfig["authorization_endpoint"].(string)
			createAccessTypeOpts.ResponseType = accessConfig["response_type"].(string)
			createAccessTypeOpts.ResponseMode = accessConfig["response_mode"].(string)
		}
		log.Printf("[DEBUG] Create access type of provider: %#v", opts)

		_, err = providers.CreateOIDC(clientAdmin, createAccessTypeOpts)
		if err != nil {
			return fmterr.Errorf("Error creating the provider access config: %s", err)
		}
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

// importMetadata import metadata to provider, overwrite if it exists.
func importMetadata(conf *cfg.Config, d *schema.ResourceData, meta interface{}) error {
	metaData := d.Get("metadata").(string)
	if len(metaData) == 0 {
		return nil
	}
	clientExt, err := identityExtClient(d, meta)
	if err != nil {
		return err
	}
	client, err := conf.IdentityV3Client(conf.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	providerID := d.Get("name").(string)
	domainID, err := getDomainID(conf, client)
	if err != nil {
		return fmt.Errorf("error getting the domain id, err=%s", err)
	}
	opts := metadata.ImportOpts{
		DomainID: domainID,
		Metadata: metaData,
	}
	_, err = metadata.Import(clientExt, providerID, protocolSAML, opts).Extract()
	if err != nil {
		return fmt.Errorf("failed to import metadata: %s", err)
	}
	return nil
}

// createProtocol create protocol and default mapping
func createProtocol(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	providerID := d.Get("name").(string)
	var mapping *mappings.Mapping
	var err error
	conversionRuleID := "mapping_" + providerID
	if rulesRaw := d.Get("mapping_rules").(string); rulesRaw == "" {
		// Create default mapping
		defaultConversionRules := getDefaultConversionOpts()
		mapping, err = mappings.Create(client, conversionRuleID, *defaultConversionRules).Extract()
		if err != nil {
			return fmt.Errorf("error in creating default conversion rule: %s", err)
		}
	} else {
		// Create custom mapping
		rulesBytes := []byte(rulesRaw)
		rules := make([]mappings.RuleOpts, 1)
		if err = json.Unmarshal(rulesBytes, &rules); err != nil {
			return err
		}

		createOpts := mappings.CreateOpts{
			Rules: rules,
		}
		mapping, err = mappings.Create(client, conversionRuleID, createOpts).Extract()
		if err != nil {
			return fmt.Errorf(mappingError, "creating", err)
		}
	}

	// Create protocol
	protocolName := d.Get("protocol").(string)
	_, err = protocols.Create(client, providerID, protocolName,
		protocols.CreateOpts{MappingID: mapping.ID},
	).Extract()
	if err != nil {
		// If fails to create protocols, then delete the mapping.
		mErr := multierror.Append(
			nil,
			err,
			mappings.Delete(client, conversionRuleID).ExtractErr(),
		)
		log.Printf("[ERROR] Error creating protocol, and the mapping that has been created. Error: %s", mErr)
		return fmt.Errorf("error creating identity provider protocol: %s", mErr.Error())
	}

	return nil
}

func getDefaultConversionOpts() *mappings.CreateOpts {
	localRules := []mappings.LocalRuleOpts{
		{
			User: &mappings.UserOpts{
				Name: "FederationUser",
			},
		},
	}
	remoteRules := []mappings.RemoteRuleOpts{
		{
			Type: "__NAMEID__",
		},
	}

	opts := mappings.CreateOpts{
		Rules: []mappings.RuleOpts{
			{
				Local:  localRules,
				Remote: remoteRules,
			},
		},
	}
	return &opts
}

func resourceIdentityProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	clientExt, err := identityExtClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	p, err := providers.Get(client, d.Id()).Extract()
	if err != nil {
		log.Printf("[ERROR] Error obtaining identity provider: %s", err)
		return common.CheckDeletedDiag(d, err, "error obtaining identity provider")
	}

	// Query the protocol name from OpenTelekomCloud.
	prot := queryProtocolName(client, d)
	domainID, err := getDomainID(config, client)
	if err != nil {
		return fmterr.Errorf("error getting the domain id, err=%s", err)
	}
	url := generateLoginLink(domainID, p.ID, prot)

	mErr := multierror.Append(err,
		d.Set("name", p.ID),
		d.Set("protocol", prot),
		d.Set("status", p.Enabled),
		d.Set("login_link", url),
		d.Set("description", p.Description),
	)

	// Query and set conversion rules
	conversionRuleID := "mapping_" + d.Id()
	conversions, err := mappings.Get(client, conversionRuleID).Extract()
	if err == nil {
		conversionRules := buildConversionRulesAttr(conversions)
		mErr = multierror.Append(mErr,
			d.Set("conversion_rules", conversionRules),
			d.Set("links", conversions.Links),
		)
	}

	// Query and set metadata of the protocol SAML provider
	if prot == protocolSAML {
		r, err := metadata.Get(clientExt, d.Id(), protocolSAML).Extract()
		if err == nil {
			err = d.Set("metadata", common.GetHashOrEmpty(r.Data))
			mErr = multierror.Append(mErr, err)
		}
	}

	// Query and set access type of the protocol OIDC provider
	if prot == protocolOIDC {
		clientOIDC, err := config.IdentityV30Client()
		if err != nil {
			return fmterr.Errorf(clientCreationFail, err)
		}
		accessType, err := providers.GetOIDC(clientOIDC, d.Id())
		if err == nil {
			scopes := strings.Split(accessType.Scope, scopeSpilt)
			accessTypeConfig := []interface{}{
				map[string]interface{}{
					"access_type":            accessType.AccessMode,
					"provider_url":           accessType.IdpUrl,
					"client_id":              accessType.ClientId,
					"signing_key":            accessType.SigningKey,
					"scopes":                 scopes,
					"response_mode":          accessType.ResponseMode,
					"authorization_endpoint": accessType.AuthEndpoint,
					"response_type":          accessType.ResponseType,
				},
			}

			mErr = multierror.Append(
				mErr,
				d.Set("access_config", accessTypeConfig),
			)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		log.Printf("[ERROR] Error setting identity provider attributes %s: %s", d.Id(), err)
		return fmterr.Errorf("Error setting identity provider attributes: %s", err)
	}
	return nil
}

func buildConversionRulesAttr(conversions *mappings.Mapping) []interface{} {
	conversionRules := make([]interface{}, 0, len(conversions.Rules))
	for _, v := range conversions.Rules {
		localRules := make([]map[string]interface{}, 0, len(v.Local))
		for _, localRule := range v.Local {
			r := map[string]interface{}{}
			if localRule.User != nil {
				r["username"] = localRule.User.Name
			}

			if localRule.Group != nil {
				r["group"] = localRule.Group.Name
			}

			if localRule.Groups != "" {
				r["groups"] = localRule.Groups
			}
			localRules = append(localRules, r)
		}

		remoteRules := make([]map[string]interface{}, 0, len(v.Remote))
		for _, remoteRule := range v.Remote {
			r := map[string]interface{}{
				"attribute": remoteRule.Type,
			}
			if len(remoteRule.NotAnyOf) > 0 {
				r["condition"] = "not_any_of"
				r["value"] = remoteRule.NotAnyOf
			} else if len(remoteRule.AnyOneOf) > 0 {
				r["condition"] = "any_one_of"
				r["value"] = remoteRule.AnyOneOf
			}
			remoteRules = append(remoteRules, r)
		}

		rule := map[string]interface{}{
			"local":  localRules,
			"remote": remoteRules,
		}
		conversionRules = append(conversionRules, rule)
	}
	return conversionRules
}

// generateLoginLink generate login link base on config.domainID.
func generateLoginLink(domainID, id, protocol string) string {
	// The domain name is the same as that of the console
	url := fmt.Sprintf("https://auth.otc.t-systems.com/authui/federation/websso?domain_id=%s&idp=%s&protocol=%s",
		domainID, id, protocol)
	return url
}

func queryProtocolName(client *golangsdk.ServiceClient, d *schema.ResourceData) string {
	allPages, err := protocols.List(client, d.Id()).AllPages()
	if err != nil {
		return err.Error()
	}
	allProtocols, err := protocols.ExtractProtocols(allPages)
	if err != nil {
		return fmt.Sprintf("unable to retrieve projects: %s", err.Error())
	}
	// The SAML protocol provider may not have protocol data,
	// so the default value is set to SAML.
	protocolName := protocolSAML
	if len(allProtocols) > 0 {
		protocolName = allProtocols[0].ID
	}
	return protocolName
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	mErr := &multierror.Error{}
	if d.HasChanges("status", "description") {
		status := d.Get("status").(bool)
		description := d.Get("description").(string)
		opts := providers.UpdateOpts{
			Enabled:     &status,
			Description: description,
		}
		log.Printf("[DEBUG] Update identity options %s : %#v", d.Id(), opts)

		_, err = providers.Update(client, d.Id(), opts).Extract()
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("metadata") {
		err = importMetadata(config, d, meta)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("access_config") && d.Get("protocol") == protocolOIDC {
		err = updateAccessConfig(d, meta)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("mapping_rules") {
		err = updateMappingRules(d, meta)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error in updating provider: %s", err)
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func updateAccessConfig(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	clientAdmin, err := config.IdentityV30Client()
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}
	accessConfigArr := d.Get("access_config").([]interface{})
	if len(accessConfigArr) == 0 {
		return fmt.Errorf("the access_config is required for the OIDC provider")
	}
	accessConfig := accessConfigArr[0].(map[string]interface{})

	accessType := accessConfig["access_type"].(string)
	opts := providers.UpdateOIDCOpts{
		AccessMode: accessType,
		IdpUrl:     accessConfig["provider_url"].(string),
		ClientId:   accessConfig["client_id"].(string),
		SigningKey: accessConfig["signing_key"].(string),
	}

	if accessType == "program_console" {
		scopes := common.ExpandToStringSlice(accessConfig["scopes"].([]interface{}))
		opts.Scope = strings.Join(scopes, scopeSpilt)
		opts.AuthEndpoint = accessConfig["authorization_endpoint"].(string)
		opts.ResponseType = accessConfig["response_type"].(string)
		opts.ResponseMode = accessConfig["response_mode"].(string)
	}
	log.Printf("[DEBUG] Update access type of provider: %#v", opts)

	subId := strings.Split(d.Id(), "/")

	switch len(subId) {
	case 1:
		opts.IdpIp = d.Id()
	case 2:
		opts.IdpIp = subId[0]
	}
	_, err = providers.UpdateOIDC(clientAdmin, opts)
	return err
}

func updateMappingRules(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}
	updateOpts := mappings.UpdateOpts{}
	rulesRaw := d.Get("mapping_rules").(string)
	rulesBytes := []byte(rulesRaw)
	rules := make([]mappings.RuleOpts, 1)
	if err = json.Unmarshal(rulesBytes, &rules); err != nil {
		return err
	}
	updateOpts.Rules = rules
	_, err = mappings.Update(client, "mapping_"+d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf(mappingError, "updating", err)
	}
	return nil
}

func resourceIdentityProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	err = providers.Delete(client, d.Id()).Err
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting OpenTelekomCloud identity provider")
	}
	d.SetId("")
	return nil
}
