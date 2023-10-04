package vpc

import "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/mutexkv"

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

var defaultDNS = []string{"100.125.4.25", "100.125.129.199"}
