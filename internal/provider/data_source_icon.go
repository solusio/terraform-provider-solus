package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceIcon() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("Icon", dataSourceIconRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the Icon",
				ValidateFunc: validation.IntAtLeast(1),
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the Icon",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceIconRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	var (
		res solus.Icon
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	switch {
	case hasID:
		res, err = client.Icons.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasName:
		res, err = dataSourceIconReadByName(ctx, client, name.(string))
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
}

func dataSourceIconReadByName(ctx context.Context, client *client, name string) (solus.Icon, error) {
	res, err := client.Icons.List(ctx, new(solus.FilterIcons).ByName(name))
	if err != nil {
		return solus.Icon{}, err
	}

	if len(res.Data) == 1 {
		return res.Data[0], nil
	}

	err = errResourceNotFound
	if len(res.Data) > 1 {
		err = errTooManyResults
	}
	return solus.Icon{}, err
}
