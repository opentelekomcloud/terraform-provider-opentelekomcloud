package gemini

import (
	"context"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/spec"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceGeminiDatastoresV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDatastoresRead,

		Schema: map[string]*schema.Schema{
			"engine_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cassandra", "influxdb",
				}, false),
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"storage_engines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"modes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGeminiDatastoresRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	engineName := d.Get("engine_name").(string)

	versions, err := spec.GetVersions(client, engineName)
	if err != nil {
		return diag.Errorf("error getting GeminiDB %s versions: %s", engineName, err)
	}

	dbVersions := versions.Versions
	if len(dbVersions) == 0 {
		// The version API doesn't report the GeminiDB Influx versions, so fall back to the
		// versions advertised by the specification API.
		dbVersions, err = versionsFromFlavors(config, d, engineName)
		if err != nil {
			return diag.Errorf("error getting GeminiDB %s versions from the flavors: %s", engineName, err)
		}
	}

	d.SetId(hashcode.Strings(append([]string{engineName}, dbVersions...)))

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("versions", dbVersions),
		// The storage engine and the instance types are not exposed by the API, they are
		// fixed per database engine.
		d.Set("storage_engines", []string{"rocksDB"}),
		d.Set("modes", datastoreModes(engineName)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func versionsFromFlavors(config *cfg.Config, d *schema.ResourceData, engineName string) ([]string, error) {
	client, err := config.GeminiDBV31Client(config.GetRegion(d))
	if err != nil {
		return nil, err
	}

	flavors, err := listAllFlavors(client, spec.ListFlavorsOpts{
		EngineName: engineName,
		Mode:       defaultFlavorMode(engineName),
	})
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, len(flavors))
	versions := make([]string, 0, len(flavors))
	for _, flavor := range flavors {
		if flavor.EngineVersion == "" {
			continue
		}
		if _, ok := seen[flavor.EngineVersion]; ok {
			continue
		}
		seen[flavor.EngineVersion] = struct{}{}
		versions = append(versions, flavor.EngineVersion)
	}
	sort.Strings(versions)

	return versions, nil
}

func datastoreModes(engineName string) []string {
	if engineName == "influxdb" {
		return []string{"CloudNativeCluster"}
	}
	return []string{"Cluster"}
}
