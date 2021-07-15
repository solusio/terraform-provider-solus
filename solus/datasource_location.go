package solus

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
	"strconv"
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
	client := m.(metadata).Client
	timeout := m.(metadata).RequestTimeout

	var (
		l   solus.Location
		err error
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rawID, hasRawID := d.GetOk("id")
	rawName, hasRawName := d.GetOk("name")

	// These two properties are mutually exclusive.
	switch {
	case hasRawID:
		id, ok := rawID.(int)
		if !ok {
			return diagErr("Failed to get location by id", "Location ID isn't an integer")
		}

		l, err = client.Locations.Get(ctx, id)
		if err != nil {
			return diagErr("Failed to get location by id", err.Error())
		}

	case hasRawName:
		name, ok := rawName.(string)
		if !ok {
			return diagErr("Failed to get location by name", "Location name isn't a string")
		}

		p, err := client.Locations.List(ctx, new(solus.FilterLocations).ByName(name))
		if err != nil {
			return diagErr("Failed to get location by name", err.Error())
		}

		if len(p.Data) == 0 {
			return diagErr("Failed to get location by name", "Got zero response")
		}

		if len(p.Data) > 1 {
			return diagErr("Failed to get location by name", "Got more than one result")
		}

		l = p.Data[0]
	}

	d.SetId(strconv.Itoa(l.ID))
	if err := locationToResourceData(l, d); err != nil {
		return diagErr("Failed to map location response to resource", err.Error())
	}

	return nil
}
