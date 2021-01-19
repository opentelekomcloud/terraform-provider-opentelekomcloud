package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func validateRDSv3Version(argumentName string) schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		config, ok := meta.(*Config)
		if !ok {
			return fmt.Errorf("error retreiving configuration: can't convert %v to Config", meta)
		}

		rdsClient, err := config.rdsV3Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud RDSv3 Client: %s", err)
		}

		dataStoreInfo := d.Get(argumentName).([]interface{})[0].(map[string]interface{})
		datastoreVersions, err := getRdsV3VersionList(rdsClient, dataStoreInfo["type"].(string))
		if err != nil {
			return fmt.Errorf("unable to get datastore versions: %s", err)
		}

		var matches = false
		for _, datastore := range datastoreVersions {
			if datastore == dataStoreInfo["version"] {
				matches = true
				break
			}
		}
		if !matches {
			return fmt.Errorf("can't find version `%s`", dataStoreInfo["version"])
		}

		return nil
	}
}
