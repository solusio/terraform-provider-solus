package solus

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "id of the location",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "name of the location",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	var (
		res solus.Location
		err error
	)

	rawID, hasRawID := d.GetOk("id")
	rawName, hasRawName := d.GetOk("name")

	// These two properties are mutually exclusive.
	switch {
	case hasRawID:
		id, ok := rawID.(int)
		if !ok {
			return diag.Errorf("id isn't an integer")
		}

		res, err = client.Locations.Get(ctx, id)
		if err != nil {
			return diag.Errorf("failed to get location by id %d: %s", id, err)
		}

	case hasRawName:
		name, ok := rawName.(string)
		if !ok {
			return diag.Errorf("name isn't a string")
		}

		p, err := client.Locations.List(ctx, new(solus.FilterLocations).ByName(name))
		if err != nil {
			return diag.Errorf("failed to get location by name %q: %s", name, err)
		}

		if len(p.Data) == 0 {
			return diag.Errorf("location not found")
		}

		if len(p.Data) > 1 {
			return diag.Errorf("find more than one location")
		}

		res = p.Data[0]
	}

	d.SetId(strconv.Itoa(res.ID))

	return nil
}
