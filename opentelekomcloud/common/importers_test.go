package common

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

var res = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"attr1": {
			Type: schema.TypeString,
		},
		"attr2": {
			Type: schema.TypeString,
		},
	},
}

func TestSetComplexID(t *testing.T) {
	d := res.TestResourceData()

	a1 := tools.RandomString("id-", 5)
	th.AssertNoErr(t, d.Set("attr1", a1))
	a2 := tools.RandomString("id-", 5)
	th.AssertNoErr(t, d.Set("attr2", a2))

	th.AssertNoErr(t, SetComplexID(d, "attr1", "attr2"))
	th.AssertEquals(t, fmt.Sprintf("%s/%s", a1, a2), d.Id())
}

func TestImportByPath(t *testing.T) {
	d := res.TestResourceData()
	fnc := ImportByPath("attr1", "attr2")

	a1 := tools.RandomString("id-", 5)
	a2 := tools.RandomString("id-", 5)
	d.SetId(fmt.Sprintf("%s/%s", a1, a2))

	th.AssertEquals(t, "", d.Get("attr1"))
	th.AssertEquals(t, "", d.Get("attr2"))

	dataSlice, err := fnc(context.TODO(), d, nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, a1, d.Get("attr1"))
	th.AssertEquals(t, a2, d.Get("attr2"))
	th.AssertDeepEquals(t, d, dataSlice[0])
}
