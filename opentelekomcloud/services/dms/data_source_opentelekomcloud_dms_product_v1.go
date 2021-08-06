package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/products"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDmsProductV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsProductV1Read,

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partition_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"io_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func getIObyIOtype(d *schema.ResourceData, ios []products.IO) []products.IO {
	ioType := d.Get("io_type").(string)
	storageSpecCode := d.Get("storage_spec_code").(string)

	if ioType != "" || storageSpecCode != "" {
		var getIOs []products.IO
		for _, io := range ios {
			if ioType == io.IOType || storageSpecCode == io.StorageSpecCode {
				getIOs = append(getIOs, io)
			}
		}
		return getIOs
	}

	return ios
}

func dataSourceDmsProductV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error get OpenTelekomCloud dms product client: %s", err)
	}

	instanceEngine := d.Get("engine").(string)
	if instanceEngine != "rabbitmq" && instanceEngine != "kafka" {
		return fmterr.Errorf("the instance_engine value should be 'rabbitmq' or 'kafka', not: %s", instanceEngine)
	}

	v, err := products.Get(DmsV1Client, instanceEngine).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	Products := v.Hourly

	log.Printf("[DEBUG] Dms get products : %+v", Products)
	instanceType := d.Get("instance_type").(string)
	if instanceType != "single" && instanceType != "cluster" {
		return fmterr.Errorf("the instance_type value should be 'single' or 'cluster', not: %s", instanceType)
	}
	var filteredPd []products.Detail
	var FilteredPdInfo []products.ProductInfo
	for _, pd := range Products {
		version := d.Get("version").(string)
		if version != "" && pd.Version != version {
			continue
		}

		for _, value := range pd.Values {
			if value.Name != instanceType {
				continue
			}
			for _, detail := range value.Details {
				vmSpecification := d.Get("vm_specification").(string)
				if vmSpecification != "" && detail.VMSpecification != vmSpecification {
					continue
				}
				bandwidth := d.Get("bandwidth").(string)
				if bandwidth != "" && detail.Bandwidth != bandwidth {
					continue
				}

				if instanceType == "single" || instanceEngine == "kafka" {
					storage := d.Get("storage").(string)
					if storage != "" && detail.Storage != storage {
						continue
					}
					IOs := getIObyIOtype(d, detail.IOs)
					if len(IOs) == 0 {
						continue
					}
					detail.IOs = IOs
				} else {
					for _, pdInfo := range detail.ProductInfos {
						storage := d.Get("storage").(string)
						if storage != "" && pdInfo.Storage != storage {
							continue
						}
						nodeNum := d.Get("node_num").(string)
						if nodeNum != "" && pdInfo.NodeNum != nodeNum {
							continue
						}

						ios := getIObyIOtype(d, pdInfo.IOs)
						if len(ios) == 0 {
							continue
						}
						pdInfo.IOs = ios
						FilteredPdInfo = append(FilteredPdInfo, pdInfo)
					}
					if len(FilteredPdInfo) == 0 {
						continue
					}
					detail.ProductInfos = FilteredPdInfo
				}
				filteredPd = append(filteredPd, detail)
			}
		}
	}

	if len(filteredPd) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your filters and try again.")
	}

	pd := filteredPd[0]
	mErr := multierror.Append(d.Set("vm_specification", pd.VMSpecification))

	if instanceType == "single" || instanceEngine == "kafka" {
		d.SetId(pd.ProductID)
		mErr = multierror.Append(mErr,
			d.Set("storage", pd.Storage),
			d.Set("partition_num", pd.PartitionNum),
			d.Set("bandwidth", pd.Bandwidth),
			d.Set("storage_spec_code", pd.IOs[0].StorageSpecCode),
			d.Set("io_type", pd.IOs[0].IOType),
		)
		log.Printf("[DEBUG] Dms product : %+v", pd)
	} else {
		if len(pd.ProductInfos) < 1 {
			return fmterr.Errorf("your query returned no results. Please change your filters and try again.")
		}
		pdInfo := pd.ProductInfos[0]
		d.SetId(pdInfo.ProductID)
		mErr = multierror.Append(mErr,
			d.Set("storage", pdInfo.Storage),
			d.Set("io_type", pdInfo.IOs[0].IOType),
			d.Set("node_num", pdInfo.NodeNum),
			d.Set("storage_spec_code", pdInfo.IOs[0].StorageSpecCode),
		)
		log.Printf("[DEBUG] Dms product : %+v", pdInfo)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
