package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	pg_ext "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/postgres-extensions"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsPostgresExtensionV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsPostgresExtensionV3Create,
		ReadContext:   resourceRdsPostgresExtensionV3Read,
		DeleteContext: resourceRdsPostgresExtensionV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRdsPostgresExtensionV3ImportState,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"extension_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_preload_libraries": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_install": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceRdsPostgresExtensionV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	instanceId := d.Get("instance_id").(string)
	opts := pg_ext.PostgresExtensionOpts{
		DatabaseName:  d.Get("database_name").(string),
		ExtensionName: d.Get("extension_name").(string),
	}
	if err := pg_ext.Create(client, instanceId, opts); err != nil {
		return fmterr.Errorf("error creating RDS Postgres extension: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, opts.DatabaseName, opts.ExtensionName))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceRdsPostgresExtensionV3Read(clientCtx, d, meta)
}

func getPostgresExtension(client *golangsdk.ServiceClient, instanceId, databaseName, extensionName string) (*pg_ext.Extension, error) {
	offset := 0
	for {
		resp, err := pg_ext.List(client, instanceId, pg_ext.ListOpts{
			DatabaseName: databaseName,
			Offset:       offset,
			Limit:        100,
		})
		if err != nil {
			return nil, err
		}
		for _, extension := range resp.Extensions {
			if extension.Name == extensionName {
				return &extension, nil
			}
		}
		offset += 100
		if offset >= resp.TotalCount {
			return nil, nil
		}
	}
}

func resourceRdsPostgresExtensionV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	instanceId := d.Get("instance_id").(string)
	databaseName := d.Get("database_name").(string)
	extensionName := d.Get("extension_name").(string)

	ext, err := getPostgresExtension(client, instanceId, databaseName, extensionName)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error listing RDS Postgres extensions for instance %s", instanceId))
	}
	if ext == nil || !ext.Created {
		d.SetId("")
		return fmterr.Errorf("Extension not created.")
	}

	mErr := multierror.Append(nil,
		d.Set("version", ext.Version),
		d.Set("version_update", ext.VersionUpdate),
		d.Set("description", ext.Description),
		d.Set("shared_preload_libraries", ext.SharedPreloadLibraries),
		d.Set("enable_install", ext.EnableInstall),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	return nil
}

func resourceRdsPostgresExtensionV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := pg_ext.PostgresExtensionOpts{
		DatabaseName:  d.Get("database_name").(string),
		ExtensionName: d.Get("extension_name").(string),
	}
	if err := pg_ext.Delete(client, d.Get("instance_id").(string), opts); err != nil {
		return fmterr.Errorf("error deleting RDS Postgres extension: %w", err)
	}

	d.SetId("")
	return nil
}

func resourceRdsPostgresExtensionV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<database_name>/<extension_name>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("database_name", parts[1]),
		d.Set("extension_name", parts[2]),
	)
	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}
	return []*schema.ResourceData{d}, nil
}
