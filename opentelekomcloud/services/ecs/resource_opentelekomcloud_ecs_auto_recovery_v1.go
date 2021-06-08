package ecs

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/auto_recovery"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func resourceECSAutoRecoveryV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}, instanceID string) (bool, error) {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return false, fmt.Errorf("error creating OpenTelekomCloud client: %s", err)
	}

	rId := instanceID

	r, err := auto_recovery.Get(client, rId).Extract()
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] Retrieved ECS-AutoRecovery:%#v of instance:%s", rId, r)
	return strconv.ParseBool(r.SupportAutoRecovery)
}

func setAutoRecoveryForInstance(ctx context.Context, d *schema.ResourceData, meta interface{}, instanceID string, ar bool) error {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud client: %s", err)
	}

	rId := instanceID

	updateOpts := auto_recovery.UpdateOpts{SupportAutoRecovery: strconv.FormatBool(ar)}

	timeout := d.Timeout(schema.TimeoutUpdate)

	log.Printf("[DEBUG] Setting ECS-AutoRecovery for instance:%s with options: %#v", rId, updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := auto_recovery.Update(client, rId, updateOpts)
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error setting ECS-AutoRecovery for instance%s: %s", rId, err)
	}
	return nil
}
