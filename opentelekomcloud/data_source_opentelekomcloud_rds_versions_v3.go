package opentelekomcloud

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/datastores"
)

func dataSourceRdsVersionsV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRdsVersionsV3Read,
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MySQL", "PostgreSQL", "SQLServer",
				}, true),
			},

			"versions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceRdsVersionsV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.rdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating RDSv3 client: %s", err)
	}
	name := d.Get("database_name").(string)
	stores, err := getRdsV3VersionList(client, name)

	if err := d.Set("versions", stores); err != nil {
		return fmt.Errorf("error setting version list: %s", err)
	}
	d.SetId(fmt.Sprintf("%s_versions", name))

	return nil
}

type DataStoresSorting struct {
	sort.StringSlice
}

type Version struct {
	Major int
	Minor int
}

func (d DataStoresSorting) Less(i, j int) bool {
	v1 := strings.Split(d.StringSlice[i], ".")
	v2 := strings.Split(d.StringSlice[j], ".")
	v1Maj, err := strconv.Atoi(v1[0])
	if err != nil {
		return d.StringSlice.Less(i, j)
	}
	v2Maj, err := strconv.Atoi(v2[0])
	if err != nil {
		return d.StringSlice.Less(i, j)
	}
	if v1Maj < v2Maj {
		return true
	}
	if v1Maj > v2Maj {
		return false
	}
	v1Min, err := strconv.Atoi(v1[1])
	if err != nil {
		return d.StringSlice.Less(i, j)
	}
	v2Min, err := strconv.Atoi(v2[1])
	if err != nil {
		return d.StringSlice.Less(i, j)
	}
	return v1Min < v2Min
}

func getRdsV3VersionList(client *golangsdk.ServiceClient, dbName string) ([]string, error) {
	pages, err := datastores.List(client, dbName).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error listing RDSv3 versions: %s", err)
	}
	stores, err := datastores.ExtractDataStores(pages)
	if err != nil {
		return nil, fmt.Errorf("error extracting RDSv3 versions: %s", err)
	}
	result := make([]string, len(stores.DataStores))
	for i, store := range stores.DataStores {
		result[i] = store.Name
	}
	resultSorted := DataStoresSorting{StringSlice: result}
	sort.Sort(sort.Reverse(resultSorted))
	return resultSorted.StringSlice, nil
}
