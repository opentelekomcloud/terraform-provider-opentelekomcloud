package vpc

import "github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

var defaultDNS = []string{"100.125.4.25", "1.1.1.1"}
