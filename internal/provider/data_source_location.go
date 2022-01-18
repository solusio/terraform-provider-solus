package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("Location", dataSourceLocationRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the Location",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the Location",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceLocationRead(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	var (
		res solus.Location
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	switch {
	case hasID:
		res, err = client.Locations.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasName:
		res, err = dataSourceLocationByName(ctx, client, name.(string))
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
}

func dataSourceLocationByName(ctx context.Context, client *solus.Client, name string) (solus.Location, error) {
	res, err := client.Locations.List(ctx, new(solus.FilterLocations).ByName(name))
	if err != nil {
		return solus.Location{}, err
	}

	if len(res.Data) == 1 {
		return res.Data[0], nil
	}

	err = errResourceNotFound
	if len(res.Data) > 1 {
		err = errTooManyResults
	}
	return solus.Location{}, err
}
