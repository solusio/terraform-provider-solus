package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceOsImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("OS Image", dataSourceOsImageRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the OS Image",
				ValidateFunc: validation.IntAtLeast(1),
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the OS Image",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceOsImageRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	var (
		res solus.OsImage
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	switch {
	case hasID:
		res, err = client.OsImages.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasName:
		res, err = dataSourceOSImageByName(ctx, client, name.(string))
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
}

func dataSourceOSImageByName(ctx context.Context, client *client, name string) (solus.OsImage, error) {
	res, err := client.OsImages.List(ctx, new(solus.FilterOsImages).ByName(name))
	if err != nil {
		return solus.OsImage{}, err
	}

	if len(res.Data) == 1 {
		return res.Data[0], nil
	}

	err = errResourceNotFound
	if len(res.Data) > 1 {
		err = errTooManyResults
	}
	return solus.OsImage{}, err
}
