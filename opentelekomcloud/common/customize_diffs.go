package common

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v1/volumetypes"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	argMissingMsg = "schema missing %s argument"
)

var (
	elementListRegex = regexp.MustCompile(`^(.+?)\.\*\.(.+)$`)
)

func checkVolumeTypeAvailable(d cfg.SchemaOrDiff, argName, expectedAZ string, typeAZs map[string][]string) error {
	volumeType := d.Get(argName)
	if volumeType == nil {
		return fmt.Errorf(argMissingMsg, argName)
	}
	resourceVolType := strings.ToLower(volumeType.(string))
	if resourceVolType == "" {
		return nil
	}
	var validAZs []string
	for typeName, azs := range typeAZs {
		if typeName == resourceVolType {
			validAZs = azs
			break
		}
	}
	if len(validAZs) == 0 {
		return fmt.Errorf("volume type `%s` doesn't exist", resourceVolType)
	}
	if !StringInSlice(expectedAZ, validAZs) {
		return fmt.Errorf(
			"volume type `%v` is not supported in AZ `%s`.\nSupported AZs: %v",
			volumeType, expectedAZ, validAZs,
		)
	}
	return nil
}

func ValidateVolumeType(argName string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		expectedAZ := d.Get("availability_zone").(string)
		if expectedAZ == "" || expectedAZ == "random" {
			log.Printf("[DEBUG] No AZ provided, can't define available volume types")
			return nil
		}
		config := meta.(*cfg.Config)
		client, err := config.BlockStorageV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating blockstorage v3 client: %s", err)
		}

		types, err := volumetypes.List(client)
		if err != nil {
			return fmt.Errorf("error retrieving volume types: %s", err)
		}
		typeAZs := make(map[string][]string) // map of type name (lower case) -> az list
		for _, volumeType := range types {
			typeName := strings.ToLower(volumeType.Name)
			typeAZs[typeName] = getZonesFromVolumeType(volumeType)
		}

		if !strings.Contains(argName, ".*") {
			return checkVolumeTypeAvailable(d, argName, expectedAZ, typeAZs)
		}

		reGroups := elementListRegex.FindStringSubmatch(argName)
		countExpr := fmt.Sprintf("%s.#", reGroups[1])
		count := d.Get(countExpr).(int)
		for i := 0; i < count; i++ {
			exactItemExpr := fmt.Sprintf("%s.%d.%s", reGroups[1], i, reGroups[2])
			if err := checkVolumeTypeAvailable(d, exactItemExpr, expectedAZ, typeAZs); err != nil {
				return err
			}
		}
		return nil
	}
}

func getZonesFromVolumeType(t volumetypes.VolumeType) []string {
	if t.ExtraSpecs == nil || len(t.ExtraSpecs) == 0 {
		return []string{}
	}
	zonesStr := t.ExtraSpecs["RESKEY:availability_zones"].(string)
	return strings.Split(zonesStr, ",")
}

func ValidateVPC(argName string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		vpcID := d.Get(argName)
		if vpcID == nil {
			return fmt.Errorf(argMissingMsg, argName)
		}
		if vpcID == "" {
			return nil
		}
		config := meta.(*cfg.Config)
		vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud CCE Client: %s", err)
		}
		if err := vpcs.Get(vpcClient, vpcID.(string)).Err; err != nil {
			return fmt.Errorf("can't find VPC `%s`: %s", vpcID, err)
		}
		return nil
	}
}

func ValidateSubnet(argName string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		subnetId := d.Get(argName)
		if subnetId == nil {
			return fmt.Errorf(argMissingMsg, argName)
		}
		if subnetId == "" {
			return nil
		}
		config := meta.(*cfg.Config)
		subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud CCE Client: %s", err)
		}
		if err := subnets.Get(subnetClient, subnetId.(string)).Err; err != nil {
			return fmt.Errorf("can't find Subnet `%s`: %s", subnetId, err)
		}
		return nil
	}
}

func MultipleCustomizeDiffs(funcs ...schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		mErr := &multierror.Error{}
		for _, fn := range funcs {
			mErr = multierror.Append(mErr, fn(ctx, d, meta))
		}
		return mErr.ErrorOrNil()
	}
}

func ValidateDiskType(v interface{}, path cty.Path) diag.Diagnostics {
	diskType := v.(string)
	if diskType != "SATA" {
		return nil
	}
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       "[DEPRECATION WARNING]",
		Detail:        "Common I/O (SATA) will reach end of life, end of 2025.",
		AttributePath: path,
	}}
}

func ValidateVpnRegion(v interface{}, path cty.Path) diag.Diagnostics {
	region := v.(string)
	if region != "eu-de" {
		return nil
	}
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       "[DEPRECATION WARNING]",
		Detail:        "Classic VPN reach end of life for eu-de region, end of may 2025.",
		AttributePath: path,
	}}
}

func ValidateDcsEngineVersion(v interface{}, path cty.Path) diag.Diagnostics {
	engine := v.(string)
	if engine != "3.0" {
		return nil
	}
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       "[DEPRECATION WARNING]",
		Detail:        "Redis 3.x versions in DCS have reached their End of Sale status on 21st June 2024.",
		AttributePath: path,
	}}
}

func ValidateDmsEngineVersion(v interface{}, path cty.Path) diag.Diagnostics {
	engine := v.(string)
	if !strings.HasPrefix(engine, "1") {
		return nil
	}
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       "[DEPRECATION WARNING]",
		Detail:        "Kafka 1.x versions in DMS has reached their End of Sale status on 21st June 2024.",
		AttributePath: path,
	}}
}

// FlexibleForceNew make the ForceNew of parameters configurable
// this func accepts a list of non-updatable parameters
// when non-updatable parameters are changed
// if ForceNew is enabled, the resource will be recreated
// if ForceNew is not enabled, an error will be raised
// if there is DiffSuppressFunc in the schema, this func need resource schema to make DiffSuppressFunc work
func FlexibleForceNew(keys []string, resourceSchemas ...map[string]*schema.Schema) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
		var resourceSchema map[string]*schema.Schema
		if len(resourceSchemas) > 0 {
			resourceSchema = resourceSchemas[0]
		}

		c := meta.(*cfg.Config)
		var err error
		forceNew := c.GetForceNew(d)
		keysExpand := expandKeys(keys, d)

		if forceNew {
			for _, k := range keysExpand {
				if err := d.ForceNew(k); err != nil {
					log.Printf("[WARN] unable to require attribute replacement of %s: %s", k, err)
				}
			}
		} else {
			for _, k := range keysExpand {
				if d.Id() != "" && d.HasChange(k) {
					oldValue, newValue := d.GetChange(k)
					if cmp.Equal(oldValue, newValue) {
						continue
					}

					schemaAttr, schemaOk := resourceSchema[k]
					if schemaOk && schemaAttr != nil && schemaAttr.DiffSuppressFunc != nil &&
						schemaAttr.DiffSuppressFunc(k, oldValue.(string), newValue.(string), nil) {
						if schemaAttr.Sensitive {
							log.Printf("[DEBUG] ignoring change of %s due to DiffSuppressFunc, %v", k, "(sensitive value)")
						} else {
							log.Printf("[DEBUG] ignoring change of %s due to DiffSuppressFunc, %v -> %v", k, oldValue, newValue)
						}
					} else {
						if schemaOk && schemaAttr != nil && schemaAttr.Sensitive {
							err = multierror.Append(err, fmt.Errorf("%s can't be updated, %v", k, "(sensitive value)"))
						} else {
							err = multierror.Append(err, fmt.Errorf("%s can't be updated, %v -> %v", k, oldValue, newValue))
						}
					}
				}
			}
		}

		return err
	}
}

func expandKeys(keys []string, d *schema.ResourceDiff) []string {
	res := []string{}
	for _, k := range keys {
		if strings.Contains(k, "*") {
			parts := strings.SplitN(k, ".*.", 2)
			l := len(d.Get(parts[0]).([]interface{}))
			i := 0
			var tempKeys []string
			for i < l {
				tempKeys = append(tempKeys, strings.Join([]string{parts[0], parts[1]}, fmt.Sprintf(".%s.", strconv.Itoa(i))))
				i++
			}
			res = append(res, expandKeys(tempKeys, d)...)
		} else {
			res = append(res, k)
		}
	}
	return res
}
