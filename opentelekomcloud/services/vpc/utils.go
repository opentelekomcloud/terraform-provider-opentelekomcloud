package vpc

import "github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()
