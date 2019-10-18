package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/policies"
)

func resourceWafPolicyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafPolicyV1Create,
		Read:   resourceWafPolicyV1Read,
		Update: resourceWafPolicyV1Update,
		Delete: resourceWafPolicyV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"webattack": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"common": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_engine": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_scanner": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_script": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_other": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"webshell": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"cc": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"custom": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"whiteblackip": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"privacy": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"ignore": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"antitamper": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
			},
			"full_detection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getOptions(d *schema.ResourceData) *policies.Options {
	optionsRaw := d.Get("options").([]interface{})
	rawMap := optionsRaw[0].(map[string]interface{})
	webattack := rawMap["webattack"].(bool)
	common := rawMap["common"].(bool)
	crawler := rawMap["crawler"].(bool)
	crawler_engine := rawMap["crawler_engine"].(bool)
	crawler_scanner := rawMap["crawler_scanner"].(bool)
	crawler_script := rawMap["crawler_script"].(bool)
	crawler_other := rawMap["crawler_other"].(bool)
	webshell := rawMap["webshell"].(bool)
	cc := rawMap["cc"].(bool)
	custom := rawMap["custom"].(bool)
	whiteblackip := rawMap["whiteblackip"].(bool)
	privacy := rawMap["privacy"].(bool)
	ignore := rawMap["ignore"].(bool)
	antitamper := rawMap["antitamper"].(bool)

	options := &policies.Options{
		WebAttack:      &webattack,
		Common:         &common,
		Crawler:        &crawler,
		CrawlerEngine:  &crawler_engine,
		CrawlerScanner: &crawler_scanner,
		CrawlerScript:  &crawler_script,
		CrawlerOther:   &crawler_other,
		WebShell:       &webshell,
		Cc:             &cc,
		Custom:         &custom,
		WhiteblackIp:   &whiteblackip,
		Privacy:        &privacy,
		Ignore:         &ignore,
		AntiTamper:     &antitamper,
	}

	log.Printf("[DEBUG] getOptions: %#v", options)
	return options
}

func getAction(d *schema.ResourceData) *policies.Action {
	actionRaw := d.Get("action").([]interface{})
	rawMap := actionRaw[0].(map[string]interface{})

	action := &policies.Action{
		Category: rawMap["category"].(string),
	}

	log.Printf("[DEBUG] getAction: %#v", action)
	return action
}

func resourceWafPolicyV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name: d.Get("name").(string),
	}

	policy, err := policies.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Policy: %s", err)
	}

	log.Printf("[DEBUG] Waf policy created: %#v", policy)
	d.SetId(policy.Id)

	// Update the policy as POST API only supports Name argument
	var updateOpts policies.UpdateOpts
	if hasFilledOpt(d, "action") {
		action := getAction(d)
		if action.Category != "" {
			updateOpts.Action = getAction(d)
		}
	}
	if hasFilledOpt(d, "options") {
		updateOpts.Options = getOptions(d)
	}
	if hasFilledOpt(d, "level") {
		updateOpts.Level = d.Get("level").(int)
	}
	if hasFilledOpt(d, "full_detection") {
		full_detection := d.Get("full_detection").(bool)
		updateOpts.FullDetection = &full_detection
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	if updateOpts != (policies.UpdateOpts{}) {
		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF Policy: %s", err)
		}
	}

	if hasFilledOpt(d, "hosts") {
		var updateHostsOpts policies.UpdateHostsOpts
		v := d.Get("hosts").([]interface{})
		hosts := make([]string, len(v))
		for i, v := range v {
			hosts[i] = v.(string)
		}
		updateHostsOpts.Hosts = hosts

		_, err = policies.UpdateHosts(wafClient, d.Id(), updateHostsOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF Policy Hosts: %s", err)
		}
	}

	return resourceWafPolicyV1Read(d, meta)
}

func resourceWafPolicyV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	n, err := policies.Get(wafClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf Policy: %s", err)
	}

	d.SetId(n.Id)
	d.Set("name", n.Name)
	d.Set("level", n.Level)
	d.Set("full_detection", n.FullDetection)
	d.Set("hosts", n.Hosts)

	action := []map[string]string{
		{
			"category": n.Action.Category,
		},
	}
	d.Set("action", action)

	options := []map[string]interface{}{
		{
			"webattack":       *n.Options.WebAttack,
			"common":          *n.Options.Common,
			"crawler":         *n.Options.Crawler,
			"crawler_engine":  *n.Options.CrawlerEngine,
			"crawler_scanner": *n.Options.CrawlerScanner,
			"crawler_script":  *n.Options.CrawlerScript,
			"crawler_other":   *n.Options.CrawlerOther,
			"webshell":        *n.Options.WebShell,
			"cc":              *n.Options.Cc,
			"custom":          *n.Options.Custom,
			"whiteblackip":    *n.Options.WhiteblackIp,
			"privacy":         *n.Options.Privacy,
			"ignore":          *n.Options.Ignore,
			"antitamper":      *n.Options.AntiTamper,
		},
	}
	d.Set("options", options)
	return nil
}

func resourceWafPolicyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts policies.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("action") {
		action := getAction(d)
		if action.Category != "" {
			updateOpts.Action = getAction(d)
		}
	}
	if d.HasChange("options") {
		updateOpts.Options = getOptions(d)
	}
	if d.HasChange("level") {
		updateOpts.Level = d.Get("level").(int)
	}
	if d.HasChange("full_detection") {
		full_detection := d.Get("full_detection").(bool)
		updateOpts.FullDetection = &full_detection
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	if updateOpts != (policies.UpdateOpts{}) {
		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF Policy: %s", err)
		}
	}

	if d.HasChange("hosts") {
		var updateHostsOpts policies.UpdateHostsOpts
		v := d.Get("hosts").([]interface{})
		hosts := make([]string, len(v))
		for i, v := range v {
			hosts[i] = v.(string)
		}
		updateHostsOpts.Hosts = hosts

		_, err = policies.UpdateHosts(wafClient, d.Id(), updateHostsOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF Policy Hosts: %s", err)
		}
	}
	return resourceWafPolicyV1Read(d, meta)
}

func resourceWafPolicyV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	if hosts, ok := d.GetOk("hosts"); ok {
		log.Printf("[DEBUG] Policies already used by domain: %#v", hosts)
		var updateHostsOpts policies.UpdateHostsOpts
		updateHostsOpts.Hosts = make([]string, 0)

		_, err = policies.UpdateHosts(wafClient, d.Id(), updateHostsOpts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error updating OpenTelekomCloud WAF Policy Hosts: %s", err)
		}
	}
	err = policies.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF Policy: %s", err)
	}

	d.SetId("")
	return nil
}
