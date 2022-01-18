package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceIPBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("IP Block", dataSourceIPBlockRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the IP Block",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the IP Block",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceIPBlockRead(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	var (
		res solus.IPBlock
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	switch {
	case hasID:
		res, err = client.IPBlocks.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasName:
		res, err = dataSourceIPBlockByName(ctx, client, name.(string))
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
}

func dataSourceIPBlockByName(ctx context.Context, client *solus.Client, name string) (solus.IPBlock, error) {
	res, err := client.IPBlocks.List(ctx, new(solus.FilterIPBlocks).ByName(name))
	if err != nil {
		return solus.IPBlock{}, normalizeAPIError(err)
	}

	if len(res.Data) == 1 {
		return res.Data[0], nil
	}

	err = errResourceNotFound
	if len(res.Data) > 1 {
		err = errTooManyResults
	}
	return solus.IPBlock{}, err
}
