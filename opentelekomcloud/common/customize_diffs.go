package common

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
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
		client, err := config.BlockStorageV3Client(config.GetRegion(d))
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
